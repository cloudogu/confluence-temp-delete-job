package cmd

import (
	"bytes"
	"flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func Test_deleteFiles(t *testing.T) {
	realStdout := os.Stdout

	t.Run("should fail with missing directory parameter", func(t *testing.T) {
		defer restoreOriginalStdout(realStdout)

		sut := DeleteFilesCommand
		fakeReaderPipe, fakeWriterPipe := routeStdoutToReplacement()

		var cliArgs []string
		c := getTestCliContextForCommand(t, sut, fakeWriterPipe, cliArgs)

		// when
		err := sut.Action(c)

		// then
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected directory")

		actualOutput := captureOutput(fakeReaderPipe, fakeWriterPipe, realStdout)
		assert.Contains(t, actualOutput, "NAME:")
		assert.Contains(t, actualOutput, "USAGE:")
		assert.Contains(t, actualOutput, "OPTIONS:")
		assert.Contains(t, actualOutput, flagMaxAgeHoursLong)
	})
	t.Run("should fail with superfluous parameters", func(t *testing.T) {
		defer restoreOriginalStdout(realStdout)

		sut := DeleteFilesCommand
		fakeReaderPipe, fakeWriterPipe := routeStdoutToReplacement()

		cliArgs := []string{"/path/to/dir", "oops"}
		c := getTestCliContextForCommand(t, sut, fakeWriterPipe, cliArgs)

		// when
		err := sut.Action(c)

		// then
		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected argument(s) found: [oops]")

		actualOutput := captureOutput(fakeReaderPipe, fakeWriterPipe, realStdout)
		assert.Contains(t, actualOutput, "NAME:")
		assert.Contains(t, actualOutput, "USAGE:")
		assert.Contains(t, actualOutput, "OPTIONS:")
		assert.Contains(t, actualOutput, flagMaxAgeHoursLong)
	})
	t.Run("should succeed", func(t *testing.T) {
		dir, _ := ioutil.TempDir(os.TempDir(), "tempdel-")
		defer func() { _ = os.RemoveAll(dir) }()
		defer restoreOriginalStdout(realStdout)

		sut := DeleteFilesCommand
		fakeReaderPipe, fakeWriterPipe := routeStdoutToReplacement()

		cliArgs := []string{dir}
		c := getTestCliContextForCommand(t, sut, fakeWriterPipe, cliArgs)

		// when
		err := sut.Action(c)

		// then
		require.NoError(t, err)

		actualOutput := captureOutput(fakeReaderPipe, fakeWriterPipe, realStdout)
		assert.Empty(t, actualOutput)
	})
}

func getTestCliContextForCommand(t *testing.T, command *cli.Command, stdout *os.File, args []string) *cli.Context {
	t.Helper()

	// setup CLI internals so urfave usage help prints like a production cesapp
	app := cli.NewApp()
	app.Name = "tempdel"
	app.HelpName = app.Name
	app.Usage = "usage"

	// overwrite with a writer pipe to conveniently capture urfave usage help output
	app.Writer = stdout

	// finishing touches to CLI app and context so urfave can properly apply all given information
	app.Commands = []*cli.Command{command}
	app.Setup()
	flags := flag.NewFlagSet(app.Name, flag.ContinueOnError)
	_ = flags.Parse(args)

	c := cli.NewContext(app, flags, nil)
	c.Command = command

	return c
}

func routeStdoutToReplacement() (readerPipe, writerPipe *os.File) {
	r, w, _ := os.Pipe()
	os.Stdout = w

	return r, w
}

func captureOutput(fakeReaderPipe, fakeWriterPipe, originalStdout *os.File) string {
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, fakeReaderPipe)
		outC <- buf.String()
	}()

	// back to normal state
	fakeWriterPipe.Close()
	restoreOriginalStdout(originalStdout)

	actualOutput := <-outC

	return actualOutput
}

func restoreOriginalStdout(stdout *os.File) {
	os.Stdout = stdout
}
