package repository

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/dbmodel"
	"time"
)

type hotelRepository struct {
	DB *gorm.DB
}

func (r *hotelRepository) FindByID(hotelId string) (*dbmodel.Hotel, error) {
	var hotel dbmodel.Hotel

	if r.DB.Preload("Amenities").Preload("Places").Preload("FAQList").
		Preload("Badges").Preload("Amenities.AmenityCategory").
		Find(&hotel, "PlaceId=? or NameEn=? COLLATE SQL_Latin1_General_CP1_CS_AS ", hotelId, hotelId).RecordNotFound() {
		return nil, common.HotelNotFound
	}
	return &hotel, nil
}

func (r *hotelRepository) FindByIDForSync(hotelId string) (*dbmodel.Hotel, error) {
	var hotel dbmodel.Hotel

	if r.DB.Preload("Amenities").Preload("Places").Preload("Badges").
		Find(&hotel, "PlaceId=?", hotelId).RecordNotFound() {
		return nil, common.HotelNotFound
	}
	return &hotel, nil
}

func (r *hotelRepository) GetHotel(hotelId string) (*dbmodel.Hotel, error) {
	var hotel dbmodel.Hotel
	if r.DB.
		Find(&hotel, "PlaceId=?", hotelId).RecordNotFound() {
		return nil, common.HotelNotFound
	}
	return &hotel, nil
}

func (r *hotelRepository) StoreOrUpdate(hotel *dbmodel.Hotel) error {
	return r.DB.Save(&hotel).Error
}

func (r *hotelRepository) Delete(hotel dbmodel.Hotel) error {
	db := r.DB.Delete(&hotel)
	if db.RecordNotFound() {
		return common.HotelNotFound
	}
	return db.Error
}

func (r *hotelRepository) GetAllHotelIds() ([]string, error) {
	var hotels []dbmodel.Hotel
	db := r.DB.Select("PlaceId").Find(&hotels)

	hotelIds := make([]string, 0, len(hotels))
	for _, hotel := range hotels {
		hotelIds = append(hotelIds, hotel.PlaceID)
	}

	return hotelIds, db.Error
}

func (r *hotelRepository) GetHotelsPageForSync(page int, size int) ([]dbmodel.Hotel, error) {
	var hotels []dbmodel.Hotel
	db := r.DB.Preload("Amenities").Preload("Badges").Preload("Amenities.AmenityCategory").
		Order("id desc").Limit(size).Offset(size * (page - 1)).Find(&hotels)
	return hotels, db.Error
}

func (r *hotelRepository) GetHotels(ids []string) ([]dbmodel.Hotel, error) {
	var hotels []dbmodel.Hotel
	db := r.DB.Preload("Places").
		Preload("Badges").Preload("Amenities").
		Preload("Amenities.AmenityCategory").Where("PlaceId IN (?)", ids).Find(&hotels)
	return hotels, db.Error
}

func (r *hotelRepository) GetAllHotels() ([]dbmodel.Hotel, error) {
	var hotels []dbmodel.Hotel
	db := r.DB.Select("PlaceId, Price, created_at, updated_at").Find(&hotels)
	return hotels, db.Error
}

func (r *hotelRepository) HasBeenSynced() (bool, error) {
	now := time.Now()
	year, month, day := now.Date()
	startDate := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	var count int64
	r.DB.Model(&dbmodel.Hotel{}).Where("updated_at >= ? AND RoomId != '0'", startDate).Count(&count)
	return count >= 100, r.DB.Error
}

func (r *hotelRepository) GetHotelsList(page int, size int, search string) ([]dbmodel.Hotel, int, error) {
	data := make(chan []dbmodel.Hotel)
	query := r.DB.Model(dbmodel.Hotel{})
	if search != "" {
		query = query.Where("PlaceId like ? or NameEn like ? COLLATE SQL_Latin1_General_CP1_CS_AS", "%"+search+"%", "%"+search+"%")
	}

	go func(channel chan<- []dbmodel.Hotel) {
		var hotels []dbmodel.Hotel
		query.Preload("Amenities").Preload("Badges").
			Preload("Amenities.AmenityCategory").Order("id desc").
			Limit(size).Offset(size * (page - 1)).Find(&hotels)
		channel <- hotels
	}(data)

	var total int
	db := query.Count(&total)
	return <-data, total, db.Error
}

func (r *hotelRepository) RemoveFAQ(hotelId string, faq *dbmodel.FAQ) (*dbmodel.Hotel, error) {
	var hotel dbmodel.Hotel
	err := r.DB.Find(&hotel, "PlaceId=? or NameEn=? COLLATE SQL_Latin1_General_CP1_CS_AS ", hotelId, hotelId).
		Association("FAQList").Delete(faq).Error
	return &hotel, err
}

func newHotelRepository(DB *gorm.DB) core.HotelRepository {
	return &hotelRepository{DB: DB}
}
