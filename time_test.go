package kocto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHourAnchor(t *testing.T) {
	now := time.Date(2021, 10, 20, 0, 0, 0, 0, time.UTC)

	result := HourAnchor(now)

	t.Log("time: ", now)
	t.Log("result: ", result)

	assert.True(t, result.After(now))
	assert.LessOrEqual(t, now.Sub(result), time.Hour)

	assert.Equal(t, now.Year(), result.Year())
	assert.Equal(t, now.Month(), result.Month())
	assert.Equal(t, now.Day(), result.Day())
	assert.Equal(t, now.Hour()+1, result.Hour())
	assert.Equal(t, result.Minute(), 0)
}

func TestMidnightAnchor(t *testing.T) {
	now := time.Date(2021, 10, 20, 0, 0, 0, 0, time.UTC)

	result := MidnightAnchor(now)

	t.Log("time: ", now)
	t.Log("result: ", result)

	assert.True(t, result.After(now))

	assert.Equal(t, now.Year(), result.Year())
	assert.Equal(t, now.Month(), result.Month())
	assert.Equal(t, now.Day()+1, result.Day())
    assert.Equal(t, result.Hour(), 0)
}

func TestAnchorsRespectsLocation(t *testing.T) {
	now := time.Date(2021, 10, 20, 0, 0, 0, 0, time.FixedZone("TEST", 20))

    r1 := HourAnchor(now)
    r2 := MidnightAnchor(now)

    assert.Equal(t, now.Location(), r1.Location())
    assert.Equal(t, now.Location(), r2.Location())
}
