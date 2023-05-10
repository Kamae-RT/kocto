package kocto

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestHourAnchor(t *testing.T) {
	now := time.Date(2021, 10, 20, 0, 0, 0, 0, time.UTC)

	result := HourAnchor(now)

	t.Log("time: ", now)
	t.Log("result: ", result)

	is := is.New(t)
	is.True(result.After(now))
	is.True(now.Sub(result) <= time.Hour)

	is.Equal(now.Year(), result.Year())
	is.Equal(now.Month(), result.Month())
	is.Equal(now.Day(), result.Day())
	is.Equal(now.Hour()+1, result.Hour())
	is.Equal(result.Minute(), 0)
}

func TestMidnightAnchor(t *testing.T) {
	now := time.Date(2021, 10, 20, 0, 0, 0, 0, time.UTC)

	result := MidnightAnchor(now)

	t.Log("time: ", now)
	t.Log("result: ", result)

	is := is.New(t)

	is.True(result.After(now))

	is.Equal(now.Year(), result.Year())
	is.Equal(now.Month(), result.Month())
	is.Equal(now.Day()+1, result.Day())
	is.Equal(result.Hour(), 0)
}

func TestAnchorsRespectsLocation(t *testing.T) {
	now := time.Date(2021, 10, 20, 0, 0, 0, 0, time.FixedZone("TEST", 20))

	r1 := HourAnchor(now)
	r2 := MidnightAnchor(now)

	is := is.New(t)

	is.Equal(now.Location(), r1.Location())
	is.Equal(now.Location(), r2.Location())
}
