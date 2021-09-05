package dbmodel

type Place struct {
	ID          string  `gorm:"primary_key;type:nvarchar(50)"`
	Name        string  `gorm:"column:Name;type:nvarchar(100);not null"`
	GeoLocation string  `gorm:"column:GeoLocation;type:nvarchar(200);not null"`
	Distance    float64 `gorm:"column:Distance;not null"`
	Priority    int     `gorm:"column:Priority;not null"`
}
