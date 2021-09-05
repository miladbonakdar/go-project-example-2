package jobs

import (
	"hotel-engine/core/common/logtags"
	"hotel-engine/infrastructure/config"
	"hotel-engine/infrastructure/hotelproviderinterface"
	"hotel-engine/infrastructure/logger"
)

type syncTokenCronJob struct {
}

func (*syncTokenCronJob) do() {
	_, err := hotelProviderInterface.UpdateToken()
	if err != nil {
		logger.WithName(logtags.CannotGetAccessTokenError).ErrorException(err, "error while syncing access token")
	}
}

func (*syncTokenCronJob) cronTab() string {
	return config.Get().SyncTokenCronTab
}

func newSyncTokenCronJob() job {
	return &syncTokenCronJob{}
}
