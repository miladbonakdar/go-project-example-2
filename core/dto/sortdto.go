package dto

type SortDto struct {
	Field     string `json:"field"`
	Direction bool   `json:"direction"`
	Name      string `json:"name"`
}
