package logic

import (
	"fmt"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/logger"
	"hotel-engine/utils/atomicflag"
)

var (
	inSyncingCities = atomicflag.NewAtomicFlag()
)

type publicService struct {
	unitOfWork          core.UnitOfWork
	mapper              core.Mapper
	cacheStore          core.CacheStore
	informationProvider core.BasicInformationProvider
}

func (s *publicService) GetAllCities() []dto.CityDto {
	return s.cacheStore.CityStore().GetAll()
}

func (s *publicService) GetAllAmenities() []dto.HotelAmenityDto {
	return s.cacheStore.AmenityStore().GetAll()
}

func (s *publicService) GetAllSorts() []dto.SortDto {
	return []dto.SortDto{
		{
			Field:     "score",
			Direction: false,
			Name:      "بالاترین امتیاز",
		}, {
			Field:     "minPrice",
			Direction: false,
			Name:      "بیشترین قیمت",
		}, {
			Field:     "minPrice",
			Direction: true,
			Name:      "کمترین قیمت",
		},
	}
}

func (s *publicService) GetAllFilters() *dto.FiltersDto {
	return dto.CreateDefaultFiltersDto(nil)
}

func (s *publicService) SyncAllCities() ([]dto.CityDto, error) {
	if inSyncingCities.Get() {
		return nil, common.AlreadyUpdatingHotels
	}
	inSyncingCities.Set(true)
	defer inSyncingCities.Set(false)
	cities, err := s.informationProvider.GetCities()
	if err != nil {
		logger.WithName(logtags.SyncingCitiesError).ErrorException(err, "problem wile syncing cities")
	}
	err = s.unitOfWork.City().BulkInsert(s.mapper.ToCitiesModel(cities))
	s.cacheStore.UpdateCityStore()
	return cities, err
}

func (s *publicService) GetAllPlaces() []dto.HotelPlaceDto {
	places := s.unitOfWork.Place().GetAll()
	return s.mapper.ToPlacesDto(places)
}

func (s *publicService) GetChildAgeRanges() []dto.ChildAgeRangeDto {
	ranges := make([]dto.ChildAgeRangeDto, common.MaxAgeAsAChild)
	for i := 0; i < 12; i++ {
		ranges[i] = dto.ChildAgeRangeDto{
			Name:  fmt.Sprintf("از %d تا %d سال", i, i+1),
			Value: i,
		}
	}
	return ranges
}

func (s *publicService) CreateAmenityCategory(dto dto.AmenityCategoryDto) (dto.AmenityCategoryDto, error) {
	cat := s.mapper.ToAmenityCategoryModel(dto)
	err := s.unitOfWork.AmenityCategory().StoreOrUpdate(cat)
	dto.ID = cat.ID
	return dto, err
}

func (s *publicService) GetAmenityCategory(id uint) (*dto.AmenityCategoryDto, error) {
	cat, err := s.unitOfWork.AmenityCategory().GetOneById(id)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToAmenityCategoryDto(cat), nil
}

func (s *publicService) DeleteAmenityCategory(id int) error {
	return s.unitOfWork.AmenityCategory().Delete(id)
}

func (s *publicService) GetAmenityCategories() []*dto.AmenityCategoryDto {
	return s.mapper.ToAmenityCategoriesDto(s.unitOfWork.AmenityCategory().GetAll())
}

func (s *publicService) UpdateAmenityCategory(item dto.AmenityCategoryDto) (dto.AmenityCategoryDto, error) {
	cat, err := s.unitOfWork.AmenityCategory().GetOneById(item.ID)
	if err != nil {
		return dto.AmenityCategoryDto{}, err
	}
	cat.UpdateIconUrl(item.IconUrl)
	cat.UpdateName(item.Name, item.NameEn)
	cat.UpdateOrder(item.Order)
	err = s.unitOfWork.AmenityCategory().StoreOrUpdate(*cat)
	return item, err
}

func (s *publicService) SetBadgeIcon(request dto.SetBadgeIconDto) (dto.HotelBadgeDto, error) {
	badge, err := s.unitOfWork.Badge().UpdateIcon(request.BadgeId, request.IconUrl)
	if err != nil {
		return dto.HotelBadgeDto{}, err
	}
	return *s.mapper.ToBadgeDto(*badge), nil
}

func (s *publicService) GetAllBadges() []dto.HotelBadgeDto {
	badges := s.unitOfWork.Badge().GetAll()
	return s.mapper.ToBadgesDto(badges)
}

func NewPublicService(unit core.UnitOfWork, mapper core.Mapper,
	cacheStore core.CacheStore, informationProvider core.BasicInformationProvider) core.PublicService {
	return &publicService{
		unitOfWork:          unit,
		mapper:              mapper,
		cacheStore:          cacheStore,
		informationProvider: informationProvider,
	}
}
