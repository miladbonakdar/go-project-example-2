package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"hotel-engine/utils/indraframework"
)

type AmenityCategoryDto struct {
	ID      uint                           `json:"id"`
	Name    string                         `json:"name"`
	NameEn  string                         `json:"nameEn"`
	Order   uint                           `json:"order"`
	IconUrl string                         `json:"iconUrl"`
	Error   *indraframework.IndraException `json:"error"`
}

func (a *AmenityCategoryDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}

func (a AmenityCategoryDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.NameEn, validation.Required),
	)
}
