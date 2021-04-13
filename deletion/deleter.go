package deletion

import "errors"

// Args transports CLI parameters to the business package.
type Args struct {
	// Directory names the starting directory which the deleter will recursively inspect for old files.
	Directory string
	// MaxAgeInHours sets how old at least a file or directory must be before it will be selected for deletion.
	MaxAgeInHours int
}

type deleter struct {
	Args
}

func New(args Args) (*deleter, error) {
	if args.Directory == "" {
		return nil, errors.New("directory must not be empty")
	}
	if args.MaxAgeInHours < 0 {
		return nil, errors.New("file age must zero or positive")
	}
	return &deleter{args}, nil
}

func (d *deleter) Execute() (*Results, error) {
	return &Results{}, nil
}

type Results struct {
}

func (r *Results) Print() {

}
