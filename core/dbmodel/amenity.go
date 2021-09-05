package dbmodel

type Amenity struct {
	ID                int    `gorm:"primary_key;auto_increment:false"`
	Name              string `gorm:"column:Name;type:nvarchar(100);not null"`
	NameEn            string `gorm:"column:NameEn;type:nvarchar(100);null"`
	GroupId           int    `gorm:"column:GroupId;not null"`
	IconUrl           string `gorm:"column:IconUrl;type:nvarchar(2500);not null"`
	AmenityCategoryID *uint
	AmenityCategory   *AmenityCategory
}

func (a *Amenity) UpdateIconUrl(iconUrl string) {
	a.IconUrl = iconUrl
}

func (a *Amenity) UpdateCategory(catId uint) {
	a.AmenityCategoryID = &catId
}
