package dto

import (
	"hotel-engine/utils/indraframework"
)

type DeleteAmenityCategoryResponse struct {
	ID    uint                           `json:"id"`
	Error *indraframework.IndraException `json:"error"`
}

func (a *DeleteAmenityCategoryResponse) SetError(exc *indraframework.IndraException) {
	a.Error = exc
}

func NewDeleteAmenityCategoryResponse(id uint) DeleteAmenityCategoryResponse {
	return DeleteAmenityCategoryResponse{ID: id}
}
