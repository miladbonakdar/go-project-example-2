package dto

import "hotel-engine/utils/indraframework"

type TaskRunningResult struct {
	Message string                         `json:"message"`
	Success bool                           `json:"success"`
	Error   *indraframework.IndraException `json:"error"`
}

func (a *TaskRunningResult) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}
