package main

import (
	"fmt"
	"github.com/drognisep/initbuf/files"
	flag "github.com/spf13/pflag"
	"os"
	"path/filepath"
)

func main() {
	config := struct {
		goOut             string
		genGrpc           bool
		genJson           bool
		managedGeneration bool
		protoDirs         []string
	}{}

	flag.StringVar(&config.goOut, "go-out", "", "Sets the Go output directory")
	flag.BoolVar(&config.genGrpc, "gen-grpc", false, "Use the protoc-gen-go-grpc plugin")
	flag.BoolVar(&config.genJson, "gen-json", false, "Use the protoc-gen-go-json plugin")
	flag.BoolVar(&config.managedGeneration, "managed", true, "Sets whether code generation should be managed")
	flag.StringSliceVar(&config.protoDirs, "proto-dir", nil, "Adds a protobuf root directory. Should be relative to the module root.")
	flag.Parse()

	root, err := files.ModuleRoot()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err := files.WriteBufWorkYaml(root, config.protoDirs...); err != nil {
		fmt.Printf("Error generating buf.work.yaml in %s: %v", root, err)
		os.Exit(1)
	}

	for _, dir := range config.protoDirs {
		bufPath := filepath.Join(root, dir)
		bufPathDir := filepath.Dir(bufPath)
		if err := os.MkdirAll(bufPathDir, 0750); err != nil {
			fmt.Printf("Failed to create %s: %v\n", bufPathDir, err)
			os.Exit(1)
		}
		if err := files.WriteBufYaml(bufPath); err != nil {
			fmt.Printf("Error generating buf.yaml %s: %v\n", bufPath, err)
			os.Exit(1)
		}
	}

	genYaml := files.BufGenYaml{
		Version: "v1",
		Managed: files.ManagedRules{},
		Plugins: []files.GenerationPlugin{
			{
				Name:       "go",
				OutputPath: config.goOut,
				Options:    []string{"paths=source_relative"},
			},
		},
		GoOut:     config.goOut,
		UseGoGrpc: config.genGrpc,
		UseGoJson: config.genJson,
	}
	if config.managedGeneration {
		genYaml.Managed["enabled"] = true
		if err := genYaml.SetDefaultGoPrefix(config.goOut); err != nil {
			fmt.Printf("Error getting default go prefix: %v\n", err)
			os.Exit(1)
		}
	}
	if err := files.WriteBufGenYaml(root, &genYaml); err != nil {
		fmt.Printf("Error generating buf.gen.yaml: %v\n", err)
	}

	deps := `Successfully generated base configuration!

Make sure these dependencies are installed:
 - buf (https://docs.buf.build/installation)
 - google.golang.org/protobuf (add to this module)
 - google.golang.org/protobuf/cmd/protoc-gen-go
`
	if config.genGrpc {
		deps += " - google.golang.org/grpc (add to this module)\n"
		deps += " - google.golang.org/grpc/cmd/protoc-gen-go-grpc\n"
	}
	if config.genJson {
		deps += " - github.com/mitchellh/protoc-gen-go-json\n"
	}
	fmt.Println(deps)
}
