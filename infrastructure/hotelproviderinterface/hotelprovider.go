package hotelProviderInterface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/config"
	"hotel-engine/infrastructure/constants"
	"hotel-engine/infrastructure/hotelproviderinterface/dtos"
	"hotel-engine/infrastructure/logger"
	"hotel-engine/utils/httphelper"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"

	"github.com/tidwall/gjson"
)

const (
	TokenEndpoint               = "/api/v1/hoteladmin/login"
	HotelListEndpoint           = "/api/v1/hoteladmin/get-hotels"
	SearchDirectEndpoint        = "/api/v1/hotel/search/direct"
	HotelPriceEndpoint          = "/api/v1/hotel/rate/room"
	HotelOptionInfoEndpoint     = "/api/v1/hotel/general/info"
	AvailableEndpoint           = "/api/v2/hotel/book/available"
	SearchEndpoint              = "/api/v1/hotel/search"
	SearchResultEndpoint        = "/api/v1/hotel/result"
	RoomCancellationPolicy      = "/api/v1/hotel/rooms/getCancellationPolicy"
	ConfirmOrderEndPoint        = "/api/v1/coordinator/order/{orderId}/confirm"
	PayByBankAndAccountEndpoint = "/api/v1/coordinator/order/{orderId}/pay-by-bank-and-account"
	GetOrderStatusEndpoint      = "/api/v1/coordinator/order/{orderId}/status"
	OrderRefundEndpoint         = "/api/v1/profile/refunds"
	EnquiryOrderEndpoint        = "/api/v1/profile/refunds/enquiry/{orderId}"
	OrdersRefundStatusEndpoint  = "/api/v1/management/refunds"
)

type hotelProvider struct {
	password            string
	username            string
	baseHotelUrl        string
	baseOrderUrl        string
	client              *http.Client
	baseOrderServiceUrl string
}

func (p *hotelProvider) GetAccessToken() (string, error) {
	logger.Debug("trying to get a new token")
	body, err := json.Marshal(map[string]string{
		"username": p.username,
		"password": p.password,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", p.baseHotelUrl+TokenEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Close = true
	req.Header.Add("ab-channel", common.ABChannelName)
	logger.WithName(logtags.GetAccessTokenRequest).WithData(body).Info("Get access token request log")
	res, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	err = httphelper.GetResponseError(res)
	if err != nil {
		return "", err
	}
	var result dtos.AccessTokenResponse
	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {
		return "", err
	}
	logger.Debug("new token has been taken successfully")
	return result.Result.Token, nil
}

func (p *hotelProvider) SearchHotels(limit, skip uint32, hotelGiataId, cityId, hotelName string, cityBaseId int64) (*dtos.HotelsListResult, error) {
	requestBody := map[string]interface{}{
		"hotelGiataId": hotelGiataId,
		"cityId":       cityId,
		"hotelName":    hotelName,
	}
	if cityBaseId >= 0 {
		requestBody["cityBaseId"] = cityBaseId
	}

	body, err := json.Marshal(requestBody)

	if err != nil {
		return nil, err
	}
	urlParam := fmt.Sprintf("?limit=%d&skip=%d", limit, skip)
	req, err := http.NewRequest("POST", p.baseHotelUrl+HotelListEndpoint+urlParam, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}
	token, err := Token()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("ab-channel", common.ABChannelName)

	req.Close = true

	res, err := p.client.Do(req)

	if err != nil {
		return nil, err
	}
	var result dtos.HotelsListResponse
	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {
		return nil, err
	}
	return &result.Result, nil
}

func (p *hotelProvider) GetHotels(limit, skip uint32) (*dtos.HotelsListResult, error) {
	return p.SearchHotels(limit, skip, "", "", "", -1)
}

func (p *hotelProvider) GetHotelData(hotelId string, checkIn, checkout time.Time) (*dto.HotelDto, error) {
	body := dtos.NewSearchDirectRequest(hotelId, checkIn, checkout).ToJson()
	req, err := http.NewRequest("POST", p.baseHotelUrl+SearchDirectEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("ab-channel", common.ABChannelName)

	req.Close = true

	res, err := p.client.Do(req)

	if err != nil {
		return nil, err
	}
	err = httphelper.GetResponseError(res)
	if err != nil {
		return nil, err
	}
	var searchRes dtos.SearchDirectResponse
	err = json.NewDecoder(res.Body).Decode(&searchRes)

	if err != nil {
		return nil, err
	}

	roomReqBody, err := json.Marshal(map[string]string{
		"sessionId": searchRes.Result.SessionId,
		"hotelId":   searchRes.Result.HotelId,
	})

	delays := []uint{1, 2, 4, 8, 16, 32, 64, 128}
	maxTry := len(delays)
	for i := 0; i < maxTry; i++ {
		req, err = http.NewRequest("POST", p.baseHotelUrl+HotelPriceEndpoint, bytes.NewBuffer(roomReqBody))
		if err != nil {
			return nil, err
		}
		req.Header.Add("ab-channel", common.ABChannelName)
		req.Close = true

		time.Sleep(time.Duration(delays[i]) * time.Second)
		res, err = p.client.Do(req)
		if err != nil {
			return nil, err
		}
		err = httphelper.GetResponseError(res)
		if err != nil {
			return nil, err
		}
		jsonData, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return nil, err
		}
		jsonString := string(jsonData)
		finalResult := gjson.Get(jsonString, "result.finalResult").Bool()
		if !finalResult {
			continue
		}
		return getHotelDto(jsonString, checkIn)
	}
	return nil, common.CannotGetHotelDataForSync
}

func (p *hotelProvider) CreateHotelEnumerable(cities []dto.CityDto) core.HotelEnumerable {
	return NewHotelEnumerable(p, cities)
}

func getHotelDto(jsonString string, date time.Time) (*dto.HotelDto, error) {
	if !gjson.Valid(jsonString) {
		logger.WithName(logtags.InvalidJsonResponseError).WithData(map[string]interface{}{
			"response": jsonString,
		}).Error(common.JsonDataIsNotValid.Error())
		return nil, common.JsonDataIsNotValid
	}

	hotel := gjson.Get(jsonString, "result.hotel")
	rooms := gjson.Get(jsonString, "result.rooms").Array()

	placeID := hotel.Get("hotelId").String()
	if placeID == "" {
		placeID = hotel.Get("_id").String()
	}

	price := int64(0)
	oldPrice := int64(0)
	discountPercent := int64(0)
	discountPrice := int64(0)
	roomID := "0"
	if len(rooms) != 0 {
		roomID = rooms[0].Get("id").String()
		price = rooms[0].Get("price").Int()
		if rooms[0].Get("oldPrice").Exists() {
			oldPrice = rooms[0].Get("oldPrice").Int()
		}
	}
	if oldPrice != 0 && oldPrice > price {
		discountPercent = int64(math.Round((float64(oldPrice-price) / float64(oldPrice)) * 100))
		discountPrice = oldPrice - price
	}

	//geoLocation
	geoLocation := ""
	coordinates := hotel.Get("location.coordinates").Array()
	if len(coordinates) == 2 {
		geoLocation = fmt.Sprintf("%f", coordinates[1].Float()) + "," +
			fmt.Sprintf("%f", coordinates[0].Float())
	}

	//images
	images := []string{}
	for _, i := range hotel.Get("images.#.url").Array() {
		images = append(images, i.String())
	}

	//amenities
	amenities := []dto.HotelAmenityDto{}
	for _, f := range hotel.Get("facilities").Array() {
		amenities = append(amenities, dto.HotelAmenityDto{
			ID:      int(f.Get("id").Int()),
			Name:    f.Get("name.fa").String(),
			NameEn:  strings.ToLower(f.Get("name.en").String()),
			GroupId: int(f.Get("groupId").Int()),
		})
	}

	places := make([]dto.HotelPlaceDto, 0)
	for _, f := range hotel.Get("places").Array() {

		placeLocation := ""
		placeLocations := f.Get("location").Array()

		if len(placeLocations) == 2 {
			placeLocation = fmt.Sprintf("%f", placeLocations[1].Float()) + "," +
				fmt.Sprintf("%f", placeLocations[0].Float())
		}
		places = append(places, dto.HotelPlaceDto{
			ID:          f.Get("_id").String(),
			Name:        f.Get("name").String(),
			GeoLocation: placeLocation,
			Distance:    f.Get("distance").Float(),
			Priority:    int(f.Get("priority").Int()),
		})
	}

	badges := make([]dto.HotelBadgeDto, 0)
	for _, f := range hotel.Get("badges").Array() {
		badges = append(badges, dto.HotelBadgeDto{
			ID:              f.Get("_id").String(),
			Text:            f.Get("text").String(),
			Icon:            f.Get("icon").String(),
			TextColor:       f.Get("color.text").String(),
			BackgroundColor: f.Get("color.background").String(),
		})
	}

	return &dto.HotelDto{
		PlaceID:         placeID,
		RoomID:          roomID,
		Type:            common.HotelType_Hotel,
		Kind:            constants.DefaultHotelSyncKind,
		Region:          constants.DefaultHotelSyncRegion,
		ReservationType: constants.DefaultHotelSyncReservationType,
		PaymentType:     constants.DefaultHotelSyncPaymentType,
		MinNight:        1,
		Tags:            []string{"hotel"},
		SuitableFor:     []string{},
		Verified:        true,
		Name:            hotel.Get("name.fa").String(),
		NameEn:          hotel.Get("name.en").String(),
		Description:     strip.StripTags(hotel.Get("description.fa").String()),
		Images:          images,
		Amenities:       amenities,
		Places:          places,
		Address:         hotel.Get("address").String(),
		CheckInTime:     hotel.Get("checkinTime").String(),
		CheckOutTime:    hotel.Get("checkoutTime").String(),
		City:            hotel.Get("city.fa").String(),
		CityEn:          hotel.Get("city.en").String(),
		Province:        hotel.Get("state.fa").String(),
		ProvinceEn:      hotel.Get("state.en").String(),
		Capacity:        2,
		CheckIn:         date,
		CheckOut:        date.Add(time.Hour * 24),
		GeoLocation:     geoLocation,
		Price:           price,
		RateReviewCount: 0,
		RateReviewScore: 0,
		Star:            int(hotel.Get("star").Int()),
		Country:         hotel.Get("country.fa").String(),
		CountryEn:       hotel.Get("country.en").String(),
		CountryCode:     hotel.Get("country.code").String(),
		OldPrice:        oldPrice,
		DiscountPercent: discountPercent,
		DiscountPrice:   discountPrice,
		Badges:          badges,
		Sort:            hotel.Get("score").Float(),
	}, nil
}

func (p *hotelProvider) SearchResult(searchDto dtos.ProviderSearchDto) ([]dtos.Result, int, error) {
	if searchDto.RequestSearchHotels.SessionId == "" {
		sessionId, err := p.getSearchSessionId(searchDto.RequestSession)
		if err != nil {
			return nil, 0, err
		}
		searchDto.RequestSearchHotels.SessionId = sessionId
	}

	body := searchDto.RequestSearchHotels.ToJson()
	delays := []uint{0, 200, 300, 500, 800, 1300, 2100, 3400}
	maxTry := len(delays)
	for i := 0; i < maxTry; i++ {
		req, err := http.NewRequest("POST", p.baseOrderUrl+SearchResultEndpoint, bytes.NewBuffer(body))
		if err != nil {
			return nil, 0, err
		}
		req.Header.Add("ab-channel", common.ABChannelName)
		req.Close = true

		time.Sleep(time.Duration(delays[i]) * time.Millisecond)
		res, err := p.client.Do(req)
		if err != nil {
			return nil, 0, err
		}
		if err = httphelper.GetResponseError(res); err != nil {
			return nil, 0, err
		}
		var result dtos.ResultResponse

		if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
			return nil, 0, err
		}
		if result.Result.Progress >= 100 {
			return result.Result.Result, result.Result.Info.ResultNo, nil
		}
	}
	return nil, 0, common.GetHotelListMaxTryLimitReached
}

func (p *hotelProvider) getSearchSessionId(data dtos.ProviderRequestSessionDto) (string, error) {
	body := data.ToJson()
	req, err := http.NewRequest("POST", p.baseOrderUrl+SearchEndpoint, bytes.NewBuffer(body))

	if err != nil {
		return "", err
	}
	req.Header.Add("ab-channel", common.ABChannelName)

	req.Close = true

	res, err := p.client.Do(req)

	if err != nil {
		return "", err
	}
	err = httphelper.GetResponseError(res)
	if err != nil {
		return "", err
	}

	var searchRes dtos.SearchResponse
	err = json.NewDecoder(res.Body).Decode(&searchRes)

	if err != nil {
		return "", err
	}
	return searchRes.Result.SessionId, nil
}

func (p *hotelProvider) DirectHotel(data dto.HotelRoomsDto) (*dto.DirectResponseDto, error) {
	requestModel := dtos.NewDirectRequestDto(data.HotelId,
		data.CheckIn, data.CheckOut, data.Rooms)
	body, err := json.Marshal(requestModel)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", p.baseOrderUrl+SearchDirectEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("ab-channel", common.ABChannelName)

	req.Close = true

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	err = httphelper.GetResponseError(res)
	if err != nil {
		return nil, err
	}

	var searchRes dtos.SearchDirectResponse
	err = json.NewDecoder(res.Body).Decode(&searchRes)

	if err != nil {
		return nil, err
	}

	return &dto.DirectResponseDto{
		SessionId: searchRes.Result.SessionId,
		HotelId:   searchRes.Result.HotelId,
		Error:     nil,
	}, nil
}

func (p *hotelProvider) GetHotelRooms(data dto.HotelRoomsDto) (*dto.RateRoomResponseDto, error) {
	sessionDto, err := p.DirectHotel(data)
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(map[string]string{
		"sessionId": sessionDto.SessionId,
		"hotelId":   sessionDto.HotelId,
	})

	if err != nil {
		return nil, err
	}
	rooms, err := p.getHotelRoomsWithDelay(body)
	return &dto.RateRoomResponseDto{
		Rooms:     rooms,
		HotelId:   data.HotelId,
		SessionId: sessionDto.SessionId,
		Error:     nil,
	}, err
}
func (p *hotelProvider) GetHotelRoomsWithSession(data dto.HotelRoomsWithSessionDto) (*dto.RateRoomResponseDto, error) {

	body, err := json.Marshal(map[string]string{
		"sessionId": data.SessionId,
		"hotelId":   data.HotelId,
	})

	if err != nil {
		return nil, err
	}
	rooms, err := p.getHotelRoomsWithDelay(body)
	return &dto.RateRoomResponseDto{
		Rooms:     rooms,
		HotelId:   data.HotelId,
		SessionId: data.SessionId,
		Error:     nil,
	}, err
}

func (p *hotelProvider) GetHotelType(hotelId string) (string, error) {
	result, err := p.SearchHotels(1, 0, hotelId, "", "", -1)
	if err != nil {
		return common.HotelType_Hotel, err
	}

	if len(result.HotelsList) == 0 {
		return common.HotelType_Hotel, common.HotelNotFound
	}
	typeCode, err := strconv.Atoi(result.HotelsList[0].AccommodationType)
	if err != nil {
		typeCode = common.HotelTypes[common.HotelType_Hotel]
	}
	typeString := common.HotelTypesString[typeCode]
	if typeString == "" {
		typeString = common.HotelType_Hotel
	}
	return typeString, nil
}

func (p *hotelProvider) GetHotelOptionInfo(infoDto dto.OptionInfoRequestDto) (*dto.OptionInfoResponseDto, error) {
	data := map[string]string{
		"sessionId": infoDto.SessionId,
		"hotelId":   infoDto.HotelId,
		"optionId":  infoDto.OptionId,
	}
	body, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", p.baseOrderUrl+HotelOptionInfoEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("ab-channel", common.ABChannelName)
	logger.WithName(logtags.GetHotelOptionInfoRequest).WithData(data).Info("Get hotel options info request log")
	req.Close = true

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	err = httphelper.GetResponseError(res)
	if err != nil {
		return nil, err
	}

	jsonData, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}
	jsonString := string(jsonData)

	policiesArray := gjson.Get(jsonString, "result.policy.general.policies").Array()
	policies := make([]string, 0, len(policiesArray))
	for _, policy := range policiesArray {
		policies = append(policies, policy.String())
	}

	roomsArray := gjson.Get(jsonString, "result.detail.rooms").Array()
	rooms := make([]dto.OptionInfoDetailRoomDto, 0, len(roomsArray))
	for _, room := range roomsArray {
		rooms = append(rooms, dto.OptionInfoDetailRoomDto{
			Name: room.Get("name").String(),
		})
	}

	return &dto.OptionInfoResponseDto{
		Detail: dto.OptionInfoDetailDto{
			Price:    gjson.Get(jsonString, "result.detail.price").Int(),
			HotelId:  infoDto.HotelId,
			Provider: gjson.Get(jsonString, "result.detail.provider").String(),
			Currency: gjson.Get(jsonString, "result.detail.currency").String(),
			RestrictedMarkup: dto.OptionInfoDetailRestrictedMarkupDto{
				Amount: gjson.Get(jsonString, "result.detail.restrictedMarkup.amount").Int(),
				Type:   gjson.Get(jsonString, "result.detail.restrictedMarkup.type").String(),
			},
			MealPlan: gjson.Get(jsonString, "result.detail.mealPlan").String(),
			Rooms:    rooms,
		},
		Policy: dto.OptionCancellationDto{
			NonRefundable:   gjson.Get(jsonString, "result.policy.nonRefundable").Bool(),
			GeneralPolicies: policies,
		},
		Error: nil,
	}, err
}

func (p *hotelProvider) getHotelRoomsWithDelay(body []byte) ([]dto.RoomOptionDto, error) {
	delays := []uint{800, 50, 50, 150, 200, 300}
	maxTry := len(delays)
	for i := 0; i < maxTry; i++ {
		req, err := http.NewRequest("POST", p.baseOrderUrl+HotelPriceEndpoint, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		req.Header.Add("ab-channel", common.ABChannelName)
		req.Close = true

		time.Sleep(time.Duration(delays[i]) * time.Millisecond)
		res, err := p.client.Do(req)
		if err != nil {
			return nil, err
		}
		err = httphelper.GetResponseError(res)
		if err != nil {
			return nil, err
		}

		jsonData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		jsonString := string(jsonData)
		finalResult := gjson.Get(jsonString, "result.finalResult").Bool()
		if !finalResult {
			continue
		}
		roomsJsonArray := gjson.Get(jsonString, "result.rooms").Array()
		if len(roomsJsonArray) == 0 {
			continue
		}
		logger.Debug(fmt.Sprintf("rooms data fetched after %d tries", i+1))
		rooms := make([]dto.RoomOptionDto, 0)
		for i, r := range roomsJsonArray {
			oldPrice := int64(0)
			discountPrice := int64(0)
			discountPercent := int64(0)
			price := int64(0)
			roomList := make([]dto.RoomDto, 0)
			roomsArray := r.Get("rooms").Array()
			price = r.Get("price").Int()
			for j, each := range roomsArray {
				roomList = append(roomList, dto.RoomDto{
					ExtraCharge:   make([]interface{}, 0),
					PricePerNight: each.Get("pricePerNight").Int(),
					Name:          each.Get("name").String(),
					NameEn:        each.Get("name_en").String(),
					Price:         each.Get("price").Int(),
					Number:        j + 1,
				})
			}
			mealPlan := dto.MealPlans["Unknown"]
			if val, ok := dto.MealPlans[r.Get("mealPlan").String()]; ok {
				mealPlan = val
			}

			if r.Get("oldPrice").Exists() {
				oldPrice = r.Get("oldPrice").Int()
			}
			if oldPrice != 0 && oldPrice > price {
				discountPercent = int64(math.Round((float64(oldPrice-price) / float64(oldPrice)) * 100))
				discountPrice = oldPrice - price
			}

			rooms = append(rooms, dto.RoomOptionDto{
				Id:              r.Get("id").String(),
				ProviderName:    r.Get("providerName").String(),
				Provider:        r.Get("provider").String(),
				NonRefundable:   r.Get("nonRefundable").Bool(),
				MealPlan:        mealPlan,
				Currency:        r.Get("currency").String(),
				Price:           price,
				Rooms:           roomList,
				Number:          i + 1,
				OldPrice:        oldPrice,
				DiscountPercent: discountPercent,
				DiscountPrice:   discountPrice,
			})
		}
		return rooms, nil
	}
	return []dto.RoomOptionDto{}, nil
}

func (p *hotelProvider) HotelAvailable(data dto.AvailableDto) (*dto.AvailableResponseDto, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", p.baseOrderUrl+AvailableEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	token, err := Token()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("ab-channel", common.ABChannelName)
	logger.WithName(logtags.HotelAvailableRequest).WithData(data).Info("hotel available request log")
	req.Close = true

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	err = httphelper.GetResponseError(res)
	if err != nil {
		return nil, err
	}

	jsonData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	jsonString := string(jsonData)

	return &dto.AvailableResponseDto{
		OrderId:      gjson.Get(jsonString, "result.details.orderId").String(),
		TotalPrice:   gjson.Get(jsonString, "result.totalPrice").Int(),
		Status:       gjson.Get(jsonString, "result.details.status").String(),
		IndraOrderId: gjson.Get(jsonString, "result.id").Int(),
		CheckIn:      gjson.Get(jsonString, "result.details.detail.checkIn").String(),
		CheckOut:     gjson.Get(jsonString, "result.details.detail.checkOut").String(),
		OptionId:     data.OptionId,
		Error:        nil,
	}, nil
}

func (p *hotelProvider) GetOrderDetail(hotelId, sessionId, optionId string) (*dto.OrderDetailDto, error) {
	info, err := p.GetHotelOptionInfo(dto.OptionInfoRequestDto{
		OptionId:  optionId,
		HotelId:   hotelId,
		SessionId: sessionId,
	})
	if err != nil {
		return nil, err
	}
	roomsList, err := p.GetHotelRoomsWithSession(dto.HotelRoomsWithSessionDto{
		HotelId:   hotelId,
		SessionId: sessionId,
	})
	if err != nil {
		return nil, err
	}
	var option dto.RoomOptionDto
	for _, room := range roomsList.Rooms {
		if room.Id == optionId {
			option = room
		}
	}
	rooms := make([]dto.OrderDetailRoomDto, 0)

	for _, room := range option.Rooms {
		rooms = append(rooms, dto.OrderDetailRoomDto{
			OrderRoomId:   0,
			OrderId:       0,
			PricePerNight: room.PricePerNight,
			Price:         room.Price,
			Name:          room.Name,
			NameEn:        room.NameEn,
		})
	}

	return &dto.OrderDetailDto{
		OrderId:                0,
		ProviderOrderId:        "",
		HotelID:                0,
		ProviderHotelId:        hotelId,
		NonRefundable:          info.Policy.NonRefundable,
		GeneralPolicies:        info.Policy.GeneralPolicies,
		TotalPrice:             option.Price,
		Provider:               option.Provider,
		ProviderName:           option.ProviderName,
		Currency:               option.Currency,
		MealPlan:               option.MealPlan.Key,
		RestrictedMarkupAmount: info.Detail.RestrictedMarkup.Amount,
		RestrictedMarkupType:   info.Detail.RestrictedMarkup.Type,
		Status:                 "Draft",
		Rooms:                  rooms,
	}, nil
}

func (p *hotelProvider) GetRoomCancellationPolicy(hotelId, roomId, sessionId string) (*dto.RoomCancellationPolicyDto, error) {
	body, err := json.Marshal(map[string]string{
		"hotelId":   hotelId,
		"optionId":  roomId,
		"sessionId": sessionId,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", p.baseOrderUrl+RoomCancellationPolicy, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("ab-channel", common.ABChannelName)
	req.Close = true
	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	err = httphelper.GetResponseError(res)
	if err != nil {
		logger.WithName(logtags.GettingCancellationPolicyError).ErrorException(err, "error while trying to get cancellation policy")
		return nil, common.ProviderUnknownProblem
	}
	jsonData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	jsonString := string(jsonData)
	if !gjson.Valid(jsonString) {
		logger.WithName(logtags.InvalidJsonResponseError).WithData(map[string]interface{}{
			"response": jsonString,
		}).Error(common.JsonDataIsNotValid.Error())
		return nil, common.JsonDataIsNotValid
	}
	policiesArray := gjson.Get(jsonString, "policy.general.policies").Array()
	policies := make([]string, 0, len(policiesArray))
	for _, policy := range policiesArray {
		policies = append(policies, policy.String())
	}

	nonRefundable := gjson.Get(jsonString, "policy.nonRefundable").Bool()
	if nonRefundable {
		policies = []string{common.NonRefundableDefaultMessage}
	}
	return &dto.RoomCancellationPolicyDto{
		NonRefundable:   nonRefundable,
		GeneralPolicies: policies,
		Error:           nil,
	}, nil
}

func (p *hotelProvider) ConfirmOrder(orderId string) (dto.ConfirmResponseDto, error) {
	req, err := http.NewRequest("POST", p.baseOrderUrl+
		strings.Replace(ConfirmOrderEndPoint, "{orderId}", orderId, 1), nil)
	if err != nil {
		return dto.ConfirmResponseDto{}, err
	}
	req.Header.Add("ab-channel", common.ABChannelName)
	logger.WithName(logtags.ConfirmOrderRequest).WithData(map[string]interface{}{
		"orderId": orderId,
	}).Info("Confirm order request log")
	req.Close = true

	stringData, err := p.requestToUrl(req, true)
	if err != nil {
		return dto.ConfirmResponseDto{}, err
	}
	if gjson.Get(stringData, "result").Bool() {
		return dto.ConfirmResponseDto{
			OrderId: orderId,
			Error:   nil,
		}, nil
	}
	return dto.ConfirmResponseDto{}, common.ErrorInConfirmingOrder
}

func (p *hotelProvider) PayByAccount(orderId string) (dto.OrderPayByAccountResponseDto, error) {
	data := dtos.NewPayByAccountRequest("")
	req, err := http.NewRequest("POST", p.baseOrderUrl+
		strings.Replace(PayByBankAndAccountEndpoint, "{orderId}", orderId, 1), bytes.NewBuffer(data.ToJson()))
	if err != nil {
		return dto.OrderPayByAccountResponseDto{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("ab-channel", common.ABChannelName)

	logger.WithName(logtags.PayByAccountRequest).WithData(map[string]interface{}{
		"orderId": orderId,
		"reqBody": data,
	}).Info("Pay by account request log")
	req.Close = true

	stringData, err := p.requestToUrl(req, true)
	if err != nil {
		return dto.OrderPayByAccountResponseDto{}, err
	}
	result := gjson.Get(stringData, "result")
	transactionIds := make([]string, 0)
	for _, item := range result.Get("transactionIds").Array() {
		transactionIds = append(transactionIds, item.String())
	}

	return dto.OrderPayByAccountResponseDto{
		TransactionStatus: result.Get("transactionStatus").String(),
		RequestId:         result.Get("requestId").String(),
		TransactionIds:    transactionIds,
		ResultMessage:     result.Get("resultMessage").String(),
		Error:             nil,
	}, nil
}

func (p *hotelProvider) GetOrderStatus(orderId string) (dto.OrderStatusResponseDto, error) {
	stringData, err := p.getFromUrl(p.baseOrderUrl+
		strings.Replace(GetOrderStatusEndpoint, "{orderId}", orderId, 1),
		true)
	if err != nil {
		return dto.OrderStatusResponseDto{}, err
	}
	status := gjson.Get(stringData, "result").String()
	return dto.OrderStatusResponseDto{
		OrderId: orderId,
		Status:  status,
		Error:   nil,
	}, nil
}

func (p *hotelProvider) GetOrderEnquiry(orderId, providerId string) (dto.OrderEnquiryResponseDto, error) {
	parameters := "?providerId=" + providerId

	logger.WithName(logtags.GetOrderEnquiryRequest).WithData(map[string]interface{}{
		"orderId": orderId,
	}).Info("Get order enquiry request log")

	stringData, err := p.getFromUrl(p.baseOrderUrl+
		strings.Replace(EnquiryOrderEndpoint, "{orderId}", orderId, 1)+
		parameters, true)
	if err != nil {
		return dto.OrderEnquiryResponseDto{}, err
	}
	result := gjson.Get(stringData, "result")
	allowedRefundPaymentMethods := make([]string, 0)
	for _, item := range result.Get("allowedRefundPaymentMethods").Array() {
		allowedRefundPaymentMethods = append(allowedRefundPaymentMethods, item.String())
	}
	//amenities
	var items []dto.OrderEnquiryItemDto
	for _, f := range result.Get("items").Array() {
		var itemOptions []dto.OrderEnquiryItemOptionDto
		for _, o := range f.Get("items").Array() {
			itemOptions = append(itemOptions, dto.OrderEnquiryItemOptionDto{
				ReferenceCode:      o.Get("referenceCode").String(),
				IsRefundable:       o.Get("isRefundable").Bool(),
				PaidAmount:         o.Get("paidAmount").Int(),
				TotalPenaltyAmount: o.Get("totalPenaltyAmount").Int(),
				RefundableAmount:   o.Get("refundableAmount").Int(),
				RefundableType:     o.Get("refundableType").String(),
				RefundableStatus:   o.Get("refundableStatus").String(),
				RefundStatus:       o.Get("refundStatus").String(),
				PassengerInformation: dto.OrderEnquiryItemOptionPassengerInformation{
					Title:           o.Get("passengerInformation.title").String(),
					Name:            o.Get("passengerInformation.name").String(),
					LastName:        o.Get("passengerInformation.lastName").String(),
					NamePersian:     o.Get("passengerInformation.namePersian").String(),
					LastNamePersian: o.Get("passengerInformation.lastNamePersian").String(),
				},
			})
		}
		items = append(items, dto.OrderEnquiryItemDto{
			ProviderId:          f.Get("providerId").String(),
			ProductProviderType: f.Get("productProviderType").String(),
			Destination:         f.Get("destination").String(),
			DestinationName:     f.Get("destinationName").String(),
			Items:               itemOptions,
		})
	}
	return dto.OrderEnquiryResponseDto{
		AllowedRefundPaymentMethods: allowedRefundPaymentMethods,
		Items:                       items,
		Error:                       nil,
	}, nil
}

func (p *hotelProvider) GetOrdersRefundStatus(ids []string, size int, page int) (*dto.OrdersRefundStatusResponseDto, error) {
	parameters := fmt.Sprintf("?page_no=%d&page_size=%d", page, size)
	for _, id := range ids {
		parameters += "&orderIds=" + id
	}
	req, err := http.NewRequest("GET", p.baseOrderServiceUrl+OrdersRefundStatusEndpoint+parameters, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("ab-channel", common.ABChannelName)

	logger.WithName(logtags.GetOrdersRefundStatusRequest).WithData(map[string]interface{}{
		"ids": ids,
	}).Info("Get orders refund status request log")
	req.Close = true

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	err = httphelper.GetResponseError(res)
	if err != nil {
		return nil, err
	}

	var result dto.OrdersRefundStatusResponseDto
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *hotelProvider) RefundOrder(orderId, referenceCode string) (dto.OrderRefundResponseDto, error) {
	body := dtos.NewRefundRequestDto(orderId, referenceCode)
	req, err := http.NewRequest("POST", p.baseOrderUrl+OrderRefundEndpoint, bytes.NewBuffer(body.ToJson()))
	if err != nil {
		return dto.OrderRefundResponseDto{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("ab-channel", common.ABChannelName)

	logger.WithName(logtags.RefundOrderRequest).WithData(map[string]interface{}{
		"orderId": orderId,
		"reqBody": body,
	}).Info("Refund order request log")
	req.Close = true

	stringData, err := p.requestToUrl(req, true)
	if err != nil {
		return dto.OrderRefundResponseDto{}, err
	}

	result := gjson.Get(stringData, "result")

	return dto.OrderRefundResponseDto{
		OrderId:         orderId,
		RefundRequestId: result.Get("refundRequestId").Int(),
		Error:           nil,
	}, nil
}

func (p *hotelProvider) getFromUrl(url string, needAuthentication bool) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("ab-channel", common.ABChannelName)
	req.Close = true

	return p.requestToUrl(req, needAuthentication)
}

func (p *hotelProvider) requestToUrl(req *http.Request, needAuthentication bool) (string, error) {
	if needAuthentication {
		token, err := Token()
		if err != nil {
			return "", err
		}
		req.Header.Add("Authorization", "Bearer "+token)
	}
	res, err := p.client.Do(req)
	if err != nil {
		return "", err
	}

	err = httphelper.GetResponseError(res)
	if err != nil {
		return "", err
	}

	jsonData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	jsonString := string(jsonData)

	if !gjson.Valid(jsonString) {
		logger.WithName(logtags.InvalidJsonResponseError).WithData(map[string]interface{}{
			"response": jsonString,
		}).Error(common.JsonDataIsNotValid.Error())
	}
	return jsonString, nil
}

func NewHotelProvider() core.HotelProvider {
	conf := config.Get()
	provider := &hotelProvider{
		password:            conf.ProviderPassword,
		username:            conf.ProviderUsername,
		baseHotelUrl:        conf.ProviderEndpoint,
		baseOrderUrl:        conf.OrderEndpoint,
		baseOrderServiceUrl: conf.OrderServiceEndpoint,
		client:              &http.Client{},
	}
	return provider
}
