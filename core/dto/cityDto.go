package dto

type CityDto struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	State   string `json:"state"`
	BaseId  int64  `json:"baseId"`
}
