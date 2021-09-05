package dbmodel

import "github.com/jinzhu/gorm"

type AmenityCategory struct {
	gorm.Model
	Name    string `gorm:"column:Name;type:nvarchar(100);not null;unique_index"`
	NameEn  string `gorm:"column:NameEn;type:nvarchar(100);null"`
	IconUrl string `gorm:"column:IconUrl;type:nvarchar(2500);not null"`
	Order   uint   `gorm:"column:Order;not null;default:0"`
}

func (a *AmenityCategory) UpdateIconUrl(iconUrl string) {
	a.IconUrl = iconUrl
}

func (a *AmenityCategory) UpdateOrder(order uint) {
	a.Order = order
}

func (a *AmenityCategory) UpdateName(name, nameEn string) {
	a.Name = name
	a.NameEn = nameEn
}
