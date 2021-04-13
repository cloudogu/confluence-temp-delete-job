package deletion

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testMaxAgeInHours = 12

func TestNew(t *testing.T) {
	t.Run("should set args", func(t *testing.T) {
		input := Args{Directory: "/test", MaxAgeInHours: 42}

		sut, _ := New(input)

		assert.Equal(t, "/test", sut.Directory)
		assert.Equal(t, 42, sut.MaxAgeInHours)
	})

	type args struct {
		args Args
	}
	tests := []struct {
		name        string
		args        args
		wantDeleter bool
		wantErr     bool
	}{
		{"should pass", args{Args{Directory: "/test", MaxAgeInHours: 12}}, true, false},
		{"should pass with 0 age", args{Args{Directory: "/test", MaxAgeInHours: 0}}, true, false},
		{"should fail with dir", args{Args{Directory: "", MaxAgeInHours: 12}}, false, true},
		{"should fail with age", args{Args{Directory: "/a", MaxAgeInHours: -1}}, false, true},
		{"should fail with both dir and age", args{Args{Directory: "", MaxAgeInHours: -1}}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got != nil) != tt.wantDeleter {
				t.Errorf("New() deleter = %v, wantDeleter %v", got, tt.wantDeleter)
				return
			}
		})
	}
}

func Test_fileOlderThan(t *testing.T) {
	nowClock = &testClock{time.Now()}

	t.Run("should return true for ancient files", func(t *testing.T) {
		theBeginningOfComputing := time.Unix(0, 0)

		// when
		actual := fileOlderThan(testMaxAgeInHours, theBeginningOfComputing)

		// then
		assert.True(t, actual)
	})
	t.Run("should return true for 12 hours and 1 second old files", func(t *testing.T) {
		minus12HoursTime := nowClock.Now().Add(-12 * time.Hour).Add(-1 * time.Second)

		// when
		actual := fileOlderThan(testMaxAgeInHours, minus12HoursTime)

		// then
		assert.True(t, actual)
	})
	t.Run("should return false for exactly 12 hours old files", func(t *testing.T) {
		minus12HoursTime := nowClock.Now().Add(-12 * time.Hour)

		// when
		actual := fileOlderThan(testMaxAgeInHours, minus12HoursTime)

		// then
		assert.False(t, actual)
	})
	t.Run("should return false for 11 hours and 59 minutes old files", func(t *testing.T) {
		minus11Hours59SecTime := nowClock.Now().Add(-11 * time.Hour).Add(-59 * time.Second)

		// when
		actual := fileOlderThan(testMaxAgeInHours, minus11Hours59SecTime)

		// then
		assert.False(t, actual)
	})
	t.Run("should return false for 11 hours old files", func(t *testing.T) {
		minus11HoursTime := nowClock.Now().Add(-11 * time.Hour)

		// when
		actual := fileOlderThan(testMaxAgeInHours, minus11HoursTime)

		// then
		assert.False(t, actual)
	})
	t.Run("should return false for 0 hours old files", func(t *testing.T) {
		// when
		actual := fileOlderThan(testMaxAgeInHours, nowClock.Now())

		// then
		assert.False(t, actual)
	})
	t.Run("should return false for files back from the future", func(t *testing.T) {
		theFuture := nowClock.Now().Add(24 * 365 * time.Hour)

		// when
		actual := fileOlderThan(testMaxAgeInHours, theFuture)

		// then
		assert.False(t, actual)
	})
}

type testClock struct {
	desiredTime time.Time
}

func (t *testClock) Now() time.Time {
	return t.desiredTime
}
