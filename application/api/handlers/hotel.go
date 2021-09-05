package handlers

import (
	"errors"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/config"
	_ "hotel-engine/infrastructure/hotelproviderinterface/dtos"
	"hotel-engine/utils/date"
	_ "hotel-engine/utils/indraframework"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HotelHandler interface {
	GetHotelDetails(c *gin.Context)
	FindHotels(c *gin.Context)
	GetRoomCancellationPolicy(c *gin.Context)
	Search(c *gin.Context)
	UpdateAllHotels(c *gin.Context)
	SyncAllHotels(c *gin.Context)
	UpdateAndSyncElastic(c *gin.Context)
	SyncElastic(c *gin.Context)
	SyncSomeHotels(c *gin.Context)
	SyncedHotels(c *gin.Context)
	Rooms(c *gin.Context)
	Info(c *gin.Context)
	RoomsWithSession(c *gin.Context)
	Available(c *gin.Context)
	FinalizeOrder(c *gin.Context)
	OrderDetail(c *gin.Context)
	SetAmenityIcon(c *gin.Context)
	SetAmenityCategory(c *gin.Context)

	ConfirmOrder(c *gin.Context)
	PayByAccount(c *gin.Context)
	GetOrderStatus(c *gin.Context)
	GetOrderEnquiry(c *gin.Context)
	RefundOrder(c *gin.Context)

	GetHotelsList(c *gin.Context)
	GetHotelById(c *gin.Context)
	SetHotelSeoTags(c *gin.Context)

	SetHotelFaq(c *gin.Context)
	DeleteHotelFaq(c *gin.Context)
}

type hotelHandler struct {
	service     core.HotelService
	syncService core.SyncService
}

// GetHotelDetails godoc
// @Summary Get a hotel by id
// @Description find a hotel from the db
// @ID GetHotelDetails
// @tags Hotel
// @Accept  json
// @Produce  json
// @Param pdpRequestDto body dto.HotelDetailsDto true "the request body"
// @Success 200 {object} dto.HotelPDPDto
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/hotel-pdp [post]
func (h *hotelHandler) GetHotelDetails(c *gin.Context) {
	var detailsRequest dto.HotelDetailsDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&detailsRequest), &dto.HotelPDPDto{} },
		func() (error error, data dto.Dto) { return detailsRequest.Validate(), &dto.HotelPDPDto{} }); !success {
		return
	}
	detailsRequest.SetDefaults()
	hotel, err := h.service.GetHotelDetails(detailsRequest)

	if err == common.HotelNotFound {
		jsonNotFound(c, &dto.HotelPDPDto{}, err)
		return
	}

	if err != nil {
		jsonBadRequest(c, &dto.HotelPDPDto{}, err)
		return
	}
	jsonSuccess(c, hotel)
}

// FindHotels godoc
// @Summary Get a hotel by id
// @Description find a hotel from the db
// @ID FindHotels
// @tags Hotel
// @Accept  json
// @Produce  json
// @Param pdpRequestDto body dto.HotelInput true "the request body"
// @Success 200 {object} dto.FindHotelResponseDto
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/find-hotels [post]
func (h *hotelHandler) FindHotels(c *gin.Context) {
	var detailsRequest dto.HotelInput
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&detailsRequest), &dto.FindHotelResponseDto{} }); !success {
		return
	}
	hotel, err := h.service.GetHotels(*detailsRequest.HotelIDs)

	if err == common.HotelNotFound {
		jsonNotFound(c, &dto.FindHotelResponseDto{}, err)
		return
	}

	if err != nil {
		jsonBadRequest(c, &dto.FindHotelResponseDto{}, err)
		return
	}
	jsonSuccess(c, dto.FindHotelResponseDto{Hotels: &hotel})
}

// GetRoomCancellationPolicy godoc
// @Summary Get room cancellation policy
// @Description Get room cancellation policy
// @ID GetRoomCancellationPolicy
// @tags Hotel
// @Accept  json
// @Produce  json
// @Param hotelId path string true "hotel id"
// @Param roomId path string true "room id"
// @Param sessionId path string true "session id"
// @Success 200 {object} dto.RoomCancellationPolicyDto
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/room-cancellation-policy/{hotelId}/{roomId}/{sessionId} [get]
func (h *hotelHandler) GetRoomCancellationPolicy(c *gin.Context) {
	policyDto, err := h.service.GetRoomCancellationPolicy(c.Param("hotelId"),
		c.Param("roomId"), c.Param("sessionId"))
	if err != nil {
		jsonBadRequest(c, &dto.TaskRunningResult{}, err)
		return
	}
	jsonSuccess(c, policyDto)
}

// UpdateAllHotels godoc
// @Summary Update all the hotels
// @Description Update all the hotels from provider
// @ID update-all
// @tags Hotel - management
// @Accept  json
// @Produce  json
// @Param searchDto body dto.UpdateHotelsDto true "the request body"
// @Success 200 {object} dto.TaskRunningResult
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/update-hotels [put]
func (h *hotelHandler) UpdateAllHotels(c *gin.Context) {

	var updateDto dto.UpdateHotelsDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&updateDto), &dto.TaskRunningResult{} },
		func() (error error, data dto.Dto) { return updateDto.Validate(), &dto.TaskRunningResult{} }); !success {
		return
	}
	item, err := h.service.UpdateAll(date.StringToDateOrDefault(updateDto.Date))
	if err != nil {
		jsonBadRequest(c, &dto.TaskRunningResult{}, err)
		return
	}
	//fmt.Printf("update hotels completed for %d hotels in %d nanoseconds", item.ItemsCount, item.TimeTaken)

	jsonSuccess(c, item)
}

// Search godoc
// @Summary search available hotels
// @Description
// @ID Search
// @tags Hotel
// @Accept  json
// @Produce  json
// @Param searchDto body dto.SearchDto true "the request body"
// @Success 200 {object} dto.SearchResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/search [post]
func (h *hotelHandler) Search(c *gin.Context) {

	var searchDto dto.SearchDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&searchDto), &dto.SearchResponseDto{} },
		func() (error error, data dto.Dto) { return searchDto.Validate(), &dto.SearchResponseDto{} }); !success {
		return
	}

	result, err := h.service.SearchResult(searchDto)

	if err != nil {
		jsonBadRequest(c, &dto.SearchResponseDto{}, err)
		return
	}

	jsonSuccess(c, result)
}

// Room godoc
// @Summary get hotel rooms
// @Description get hotel rooms
// @ID room
// @tags Hotel
// @Accept  json
// @Produce  json
// @Param roomDto body dto.HotelRoomsDto true "get hotel rooms"
// @Success 200 {object} dto.RateRoomResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/rooms [post]
func (h *hotelHandler) Rooms(c *gin.Context) {
	var roomDto dto.HotelRoomsDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&roomDto), &dto.RateRoomResponseDto{} },
		func() (error error, data dto.Dto) { return roomDto.Validate(), &dto.RateRoomResponseDto{} }); !success {
		return
	}
	roomDto.SetDefaults()
	res, err := h.service.GetHotelRooms(roomDto)

	if err != nil {
		jsonBadRequest(c, &dto.RateRoomResponseDto{}, err)
		return
	}

	jsonSuccess(c, res)
}

// Info godoc
// @Summary get hotel option Info
// @Description get hotel option Info
// @ID info
// @tags Hotel
// @Accept  json
// @Produce  json
// @Param infoDto body dto.OptionInfoRequestDto true "get hotel option info"
// @Success 200 {object} dto.OptionInfoResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/info [post]
func (h *hotelHandler) Info(c *gin.Context) {
	var infoDto dto.OptionInfoRequestDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&infoDto), &dto.OptionInfoResponseDto{} },
		func() (error error, data dto.Dto) { return infoDto.Validate(), &dto.OptionInfoResponseDto{} }); !success {
		return
	}

	res, err := h.service.GetHotelOptionInfo(infoDto)

	if err != nil {
		jsonBadRequest(c, &dto.OptionInfoResponseDto{}, err)
		return
	}

	jsonSuccess(c, res)
}

// Room godoc
// @Summary get hotel rooms
// @Description get hotel rooms
// @ID RoomsWithSession
// @tags Hotel
// @Accept  json
// @Produce  json
// @Param roomDto body dto.HotelRoomsWithSessionDto true "get hotel rooms"
// @Success 200 {object} dto.RateRoomResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/rooms-by-session [post]
func (h *hotelHandler) RoomsWithSession(c *gin.Context) {
	var roomDto dto.HotelRoomsWithSessionDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&roomDto), &dto.RateRoomResponseDto{} }); !success {
		return
	}

	res, err := h.service.GetHotelRoomsWithSession(roomDto)

	if err != nil {
		jsonBadRequest(c, &dto.RateRoomResponseDto{}, err)
		return
	}

	jsonSuccess(c, res)
}

// Available godoc
// @Summary get order id if hotel is available
// @Description get order id if hotel is available
// @ID available
// @tags Hotel - Order
// @Accept  json
// @Produce  json
// @Param availableDto body dto.AvailableDto true "get order id if hotel is available"
// @Success 200 {object} dto.AvailableResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/available [post]
func (h *hotelHandler) Available(c *gin.Context) {
	var availableDto dto.AvailableDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&availableDto), &dto.AvailableResponseDto{} },
		func() (error error, data dto.Dto) { return availableDto.Validate(), &dto.AvailableResponseDto{} }); !success {
		return
	}

	res, err := h.service.HotelAvailable(availableDto)

	if err != nil {
		jsonBadRequest(c, &dto.AvailableResponseDto{}, err)
		return
	}

	jsonSuccess(c, res)
}

// FinalizeOrder godoc
// @Summary Finalize Order
// @Description Finalize Order
// @ID FinalizeOrder
// @tags Hotel - Order
// @Accept  json
// @Produce  json
// @Param finalizeOrderDto body dto.FinalizeOrderDto true "get hotel rooms"
// @Success 200 {object} dto.FinalizeOrderResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/finalize-order [post]
func (h *hotelHandler) FinalizeOrder(c *gin.Context) {
	var finalizeOrderDto dto.FinalizeOrderDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) {
			return c.BindJSON(&finalizeOrderDto), &dto.FinalizeOrderResponseDto{}
		},
		func() (error error, data dto.Dto) {
			return finalizeOrderDto.Validate(), &dto.FinalizeOrderResponseDto{}
		}); !success {
		return
	}

	res, err := h.service.FinalizeHotelOrder(finalizeOrderDto)

	if err != nil {
		jsonBadRequest(c, &dto.FinalizeOrderResponseDto{}, err)
		return
	}

	jsonSuccess(c, res)
}

// OrderDetail godoc
// @Summary Get Order Detail
// @Description Get order detail
// @ID OrderDetail
// @tags Hotel - Order
// @Accept  json
// @Produce  json
// @Param id path string true "order id"
// @Success 200 {object} dto.OrderDetailDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/order-detail/{id} [get]
func (h *hotelHandler) OrderDetail(c *gin.Context) {
	id, found := c.Params.Get("id")
	if !found {
		jsonBadRequest(c, &dto.OrderDetailDto{}, errors.New("order id is required"))
		return
	}
	detail, err := h.service.GetAnOrderDetail(id)

	if err == common.OrderNotFound {
		jsonNotFound(c, &dto.OrderDetailDto{}, err)
		return
	}

	if err != nil {
		jsonBadRequest(c, &dto.OrderDetailDto{}, err)
		return
	}
	jsonSuccess(c, detail)
}

// SyncAllHotels godoc
// @Summary sync all hotels
// @Description sync all hotels from provider
// @ID sync-all-hotels
// @tags Hotel - management
// @Produce  json
// @Param secret path string true "the sync secret"
// @Success 200 {object} dto.TaskRunningResult
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/sync-hotels/{secret} [put]
func (h *hotelHandler) SyncAllHotels(c *gin.Context) {
	secret := c.Param("secret")
	if config.Get().SyncSecret != secret {
		jsonForbiddenRequest(c, &dto.TaskRunningResult{}, errors.New("secret key is not correct"))
		return
	}
	taskRes, err := h.service.SyncAllHotels()
	if err != nil {
		jsonBadRequest(c, &dto.TaskRunningResult{}, err)
		return
	}
	jsonSuccess(c, taskRes)
}

// UpdateAndSyncElastic godoc
// @Summary Update And Sync Elastic
// @Description Update And Sync Elastic
// @ID update-and-sync-elastic
// @tags Hotel - management
// @Produce  json
// @Param secret path string true "the sync secret"
// @Success 200 {object} dto.TaskRunningResult
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/update-and-sync-elastic/{secret} [put]
func (h *hotelHandler) UpdateAndSyncElastic(c *gin.Context) {
	secret := c.Param("secret")
	if config.Get().SyncSecret != secret {
		jsonForbiddenRequest(c, &dto.TaskRunningResult{}, errors.New("secret key is not correct"))
		return
	}
	go func() {
		h.syncService.UpdateAndSyncElastic()
	}()
	jsonSuccess(c, &dto.TaskRunningResult{
		Message: "updating database and syncing elastic task is now running in background",
		Success: true,
		Error:   nil,
	})
}

// SyncElastic godoc
// @Summary sync elastic
// @Description sync elastic
// @ID sync-elastic
// @tags Hotel - management
// @Produce  json
// @Param secret path string true "the sync secret"
// @Success 200 {object} dto.TaskRunningResult
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/sync-elastic/{secret} [put]
func (h *hotelHandler) SyncElastic(c *gin.Context) {
	secret := c.Param("secret")
	if config.Get().SyncSecret != secret {
		jsonForbiddenRequest(c, &dto.TaskRunningResult{}, errors.New("secret key is not correct"))
		return
	}

	go func() {
		h.syncService.SyncElastic()
	}()

	jsonSuccess(c, &dto.TaskRunningResult{
		Message: "syncing elastic task is now running in background",
		Success: true,
		Error:   nil,
	})
}

// SyncSomeHotels godoc
// @Summary sync some hotels
// @Description sync all given hotels from provider
// @ID sync-some
// @tags Hotel - management
// @Accept  json
// @Produce  json
// @Param searchDto body dto.SyncSomeHotelsDto true "the request body"
// @Success 200 {object} dto.UpdateResultDto
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/sync-some-hotels [put]
func (h *hotelHandler) SyncSomeHotels(c *gin.Context) {

	var syncDto dto.SyncSomeHotelsDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&syncDto), &dto.UpdateResultDto{} },
		func() (error error, data dto.Dto) { return syncDto.Validate(), &dto.UpdateResultDto{} }); !success {
		return
	}

	item, err := h.service.SyncSomeHotels(syncDto)
	if err != nil {
		jsonBadRequest(c, &dto.UpdateResultDto{}, err)
		return
	}
	jsonSuccess(c, item)
}

// SyncedHotels godoc
// @Summary get synced hotels
// @Description get synced hotels
// @ID synced-hotels
// @tags Hotel - management
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.SyncedHotelsDetail
// @Router /v1/hotel/synced-hotel [get]
func (h *hotelHandler) SyncedHotels(c *gin.Context) {
	details, err := h.service.SyncedHotels()
	if err != nil {
		jsonBadRequest(c, &dto.SyncedHotelsDetail{}, err)
		return
	}
	jsonSuccess(c, details)
}

// SetAmenityIcon godoc
// @Summary update amenity icon
// @Description update amenity icon
// @ID SetAmenityIcon
// @tags Hotel - Admin
// @Accept  json
// @Produce  json
// @Param amenityIconDto body dto.SetAmenityIconDto true "the request body"
// @Success 200 {object} dto.HotelAmenityDto
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/set-amenity-icon [put]
func (h *hotelHandler) SetAmenityIcon(c *gin.Context) {

	var body dto.SetAmenityIconDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&body), &dto.HotelAmenityDto{} },
		func() (error error, data dto.Dto) { return body.Validate(), &dto.HotelAmenityDto{} }); !success {
		return
	}

	item, err := h.service.SetAmenityIcon(body)
	if err == common.AmenityNotFound {
		jsonNotFound(c, &dto.HotelAmenityDto{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.HotelAmenityDto{}, err)
		return
	}
	jsonSuccess(c, item)
}

// SetAmenityCategory godoc
// @Summary update amenity category
// @Description update amenity category
// @ID SetAmenityCategory
// @tags Hotel - Admin
// @Accept  json
// @Produce  json
// @Param amenityIconDto body dto.SetAmenityCategoryDto true "the request body"
// @Success 200 {object} dto.HotelAmenityDto
// @Failure 404 {object}  indraframework.IndraException
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/set-amenity-category [put]
func (h *hotelHandler) SetAmenityCategory(c *gin.Context) {

	var body dto.SetAmenityCategoryDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&body), &dto.HotelAmenityDto{} },
		func() (error error, data dto.Dto) { return body.Validate(), &dto.HotelAmenityDto{} }); !success {
		return
	}

	item, err := h.service.SetAmenityCategory(body)
	if err == common.AmenityNotFound {
		jsonNotFound(c, &dto.HotelAmenityDto{}, err)
		return
	}
	if err == common.AmentityCategoryNotFound {
		jsonNotFound(c, &dto.HotelAmenityDto{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.HotelAmenityDto{}, err)
		return
	}
	jsonSuccess(c, item)
}

// ConfirmOrder godoc
// @Summary confirm an order
// @Description confirm an order
// @ID ConfirmOrder
// @tags Hotel - Order
// @Accept  json
// @Produce  json
// @Param orderId path string true "order id"
// @Success 200 {object} dto.ConfirmResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/order/confirm/{orderId} [put]
func (h *hotelHandler) ConfirmOrder(c *gin.Context) {
	orderId := c.Param("orderId")
	res, err := h.service.ConfirmOrder(orderId)

	if err == common.OrderNotFound {
		jsonNotFound(c, &dto.ConfirmResponseDto{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.ConfirmResponseDto{}, err)
		return
	}
	jsonSuccess(c, res)
}

// PayByAccount godoc
// @Summary pay order by jabama account
// @Description pay order by jabama account
// @ID PayByAccount
// @tags Hotel - Order
// @Accept  json
// @Produce  json
// @Param orderId path string true "order id"
// @Success 200 {object} dto.OrderPayByAccountResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/order/pay-by-account/{orderId} [put]
func (h *hotelHandler) PayByAccount(c *gin.Context) {
	orderId := c.Param("orderId")
	res, err := h.service.PayByAccount(orderId)

	if err == common.OrderNotFound {
		jsonNotFound(c, &dto.OrderPayByAccountResponseDto{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.OrderPayByAccountResponseDto{}, err)
		return
	}
	jsonSuccess(c, res)
}

// GetOrderStatus godoc
// @Summary get order status
// @Description get order status
// @ID GetOrderStatus
// @tags Hotel - Order
// @Accept  json
// @Produce  json
// @Param orderId path string true "order id"
// @Success 200 {object} dto.OrderStatusResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/order/status/{orderId} [get]
func (h *hotelHandler) GetOrderStatus(c *gin.Context) {
	orderId := c.Param("orderId")
	res, err := h.service.GetOrderStatus(orderId)

	if err == common.OrderNotFound {
		jsonNotFound(c, &dto.OrderStatusResponseDto{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.OrderStatusResponseDto{}, err)
		return
	}
	jsonSuccess(c, res)
}

// GetOrderEnquiry godoc
// @Summary get order enquiry
// @Description get order enquiry
// @ID GetOrderEnquiry
// @tags Hotel - Order
// @Accept  json
// @Produce  json
// @Param orderId path string true "order id"
// @Success 200 {object} dto.OrderEnquiryResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/order/enquiry/{orderId} [get]
func (h *hotelHandler) GetOrderEnquiry(c *gin.Context) {
	orderId := c.Param("orderId")
	res, err := h.service.GetOrderEnquiry(orderId)

	if err == common.OrderNotFound {
		jsonNotFound(c, &dto.OrderEnquiryResponseDto{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.OrderEnquiryResponseDto{}, err)
		return
	}
	jsonSuccess(c, res)
}

// RefundOrder godoc
// @Summary refund an order
// @Description refund an order
// @ID RefundOrder
// @tags Hotel - Order
// @Accept  json
// @Produce  json
// @Param infoDto body dto.OrderRefundRequestDto true "refund request dto"
// @Success 200 {object} dto.OrderRefundResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/order/refund [post]
func (h *hotelHandler) RefundOrder(c *gin.Context) {
	var body dto.OrderRefundRequestDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&body), &dto.OrderRefundResponseDto{} },
		func() (error error, data dto.Dto) { return body.Validate(), &dto.OrderRefundResponseDto{} }); !success {
		return
	}

	item, err := h.service.RefundOrder(body)
	if err == common.OrderNotFound {
		jsonNotFound(c, &dto.OrderRefundResponseDto{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.OrderRefundResponseDto{}, err)
		return
	}
	jsonSuccess(c, item)
}

// GetHotelsList godoc
// @Summary get hotel lists
// @Description get hotel lists
// @ID GetHotelsList
// @tags Hotel - Admin
// @Accept  json
// @Produce  json
// @Param infoDto body dto.HotelsPageRequestDto true "hotels list request dto"
// @Success 200 {object} dto.HotelsPageResponseDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/list [post]
func (h *hotelHandler) GetHotelsList(c *gin.Context) {
	var body dto.HotelsPageRequestDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&body), &dto.HotelsPageResponseDto{} },
		func() (error error, data dto.Dto) { return body.Validate(), &dto.HotelsPageResponseDto{} }); !success {
		return
	}

	res, err := h.service.GetHotelsList(body)
	if err != nil {
		jsonBadRequest(c, &dto.HotelsPageResponseDto{}, err)
		return
	}
	jsonSuccess(c, res)
}

// GetHotelById godoc
// @Summary get hotel by id
// @Description get hotel by id
// @ID GetHotelById
// @tags Hotel - Admin
// @Accept  json
// @Produce  json
// @Param hotelId path string true "hotel id"
// @Success 200 {object} dto.HotelDto
// @Failure 400 {object} indraframework.IndraException
// @Router /v1/hotel/find/{hotelId} [get]
func (h *hotelHandler) GetHotelById(c *gin.Context) {
	res, err := h.service.FindHotelById(c.Param("hotelId"))
	if err != nil {
		jsonBadRequest(c, &dto.HotelDto{}, err)
		return
	}
	jsonSuccess(c, res)
}

// SetHotelSeoTags godoc
// @Summary set hotel meta tags
// @Description set hotel meta tags
// @ID SetHotelSeoTags
// @tags Hotel - Admin
// @Accept  json
// @Produce  json
// @Param hotelSeoRequestDto body dto.SetHotelSeoRequestDto true "the request body"
// @Success 200 {object} dto.HotelDto
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/set-hotel-meta-tags [put]
func (h *hotelHandler) SetHotelSeoTags(c *gin.Context) {

	var hotelSeoRequestDto dto.SetHotelSeoRequestDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&hotelSeoRequestDto), &dto.HotelDto{} },
		func() (error error, data dto.Dto) { return hotelSeoRequestDto.Validate(), &dto.HotelDto{} }); !success {
		return
	}

	item, err := h.service.SetHotelSeoDetails(hotelSeoRequestDto)
	if err == common.HotelNotFound {
		jsonNotFound(c, &dto.HotelDto{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.HotelDto{}, err)
		return
	}
	jsonSuccess(c, item)
}

// SetHotelFaq godoc
// @Summary set hotel faq
// @Description set hotel faq
// @ID SetHotelFaq
// @tags Hotel - Admin
// @Accept  json
// @Produce  json
// @Param hotelFaqRequestDto body dto.SetHotelFaqRequestDto true "the request body"
// @Success 200 {object} dto.HotelDto
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/set-hotel-faq [put]
func (h *hotelHandler) SetHotelFaq(c *gin.Context) {

	var hotelFaqRequestDto dto.SetHotelFaqRequestDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&hotelFaqRequestDto), &dto.HotelDto{} },
		func() (error error, data dto.Dto) { return hotelFaqRequestDto.Validate(), &dto.HotelDto{} }); !success {
		return
	}

	item, err := h.service.SetHotelFaq(hotelFaqRequestDto)
	if err == common.HotelNotFound {
		jsonNotFound(c, &dto.HotelDto{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.HotelDto{}, err)
		return
	}
	jsonSuccess(c, item)
}

// DeleteHotelFaq godoc
// @Summary delete hotel faq
// @Description delete hotel faq
// @ID DeleteHotelFaq
// @tags Hotel - Admin
// @Accept  json
// @Produce  json
// @Param hotelId path string true "hotel id"
// @Param faqId path integer true "faq id"
// @Success 200 {object} dto.HotelDto
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/delete-hotel-faq/{hotelId}/{faqId} [delete]
func (h *hotelHandler) DeleteHotelFaq(c *gin.Context) {
	hotelId := c.Param("hotelId")
	faqIdString := c.Param("faqId")
	faqId, err := strconv.Atoi(faqIdString)
	if err != nil {
		jsonBadRequest(c, &dto.HotelDto{}, err)
		return
	}

	item, err := h.service.RemoveHotelFaq(hotelId, uint(faqId))
	if err == common.HotelNotFound || err == common.FAQNotFound {
		jsonNotFound(c, &dto.HotelDto{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.HotelDto{}, err)
		return
	}
	jsonSuccess(c, item)
}

func NewHotelHandler(service core.HotelService, syncService core.SyncService) HotelHandler {
	return &hotelHandler{service: service,
		syncService: syncService}
}
