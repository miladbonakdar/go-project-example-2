package stores

import (
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/dto"
	"strings"
)

type cityCacheStore struct {
	cities []dto.CityDto
}

func (s *cityCacheStore) GetAll() []dto.CityDto {
	return s.cities
}

func (s *cityCacheStore) FindOne(name string) (dto.CityDto, error) {
	for _, city := range s.cities {
		if strings.Contains(city.Name, name) {
			return city, nil
		}
	}
	return dto.CityDto{}, common.CityNotFound
}

func (s *cityCacheStore) Find(name string) []dto.CityDto {
	cities := make([]dto.CityDto, 0)
	for _, city := range s.cities {
		if strings.Contains(city.Name, name) {
			cities = append(cities, city)
		}
	}
	return cities
}

func (s *cityCacheStore) Get(id string) (dto.CityDto, error) {
	for _, city := range s.cities {
		if city.Id == id {
			return city, nil
		}
	}
	return dto.CityDto{}, common.CityNotFound
}

func NewCityCacheStore(cities []dto.CityDto) core.CityCacheStore {
	return &cityCacheStore{
		cities: cities,
	}
}
