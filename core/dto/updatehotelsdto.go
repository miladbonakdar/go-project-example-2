package dto

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"hotel-engine/utils/date"
	"time"
)

type UpdateHotelsDto struct {
	Date string `json:"date"`
}

func (a UpdateHotelsDto) Validate() error {
	if err := CheckForDate(a.Date); err != nil{
		return err
	}
	updateFor := date.StringToDateOrDefault(a.Date)
	if  time.Now().Add(time.Hour * -24).After(updateFor) {
		return errors.New("update date at least should be today")
	}
	return validation.ValidateStruct(&a,
		validation.Field(&a.Date, validation.Required))
}
