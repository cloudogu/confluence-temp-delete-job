package deletion

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const testMaxAgeInHours = 12

func TestNew(t *testing.T) {
	t.Run("should set args", func(t *testing.T) {
		input := Args{Directory: "/test", MaxAgeInHours: 42}

		sut, _ := New(input)

		assert.Equal(t, "/test", sut.Directory)
		assert.Equal(t, 42, sut.MaxAgeInHours)
	})

	type args struct {
		args Args
	}
	tests := []struct {
		name        string
		args        args
		wantDeleter bool
		wantErr     bool
	}{
		{"should pass", args{Args{Directory: "/test", MaxAgeInHours: 12}}, true, false},
		{"should pass with 0 age", args{Args{Directory: "/test", MaxAgeInHours: 0}}, true, false},
		{"should fail with dir", args{Args{Directory: "", MaxAgeInHours: 12}}, false, true},
		{"should fail with age", args{Args{Directory: "/a", MaxAgeInHours: -1}}, false, true},
		{"should fail with both dir and age", args{Args{Directory: "", MaxAgeInHours: -1}}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got != nil) != tt.wantDeleter {
				t.Errorf("New() deleter = %v, wantDeleter %v", got, tt.wantDeleter)
				return
			}
		})
	}
}

func Test_deleter_deleteFile(t *testing.T) {
	t.Run("should delete file and update stats", func(t *testing.T) {
		// given
		dir, _ := ioutil.TempDir(os.TempDir(), "tempdel-")
		defer func() { _ = os.RemoveAll(dir) }()
		oldness := nowClock.Now().Add(-20 * time.Hour)
		file := createFileWithTime(t, dir, "", oldness)

		sut, _ := New(Args{Directory: dir, MaxAgeInHours: testMaxAgeInHours})
		assert.Empty(t, sut.Results)

		// when
		err := sut.deleteFile(file)

		// then
		require.NoError(t, err)
		assert.Equal(t, 1, sut.Results.passed)
		assert.Equal(t, 0, sut.Results.failed)
		assert.Equal(t, 0, sut.Results.skipped)
	})

	t.Run("should error on file error and update stats", func(t *testing.T) {
		// given
		const path = "/some/error/will/occur/here"
		removerMock := &mockFileRemover{}
		removerMock.On("Remove", path).Return(os.ErrNotExist)
		remover = removerMock
		defer func() { remover = &realFileRemover{} }()

		sut, _ := New(Args{Directory: "dir", MaxAgeInHours: testMaxAgeInHours})

		// when
		err := sut.deleteFile(path)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "file does not exist")
		assert.Equal(t, 0, sut.Results.passed)
		assert.Equal(t, 1, sut.Results.failed)
		assert.Equal(t, 0, sut.Results.skipped)
		(remover).(*mockFileRemover).AssertExpectations(t)
	})
}

func Test_fileOlderThan(t *testing.T) {
	nowClock = &testClock{time.Now()}

	t.Run("should return true for ancient files", func(t *testing.T) {
		theBeginningOfComputing := time.Unix(0, 0)

		// when
		actual := fileOlderThan(testMaxAgeInHours, theBeginningOfComputing)

		// then
		assert.True(t, actual)
	})
	t.Run("should return true for 12 hours and 1 second old files", func(t *testing.T) {
		minus12HoursTime := nowClock.Now().Add(-12 * time.Hour).Add(-1 * time.Second)

		// when
		actual := fileOlderThan(testMaxAgeInHours, minus12HoursTime)

		// then
		assert.True(t, actual)
	})
	t.Run("should return false for exactly 12 hours old files", func(t *testing.T) {
		minus12HoursTime := nowClock.Now().Add(-12 * time.Hour)

		// when
		actual := fileOlderThan(testMaxAgeInHours, minus12HoursTime)

		// then
		assert.False(t, actual)
	})
	t.Run("should return false for 11 hours and 59 minutes old files", func(t *testing.T) {
		minus11Hours59SecTime := nowClock.Now().Add(-11 * time.Hour).Add(-59 * time.Second)

		// when
		actual := fileOlderThan(testMaxAgeInHours, minus11Hours59SecTime)

		// then
		assert.False(t, actual)
	})
	t.Run("should return false for 11 hours old files", func(t *testing.T) {
		minus11HoursTime := nowClock.Now().Add(-11 * time.Hour)

		// when
		actual := fileOlderThan(testMaxAgeInHours, minus11HoursTime)

		// then
		assert.False(t, actual)
	})
	t.Run("should return false for 0 hours old files", func(t *testing.T) {
		// when
		actual := fileOlderThan(testMaxAgeInHours, nowClock.Now())

		// then
		assert.False(t, actual)
	})
	t.Run("should return false for files back from the future", func(t *testing.T) {
		theFuture := nowClock.Now().Add(24 * 365 * time.Hour)

		// when
		actual := fileOlderThan(testMaxAgeInHours, theFuture)

		// then
		assert.False(t, actual)
	})
}

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

	// when
	for i := 0; i < 9; i++ {
		sut.pass(fmt.Sprintf("file /file_%d", i))
	}

	assert.Equal(t, 9, sut.passed)
}

func TestResults_skip(t *testing.T) {
	sut := &Results{}

	// when
	for i := 0; i < 9; i++ {
		sut.skip(fmt.Sprintf("file /file_%d", i))
	}

	assert.Equal(t, 9, sut.skipped)
}

func Test_deleter_filterOldFiles(t *testing.T) {
	t.Run("should skip directories", func(t *testing.T) {
		dir, _ := ioutil.TempDir(os.TempDir(), "tempdel-")
		defer func() { _ = os.RemoveAll(dir) }()
		innerDir, _ := ioutil.TempDir(dir, "tempdel-")

		sut, _ := New(Args{Directory: dir, MaxAgeInHours: testMaxAgeInHours})
		innerDirStats, _ := os.Stat(innerDir)

		// when
		err := sut.filterOldFiles(innerDir, innerDirStats, nil)

		require.NoError(t, err)
		_, err = os.Stat(innerDir)
		assert.NoError(t, err)
	})
}

func Test_deleter_Execute(t *testing.T) {
	t.Run("should delete old files and leave new files", func(t *testing.T) {
		// given
		startDir, _ := ioutil.TempDir(os.TempDir(), "tempdel-")
		defer func() { _ = os.RemoveAll(startDir) }()
		// Name files ABC... because fileWalk iterates files alphabetically
		oldTime := nowClock.Now().Add(-testMaxAgeInHours - 20*time.Hour)
		newTime := nowClock.Now().Add(-2 * time.Hour)
		deleteFile1 := createFileWithTime(t, startDir, "a-", oldTime)
		leaveFile1 := createFileWithTime(t, startDir, "b-", newTime)
		deleteFile2 := createFileWithTime(t, startDir, "c-", oldTime)
		leaveFile2 := createFileWithTime(t, startDir, "d-", newTime)

		sut, _ := New(Args{Directory: startDir, MaxAgeInHours: testMaxAgeInHours})

		// when
		actual, err := sut.Execute()

		// then
		require.NoError(t, err)
		expectedStats := Results{
			passed:  2,
			failed:  0,
			skipped: 2,
		}
		assert.Equal(t, expectedStats, *actual)
		assertFileNotExists(t, deleteFile1)
		assertFileNotExists(t, deleteFile2)
		assertFileExists(t, leaveFile1)
		assertFileExists(t, leaveFile2)
	})
	t.Run("should delete empty directories and leave new files", func(t *testing.T) {
		// given
		startDir, _ := ioutil.TempDir(os.TempDir(), "tempdel-")
		defer func() { _ = os.RemoveAll(startDir) }()
		// Name files ABC... because fileWalk iterates files alphabetically
		oldTime := nowClock.Now().Add(-20 * time.Hour)
		newTime := nowClock.Now().Add(-2 * time.Hour)

		// note that the dir timestamp will change upon file deletion
		deleteDir1, _ := ioutil.TempDir(startDir, "a-del-dir-")
		_ = os.Chtimes(deleteDir1, oldTime, oldTime)

		deleteFile2 := createFileWithTime(t, deleteDir1, "a-del-file", oldTime)
		leaveFile1 := createFileWithTime(t, startDir, "b-stay-file", newTime)
		deleteFile3 := createFileWithTime(t, startDir, "c-del-file", oldTime)
		leaveFile2 := createFileWithTime(t, startDir, "d-stay", newTime)

		sut, _ := New(Args{Directory: startDir, MaxAgeInHours: testMaxAgeInHours})

		// when
		actual, err := sut.Execute()

		// then
		require.NoError(t, err)
		expectedStats := Results{
			passed:  3,
			failed:  0,
			skipped: 2,
		}
		assert.Equal(t, expectedStats, *actual)
		assertFileNotExists(t, deleteDir1)
		assertFileNotExists(t, deleteFile2)
		assertFileNotExists(t, deleteFile3)
		assertFileExists(t, leaveFile1)
		assertFileExists(t, leaveFile2)
	})
	t.Run("should not delete non-empty directories that contain new files", func(t *testing.T) {
		// given
		startDir, _ := ioutil.TempDir(os.TempDir(), "tempdel-")
		defer func() { _ = os.RemoveAll(startDir) }()
		// Name files ABC... because fileWalk iterates files alphabetically
		oldTime := nowClock.Now().Add(-20 * time.Hour)
		newTime := nowClock.Now().Add(-2 * time.Hour)

		// note that the dir timestamp will change upon file deletion
		leaveDir1, _ := ioutil.TempDir(startDir, "a-stay-dir-")
		_ = os.Chtimes(leaveDir1, oldTime, oldTime)

		deleteFile2 := createFileWithTime(t, leaveDir1, "a-del-file", oldTime)
		leaveFile1 := createFileWithTime(t, leaveDir1, "b-stay-file", newTime)
		deleteFile3 := createFileWithTime(t, startDir, "c-del-file", oldTime)
		leaveFile2 := createFileWithTime(t, startDir, "d-stay", newTime)

		sut, _ := New(Args{Directory: startDir, MaxAgeInHours: testMaxAgeInHours})

		// when
		actual, err := sut.Execute()

		// then
		require.NoError(t, err)
		expectedStats := Results{
			passed:  2,
			failed:  0,
			skipped: 3,
		}
		assert.Equal(t, expectedStats, *actual)
		assertFileNotExists(t, deleteFile2)
		assertFileNotExists(t, deleteFile3)
		assertFileExists(t, leaveDir1)
		assertFileExists(t, leaveFile1)
		assertFileExists(t, leaveFile2)
	})
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()

	_, err := os.Stat(path)
	assert.NoError(t, err)
}

func assertFileNotExists(t *testing.T, path string) {
	t.Helper()

	_, err := os.Stat(path)
	assert.Error(t, err, "a NotExistError is expected to show that the file object was properly deleted")
	// expect only NotExistErr and fail all other errors
	if !os.IsNotExist(err) {
		assert.NoError(t, err)
	}
}

type testClock struct {
	desiredTime time.Time
}

func (t *testClock) Now() time.Time {
	return t.desiredTime
}

func createFileWithTime(t *testing.T, directory string, filenamePrefix string, fileTime time.Time) string {
	t.Helper()

	file, _ := ioutil.TempFile(directory, filenamePrefix+"tempdel-")
	filePath := file.Name()

	err := os.Chtimes(filePath, fileTime, fileTime)
	assert.NoError(t, err)

	return filePath
}

type mockFileRemover struct {
	mock.Mock
}

func (m *mockFileRemover) Remove(path string) error {
	args := m.Called(path)
	return args.Error(0)
}
