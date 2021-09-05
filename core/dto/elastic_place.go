package dto

type ElasticUpdateRequest struct {
	Places     []ElasticHotel `json:"places"`
	ClearCache int            `json:"clearCache"`
}

//ElasticHotel ...
type ElasticHotel struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	NameEn          string            `json:"nameEn"`
	Description     string            `json:"description"`
	Type            string            `json:"type"`
	Kind            string            `json:"kind"`
	PlaceID         string            `json:"place_id"`
	RoomID          string            `json:"room_id"`
	Location        ElasticLocation   `json:"location"`
	SuitableFor     []string          `json:"suitable_for"`
	Region          string            `json:"region"`
	Image           string            `json:"image"`
	Images          []string          `json:"images"`
	Tags            []string          `json:"tags"`
	Amenities       []ElasticAmenity  `json:"amenities"`
	MinPrice        int               `json:"min_price"`
	Calendar        []ElasticCalendar `json:"calendar"`
	ReservationType string            `json:"reservation_type"`
	PaymentType     string            `json:"payment_type"`
	RateReview      ElasticRateReview `json:"rate_review"`
	Verified        bool              `json:"verified"`
	MinNight        int               `json:"min_night"`
	Star            int               `json:"star"`
	Code            int               `json:"code"`
	Status          string            `json:"status"`
	DiscountPercent int64             `json:"discount_percent_hotel"`
	OldPrice        int64             `json:"old_price_hotel"`
	DiscountPrice   int64             `json:"discount_price_hotel"`
	Sort            float64           `json:"sort"`
	Badges          []ElasticBadge    `json:"badges"`
}

type ElasticAmenity struct {
	Category     string `json:"category"`
	CategoryName string `json:"category_name"`
	Name         string `json:"name"`
}

type ElasticBadge struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type ElasticCalendar struct {
	Year      int                     `json:"year"`
	Month     int                     `json:"month"`
	Day       int                     `json:"day"`
	Date      int                     `json:"date"`
	Price     int                     `json:"price"`
	Capacity  ElasticCalendarCapacity `json:"capacity"`
	Available int                     `json:"available"`
}

type ElasticCalendarCapacity struct {
	Base  int `json:"base"`
	Extra int `json:"extra"`
}

type ElasticRateReview struct {
	Score float64 `json:"score"`
	Count int     `json:"count"`
}

type ElasticLocation struct {
	City       string             `json:"city"`
	Province   string             `json:"province"`
	Geo        ElasticLocationGeo `json:"geo"`
	CityEn     string             `json:"cityEn"`
	ProvinceEn string             `json:"provinceEn"`
}

type ElasticLocationGeo struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
