package cmd

import (
	"bytes"
	"github.com/cloudogu/confluence-temp-delete-job/deletion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_deleteFilesWithArgs(t *testing.T) {
	realStdout := os.Stdout

	t.Run("should fail with missing directory parameter", func(t *testing.T) {
		// when
		err := deleteFilesWithArgs(deletion.Args{
			Directory:     "",
			MaxAgeInHours: 0,
		})

		// then
		require.Error(t, err)
		require.Contains(t, err.Error(), "directory must not be empty")
	})
	t.Run("should succeed", func(t *testing.T) {
		dir, _ := ioutil.TempDir(os.TempDir(), "tempdel-")
		defer func() { _ = os.RemoveAll(dir) }()
		defer restoreOriginalStdout(realStdout)

		fakeReaderPipe, fakeWriterPipe := routeStdoutToReplacement()

		// when
		err := deleteFilesWithArgs(deletion.Args{
			Directory:     dir,
			MaxAgeInHours: 12,
		})

		// then
		require.NoError(t, err)

		actualOutput := captureOutput(fakeReaderPipe, fakeWriterPipe, realStdout)
		assert.Equal(t, "[tempdel] deleted: 0 (0 MB), skipped: 0, failed: 0\n", actualOutput)
	})
}

func Test_runDeletionLoop(t *testing.T) {
	realStdout := os.Stdout

	t.Run("should succeed", func(t *testing.T) {
		dir, _ := ioutil.TempDir(os.TempDir(), "tempdel-")
		defer func() { _ = os.RemoveAll(dir) }()
		defer restoreOriginalStdout(realStdout)
		stopChan := make(chan bool, 1)

		fakeReaderPipe, fakeWriterPipe := routeStdoutToReplacement()
		intervalInSec := 1 * time.Second
		args := deletion.Args{
			Directory:     dir,
			MaxAgeInHours: 12,
		}

		// when
		go runDeletionLoop(args, intervalInSec, stopChan)

		// stop when loop ran 1x
		time.Sleep(intervalInSec + cpuLoadSleepInSec*time.Second + 1*time.Second)
		stopChan <- true

		// then
		actualOutput := captureOutput(fakeReaderPipe, fakeWriterPipe, realStdout)
		assert.Contains(t, actualOutput, "[tempdel] deleted: 0 (0 MB), skipped: 0, failed: 0\n")
	})
}

func Test_registerUnixSignals(t *testing.T) {
	realStdout := os.Stdout

	t.Run("should stop listening and send true through the stop channel", func(t *testing.T) {
		defer restoreOriginalStdout(realStdout)
		fakeReaderPipe, fakeWriterPipe := routeStdoutToReplacement()

		thisTestProcess, procErr := os.FindProcess(os.Getpid())
		if procErr != nil {
			t.Fatal(procErr)
		}

		// when
		stopChan := registerUnixSignals()

		// then
		sigErr := thisTestProcess.Signal(os.Interrupt)
		println("Sending SIGINT")
		time.Sleep(2 * time.Second)
		assert.NoError(t, sigErr)
		assert.True(t, <-stopChan)
		actualOutput := captureOutput(fakeReaderPipe, fakeWriterPipe, realStdout)
		assert.Contains(t, actualOutput, "[tempdel] Caught signal...\n")
	})
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
	_ = fakeWriterPipe.Close()
	restoreOriginalStdout(originalStdout)

	actualOutput := <-outC

	return actualOutput
}

func restoreOriginalStdout(stdout *os.File) {
	os.Stdout = stdout
}

func Test_minToSec(t *testing.T) {
	t.Run("should convert 1 minute to a duration of 60 seconds", func(t *testing.T) {
		actual := minuteToDuration(1)

		assert.Equal(t, time.Duration(60)*time.Second, actual)
	})
	t.Run("should convert 60 minutes to a duration of 1 hour", func(t *testing.T) {
		actual := minuteToDuration(60)

		assert.Equal(t, time.Duration(1)*time.Hour, actual)
	})
}
