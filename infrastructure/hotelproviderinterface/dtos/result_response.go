package dtos

import "time"

type ResultResponse struct {
	BaseResponse
	Result Result2 `json:"result"`
}

type Name struct {
	En string `json:"en"`
	Fa string `json:"fa"`
}
type Accommodation struct {
	ID   int  `json:"id"`
	Name Name `json:"name"`
}
type Location struct {
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type Badge struct {
	Text  string     `json:"text"`
	Icon  string     `json:"icon"`
	Id    string     `json:"_id"`
	Color BadgeColor `json:"color"`
}

type BadgeColor struct {
	Text       string `json:"text"`
	Background string `json:"background"`
}

type Country struct {
	Code string `json:"code"`
}

type Result struct {
	ID                   string        `json:"id"`
	Score                int           `json:"score"`
	Star                 int           `json:"star"`
	Accommodation        Accommodation `json:"accommodation"`
	Facilities           []int         `json:"facilities"`
	MinPriceProviderName string        `json:"minPriceProviderName"`
	MinPrice             float64       `json:"minPrice"`
	PricePerNight        float64       `json:"pricePerNight"`
	Area                 interface{}   `json:"area"`
	Special              bool          `json:"special"`
	Usable               bool          `json:"usable"`
	Currency             string        `json:"currency"`
	Providers            []interface{} `json:"providers"`
	Name                 Name          `json:"name"`
	Location             Location      `json:"location"`
	Image                string        `json:"image"`
	Images               []string      `json:"images"`
	Link                 string        `json:"link"`
	Country              Country       `json:"country"`
	Badges               []Badge       `json:"badges"`
	Places               []interface{} `json:"places"`
	SpecialOffer         bool          `json:"specialOffer"`
	OldPrice             *float64      `json:"oldPrice,omitempty"`
}
type Destination struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type Rooms struct {
	Adults   []int         `json:"adults"`
	Children []interface{} `json:"children"`
	ID       string        `json:"_id"`
}
type Request struct {
	Destination Destination `json:"destination"`
	CheckIn     time.Time   `json:"checkIn"`
	CheckOut    time.Time   `json:"checkOut"`
	Rooms       []Rooms     `json:"rooms"`
	StayNo      int         `json:"stayNo"`
	Nationality string      `json:"nationality"`
}
type Star struct {
	Num0 int `json:"0"`
	Num1 int `json:"1"`
	Num2 int `json:"2"`
	Num3 int `json:"3"`
	Num4 int `json:"4"`
	Num5 int `json:"5"`
}
type Score struct {
	Num0 int `json:"0"`
	Num3 int `json:"3"`
	Num4 int `json:"4"`
	Num5 int `json:"5"`
	Num6 int `json:"6"`
	Num7 int `json:"7"`
	Num8 int `json:"8"`
	Num9 int `json:"9"`
}
type FacilitiesStat struct {
	Num1027 int `json:"1027"`
	Num1048 int `json:"1048"`
	Num1056 int `json:"1056"`
	Num1062 int `json:"1062"`
}
type Price struct {
	Num0        int `json:"0"`
	Num5000000  int `json:"5000000"`
	Num10000000 int `json:"10000000"`
}
type AccommodationsStat struct {
	Num204 int `json:"204"`
}
type Group struct {
	Name Name `json:"name"`
}
type Facilities struct {
	ID       string `json:"id"`
	Name     Name   `json:"name"`
	GroupID  int    `json:"groupId"`
	Priority int    `json:"priority"`
	Group    Group  `json:"group"`
}
type Accommodations struct {
	ID   int  `json:"id"`
	Name Name `json:"name"`
}
type Info struct {
	Expire             int                `json:"expire"`
	Area               []interface{}      `json:"area"`
	ResultNo           int                `json:"resultNo"`
	Star               Star               `json:"star"`
	Score              Score              `json:"score"`
	FacilitiesStat     FacilitiesStat     `json:"facilitiesStat"`
	Price              Price              `json:"price"`
	AccommodationsStat AccommodationsStat `json:"accommodationsStat"`
	MinPrice           int                `json:"minPrice"`
	MaxPrice           int                `json:"maxPrice"`
	Facilities         []Facilities       `json:"facilities"`
	Accommodations     []Accommodations   `json:"accommodations"`
	TotalResultNo      int                `json:"totalResultNo"`
}
type Result2 struct {
	Result        []Result `json:"result"`
	Request       Request  `json:"request"`
	Info          Info     `json:"info"`
	Progress      int      `json:"progress"`
	LastChunk     bool     `json:"lastChunk"`
	SessionExpire int      `json:"sessionExpire"`
}
