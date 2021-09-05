package dto

import (
	"hotel-engine/utils/indraframework"
)

type UpdateResultDto struct {
	Message    string                         `json:"message"`
	TimeTaken  int64                          `json:"time_taken"`
	ItemsCount int                            `json:"items_count"`
	Error      *indraframework.IndraException `json:"error"`
}

func (a *UpdateResultDto) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
