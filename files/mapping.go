package files

import (
	"path/filepath"
	"strings"
)

type GenerationPlugin struct {
	Name       string   `yaml:"name"`
	OutputPath string   `yaml:"out"`
	Options    []string `yaml:"opt,omitempty"`
}

func (p *GenerationPlugin) OptionString() string {
	return strings.Join(p.Options, ",")
}

type ManagedRules = map[string]interface{}

type BufGenYaml struct {
	Version   string
	Managed   ManagedRules
	Plugins   []GenerationPlugin
	UseGoGrpc bool
	UseGoJson bool
	GoOut     string
}

func (r *BufGenYaml) HasPlugins() bool {
	return len(r.Plugins) > 0 || r.UseGoGrpc || r.UseGoJson
}

func (r *BufGenYaml) IsManaged() bool {
	value, _ := r.BoolManagedValue("enabled")
	return value
}

func (r *BufGenYaml) BoolManagedValue(key string) (bool, bool) {
	iVal, hasKey := r.Managed[key]
	if !hasKey {
		return false, false
	}
	val, ok := iVal.(bool)
	if !ok {
		return false, false
	}
	return val, true
}

func (r *BufGenYaml) StringManagedValue(key string) (string, bool) {
	iVal, hasKey := r.Managed[key]
	if !hasKey {
		return "", false
	}
	val, ok := iVal.(string)
	if !ok {
		return "", false
	}
	return val, true
}

func (r *BufGenYaml) SetDefaultGoPrefix(outputDir string) error {
	module, err := getCurrentGoModuleCwd()
	if err != nil {
		return err
	}
	module = module + "/" + filepath.ToSlash(outputDir)
	r.Managed["go_package_prefix"] = ManagedRules{
		"default": module,
	}
	return nil
}

func (r *BufGenYaml) DefaultGoPrefix() string {
	const defaultVal = ""

	if r.Managed == nil {
		return defaultVal
	}
	iprefix, hasPrefix := r.Managed["go_package_prefix"]
	if !hasPrefix {
		return defaultVal
	}
	prefix, ok := iprefix.(ManagedRules)
	if !ok {
		return defaultVal
	}
	imodule, hasDefault := prefix["default"]
	if !hasDefault {
		return defaultVal
	}
	module, ok := imodule.(string)
	if !ok {
		return defaultVal
	}
	return module
}

type BufWorkYaml struct {
	Version     string   `yaml:"version"`
	Directories []string `yaml:"directories"`
}

type LintRuleSet struct {
	RulesUsed     []string `yaml:"use"`
	RulesExcluded []string `yaml:"except"`
	IgnoredPaths  []string `yaml:"ignore"`
}

type BreakingRuleSet struct {
	RulesUsed []string `yaml:"use"`
}

type BufYaml struct {
	Version       string          `yaml:"version"`
	LintRules     LintRuleSet     `yaml:"lint"`
	BreakingRules BreakingRuleSet `yaml:"breaking"`
}
