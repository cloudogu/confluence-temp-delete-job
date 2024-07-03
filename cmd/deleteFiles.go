package cmd

import (
	"fmt"
	"github.com/cloudogu/confluence-temp-delete-job/deletion"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	flagMaxAgeHoursLong          = "age"
	flagMaxAgeHoursShort         = "a"
	flagLoopIntervalMinutesLong  = "interval"
	flagLoopIntervalMinutesShort = "i"
)

const cpuLoadSleepInSec = 2

var log = logging.MustGetLogger("cmd")

// DeleteFilesCommand provides CLI entry logic for deleting files..
var DeleteFilesCommand = &cli.Command{
	Name:  "delete-loop",
	Usage: "Endless loop that recursively deletes files and directories according the given parameters",
	Description: "This command recursively walks the given start directory and deletes files older than the given `age`. " +
		"Directories will only be deleted last and only if there are no files left to be contained. The loop will run " +
		"eternally until it receives the following signals: SIGHUP, SIGINT (Strg+C), SIGTERM, SIGKILL.",
	Action:    deleteFiles,
	ArgsUsage: "directory",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    flagMaxAgeHoursLong,
			Usage:   "Sets the max. age of files and directories in hours that will be deleted. Must be larger than zero.",
			Value:   12,
			Aliases: []string{flagMaxAgeHoursShort},
		},
		&cli.IntFlag{
			Name:    flagLoopIntervalMinutesLong,
			Usage:   "Sets the interval in minutes to run the deletion routine. Must be larger than zero.",
			Value:   60,
			Aliases: []string{flagLoopIntervalMinutesShort},
		},
	},
}

func deleteFiles(c *cli.Context) error {
	directory := ""
	maxAgeInHours := c.Int(flagMaxAgeHoursLong)
	loopIntervalInMin := c.Int(flagLoopIntervalMinutesLong)
	loopInterval := minuteToDuration(loopIntervalInMin)

	switch c.Args().Len() {
	case 1:
		directory = c.Args().First()
	case 0:
		_ = cli.ShowAppHelp(c)
		return fmt.Errorf("expected directory as argument")
	default:
		_ = cli.ShowAppHelp(c)
		return fmt.Errorf("unexpected argument(s) found: %v", c.Args().Slice()[1:])
	}

	args := deletion.Args{Directory: directory, MaxAgeInHours: maxAgeInHours}

	loopStopper := registerUnixSignals()
	defer close(loopStopper)

	fmt.Println("[tempdel] Start delete-loop...")
	runDeletionLoop(args, loopInterval, loopStopper)

	return nil
}

func minuteToDuration(min int) time.Duration {
	return time.Duration(min) * time.Minute
}

// registerUnixSignals listens to different unix signals (that Docker or a user might cause) and returns a semaphore
// channel. A value sent through this channel should stop the deletion loop.
func registerUnixSignals() (loopStopper chan bool) {
	loopStopper = make(chan bool, 1)
	procSignals := make(chan os.Signal, 1)

	signal.Notify(procSignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		for range procSignals {
			loopStopper <- true
			break
		}
		fmt.Println("[tempdel] Caught signal...")
		close(procSignals)
	}()

	return loopStopper
}

// the interval is here chosen for seconds for reasons of unit test duration
func runDeletionLoop(args deletion.Args, intervalInSecs time.Duration, loopStopper chan bool) {
	var err error
	ticker := time.NewTicker(intervalInSecs)

	for {
		select {
		case <-loopStopper:
			ticker.Stop()
			fmt.Println("[tempdel] Exiting tempdel...")
			return
		case <-ticker.C:
			log.Debug("[tempdel] Start deletion run...")
			err = deleteFilesWithArgs(args)
			if err != nil {
				log.Errorf("[tempdel] Deleting files failed with this error: %s", err.Error())
			}
			log.Debug("[tempdel] End deletion run.")
		default:
			// reduces CPU load but stretches reaction to unix signals
			time.Sleep(cpuLoadSleepInSec * time.Second)
		}
	}
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
