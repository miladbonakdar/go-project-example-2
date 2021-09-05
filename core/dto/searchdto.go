package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type SearchDto struct {
	SessionId     string           `json:"sessionId,omitempty"`
	Rooms         []RequestRoomDto `json:"rooms"`
	Date          SearchDateDto    `json:"date"`
	Keyword       string           `json:"keyword,omitempty"`
	City          string           `json:"city"`
	PageNumber    int64            `json:"page-number"`
	PageSize      int64            `json:"page-size"`
	Price         SearchPriceDto   `json:"price,omitempty"`
	Region        []string         `json:"region,omitempty"`
	Score         SearchScoreDto   `json:"score,omitempty"`
	Sort          string           `json:"sort,omitempty"`
	Stars         []int            `json:"stars,omitempty"`
	SortDirection bool             `json:"sort-direction,omitempty"`
	HotelTypes    []string         `json:"hotel-types,omitempty"`
}

type SearchDateDto struct {
	End   string `json:"end"`
	Start string `json:"start"`
}

type SearchPriceDto struct {
	End   int64 `json:"end"`
	Start int64 `json:"start"`
}
type SearchScoreDto struct {
	End   int64 `json:"end"`
	Start int64 `json:"start"`
}

func (a SearchDto) Validate() error {
	if a.Date.End != "" {
		if err := CheckForDate(a.Date.End); err != nil {
			return err
		}
	}
	if a.Date.Start != "" {
		if err := CheckForDate(a.Date.Start); err != nil {
			return err
		}
	}

	for _, room := range a.Rooms {
		if err := room.Validate(); err != nil {
			return err
		}
	}

	return validation.ValidateStruct(&a,
		validation.Field(&a.City, validation.Required),
		validation.Field(&a.Date, validation.Required),
	)
}
