package log

var formatterCreators = map[string]CreateFormatterFunc{}

// CreateFormatterFunc is a function which creates a new instance
// of a Formatter.
type CreateFormatterFunc func(name, kind string) (Formatter, error)

// createFormatter creates formatters. It accepts a kind in {"text", "JSON"}
// which correspond to TextFormatter and JSONFormatter, and the name of the
// logger.
func createFormatter(name string, kind string) (Formatter, error) {
	if kind == FormatEnv {
		kind = logxiFormat
	}
	if kind == "" {
		kind = FormatText
	}

	fn := formatterCreators[kind]
	if fn == nil {
		fn = formatterCreators[FormatText]
	}

	formatter, err := fn(name, kind)
	if err != nil {
		return nil, err
	}
	// custom formatter may have not returned a formatter
	if formatter == nil {
		formatter, err = formatFactory(name, FormatText)
	}
	return formatter, err
}

func formatFactory(name string, kind string) (Formatter, error) {
	var formatter Formatter
	var err error
	switch kind {
	default:
		formatter = NewTextFormatter(name)
	case FormatHappy:
		formatter = NewHappyDevFormatter(name)
	case FormatText:
		formatter = NewTextFormatter(name)
	case FormatJSON:
		formatter = NewJSONFormatter(name)
	}
	return formatter, err
}

// RegisterFormatFactory registers a format factory function.
func RegisterFormatFactory(kind string, fn CreateFormatterFunc) {
	if kind == "" {
		panic("kind is empty string")
	}
	if fn == nil {
		panic("creator is nil")
	}
	formatterCreators[kind] = fn
}
