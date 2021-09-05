package repository

import (
	"github.com/jinzhu/gorm"
	"hotel-engine/core"
)

type unitOfWork struct {
	hotel           core.HotelRepository
	amenity         core.AmenityRepository
	city            core.CityRepository
	order           core.OrderRepository
	place           core.PlaceRepository
	amenityCategory core.AmenityCategoryRepository
	badge           core.BadgeRepository
}

func (u *unitOfWork) Hotel() core.HotelRepository {
	return u.hotel
}
func (u *unitOfWork) City() core.CityRepository {
	return u.city
}

func (u *unitOfWork) Order() core.OrderRepository {
	return u.order
}

func (u *unitOfWork) Amenity() core.AmenityRepository {
	return u.amenity
}

func (u *unitOfWork) Place() core.PlaceRepository {
	return u.place
}

func (u *unitOfWork) AmenityCategory() core.AmenityCategoryRepository {
	return u.amenityCategory
}

func (u *unitOfWork) Badge() core.BadgeRepository {
	return u.badge
}

func NewUnitOfWork(DB *gorm.DB) core.UnitOfWork {
	return &unitOfWork{
		hotel:           newHotelRepository(DB),
		amenity:         newAmenityRepository(DB),
		city:            newCityRepository(DB),
		order:           newOrderRepository(DB),
		place:           newPlaceRepository(DB),
		amenityCategory: newAmenityCategory(DB),
		badge:           newBadgeRepository(DB),
	}
}
