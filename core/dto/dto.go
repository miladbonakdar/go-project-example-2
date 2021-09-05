package dto

import (
	"hotel-engine/utils/date"
	"hotel-engine/utils/indraframework"
)

type Dto interface {
	SetError(exc *indraframework.IndraException)
}

func CheckForDate(dateToCheck string) error {
	_, err := date.StringToDate(dateToCheck)
	return err
}
