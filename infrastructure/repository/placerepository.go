package repository

import (
	"github.com/jinzhu/gorm"
	"hotel-engine/core"
	"hotel-engine/core/dbmodel"
)

type placeRepository struct {
	DB *gorm.DB
}

func (r *placeRepository) GetAll() []dbmodel.Place {
	var places []dbmodel.Place
	r.DB.Find(&places)
	return places
}

func newPlaceRepository(DB *gorm.DB) core.PlaceRepository {
	return &placeRepository{DB: DB}
}
