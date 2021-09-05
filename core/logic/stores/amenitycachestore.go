package stores

import (
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/dto"
	"strings"
)

type amenityCacheStore struct {
	amenities []dto.HotelAmenityDto
}

func (s *amenityCacheStore) GetAll() []dto.HotelAmenityDto {
	return s.amenities
}

func (s *amenityCacheStore) FindOne(name string) (dto.HotelAmenityDto, error) {
	for _, amenity := range s.amenities {
		if strings.Contains(amenity.Name, name) {
			return amenity, nil
		}
	}
	return dto.HotelAmenityDto{}, common.AmenityNotFound
}

func (s *amenityCacheStore) FindByIds(ids []int) []dto.HotelAmenityDto {
	amenities := make([]dto.HotelAmenityDto, 0)
	for _, id := range ids {
		if amenity, err := s.Get(id); err == nil {
			amenities = append(amenities, amenity)
		}
	}
	return amenities
}

func (s *amenityCacheStore) Find(name string) []dto.HotelAmenityDto {
	amenities := make([]dto.HotelAmenityDto, 0)
	for _, amenity := range s.amenities {
		if strings.Contains(amenity.Name, name) {
			amenities = append(amenities, amenity)
		}
	}
	return amenities
}

func (s *amenityCacheStore) Get(id int) (dto.HotelAmenityDto, error) {
	for _, amenity := range s.amenities {
		if amenity.ID == id {
			return amenity, nil
		}
	}
	return dto.HotelAmenityDto{}, common.AmenityNotFound
}

func NewAmenityCacheStore(amenities []dto.HotelAmenityDto) core.AmenityCacheStore {
	return &amenityCacheStore{
		amenities: amenities,
	}
}
