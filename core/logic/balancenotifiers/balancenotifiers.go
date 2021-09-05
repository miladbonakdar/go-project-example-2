package balancenotifiers

import "hotel-engine/core"

func CreateBalanceAlertNotifiers() []core.BalanceAlertNotifier {
	return []core.BalanceAlertNotifier{
		newSmsBalanceNotifier(),
		newLoggerBalanceNotifier(),
	}
}
