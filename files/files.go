// Package files provides config files and customizations for buf.
package files

import (
	_ "embed"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	//go:embed buf.yaml.tmpl
	baseBufYaml string
	//go:embed buf.work.yaml.tmpl
	baseBufWorkYaml string
	bufWorkYamlTmpl = template.Must(template.New("buf.work.yaml").Parse(baseBufWorkYaml))
	//go:embed buf.gen.yaml.tmpl
	baseBufGenYaml string
	bufGenYamlTmpl = template.Must(template.New("buf.gen.yaml").Parse(baseBufGenYaml))
)

type params = map[string]interface{}

func WriteBufYaml(path string) error {
	file, err := os.OpenFile(filepath.Join(path, "buf.yaml"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := io.Copy(file, strings.NewReader(baseBufYaml)); err != nil {
		return err
	}
	return nil
}

func WriteBufWorkYaml(path string, dirs ...string) error {
	file, err := os.OpenFile(filepath.Join(path, "buf.work.yaml"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	param := params{
		"Directories": dirs,
	}

	if err := bufWorkYamlTmpl.Execute(file, &param); err != nil {
		return err
	}
	return nil
}

func WriteBufGenYaml(path string, config *BufGenYaml) error {
	file, err := os.OpenFile(filepath.Join(path, "buf.gen.yaml"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	reader, err := genBufGenYaml(config)
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, reader); err != nil {
		return err
	}
	return nil
}

func genBufGenYaml(config *BufGenYaml) (io.Reader, error) {
	var buf strings.Builder
	if err := bufGenYamlTmpl.Execute(&buf, config); err != nil {
		return nil, err
	}
	return strings.NewReader(buf.String()), nil
}
