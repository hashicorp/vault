package lintersdb

import (
	"sort"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type EnabledSet struct {
	m   *Manager
	v   *Validator
	log logutils.Log
	cfg *config.Config
}

func NewEnabledSet(m *Manager, v *Validator, log logutils.Log, cfg *config.Config) *EnabledSet {
	return &EnabledSet{
		m:   m,
		v:   v,
		log: log,
		cfg: cfg,
	}
}

// nolint:gocyclo
func (es EnabledSet) build(lcfg *config.Linters, enabledByDefaultLinters []*linter.Config) map[string]*linter.Config {
	resultLintersSet := map[string]*linter.Config{}
	switch {
	case len(lcfg.Presets) != 0:
		break // imply --disable-all
	case lcfg.EnableAll:
		resultLintersSet = linterConfigsToMap(es.m.GetAllSupportedLinterConfigs())
	case lcfg.DisableAll:
		break
	default:
		resultLintersSet = linterConfigsToMap(enabledByDefaultLinters)
	}

	// --presets can only add linters to default set
	for _, p := range lcfg.Presets {
		for _, lc := range es.m.GetAllLinterConfigsForPreset(p) {
			lc := lc
			resultLintersSet[lc.Name()] = lc
		}
	}

	// --fast removes slow linters from current set.
	// It should be after --presets to be able to run only fast linters in preset.
	// It should be before --enable and --disable to be able to enable or disable specific linter.
	if lcfg.Fast {
		for name := range resultLintersSet {
			if es.m.GetLinterConfig(name).NeedsSSARepr {
				delete(resultLintersSet, name)
			}
		}
	}

	metaLinters := es.m.GetMetaLinters()

	for _, name := range lcfg.Enable {
		if metaLinter := metaLinters[name]; metaLinter != nil {
			// e.g. if we use --enable=megacheck we should add staticcheck,unused and gosimple to result set
			for _, childLinter := range metaLinter.DefaultChildLinterNames() {
				resultLintersSet[childLinter] = es.m.GetLinterConfig(childLinter)
			}
			continue
		}

		lc := es.m.GetLinterConfig(name)
		// it's important to use lc.Name() nor name because name can be alias
		resultLintersSet[lc.Name()] = lc
	}

	for _, name := range lcfg.Disable {
		if metaLinter := metaLinters[name]; metaLinter != nil {
			// e.g. if we use --disable=megacheck we should remove staticcheck,unused and gosimple from result set
			for _, childLinter := range metaLinter.DefaultChildLinterNames() {
				delete(resultLintersSet, childLinter)
			}
			continue
		}

		lc := es.m.GetLinterConfig(name)
		// it's important to use lc.Name() nor name because name can be alias
		delete(resultLintersSet, lc.Name())
	}

	return resultLintersSet
}

func (es EnabledSet) optimizeLintersSet(linters map[string]*linter.Config) {
	for _, metaLinter := range es.m.GetMetaLinters() {
		var children []string
		for _, child := range metaLinter.AllChildLinterNames() {
			if _, ok := linters[child]; ok {
				children = append(children, child)
			}
		}

		if len(children) <= 1 {
			continue
		}

		for _, child := range children {
			delete(linters, child)
		}
		builtLinterConfig, err := metaLinter.BuildLinterConfig(children)
		if err != nil {
			panic("shouldn't fail during linter building: " + err.Error())
		}
		linters[metaLinter.Name()] = builtLinterConfig
		es.log.Infof("Optimized sublinters %s into metalinter %s", children, metaLinter.Name())
	}
}

func (es EnabledSet) Get(optimize bool) ([]*linter.Config, error) {
	if err := es.v.validateEnabledDisabledLintersConfig(&es.cfg.Linters); err != nil {
		return nil, err
	}

	resultLintersSet := es.build(&es.cfg.Linters, es.m.GetAllEnabledByDefaultLinters())
	es.verbosePrintLintersStatus(resultLintersSet)
	if optimize {
		es.optimizeLintersSet(resultLintersSet)
	}

	var resultLinters []*linter.Config
	for _, lc := range resultLintersSet {
		resultLinters = append(resultLinters, lc)
	}

	return resultLinters, nil
}

func (es EnabledSet) verbosePrintLintersStatus(lcs map[string]*linter.Config) {
	var linterNames []string
	for _, lc := range lcs {
		linterNames = append(linterNames, lc.Name())
	}
	sort.StringSlice(linterNames).Sort()
	es.log.Infof("Active %d linters: %s", len(linterNames), linterNames)

	if len(es.cfg.Linters.Presets) != 0 {
		sort.StringSlice(es.cfg.Linters.Presets).Sort()
		es.log.Infof("Active presets: %s", es.cfg.Linters.Presets)
	}
}
