package deletion

import (
	"errors"
	errors2 "github.com/pkg/errors"
	"os"
	"path/filepath"
	"time"
)

var nowClock clock = &realClock{}

type clock interface {
	Now() time.Time
}

type realClock struct{}

func (r realClock) Now() time.Time {
	return time.Now()
}

// Args transports CLI parameters to the business package.
type Args struct {
	// Directory names the starting directory which the deleter will recursively inspect for old files.
	Directory string
	// MaxAgeInHours sets how old at least a file or directory must be before it will be selected for deletion.
	MaxAgeInHours int
}

type deleter struct {
	Args
	Results *Results
}

func New(args Args) (*deleter, error) {
	if args.Directory == "" {
		return nil, errors.New("directory must not be empty")
	}
	if args.MaxAgeInHours < 0 {
		return nil, errors.New("file age must zero or positive")
	}

	return &deleter{args, &Results{}}, nil
}

func (d *deleter) Execute() (*Results, error) {
	err := filepath.Walk(d.Directory, d.filterOldFiles)
	if err != nil {
		return d.Results, err
	}

	return d.Results, nil
}

func (d *deleter) filterOldFiles(path string, info os.FileInfo, err error) error {
	if err != nil {
		return errors2.Wrapf(err, "error while visiting path %q", path)
	}

	if info.IsDir() {
		return nil
	}

	if fileOlderThan(d.MaxAgeInHours, info.ModTime()) {
		return d.deleteFile(path)
	}

	return nil
}

func (d *deleter) deleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		d.Results.Fail(path, err)
	} else {
		d.Results.Pass(path)
	}

	return err
}

func fileOlderThan(maxAgeHours int, fileTime time.Time) bool {
	ageCutOff := time.Duration(maxAgeHours) * time.Hour
	now := nowClock.Now()

	if diff := now.Sub(fileTime); diff <= ageCutOff {
		return false
	}
	return true
}

type Results struct {
}

// PrintStats prints deletion statistics as one-liner.
func (r *Results) PrintStats() {

}

func (r *Results) Fail(path string, err error) {

}

func (r *Results) Pass(path string) {

}
