package repository

import (
	"github.com/jinzhu/gorm"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/dbmodel"
	"strconv"
	"time"
)

type orderRepository struct {
	DB *gorm.DB
}

func (r *orderRepository) getHotel(hotelId string) (*dbmodel.Hotel, error) {
	var hotel dbmodel.Hotel
	if r.DB.
		Find(&hotel, "PlaceId=?", hotelId).RecordNotFound() {
		return nil, common.HotelNotFound
	}
	return &hotel, nil
}

func (r *orderRepository) Insert(order dbmodel.Order) (*dbmodel.Order, error) {
	hotel, err := r.getHotel(order.ProviderHotelId)
	if err != nil {
		return nil, err
	}
	order.HotelID = hotel.ID
	result := r.DB.Save(&order)
	return &order, result.Error
}

func (r *orderRepository) GetOneByIndraId(indraId string) (*dbmodel.Order, error) {
	var order dbmodel.Order
	if r.DB.Preload("Rooms").
		Find(&order, "IndraOrderId=?", indraId).RecordNotFound() {
		return nil, common.OrderNotFound
	}
	return &order, nil
}

func (r *orderRepository) StoreOrUpdate(order dbmodel.Order) error {
	return r.DB.Save(&order).Error
}

func (r *orderRepository) GetProperOrderIdsForRefundUpdateStatus(
	fromDate time.Time) ([]string, error) {

	var orders []dbmodel.Order
	db := r.DB.Select("IndraOrderId").
		Where("RefundRequestId != '' and updated_at > ? and RefundStatus != ?", fromDate, common.RefundStatus_PaymentFinalized).
		Order("id desc").Find(&orders)

	indraOrderIds := make([]string, 0, len(orders))
	for _, order := range orders {
		indraOrderIds = append(indraOrderIds, strconv.FormatInt(order.IndraOrderId, 10))
	}

	return indraOrderIds, db.Error
}

func newOrderRepository(DB *gorm.DB) core.OrderRepository {
	return &orderRepository{DB: DB}
}
