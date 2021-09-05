package api

import (
	"github.com/gin-gonic/gin"
	"hotel-engine/application/api/handlers"
	"hotel-engine/application/middlewares"
	"net/http"

	"github.com/gin-contrib/pprof"
	_ "net/http/pprof"
)

func CreateRoute(hotelHandler handlers.HotelHandler, publicHandler handlers.PublicHandler) *gin.Engine {

	route := gin.Default()

	route.Use(middlewares.GinBodyLogMiddleware)

	hotelV1 := route.Group("v1/hotel")
	{
		hotelV1.POST("/hotel-pdp", hotelHandler.GetHotelDetails)
		hotelV1.POST("/find-hotels", hotelHandler.FindHotels)
		hotelV1.GET("/room-cancellation-policy/:hotelId/:roomId/:sessionId",
			hotelHandler.GetRoomCancellationPolicy)
		hotelV1.POST("/search", hotelHandler.Search)

		hotelV1.PUT("/update-hotels", hotelHandler.UpdateAllHotels)
		hotelV1.PUT("/sync-hotels/:secret", hotelHandler.SyncAllHotels)
		hotelV1.PUT("/sync-some-hotels", hotelHandler.SyncSomeHotels)
		hotelV1.GET("/synced-hotel", hotelHandler.SyncedHotels)
		hotelV1.PUT("/update-and-sync-elastic/:secret", hotelHandler.UpdateAndSyncElastic)
		hotelV1.PUT("/sync-elastic/:secret", hotelHandler.SyncElastic)

		hotelV1.POST("/rooms", hotelHandler.Rooms)
		hotelV1.POST("/info", hotelHandler.Info)
		hotelV1.POST("/rooms-by-session", hotelHandler.RoomsWithSession)
		hotelV1.POST("/available", hotelHandler.Available)
		hotelV1.POST("/finalize-order", hotelHandler.FinalizeOrder)
		hotelV1.PUT("/order/confirm/:orderId", hotelHandler.ConfirmOrder)
		hotelV1.PUT("/order/pay-by-account/:orderId", hotelHandler.PayByAccount)
		hotelV1.GET("/order/status/:orderId", hotelHandler.GetOrderStatus)
		hotelV1.GET("/order/enquiry/:orderId", hotelHandler.GetOrderEnquiry)
		hotelV1.POST("/order/refund", hotelHandler.RefundOrder)
		hotelV1.GET("/order-detail/:id", hotelHandler.OrderDetail)

		hotelV1.PUT("/set-amenity-icon", hotelHandler.SetAmenityIcon)
		hotelV1.PUT("/set-amenity-category", hotelHandler.SetAmenityCategory)

		hotelV1.POST("/list", hotelHandler.GetHotelsList)
		hotelV1.GET("/find/:hotelId", hotelHandler.GetHotelById)
		hotelV1.PUT("/set-hotel-meta-tags", hotelHandler.SetHotelSeoTags)

		hotelV1.PUT("/set-hotel-faq", hotelHandler.SetHotelFaq)
		hotelV1.DELETE("/delete-hotel-faq/:hotelId/:faqId", hotelHandler.DeleteHotelFaq)
	}

	publicV1 := route.Group("v1/public")
	{
		publicV1.GET("/amenities", publicHandler.Amenities)
		publicV1.GET("/cities", publicHandler.Cities)
		publicV1.GET("/filters", publicHandler.Filters)
		publicV1.GET("/sorts", publicHandler.Sorts)
		publicV1.GET("/places", publicHandler.Places)
		publicV1.GET("/child-age-ranges", publicHandler.ChildAgeRanges)
		publicV1.GET("/load-inner-types", publicHandler.LoadInnerTypes)
		publicV1.GET("/provider-token", publicHandler.GetProviderToken)
		publicV1.GET("/provider-balance", publicHandler.GetProviderBalance)

		hotelV1.PUT("/sync-cities/:secret", publicHandler.SyncCities)

		publicV1.POST("/amenity-category", publicHandler.CreateAmenityCategory)
		publicV1.GET("/amenity-category/:id", publicHandler.GetAmenityCategory)
		publicV1.GET("/amenity-category", publicHandler.GetAmenityCategories)
		publicV1.PUT("/amenity-category", publicHandler.UpdateAmenityCategory)
		publicV1.DELETE("/amenity-category/:id", publicHandler.DeleteAmenityCategory)

		publicV1.GET("/badges", publicHandler.GetBadges)
		publicV1.PUT("/set-badge-icon", publicHandler.UpdateBadgeIcon)
	}
	route.GET("/health", publicHandler.Health)
	route.GET("/info", publicHandler.Info)
	pprof.Register(route, pprof.DefaultPrefix)

	swaggerRedirectHandler := func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	}
	route.GET("/swagger", swaggerRedirectHandler)
	route.GET("/swagg", swaggerRedirectHandler)

	return route
}
