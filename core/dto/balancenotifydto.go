package dto

type BalanceNotifyDto struct {
	Balance           int64
	BalanceAlertLimit int64
	Account           string
}

func NewBalanceNotifyDto(balance, balanceLimit int64, account string) BalanceNotifyDto {
	return BalanceNotifyDto{
		Balance:           balance,
		BalanceAlertLimit: balanceLimit,
		Account:           account,
	}
}
