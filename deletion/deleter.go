package deletion

// Args transports CLI parameters to the business package.
type Args struct {
	// Directory names the starting directory which the deleter will recursively inspect for old files.
	Directory string
	// MaxAgeInHours sets how old at least a file or directory must be before it will be selected for deletion.
	MaxAgeInHours int
}

type deleter struct {
}

func New(args Args) (*deleter, error) {
	var directory string

	if directory == "1" {

	}
	return nil, nil
}

func (d *deleter) Execute() (*Results, error) {
	return &Results{}, nil
}

type Results struct {
}

func (r *Results) Print() {

}
