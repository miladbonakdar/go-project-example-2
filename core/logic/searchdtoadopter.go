package logic

import (
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dbmodel"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/constants"
	"hotel-engine/infrastructure/hotelproviderinterface/dtos"
	"hotel-engine/infrastructure/logger"
	"math"
)

type searchDtoAdopter struct {
	cacheStore core.CacheStore
}

func (p *searchDtoAdopter) CreateProviderDto(searchDto dto.SearchDto) (*dtos.ProviderSearchDto, error) {
	city, err := p.cacheStore.CityStore().FindOne(searchDto.City)
	if err != nil {
		return nil, err
	}
	if searchDto.PageNumber <= 0 {
		searchDto.PageNumber = 1
	}
	sortOrder := -1
	sort := "score"
	if searchDto.SortDirection {
		sortOrder = 1
	}
	if searchDto.Sort != "" {
		sort = searchDto.Sort
	}
	providerDto := &dtos.ProviderSearchDto{
		RequestSession: dtos.ProviderRequestSessionDto{
			CheckIn:  searchDto.Date.Start,
			CheckOut: searchDto.Date.End,
			Rooms:    createRoomsData(searchDto),
			Destination: dtos.ProviderRequestSessionDestinationDto{
				Type: "City",
				Id:   city.Id,
			},
		},
		RequestSearchHotels: dtos.ProviderSearchHotelsRequestDto{
			SessionId: searchDto.SessionId,
			Limit:     searchDto.PageSize,
			Skip:      searchDto.PageSize * (searchDto.PageNumber - 1),
			Sort: dtos.ProviderSearchHotelsSortDto{
				Field: sort,
				Order: sortOrder,
			},
			Filters: createSearchFilters(searchDto),
		},
	}

	return providerDto, nil
}

func createSearchFilters(searchDto dto.SearchDto) []interface{} {
	filters := make([]interface{}, 0)
	if searchDto.Keyword != "" {
		filters = append(filters, dtos.SearchHotelsStringFilter{
			Field: "name",
			Value: []string{searchDto.Keyword},
		})
	}

	if searchDto.HotelTypes != nil {
		typesFilter := make([]int, 0)
		for _, t := range searchDto.HotelTypes {
			typesFilter = append(typesFilter, common.HotelTypes[t])
		}
		if len(typesFilter) > 0 {
			filters = append(filters, dtos.SearchHotelsIntFilter{
				Field: "accommodation",
				Value: typesFilter,
			})
		}
	}

	if searchDto.Score.Start != 0 || searchDto.Score.End != 0 {
		filters = append(filters, dtos.SearchHotelsFloatFilter{
			Field: "score",
			Value: []float32{float32(searchDto.Score.Start), float32(searchDto.Score.End)},
		})
	}

	if searchDto.Price.Start != 0 || searchDto.Price.End != 0 {
		filters = append(filters, dtos.SearchHotelsIntRangeFilter{
			Field: "price",
			Value: [][]int64{
				{searchDto.Price.Start, searchDto.Price.End},
			},
		})
	}

	if searchDto.Stars != nil && len(searchDto.Stars) > 0 {
		var Values []float32

		for _, t := range searchDto.Stars {
			switch t {
			case 1:
				Values = append(Values, 1.0, 1.5)
				break
			case 2:
				Values = append(Values, 2.0, 2.5)
				break
			case 3:
				Values = append(Values, 3.0, 3.5)
				break
			case 4:
				Values = append(Values, 4.0, 4.5)
				break
			case 5:
				Values = append(Values, 5.0, 5.5)
				break
			default:
				break
			}
		}

		filters = append(filters, dtos.SearchHotelsFloatFilter{
			Field: "star",
			Value: Values,
		})
	}

	return filters
}

func createRoomsData(searchDto dto.SearchDto) []dtos.ProviderRequestSessionRoomDto {
	rooms := make([]dtos.ProviderRequestSessionRoomDto, len(searchDto.Rooms))
	for i, room := range searchDto.Rooms {
		rooms[i] = dtos.ProviderRequestSessionRoomDto{
			Adults:   room.Adults,
			Children: room.Children,
		}
	}
	return rooms
}

func (p *searchDtoAdopter) CreateResultDto(results []dtos.Result, hotels []dbmodel.Hotel, cityId string, days int) (*dto.SearchResponseDto, error) {
	city, err := p.cacheStore.CityStore().Get(cityId)
	if err != nil {
		return nil, err
	}
	resultsDto := make([]dto.SearchResponseHotelDto, 0, len(results))

	for _, result := range results {
		hotel := getHotel(hotels, result.ID)
		if hotel == nil {
			continue
		}
		newRes := *p.createSearchResponseHotelDto(result, city, days)
		newRes.Description = hotel.Description
		newRes.Location.CityEn = hotel.CityEn
		newRes.RateReview.Count = hotel.RateReviewCount
		newRes.RateReview.Score = hotel.RateReviewScore
		newRes.Type = hotel.Type

		newRes.Badges = make([]dto.SearchResponseBadgeDto, 0)
		for _, badge := range hotel.Badges {
			newRes.Badges = append(newRes.Badges, dto.SearchResponseBadgeDto{
				Name: badge.Text,
				Icon: badge.Icon,
			})
		}
		resultsDto = append(resultsDto, newRes)
	}

	return &dto.SearchResponseDto{
		Result: resultsDto,
		Error:  nil,
	}, nil
}

func getHotel(hotels []dbmodel.Hotel, id string) *dbmodel.Hotel {
	for i := range hotels {
		if hotels[i].PlaceID == id {
			return &hotels[i]
		}
	}
	logger.WithName(logtags.NewHotelIdDetected).WithData(id).Warn("Cannot find this hotel by its id. maybe you have to sync this one from provider")
	return nil
}

func (p *searchDtoAdopter) createSearchResponseHotelDto(result dtos.Result, city dto.CityDto, days int) *dto.SearchResponseHotelDto {
	amenities := p.cacheStore.AmenityStore().FindByIds(result.Facilities)
	amenitiesName := make([]string, 0)
	for _, amenity := range amenities {
		amenitiesName = append(amenitiesName, amenity.Name)
	}
	var discount = 0
	var discountPrice = 0
	var oldPricePerNight = math.Round(result.PricePerNight/10000) * 10000
	if result.OldPrice != nil {
		discount = int(math.Round(((*result.OldPrice - result.MinPrice) / *result.OldPrice) * 100))
		if days != 0 {
			oldPricePerNight = (*result.OldPrice) / float64(days)
		}
		discountPrice = int((*result.OldPrice - result.MinPrice) / float64(days))
	}

	var oldPriceInt = int64(oldPricePerNight)

	return &dto.SearchResponseHotelDto{
		PlaceID:         result.ID,
		Id:              result.ID,
		Kind:            constants.DefaultHotelSyncKind,
		MinNight:        1,
		ReservationType: constants.DefaultHotelSyncReservationType,
		PaymentType:     constants.DefaultHotelSyncPaymentType,
		Name:            result.Name.Fa,
		NameEn:          result.Name.En,
		Region:          "",
		Images:          append([]string{result.Image}, result.Images...),
		Image:           result.Image,
		MinPrice:        result.MinPrice,
		PricePerNight:   math.Round(result.PricePerNight/10000) * 10000,
		Location: dto.SearchResponseLocationDto{
			City: city.Name,
			Geo: dto.SearchResponseLocationGeoDto{
				Lat: result.Location.Coordinates[1],
				Lon: result.Location.Coordinates[0],
			},
			Province: city.State,
		},
		Tags:     nil,
		Verified: true,
		RateReview: dto.SearchResponseRateReviewDto{
			Count: 0,
			Score: 0.0,
		},
		Amenities:     amenitiesName,
		Star:          result.Star,
		OldPrice:      &oldPriceInt,
		Discount:      discount,
		DiscountPrice: discountPrice,
	}
}

func NewProviderSearchDtoFactory(cacheStore core.CacheStore) core.SearchDtoAdopter {
	return &searchDtoAdopter{cacheStore: cacheStore}
}
