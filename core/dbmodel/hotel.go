package dbmodel

import (
	"hotel-engine/core/common"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Hotel struct {
	gorm.Model
	PlaceID         string `gorm:"column:PlaceId;type:nvarchar(50);not null;unique_index"`
	HotelCode       string `gorm:"column:HotelCode;type:nvarchar(5);not null;unique_index"`
	RoomID          string `gorm:"column:RoomId;type:nvarchar(50);not null"`
	Type            string `gorm:"column:Type;type:nvarchar(50);not null"`
	Kind            string `gorm:"column:Kind;type:nvarchar(50);not null"`
	MinNight        int    `gorm:"column:MinNight;not null;default:0"`
	ReservationType string `gorm:"column:ReservationType;type:nvarchar(50);not null"`
	PaymentType     string `gorm:"column:PaymentType;type:nvarchar(50);not null"`
	Name            string `gorm:"column:Name;type:nvarchar(100);not null"`
	NameEn          string `gorm:"column:NameEn;type:nvarchar(100);null"`
	Description     string `gorm:"column:Description;type:nvarchar(4000);not null"`
	City            string `gorm:"column:City;type:nvarchar(100);not null"`
	CityEn          string `gorm:"column:CityEn;type:nvarchar(100);not null"`
	Province        string `gorm:"column:Province;type:nvarchar(100);not null"`
	ProvinceEn      string `gorm:"column:ProvinceEn;type:nvarchar(100);not null"`
	Country         string `gorm:"column:Country;type:nvarchar(100);not null"`
	CountryEn       string `gorm:"column:CountryEn;type:nvarchar(100);not null"`
	CountryCode     string `gorm:"column:CountryCode;type:nvarchar(100);not null"`
	GeoLocation     string `gorm:"column:GeoLocation;type:nvarchar(200);not null"`
	Region          string `gorm:"column:Region;type:nvarchar(50);not null"`
	SuitableFor     string `gorm:"type:nvarchar(1024);column:SuitableFor;not null"`
	Images          string `gorm:"type:nvarchar(4000);column:Images;not null"`

	CheckInTime     string    `gorm:"column:CheckInTime;not null"`
	CheckOutTime    string    `gorm:"column:CheckOutTime;not null"`
	Address         string    `gorm:"column:Address;type:nvarchar(2500);not null"`
	CheckIn         time.Time `gorm:"column:CheckIn;not null"`
	CheckOut        time.Time `gorm:"column:CheckOut;not null"`
	Price           int64     `gorm:"column:Price;not null"`
	OldPrice        int64     `gorm:"column:OldPrice;not null;default:0"`
	DiscountPrice   int64     `gorm:"column:DiscountPrice;not null;default:0"`
	DiscountPercent int64     `gorm:"column:DiscountPercent;not null;default:0"`
	Capacity        int       `gorm:"column:Capacity;not null"`

	Tags     string `gorm:"type:nvarchar(2500);column:Tags;not null"`
	Verified bool   `gorm:"column:Verified;not null"`

	RateReviewScore float64 `gorm:"column:RateReview_Score;not null;default:0"`
	Star            int     `gorm:"column:Star;not null"`
	RateReviewCount int     `gorm:"column:RateReview_Count;not null;default:0"`

	Amenities []*Amenity `gorm:"many2many:hotel_amenity;"`
	Places    []*Place   `gorm:"many2many:hotel_place;"`
	Badges    []*Badge   `gorm:"many2many:hotel_badge;"`
	Orders    []Order    `gorm:"foreignKey:HotelID"`
	Sort      float64    `gorm:"column:Sort;not null;default:0"`
	FAQList   []*FAQ     `gorm:"many2many:hotel_faq;"`
	FAQTitle  string     `gorm:"type:nvarchar(4000);column:FAQTitle"`

	SeoTitle           string `gorm:"column:SeoTitle;type:nvarchar(500)"`
	SeoH1              string `gorm:"column:SeoH1;type:nvarchar(500)"`
	SeoDescription     string `gorm:"column:SeoDescription;type:nvarchar(4000)"`
	SeoRobots          string `gorm:"column:SeoRobots;type:nvarchar(4000)"`
	SeoCanonical       string `gorm:"column:SeoCanonical;type:nvarchar(4000)"`
	SeoMetaDescription string `gorm:"column:SeoMetaDescription;type:nvarchar(4000)"`
	Code               int    `gorm:"column:Code;not null;default:0"`
}

func (h *Hotel) UpdateWith(newHotel Hotel) *Hotel {
	h.Tags = newHotel.Tags
	h.RoomID = newHotel.RoomID
	h.Kind = newHotel.Kind
	h.MinNight = newHotel.MinNight
	h.ReservationType = newHotel.ReservationType
	h.PaymentType = newHotel.PaymentType
	h.Description = newHotel.Description
	h.Name = newHotel.Name
	h.NameEn = strings.ToLower(newHotel.NameEn)
	h.GeoLocation = newHotel.GeoLocation
	h.City = newHotel.City
	h.Province = newHotel.Province
	h.CityEn = newHotel.CityEn
	h.ProvinceEn = newHotel.ProvinceEn
	h.Country = newHotel.Country
	h.CountryEn = newHotel.CountryEn
	h.CountryCode = newHotel.CountryCode
	h.Region = newHotel.Region
	h.SuitableFor = newHotel.SuitableFor
	h.Images = newHotel.Images

	h.Price = newHotel.Price
	h.OldPrice = newHotel.OldPrice
	h.DiscountPercent = newHotel.DiscountPercent
	h.DiscountPrice = newHotel.DiscountPrice
	h.CheckIn = newHotel.CheckIn
	h.CheckOut = newHotel.CheckOut

	h.Capacity = newHotel.Capacity

	h.mergeHotelAmenities(newHotel.Amenities)
	h.mergeHotelBadges(newHotel.Badges)
	h.Places = newHotel.Places
	h.Tags = newHotel.Tags
	h.Verified = newHotel.Verified
	h.Star = newHotel.Star
	h.CheckInTime = newHotel.CheckInTime
	h.CheckOutTime = newHotel.CheckOutTime
	h.Address = newHotel.Address
	h.Type = newHotel.Type
	h.Sort = newHotel.Sort

	return h
}

func (h *Hotel) UpdateRateAndReview(rateCount int, rate float64) *Hotel {
	h.RateReviewScore = rate
	h.RateReviewCount = rateCount
	return h
}

func (h *Hotel) mergeHotelAmenities(amenities []*Amenity) {
	if h.Amenities == nil {
		h.Amenities = make([]*Amenity, 0)
	}
	for _, amenity := range amenities {
		found := false
		for _, hAmenity := range h.Amenities {
			if amenity.ID == hAmenity.ID {
				found = true
				hAmenity.NameEn = amenity.NameEn
				hAmenity.Name = amenity.Name
				hAmenity.GroupId = amenity.GroupId
				break
			}
		}
		if !found {
			h.Amenities = append(h.Amenities, amenity)
		}
	}
}

func (h *Hotel) mergeHotelBadges(badges []*Badge) {
	if h.Badges == nil {
		h.Badges = make([]*Badge, 0)
	}
	for _, badge := range badges {
		found := false
		for _, hBadge := range h.Badges {
			if badge.ID == hBadge.ID {
				found = true
				hBadge.TextColor = badge.TextColor
				hBadge.Text = badge.Text
				hBadge.BackgroundColor = badge.BackgroundColor
				break
			}
		}
		if !found {
			h.Badges = append(h.Badges, badge)
		}
	}
}

func (h *Hotel) UpdateSeoTags(seoTitle, seoH1, seoDescription, seoRobots, seoCanonical string, seoMetaDescription string) *Hotel {
	h.SeoTitle = seoTitle
	h.SeoCanonical = seoCanonical
	h.SeoDescription = seoDescription
	h.SeoH1 = seoH1
	h.SeoRobots = seoRobots
	h.SeoMetaDescription = seoMetaDescription
	return h
}

func (h *Hotel) UpdateFaqDetails(faqTitle string, faqList []*FAQ) {
	if faqTitle != "" {
		h.FAQTitle = faqTitle
	}
	if len(faqList) == 0 {
		return
	}
	if h.FAQList == nil {
		h.FAQList = make([]*FAQ, 0)
	}
	for _, faq := range faqList {
		if faq.ID == 0 {
			h.FAQList = append(h.FAQList, faq)
			continue
		}
		found := false
		for _, hFaq := range h.FAQList {
			if faq.ID == hFaq.ID {
				found = true
				hFaq.Answer = faq.Answer
				hFaq.Question = faq.Question
				break
			}
		}
		if !found {
			h.FAQList = append(h.FAQList, faq)
		}
	}
}

func (h *Hotel) GetHotelFAQ(faqId uint) (FAQ, error) {
	index := -1
	for i, hFaq := range h.FAQList {
		if faqId == hFaq.ID {
			index = i
			break
		}
	}
	if index == -1 {
		return FAQ{}, common.FAQNotFound
	}
	return *h.FAQList[index], nil
}
