package repository

import (
	"github.com/jinzhu/gorm"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/dbmodel"
)

type badgeRepository struct {
	DB *gorm.DB
}

func (r *badgeRepository) GetAll() []dbmodel.Badge {
	var badges []dbmodel.Badge
	r.DB.Find(&badges)
	return badges
}

func (r *badgeRepository) UpdateIcon(badgeId string, icon string) (*dbmodel.Badge, error) {
	var badge dbmodel.Badge

	if r.DB.Find(&badge, "id=?", badgeId).RecordNotFound() {
		return nil, common.BadgeNotFound
	}
	badge.UpdateIconUrl(icon)
	return &badge, r.DB.Save(&badge).Error
}

func newBadgeRepository(DB *gorm.DB) core.BadgeRepository {
	return &badgeRepository{DB: DB}
}
