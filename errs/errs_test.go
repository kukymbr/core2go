package errs_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/kukymbr/core2go/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExposable(t *testing.T) {
	tests := []struct {
		Code         int
		Messages     []any
		ExpectedText string
		ExpectedCode int
	}{
		{http.StatusBadRequest, nil, "Bad Request", 400},
		{0, nil, "Internal Server Error", 500},
		{0, []any{"Test message"}, "Test message", 500},
		{http.StatusTeapot, nil, "I'm a teapot", 418},
		{http.StatusBadRequest, []any{"test1: ", "test2 ", 300}, "test1: test2 300", 400},
	}

	for _, test := range tests {
		err := errs.NewExposable(test.Code, test.Messages...)

		var expErr *errs.ExposableError
		ok := errors.As(err, &expErr)
		require.True(t, ok)

		assert.Error(t, err)
		assert.Equal(t, test.ExpectedText, err.Error())
		assert.Equal(t, test.ExpectedCode, expErr.Code)
	}
}
