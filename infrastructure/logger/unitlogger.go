package logger

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type EntryLog struct {
	devMessage *string
	url        *string
	statusCode *int
	exception  error
	name       string
	data       interface{}
}

func newUnitLogger() *EntryLog {
	return &EntryLog{
		devMessage: nil,
		exception:  nil,
		name:       defaultLogName,
		data:       nil,
	}
}

func (u *EntryLog) log(level log.Level, message string) {
	entry := log.WithField(messageKey, message)
	entry = log.WithField(nameKey, u.name)
	if u.devMessage != nil {
		entry = entry.WithField(devMessageKey, *u.devMessage)
	}
	if u.url != nil {
		entry = entry.WithField(url, *u.url)
	}
	if u.statusCode != nil {
		entry = entry.WithField(status, *u.statusCode)
	}
	if u.exception != nil {
		entry = entry.WithField(exceptionKey, u.exception.Error())
	}
	if u.data != nil {
		entry = entry.WithField(dataKey, u.data)
	}
	entry.Log(level, message)
}

func (u *EntryLog) Trace(message string)   { u.log(log.TraceLevel, message) }
func (u *EntryLog) Debug(message string)   { u.log(log.DebugLevel, message) }
func (u *EntryLog) Print(message string)   { u.Info(message) }
func (u *EntryLog) Info(message string)    { u.log(log.InfoLevel, message) }
func (u *EntryLog) Warn(message string)    { u.log(log.WarnLevel, message) }
func (u *EntryLog) Warning(message string) { u.Warn(message) }
func (u *EntryLog) Error(message string)   { u.log(log.ErrorLevel, message) }
func (u *EntryLog) Fatal(message string)   { u.log(log.FatalLevel, message) }
func (u *EntryLog) Panic(message string)   { u.log(log.PanicLevel, message) }

func (u *EntryLog) PanicException(exception error, message string) {
	u.exception = exception
	u.Panic(message)
}
func (u *EntryLog) WarnException(exception error, message string) {
	u.exception = exception
	u.Warn(message)
}
func (u *EntryLog) WarningException(exception error, message string) {
	u.WarnException(exception, message)
}
func (u *EntryLog) ErrorException(exception error, message string) {
	u.exception = exception
	u.Error(message)
}
func (u *EntryLog) FatalException(exception error, message string) {
	u.exception = exception
	u.Fatal(message)
}

func (u *EntryLog) WithException(exception error) *EntryLog {
	u.exception = exception
	return u
}

func (u *EntryLog) WithName(name string) *EntryLog {
	u.name = name
	return u
}

func (u *EntryLog) WithDevMessage(devMessage string) *EntryLog {
	u.devMessage = &devMessage
	return u
}

func (u *EntryLog) WithStatusCode(code int) *EntryLog {
	u.statusCode = &code
	return u
}

func (u *EntryLog) OkResponse() *EntryLog {
	code := 200
	u.statusCode = &code
	return u
}

func (u *EntryLog) WithUrl(url string) *EntryLog {
	u.url = &url
	return u
}

func (u *EntryLog) WithData(data interface{}) *EntryLog {
	if reflect.ValueOf(data).Kind() == reflect.Struct || reflect.ValueOf(data).Kind() == reflect.Map {
		b, err := json.Marshal(data)
		if err != nil {
			u.data = data
			return u
		}
		u.data = string(b)
	} else {
		u.data = data
	}
	return u
}
