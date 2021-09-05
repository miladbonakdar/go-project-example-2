package balancenotifiers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/config"
	"hotel-engine/infrastructure/logger"
	"hotel-engine/utils/httphelper"
	"net/http"
	"strings"
)

const notifySMSEndPoint = "/api/v2/sms/messages"
const smsProvidersEndpoint = "/api/v1/management/sms/providers"

type smsBalanceNotifier struct {
	client       *http.Client
	smsBaseUrl   string
	phoneNumbers []string
	smsProvider  *smsProvider
}

func (n *smsBalanceNotifier) Notify(alert dto.BalanceNotifyDto) {
	if len(n.phoneNumbers) == 0 && n.smsProvider == nil {
		return
	}
	var numbers []map[string]interface{}

	for _, number := range n.phoneNumbers {
		message := strings.Replace(common.SmsNotificationMessage, "{account}", alert.Account, 1)
		message = strings.Replace(message, "{currentBalance}", fmt.Sprintf("%d", alert.Balance), 1)
		message = strings.Replace(message, "{balanceLimit}", fmt.Sprintf("%d", alert.BalanceAlertLimit), 1)
		numbers = append(numbers, map[string]interface{}{
			"to":         number,
			"text":       message,
			"providerId": n.smsProvider.Id,
			"smsType":    "Normal",
		})
	}

	body, err := json.Marshal(map[string]interface{}{
		"messages": numbers,
	})
	if err != nil {
		logger.WithName(logtags.CannotCreateHttpBody).ErrorException(err, "cannot create http body for sms provider send sms endpoint")
		return
	}
	req, err := http.NewRequest("POST", n.smsBaseUrl+notifySMSEndPoint,
		bytes.NewBuffer(body))
	if err != nil {
		logger.WithName(logtags.CannotCreateRequestObject).ErrorException(err, "cannot create request object for sms provider send sms endpoint")
		return
	}
	req.Header.Add("Content-Type", "application/json-patch+json")
	req.Header.Add("Accept", "application/json")

	req.Close = true

	res, err := n.client.Do(req)
	if err != nil {
		logger.WithName(logtags.SmsProviderError).
			ErrorException(err, "error while calling provider send sms endpoint")
		return
	}

	err = httphelper.GetResponseError(res)
	if err != nil {
		logger.WithName(logtags.SmsProviderError).ErrorException(err, "error manipulating provider send sms endpoint response")
		return
	}

}

func newSmsBalanceNotifier() core.BalanceAlertNotifier {
	con := config.Get()
	client := &http.Client{}
	smsProvider, err := findSmsProvider(con.SmsProviderEndpoint, client)

	if err != nil {
		logger.WithName(logtags.SmsProviderError).ErrorException(err, "cannot find a provider for sms notifications")
	} else {
		logger.Info("sms notification provider founded successfully")
	}
	return &smsBalanceNotifier{
		client:       client,
		smsBaseUrl:   con.SmsProviderEndpoint,
		phoneNumbers: con.BalanceAlertPhoneNumbers,
		smsProvider:  smsProvider,
	}
}
