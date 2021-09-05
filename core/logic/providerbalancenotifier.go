package logic

import (
	"hotel-engine/core"
	"hotel-engine/core/common"
	"hotel-engine/core/common/logtags"
	"hotel-engine/core/dto"
	"hotel-engine/infrastructure/config"
	hotelProviderInterface "hotel-engine/infrastructure/hotelproviderinterface"
	"hotel-engine/infrastructure/logger"
	"hotel-engine/utils/atomicflag"
	"hotel-engine/utils/httphelper"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

const balanceEndPoint = "/api/v1/profile/account/balance"

var notified = atomicflag.NewAtomicFlag()

type alibabaBalanceChecker struct {
	client          *http.Client
	baseProviderUri string
	notifiers       []core.BalanceAlertNotifier
	balanceLimit    float64
	providerAccount string
}

func (p *alibabaBalanceChecker) getFromUrl(url string, needAuthentication bool) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Close = true

	return p.requestToUrl(req, needAuthentication)
}

func (p *alibabaBalanceChecker) requestToUrl(req *http.Request, needAuthentication bool) (string, error) {
	if needAuthentication {
		token, err := hotelProviderInterface.Token()
		if err != nil {
			return "", err
		}
		req.Header.Add("Authorization", "Bearer "+token)
	}
	res, err := p.client.Do(req)
	if err != nil {
		return "", err
	}

	err = httphelper.GetResponseError(res)
	if err != nil {
		return "", err
	}

	jsonData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	jsonString := string(jsonData)

	if !gjson.Valid(jsonString) {
		logger.WithName(logtags.SmsProviderError).WithData(map[string]interface{}{
			"response": jsonString,
		}).Error(common.JsonDataIsNotValid.Error())
	}
	return jsonString, nil
}

func (a *alibabaBalanceChecker) GetBalance() (float64, error) {
	stringData, err := a.getFromUrl(a.baseProviderUri+balanceEndPoint,
		true)
	if err != nil {
		return 0.0, err
	}
	return gjson.Get(stringData, "result.balance").Float(), nil
}

func (a *alibabaBalanceChecker) CheckAdequateBalance() {
	balance, err := a.GetBalance()
	if err != nil {
		logger.WithName(logtags.SmsProviderError).ErrorException(err, "cannot get balance from provider")
		return
	}
	if balance > a.balanceLimit {
		if notified.Get() {
			notified.Set(false)
		}
		return
	}

	if notified.Get() {
		return
	}

	notificationObject := dto.NewBalanceNotifyDto(int64(balance), int64(a.balanceLimit), a.providerAccount)
	for _, notifier := range a.notifiers {
		go notifier.Notify(notificationObject)
	}
	notified.Set(true)
}

func NewProviderBalanceChecker(notifiers []core.BalanceAlertNotifier) core.ProviderBalanceChecker {
	con := config.Get()
	return &alibabaBalanceChecker{
		client:          &http.Client{},
		baseProviderUri: con.OrderEndpoint,
		notifiers:       notifiers,
		balanceLimit:    con.BalanceAlertLimit,
		providerAccount: con.ProviderUsername,
	}
}
