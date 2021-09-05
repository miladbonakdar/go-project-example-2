package dbmodel

type FAQ struct {
	ID       uint   `gorm:"primary_key"`
	Question string `gorm:"column:Question;type:nvarchar(4000);not null"`
	Answer   string `gorm:"column:Answer;type:nvarchar(4000);not null"`
}
