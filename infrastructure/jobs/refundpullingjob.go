package jobs

import (
	"hotel-engine/core"
	"hotel-engine/infrastructure/config"
	"time"
)

type refundPullingCronJob struct {
	cron                    string
	refundPullingFromDate   time.Time
	refundPullingMaxTryDays int
	service                 core.HotelService
}

func (r *refundPullingCronJob) do() {
	fromDate := time.Now().AddDate(0, 0, -r.refundPullingMaxTryDays)
	if fromDate.Before(r.refundPullingFromDate) {
		fromDate = r.refundPullingFromDate
	}
	r.service.UpdateRefundedOrdersPaymentStatus(fromDate)
}

func (r *refundPullingCronJob) cronTab() string {
	return r.cron
}

func newRefundPullingCronJob(hotelService core.HotelService) job {
	con := config.Get()
	return &refundPullingCronJob{
		cron:                    con.RefundPullingCronTab,
		refundPullingFromDate:   con.RefundPullingFromDate,
		refundPullingMaxTryDays: con.RefundPullingMaxTryDays,
		service:                 hotelService,
	}
}
