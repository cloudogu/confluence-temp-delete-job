package deletion

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
