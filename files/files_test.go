package files

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBufWorkYaml(t *testing.T) {
	param := params{
		"Directories": []string{"a", "b"},
	}

	var out strings.Builder
	assert.NoError(t, bufWorkYamlTmpl.Execute(&out, param))

	expected := `version: v1
directories:
  - a
  - b
`
	assert.Equal(t, expected, out.String())
}

func TestBufGenYaml(t *testing.T) {
	param := params{
		"Plugins": []params{
			{
				"Name": "test",
				"Out":  "some/dir",
			},
			{
				"Name": "test2",
				"Out":  "some/dir",
			},
		},
	}

	expected := `version: v1
plugins:
  - name: test
    out: some/dir
  - name: test2
    out: some/dir
`

	var buf strings.Builder
	assert.NoError(t, bufGenYamlTmpl.Execute(&buf, &param))

	assert.Equal(t, expected, buf.String())
}
