package health

import (
	"hotel-engine/core/common/logtags"
	"hotel-engine/infrastructure/logger"
	"net/http"
	"time"
)

type serviceHealthChecker struct {
	client      *http.Client
	endpoint    string
	tag         string
	lastAttempt time.Time
	lastResult  *HealthResultDto
	Config      struct {
		Threshold time.Duration
	}
}

var healthyResult = HealthResultDto{
	Status:   Healthy,
	Duration: defaultTimeStampFormat,
	Data:     map[string]string{},
}

func (c *serviceHealthChecker) Check() HealthResultDto {

	if time.Now().Sub(c.lastAttempt) <= c.Config.Threshold {
		res := c.check()
		c.lastResult = &res
		c.lastAttempt = time.Time{}
	}

	return *c.lastResult
}

func (c *serviceHealthChecker) check() HealthResultDto {
	re, err := c.client.Get(c.endpoint)
	if err != nil {
		logger.WithName(logtags.ServiceHealthCheckFailed).
			WithException(err).WithDevMessage("error while checking service endpoint : " + c.endpoint).
			Error("service is not responding properly")
		return HealthResultDto{
			Status:      UnHealthy,
			Duration:    defaultTimeStampFormat,
			Exception:   err.Error(),
			Description: err.Error(),
			Data:        map[string]string{},
		}
	}
	if re.StatusCode < 200 || re.StatusCode > 260 {
		logger.WithName(logtags.ServiceHealthCheckFailed).WithData(map[string]interface{}{
			"statusCode": re.StatusCode,
			"endpoint":   c.endpoint,
		}).Error("calling service endpoint resulted with invalid status code")
		return HealthResultDto{
			Status:      UnHealthy,
			Duration:    defaultTimeStampFormat,
			Exception:   "calling service endpoint resulted with invalid status code",
			Description: "calling service endpoint resulted with invalid status code",
			Data:        map[string]string{},
		}
	}
	return healthyResult
}

func (c *serviceHealthChecker) Tag() string {
	return c.tag
}

func NewServiceHealthChecker(tag, endpoint string, thresholdInSecond int64) Checker {
	return &serviceHealthChecker{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		endpoint:    endpoint,
		tag:         tag,
		lastAttempt: time.Time{},
		Config: struct{ Threshold time.Duration }{
			Threshold: time.Duration(thresholdInSecond) * time.Second,
		},
	}
}
