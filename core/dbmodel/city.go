package dbmodel

type City struct {
	ID      uint   `gorm:"primary_key;"`
	CityID  string `gorm:"column:CityID;type:nvarchar(50);not null;unique_index"`
	Name    string `gorm:"column:Name;type:nvarchar(100);not null"`
	Country string `gorm:"column:Country;type:nvarchar(100);not null"`
	State   string `gorm:"column:State;type:nvarchar(100);not null"`
	BaseId  int64  `gorm:"column:BaseId;not null"`
}
