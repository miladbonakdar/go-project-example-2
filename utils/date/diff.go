package date

import "time"

func DaysDiff(start, end string) (int, error) {
	t1, err := StringToDate(start)
	if err != nil {
		return 0, err
	}
	t2, err := StringToDate(end)
	if err != nil {
		return 0, err
	}
	t2 = t2.Add(time.Second)
	return int(t2.Sub(t1).Hours() / 24), nil
}
