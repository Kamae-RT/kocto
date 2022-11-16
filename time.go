package kocto

import (
	"math"
	"time"
)

const TimeLayout = "2006-01-02T15:04:05-0700"

func TimeParse(s string) (time.Time, error) {
	return time.Parse(TimeLayout, s)
}

func TicksFromTime(t time.Time) uint64 {
	ticksPerMilli := 10_000
	base := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	baseTicks := ticksPerMilli * int(math.Abs(float64(base)))
	ticks := t.UnixMilli() * int64(ticksPerMilli)

	return uint64(baseTicks) + uint64(ticks)
}
