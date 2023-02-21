package n

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_pathFixer(t *testing.T) {
	require.NotNil(t, _pathFixer)
	assert.Equal(t, "/base/path", _pathFixer.Replace("/base/path"))
	assert.Equal(t, "/base/path", _pathFixer.Replace("\\base\\path"))
}

func TestNewDirRoute(t *testing.T) {
	var (
		fullPath string
	)

	r := NewDirRoute("/base/path/", func(path string, _ Handler) Route {
		fullPath = path
		return nil
	})

	r.Handle("/route", HandlerFunc(func(ctx Context) error {
		return nil
	}))

	assert.Equal(t, "/base/path/route", fullPath)
}

func TestNDirRoute_Handle(t *testing.T) {

}
