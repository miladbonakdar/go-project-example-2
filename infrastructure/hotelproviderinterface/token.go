package hotelProviderInterface

import (
	"hotel-engine/core/common/logtags"
	"hotel-engine/infrastructure/logger"
	"sync"
)

var token string
var mutex = &sync.Mutex{}

func Token() (string, error) {
	if token != "" {
		return token, nil
	}
	mutex.Lock()
	defer mutex.Unlock()
	if token != "" {
		return token, nil
	}
	_, err := UpdateToken()
	if err != nil {
		logger.WithName(logtags.CannotGetAccessTokenError).
			ErrorException(err, "error in getting access token")
		return "", err
	}
	return token, nil
}

func UpdateToken() (string, error) {
	provider := NewHotelProvider()
	t, err := provider.GetAccessToken()
	if err != nil {
		return "", err
	}
	token = t
	logger.Info("token has been updated")
	return t, nil
}
