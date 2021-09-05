package jobs

import (
	"hotel-engine/core"
	"hotel-engine/core/common/logtags"
	"hotel-engine/infrastructure/config"
	"hotel-engine/infrastructure/logger"
	"time"
)

type syncHotelsCronJob struct {
	service      core.SyncService
	syncLockKey  string
	locker       core.DistributedLocker
	trySyncUntil int
}

func (s *syncHotelsCronJob) do() {
	s.endeavorToWin(s.syncHotels)
}

func (s *syncHotelsCronJob) syncHotels() {
	if !s.canSync() {
		return
	}
	s.service.UpdateDatabase()
	s.service.SyncElastic()
	logger.WithName(logtags.SyncingHotelsJobCompleted).WithDevMessage("sync hotels job -> do").
		Info("updating database and syncing elastic search finished")
}

func (s *syncHotelsCronJob) canSync() bool {
	nowHour := time.Now().Hour()
	if nowHour >= s.trySyncUntil {
		return false
	}
	synced, err := s.service.HasBeenSynced()
	if err != nil {
		logger.ErrorException(err, "got error while calling `HasBeenSynced` function")
		return false
	}
	return !synced
}

func (s *syncHotelsCronJob) endeavorToWin(runIfWon func()) {
	err := s.locker.Lock(s.syncLockKey, time.Second*10, runIfWon)
	if err != nil {
		logger.ErrorException(err, "error while trying to obtain a lock")
		return
	}
}

func (*syncHotelsCronJob) cronTab() string {
	return config.Get().SyncHotelsCronTab
}

func newSyncHotelsCronJob(service core.SyncService, locker core.DistributedLocker) job {
	con := config.Get()
	job := &syncHotelsCronJob{
		service:      service,
		locker:       locker,
		syncLockKey:  con.SyncLockKey,
		trySyncUntil: con.TrySyncUntil,
	}
	return job
}
