package logic

import (
	"fmt"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dbmodel"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/config"
	"hotel-engine/infrastructure/hotelproviderinterface/dtos"
	"hotel-engine/infrastructure/logger"
	"hotel-engine/utils/array"
	"hotel-engine/utils/atomicflag"
	"hotel-engine/utils/date"
	"strconv"
	"strings"
	"time"
)

var (
	inSyncing  = atomicflag.NewAtomicFlag()
	inUpdating = atomicflag.NewAtomicFlag()
)

type hotelService struct {
	mapper               core.Mapper
	provider             core.HotelProvider
	unitOfWork           core.UnitOfWork
	searchDtoAdopter     core.SearchDtoAdopter
	publicService        core.PublicService
	cacheStore           core.CacheStore
	balanceChecker       core.ProviderBalanceChecker
	syncChunkSize        int
	orderEventDispatcher core.OrderEventDispatcher
}

func (g *hotelService) FindHotelById(id string) (*dto.HotelDto, error) {
	hotel, err := g.unitOfWork.Hotel().FindByID(id)
	if err != nil {
		return nil, err
	}
	return g.mapper.ToHotelDto(*hotel), nil
}

func (g *hotelService) GetHotelDetails(request dto.HotelDetailsDto) (*dto.HotelPDPDto, error) {
	hotel, err := g.FindHotelById(request.HotelId)
	if err != nil {
		return nil, err
	}
	hotel.UnavailableAmenities = g.getUnavailableAmenities(hotel.Amenities)
	if !strings.Contains(hotel.Name, "هتل") {
		hotel.Name = fmt.Sprintf("هتل %s", hotel.Name)
	}
	if request.CanSkipOptions() && hotel.Price != 0 {
		return g.mapper.ToHotelPDPDto(*hotel, request.Rooms), nil
	}
	hotelRooms, err := g.GetHotelRooms(dto.HotelRoomsDto{
		HotelId:  hotel.PlaceID,
		CheckIn:  request.CheckIn,
		CheckOut: request.CheckOut,
		Rooms:    request.Rooms,
	})
	if err != nil {
		return nil, err
	}
	hotel.Rooms = hotelRooms.Rooms
	hotel.SessionId = hotelRooms.SessionId
	hotel.Price = getRoomsMinimumPrice(hotelRooms.Rooms)
	return g.mapper.ToHotelPDPDto(*hotel, request.Rooms), nil
}

func getRoomsMinimumPrice(rooms []dto.RoomOptionDto) int64 {
	if len(rooms) == 0 {
		return 0
	}
	minPrice := common.MaxInt64
	for _, room := range rooms {
		if len(room.Rooms) == 0 {
			continue
		}
		var roomMinPricePerNight int64
		roomMinPricePerNight = 0
		for _, perRoom := range room.Rooms {
			roomMinPricePerNight += perRoom.PricePerNight
		}
		if roomMinPricePerNight < minPrice {
			minPrice = roomMinPricePerNight
		}
	}
	return minPrice
}

func (g *hotelService) getUnavailableAmenities(amenities []dto.HotelAmenityDto) []dto.HotelAmenityDto {

	unavailableAmenities := make([]dto.HotelAmenityDto, 0)
	allAmenities := g.cacheStore.AmenityStore().GetAll()
	for _, candidate := range allAmenities {
		found := false
		for _, amenity := range amenities {
			if candidate.ID == amenity.ID {
				found = true
				break
			}
		}
		if !found {
			unavailableAmenities = append(unavailableAmenities, candidate)
		}
	}
	return unavailableAmenities
}

func (g *hotelService) GetRoomCancellationPolicy(hotelId, roomId, sessionId string) (*dto.RoomCancellationPolicyDto, error) {
	return g.provider.GetRoomCancellationPolicy(hotelId, roomId, sessionId)
}

func (g *hotelService) UpdateSomeHotels(hotelsDto dto.SyncSomeHotelsDto, date time.Time) (*dto.UpdateResultDto, error) {

	start := time.Now()

	length := len(hotelsDto.HotelIds)
	hotelChannel := make(chan *dbmodel.Hotel, length)
	for i := 0; i < length; i++ {
		go g.updateHotel(hotelsDto.HotelIds[i], "", date, hotelChannel)
	}
	for i := 0; i < length; i++ {
		hotel := <-hotelChannel
		if hotel == nil {
			continue
		}
		err := g.unitOfWork.Hotel().StoreOrUpdate(hotel)
		if err != nil {
			logger.WithName(logtags.CannotCreateOrUpdateHotel).Print(err.Error())
		}
	}

	g.cacheStore.UpdateAmenityStore()
	return &dto.UpdateResultDto{
		Message:    "Update some hotels completed",
		TimeTaken:  int64(time.Since(start)),
		ItemsCount: len(hotelsDto.HotelIds),
		Error:      nil,
	}, nil
}

func (g *hotelService) UpdateAll(date time.Time) (*dto.TaskRunningResult, error) {
	if inUpdating.Get() {
		return nil, common.AlreadyUpdatingHotels
	}
	inUpdating.Set(true)
	go func() {
		defer inUpdating.Set(false)
		_, err := g.UpdateAllSync(date)
		if err != nil {
			logger.WithName(logtags.UpdatingHotelsError).ErrorException(err, "error wile updating hotels")
			return
		}
		g.cacheStore.UpdateAmenityStore()
	}()
	return &dto.TaskRunningResult{
		Message: "Update all hotels task is now running in background",
		Success: true,
		Error:   nil,
	}, nil
}

func (g *hotelService) UpdateAllSync(date time.Time) (*dto.UpdateResultDto, error) {
	start := time.Now()

	ids, err := g.unitOfWork.Hotel().GetAllHotelIds()
	if err != nil {
		return nil, err
	}
	length := len(ids)

	for _, hotels := range array.Chunks(ids, g.syncChunkSize) {
		err := g.updateSync(date, hotels)
		if err != nil {
			logger.WithName(logtags.UpdatingHotelsError).ErrorException(err, "error while updating chunk of hotel ids")
		}
	}

	logger.WithName(logtags.UpdatingHotelsCompleted).WithData(fmt.Sprintf("updating hotels completed in %d nano seconds", time.Since(start))).Info("updating hotels completed")

	g.cacheStore.UpdateAmenityStore()

	return &dto.UpdateResultDto{
		Message:    fmt.Sprintf("updating hotels completed in %d nano seconds", time.Since(start)),
		TimeTaken:  int64(time.Since(start)),
		ItemsCount: length,
		Error:      nil,
	}, nil
}

func (g *hotelService) updateSync(date time.Time, ids []string) error {
	length := len(ids)
	hotelChannel := make(chan *dbmodel.Hotel, length)
	for i := 0; i < length; i++ {
		go g.updateHotel(ids[i], "", date, hotelChannel)
	}
	for i := 0; i < length; i++ {
		hotel := <-hotelChannel
		if hotel == nil {
			continue
		}
		err := g.unitOfWork.Hotel().StoreOrUpdate(hotel)
		if err != nil {
			logger.WithName(logtags.UpdatingHotelsError).WithException(err).
				Error("error wile updating all hotels sync")
		}
	}

	return nil
}

func (g *hotelService) SyncSomeHotels(hotelsDto dto.SyncSomeHotelsDto) (*dto.UpdateResultDto, error) {
	return g.UpdateSomeHotels(hotelsDto, time.Now())
}

func (g *hotelService) SyncedHotels() (*dto.SyncedHotelsDetail, error) {
	hotels, err := g.unitOfWork.Hotel().GetAllHotels()
	if err != nil {
		return nil, err
	}
	return g.mapper.ToHotelsDetail(hotels), nil
}

func (g *hotelService) SyncAllHotels() (*dto.TaskRunningResult, error) {
	if inSyncing.Get() {
		return nil, common.AlreadyInSyncing
	}
	inSyncing.Set(true)
	go func() {
		defer inSyncing.Set(false)
		start := time.Now()
		cities, _ := g.publicService.SyncAllCities()
		he := g.provider.CreateHotelEnumerable(cities)
		for he.MoveNext() {
			res, err := he.Current()
			if err != nil {
				logger.WithName(logtags.SyncHotelsError).WithException(err).
					Error("problem in getting list of hotels")
				continue
			}
			length := len(res)
			hotelChannel := make(chan *dbmodel.Hotel, length)
			for i := 0; i < length; i++ {
				go g.updateHotel(res[i].Id, res[i].HotelType, start, hotelChannel)
			}
			for i := 0; i < length; i++ {
				hotel := <-hotelChannel
				if hotel == nil {
					continue
				}
				err := g.unitOfWork.Hotel().StoreOrUpdate(hotel)
				if err != nil {
					logger.WithName(logtags.SyncHotelsError).WithException(err).
						Error("problem wile updating hotel information")
				}
			}
		}
		logger.WithName(logtags.SyncingHotelsCompleted).WithData(fmt.Sprintf("syncing hotels completed in %d nanoseconds", time.Since(start))).
			Info("syncing hotels completed")
		g.cacheStore.UpdateAmenityStore()
	}()

	return &dto.TaskRunningResult{
		Message: "syncing all hotels task is now running in background",
		Success: true,
		Error:   nil,
	}, nil
}

func (g *hotelService) updateHotel(hotelId, hotelType string, date time.Time, hotelChannel chan<- *dbmodel.Hotel) {
	hotelData, err := g.provider.GetHotelData(hotelId, date, date.Add(time.Hour*24))
	if err != nil {
		logger.WithException(err).
			WithName(logtags.GettingHotelDetailError).
			WithData(hotelId).
			Error("problem while getting hotel data")
		hotelChannel <- nil
		return
	}
	city, _ := g.cacheStore.CityStore().FindOne(hotelData.City)
	if hotelType == "" {
		t, err := g.provider.GetHotelType(hotelId)
		if err != nil {
			logger.WithName(logtags.GettingHotelTypeError).WithException(err).
				Error("problem while getting hotel type for update")
		}
		hotelType = t
	}
	hotelData.SetProvince(city.State)
	hotelData.SetType(hotelType)
	hotelModel := g.mapper.ToHotelModel(*hotelData)
	hotel, err := g.unitOfWork.Hotel().FindByIDForSync(hotelModel.PlaceID)
	if err == nil {
		hotelModel = hotel.UpdateWith(*hotelModel)
	}

	if err != nil && err != common.HotelNotFound {
		logger.WithException(err).WithName(logtags.GettingHotelDetailError).
			Error("problem while updating hotel information")
		hotelChannel <- nil
		return
	}
	hotelChannel <- hotelModel
	return
}

func (g *hotelService) GetHotelRooms(dto dto.HotelRoomsDto) (*dto.RateRoomResponseDto, error) {
	res, err := g.provider.GetHotelRooms(dto)
	if err != nil {
		return nil, err
	}
	res.RequestedRooms = dto.Rooms
	return res, nil
}
func (g *hotelService) GetHotelRoomsWithSession(dto dto.HotelRoomsWithSessionDto) (*dto.RateRoomResponseDto, error) {
	return g.provider.GetHotelRoomsWithSession(dto)
}
func (g *hotelService) HotelAvailable(body dto.AvailableDto) (*dto.AvailableResponseDto, error) {
	if err := hotelAvailableGuard(body.HotelId, body.PhoneNumber); err != nil {
		return nil, err
	}
	detail, err := g.provider.GetOrderDetail(body.HotelId, body.SessionId,
		body.OptionId)
	if err != nil {
		return nil, err
	}
	available, err := g.provider.HotelAvailable(body)
	if err != nil {
		if strings.Contains(err.Error(), "Room is not available") {
			return nil, common.RoomIsNotAvailable
		}
		return nil, err
	}
	err = CompareProviderAndResellerDates(available.CheckIn, available.CheckOut,
		body.CheckIn, body.CheckOut)
	if err != nil {
		return nil, err
	}

	detail.ProviderOrderId = available.OrderId
	detail.Status = available.Status
	detail.TotalPrice = available.TotalPrice
	detail.IndraOrderId = available.IndraOrderId
	_, err = g.unitOfWork.Order().Insert(g.mapper.ToOrderModel(*detail))
	if err != nil {
		return nil, err
	}
	logger.WithName(logtags.HotelAvailableCompleted).WithData(available).
		Info("hotel available completed successfully")
	return available, nil
}

func CompareProviderAndResellerDates(pCheckIn, pCheckOut, rCheckIn, rCheckOut string) error {
	isCheckInValid, err := date.CompareTwoDates(rCheckIn, pCheckIn)
	if err != nil {
		return err
	}
	isCheckOutValid, err := date.CompareTwoDates(rCheckOut, pCheckOut)
	if err != nil {
		return err
	}
	if isCheckOutValid && isCheckInValid {
		return nil
	}
	logger.WithName(logtags.DatesNotMatchError).WithData(map[string]string{
		"resellerCheckIn":  rCheckIn,
		"resellerCheckOut": rCheckOut,
		"providerCheckIn":  pCheckIn,
		"providerCheckOut": pCheckOut,
	}).Error("check in or check out dates are not match")
	return common.DatesNotMatchError
}

func (g *hotelService) FinalizeHotelOrder(request dto.FinalizeOrderDto) (*dto.FinalizeOrderResponseDto, error) {
	_, err := g.ConfirmOrder(request.OrderId)
	if err != nil {
		return nil, err
	}
	payResult, err := g.PayByAccount(request.OrderId)
	if err != nil {
		return nil, err
	}
	statusResult, err := g.GetOrderStatus(request.OrderId)
	if err != nil {
		return nil, err
	}
	successRes := &dto.FinalizeOrderResponseDto{
		PaymentResult: &payResult,
		StatusResult:  &statusResult,
		Error:         nil,
	}
	logger.WithName(logtags.OrderFinalizedCompleted).WithData(successRes).
		Info("order finalized successfully")
	return successRes, nil
}

func (g *hotelService) GetAnOrderDetail(orderId string) (*dto.OrderDetailDto, error) {
	order, err := g.unitOfWork.Order().GetOneByIndraId(orderId)
	if err != nil {
		return nil, err
	}
	hotel, err := g.unitOfWork.Hotel().GetHotel(order.ProviderHotelId)
	if err != nil {
		return nil, err
	}
	orderDto := g.mapper.ToOrderDto(*order)
	orderDto.Hotel = *g.mapper.ToHotelDto(*hotel)
	return &orderDto, nil
}

func (g *hotelService) SearchResult(request dto.SearchDto) (*dto.SearchResponseDto, error) {
	daysDiff, err := date.DaysDiff(request.Date.Start, request.Date.End)
	if err != nil {
		return nil, err
	}
	searchDto, err := g.searchDtoAdopter.CreateProviderDto(request)
	if err != nil {
		return nil, err
	}
	results, totalHits, err := g.provider.SearchResult(*searchDto)
	if err != nil {
		return nil, err
	}
	hotels, err := g.getHotelsFromResults(results)
	if err != nil {
		return nil, err
	}
	res, err := g.searchDtoAdopter.CreateResultDto(results, hotels,
		searchDto.RequestSession.Destination.Id, daysDiff)
	if err != nil {
		return nil, err
	}
	filters := g.publicService.GetAllFilters()
	res.Sorts = g.publicService.GetAllSorts()
	res.Filters = []dto.FilterOutput{
		filters.PriceFilter,
		filters.AmenityFilters,
		filters.ScoreFilters,
		filters.StarFilters,
		filters.HotelTypes,
	}
	res.TotalHits = totalHits
	logger.WithName(logtags.SearchCompleted).Info("Search hotels completed successfully")
	return res, nil
}

func (g *hotelService) SetAmenityIcon(request dto.SetAmenityIconDto) (dto.HotelAmenityDto, error) {
	amenity, err := g.unitOfWork.Amenity().UpdateIcon(request.AmenityId, request.IconUrl)
	if err != nil {
		return dto.HotelAmenityDto{}, err
	}
	g.cacheStore.UpdateAmenityStore()
	return *g.mapper.ToAmenityDto(*amenity), nil
}

func (g *hotelService) SetAmenityCategory(request dto.SetAmenityCategoryDto) (dto.HotelAmenityDto, error) {
	_, err := g.unitOfWork.AmenityCategory().GetOneById(request.AmenityCategoryId)
	if err != nil {
		return dto.HotelAmenityDto{}, err
	}
	amenity, err := g.unitOfWork.Amenity().UpdateCategory(request.AmenityId, request.AmenityCategoryId)
	if err != nil {
		return dto.HotelAmenityDto{}, err
	}
	g.cacheStore.UpdateAmenityStore()
	return *g.mapper.ToAmenityDto(*amenity), nil
}

func (g *hotelService) getHotelsFromResults(results []dtos.Result) ([]dbmodel.Hotel, error) {
	ids := make([]string, 0)
	for _, result := range results {
		ids = append(ids, result.ID)
	}
	return g.unitOfWork.Hotel().GetHotels(ids)
}

func (g *hotelService) GetHotelOptionInfo(infoDto dto.OptionInfoRequestDto) (*dto.OptionInfoResponseDto, error) {
	return g.provider.GetHotelOptionInfo(infoDto)
}

func (g *hotelService) GetHotels(ids []string) ([]dto.HotelDto, error) {
	hotels, _ := g.unitOfWork.Hotel().GetHotels(ids)
	result := make([]dto.HotelDto, 0)
	for _, hotel := range hotels {
		result = append(result, *g.mapper.ToHotelDto(hotel))
	}
	return result, nil
}

func (g *hotelService) ConfirmOrder(orderId string) (dto.ConfirmResponseDto, error) {
	order, err := g.unitOfWork.Order().GetOneByIndraId(orderId)
	if err != nil {
		return dto.ConfirmResponseDto{}, err
	}
	if order.Confirmed {
		return dto.ConfirmResponseDto{
			OrderId: orderId,
			Error:   nil,
		}, nil
	}
	order.UpdateConfirmed(true)
	res, err := g.provider.ConfirmOrder(orderId)
	if err != nil {
		return dto.ConfirmResponseDto{}, err
	}
	err = g.unitOfWork.Order().StoreOrUpdate(*order)
	logger.WithName(logtags.ConfirmOrderCompleted).WithData(res).
		Info("Confirm order completed successfully")
	return res, err
}

func (g *hotelService) PayByAccount(orderId string) (dto.OrderPayByAccountResponseDto, error) {
	order, err := g.unitOfWork.Order().GetOneByIndraId(orderId)
	if err != nil {
		return dto.OrderPayByAccountResponseDto{}, err
	}
	if order.TransactionRequestId != "" {
		return dto.OrderPayByAccountResponseDto{
			TransactionStatus: order.TransactionStatus,
			RequestId:         order.TransactionRequestId,
			TransactionIds:    strings.Split(order.TransactionIds, ","),
			ResultMessage:     "",
			Error:             nil,
		}, nil
	}
	res, err := g.provider.PayByAccount(orderId)
	go g.balanceChecker.CheckAdequateBalance()
	if err != nil {
		return res, err
	}
	if res.TransactionStatus == "Pending" {
		logger.WithName(logtags.CannotCompleteOrderPayment).Error("hotel payment status is Pending. please check if the alibaba hotel client has adequate balance for payment")
	}
	order.UpdateTransaction(res.TransactionStatus, strings.Join(res.TransactionIds, ","),
		res.RequestId)
	err = g.unitOfWork.Order().StoreOrUpdate(*order)

	logger.WithName(logtags.PayByAccountCompleted).WithData(res).
		Info("Pay order by account completed successfully")
	return res, err
}

func (g *hotelService) GetOrderStatus(orderId string) (dto.OrderStatusResponseDto, error) {
	order, err := g.unitOfWork.Order().GetOneByIndraId(orderId)
	if err != nil {
		return dto.OrderStatusResponseDto{}, err
	}
	res, err := g.provider.GetOrderStatus(orderId)
	if err != nil {
		return res, err
	}
	order.UpdateStatus(res.Status)
	err = g.unitOfWork.Order().StoreOrUpdate(*order)
	return res, err
}

func (g *hotelService) GetOrderEnquiry(orderId string) (dto.OrderEnquiryResponseDto, error) {
	order, err := g.unitOfWork.Order().GetOneByIndraId(orderId)
	if err != nil {
		return dto.OrderEnquiryResponseDto{}, err
	}
	return g.provider.GetOrderEnquiry(orderId, order.ProviderOrderId)
}

func (g *hotelService) RefundOrder(refundRequest dto.OrderRefundRequestDto) (dto.OrderRefundResponseDto, error) {
	order, err := g.unitOfWork.Order().GetOneByIndraId(refundRequest.OrderId)
	if err != nil {
		return dto.OrderRefundResponseDto{}, err
	}
	if order.RefundRequestId != 0 {
		return dto.OrderRefundResponseDto{
			OrderId:         refundRequest.OrderId,
			RefundRequestId: order.RefundRequestId,
			Error:           nil,
		}, nil
	}
	res, err := g.provider.RefundOrder(refundRequest.OrderId, order.ProviderOrderId)
	if err != nil {
		return res, err
	}
	order.UpdateRefundRequestId(res.RefundRequestId, refundRequest.RefundRequestID, refundRequest.JabamaOrderID)
	err = g.unitOfWork.Order().StoreOrUpdate(*order)
	logger.WithName(logtags.NewRefundRequest).WithData(refundRequest).
		Info(fmt.Sprintf("order with id %s commited a refund request", refundRequest.OrderId))
	return res, err
}

func (g *hotelService) UpdateHotelRateReview(rateDto dto.RateReviewEventDto) error {
	hotel, err := g.unitOfWork.Hotel().GetHotel(rateDto.PlaceId)
	if err != nil {
		logger.WithName(logtags.GettingHotelDetailError).ErrorException(err, "error while trying to find hotel to update rate and review details")
		return err
	}
	hotel.UpdateRateAndReview(rateDto.ReviewsCount, rateDto.Rating)
	err = g.unitOfWork.Hotel().StoreOrUpdate(hotel)
	if err != nil {
		logger.WithName(logtags.UpdatingRateAndReviewError).ErrorException(err, "error while trying to updating rate and review details")
		return err
	}
	return nil
}

func (g *hotelService) UpdateRefundedOrdersPaymentStatus(fromDate time.Time) {
	ids, err := g.unitOfWork.Order().GetProperOrderIdsForRefundUpdateStatus(fromDate)
	if err != nil {
		logger.WithName(logtags.GettingListOfRefundableOrdersError).
			ErrorException(err, "error while trying to get list of proper order ids for updating refund status")
		return
	}
	for _, orderIds := range array.Chunks(ids, g.syncChunkSize) {
		length := len(orderIds)
		orderChannel := make(chan *dbmodel.Order, length)
		g.tryUpdatingOrdersRefundStatus(orderIds, orderChannel)
		for i := 0; i < length; i++ {
			order := <-orderChannel
			if order == nil {
				continue
			}
			err := g.unitOfWork.Order().StoreOrUpdate(*order)
			if err != nil {
				logger.WithName(logtags.CannotCreateOrUpdateHotel).WithException(err).
					Error("error wile updating an order refund status")
				continue
			}
			g.orderEventDispatcher.OrderRefundRequestFinalized(dto.OrderRefundRequestFinalizedDto{
				ApplicantRefundRequestId: order.ApplicantRefundRequestId,
				ApplicantOrderId:         order.ApplicantOrderId,
				PaidAmount:               order.PaidAmount,
				RefundStatus:             order.RefundStatus,
				RefundableAmount:         order.RefundableAmount,
				TotalPenaltyAmount:       order.TotalPenaltyAmount,
				ProviderOrderId:          strconv.FormatInt(order.IndraOrderId, 10),
			})
		}
	}
}

func (g *hotelService) tryUpdatingOrdersRefundStatus(orderIds []string, orderChannel chan<- *dbmodel.Order) {
	ordersLength := len(orderIds)
	ordersRefundStatus, err := g.provider.GetOrdersRefundStatus(orderIds, len(orderIds), 1)
	if err != nil {
		for i := 0; i < ordersLength; i++ {
			orderChannel <- nil
		}
		return
	}
	for _, id := range orderIds {
		details, founded := getOrderStatusDetails(ordersRefundStatus.Result.Items, id)
		if !founded {
			orderChannel <- nil
		}
		go func(details dto.OrdersRefundStatusResponseDtoResultItem, orderId string) {
			order, err := g.unitOfWork.Order().GetOneByIndraId(orderId)
			if err != nil {
				logger.WithName(logtags.GettingHotelDetailError).ErrorException(err, err.Error())
				orderChannel <- nil
				return
			}
			if len(details.Items) == 0 {
				orderChannel <- nil
				return
			}
			item := details.Items[0]
			if details.RefundStatus != common.RefundStatus_PaymentFinalized {
				orderChannel <- nil
				return
			}
			order.UpdateRefundResult(item.PaidAmount, item.ReferenceCode, details.RefundStatus,
				item.RefundableAmount, item.TotalPenaltyAmount)
			orderChannel <- order
		}(details, id)
	}
}

func (g *hotelService) GetHotelsList(body dto.HotelsPageRequestDto) (dto.HotelsPageResponseDto, error) {
	hotels, total, err := g.unitOfWork.Hotel().GetHotelsList(body.PageNumber, body.PageSize, body.Search)
	if err != nil {
		logger.WithName(logtags.GettingHotelsListError).ErrorException(err, err.Error())
		return dto.HotelsPageResponseDto{}, err
	}
	return dto.HotelsPageResponseDto{
		PageNumber: body.PageNumber,
		PageSize:   body.PageSize,
		Total:      total,
		Hotels:     g.mapper.ToHotelsDto(hotels),
		Error:      nil,
	}, nil
}

func (g *hotelService) SetHotelSeoDetails(body dto.SetHotelSeoRequestDto) (dto.HotelDto, error) {
	hotel, err := g.unitOfWork.Hotel().FindByID(body.HotelId)
	if err != nil {
		logger.WithName(logtags.SearchHotelsRequest).ErrorException(err, err.Error())
		return dto.HotelDto{}, err
	}
	hotel.UpdateSeoTags(body.SeoDetails.Title, body.SeoDetails.H1, body.SeoDetails.Description,
		body.SeoDetails.Robots, body.SeoDetails.Canonical, body.SeoDetails.MetaDescription)
	err = g.unitOfWork.Hotel().StoreOrUpdate(hotel)
	if err != nil {
		logger.WithName(logtags.CannotCreateOrUpdateHotel).ErrorException(err, err.Error())
		return dto.HotelDto{}, err
	}
	return *g.mapper.ToHotelDto(*hotel), nil
}

func (g *hotelService) SetHotelFaq(body dto.SetHotelFaqRequestDto) (dto.HotelDto, error) {
	hotel, err := g.unitOfWork.Hotel().FindByID(body.HotelId)
	if err != nil {
		logger.WithName(logtags.SearchHotelsRequest).ErrorException(err, err.Error())
		return dto.HotelDto{}, err
	}
	hotel.UpdateFaqDetails(body.FAQTitle, g.mapper.ToFAQsModel(body.FAQList))
	err = g.unitOfWork.Hotel().StoreOrUpdate(hotel)
	if err != nil {
		logger.WithName(logtags.CannotCreateOrUpdateHotel).ErrorException(err, err.Error())
		return dto.HotelDto{}, err
	}
	return *g.mapper.ToHotelDto(*hotel), nil
}

func (g *hotelService) RemoveHotelFaq(hotelId string, faqId uint) (dto.HotelDto, error) {

	hotel, err := g.unitOfWork.Hotel().FindByID(hotelId)
	if err != nil {
		logger.WithName(logtags.SearchHotelsRequest).ErrorException(err, err.Error())
		return dto.HotelDto{}, err
	}
	faq, err := hotel.GetHotelFAQ(faqId)
	if err != nil {
		return dto.HotelDto{}, err
	}
	hotel, err = g.unitOfWork.Hotel().RemoveFAQ(hotelId, &faq)
	if err != nil {
		logger.WithName(logtags.CannotRemoveHotelFAQ).ErrorException(err, err.Error())
		return dto.HotelDto{}, err
	}
	return *g.mapper.ToHotelDto(*hotel), nil
}

func getOrderStatusDetails(items []dto.OrdersRefundStatusResponseDtoResultItem, orderId string) (dto.OrdersRefundStatusResponseDtoResultItem, bool) {
	id, err := strconv.Atoi(orderId)
	if err != nil {
		return dto.OrdersRefundStatusResponseDtoResultItem{}, false
	}
	id64 := int64(id)
	for _, item := range items {
		if item.OrderId == id64 {
			return item, true
		}
	}
	return dto.OrdersRefundStatusResponseDtoResultItem{}, false
}

func hotelAvailableGuard(hotelId, phone string) error {
	c := config.Get()

	if c.IsProduction() {
		return nil
	}

	for _, s := range c.AvailableHotelsWhiteList {
		if s == hotelId {
			return nil
		}
	}

	return common.HotelReserveForbidden
}

func NewHotelService(unit core.UnitOfWork, mapper core.Mapper, provider core.HotelProvider,
	searchDtoAdopter core.SearchDtoAdopter, publicService core.PublicService,
	cacheStore core.CacheStore, balanceChecker core.ProviderBalanceChecker,
	orderEventDispatcher core.OrderEventDispatcher) core.HotelService {
	con := config.Get()
	return &hotelService{
		mapper:               mapper,
		provider:             provider,
		unitOfWork:           unit,
		searchDtoAdopter:     searchDtoAdopter,
		publicService:        publicService,
		cacheStore:           cacheStore,
		balanceChecker:       balanceChecker,
		syncChunkSize:        con.SyncChunkSize,
		orderEventDispatcher: orderEventDispatcher,
	}
}
