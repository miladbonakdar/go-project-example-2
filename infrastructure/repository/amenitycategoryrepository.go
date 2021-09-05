package repository

import (
	"github.com/jinzhu/gorm"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/dbmodel"
)

type amenityCategoryRepository struct {
	DB *gorm.DB
}

func (r *amenityCategoryRepository) GetOneById(id uint) (*dbmodel.AmenityCategory, error) {
	var cat dbmodel.AmenityCategory
	if r.DB.Find(&cat, "id=?", id).RecordNotFound() {
		return nil, common.AmentityCategoryNotFound
	}
	return &cat, nil
}

func (r *amenityCategoryRepository) StoreOrUpdate(cat dbmodel.AmenityCategory) error {
	return r.DB.Save(&cat).Error
}

func (r *amenityCategoryRepository) Delete(id int) error {
	var cat dbmodel.AmenityCategory
	deleteResult := r.DB.Delete(&cat, id)
	if deleteResult.RecordNotFound() {
		return common.AmentityCategoryNotFound
	}
	return deleteResult.Error
}

func (r *amenityCategoryRepository) GetAll() []*dbmodel.AmenityCategory {
	var cats []*dbmodel.AmenityCategory
	r.DB.Find(&cats)
	return cats
}

func newAmenityCategory(DB *gorm.DB) core.AmenityCategoryRepository {
	return &amenityCategoryRepository{DB: DB}
}
