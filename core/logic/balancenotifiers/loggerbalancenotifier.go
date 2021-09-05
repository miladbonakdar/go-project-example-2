package balancenotifiers

import (
	"hotel-engine/core"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/logger"
)

type loggerBalanceNotifier struct {
}

func (n *loggerBalanceNotifier) Notify(alert dto.BalanceNotifyDto) {
	logger.WithData(alert).WithName(logtags.BalanceLimitReached).
		Error("!!! Provider balance is lower that balance alert limit." +
			" please check the balance and charge the account ASAP")
}

func newLoggerBalanceNotifier() core.BalanceAlertNotifier {
	return &loggerBalanceNotifier{}
}
