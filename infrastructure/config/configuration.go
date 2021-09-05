package config

import (
	"hotel-engine/utils/date"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Configuration struct {
	HealthCheckThresholdInSecond int64
	ProviderEndpoint             string
	OrderEndpoint                string
	SyncSecret                   string
	ServerPort                   int
	ContainerPort                int
	SyncChunkSize                int
	ContainerName                string
	SyncHotelsCronTab            string
	SyncTokenCronTab             string
	ProviderUsername             string
	ProviderPassword             string
	ConnectionString             string
	ElasticUrl                   string
	ServiceName                  string
	BalanceAlertLimit            float64
	BalanceAlertPhoneNumbers     []string
	Rabbitmq                     struct {
		ConnectionString string
		Feeder           struct {
			Feed  string
			Seed  string
			Alias string
		}
	}
	Environment               string
	SmsProviderEndpoint       string
	SyncLockKey               string
	MemoryStorageConnection   string
	RefundPullingMaxTryDays   int
	RefundPullingFromDate     time.Time
	RefundPullingCronTab      string
	RefundEventTopic          string
	RateReviewSubscribeString string
	OrderServiceEndpoint      string
	AvailableHotelsWhiteList  []string
	AvailablePhonesWhiteList  []string
	TrySyncUntil              int
}

func (l Configuration) IsProduction() bool {
	return l.Environment == "Production" || l.Environment == "production"
}

func (l Configuration) IsStaging() bool {
	return l.Environment == "Staging" || l.Environment == "staging"
}

func (l Configuration) IsDevelopment() bool {
	return l.Environment == "Development" || l.Environment == "development"
}

func CreateDefaultConfiguration() Configuration {
	port, err := strconv.Atoi(os.Getenv("HOTEL_ENGINE_SERVER_PORT"))
	if err != nil {
		log.Fatalln("The port number is not valid")
	}

	chunkSize, err := strconv.Atoi(os.Getenv("HOTEL_ENGINE_SYNC_CHUNK_SIZE"))
	if err != nil {
		log.Fatalln("The syn chunk size number is not valid")
	}

	balanceLimit, err := strconv.ParseFloat(os.Getenv("HOTEL_ENGINE_BALANCE_ALERT_LIMIT"), 8)
	if err != nil {
		log.Fatalln("The balance limit value is not valid")
	}

	outSideOfContainerPort, err := strconv.Atoi(os.Getenv("HOTEL_ENGINE_CONTAINER_PORT"))
	if err != nil {
		log.Fatalln("The container port number is not valid")
	}

	refundPullingFromDate, err := date.StringToDate(os.Getenv("HOTEL_ENGINE_REFUND_PULLING_FROM_DATE"))
	if err != nil {
		log.Fatalln("The refund pulling from date is not valid date format")
	}

	refundPullingMaxTryDays, err := strconv.Atoi(os.Getenv("HOTEL_ENGINE_REFUND_PULLING_MAX_TRY_DAYS"))
	if err != nil {
		log.Fatalln("The refund pulling max try days number is not valid")
	}

	trySyncUntil, err := strconv.Atoi(os.Getenv("HOTEL_ENGINE_TRY_SYNCING_UNTIL"))
	if err != nil {
		log.Fatalln("The TrySyncUntil number is not valid")
	}

	healthCheckThresholdInSecond, err := strconv.ParseInt(os.Getenv("HOTEL_ENGINE_HealthCheck_ThresholdInSecond"), 10, 64)
	if err != nil {
		log.Fatalln("The healthCheckAttempts number is not valid")
	}

	return Configuration{
		ContainerName:             os.Getenv("HOTEL_ENGINE_CONTAINER_NAME"),
		ContainerPort:             outSideOfContainerPort,
		ProviderEndpoint:          os.Getenv("HOTEL_ENGINE_HOTEL_PROVIDER_ENDPOINT"),
		SmsProviderEndpoint:       os.Getenv("HOTEL_ENGINE_SMS_PROVIDER_ENDPOINT"),
		OrderEndpoint:             os.Getenv("HOTEL_ENGINE_ORDER_ENDPOINT"),
		OrderServiceEndpoint:      os.Getenv("HOTEL_ENGINE_ORDER_SERVICE_ENDPOINT"),
		ServerPort:                port,
		SyncHotelsCronTab:         os.Getenv("HOTEL_ENGINE_UPDATE_HOTELS_CRON_TAB"),
		SyncTokenCronTab:          os.Getenv("HOTEL_ENGINE_UPDATE_TOKEN_CRON_TAB"),
		ProviderUsername:          os.Getenv("HOTEL_ENGINE_SYNC_USERNAME"),
		ProviderPassword:          os.Getenv("HOTEL_ENGINE_SYNC_PASSWORD"),
		ConnectionString:          os.Getenv("HOTEL_ENGINE_DB_CONNECTION_STRING"),
		SyncSecret:                os.Getenv("HOTEL_ENGINE_SYNC_SECRET"),
		Environment:               os.Getenv("HOTEL_ENGINE_ENVIRONMENT"),
		ElasticUrl:                os.Getenv("HOTEL_ENGINE_ELASTIC_URL"),
		SyncLockKey:               os.Getenv("HOTEL_ENGINE_SYNC_LOCK_KEY"),
		MemoryStorageConnection:   os.Getenv("HOTEL_ENGINE_MEMORY_STORAGE_CONNECTION"),
		RefundPullingCronTab:      os.Getenv("HOTEL_ENGINE_REFUND_PULLING_CRON_TAB"),
		RefundPullingFromDate:     refundPullingFromDate,
		RefundPullingMaxTryDays:   refundPullingMaxTryDays,
		RefundEventTopic:          os.Getenv("HOTEL_ENGINE_REFUND_EVENT_TOPIC"),
		ServiceName:               os.Getenv("APP_NAME"),
		RateReviewSubscribeString: os.Getenv("HOTEL_ENGINE_RATE_REVIEW_SUBSCRIBE_STRING"),
		BalanceAlertLimit:         balanceLimit,
		SyncChunkSize:             chunkSize,
		BalanceAlertPhoneNumbers:  strings.Split(os.Getenv("HOTEL_ENGINE_BALANCE_ALERT_PHONES"), ","),
		AvailableHotelsWhiteList:  strings.Split(os.Getenv("HOTEL_ENGINE_STAGE_AVAILABLE_HOTELS_WHITE_LIST"), ","),
		AvailablePhonesWhiteList:  strings.Split(os.Getenv("HOTEL_ENGINE_STAGE_AVAILABLE_PHONES_WHITE_LIST"), ","),
		TrySyncUntil:              trySyncUntil,
		Rabbitmq: struct {
			ConnectionString string
			Feeder           struct {
				Feed  string
				Seed  string
				Alias string
			}
		}{
			ConnectionString: os.Getenv("Rabbitmq_ConnectionString"),
			Feeder: struct {
				Feed  string
				Seed  string
				Alias string
			}{
				Feed:  os.Getenv("Rabbitmq_Feeder_Feed"),
				Seed:  os.Getenv("Rabbitmq_Feeder_Seed"),
				Alias: os.Getenv("Rabbitmq_Feeder_Alias"),
			},
		},
		HealthCheckThresholdInSecond: healthCheckThresholdInSecond,
	}
}
