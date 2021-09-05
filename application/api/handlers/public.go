package handlers

import (
	"errors"
	"fmt"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/config"
	"hotel-engine/infrastructure/health"
	hotelProviderInterface "hotel-engine/infrastructure/hotelproviderinterface"
	"hotel-engine/infrastructure/logger"
	"hotel-engine/utils/parser"
	"net/http"
	"time"

	_ "hotel-engine/core/dto"
	_ "hotel-engine/utils/indraframework"

	"github.com/gin-gonic/gin"
)

type PublicHandler interface {
	Health(c *gin.Context)
	Amenities(c *gin.Context)
	Cities(c *gin.Context)
	Sorts(c *gin.Context)
	Places(c *gin.Context)
	ChildAgeRanges(c *gin.Context)
	Filters(c *gin.Context)
	LoadInnerTypes(c *gin.Context)
	SyncCities(c *gin.Context)
	Info(c *gin.Context)
	GetProviderToken(c *gin.Context)
	GetProviderBalance(context *gin.Context)

	CreateAmenityCategory(context *gin.Context)
	GetAmenityCategory(context *gin.Context)
	DeleteAmenityCategory(context *gin.Context)
	GetAmenityCategories(context *gin.Context)
	UpdateAmenityCategory(context *gin.Context)

	GetBadges(context *gin.Context)
	UpdateBadgeIcon(context *gin.Context)
}

type publicHandler struct {
	service        core.PublicService
	balanceService core.ProviderBalanceChecker
}

// Health godoc
// @Summary Api Health check
// @Description Api Health check
// @tags Health
// @ID health
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Router /health [get]
func (h *publicHandler) Health(c *gin.Context) {
	c.JSON(200, health.GetHealth())
}

// Info godoc
// @Summary Api Info
// @Description Info
// @tags Health
// @ID Info
// @Accept  json
// @Produce  json
// @Success 200 {string} string
// @Router /info [get]
func (h *publicHandler) Info(c *gin.Context) {
	con := config.Get()
	c.JSON(http.StatusOK, map[string]interface{}{
		"service":     "Hotel wrapper",
		"environment": con.Environment,
		"appName":     con.ServiceName,
		"time":        time.Now(),
	})
}

// Amenities godoc
// @Summary get all amenities
// @Description get all amenities
// @ID Amenities
// @tags Public
// @Produce  json
// @Success 200 {array} dto.HotelAmenityDto
// @Router /v1/public/amenities [get]
func (h *publicHandler) Amenities(c *gin.Context) {
	amenities := h.service.GetAllAmenities()
	jsonSuccess(c, amenities)
}

// Cities godoc
// @Summary get all cities
// @Description get all cities
// @ID Cities
// @tags Public
// @Produce  json
// @Success 200 {array} dto.CityDto
// @Router /v1/public/cities [get]
func (h *publicHandler) Cities(c *gin.Context) {
	cities := h.service.GetAllCities()
	jsonSuccess(c, cities)
}

// Places godoc
// @Summary get all Places
// @Description get all Places
// @ID Places
// @tags Public
// @Produce  json
// @Success 200 {array} dto.HotelPlaceDto
// @Router /v1/public/places [get]
func (h *publicHandler) Places(c *gin.Context) {
	places := h.service.GetAllPlaces()
	jsonSuccess(c, places)
}

// GetProviderToken godoc
// @Summary get provider token
// @Description get provider token
// @ID GetProviderToken
// @tags Public
// @Produce  json
// @Success 200 {object} dto.TokenDto
// @Router /v1/public/provider-token [get]
func (h *publicHandler) GetProviderToken(c *gin.Context) {
	token, err := hotelProviderInterface.Token()
	if err != nil {
		jsonInternalServerError(c, dto.NewTokenDto(token), err)
		return
	}
	jsonSuccess(c, dto.NewTokenDto(token))
}

// ChildAgeRanges godoc
// @Summary get all ChildAgeRanges
// @Description get all ChildAgeRanges
// @ID ChildAgeRanges
// @tags Public
// @Produce  json
// @Success 200 {array} dto.ChildAgeRangeDto
// @Router /v1/public/child-age-ranges [get]
func (h *publicHandler) ChildAgeRanges(c *gin.Context) {
	ranges := h.service.GetChildAgeRanges()
	jsonSuccess(c, ranges)
}

// Sorts godoc
// @Summary get all Sorts
// @Description get all Sorts
// @ID Sorts
// @tags Public
// @Produce  json
// @Success 200 {array} dto.SortDto
// @Router /v1/public/sorts [get]
func (h *publicHandler) Sorts(c *gin.Context) {
	sorts := h.service.GetAllSorts()
	jsonSuccess(c, sorts)
}

// Filters godoc
// @Summary get all Filters
// @Description get all Filters
// @ID Filters
// @tags Public
// @Produce  json
// @Success 200 {object} dto.FiltersDto
// @Router /v1/public/filters [get]
func (h *publicHandler) Filters(c *gin.Context) {
	filters := h.service.GetAllFilters()
	jsonSuccess(c, filters)
}

// SyncCities godoc
// @Summary sync cities
// @Description sync all cities
// @ID sync-cities
// @tags Hotel - management
// @Produce  json
// @Param secret path string true "the sync secret"
// @Success 200 {array} dto.CityDto
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/hotel/sync-cities/{secret} [put]
func (h publicHandler) SyncCities(c *gin.Context) {
	if config.Get().SyncSecret != c.Param("secret") {
		jsonForbiddenRequest(c, &dto.UpdateResultDto{}, errors.New("secret key is not correct"))
	}
	go func() {
		_, err := h.service.SyncAllCities()
		if err != nil {
			logger.WithName(logtags.SyncingCitiesError).
				ErrorException(err, "error wile syncing cities")
		}
		fmt.Printf("sync cities completed")
	}()

	jsonSuccess(c, &dto.UpdateResultDto{
		Message:    "app trying to sync cities",
		TimeTaken:  0,
		ItemsCount: 0,
		Error:      nil,
	})
}

// GetProviderBalance godoc
// @Summary get provider balance
// @Description get provider balance
// @ID GetProviderBalance
// @tags Public
// @Produce  json
// @Success 200 {object} dto.BalanceDto
// @Router /v1/public/provider-balance [get]
func (h *publicHandler) GetProviderBalance(c *gin.Context) {
	balance, err := h.balanceService.GetBalance()
	if err != nil {
		jsonBadRequest(c, &dto.BalanceDto{}, err)
		return
	}
	jsonSuccess(c, dto.NewBalanceDto(balance))
}

// CreateAmenityCategory godoc
// @Summary create amenity category
// @Description create amenity category
// @ID CreateAmenityCategory
// @tags Public - Amenity Category
// @Produce  json
// @Accept  json
// @Param catDto body dto.AmenityCategoryDto true "create category dto"
// @Success 200 {object} dto.AmenityCategoryDto
// @Failure 400 {object} indraframework.IndraException
// @Failure 404 {object} indraframework.IndraException
// @Router /v1/public/amenity-category [post]
func (h *publicHandler) CreateAmenityCategory(c *gin.Context) {
	var catDto dto.AmenityCategoryDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&catDto), &dto.AmenityCategoryDto{} },
		func() (error error, data dto.Dto) { return catDto.Validate(), &dto.AmenityCategoryDto{} }); !success {
		return
	}

	cat, err := h.service.CreateAmenityCategory(catDto)
	if err == common.AmentityCategoryNotFound {
		jsonNotFound(c, &dto.AmenityCategoryDto{}, common.AmentityCategoryNotFound)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.AmenityCategoryDto{}, err)
		return
	}
	jsonSuccess(c, cat)
}

// GetAmenityCategory godoc
// @Summary get amenity category
// @Description get amenity category
// @ID GetAmenityCategory
// @tags Public - Amenity Category
// @Produce  json
// @Param id path string true "category id"
// @Success 200 {object} dto.AmenityCategoryDto
// @Failure 400 {object} indraframework.IndraException
// @Failure 404 {object} indraframework.IndraException
// @Router /v1/public/amenity-category/{id} [get]
func (h *publicHandler) GetAmenityCategory(c *gin.Context) {
	idString := c.Param("id")
	id, err := parser.ParseNumber(idString)
	if err != nil {
		jsonBadRequest(c, &dto.AmenityCategoryDto{}, err)
		return
	}
	cat, err := h.service.GetAmenityCategory(uint(id))
	if err == common.AmentityCategoryNotFound {
		jsonNotFound(c, &dto.AmenityCategoryDto{}, common.AmentityCategoryNotFound)
		return
	}
	jsonSuccess(c, cat)
}

// DeleteAmenityCategory godoc
// @Summary delete amenity category
// @Description delete amenity category
// @ID DeleteAmenityCategory
// @tags Public - Amenity Category
// @Produce  json
// @Param id path string true "category id"
// @Success 200 {object} dto.DeleteAmenityCategoryResponse
// @Failure 400 {object} indraframework.IndraException
// @Failure 404 {object} indraframework.IndraException
// @Router /v1/public/amenity-category/{id} [delete]
func (h *publicHandler) DeleteAmenityCategory(c *gin.Context) {
	idString := c.Param("id")
	id, err := parser.ParseNumber(idString)
	if err != nil {
		jsonBadRequest(c, &dto.DeleteAmenityCategoryResponse{}, err)
		return
	}
	err = h.service.DeleteAmenityCategory(id)
	if err != nil {
		jsonBadRequest(c, &dto.DeleteAmenityCategoryResponse{}, err)
		return
	}
	jsonSuccess(c, dto.NewDeleteAmenityCategoryResponse(uint(id)))
}

// GetAmenityCategories godoc
// @Summary get amenity categories
// @Description delete amenity categories
// @ID GetAmenityCategories
// @tags Public - Amenity Category
// @Produce  json
// @Success 200 {object} dto.AmenityCategoriesResponse
// @Router /v1/public/amenity-category [get]
func (h *publicHandler) GetAmenityCategories(c *gin.Context) {
	cats := h.service.GetAmenityCategories()
	jsonSuccess(c, dto.NewAmenityCategoriesResponse(cats))
}

// UpdateAmenityCategory godoc
// @Summary update amenity category
// @Description update amenity category
// @ID UpdateAmenityCategory
// @tags Public - Amenity Category
// @Produce  json
// @Param catDto body dto.AmenityCategoryDto true "update category dto"
// @Success 200 {object} dto.AmenityCategoryDto
// @Failure 400 {object} indraframework.IndraException
// @Failure 404 {object} indraframework.IndraException
// @Router /v1/public/amenity-category [put]
func (h *publicHandler) UpdateAmenityCategory(c *gin.Context) {
	var catDto dto.AmenityCategoryDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&catDto), &dto.AmenityCategoryDto{} },
		func() (error error, data dto.Dto) { return catDto.Validate(), &dto.AmenityCategoryDto{} }); !success {
		return
	}

	cat, err := h.service.UpdateAmenityCategory(catDto)
	if err == common.AmentityCategoryNotFound {
		jsonNotFound(c, &dto.AmenityCategoryDto{}, common.AmentityCategoryNotFound)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.AmenityCategoryDto{}, err)
		return
	}
	jsonSuccess(c, cat)
}

// Badges godoc
// @Summary get all badges
// @Description get all badges
// @ID Badges
// @tags Public - Badges
// @Produce  json
// @Success 200 {array} dto.HotelBadgeDto
// @Router /v1/public/badges [get]
func (h *publicHandler) GetBadges(context *gin.Context) {
	badges := h.service.GetAllBadges()
	jsonSuccess(context, badges)
}

// UpdateBadgeIcon godoc
// @Summary Update badge icon
// @Description Update badge icon
// @ID SetBadgeIcon
// @tags Public - Badges
// @Accept  json
// @Produce  json
// @Param badgeIconDto body dto.SetBadgeIconDto true "the request body"
// @Success 200 {object} dto.HotelBadgeDto
// @Failure 400 {object}  indraframework.IndraException
// @Router /v1/public/set-badge-icon [put]
func (h *publicHandler) UpdateBadgeIcon(c *gin.Context) {

	var body dto.SetBadgeIconDto
	if success := tryActions(c,
		func() (error error, data dto.Dto) { return c.BindJSON(&body), &dto.HotelBadgeDto{} },
		func() (error error, data dto.Dto) { return body.Validate(), &dto.HotelBadgeDto{} }); !success {
		return
	}

	item, err := h.service.SetBadgeIcon(body)
	if err == common.BadgeNotFound {
		jsonNotFound(c, &dto.HotelBadgeDto{}, err)
		return
	}
	if err != nil {
		jsonBadRequest(c, &dto.HotelBadgeDto{}, err)
		return
	}
	jsonSuccess(c, item)
}

// LoadInnerTypes godoc
// @Summary this method does not do anything. it just contains the inner objects of other method's responses
// @Description this method does not do anything. it just contains the inner objects of other method's responses
// @ID LoadInnerTypes
// @Accept  json
// @Produce  json
// @Param username path string true "UserName"
// @Success 200 {object}  dto.AvailableRoomDto
// @Success 400 {object}  dto.HotelPlaceDto
// @Success 401 {object}  dto.RoomOptionDto
// @Failure 402 {object}  dto.SearchDateDto
// @Failure 403 {object}  dto.SearchPriceDto
// @Failure 404 {object}  dto.SearchScoreDto
// @Failure 405 {object}  dto.FilterOutput
// @Failure 406 {object}  dto.FilterDto
// @Failure 407 {object}  dto.AvailableRoomAdultDto
// @Failure 408 {object}  dto.AvailableRoomChildrenDto
// @Failure 409 {object}  dto.AvailableRoomPassportDto
// @Failure 410 {object}  dto.SearchResponseHotelDto
// @Failure 411 {object}  dto.RequestRoomDto
// @Failure 412 {object}  dto.SearchResponseLocationDto
// @Failure 413 {object}  dto.SearchResponseRateReviewDto
// @Failure 414 {object}  dto.SearchResponseLocationGeoDto
// @Failure 415 {object}  dto.RoomDto
// @Failure 416 {object}  dto.MealPlanTypeDto
// @Failure 417 {object}  dto.OptionInfoDetailDto
// @Failure 418 {object}  dto.OptionInfoDetailRoomDto
// @Failure 419 {object}  dto.OptionInfoDetailRestrictedMarkupDto
// @Failure 420 {object}  dto.OptionCancellationDto
// @Failure 421 {object}  dto.OrderDetailRoomDto
// @Failure 422 {object}  dto.HotelDto
// @Failure 423 {object}  dto.OrderEnquiryItemDto
// @Failure 424 {object}  dto.OrderEnquiryItemOptionDto
// @Failure 425 {object}  dto.OrderEnquiryItemOptionPassengerInformation
// @Failure 426 {object}  dto.HotelSyncDetail
// @Failure 427 {object}  dto.HotelBadgeDto
// @Failure 428 {object}  dto.SearchResponseBadgeDto
// @Failure 429 {object}  dto.HotelSeoDto
// @Failure 430 {object}  dto.HotelFAQDto
// @Failure 431 {object}  dto.HotelFAQDetailsDto
// @Router /v1/load-inner-types [get]
func (h *publicHandler) LoadInnerTypes(c *gin.Context) {
	c.JSON(http.StatusOK, "tammam")
}

func NewPublicHandler(service core.PublicService, balanceService core.ProviderBalanceChecker) PublicHandler {
	return &publicHandler{
		service:        service,
		balanceService: balanceService,
	}
}
