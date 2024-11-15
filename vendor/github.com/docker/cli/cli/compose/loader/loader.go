// FIXME(thaJeztah): remove once we are a module; the go:build directive prevents go from downgrading language version to go1.16:
//go:build go1.21

package loader

import (
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	interp "github.com/docker/cli/cli/compose/interpolation"
	"github.com/docker/cli/cli/compose/schema"
	"github.com/docker/cli/cli/compose/template"
	"github.com/docker/cli/cli/compose/types"
	"github.com/docker/cli/opts"
	"github.com/docker/docker/api/types/versions"
	"github.com/docker/go-connections/nat"
	units "github.com/docker/go-units"
	"github.com/go-viper/mapstructure/v2"
	"github.com/google/shlex"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Options supported by Load
type Options struct {
	// Skip schema validation
	SkipValidation bool
	// Skip interpolation
	SkipInterpolation bool
	// Interpolation options
	Interpolate *interp.Options
	// Discard 'env_file' entries after resolving to 'environment' section
	discardEnvFiles bool
}

// WithDiscardEnvFiles sets the Options to discard the `env_file` section after resolving to
// the `environment` section
func WithDiscardEnvFiles(options *Options) {
	options.discardEnvFiles = true
}

// ParseYAML reads the bytes from a file, parses the bytes into a mapping
// structure, and returns it.
func ParseYAML(source []byte) (map[string]any, error) {
	var cfg any
	if err := yaml.Unmarshal(source, &cfg); err != nil {
		return nil, err
	}
	cfgMap, ok := cfg.(map[any]any)
	if !ok {
		return nil, errors.Errorf("top-level object must be a mapping")
	}
	converted, err := convertToStringKeysRecursive(cfgMap, "")
	if err != nil {
		return nil, err
	}
	return converted.(map[string]any), nil
}

// Load reads a ConfigDetails and returns a fully loaded configuration
func Load(configDetails types.ConfigDetails, opt ...func(*Options)) (*types.Config, error) {
	if len(configDetails.ConfigFiles) < 1 {
		return nil, errors.Errorf("No files specified")
	}

	options := &Options{
		Interpolate: &interp.Options{
			Substitute:      template.Substitute,
			LookupValue:     configDetails.LookupEnv,
			TypeCastMapping: interpolateTypeCastMapping,
		},
	}

	for _, op := range opt {
		op(options)
	}

	configs := []*types.Config{}
	var err error

	for _, file := range configDetails.ConfigFiles {
		configDict := file.Config
		version := schema.Version(configDict)
		if configDetails.Version == "" {
			configDetails.Version = version
		}
		if configDetails.Version != version {
			return nil, errors.Errorf("version mismatched between two composefiles : %v and %v", configDetails.Version, version)
		}

		if err := validateForbidden(configDict); err != nil {
			return nil, err
		}

		if !options.SkipInterpolation {
			configDict, err = interpolateConfig(configDict, *options.Interpolate)
			if err != nil {
				return nil, err
			}
		}

		if !options.SkipValidation {
			if err := schema.Validate(configDict, configDetails.Version); err != nil {
				return nil, err
			}
		}

		cfg, err := loadSections(configDict, configDetails)
		if err != nil {
			return nil, err
		}
		cfg.Filename = file.Filename
		if options.discardEnvFiles {
			for i := range cfg.Services {
				cfg.Services[i].EnvFile = nil
			}
		}

		configs = append(configs, cfg)
	}

	return merge(configs)
}

func validateForbidden(configDict map[string]any) error {
	servicesDict, ok := configDict["services"].(map[string]any)
	if !ok {
		return nil
	}
	forbidden := getProperties(servicesDict, types.ForbiddenProperties)
	if len(forbidden) > 0 {
		return &ForbiddenPropertiesError{Properties: forbidden}
	}
	return nil
}

func loadSections(config map[string]any, configDetails types.ConfigDetails) (*types.Config, error) {
	var err error
	cfg := types.Config{
		Version: schema.Version(config),
	}

	loaders := []struct {
		key string
		fnc func(config map[string]any) error
	}{
		{
			key: "services",
			fnc: func(config map[string]any) error {
				cfg.Services, err = LoadServices(config, configDetails.WorkingDir, configDetails.LookupEnv)
				return err
			},
		},
		{
			key: "networks",
			fnc: func(config map[string]any) error {
				cfg.Networks, err = LoadNetworks(config, configDetails.Version)
				return err
			},
		},
		{
			key: "volumes",
			fnc: func(config map[string]any) error {
				cfg.Volumes, err = LoadVolumes(config, configDetails.Version)
				return err
			},
		},
		{
			key: "secrets",
			fnc: func(config map[string]any) error {
				cfg.Secrets, err = LoadSecrets(config, configDetails)
				return err
			},
		},
		{
			key: "configs",
			fnc: func(config map[string]any) error {
				cfg.Configs, err = LoadConfigObjs(config, configDetails)
				return err
			},
		},
	}
	for _, loader := range loaders {
		if err := loader.fnc(getSection(config, loader.key)); err != nil {
			return nil, err
		}
	}
	cfg.Extras = getExtras(config)
	return &cfg, nil
}

func getSection(config map[string]any, key string) map[string]any {
	section, ok := config[key]
	if !ok {
		return make(map[string]any)
	}
	return section.(map[string]any)
}

// GetUnsupportedProperties returns the list of any unsupported properties that are
// used in the Compose files.
func GetUnsupportedProperties(configDicts ...map[string]any) []string {
	unsupported := map[string]bool{}

	for _, configDict := range configDicts {
		for _, service := range getServices(configDict) {
			serviceDict := service.(map[string]any)
			for _, property := range types.UnsupportedProperties {
				if _, isSet := serviceDict[property]; isSet {
					unsupported[property] = true
				}
			}
		}
	}

	return sortedKeys(unsupported)
}

func sortedKeys(set map[string]bool) []string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// GetDeprecatedProperties returns the list of any deprecated properties that
// are used in the compose files.
func GetDeprecatedProperties(configDicts ...map[string]any) map[string]string {
	deprecated := map[string]string{}

	for _, configDict := range configDicts {
		deprecatedProperties := getProperties(getServices(configDict), types.DeprecatedProperties)
		for key, value := range deprecatedProperties {
			deprecated[key] = value
		}
	}

	return deprecated
}

func getProperties(services map[string]any, propertyMap map[string]string) map[string]string {
	output := map[string]string{}

	for _, service := range services {
		if serviceDict, ok := service.(map[string]any); ok {
			for property, description := range propertyMap {
				if _, isSet := serviceDict[property]; isSet {
					output[property] = description
				}
			}
		}
	}

	return output
}

// ForbiddenPropertiesError is returned when there are properties in the Compose
// file that are forbidden.
type ForbiddenPropertiesError struct {
	Properties map[string]string
}

func (e *ForbiddenPropertiesError) Error() string {
	return "Configuration contains forbidden properties"
}

func getServices(configDict map[string]any) map[string]any {
	if services, ok := configDict["services"]; ok {
		if servicesDict, ok := services.(map[string]any); ok {
			return servicesDict
		}
	}

	return map[string]any{}
}

// Transform converts the source into the target struct with compose types transformer
// and the specified transformers if any.
func Transform(source any, target any, additionalTransformers ...Transformer) error {
	data := mapstructure.Metadata{}
	config := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			createTransformHook(additionalTransformers...),
			mapstructure.StringToTimeDurationHookFunc()),
		Result:   target,
		Metadata: &data,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(source)
}

// TransformerFunc defines a function to perform the actual transformation
type TransformerFunc func(any) (any, error)

// Transformer defines a map to type transformer
type Transformer struct {
	TypeOf reflect.Type
	Func   TransformerFunc
}

func createTransformHook(additionalTransformers ...Transformer) mapstructure.DecodeHookFuncType {
	transforms := map[reflect.Type]func(any) (any, error){
		reflect.TypeOf(types.External{}):                         transformExternal,
		reflect.TypeOf(types.HealthCheckTest{}):                  transformHealthCheckTest,
		reflect.TypeOf(types.ShellCommand{}):                     transformShellCommand,
		reflect.TypeOf(types.StringList{}):                       transformStringList,
		reflect.TypeOf(map[string]string{}):                      transformMapStringString,
		reflect.TypeOf(types.UlimitsConfig{}):                    transformUlimits,
		reflect.TypeOf(types.UnitBytes(0)):                       transformSize,
		reflect.TypeOf([]types.ServicePortConfig{}):              transformServicePort,
		reflect.TypeOf(types.ServiceSecretConfig{}):              transformStringSourceMap,
		reflect.TypeOf(types.ServiceConfigObjConfig{}):           transformStringSourceMap,
		reflect.TypeOf(types.StringOrNumberList{}):               transformStringOrNumberList,
		reflect.TypeOf(map[string]*types.ServiceNetworkConfig{}): transformServiceNetworkMap,
		reflect.TypeOf(types.Mapping{}):                          transformMappingOrListFunc("=", false),
		reflect.TypeOf(types.MappingWithEquals{}):                transformMappingOrListFunc("=", true),
		reflect.TypeOf(types.Labels{}):                           transformMappingOrListFunc("=", false),
		reflect.TypeOf(types.MappingWithColon{}):                 transformMappingOrListFunc(":", false),
		reflect.TypeOf(types.HostsList{}):                        transformHostsList,
		reflect.TypeOf(types.ServiceVolumeConfig{}):              transformServiceVolumeConfig,
		reflect.TypeOf(types.BuildConfig{}):                      transformBuildConfig,
		reflect.TypeOf(types.Duration(0)):                        transformStringToDuration,
	}

	for _, transformer := range additionalTransformers {
		transforms[transformer.TypeOf] = transformer.Func
	}

	return func(_ reflect.Type, target reflect.Type, data any) (any, error) {
		transform, ok := transforms[target]
		if !ok {
			return data, nil
		}
		return transform(data)
	}
}

// keys needs to be converted to strings for jsonschema
func convertToStringKeysRecursive(value any, keyPrefix string) (any, error) {
	if mapping, ok := value.(map[any]any); ok {
		dict := make(map[string]any)
		for key, entry := range mapping {
			str, ok := key.(string)
			if !ok {
				return nil, formatInvalidKeyError(keyPrefix, key)
			}
			var newKeyPrefix string
			if keyPrefix == "" {
				newKeyPrefix = str
			} else {
				newKeyPrefix = fmt.Sprintf("%s.%s", keyPrefix, str)
			}
			convertedEntry, err := convertToStringKeysRecursive(entry, newKeyPrefix)
			if err != nil {
				return nil, err
			}
			dict[str] = convertedEntry
		}
		return dict, nil
	}
	if list, ok := value.([]any); ok {
		var convertedList []any
		for index, entry := range list {
			newKeyPrefix := fmt.Sprintf("%s[%d]", keyPrefix, index)
			convertedEntry, err := convertToStringKeysRecursive(entry, newKeyPrefix)
			if err != nil {
				return nil, err
			}
			convertedList = append(convertedList, convertedEntry)
		}
		return convertedList, nil
	}
	return value, nil
}

func formatInvalidKeyError(keyPrefix string, key any) error {
	var location string
	if keyPrefix == "" {
		location = "at top level"
	} else {
		location = "in " + keyPrefix
	}
	return errors.Errorf("non-string key %s: %#v", location, key)
}

// LoadServices produces a ServiceConfig map from a compose file Dict
// the servicesDict is not validated if directly used. Use Load() to enable validation
func LoadServices(servicesDict map[string]any, workingDir string, lookupEnv template.Mapping) ([]types.ServiceConfig, error) {
	services := make([]types.ServiceConfig, 0, len(servicesDict))

	for name, serviceDef := range servicesDict {
		serviceConfig, err := LoadService(name, serviceDef.(map[string]any), workingDir, lookupEnv)
		if err != nil {
			return nil, err
		}
		services = append(services, *serviceConfig)
	}

	return services, nil
}

// LoadService produces a single ServiceConfig from a compose file Dict
// the serviceDict is not validated if directly used. Use Load() to enable validation
func LoadService(name string, serviceDict map[string]any, workingDir string, lookupEnv template.Mapping) (*types.ServiceConfig, error) {
	serviceConfig := &types.ServiceConfig{}
	if err := Transform(serviceDict, serviceConfig); err != nil {
		return nil, err
	}
	serviceConfig.Name = name

	if err := resolveEnvironment(serviceConfig, workingDir, lookupEnv); err != nil {
		return nil, err
	}

	if err := resolveVolumePaths(serviceConfig.Volumes, workingDir, lookupEnv); err != nil {
		return nil, err
	}

	serviceConfig.Extras = getExtras(serviceDict)

	return serviceConfig, nil
}

func loadExtras(name string, source map[string]any) map[string]any {
	if dict, ok := source[name].(map[string]any); ok {
		return getExtras(dict)
	}
	return nil
}

func getExtras(dict map[string]any) map[string]any {
	extras := map[string]any{}
	for key, value := range dict {
		if strings.HasPrefix(key, "x-") {
			extras[key] = value
		}
	}
	if len(extras) == 0 {
		return nil
	}
	return extras
}

func updateEnvironment(environment map[string]*string, vars map[string]*string, lookupEnv template.Mapping) {
	for k, v := range vars {
		interpolatedV, ok := lookupEnv(k)
		if (v == nil || *v == "") && ok {
			// lookupEnv is prioritized over vars
			environment[k] = &interpolatedV
		} else {
			environment[k] = v
		}
	}
}

func resolveEnvironment(serviceConfig *types.ServiceConfig, workingDir string, lookupEnv template.Mapping) error {
	environment := make(map[string]*string)

	if len(serviceConfig.EnvFile) > 0 {
		var envVars []string

		for _, file := range serviceConfig.EnvFile {
			filePath := absPath(workingDir, file)
			fileVars, err := opts.ParseEnvFile(filePath)
			if err != nil {
				return err
			}
			envVars = append(envVars, fileVars...)
		}
		updateEnvironment(environment,
			opts.ConvertKVStringsToMapWithNil(envVars), lookupEnv)
	}

	updateEnvironment(environment, serviceConfig.Environment, lookupEnv)
	serviceConfig.Environment = environment
	return nil
}

func resolveVolumePaths(volumes []types.ServiceVolumeConfig, workingDir string, lookupEnv template.Mapping) error {
	for i, volume := range volumes {
		if volume.Type != "bind" {
			continue
		}

		if volume.Source == "" {
			return errors.New(`invalid mount config for type "bind": field Source must not be empty`)
		}

		filePath := expandUser(volume.Source, lookupEnv)
		// Check if source is an absolute path (either Unix or Windows), to
		// handle a Windows client with a Unix daemon or vice-versa.
		//
		// Note that this is not required for Docker for Windows when specifying
		// a local Windows path, because Docker for Windows translates the Windows
		// path into a valid path within the VM.
		if !path.IsAbs(filePath) && !isAbs(filePath) {
			filePath = absPath(workingDir, filePath)
		}
		volume.Source = filePath
		volumes[i] = volume
	}
	return nil
}

// TODO: make this more robust
func expandUser(srcPath string, lookupEnv template.Mapping) string {
	if strings.HasPrefix(srcPath, "~") {
		home, ok := lookupEnv("HOME")
		if !ok {
			logrus.Warn("cannot expand '~', because the environment lacks HOME")
			return srcPath
		}
		return strings.Replace(srcPath, "~", home, 1)
	}
	return srcPath
}

func transformUlimits(data any) (any, error) {
	switch value := data.(type) {
	case int:
		return types.UlimitsConfig{Single: value}, nil
	case map[string]any:
		ulimit := types.UlimitsConfig{}
		ulimit.Soft = value["soft"].(int)
		ulimit.Hard = value["hard"].(int)
		return ulimit, nil
	default:
		return data, errors.Errorf("invalid type %T for ulimits", value)
	}
}

// LoadNetworks produces a NetworkConfig map from a compose file Dict
// the source Dict is not validated if directly used. Use Load() to enable validation
func LoadNetworks(source map[string]any, version string) (map[string]types.NetworkConfig, error) {
	networks := make(map[string]types.NetworkConfig)
	err := Transform(source, &networks)
	if err != nil {
		return networks, err
	}
	for name, network := range networks {
		if !network.External.External {
			continue
		}
		switch {
		case network.External.Name != "":
			if network.Name != "" {
				return nil, errors.Errorf("network %s: network.external.name and network.name conflict; only use network.name", name)
			}
			if versions.GreaterThanOrEqualTo(version, "3.5") {
				logrus.Warnf("network %s: network.external.name is deprecated in favor of network.name", name)
			}
			network.Name = network.External.Name
			network.External.Name = ""
		case network.Name == "":
			network.Name = name
		}
		network.Extras = loadExtras(name, source)
		networks[name] = network
	}
	return networks, nil
}

func externalVolumeError(volume, key string) error {
	return errors.Errorf(
		"conflicting parameters \"external\" and %q specified for volume %q",
		key, volume)
}

// LoadVolumes produces a VolumeConfig map from a compose file Dict
// the source Dict is not validated if directly used. Use Load() to enable validation
func LoadVolumes(source map[string]any, version string) (map[string]types.VolumeConfig, error) {
	volumes := make(map[string]types.VolumeConfig)
	if err := Transform(source, &volumes); err != nil {
		return volumes, err
	}

	for name, volume := range volumes {
		if !volume.External.External {
			continue
		}
		switch {
		case volume.Driver != "":
			return nil, externalVolumeError(name, "driver")
		case len(volume.DriverOpts) > 0:
			return nil, externalVolumeError(name, "driver_opts")
		case len(volume.Labels) > 0:
			return nil, externalVolumeError(name, "labels")
		case volume.External.Name != "":
			if volume.Name != "" {
				return nil, errors.Errorf("volume %s: volume.external.name and volume.name conflict; only use volume.name", name)
			}
			if versions.GreaterThanOrEqualTo(version, "3.4") {
				logrus.Warnf("volume %s: volume.external.name is deprecated in favor of volume.name", name)
			}
			volume.Name = volume.External.Name
			volume.External.Name = ""
		case volume.Name == "":
			volume.Name = name
		}
		volume.Extras = loadExtras(name, source)
		volumes[name] = volume
	}
	return volumes, nil
}

// LoadSecrets produces a SecretConfig map from a compose file Dict
// the source Dict is not validated if directly used. Use Load() to enable validation
func LoadSecrets(source map[string]any, details types.ConfigDetails) (map[string]types.SecretConfig, error) {
	secrets := make(map[string]types.SecretConfig)
	if err := Transform(source, &secrets); err != nil {
		return secrets, err
	}
	for name, secret := range secrets {
		obj, err := loadFileObjectConfig(name, "secret", types.FileObjectConfig(secret), details)
		if err != nil {
			return nil, err
		}
		secretConfig := types.SecretConfig(obj)
		secretConfig.Extras = loadExtras(name, source)
		secrets[name] = secretConfig
	}
	return secrets, nil
}

// LoadConfigObjs produces a ConfigObjConfig map from a compose file Dict
// the source Dict is not validated if directly used. Use Load() to enable validation
func LoadConfigObjs(source map[string]any, details types.ConfigDetails) (map[string]types.ConfigObjConfig, error) {
	configs := make(map[string]types.ConfigObjConfig)
	if err := Transform(source, &configs); err != nil {
		return configs, err
	}
	for name, config := range configs {
		obj, err := loadFileObjectConfig(name, "config", types.FileObjectConfig(config), details)
		if err != nil {
			return nil, err
		}
		configConfig := types.ConfigObjConfig(obj)
		configConfig.Extras = loadExtras(name, source)
		configs[name] = configConfig
	}
	return configs, nil
}

func loadFileObjectConfig(name string, objType string, obj types.FileObjectConfig, details types.ConfigDetails) (types.FileObjectConfig, error) {
	// if "external: true"
	switch {
	case obj.External.External:
		// handle deprecated external.name
		if obj.External.Name != "" {
			if obj.Name != "" {
				return obj, errors.Errorf("%[1]s %[2]s: %[1]s.external.name and %[1]s.name conflict; only use %[1]s.name", objType, name)
			}
			if versions.GreaterThanOrEqualTo(details.Version, "3.5") {
				logrus.Warnf("%[1]s %[2]s: %[1]s.external.name is deprecated in favor of %[1]s.name", objType, name)
			}
			obj.Name = obj.External.Name
			obj.External.Name = ""
		} else if obj.Name == "" {
			obj.Name = name
		}
		// if not "external: true"
	case obj.Driver != "":
		if obj.File != "" {
			return obj, errors.Errorf("%[1]s %[2]s: %[1]s.driver and %[1]s.file conflict; only use %[1]s.driver", objType, name)
		}
	default:
		obj.File = absPath(details.WorkingDir, obj.File)
	}

	return obj, nil
}

func absPath(workingDir string, filePath string) string {
	if filepath.IsAbs(filePath) {
		return filePath
	}
	return filepath.Join(workingDir, filePath)
}

var transformMapStringString TransformerFunc = func(data any) (any, error) {
	switch value := data.(type) {
	case map[string]any:
		return toMapStringString(value, false), nil
	case map[string]string:
		return value, nil
	default:
		return data, errors.Errorf("invalid type %T for map[string]string", value)
	}
}

var transformExternal TransformerFunc = func(data any) (any, error) {
	switch value := data.(type) {
	case bool:
		return map[string]any{"external": value}, nil
	case map[string]any:
		return map[string]any{"external": true, "name": value["name"]}, nil
	default:
		return data, errors.Errorf("invalid type %T for external", value)
	}
}

var transformServicePort TransformerFunc = func(data any) (any, error) {
	switch entries := data.(type) {
	case []any:
		// We process the list instead of individual items here.
		// The reason is that one entry might be mapped to multiple ServicePortConfig.
		// Therefore we take an input of a list and return an output of a list.
		ports := []any{}
		for _, entry := range entries {
			switch value := entry.(type) {
			case int:
				v, err := toServicePortConfigs(strconv.Itoa(value))
				if err != nil {
					return data, err
				}
				ports = append(ports, v...)
			case string:
				v, err := toServicePortConfigs(value)
				if err != nil {
					return data, err
				}
				ports = append(ports, v...)
			case map[string]any:
				ports = append(ports, value)
			default:
				return data, errors.Errorf("invalid type %T for port", value)
			}
		}
		return ports, nil
	default:
		return data, errors.Errorf("invalid type %T for port", entries)
	}
}

var transformStringSourceMap TransformerFunc = func(data any) (any, error) {
	switch value := data.(type) {
	case string:
		return map[string]any{"source": value}, nil
	case map[string]any:
		return data, nil
	default:
		return data, errors.Errorf("invalid type %T for secret", value)
	}
}

var transformBuildConfig TransformerFunc = func(data any) (any, error) {
	switch value := data.(type) {
	case string:
		return map[string]any{"context": value}, nil
	case map[string]any:
		return data, nil
	default:
		return data, errors.Errorf("invalid type %T for service build", value)
	}
}

var transformServiceVolumeConfig TransformerFunc = func(data any) (any, error) {
	switch value := data.(type) {
	case string:
		return ParseVolume(value)
	case map[string]any:
		return data, nil
	default:
		return data, errors.Errorf("invalid type %T for service volume", value)
	}
}

var transformServiceNetworkMap TransformerFunc = func(value any) (any, error) {
	if list, ok := value.([]any); ok {
		mapValue := map[any]any{}
		for _, name := range list {
			mapValue[name] = nil
		}
		return mapValue, nil
	}
	return value, nil
}

var transformStringOrNumberList TransformerFunc = func(value any) (any, error) {
	list := value.([]any)
	result := make([]string, len(list))
	for i, item := range list {
		result[i] = fmt.Sprint(item)
	}
	return result, nil
}

var transformStringList TransformerFunc = func(data any) (any, error) {
	switch value := data.(type) {
	case string:
		return []string{value}, nil
	case []any:
		return value, nil
	default:
		return data, errors.Errorf("invalid type %T for string list", value)
	}
}

var transformHostsList TransformerFunc = func(data any) (any, error) {
	hl := transformListOrMapping(data, ":", false, []string{"=", ":"})

	// Remove brackets from IP addresses if present (for example "[::1]" -> "::1").
	result := make([]string, 0, len(hl))
	for _, hip := range hl {
		host, ip, _ := strings.Cut(hip, ":")
		if len(ip) > 2 && ip[0] == '[' && ip[len(ip)-1] == ']' {
			ip = ip[1 : len(ip)-1]
		}
		result = append(result, fmt.Sprintf("%s:%s", host, ip))
	}
	return result, nil
}

// transformListOrMapping transforms pairs of strings that may be represented as
// a map, or a list of '=' or ':' separated strings, into a list of ':' separated
// strings.
func transformListOrMapping(listOrMapping any, sep string, allowNil bool, allowSeps []string) []string {
	switch value := listOrMapping.(type) {
	case map[string]any:
		return toStringList(value, sep, allowNil)
	case []any:
		result := make([]string, 0, len(value))
		for _, entry := range value {
			for i, allowSep := range allowSeps {
				entry := fmt.Sprint(entry)
				k, v, ok := strings.Cut(entry, allowSep)
				if ok {
					// Entry uses this allowed separator. Add it to the result, using
					// sep as a separator.
					result = append(result, fmt.Sprintf("%s%s%s", k, sep, v))
					break
				} else if i == len(allowSeps)-1 {
					// No more separators to try, keep the entry if allowNil.
					if allowNil {
						result = append(result, k)
					}
				}
			}
		}
		return result
	}
	panic(errors.Errorf("expected a map or a list, got %T: %#v", listOrMapping, listOrMapping))
}

func transformMappingOrListFunc(sep string, allowNil bool) TransformerFunc {
	return func(data any) (any, error) {
		return transformMappingOrList(data, sep, allowNil), nil
	}
}

func transformMappingOrList(mappingOrList any, sep string, allowNil bool) any {
	switch values := mappingOrList.(type) {
	case map[string]any:
		return toMapStringString(values, allowNil)
	case []any:
		result := make(map[string]any)
		for _, v := range values {
			key, val, hasValue := strings.Cut(v.(string), sep)
			switch {
			case !hasValue && allowNil:
				result[key] = nil
			case !hasValue && !allowNil:
				result[key] = ""
			default:
				result[key] = val
			}
		}
		return result
	}
	panic(errors.Errorf("expected a map or a list, got %T: %#v", mappingOrList, mappingOrList))
}

var transformShellCommand TransformerFunc = func(value any) (any, error) {
	if str, ok := value.(string); ok {
		return shlex.Split(str)
	}
	return value, nil
}

var transformHealthCheckTest TransformerFunc = func(data any) (any, error) {
	switch value := data.(type) {
	case string:
		return append([]string{"CMD-SHELL"}, value), nil
	case []any:
		return value, nil
	default:
		return value, errors.Errorf("invalid type %T for healthcheck.test", value)
	}
}

var transformSize TransformerFunc = func(value any) (any, error) {
	switch value := value.(type) {
	case int:
		return int64(value), nil
	case string:
		return units.RAMInBytes(value)
	}
	panic(errors.Errorf("invalid type for size %T", value))
}

var transformStringToDuration TransformerFunc = func(value any) (any, error) {
	switch value := value.(type) {
	case string:
		d, err := time.ParseDuration(value)
		if err != nil {
			return value, err
		}
		return types.Duration(d), nil
	default:
		return value, errors.Errorf("invalid type %T for duration", value)
	}
}

func toServicePortConfigs(value string) ([]any, error) {
	var portConfigs []any

	ports, portBindings, err := nat.ParsePortSpecs([]string{value})
	if err != nil {
		return nil, err
	}
	// We need to sort the key of the ports to make sure it is consistent
	keys := []string{}
	for port := range ports {
		keys = append(keys, string(port))
	}
	sort.Strings(keys)

	for _, key := range keys {
		// Reuse ConvertPortToPortConfig so that it is consistent
		portConfig, err := opts.ConvertPortToPortConfig(nat.Port(key), portBindings)
		if err != nil {
			return nil, err
		}
		for _, p := range portConfig {
			portConfigs = append(portConfigs, types.ServicePortConfig{
				Protocol:  string(p.Protocol),
				Target:    p.TargetPort,
				Published: p.PublishedPort,
				Mode:      string(p.PublishMode),
			})
		}
	}

	return portConfigs, nil
}

func toMapStringString(value map[string]any, allowNil bool) map[string]any {
	output := make(map[string]any)
	for key, value := range value {
		output[key] = toString(value, allowNil)
	}
	return output
}

func toString(value any, allowNil bool) any {
	switch {
	case value != nil:
		return fmt.Sprint(value)
	case allowNil:
		return nil
	default:
		return ""
	}
}

func toStringList(value map[string]any, separator string, allowNil bool) []string {
	output := []string{}
	for key, value := range value {
		if value == nil && !allowNil {
			continue
		}
		output = append(output, fmt.Sprintf("%s%s%s", key, separator, value))
	}
	sort.Strings(output)
	return output
}
