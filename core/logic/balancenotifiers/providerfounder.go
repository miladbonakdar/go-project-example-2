package balancenotifiers

import (
	"errors"
	"github.com/tidwall/gjson"
	"hotel-engine/core/common"
	"hotel-engine/infrastructure/logger"
	"hotel-engine/utils/httphelper"
	"io/ioutil"
	"net/http"
)

type smsProvider struct {
	Id                int
	providerName      string
	lineNumber        string
	defaultProviderId int
}

func findSmsProvider(baseSmsUrl string, client *http.Client) (*smsProvider, error) {
	req, err := http.NewRequest("GET", baseSmsUrl+smsProvidersEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Close = true

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	err = httphelper.GetResponseError(res)
	if err != nil {
		return nil, err
	}

	jsonData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	jsonString := string(jsonData)

	if !gjson.Valid(jsonString) {
		logger.WithData(map[string]interface{}{
			"response": jsonString,
		}).Error(common.JsonDataIsNotValid.Error())
	}
	providers := gjson.Get(jsonString, "result").Array()

	if len(providers) == 0 {
		return nil, errors.New("there is no available provider right now")
	}

	return &smsProvider{
		Id:                int(providers[0].Get("id").Int()),
		providerName:      providers[0].Get("providerName").String(),
		lineNumber:        providers[0].Get("lineNumber").String(),
		defaultProviderId: int(providers[0].Get("defaultProviderId").Int()),
	}, nil
}
