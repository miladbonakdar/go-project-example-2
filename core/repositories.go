package core

import (
	"hotel-engine/core/dbmodel"
	"time"
)

type UnitOfWork interface {
	Hotel() HotelRepository
	City() CityRepository
	Order() OrderRepository
	Amenity() AmenityRepository
	AmenityCategory() AmenityCategoryRepository
	Place() PlaceRepository
	Badge() BadgeRepository
}

type CityRepository interface {
	GetAll() []dbmodel.City
	BulkInsert(cities []dbmodel.City) error
}

type OrderRepository interface {
	Insert(order dbmodel.Order) (*dbmodel.Order, error)
	GetOneByIndraId(indraId string) (*dbmodel.Order, error)
	StoreOrUpdate(order dbmodel.Order) error
	GetProperOrderIdsForRefundUpdateStatus(fromDate time.Time) ([]string, error)
}

type AmenityRepository interface {
	GetAll() []dbmodel.Amenity
	UpdateIcon(amenityId int, amenityIcon string) (*dbmodel.Amenity, error)
	UpdateCategory(amenityId int, amenityCategoryId uint) (*dbmodel.Amenity, error)
}

type BadgeRepository interface {
	GetAll() []dbmodel.Badge
	UpdateIcon(badgeId string, icon string) (*dbmodel.Badge, error)
}

type AmenityCategoryRepository interface {
	GetOneById(id uint) (*dbmodel.AmenityCategory, error)
	StoreOrUpdate(cat dbmodel.AmenityCategory) error
	Delete(id int) error
	GetAll() []*dbmodel.AmenityCategory
}

type PlaceRepository interface {
	GetAll() []dbmodel.Place
}

type HotelRepository interface {
	FindByID(hotelId string) (*dbmodel.Hotel, error)
	FindByIDForSync(hotelId string) (*dbmodel.Hotel, error)
	StoreOrUpdate(hotel *dbmodel.Hotel) error
	Delete(hotel dbmodel.Hotel) error
	GetAllHotelIds() ([]string, error)
	GetHotelsPageForSync(page int, size int) ([]dbmodel.Hotel, error)
	GetHotels(ids []string) ([]dbmodel.Hotel, error)
	GetHotel(hotelId string) (*dbmodel.Hotel, error)
	GetAllHotels() ([]dbmodel.Hotel, error)
	HasBeenSynced() (bool, error)
	GetHotelsList(page int, size int, search string) ([]dbmodel.Hotel, int, error)
	RemoveFAQ(hotelId string, faq *dbmodel.FAQ) (*dbmodel.Hotel, error)
}
