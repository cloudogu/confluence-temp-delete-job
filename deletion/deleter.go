package deletion

import "github.com/urfave/cli/v2"

// Args transports CLI parameters to the business package.
type Args struct {
	// CliArgs contains all CLI arguments except flags and switches.
	CliArgs cli.Args
	// MaxAgeInHours sets how old at least a file or directory must be before it will be selected for deletion.
	MaxAgeInHours int
}
