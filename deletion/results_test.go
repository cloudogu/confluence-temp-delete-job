package deletion

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestResults_fail(t *testing.T) {
	sut := &Results{}

	// when
	for i := 0; i < 9; i++ {
		sut.fail(fmt.Sprintf("file /file_%d", i), assert.AnError)
	}

	assert.Equal(t, 9, sut.failed)
}

func TestResults_pass(t *testing.T) {
	sut := &Results{}
	arbitraryInfo, _ := os.Stat(".")

	// when
	for i := 0; i < 9; i++ {
		sut.pass(fmt.Sprintf("file /file_%d", i), arbitraryInfo)
	}

	assert.Equal(t, 9, sut.deleted)
}

func TestResults_skip(t *testing.T) {
	sut := &Results{}

	// when
	for i := 0; i < 9; i++ {
		sut.skip(fmt.Sprintf("file /file_%d", i))
	}

	assert.Equal(t, 9, sut.skipped)
}

func TestResults_PrintStats(t *testing.T) {
	t.Run("should print stats nicely", func(t *testing.T) {
		realStdout := os.Stdout
		defer restoreOriginalStdout(realStdout)
		fakeReaderPipe, fakeWriterPipe := routeStdoutToReplacement()

		sut := &Results{
			deleted:       20,
			deletedSizeKB: 24_890,
			failed:        1,
			skipped:       8,
		}

		// when
		sut.PrintStats()

		// then
		actual := captureOutput(fakeReaderPipe, fakeWriterPipe, realStdout)
		assert.Equal(t, "[tempdel] deleted: 20 (24 MB), skipped: 8, failed: 1\n", actual)
	})
}

func TestResults(t *testing.T) {
	oldTime := time.Now().Add(-20 * time.Hour)

	t.Run("should count small stats correctly", func(t *testing.T) {
		startDir, _ := ioutil.TempDir(os.TempDir(), "tempdel-")
		defer func() { _ = os.RemoveAll(startDir) }()
		file1 := createFileWithTime(t, startDir, "del-file", oldTime)
		file2 := createFileWithTime(t, startDir, "del-file", oldTime)
		file3 := createFileWithTime(t, startDir, "del-file", oldTime)
		writeBytesToFile(t, file1, 1234)
		writeBytesToFile(t, file2, 1023)
		writeBytesToFile(t, file3, 2048*1024)
		info1 := fileInfo(t, file1)
		info2 := fileInfo(t, file2)
		info3 := fileInfo(t, file3)

		sut := Results{}

		// when
		sut.pass(file1, info1)
		sut.pass(file2, info2)
		sut.pass(file3, info3)

		// then
		expected := Results{
			deleted:       3,
			deletedSizeKB: 2049,
			failed:        0,
			skipped:       0,
		}
		assert.Equal(t, expected, sut)
	})
	t.Run("should count large stats correctly", func(t *testing.T) {
		startDir, _ := ioutil.TempDir(os.TempDir(), "tempdel-")
		defer func() { _ = os.RemoveAll(startDir) }()
		file1 := createFileWithTime(t, startDir, "del-file", oldTime)
		file2 := createFileWithTime(t, startDir, "del-file", oldTime)
		file3 := createFileWithTime(t, startDir, "del-file", oldTime)
		writeBytesToFile(t, file1, 1*1234*1024)
		writeBytesToFile(t, file2, 1*1023*1024)
		writeBytesToFile(t, file3, 10*1024*1024)
		info1 := fileInfo(t, file1)
		info2 := fileInfo(t, file2)
		info3 := fileInfo(t, file3)

		sut := Results{}

		// when
		sut.pass(file1, info1)
		sut.pass(file2, info2)
		sut.pass(file3, info3)

		// then
		expected := Results{
			deleted:       3,
			deletedSizeKB: 12497,
			failed:        0,
			skipped:       0,
		}
		assert.Equal(t, expected, sut)
	})
}

func writeBytesToFile(t *testing.T, path string, amount int) {
	t.Helper()

	content := make([]byte, amount)
	err := ioutil.WriteFile(path, content, 0644)
	assert.NoError(t, err, "error while writing bytes to "+path)
}

func fileInfo(t *testing.T, path string) os.FileInfo {
	t.Helper()

	info, err := os.Stat(path)
	assert.NoError(t, err, "error while getting file info on "+path)
	return info
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
