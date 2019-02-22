package lintersdb

import (
	"os"

	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

type Manager struct {
	nameToLC map[string]*linter.Config
}

func NewManager() *Manager {
	m := &Manager{}
	nameToLC := make(map[string]*linter.Config)
	for _, lc := range m.GetAllSupportedLinterConfigs() {
		for _, name := range lc.AllNames() {
			nameToLC[name] = lc
		}
	}

	m.nameToLC = nameToLC
	return m
}

func (Manager) AllPresets() []string {
	return []string{linter.PresetBugs, linter.PresetUnused, linter.PresetFormatting,
		linter.PresetStyle, linter.PresetComplexity, linter.PresetPerformance}
}

func (m Manager) allPresetsSet() map[string]bool {
	ret := map[string]bool{}
	for _, p := range m.AllPresets() {
		ret[p] = true
	}
	return ret
}

func (m Manager) GetMetaLinter(name string) linter.MetaLinter {
	return m.GetMetaLinters()[name]
}

func (m Manager) GetLinterConfig(name string) *linter.Config {
	lc, ok := m.nameToLC[name]
	if !ok {
		return nil
	}

	return lc
}

func enableLinterConfigs(lcs []*linter.Config, isEnabled func(lc *linter.Config) bool) []*linter.Config {
	var ret []*linter.Config
	for _, lc := range lcs {
		lc := lc
		lc.EnabledByDefault = isEnabled(lc)
		ret = append(ret, lc)
	}

	return ret
}

func (Manager) GetMetaLinters() map[string]linter.MetaLinter {
	metaLinters := []linter.MetaLinter{
		golinters.MegacheckMetalinter{},
	}

	ret := map[string]linter.MetaLinter{}
	for _, metaLinter := range metaLinters {
		ret[metaLinter.Name()] = metaLinter
	}

	return ret
}

func (Manager) GetAllSupportedLinterConfigs() []*linter.Config {
	lcs := []*linter.Config{
		linter.NewConfig(golinters.Govet{}).
			WithTypeInfo().
			WithPresets(linter.PresetBugs).
			WithSpeed(4).
			WithAlternativeNames("vet", "vetshadow").
			WithURL("https://golang.org/cmd/vet/"),
		linter.NewConfig(golinters.Errcheck{}).
			WithTypeInfo().
			WithPresets(linter.PresetBugs).
			WithSpeed(10).
			WithURL("https://github.com/kisielk/errcheck"),
		linter.NewConfig(golinters.Golint{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(3).
			WithURL("https://github.com/golang/lint"),

		linter.NewConfig(golinters.NewStaticcheck()).
			WithSSA().
			WithPresets(linter.PresetBugs).
			WithSpeed(2).
			WithURL("https://staticcheck.io/"),
		linter.NewConfig(golinters.NewUnused()).
			WithSSA().
			WithPresets(linter.PresetUnused).
			WithSpeed(5).
			WithURL("https://github.com/dominikh/go-tools/tree/master/cmd/unused"),
		linter.NewConfig(golinters.NewGosimple()).
			WithSSA().
			WithPresets(linter.PresetStyle).
			WithSpeed(5).
			WithURL("https://github.com/dominikh/go-tools/tree/master/cmd/gosimple"),
		linter.NewConfig(golinters.NewStylecheck()).
			WithSSA().
			WithPresets(linter.PresetStyle).
			WithSpeed(5).
			WithURL("https://github.com/dominikh/go-tools/tree/master/stylecheck"),

		linter.NewConfig(golinters.Gosec{}).
			WithTypeInfo().
			WithPresets(linter.PresetBugs).
			WithSpeed(8).
			WithURL("https://github.com/securego/gosec").
			WithAlternativeNames("gas"),
		linter.NewConfig(golinters.Structcheck{}).
			WithTypeInfo().
			WithPresets(linter.PresetUnused).
			WithSpeed(10).
			WithURL("https://github.com/opennota/check"),
		linter.NewConfig(golinters.Varcheck{}).
			WithTypeInfo().
			WithPresets(linter.PresetUnused).
			WithSpeed(10).
			WithURL("https://github.com/opennota/check"),
		linter.NewConfig(golinters.Interfacer{}).
			WithSSA().
			WithPresets(linter.PresetStyle).
			WithSpeed(6).
			WithURL("https://github.com/mvdan/interfacer"),
		linter.NewConfig(golinters.Unconvert{}).
			WithTypeInfo().
			WithPresets(linter.PresetStyle).
			WithSpeed(10).
			WithURL("https://github.com/mdempsky/unconvert"),
		linter.NewConfig(golinters.Ineffassign{}).
			WithPresets(linter.PresetUnused).
			WithSpeed(9).
			WithURL("https://github.com/gordonklaus/ineffassign"),
		linter.NewConfig(golinters.Dupl{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(7).
			WithURL("https://github.com/mibk/dupl"),
		linter.NewConfig(golinters.Goconst{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(9).
			WithURL("https://github.com/jgautheron/goconst"),
		linter.NewConfig(golinters.Deadcode{}).
			WithTypeInfo().
			WithPresets(linter.PresetUnused).
			WithSpeed(10).
			WithURL("https://github.com/remyoudompheng/go-misc/tree/master/deadcode"),
		linter.NewConfig(golinters.Gocyclo{}).
			WithPresets(linter.PresetComplexity).
			WithSpeed(8).
			WithURL("https://github.com/alecthomas/gocyclo"),
		linter.NewConfig(golinters.TypeCheck{}).
			WithTypeInfo().
			WithPresets(linter.PresetBugs).
			WithSpeed(10).
			WithURL(""),

		linter.NewConfig(golinters.Gofmt{}).
			WithPresets(linter.PresetFormatting).
			WithSpeed(7).
			WithURL("https://golang.org/cmd/gofmt/"),
		linter.NewConfig(golinters.Gofmt{UseGoimports: true}).
			WithPresets(linter.PresetFormatting).
			WithSpeed(5).
			WithURL("https://godoc.org/golang.org/x/tools/cmd/goimports"),
		linter.NewConfig(golinters.Maligned{}).
			WithTypeInfo().
			WithPresets(linter.PresetPerformance).
			WithSpeed(10).
			WithURL("https://github.com/mdempsky/maligned"),
		linter.NewConfig(golinters.Depguard{}).
			WithTypeInfo().
			WithPresets(linter.PresetStyle).
			WithSpeed(6).
			WithURL("https://github.com/OpenPeeDeeP/depguard"),
		linter.NewConfig(golinters.Misspell{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(7).
			WithURL("https://github.com/client9/misspell"),
		linter.NewConfig(golinters.Lll{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(10).
			WithURL("https://github.com/walle/lll"),
		linter.NewConfig(golinters.Unparam{}).
			WithPresets(linter.PresetUnused).
			WithSpeed(3).
			WithSSA().
			WithURL("https://github.com/mvdan/unparam"),
		linter.NewConfig(golinters.Nakedret{}).
			WithPresets(linter.PresetComplexity).
			WithSpeed(10).
			WithURL("https://github.com/alexkohler/nakedret"),
		linter.NewConfig(golinters.Prealloc{}).
			WithPresets(linter.PresetPerformance).
			WithSpeed(8).
			WithURL("https://github.com/alexkohler/prealloc"),
		linter.NewConfig(golinters.Scopelint{}).
			WithPresets(linter.PresetBugs).
			WithSpeed(8).
			WithURL("https://github.com/kyoh86/scopelint"),
		linter.NewConfig(golinters.Gocritic{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(5).
			WithTypeInfo().
			WithURL("https://github.com/go-critic/go-critic"),
		linter.NewConfig(golinters.Gochecknoinits{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(10).
			WithURL("https://github.com/leighmcculloch/gochecknoinits"),
		linter.NewConfig(golinters.Gochecknoglobals{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(10).
			WithURL("https://github.com/leighmcculloch/gochecknoglobals"),
	}

	isLocalRun := os.Getenv("GOLANGCI_COM_RUN") == ""
	enabledByDefault := map[string]bool{
		golinters.Govet{}.Name():       true,
		golinters.Errcheck{}.Name():    true,
		golinters.Staticcheck{}.Name(): true,
		golinters.Unused{}.Name():      true,
		golinters.Gosimple{}.Name():    true,
		golinters.Structcheck{}.Name(): true,
		golinters.Varcheck{}.Name():    true,
		golinters.Ineffassign{}.Name(): true,
		golinters.Deadcode{}.Name():    true,

		// don't typecheck for golangci.com: too many troubles
		golinters.TypeCheck{}.Name(): isLocalRun,
	}
	return enableLinterConfigs(lcs, func(lc *linter.Config) bool {
		return enabledByDefault[lc.Name()]
	})
}

func (m Manager) GetAllEnabledByDefaultLinters() []*linter.Config {
	var ret []*linter.Config
	for _, lc := range m.GetAllSupportedLinterConfigs() {
		if lc.EnabledByDefault {
			ret = append(ret, lc)
		}
	}

	return ret
}

func linterConfigsToMap(lcs []*linter.Config) map[string]*linter.Config {
	ret := map[string]*linter.Config{}
	for _, lc := range lcs {
		lc := lc // local copy
		ret[lc.Name()] = lc
	}

	return ret
}

func (m Manager) GetAllLinterConfigsForPreset(p string) []*linter.Config {
	var ret []*linter.Config
	for _, lc := range m.GetAllSupportedLinterConfigs() {
		for _, ip := range lc.InPresets {
			if p == ip {
				ret = append(ret, lc)
				break
			}
		}
	}

	return ret
}
