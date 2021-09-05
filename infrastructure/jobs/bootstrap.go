package jobs

import (
	"fmt"
	"hotel-engine/core"

	"github.com/robfig/cron/v3"
)

type job interface {
	do()
	cronTab() string
}

func RegisterCronJobs(service core.SyncService, hotelService core.HotelService, locker core.DistributedLocker) {
	c := cron.New()

	syncHotelsCronJob := newSyncHotelsCronJob(service, locker)
	refundPullingCronJob := newRefundPullingCronJob(hotelService)
	syncTokenCronJob := newSyncTokenCronJob()

	c.AddFunc(syncHotelsCronJob.cronTab(), syncHotelsCronJob.do)
	c.AddFunc(syncTokenCronJob.cronTab(), syncTokenCronJob.do)
	c.AddFunc(refundPullingCronJob.cronTab(), refundPullingCronJob.do)

	c.Start()
	fmt.Println("all cron jobs registered")
}
