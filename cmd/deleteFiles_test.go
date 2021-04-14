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
		args := deletion.Args{
			Directory:     dir,
			MaxAgeInHours: 12,
		}

		// when
		go runDeletionLoop(args, 1, stopChan)

		// stop when loop ran 1x
		time.Sleep(2 * time.Second)
		stopChan <- true

		// then
		actualOutput := captureOutput(fakeReaderPipe, fakeWriterPipe, realStdout)
		assert.Equal(t, "[tempdel] deleted: 0 (0 MB), skipped: 0, failed: 0\n", actualOutput)
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
	fakeWriterPipe.Close()
	restoreOriginalStdout(originalStdout)

	actualOutput := <-outC

	return actualOutput
}

func restoreOriginalStdout(stdout *os.File) {
	os.Stdout = stdout
}
