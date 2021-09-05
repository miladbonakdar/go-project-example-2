package date

import (
	"time"
)

const LayoutISO = "2006-01-02"

func StringToDateUTC(value string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", value, time.UTC)
	return t, err
}

func StringToDate(value string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", value)
	return t, err
}

func StringToDateUTCOrDefault(value string) time.Time {
	t, err := StringToDateUTC(value)
	if err != nil {
		return time.Now().UTC()
	}
	return t
}

func StringToDateOrDefault(value string) time.Time {
	t, err := StringToDate(value)
	if err != nil {
		return time.Now().UTC()
	}
	return t
}

func CompareTwoDates(first string, second string) (bool, error) {
	t1, err := StringToDate(first)
	if err != nil {
		return false, err
	}
	t2, err := StringToDate(second)
	if err != nil {
		return false, err
	}
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() &&
		t1.Day() == t2.Day(), nil
}
