package files

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCurrentGoModuleCwd(t *testing.T) {
	module, err := getCurrentGoModuleCwd()
	assert.NoError(t, err)
	assert.Equal(t, "github.com/drognisep/initbuf", module)
}

func TestGetCurrentGoModule(t *testing.T) {
	module, err := getCurrentGoModule("../..")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotAModule), "should have been ErrNotAModule")
	assert.Equal(t, "", module)
}
