package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SetBadgeIconDto struct {
	BadgeId string `json:"badgeId"`
	IconUrl string `json:"iconUrl"`
}

func (a SetBadgeIconDto) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.BadgeId, validation.Required),
		validation.Field(&a.IconUrl, validation.Required),
	)
}
