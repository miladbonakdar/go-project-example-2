package hotelProviderInterface

import (
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/dto"
	"hotel-engine/utils"
	"strconv"
)

type hotelEnumerable struct {
	total       int
	cities      []dto.CityDto
	current     int
	provider    core.HotelProvider
	itemPerPage int
}

func (e *hotelEnumerable) MoveNext() bool {
	e.current++
	return e.current < e.total
}

func (e *hotelEnumerable) Current() ([]dto.HotelIdentifierDto, error) {
	res, err := e.provider.SearchHotels(uint32(e.itemPerPage), 0, "", "",
		"", e.cities[e.current].BaseId)

	if err != nil {
		return nil, err
	}

	hotelIds := make([]dto.HotelIdentifierDto, 0, len(res.HotelsList))
	for _, item := range res.HotelsList {
		typeCode, err := strconv.Atoi(item.AccommodationType)
		if err != nil {
			typeCode = common.HotelTypes[common.HotelType_Hotel]
		}
		typeString := common.HotelTypesString[typeCode]
		if typeString == "" {
			typeString = common.HotelType_Hotel
		}
		hotelIds = append(hotelIds, dto.HotelIdentifierDto{
			Id:        item.Id,
			GiataId:   item.GiataId,
			HotelType: typeString,
		})
	}
	return hotelIds, nil
}

func (e *hotelEnumerable) Reset() {
	e.current = -1
}

func NewHotelEnumerable(provider core.HotelProvider, cities []dto.CityDto) core.HotelEnumerable {
	return &hotelEnumerable{
		cities:      cities,
		total:       len(cities),
		current:     -1,
		provider:    provider,
		itemPerPage: utils.GetHotelPerRequest,
	}
}
