package deletion

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	"github.com/op/go-logging"
	errors2 "github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	log                  = logging.MustGetLogger("deletion")
	nowClock clock       = &realClock{}
	remover  fileRemover = &realFileRemover{}
)

// Args transports CLI parameters to the business package.
type Args struct {
	// Directory names the starting directory which the deleter will recursively inspect for old files.
	Directory string
	// MaxAgeInHours sets how old at least a file or directory must be before it will be selected for deletion.
	MaxAgeInHours int
}

type clock interface {
	Now() time.Time
}

type realClock struct{}

func (r realClock) Now() time.Time {
	return time.Now()
}

type fileRemover interface {
	Remove(path string) error
}

type realFileRemover struct{}

func (r *realFileRemover) Remove(path string) error {
	return os.Remove(path)
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
	var err error

	log.Debug("Start recursive file deletion")
	fileErr := filepath.Walk(d.Directory, d.filterOldFiles)
	if fileErr != nil {
		err = multierror.Append(err, fileErr)
	}

	log.Debug("Start recursive directory deletion")
	// delete old and empty directories because recursive directories are complicated during the first file walk
	dirErr := filepath.Walk(d.Directory, d.filterOldDirectories)
	if dirErr != nil {
		err = multierror.Append(err, dirErr)
	}

	return d.Results, err
}

func (d *deleter) filterOldFiles(path string, info os.FileInfo, err error) error {
	if err != nil {
		return errors2.Wrapf(err, "error while visiting path %q", path)
	}

	if info.IsDir() {
		return nil
	}

	if fileOlderThan(d.MaxAgeInHours, info.ModTime()) {
		return d.deleteFile(path, info)
	}

	d.Results.skip(path)

	return nil
}

func (d *deleter) filterOldDirectories(path string, info os.FileInfo, err error) error {
	if err != nil {
		return errors2.Wrapf(err, "error while visiting path %q", path)
	}

	// do not count-in the start directory, not even into the skipping counter
	if path == d.Directory {
		return nil
	}

	if !info.IsDir() {
		log.Debugf("walk directories: skip file %s", path)
		return nil
	}

	empty, err := isDirectoryEmpty(path)
	if err != nil {
		return errors2.Wrapf(err, "error while checking directory contents for path %q", path)
	}

	// deleting empty directories is simpler than checking the timestamp because the timestamp changes on every
	// file deletion inside the directory
	if empty {
		return d.deleteFile(path, info)
	}
	d.Results.skip(path)

	return nil
}

// taken from https://stackoverflow.com/a/30708914/12529534
func isDirectoryEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer func() { _ = f.Close() }()

	const limitFileCount = 2 // we only want to know if the dir is empty and not how many objects there are
	_, err = f.Readdir(limitFileCount)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}

func (d *deleter) deleteFile(path string, info os.FileInfo) error {
	err := remover.Remove(path)
	if err != nil {
		d.Results.fail(path, err)
	} else {
		d.Results.pass(path, info)
	}

	return err
}

func fileOlderThan(maxAgeHours int, fileTime time.Time) bool {
	ageCutOff := time.Duration(maxAgeHours) * time.Hour
	now := nowClock.Now()

	diff := now.Sub(fileTime)
	return diff > ageCutOff
}
