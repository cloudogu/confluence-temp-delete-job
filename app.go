package main

import (
	"fmt"
	"github.com/cloudogu/confluence-temp-delete-job/cmd"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

// Version indicates if the program can be used in production or should only
// be used in development environments
// This information will be overwritten during the build process
var Version = "development"
var log = logging.MustGetLogger("main")
var appExiter exiter = &defaultExiter{}

func main() {
	app := cli.NewApp()
	app.Name = "tempdel"
	app.Version = Version
	app.Usage = "Delete files in a given directory and a given file age."
	app.Description = "Cloudogu confluence temp delete job"

	app.Commands = []*cli.Command{
		cmd.DeleteFilesCommand,
	}

	app.Flags = createGlobalFlags()
	app.Before = configureLogging

	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("%+v\n", err)
		appExiter.exit(1)
	}
}

func createGlobalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "log-level",
			Usage: "define log level",
			Value: "notice"},
	}
}

// logging format
var format = logging.MustStringFormatter(
	`%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
)

func configureLogging(c *cli.Context) error {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	logLevel, err := logging.LogLevel(c.String("log-level"))
	if err != nil {
		fmt.Printf("%s: invalid log level specified, please use critical, error, warning, notice, info or debug", time.Now().Format(time.RFC3339))
		return errors.Wrap(err, "failed to configure logging")
	}
	logging.SetLevel(logLevel, "")
	return nil
}

// exiter is an interface that enables testing this app without exiting any test code.
type exiter interface {
	exit(exitCode int)
}

type defaultExiter struct{}

func (*defaultExiter) exit(exitCode int) {
	os.Exit(exitCode)
}
