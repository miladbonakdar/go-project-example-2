package health

import (
	"github.com/jinzhu/gorm"
)

func ConfigureHealthChecks(db *gorm.DB) {
	NewCheckerService().
		Add(NewDbHealthChecker("defaultConnection", db))
	//Add(NewServiceHealthChecker("hotelBaseService", c.ProviderEndpoint+"/api/v1/alive", c.HealthCheckThresholdInSecond))
}
