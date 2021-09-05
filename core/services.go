package core

import (
	"hotel-engine/core/dbmodel"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/hotelproviderinterface/dtos"
	"time"
)

type HotelService interface {
	FindHotelById(id string) (*dto.HotelDto, error)
	GetRoomCancellationPolicy(hotelId, roomId, sessionId string) (*dto.RoomCancellationPolicyDto, error)
	SearchResult(request dto.SearchDto) (*dto.SearchResponseDto, error)

	UpdateSomeHotels(hotelsDto dto.SyncSomeHotelsDto, date time.Time) (*dto.UpdateResultDto, error)
	UpdateAllSync(date time.Time) (*dto.UpdateResultDto, error)
	UpdateAll(date time.Time) (*dto.TaskRunningResult, error)

	SyncAllHotels() (*dto.TaskRunningResult, error)
	SyncSomeHotels(hotelsDto dto.SyncSomeHotelsDto) (*dto.UpdateResultDto, error)
	SyncedHotels() (*dto.SyncedHotelsDetail, error)
	GetHotelRoomsWithSession(dto dto.HotelRoomsWithSessionDto) (*dto.RateRoomResponseDto, error)
	GetHotelRooms(dto dto.HotelRoomsDto) (*dto.RateRoomResponseDto, error)
	HotelAvailable(dto dto.AvailableDto) (*dto.AvailableResponseDto, error)
	FinalizeHotelOrder(dto dto.FinalizeOrderDto) (*dto.FinalizeOrderResponseDto, error)
	GetHotelDetails(request dto.HotelDetailsDto) (*dto.HotelPDPDto, error)
	SetAmenityIcon(request dto.SetAmenityIconDto) (dto.HotelAmenityDto, error)
	GetHotelOptionInfo(infoDto dto.OptionInfoRequestDto) (*dto.OptionInfoResponseDto, error)
	GetHotels(ids []string) ([]dto.HotelDto, error)
	GetAnOrderDetail(orderId string) (*dto.OrderDetailDto, error)

	ConfirmOrder(orderId string) (dto.ConfirmResponseDto, error)
	PayByAccount(orderId string) (dto.OrderPayByAccountResponseDto, error)
	GetOrderStatus(orderId string) (dto.OrderStatusResponseDto, error)
	GetOrderEnquiry(orderId string) (dto.OrderEnquiryResponseDto, error)
	RefundOrder(refundRequest dto.OrderRefundRequestDto) (dto.OrderRefundResponseDto, error)
	UpdateHotelRateReview(rateDto dto.RateReviewEventDto) error
	UpdateRefundedOrdersPaymentStatus(date time.Time)
	SetAmenityCategory(body dto.SetAmenityCategoryDto) (dto.HotelAmenityDto, error)

	GetHotelsList(body dto.HotelsPageRequestDto) (dto.HotelsPageResponseDto, error)
	SetHotelSeoDetails(body dto.SetHotelSeoRequestDto) (dto.HotelDto, error)
	SetHotelFaq(requestDto dto.SetHotelFaqRequestDto) (dto.HotelDto, error)
	RemoveHotelFaq(hotelId string, faqId uint) (dto.HotelDto, error)
}

type PublicService interface {
	GetAllCities() []dto.CityDto
	GetAllBadges() []dto.HotelBadgeDto
	GetAllAmenities() []dto.HotelAmenityDto
	GetAllSorts() []dto.SortDto
	GetAllFilters() *dto.FiltersDto
	SyncAllCities() ([]dto.CityDto, error)
	GetAllPlaces() []dto.HotelPlaceDto
	GetChildAgeRanges() []dto.ChildAgeRangeDto

	CreateAmenityCategory(dto dto.AmenityCategoryDto) (dto.AmenityCategoryDto, error)
	GetAmenityCategory(id uint) (*dto.AmenityCategoryDto, error)
	DeleteAmenityCategory(id int) error
	GetAmenityCategories() []*dto.AmenityCategoryDto
	UpdateAmenityCategory(dto dto.AmenityCategoryDto) (dto.AmenityCategoryDto, error)

	SetBadgeIcon(request dto.SetBadgeIconDto) (dto.HotelBadgeDto, error)
}

type SyncService interface {
	UpdateDatabase()
	SyncElastic()
	UpdateAndSyncElastic()
	HasBeenSynced() (bool, error)
}

type CityCacheStore interface {
	GetAll() []dto.CityDto
	FindOne(name string) (dto.CityDto, error)
	Find(name string) []dto.CityDto
	Get(id string) (dto.CityDto, error)
}

type AmenityCacheStore interface {
	GetAll() []dto.HotelAmenityDto
	FindOne(name string) (dto.HotelAmenityDto, error)
	FindByIds(ids []int) []dto.HotelAmenityDto
	Find(name string) []dto.HotelAmenityDto
	Get(id int) (dto.HotelAmenityDto, error)
}

type CacheStore interface {
	AmenityStore() AmenityCacheStore
	CityStore() CityCacheStore
	UpdateStores() error
	UpdateCityStore()
	UpdateAmenityStore()
}

type SearchDtoAdopter interface {
	CreateProviderDto(searchDto dto.SearchDto) (*dtos.ProviderSearchDto, error)
	CreateResultDto(results []dtos.Result, hotels []dbmodel.Hotel,
		cityId string, days int) (*dto.SearchResponseDto, error)
}

type HotelProvider interface {
	GetAccessToken() (string, error)
	GetRoomCancellationPolicy(hotelId, roomId, sessionId string) (*dto.RoomCancellationPolicyDto, error)
	CreateHotelEnumerable([]dto.CityDto) HotelEnumerable
	GetHotelData(hotelId string, checkIn, checkout time.Time) (*dto.HotelDto, error)
	GetHotels(limit, skip uint32) (*dtos.HotelsListResult, error)
	SearchHotels(limit, skip uint32, hotelGiataId, cityId, hotelName string, cityBaseId int64) (*dtos.HotelsListResult, error)
	SearchResult(searchDto dtos.ProviderSearchDto) ([]dtos.Result, int, error)
	DirectHotel(dto dto.HotelRoomsDto) (*dto.DirectResponseDto, error)
	GetHotelRooms(dto dto.HotelRoomsDto) (*dto.RateRoomResponseDto, error)
	GetHotelRoomsWithSession(dto dto.HotelRoomsWithSessionDto) (*dto.RateRoomResponseDto, error)
	HotelAvailable(dto dto.AvailableDto) (*dto.AvailableResponseDto, error)
	GetOrderDetail(hotelId, sessionId, optionId string) (*dto.OrderDetailDto, error)
	GetHotelOptionInfo(infoDto dto.OptionInfoRequestDto) (*dto.OptionInfoResponseDto, error)

	ConfirmOrder(orderId string) (dto.ConfirmResponseDto, error)
	PayByAccount(orderId string) (dto.OrderPayByAccountResponseDto, error)
	GetOrderStatus(orderId string) (dto.OrderStatusResponseDto, error)
	GetOrderEnquiry(orderId, providerId string) (dto.OrderEnquiryResponseDto, error)
	RefundOrder(orderId, referenceCode string) (dto.OrderRefundResponseDto, error)
	GetHotelType(hotelId string) (string, error)
	GetOrdersRefundStatus(ids []string, size int, page int) (*dto.OrdersRefundStatusResponseDto, error)
}

type BasicInformationProvider interface {
	GetCities() ([]dto.CityDto, error)
}

type HotelEnumerable interface {
	MoveNext() bool
	Current() ([]dto.HotelIdentifierDto, error)
	Reset()
}

type ProviderBalanceChecker interface {
	GetBalance() (float64, error)
	CheckAdequateBalance()
}

type BalanceAlertNotifier interface {
	Notify(alert dto.BalanceNotifyDto)
}

type DistributedLocker interface {
	Lock(key string, duration time.Duration, toDo func()) error
}

type OrderEventDispatcher interface {
	OrderRefundRequestFinalized(event dto.OrderRefundRequestFinalizedDto)
}
