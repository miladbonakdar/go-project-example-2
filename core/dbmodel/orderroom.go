package dbmodel

import (
	"github.com/jinzhu/gorm"
)

type OrderRoom struct {
	gorm.Model
	OrderID       uint   `gorm:"column:OrderID;not null"`
	PricePerNight int64  `gorm:"column:PricePerNight;not null"`
	Price         int64  `gorm:"column:Price;not null"`
	Name          string `gorm:"column:Name;type:nvarchar(100);not null"`
	NameEn        string `gorm:"column:NameEn;type:nvarchar(100);not null"`
}
