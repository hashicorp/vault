package linter

type MetaLinter interface {
	Name() string
	BuildLinterConfig(enabledChildren []string) (*Config, error)
	AllChildLinterNames() []string
	DefaultChildLinterNames() []string
}
