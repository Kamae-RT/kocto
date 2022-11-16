package kocto

import (
	"math"
	"time"
)

const TimeLayout = "2006-01-02T15:04:05-0700"

func TimeParse(s string) (time.Time, error) {
	return time.Parse(TimeLayout, s)
}

func TicksToTime(ticks int64, offset int64) time.Time {
	// ticks / ticks_per_second - number_of_seconds_since_0001-01-01 00:00:00
	us := ticks/10000000 - 62135596800
	return time.Unix(us, offset)
}

func TicksFromTime(t time.Time) uint64 {
	ticksPerMilli := 10_000
	base := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	baseTicks := ticksPerMilli * int(math.Abs(float64(base)))
	ticks := t.UnixMilli() * int64(ticksPerMilli)

	return uint64(baseTicks) + uint64(ticks)
}

func UnixMilliToTicks(millis int64) uint64 {
	return TicksFromTime(time.UnixMilli(millis))
}

func HourAnchor(t time.Time) time.Time {
	n := time.Now()
	y, m, d := n.Date()

	return time.Date(y, m, d, n.Hour()+1, 0, 0, 0, time.UTC)
}

func ToNextHour(now time.Time) time.Duration {
    return now.Sub(HourAnchor(now))
}

func MidnightAnchor(t time.Time) time.Time {
	n := time.Now()
	y, m, d := n.Date()

	return time.Date(y, m, d+1, 0, 0, 0, 0, time.UTC)
}

func ToNextDay(now time.Time) time.Duration {
    return now.Sub(MidnightAnchor(now))
}
