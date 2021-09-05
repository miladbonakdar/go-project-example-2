package logic

import (
	"hotel-engine/core"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/logic/stores"
	"hotel-engine/infrastructure/logger"
	"sync"
)

var store *cacheStore

type cacheStore struct {
	unitOfWork   core.UnitOfWork
	cityStore    core.CityCacheStore
	amenityStore core.AmenityCacheStore
	mapper       core.Mapper
}

func (s *cacheStore) AmenityStore() core.AmenityCacheStore {
	return s.amenityStore
}

func (s *cacheStore) CityStore() core.CityCacheStore {
	return s.cityStore
}

func (s *cacheStore) UpdateStores() error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		s.UpdateCityStore()
		wg.Done()
	}()
	go func() {
		s.UpdateAmenityStore()
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func (s *cacheStore) UpdateCityStore() {
	cities := s.unitOfWork.City().GetAll()
	s.cityStore = stores.NewCityCacheStore(s.mapper.ToCitiesDto(cities))
}

func (s *cacheStore) UpdateAmenityStore() {
	amenities := s.unitOfWork.Amenity().GetAll()
	s.amenityStore = stores.NewAmenityCacheStore(s.mapper.ToAmenitiesDto(amenities))
}

func NewCacheStore(unit core.UnitOfWork, mapper core.Mapper) core.CacheStore {
	store = &cacheStore{
		unitOfWork: unit,
		mapper:     mapper,
	}
	if err := store.UpdateStores(); err != nil {
		logger.WithName(logtags.InitializingCacheStoresError).FatalException(err, "error in initializing stores")
	}
	return store
}
