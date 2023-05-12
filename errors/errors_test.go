package errors_test

import (
	"net/http"
	"testing"

	"github.com/kukymbr/core2go/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewExposable(t *testing.T) {
	tests := []struct {
		Code         int
		Messages     []any
		ExpectedText string
	}{
		{http.StatusBadRequest, nil, "Bad Request"},
		{http.StatusTeapot, nil, "I'm a teapot"},
		{http.StatusBadRequest, []any{"test1: ", "test2 ", 300}, "test1: test2 300"},
	}

	for _, test := range tests {
		err := errors.NewExposable(test.Code, test.Messages...)

		assert.Error(t, err)
		assert.Equal(t, test.ExpectedText, err.Error())
	}
}
