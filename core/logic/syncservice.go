package logic

import (
	"fmt"
	"hotel-engine/core"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dbmodel"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/logger"
	"strconv"
	"strings"
	"time"
)

type syncService struct {
	service    core.HotelService
	unitOfWork core.UnitOfWork
	feeder     core.Feeder
}

func (s *syncService) UpdateDatabase() {
	logger.Debug("updating database hotels data")
	res, err := s.service.UpdateAllSync(time.Now())
	if err != nil {
		logger.WithName(logtags.UpdatingHotelsJobError).WithDevMessage("sync hotels job -> updateDatabase").
			ErrorException(err, "error while updating database hotels")
	}
	logger.WithName(logtags.UpdatingHotelsJobCompleted).
		Info(fmt.Sprintf("updating database completed in %d nanoseconds and updated %d hotels",
			res.TimeTaken, res.ItemsCount))
}

func (s *syncService) SyncElastic() {
	var size = 100
	var page = 1
	for {
		hotels, err := s.unitOfWork.Hotel().GetHotelsPageForSync(page, size)
		if err != nil {
			logger.WithName(logtags.GettingListOfHotelsError).
				WithDevMessage("sync hotels job -> syncElastic").
				WithData(map[string]interface{}{
					"page": page,
					"size": size,
				}).ErrorException(err, "error while getting hotels page")
			page++
			continue
		}
		err = s.feeder.Feed(mapper(hotels))
		if err != nil {
			logger.WithName(logtags.FeedElasticError).WithDevMessage("sync hotels service -> syncElastic -> feed").
				ErrorException(err, "error while feeding elastic")
		}
		if len(hotels) < size {
			break
		}
		page++
	}
}

func (s *syncService) UpdateAndSyncElastic() {
	s.UpdateDatabase()
	s.SyncElastic()
}

func (s *syncService) feed(index string) {
	s.SyncElastic()
	err := s.feeder.Alias(index)
	if err != nil {
		logger.WithName(logtags.CallingAliasError).WithDevMessage("sync hotels service -> feed(index string) -> alias").
			ErrorException(err, "error while calling alias")
	}
}

func mapper(hotels []dbmodel.Hotel) dto.ElasticUpdateRequest {
	var docs = []dto.ElasticHotel{}
	for _, hotel := range hotels {

		//Amenities
		var amenities []dto.ElasticAmenity
		for _, amenity := range hotel.Amenities {
			categoryId := ""
			categoryName := ""
			if amenity.AmenityCategory != nil {
				categoryId = fmt.Sprintf("%d", amenity.AmenityCategory.ID)
				categoryName = amenity.AmenityCategory.NameEn
			}
			amenities = append(amenities, dto.ElasticAmenity{
				Category:     categoryId,
				CategoryName: categoryName,
				Name:         amenity.NameEn,
			})
		}

		var badges []dto.ElasticBadge
		for _, badge := range hotel.Badges {
			badges = append(badges, dto.ElasticBadge{
				Icon: badge.Icon,
				Name: badge.Text,
			})
		}

		//Geo
		lat := float64(0)
		lon := float64(0)
		geo := strings.Split(hotel.GeoLocation, ",")
		if len(geo) == 2 {
			lat, _ = strconv.ParseFloat(geo[0], 32)
			lon, _ = strconv.ParseFloat(geo[1], 32)
		}

		//Calendar
		year, _month, day := hotel.CheckIn.Date()
		month, _ := strconv.ParseInt(fmt.Sprintf("%02d", _month), 10, 32)
		date, _ := strconv.ParseInt(fmt.Sprintf("%d%02d%02d", year, int(month), day), 10, 32)

		//Images
		image := ""
		images := strings.Split(hotel.Images, ",")
		if len(images) > 0 {
			image = images[0]
		} else {
			images = []string{}
		}
		hotelCalendar := []dto.ElasticCalendar{}

		if hotel.RoomID != "0" {
			hotelCalendar = append(hotelCalendar, dto.ElasticCalendar{
				Available: int(date),
				Date:      int(date),
				Capacity: dto.ElasticCalendarCapacity{
					Base: 1,
				},
				Price: int(hotel.Price),
				Day:   day,
				Month: int(month),
				Year:  year,
			})
		}

		hotelSort := hotel.Sort
		if hotel.Price == 0 {
			hotelSort = 0
		}

		doc := dto.ElasticHotel{
			ID:              fmt.Sprintf("ABH_%s", hotel.PlaceID),
			PlaceID:         fmt.Sprintf("%s", hotel.PlaceID),
			RoomID:          fmt.Sprintf("%s", hotel.RoomID),
			Name:            hotel.Name,
			NameEn:          hotel.NameEn,
			Kind:            "hotel",
			Type:            hotel.Type,
			Status:          "confirmed",
			Description:     "",
			Calendar:        hotelCalendar,
			MinPrice:        int(hotel.Price),
			OldPrice:        hotel.OldPrice,
			DiscountPercent: hotel.DiscountPercent,
			DiscountPrice:   hotel.DiscountPrice,
			MinNight:        1,
			PaymentType:     "full",
			RateReview: dto.ElasticRateReview{
				Score: hotel.RateReviewScore,
				Count: hotel.RateReviewCount,
			},
			Region:          "hotel",
			Tags:            []string{},
			SuitableFor:     []string{},
			Code:            hotel.Code,
			Verified:        true,
			ReservationType: "instant",
			Amenities:       amenities,
			Badges:          badges,
			Location: dto.ElasticLocation{
				City:       hotel.City,
				CityEn:     hotel.CityEn,
				Province:   hotel.Province,
				ProvinceEn: hotel.ProvinceEn,
				Geo: dto.ElasticLocationGeo{
					Lat: lat,
					Lon: lon,
				},
			},
			Image:  image,
			Images: images,
			Star:   hotel.Star,
			Sort:   hotelSort,
		}
		docs = append(docs, doc)
	}

	return dto.ElasticUpdateRequest{
		Places:     docs,
		ClearCache: 0,
	}
}

func (s *syncService) HasBeenSynced() (bool, error) {
	return s.unitOfWork.Hotel().HasBeenSynced()
}

func NewSyncService(feeder core.Feeder, unit core.UnitOfWork, service core.HotelService) core.SyncService {
	syncService := &syncService{
		unitOfWork: unit,
		service:    service,
		feeder:     feeder,
	}
	err := feeder.Seed(syncService.feed)
	if err != nil {
		logger.WithName(logtags.CreatingSeederError).WithDevMessage("sync hotels service -> create seeder").
			ErrorException(err, "error while creating seeder")
	}
	return syncService
}
