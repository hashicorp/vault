// FIXME(thaJeztah): remove once we are a module; the go:build directive prevents go from downgrading language version to go1.16:
//go:build go1.21

package loader

import (
	"reflect"
	"sort"

	"dario.cat/mergo"
	"github.com/docker/cli/cli/compose/types"
	"github.com/pkg/errors"
)

type specials struct {
	m map[reflect.Type]func(dst, src reflect.Value) error
}

func (s *specials) Transformer(t reflect.Type) func(dst, src reflect.Value) error {
	if fn, ok := s.m[t]; ok {
		return fn
	}
	return nil
}

func merge(configs []*types.Config) (*types.Config, error) {
	base := configs[0]
	for _, override := range configs[1:] {
		var err error
		base.Services, err = mergeServices(base.Services, override.Services)
		if err != nil {
			return base, errors.Wrapf(err, "cannot merge services from %s", override.Filename)
		}
		base.Volumes, err = mergeVolumes(base.Volumes, override.Volumes)
		if err != nil {
			return base, errors.Wrapf(err, "cannot merge volumes from %s", override.Filename)
		}
		base.Networks, err = mergeNetworks(base.Networks, override.Networks)
		if err != nil {
			return base, errors.Wrapf(err, "cannot merge networks from %s", override.Filename)
		}
		base.Secrets, err = mergeSecrets(base.Secrets, override.Secrets)
		if err != nil {
			return base, errors.Wrapf(err, "cannot merge secrets from %s", override.Filename)
		}
		base.Configs, err = mergeConfigs(base.Configs, override.Configs)
		if err != nil {
			return base, errors.Wrapf(err, "cannot merge configs from %s", override.Filename)
		}
	}
	return base, nil
}

func mergeServices(base, override []types.ServiceConfig) ([]types.ServiceConfig, error) {
	baseServices := mapByName(base)
	overrideServices := mapByName(override)
	specials := &specials{
		m: map[reflect.Type]func(dst, src reflect.Value) error{
			reflect.TypeOf(&types.LoggingConfig{}):           safelyMerge(mergeLoggingConfig),
			reflect.TypeOf([]types.ServicePortConfig{}):      mergeSlice(toServicePortConfigsMap, toServicePortConfigsSlice),
			reflect.TypeOf([]types.ServiceSecretConfig{}):    mergeSlice(toServiceSecretConfigsMap, toServiceSecretConfigsSlice),
			reflect.TypeOf([]types.ServiceConfigObjConfig{}): mergeSlice(toServiceConfigObjConfigsMap, toSServiceConfigObjConfigsSlice),
			reflect.TypeOf(&types.UlimitsConfig{}):           mergeUlimitsConfig,
			reflect.TypeOf([]types.ServiceVolumeConfig{}):    mergeSlice(toServiceVolumeConfigsMap, toServiceVolumeConfigsSlice),
			reflect.TypeOf(types.ShellCommand{}):             mergeShellCommand,
			reflect.TypeOf(&types.ServiceNetworkConfig{}):    mergeServiceNetworkConfig,
			reflect.PointerTo(reflect.TypeOf(uint64(1))):     mergeUint64,
		},
	}
	for name, overrideService := range overrideServices {
		overrideService := overrideService
		if baseService, ok := baseServices[name]; ok {
			if err := mergo.Merge(&baseService, &overrideService, mergo.WithAppendSlice, mergo.WithOverride, mergo.WithTransformers(specials)); err != nil {
				return base, errors.Wrapf(err, "cannot merge service %s", name)
			}
			baseServices[name] = baseService
			continue
		}
		baseServices[name] = overrideService
	}
	services := []types.ServiceConfig{}
	for _, baseService := range baseServices {
		services = append(services, baseService)
	}
	sort.Slice(services, func(i, j int) bool { return services[i].Name < services[j].Name })
	return services, nil
}

func toServiceSecretConfigsMap(s any) (map[any]any, error) {
	secrets, ok := s.([]types.ServiceSecretConfig)
	if !ok {
		return nil, errors.Errorf("not a serviceSecretConfig: %v", s)
	}
	m := map[any]any{}
	for _, secret := range secrets {
		m[secret.Source] = secret
	}
	return m, nil
}

func toServiceConfigObjConfigsMap(s any) (map[any]any, error) {
	secrets, ok := s.([]types.ServiceConfigObjConfig)
	if !ok {
		return nil, errors.Errorf("not a serviceSecretConfig: %v", s)
	}
	m := map[any]any{}
	for _, secret := range secrets {
		m[secret.Source] = secret
	}
	return m, nil
}

func toServicePortConfigsMap(s any) (map[any]any, error) {
	ports, ok := s.([]types.ServicePortConfig)
	if !ok {
		return nil, errors.Errorf("not a servicePortConfig slice: %v", s)
	}
	m := map[any]any{}
	for _, p := range ports {
		m[p.Published] = p
	}
	return m, nil
}

func toServiceVolumeConfigsMap(s any) (map[any]any, error) {
	volumes, ok := s.([]types.ServiceVolumeConfig)
	if !ok {
		return nil, errors.Errorf("not a serviceVolumeConfig slice: %v", s)
	}
	m := map[any]any{}
	for _, v := range volumes {
		m[v.Target] = v
	}
	return m, nil
}

func toServiceSecretConfigsSlice(dst reflect.Value, m map[any]any) error {
	s := []types.ServiceSecretConfig{}
	for _, v := range m {
		s = append(s, v.(types.ServiceSecretConfig))
	}
	sort.Slice(s, func(i, j int) bool { return s[i].Source < s[j].Source })
	dst.Set(reflect.ValueOf(s))
	return nil
}

func toSServiceConfigObjConfigsSlice(dst reflect.Value, m map[any]any) error {
	s := []types.ServiceConfigObjConfig{}
	for _, v := range m {
		s = append(s, v.(types.ServiceConfigObjConfig))
	}
	sort.Slice(s, func(i, j int) bool { return s[i].Source < s[j].Source })
	dst.Set(reflect.ValueOf(s))
	return nil
}

func toServicePortConfigsSlice(dst reflect.Value, m map[any]any) error {
	s := []types.ServicePortConfig{}
	for _, v := range m {
		s = append(s, v.(types.ServicePortConfig))
	}
	sort.Slice(s, func(i, j int) bool { return s[i].Published < s[j].Published })
	dst.Set(reflect.ValueOf(s))
	return nil
}

func toServiceVolumeConfigsSlice(dst reflect.Value, m map[any]any) error {
	s := []types.ServiceVolumeConfig{}
	for _, v := range m {
		s = append(s, v.(types.ServiceVolumeConfig))
	}
	sort.Slice(s, func(i, j int) bool { return s[i].Target < s[j].Target })
	dst.Set(reflect.ValueOf(s))
	return nil
}

type (
	tomapFn             func(s any) (map[any]any, error)
	writeValueFromMapFn func(reflect.Value, map[any]any) error
)

func safelyMerge(mergeFn func(dst, src reflect.Value) error) func(dst, src reflect.Value) error {
	return func(dst, src reflect.Value) error {
		if src.IsNil() {
			return nil
		}
		if dst.IsNil() {
			dst.Set(src)
			return nil
		}
		return mergeFn(dst, src)
	}
}

func mergeSlice(tomap tomapFn, writeValue writeValueFromMapFn) func(dst, src reflect.Value) error {
	return func(dst, src reflect.Value) error {
		dstMap, err := sliceToMap(tomap, dst)
		if err != nil {
			return err
		}
		srcMap, err := sliceToMap(tomap, src)
		if err != nil {
			return err
		}
		if err := mergo.Map(&dstMap, srcMap, mergo.WithOverride); err != nil {
			return err
		}
		return writeValue(dst, dstMap)
	}
}

func sliceToMap(tomap tomapFn, v reflect.Value) (map[any]any, error) {
	// check if valid
	if !v.IsValid() {
		return nil, errors.Errorf("invalid value : %+v", v)
	}
	return tomap(v.Interface())
}

func mergeLoggingConfig(dst, src reflect.Value) error {
	// Same driver, merging options
	if getLoggingDriver(dst.Elem()) == getLoggingDriver(src.Elem()) ||
		getLoggingDriver(dst.Elem()) == "" || getLoggingDriver(src.Elem()) == "" {
		if getLoggingDriver(dst.Elem()) == "" {
			dst.Elem().FieldByName("Driver").SetString(getLoggingDriver(src.Elem()))
		}
		dstOptions := dst.Elem().FieldByName("Options").Interface().(map[string]string)
		srcOptions := src.Elem().FieldByName("Options").Interface().(map[string]string)
		return mergo.Merge(&dstOptions, srcOptions, mergo.WithOverride)
	}
	// Different driver, override with src
	dst.Set(src)
	return nil
}

//nolint:unparam
func mergeUlimitsConfig(dst, src reflect.Value) error {
	if src.Interface() != reflect.Zero(reflect.TypeOf(src.Interface())).Interface() {
		dst.Elem().Set(src.Elem())
	}
	return nil
}

//nolint:unparam
func mergeShellCommand(dst, src reflect.Value) error {
	if src.Len() != 0 {
		dst.Set(src)
	}
	return nil
}

//nolint:unparam
func mergeServiceNetworkConfig(dst, src reflect.Value) error {
	if src.Interface() != reflect.Zero(reflect.TypeOf(src.Interface())).Interface() {
		dst.Elem().FieldByName("Aliases").Set(src.Elem().FieldByName("Aliases"))
		if ipv4 := src.Elem().FieldByName("Ipv4Address").Interface().(string); ipv4 != "" {
			dst.Elem().FieldByName("Ipv4Address").SetString(ipv4)
		}
		if ipv6 := src.Elem().FieldByName("Ipv6Address").Interface().(string); ipv6 != "" {
			dst.Elem().FieldByName("Ipv6Address").SetString(ipv6)
		}
	}
	return nil
}

//nolint:unparam
func mergeUint64(dst, src reflect.Value) error {
	if !src.IsNil() {
		dst.Elem().Set(src.Elem())
	}
	return nil
}

func getLoggingDriver(v reflect.Value) string {
	return v.FieldByName("Driver").String()
}

func mapByName(services []types.ServiceConfig) map[string]types.ServiceConfig {
	m := map[string]types.ServiceConfig{}
	for _, service := range services {
		m[service.Name] = service
	}
	return m
}

func mergeVolumes(base, override map[string]types.VolumeConfig) (map[string]types.VolumeConfig, error) {
	err := mergo.Map(&base, &override, mergo.WithOverride)
	return base, err
}

func mergeNetworks(base, override map[string]types.NetworkConfig) (map[string]types.NetworkConfig, error) {
	err := mergo.Map(&base, &override, mergo.WithOverride)
	return base, err
}

func mergeSecrets(base, override map[string]types.SecretConfig) (map[string]types.SecretConfig, error) {
	err := mergo.Map(&base, &override, mergo.WithOverride)
	return base, err
}

func mergeConfigs(base, override map[string]types.ConfigObjConfig) (map[string]types.ConfigObjConfig, error) {
	err := mergo.Map(&base, &override, mergo.WithOverride)
	return base, err
}
