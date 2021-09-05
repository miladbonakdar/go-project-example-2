package dbmodel

type Badge struct {
	ID              string `gorm:"primary_key;type:nvarchar(50)"`
	Text            string `gorm:"column:Text;type:nvarchar(300);not null"`
	Icon            string `gorm:"column:Icon;type:nvarchar(300)"`
	TextColor       string `gorm:"column:TextColor;type:nvarchar(300)"`
	BackgroundColor string `gorm:"column:BackgroundColor;type:nvarchar(300)"`
}

func (a *Badge) UpdateIconUrl(iconUrl string) {
	a.Icon = iconUrl
}
