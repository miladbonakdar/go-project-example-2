package httphelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"hotel-engine/core/common"
	"hotel-engine/core/common/logtags"
	"hotel-engine/infrastructure/hotelproviderinterface/dtos"
	"hotel-engine/infrastructure/logger"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetResponseError(res *http.Response) error {
	if res.StatusCode > 230 || res.StatusCode < 200 {
		var errorDto dtos.ErrorResponseDto
		s, _ := ioutil.ReadAll(res.Body)
		err := json.Unmarshal(s, &errorDto)
		if err != nil {
			resString := string(s)
			logger.WithName(logtags.IrrelevantResponseFromProviderError).WithException(err).
				WithData(map[string]interface{}{
					"providerResponse":  resString,
					"requestStatusCode": res.StatusCode,
					"url":               res.Request.URL.Path,
				}).
				WithDevMessage("provider balance notifier > getResponseError").
				Error("irrelevant response data from provider")

			if strings.Contains(resString, "ی شما تشخیص داده‌اند. یکی از موارد زیر ممکن است علت این مشکل  باشد") {
				return common.ProviderRateLimitProblem
			}
			if strings.Contains(resString, "Gateway Time-out") {
				return common.ProviderGateWayTimeOutProblem
			}

			return errors.New(fmt.Sprintf("irrelevant response data from provider, status code : %d", res.StatusCode))
		}

		return errors.New(errorDto.Message)
	}
	return nil
}
