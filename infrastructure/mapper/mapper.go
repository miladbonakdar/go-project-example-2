package mapper

import (
	"github.com/jinzhu/gorm"
	"hotel-engine/core"
	"hotel-engine/core/dbmodel"
	"hotel-engine/core/dto"
	"hotel-engine/utils/random"
	"strings"
)

type mapper struct{}

func (m *mapper) ToHotelModel(dto dto.HotelDto) *dbmodel.Hotel {

	amenities := make([]*dbmodel.Amenity, 0, len(dto.Amenities))
	for _, amenity := range dto.Amenities {
		amenities = append(amenities, m.ToAmenityModel(amenity))
	}

	places := make([]*dbmodel.Place, 0, len(dto.Places))
	for _, place := range dto.Places {
		places = append(places, m.ToPlaceModel(place))
	}

	badges := make([]*dbmodel.Badge, 0, len(dto.Badges))
	for _, badge := range dto.Badges {
		badges = append(badges, m.ToBadgeModel(badge))
	}

	return &dbmodel.Hotel{
		PlaceID:            dto.PlaceID,
		RoomID:             dto.RoomID,
		HotelCode:          random.GenerateHotelCode(),
		Code:               random.GenerateHotelNumericCode(),
		Type:               dto.Type,
		Kind:               dto.Kind,
		MinNight:           dto.MinNight,
		ReservationType:    dto.ReservationType,
		PaymentType:        dto.PaymentType,
		Name:               dto.Name,
		NameEn:             strings.ToLower(dto.NameEn),
		Description:        dto.Description,
		City:               dto.City,
		Province:           dto.Province,
		GeoLocation:        dto.GeoLocation,
		Region:             dto.Region,
		SuitableFor:        strings.Join(dto.SuitableFor, ","),
		Images:             strings.Join(dto.Images, ","),
		CheckIn:            dto.CheckIn,
		CheckOut:           dto.CheckOut,
		Price:              dto.Price,
		Capacity:           dto.Capacity,
		Amenities:          amenities,
		Places:             places,
		Tags:               strings.Join(dto.Tags, ","),
		Verified:           dto.Verified,
		RateReviewScore:    dto.RateReviewScore,
		RateReviewCount:    dto.RateReviewCount,
		CheckOutTime:       dto.CheckOutTime,
		CheckInTime:        dto.CheckInTime,
		Address:            dto.Address,
		Star:               dto.Star,
		CityEn:             dto.CityEn,
		ProvinceEn:         dto.ProvinceEn,
		Country:            dto.Country,
		CountryCode:        dto.CountryCode,
		CountryEn:          dto.CountryEn,
		Badges:             badges,
		DiscountPercent:    dto.DiscountPercent,
		DiscountPrice:      dto.DiscountPrice,
		OldPrice:           dto.OldPrice,
		Sort:               dto.Sort,
		SeoRobots:          dto.Seo.Robots,
		SeoH1:              dto.Seo.H1,
		SeoDescription:     dto.Seo.Description,
		SeoTitle:           dto.Seo.Title,
		SeoCanonical:       dto.Seo.Canonical,
		SeoMetaDescription: dto.Seo.MetaDescription,
	}
}

func (m *mapper) ToHotelDto(model dbmodel.Hotel) *dto.HotelDto {

	amenities := make([]dto.HotelAmenityDto, 0, len(model.Amenities))
	for _, amenity := range model.Amenities {
		amenities = append(amenities, *m.ToAmenityDto(*amenity))
	}

	places := make([]dto.HotelPlaceDto, 0, len(model.Places))
	for _, place := range model.Places {
		places = append(places, *m.ToPlaceDto(*place))
	}

	badges := make([]dto.HotelBadgeDto, 0, len(model.Badges))
	for _, badge := range model.Badges {
		badges = append(badges, *m.ToBadgeDto(*badge))
	}

	return &dto.HotelDto{
		PlaceID:         model.PlaceID,
		RoomID:          model.RoomID,
		Type:            model.Type,
		Kind:            model.Kind,
		MinNight:        model.MinNight,
		ReservationType: model.ReservationType,
		PaymentType:     model.PaymentType,
		Name:            model.Name,
		Description:     model.Description,
		Code:            model.Code,
		Region:          model.Region,
		SuitableFor:     strings.Split(model.SuitableFor, ","),
		Images:          strings.Split(model.Images, ","),
		Amenities:       amenities,
		Places:          places,
		Tags:            strings.Split(model.Tags, ","),
		Verified:        model.Verified,
		RateReviewCount: model.RateReviewCount,
		RateReviewScore: model.RateReviewScore,
		Capacity:        model.Capacity,
		CheckIn:         model.CheckIn,
		CheckOut:        model.CheckOut,
		City:            model.City,
		GeoLocation:     model.GeoLocation,
		Price:           model.Price,
		Province:        model.Province,
		CheckOutTime:    model.CheckOutTime,
		CheckInTime:     model.CheckInTime,
		Address:         model.Address,
		Star:            model.Star,
		NameEn:          model.NameEn,
		CityEn:          model.CityEn,
		ProvinceEn:      model.ProvinceEn,
		Country:         model.Country,
		CountryCode:     model.CountryCode,
		CountryEn:       model.CountryEn,
		UpdatedAt:       &(model.UpdatedAt),
		CreatedAt:       &(model.CreatedAt),
		DiscountPercent: model.DiscountPercent,
		DiscountPrice:   model.DiscountPrice,
		OldPrice:        model.OldPrice,
		Badges:          badges,
		Sort:            model.Sort,
		Seo: dto.HotelSeoDto{
			Title:           model.SeoTitle,
			H1:              model.SeoH1,
			Description:     model.SeoDescription,
			Robots:          model.SeoRobots,
			Canonical:       model.SeoCanonical,
			MetaDescription: model.SeoMetaDescription,
		},
		FAQ: dto.HotelFAQDetailsDto{
			Title:   model.FAQTitle,
			FAQList: m.ToFAQsDto(model.FAQList),
		},
	}
}

func (m *mapper) ToHotelPDPDto(hotel dto.HotelDto, rooms []dto.RequestRoomDto) *dto.HotelPDPDto {
	return &dto.HotelPDPDto{
		HotelDto: dto.HotelDto{
			PlaceID:              hotel.PlaceID,
			RoomID:               hotel.RoomID,
			Type:                 hotel.Type,
			Kind:                 hotel.Kind,
			MinNight:             hotel.MinNight,
			ReservationType:      hotel.ReservationType,
			PaymentType:          hotel.PaymentType,
			Name:                 hotel.Name,
			Description:          hotel.Description,
			Region:               hotel.Region,
			SuitableFor:          hotel.SuitableFor,
			Images:               hotel.Images,
			Amenities:            hotel.Amenities,
			Places:               hotel.Places,
			Tags:                 hotel.Tags,
			Verified:             hotel.Verified,
			RateReviewCount:      hotel.RateReviewCount,
			RateReviewScore:      hotel.RateReviewScore,
			Capacity:             hotel.Capacity,
			CheckIn:              hotel.CheckIn,
			CheckOut:             hotel.CheckOut,
			City:                 hotel.City,
			GeoLocation:          hotel.GeoLocation,
			Price:                hotel.Price,
			Province:             hotel.Province,
			CheckOutTime:         hotel.CheckOutTime,
			CheckInTime:          hotel.CheckInTime,
			Address:              hotel.Address,
			Star:                 hotel.Star,
			NameEn:               hotel.NameEn,
			CityEn:               hotel.CityEn,
			ProvinceEn:           hotel.ProvinceEn,
			Country:              hotel.Country,
			CountryCode:          hotel.CountryCode,
			CountryEn:            hotel.CountryEn,
			UpdatedAt:            hotel.UpdatedAt,
			CreatedAt:            hotel.CreatedAt,
			DiscountPercent:      hotel.DiscountPercent,
			DiscountPrice:        hotel.DiscountPrice,
			OldPrice:             hotel.OldPrice,
			Badges:               hotel.Badges,
			Sort:                 hotel.Sort,
			Seo:                  hotel.Seo,
			Rooms:                hotel.Rooms,
			UnavailableAmenities: hotel.UnavailableAmenities,
			SessionId:            hotel.SessionId,
			FAQ:                  hotel.FAQ,
			Code:                 hotel.Code,
		},
		Rooms: rooms,
	}
}

func (m *mapper) ToHotelsDto(hotels []dbmodel.Hotel) []dto.HotelDto {
	hotelsDto := make([]dto.HotelDto, 0, len(hotels))
	for _, hotel := range hotels {
		hotelsDto = append(hotelsDto, *m.ToHotelDto(hotel))
	}
	return hotelsDto
}

func (m *mapper) ToCitiesModel(dtos []dto.CityDto) []dbmodel.City {
	cities := make([]dbmodel.City, 0, len(dtos))
	for _, city := range dtos {
		cities = append(cities, *m.ToCityModel(city))
	}
	return cities
}

func (m *mapper) ToCitiesDto(models []dbmodel.City) []dto.CityDto {
	cities := make([]dto.CityDto, 0, len(models))
	for _, city := range models {
		cities = append(cities, *m.ToCityDto(city))
	}
	return cities
}

func (m *mapper) ToCityModel(dto dto.CityDto) *dbmodel.City {
	return &dbmodel.City{
		CityID:  dto.Id,
		Name:    dto.Name,
		Country: dto.Country,
		State:   dto.State,
		BaseId:  dto.BaseId,
	}
}

func (m *mapper) ToCityDto(model dbmodel.City) *dto.CityDto {
	return &dto.CityDto{
		Id:      model.CityID,
		Name:    model.Name,
		Country: model.Country,
		State:   model.State,
		BaseId:  model.BaseId,
	}
}

func (m *mapper) ToAmenitiesDto(models []dbmodel.Amenity) []dto.HotelAmenityDto {
	amenities := make([]dto.HotelAmenityDto, 0, len(models))
	for _, amenity := range models {
		amenities = append(amenities, *m.ToAmenityDto(amenity))

	}
	return amenities
}

func (m *mapper) ToAmenityModel(dto dto.HotelAmenityDto) *dbmodel.Amenity {
	return &dbmodel.Amenity{
		ID:                dto.ID,
		Name:              dto.Name,
		GroupId:           dto.GroupId,
		IconUrl:           dto.IconUrl,
		NameEn:            dto.NameEn,
		AmenityCategoryID: dto.AmenityCategoryID,
	}
}

func (m *mapper) ToAmenityDto(model dbmodel.Amenity) *dto.HotelAmenityDto {
	cat := m.ToAmenityCategoryDto(model.AmenityCategory)
	return &dto.HotelAmenityDto{
		ID:                model.ID,
		Name:              model.Name,
		GroupId:           model.GroupId,
		IconUrl:           model.IconUrl,
		NameEn:            model.NameEn,
		AmenityCategoryID: model.AmenityCategoryID,
		AmenityCategory:   cat,
	}
}

func (m *mapper) ToPlaceModel(dto dto.HotelPlaceDto) *dbmodel.Place {
	return &dbmodel.Place{
		ID:          dto.ID,
		Name:        dto.Name,
		GeoLocation: dto.GeoLocation,
		Distance:    dto.Distance,
		Priority:    dto.Priority,
	}
}

func (m *mapper) ToPlaceDto(model dbmodel.Place) *dto.HotelPlaceDto {
	return &dto.HotelPlaceDto{
		ID:          model.ID,
		Name:        model.Name,
		GeoLocation: model.GeoLocation,
		Distance:    model.Distance,
		Priority:    model.Priority,
	}
}

func (m *mapper) ToPlacesDto(models []dbmodel.Place) []dto.HotelPlaceDto {
	places := make([]dto.HotelPlaceDto, 0, len(models))
	for _, place := range models {
		places = append(places, *m.ToPlaceDto(place))
	}
	return places
}

func (m *mapper) ToOrderDto(model dbmodel.Order) dto.OrderDetailDto {
	rooms := make([]dto.OrderDetailRoomDto, 0)
	for _, room := range model.Rooms {
		rooms = append(rooms, m.ToOrderRoomDto(room))
	}
	return dto.OrderDetailDto{
		OrderId:                  model.ID,
		ProviderOrderId:          model.ProviderOrderId,
		IndraOrderId:             model.IndraOrderId,
		HotelID:                  model.HotelID,
		ProviderHotelId:          model.ProviderHotelId,
		NonRefundable:            model.NonRefundable,
		GeneralPolicies:          strings.Split(model.GeneralPolicies, ","),
		TotalPrice:               model.TotalPrice,
		Provider:                 model.Provider,
		ProviderName:             model.ProviderName,
		Currency:                 model.Currency,
		MealPlan:                 model.MealPlan,
		RestrictedMarkupAmount:   model.RestrictedMarkupAmount,
		RestrictedMarkupType:     model.RestrictedMarkupType,
		Status:                   model.Status,
		Rooms:                    rooms,
		ApplicantRefundRequestId: model.ApplicantRefundRequestId,
		ApplicantOrderId:         model.ApplicantOrderId,
		PaidAmount:               model.PaidAmount,
		ReferenceCode:            model.ReferenceCode,
		RefundStatus:             model.RefundStatus,
		RefundableAmount:         model.RefundableAmount,
		TotalPenaltyAmount:       model.TotalPenaltyAmount,
		RefundRequestId:          model.RefundRequestId,
		Confirmed:                model.Confirmed,
	}
}

func (m *mapper) ToOrderModel(item dto.OrderDetailDto) dbmodel.Order {
	rooms := make([]dbmodel.OrderRoom, 0)
	for _, room := range item.Rooms {
		rooms = append(rooms, m.ToOrderRoomModel(room))
	}
	return dbmodel.Order{
		Model:                  gorm.Model{},
		ProviderOrderId:        item.ProviderOrderId,
		IndraOrderId:           item.IndraOrderId,
		HotelID:                item.HotelID,
		ProviderHotelId:        item.ProviderHotelId,
		NonRefundable:          item.NonRefundable,
		GeneralPolicies:        strings.Join(item.GeneralPolicies, ","),
		TotalPrice:             item.TotalPrice,
		Provider:               item.Provider,
		ProviderName:           item.ProviderName,
		Currency:               item.Currency,
		MealPlan:               item.MealPlan,
		RestrictedMarkupAmount: item.RestrictedMarkupAmount,
		RestrictedMarkupType:   item.RestrictedMarkupType,
		Status:                 item.Status,
		Rooms:                  rooms,
	}
}

func (m *mapper) ToOrderRoomDto(model dbmodel.OrderRoom) dto.OrderDetailRoomDto {
	return dto.OrderDetailRoomDto{
		OrderRoomId:   model.ID,
		OrderId:       model.OrderID,
		PricePerNight: model.PricePerNight,
		Price:         model.Price,
		Name:          model.Name,
		NameEn:        model.NameEn,
	}
}

func (m *mapper) ToOrderRoomModel(item dto.OrderDetailRoomDto) dbmodel.OrderRoom {
	return dbmodel.OrderRoom{
		OrderID:       item.OrderId,
		PricePerNight: item.PricePerNight,
		Price:         item.Price,
		Name:          item.Name,
		NameEn:        item.NameEn,
	}
}

func (m *mapper) ToHotelsDetail(hotels []dbmodel.Hotel) *dto.SyncedHotelsDetail {
	details := &dto.SyncedHotelsDetail{}
	for _, h := range hotels {
		if h.Price == 0 {
			details.NotSyncedHotels = append(details.NotSyncedHotels, m.ToHotelSyncDetail(h))
			details.TotalNotSynced++
		} else {
			details.SyncedHotels = append(details.SyncedHotels, m.ToHotelSyncDetail(h))
			details.TotalSynced++
		}
	}
	details.TotalHotels = len(hotels)
	return details
}

func (m *mapper) ToHotelSyncDetail(hotel dbmodel.Hotel) dto.HotelSyncDetail {
	details := dto.HotelSyncDetail{
		PlaceID:   hotel.PlaceID,
		Price:     hotel.Price,
		CreatedAt: hotel.CreatedAt,
		UpdatedAt: hotel.UpdatedAt,
	}
	return details
}

func (m *mapper) ToAmenityCategoryDto(model *dbmodel.AmenityCategory) *dto.AmenityCategoryDto {
	if model == nil {
		return nil
	}
	cat := dto.AmenityCategoryDto{
		ID:      model.ID,
		Name:    model.Name,
		NameEn:  model.NameEn,
		Order:   model.Order,
		IconUrl: model.IconUrl,
		Error:   nil,
	}
	return &cat
}

func (m *mapper) ToAmenityCategoriesDto(items []*dbmodel.AmenityCategory) []*dto.AmenityCategoryDto {
	cats := make([]*dto.AmenityCategoryDto, 0)
	for _, item := range items {
		cats = append(cats, m.ToAmenityCategoryDto(item))
	}
	return cats
}

func (m *mapper) ToAmenityCategoryModel(dto dto.AmenityCategoryDto) dbmodel.AmenityCategory {
	cat := dbmodel.AmenityCategory{
		Name:    dto.Name,
		NameEn:  dto.NameEn,
		Order:   dto.Order,
		IconUrl: dto.IconUrl,
	}
	if dto.ID != 0 {
		cat.ID = 0
	}
	return cat
}

func (m *mapper) ToBadgeModel(dto dto.HotelBadgeDto) *dbmodel.Badge {
	return &dbmodel.Badge{
		ID:              dto.ID,
		Text:            dto.Text,
		Icon:            "",
		TextColor:       dto.TextColor,
		BackgroundColor: dto.BackgroundColor,
	}
}

func (m *mapper) ToBadgeDto(model dbmodel.Badge) *dto.HotelBadgeDto {
	return &dto.HotelBadgeDto{
		ID:              model.ID,
		Text:            model.Text,
		Icon:            model.Icon,
		TextColor:       model.TextColor,
		BackgroundColor: model.BackgroundColor,
	}
}

func (m *mapper) ToBadgesDto(models []dbmodel.Badge) []dto.HotelBadgeDto {
	badges := make([]dto.HotelBadgeDto, 0, len(models))
	for _, badge := range models {
		badges = append(badges, *m.ToBadgeDto(badge))
	}
	return badges
}

func (m *mapper) ToFAQModel(faq dto.HotelFAQDto) *dbmodel.FAQ {
	return &dbmodel.FAQ{
		ID:       faq.ID,
		Question: faq.Question,
		Answer:   faq.Answer,
	}
}

func (m *mapper) ToFAQDto(model dbmodel.FAQ) *dto.HotelFAQDto {
	return &dto.HotelFAQDto{
		ID:       model.ID,
		Question: model.Question,
		Answer:   model.Answer,
	}
}

func (m *mapper) ToFAQsModel(faqs []dto.HotelFAQDto) []*dbmodel.FAQ {
	faqsModel := make([]*dbmodel.FAQ, 0)
	for _, faq := range faqs {
		faqsModel = append(faqsModel, m.ToFAQModel(faq))
	}
	return faqsModel
}

func (m *mapper) ToFAQsDto(models []*dbmodel.FAQ) []dto.HotelFAQDto {
	faqsDto := make([]dto.HotelFAQDto, 0)
	for _, faq := range models {
		faqsDto = append(faqsDto, *m.ToFAQDto(*faq))
	}
	return faqsDto
}

func NewHotelMapper() core.Mapper {
	return &mapper{}
}
