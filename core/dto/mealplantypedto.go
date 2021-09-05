package dto

type MealPlanTypeDto struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	NameEn string `json:"nameEn"`
}

var MealPlans = map[string]MealPlanTypeDto{
	"RO": {
		Key:    "RO",
		Name:   "بدون وعده غذایی",
		NameEn: "Room Only",
	},
	"BB": {
		Key:    "BB",
		Name:   "با صبحانه",
		NameEn: "Bed and Breakfast",
	},
	"HB": {
		Key:    "HB",
		Name:   "صبحانه+نهار یا شام",
		NameEn: "Half Board",
	},
	"FB": {
		Key:    "FB",
		Name:   "صبحانه+نهار+شام",
		NameEn: "Full Board",
	},
	"AI": {
		Key:    "AI",
		Name:   "صبحانه+نهار+شام+میان وعده",
		NameEn: "'All Inclusive",
	},
	"UALL": {
		Key:    "UALL",
		Name:   "صبحانه+نهار+شام+میان وعده+نوشیدنی",
		NameEn: "Ultra All Inclusive",
	},
	"Unknown": {
		Key:    "Unknown",
		Name:   "نا معلوم",
		NameEn: "Unknown",
	},
}