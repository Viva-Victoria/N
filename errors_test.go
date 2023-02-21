package n

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBadRequestError_Error(t *testing.T) {
	e := NewBadRequestError(nil)
	assert.Equal(t, "", e.Error())

	e = NewBadRequestError(context.DeadlineExceeded)
	assert.Equal(t, context.DeadlineExceeded.Error(), e.Error())
}

func TestBadRequestError_Unwrap(t *testing.T) {
	e := NewBadRequestError(nil)
	assert.Nil(t, e.Unwrap())

	e = NewBadRequestError(context.Canceled)
	assert.Error(t, context.Canceled, e)
}