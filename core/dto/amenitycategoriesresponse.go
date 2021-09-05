package dto

type AmenityCategoriesResponse struct {
	Categories []*AmenityCategoryDto `json:"categories"`
}

func NewAmenityCategoriesResponse(cats []*AmenityCategoryDto) AmenityCategoriesResponse {
	return AmenityCategoriesResponse{
		Categories: cats,
	}
}
