package files

import (
	"github.com/stretchr/testify/assert"
	"io"
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

func TestGenBufGenYaml(t *testing.T) {
	const expected = `version: v2
managed:
  enabled: true
  go_package_prefix:
    default: testmodule/out
plugins:
  - name: go-grpc # go get google.golang.org/grpc@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    out: out
    opt: paths=source_relative
  - name: go-json # go install github.com/mitchellh/protoc-gen-go-json@latest
    out: out
    opt: paths=source_relative
  - name: go
    out: out
    opt: paths=source_relative
`
	config := &BufGenYaml{
		Version: "v2",
		Managed: ManagedRules{
			"enabled": true,
			"go_package_prefix": ManagedRules{
				"default": "testmodule/out",
			},
		},
		Plugins: []GenerationPlugin{
			{
				Name:       "go",
				OutputPath: "out",
				Options:    []string{"paths=source_relative"},
			},
		},
		UseGoGrpc: true,
		UseGoJson: true,
		GoOut:     "out",
	}

	var buf strings.Builder
	reader, err := genBufGenYaml(config)
	assert.NoError(t, err)
	_, err = io.Copy(&buf, reader)
	assert.NoError(t, err)
	assert.Equal(t, expected, buf.String())
}
