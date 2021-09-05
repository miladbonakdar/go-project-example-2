package dto

import "hotel-engine/core/common"

type FiltersDto struct {
	PriceFilter    FilterOutput `json:"price_filters"`
	StarFilters    FilterOutput `json:"star_filters"`
	ScoreFilters   FilterOutput `json:"score_filters"`
	AmenityFilters FilterOutput `json:"amenity_filters"`
	HotelTypes     FilterOutput `json:"hotel_types"`
}

type FilterDto struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type FilterOutput struct {
	Name    string      `json:"name"`
	Field   string      `json:"field"`
	Filters []FilterDto `json:"filters"`
}

const (
	Price_0_500000 = iota + 1
	Price_500000_1000000
	Price_1000000_1500000
	Price_1500000_2000000
	Price_over_2000000

	Score_lower_than_6
	Score_6_7
	Score_7_8
	Score_8_9
	Score_9_10
)

func CreateDefaultFiltersDto(amenities []HotelAmenityDto) *FiltersDto {
	amenityFilters := make([]FilterDto, 0)
	if amenities != nil {
		for _, amenity := range amenities {
			amenityFilters = append(amenityFilters, FilterDto{
				Key:   amenity.Name,
				Value: int64(amenity.ID),
			})
		}
	}
	return &FiltersDto{
		PriceFilter: FilterOutput{
			Name:  "رنج قیمتی هتل",
			Field: "minPrice",
			Filters: []FilterDto{
				{
					Key:   "0  تا  500,000",
					Value: Price_0_500000,
				},
				{
					Key:   "500,000  تا  1,000,000",
					Value: Price_500000_1000000,
				},
				{
					Key:   "1,000,000  تا  1,500,000",
					Value: Price_1000000_1500000,
				},
				{
					Key:   "1,500,000  تا  2,000,000",
					Value: Price_1500000_2000000,
				},
				{
					Key:   "2,000,000 به بالا",
					Value: Price_over_2000000,
				},
			},
		},
		StarFilters: FilterOutput{
			Name:  "درجه هتل",
			Field: "star",
			Filters: []FilterDto{
				{
					Key:   "یک ستاره",
					Value: 1,
				},
				{
					Key:   "دو ستاره",
					Value: 2,
				},
				{
					Key:   "سه ستاره",
					Value: 3,
				},
				{
					Key:   "چهار ستاره",
					Value: 4,
				},
				{
					Key:   "پنج ستاره",
					Value: 5,
				},
				{
					Key:   "رتبه بندی نشده",
					Value: -1,
				},
			},
		},
		ScoreFilters: FilterOutput{
			Name:  "محبوبیت",
			Field: "score",
			Filters: []FilterDto{
				{
					Key:   "فوق العاده (۹ تا ۱۰)",
					Value: Score_9_10,
				},
				{
					Key:   "بسیار عالی (۸ تا ۹)",
					Value: Score_8_9,
				},
				{
					Key:   "عالی (۷ تا ۸)",
					Value: Score_7_8,
				},
				{
					Key:   "خوب (۶ تا ۷)",
					Value: Score_6_7,
				},
				{
					Key:   "قابل قبول (کمتر از ۶)",
					Value: Score_lower_than_6,
				},
			},
		},
		AmenityFilters: FilterOutput{
			Name:    "دیگر امکانات",
			Field:   "facility",
			Filters: amenityFilters,
		},
		HotelTypes: FilterOutput{
			Name:  "براساس نوع هتل",
			Field: "types",
			Filters: []FilterDto{
				{
					Key:   common.HotelType_Hotel,
					Value: -1,
				},
				{
					Key:   common.HotelType_HotelApartment,
					Value: -1,
				},
			},
		},
	}
}
