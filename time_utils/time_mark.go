package time_utils

import "time"

type TimeRecorder struct {
	aTime time.Time
}

func MarkTime() *TimeRecorder {
	return &TimeRecorder{
		time.Now(),
	}
}

func (tr *TimeRecorder) Gap() int {
	return tr.GapMilli()
}

func (tr *TimeRecorder) GapMilli() int {
	now := time.Now()
	return int(now.Sub(tr.aTime).Seconds())
}

func (tr *TimeRecorder) GapSecond() int {
	now := time.Now()
	return int(now.Sub(tr.aTime).Seconds())
}

