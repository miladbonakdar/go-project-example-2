package repository

import (
	"github.com/jinzhu/gorm"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/dbmodel"
)

type amenityRepository struct {
	DB *gorm.DB
}

func (r *amenityRepository) GetAll() []dbmodel.Amenity {
	var amenities []dbmodel.Amenity
	r.DB.Preload("AmenityCategory").Find(&amenities)
	return amenities
}

func (r *amenityRepository) UpdateIcon(amenityId int, amenityIcon string) (*dbmodel.Amenity, error) {
	var amenity dbmodel.Amenity

	if r.DB.Find(&amenity, amenityId).RecordNotFound() {
		return nil, common.AmenityNotFound
	}
	amenity.UpdateIconUrl(amenityIcon)
	return &amenity, r.DB.Save(&amenity).Error
}

func (r *amenityRepository) UpdateCategory(amenityId int, amenityCategoryId uint) (*dbmodel.Amenity, error) {
	var amenity dbmodel.Amenity
	if r.DB.Preload("AmenityCategory").Find(&amenity, amenityId).RecordNotFound() {
		return nil, common.AmenityNotFound
	}
	amenity.UpdateCategory(amenityCategoryId)
	return &amenity, r.DB.Save(&amenity).Error
}

func newAmenityRepository(DB *gorm.DB) core.AmenityRepository {
	return &amenityRepository{DB: DB}
}
