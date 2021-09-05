package main

import (
	"fmt"
	"hotel-engine/application/api"
	"hotel-engine/application/api/handlers"
	"hotel-engine/cmd/docs"
	"hotel-engine/core"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/logic"
	"hotel-engine/core/logic/balancenotifiers"
	"hotel-engine/core/messaging"
	"hotel-engine/infrastructure/config"
	"hotel-engine/infrastructure/feeder/rabbitmq"
	"hotel-engine/infrastructure/health"
	provider "hotel-engine/infrastructure/hotelproviderinterface"
	"hotel-engine/infrastructure/jobs"
	"hotel-engine/infrastructure/logger"
	"hotel-engine/infrastructure/mapper"
	"hotel-engine/infrastructure/repository"
	"hotel-engine/infrastructure/repository/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/gorm"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "net/http/pprof"
)

var db *gorm.DB
var feeder core.Feeder

func main() {
	c := config.Get()
	logger.ConfigureLogger(
		logger.LoggerConfiguration{
			ServiceName: c.ServiceName,
			Environment: c.Environment,
			ElasticUrl:  c.ElasticUrl,
		})
	db = sql.InitDatabase(c.ConnectionString)
	defer db.Close()
	health.ConfigureHealthChecks(db)
	unit := repository.NewUnitOfWork(db)
	hotelMapper := mapper.NewHotelMapper()
	hotelProvider := provider.NewHotelProvider()
	basicInfoProvider := provider.NewBasicInformationProvider()
	messagingClient := messaging.NewBusClient(c.Rabbitmq.ConnectionString)

	cacheStore := logic.NewCacheStore(unit, hotelMapper)
	providerSearchDtoFactory := logic.NewProviderSearchDtoFactory(cacheStore)
	balanceCheckerService := logic.NewProviderBalanceChecker(balancenotifiers.CreateBalanceAlertNotifiers())
	publicService := logic.NewPublicService(unit, hotelMapper, cacheStore, basicInfoProvider)
	orderEventDispatcher := logic.NewOrderEventDispatcher(messagingClient, c.RefundEventTopic)
	hotelService := logic.NewHotelService(unit, hotelMapper, hotelProvider, providerSearchDtoFactory,
		publicService, cacheStore, balanceCheckerService, orderEventDispatcher)

	logic.NewRateReviewEventHandler(messagingClient, hotelService, c.RateReviewSubscribeString)
	//feeder
	feeder, err := rabbitmq.NewFeeder(messagingClient, c.Rabbitmq.Feeder.Feed, c.Rabbitmq.Feeder.Seed, c.Rabbitmq.Feeder.Alias)
	if err != nil {
		logger.WithName(logtags.CreateFeederError).
			PanicException(err, "error wile creating a feeder")
	}
	defer feeder.Close()
	syncService := logic.NewSyncService(feeder, unit, hotelService)
	redisMemoryStorage := logic.NewRedisLocker()
	jobs.RegisterCronJobs(syncService, hotelService, redisMemoryStorage)

	hotelHandler := handlers.NewHotelHandler(hotelService, syncService)
	publicHandler := handlers.NewPublicHandler(publicService, balanceCheckerService)

	route := api.CreateRoute(hotelHandler, publicHandler)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%v", c.ContainerName, c.ContainerPort)
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	errServer := route.Run(fmt.Sprintf(":%d", c.ServerPort))

	if errServer != nil {
		logger.WithName(logtags.RunRouteError).
			PanicException(errServer, "error in running server")
	}
}

func close() {
	db.Close()
	feeder.Close()
}

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close()
		os.Exit(1)
	}()
}
