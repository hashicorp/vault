package log

import (
	"os"
	"strconv"
	"strings"
)

var contextLines int

// Configuration comes from environment or external services like
// consul, etcd.
type Configuration struct {
	Format string `json:"format"`
	Colors string `json:"colors"`
	Levels string `json:"levels"`
}

func readFromEnviron() *Configuration {
	conf := &Configuration{}

	var envOrDefault = func(name, val string) string {
		result := os.Getenv(name)
		if result == "" {
			result = val
		}
		return result
	}

	conf.Levels = envOrDefault("LOGXI", defaultLogxiEnv)
	conf.Format = envOrDefault("LOGXI_FORMAT", defaultLogxiFormatEnv)
	conf.Colors = envOrDefault("LOGXI_COLORS", defaultLogxiColorsEnv)
	return conf
}

// ProcessEnv (re)processes environment.
func ProcessEnv(env *Configuration) {
	// TODO: allow reading from etcd

	ProcessLogxiEnv(env.Levels)
	ProcessLogxiColorsEnv(env.Colors)
	ProcessLogxiFormatEnv(env.Format)
}

// ProcessLogxiFormatEnv parses LOGXI_FORMAT
func ProcessLogxiFormatEnv(env string) {
	logxiFormat = env
	m := parseKVList(logxiFormat, ",")
	formatterFormat := ""
	tFormat := ""
	for key, value := range m {
		switch key {
		default:
			formatterFormat = key
		case "t":
			tFormat = value
		case "pretty":
			isPretty = value != "false" && value != "0"
		case "maxcol":
			col, err := strconv.Atoi(value)
			if err == nil {
				maxCol = col
			} else {
				maxCol = defaultMaxCol
			}
		case "context":
			lines, err := strconv.Atoi(value)
			if err == nil {
				contextLines = lines
			} else {
				contextLines = defaultContextLines
			}
		case "LTSV":
			formatterFormat = "text"
			AssignmentChar = ltsvAssignmentChar
			Separator = ltsvSeparator
		}
	}
	if formatterFormat == "" || formatterCreators[formatterFormat] == nil {
		formatterFormat = defaultFormat
	}
	logxiFormat = formatterFormat
	if tFormat == "" {
		tFormat = defaultTimeFormat
	}
	timeFormat = tFormat
}

// ProcessLogxiEnv parses LOGXI variable
func ProcessLogxiEnv(env string) {
	logxiEnable := env
	if logxiEnable == "" {
		logxiEnable = defaultLogxiEnv
	}

	logxiNameLevelMap = map[string]int{}
	m := parseKVList(logxiEnable, ",")
	if m == nil {
		logxiNameLevelMap["*"] = defaultLevel
	}
	for key, value := range m {
		if strings.HasPrefix(key, "-") {
			// LOGXI=*,-foo => disable foo
			logxiNameLevelMap[key[1:]] = LevelOff
		} else if value == "" {
			// LOGXI=* => default to all
			logxiNameLevelMap[key] = LevelAll
		} else {
			// LOGXI=*=ERR => use user-specified level
			level := LevelAtoi[value]
			if level == 0 {
				InternalLog.Error("Unknown level in LOGXI environment variable", "key", key, "value", value, "LOGXI", env)
				level = defaultLevel
			}
			logxiNameLevelMap[key] = level
		}
	}

	// must always have global default, otherwise errs may get eaten up
	if _, ok := logxiNameLevelMap["*"]; !ok {
		logxiNameLevelMap["*"] = LevelError
	}
}

func getLogLevel(name string) int {
	var wildcardLevel int
	var result int

	for k, v := range logxiNameLevelMap {
		if k == name {
			result = v
		} else if k == "*" {
			wildcardLevel = v
		} else if strings.HasPrefix(k, "*") && strings.HasSuffix(name, k[1:]) {
			result = v
		} else if strings.HasSuffix(k, "*") && strings.HasPrefix(name, k[:len(k)-1]) {
			result = v
		}
	}

	if result == LevelOff {
		return LevelOff
	}

	if result > 0 {
		return result
	}

	if wildcardLevel > 0 {
		return wildcardLevel
	}

	return LevelOff
}

// ProcessLogxiColorsEnv parases LOGXI_COLORS
func ProcessLogxiColorsEnv(env string) {
	colors := env
	if colors == "" {
		colors = defaultLogxiColorsEnv
	} else if colors == "*=off" {
		// disable all colors
		disableColors = true
	}
	theme = parseTheme(colors)
}
