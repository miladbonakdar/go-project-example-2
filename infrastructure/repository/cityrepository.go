package repository

import (
	"github.com/jinzhu/gorm"
	"hotel-engine/core"
	"hotel-engine/core/dbmodel"
)

type cityRepository struct {
	DB *gorm.DB
}

func (r *cityRepository) GetAll() []dbmodel.City {
	var cities []dbmodel.City
	r.DB.Find(&cities)

	return cities
}

func (r *cityRepository) BulkInsert(cities []dbmodel.City) error {
	var err error

	for _, city := range cities {
		if !r.ExistByCityId(city.CityID) {
			err = r.DB.Save(&city).Error
		}
	}

	return err
}

func (r *cityRepository) ExistByCityId(cityId string) bool {
	var city dbmodel.City
	if r.DB.Find(&city, "CityID=?", cityId).RecordNotFound() {
		return false
	}
	return true
}

func newCityRepository(DB *gorm.DB) core.CityRepository {
	return &cityRepository{DB: DB}
}
