package core

import (
	"hotel-engine/core/dbmodel"
	"hotel-engine/core/dto"
)

type Mapper interface {
	ToHotelModel(dto dto.HotelDto) *dbmodel.Hotel
	ToHotelDto(model dbmodel.Hotel) *dto.HotelDto
	ToHotelsDto(model []dbmodel.Hotel) []dto.HotelDto
	ToHotelPDPDto(hotel dto.HotelDto, rooms []dto.RequestRoomDto) *dto.HotelPDPDto

	ToCitiesModel(dtos []dto.CityDto) []dbmodel.City
	ToCitiesDto(models []dbmodel.City) []dto.CityDto
	ToCityModel(dto dto.CityDto) *dbmodel.City
	ToCityDto(model dbmodel.City) *dto.CityDto

	ToAmenitiesDto(models []dbmodel.Amenity) []dto.HotelAmenityDto
	ToAmenityModel(dto dto.HotelAmenityDto) *dbmodel.Amenity
	ToAmenityDto(model dbmodel.Amenity) *dto.HotelAmenityDto

	ToPlaceDto(model dbmodel.Place) *dto.HotelPlaceDto
	ToPlaceModel(dto dto.HotelPlaceDto) *dbmodel.Place
	ToPlacesDto(models []dbmodel.Place) []dto.HotelPlaceDto

	ToOrderRoomModel(item dto.OrderDetailRoomDto) dbmodel.OrderRoom
	ToOrderRoomDto(model dbmodel.OrderRoom) dto.OrderDetailRoomDto
	ToOrderModel(item dto.OrderDetailDto) dbmodel.Order
	ToOrderDto(model dbmodel.Order) dto.OrderDetailDto

	ToHotelsDetail(hotels []dbmodel.Hotel) *dto.SyncedHotelsDetail
	ToHotelSyncDetail(hotel dbmodel.Hotel) dto.HotelSyncDetail

	ToAmenityCategoryDto(model *dbmodel.AmenityCategory) *dto.AmenityCategoryDto
	ToAmenityCategoriesDto(items []*dbmodel.AmenityCategory) []*dto.AmenityCategoryDto
	ToAmenityCategoryModel(dto dto.AmenityCategoryDto) dbmodel.AmenityCategory

	ToBadgeModel(dto dto.HotelBadgeDto) *dbmodel.Badge
	ToBadgeDto(model dbmodel.Badge) *dto.HotelBadgeDto
	ToBadgesDto(models []dbmodel.Badge) []dto.HotelBadgeDto

	ToFAQModel(faq dto.HotelFAQDto) *dbmodel.FAQ
	ToFAQDto(model dbmodel.FAQ) *dto.HotelFAQDto
	ToFAQsModel(faqs []dto.HotelFAQDto) []*dbmodel.FAQ
	ToFAQsDto(models []*dbmodel.FAQ) []dto.HotelFAQDto
}
