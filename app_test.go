package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_main(t *testing.T) {
	exiterMock := new(mockExiter)
	appExiter = exiterMock
	exiterMock.On("exit", 1)

	main()

	exiterMock.AssertExpectations(t)
}

func Test_newExiter(t *testing.T) {
	t.Run("should create an exiter instance", func(t *testing.T) {
		sut := &defaultExiter{}

		require.NotNil(t, sut)
		assert.Implements(t, (*exiter)(nil), sut)
	})
}

// test util stuff
type mockExiter struct {
	mock.Mock
}

func (m *mockExiter) exit(exitCode int) {
	m.Called(exitCode)
}
