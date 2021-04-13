package cmd

import (
	"fmt"
	"github.com/cloudogu/confluence-temp-delete-job/deletion"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const (
	flagMaxAgeHoursLong  = "age"
	flagMaxAgeHoursShort = "a"
)

// DeleteFilesCommand provides CLI entry logic for deleting files..
var DeleteFilesCommand = &cli.Command{
	Name:  "delete",
	Usage: "Recursively delete files and directories according the given parameters",
	Description: "This command recursively walks the given start directory and deletes files older than the given `age`. " +
		"Directories will only be deleted last and only if there are no files left to be contained.",
	Action:    deleteFiles,
	ArgsUsage: "directory",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     flagMaxAgeHoursLong,
			Usage:    "Sets the max. age of files and directories in hours that will be deleted. Must be larger than zero.",
			Required: true,
			Aliases:  []string{flagMaxAgeHoursShort},
		},
	},
}

func deleteFiles(c *cli.Context) error {
	directory := ""
	maxAgeInHours := c.Int(flagMaxAgeHoursLong)

	switch c.Args().Len() {
	case 1:
		directory = c.Args().First()
	case 0:
		_ = cli.ShowAppHelp(c)
		return fmt.Errorf("expected directory")
	default:
		_ = cli.ShowAppHelp(c)
		return fmt.Errorf("unexpected argument(s) found: %v", c.Args().Slice()[1:])
	}

	args := deletion.Args{Directory: directory, MaxAgeInHours: maxAgeInHours}
	return deleteFilesWithArgs(args)
}

func deleteFilesWithArgs(args deletion.Args) error {
	deleter, err := deletion.New(args)
	if err != nil {
		return errors.Wrap(err, "could not create deleter")
	}

	results, err := deleter.Execute()
	if err != nil {
		return errors.Wrap(err, "an error occurred during deletion")
	}
	results.PrintStats()

	return nil
}
