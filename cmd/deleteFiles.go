package cmd

import (
	"github.com/cloudogu/confluence-temp-delete-job/deletion"
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
	cliArgs := c.Args()
	maxAgeInHours := c.Int(flagMaxAgeHoursLong)

	args := deletion.Args{CliArgs: cliArgs, MaxAgeInHours: maxAgeInHours}
	return deleteFilesWithArgs(args)
}

func deleteFilesWithArgs(args deletion.Args) error {
	return nil
}
