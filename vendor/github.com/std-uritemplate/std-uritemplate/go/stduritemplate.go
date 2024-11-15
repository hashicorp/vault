package stduritemplate

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Substitutions map[string]any

// Public API
func Expand(template string, substitutions Substitutions) (string, error) {
	return expandImpl(template, substitutions)
}

// Private implementation
type Op rune

const (
	OpUndefined    Op = 0
	OpNone         Op = -1
	OpPlus         Op = '+'
	OpHash         Op = '#'
	OpDot          Op = '.'
	OpSlash        Op = '/'
	OpSemicolon    Op = ';'
	OpQuestionMark Op = '?'
	OpAmp          Op = '&'
)

const (
	SubstitutionTypeEmpty  = "EMPTY"
	SubstitutionTypeString = "STRING"
	SubstitutionTypeList   = "LIST"
	SubstitutionTypeMap    = "MAP"
)

func validateLiteral(c rune, col int) error {
	switch c {
	case '+', '#', '/', ';', '?', '&', ' ', '!', '=', '$', '|', '*', ':', '~', '-':
		return fmt.Errorf("illegal character identified in the token at col: %d", col)
	default:
		return nil
	}
}

func getMaxChar(buffer *strings.Builder, col int) (int, error) {
	if buffer == nil || buffer.Len() == 0 {
		return -1, nil
	}
	value := buffer.String()

	if value == "" {
		return -1, nil
	}

	maxChar, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("cannot parse max chars at col: %d", col)
	}
	return maxChar, nil
}

func getOperator(c rune, token *strings.Builder, col int) (Op, error) {
	switch c {
	case '+':
		return OpPlus, nil
	case '#':
		return OpHash, nil
	case '.':
		return OpDot, nil
	case '/':
		return OpSlash, nil
	case ';':
		return OpSemicolon, nil
	case '?':
		return OpQuestionMark, nil
	case '&':
		return OpAmp, nil
	default:
		err := validateLiteral(c, col)
		if err != nil {
			return OpUndefined, err
		}
		token.WriteRune(c)
		return OpNone, nil
	}
}

func expandImpl(str string, substitutions Substitutions) (string, error) {
	var result strings.Builder

	var token = &strings.Builder{}
	var toToken = false
	var operator = OpUndefined
	var composite bool
	var maxCharBuffer = &strings.Builder{}
	var toMaxCharBuffer = false
	var firstToken = true

	for i, character := range str {
		switch character {
		case '{':
			toToken = true
			token.Reset()
			firstToken = true
		case '}':
			if toToken {
				maxChar, err := getMaxChar(maxCharBuffer, i)
				if err != nil {
					return "", err
				}
				expanded, err := expandToken(operator, token.String(), composite, maxChar, firstToken, substitutions, &result, i)
				if err != nil {
					return "", err
				}
				if expanded && firstToken {
					firstToken = false
				}
				toToken = false
				token.Reset()
				operator = OpUndefined
				composite = false
				toMaxCharBuffer = false
				maxCharBuffer.Reset()
			} else {
				return "", fmt.Errorf("failed to expand token, invalid at col: %d", i)
			}
		case ',':
			if toToken {
				maxChar, err := getMaxChar(maxCharBuffer, i)
				if err != nil {
					return "", err
				}
				expanded, err := expandToken(operator, token.String(), composite, maxChar, firstToken, substitutions, &result, i)
				if err != nil {
					return "", err
				}
				if expanded && firstToken {
					firstToken = false
				}
				token.Reset()
				composite = false
				toMaxCharBuffer = false
				maxCharBuffer.Reset()
				break
			}
			// Intentional fall-through for commas outside the {}
			fallthrough
		default:
			if toToken {
				switch {
				case operator == OpUndefined:
					var err error
					operator, err = getOperator(character, token, i)
					if err != nil {
						return "", err
					}
				case toMaxCharBuffer:
					if _, err := strconv.Atoi(string(character)); err == nil {
						maxCharBuffer.WriteRune(character)
					} else {
						return "", fmt.Errorf("illegal character identified in the token at col: %d", i)
					}
				default:
					switch character {
					case ':':
						toMaxCharBuffer = true
						maxCharBuffer.Reset()
					case '*':
						composite = true
					default:
						if err := validateLiteral(character, i); err != nil {
							return "", err
						}
						token.WriteRune(character)
					}
				}
			} else {
				result.WriteRune(character)
			}
		}
	}

	if !toToken {
		return result.String(), nil
	}

	return "", fmt.Errorf("unterminated token")
}

func addPrefix(op Op, result *strings.Builder) {
	switch op {
	case OpHash, OpDot, OpSlash, OpSemicolon, OpQuestionMark, OpAmp:
		result.WriteRune(rune(op))
	default:
		return
	}
}

func addSeparator(op Op, result *strings.Builder) {
	switch op {
	case OpDot, OpSlash, OpSemicolon:
		result.WriteRune(rune(op))
	case OpQuestionMark, OpAmp:
		result.WriteByte('&')
	default:
		result.WriteByte(',')
		return
	}
}

func addValue(op Op, token string, value string, result *strings.Builder, maxChar int) {
	switch op {
	case OpPlus, OpHash:
		addExpandedValue("", value, result, maxChar, false)
	case OpQuestionMark, OpAmp:
		result.WriteString(token + "=")
		addExpandedValue("", value, result, maxChar, true)
	case OpSemicolon:
		result.WriteString(token)
		if value != "" {
			result.WriteByte('=')
		}
		addExpandedValue("", value, result, maxChar, true)
	case OpDot, OpSlash, OpNone:
		addExpandedValue("", value, result, maxChar, true)
	}
}

func addValueElement(op Op, _, value string, result *strings.Builder, maxChar int) {
	switch op {
	case OpPlus, OpHash:
		addExpandedValue("", value, result, maxChar, false)
	case OpQuestionMark, OpAmp, OpSemicolon, OpDot, OpSlash, OpNone:
		addExpandedValue("", value, result, maxChar, true)
	}
}

func isSurrogate(str string) bool {
	_, width := utf8.DecodeRuneInString(str)
	return width > 1
}

func isIprivate(cp rune) bool {
	return 0xE000 <= cp && cp <= 0xF8FF
}

func isUcschar(cp rune) bool {
	return (0xA0 <= cp && cp <= 0xD7FF) ||
		(0xF900 <= cp && cp <= 0xFDCF) ||
		(0xFDF0 <= cp && cp <= 0xFFEF)
}

func addExpandedValue(prefix string, value string, result *strings.Builder, maxChar int, replaceReserved bool) {
	max := maxChar
	if maxChar == -1 || maxChar > len(value) {
		max = len(value)
	}
	reservedBuffer := &strings.Builder{}
	toReserved := false

	if max > 0 && prefix != "" {
		result.WriteString(prefix)
	}

	for i, character := range value {
		if i >= max {
			break
		}

		if character == '%' && !replaceReserved {
			reservedBuffer.Reset()
			toReserved = true
		}

		toAppend := string(character)
		if isSurrogate(toAppend) || replaceReserved || isUcschar(character) || isIprivate(character) {
			toAppend = url.QueryEscape(toAppend)
		}

		if toReserved {
			reservedBuffer.WriteString(toAppend)

			if reservedBuffer.Len() == 3 {
				encoded := true
				reserved := reservedBuffer.String()
				unescaped, err := url.QueryUnescape(reserved)
				if err != nil {
					encoded = (reserved == unescaped)
				}

				if encoded {
					result.WriteString(reserved)
				} else {
					result.WriteString("%25")
					// only if !replaceReserved
					result.WriteString(reservedBuffer.String()[1:])
				}
				reservedBuffer.Reset()
				toReserved = false
			}
		} else {
			switch character {
			case ' ':
				result.WriteString("%20")
			case '%':
				result.WriteString("%25")
			default:
				result.WriteString(toAppend)
			}
		}
	}

	if toReserved {
		result.WriteString("%25")
		result.WriteString(reservedBuffer.String()[1:])
	}
}

func getSubstitutionType(value any, col int) string {
	switch value.(type) {
	case nil:
		return SubstitutionTypeEmpty
	case string, float32, float64, int, int8, int16, int32, int64, bool, time.Time:
		return SubstitutionTypeString
	case []string, []float32, []float64, []int, []int8, []int16, []int32, []int64, []bool, []time.Time, []any:
		return SubstitutionTypeList
	case map[string]string, map[string]float32, map[string]float64, map[string]int, map[string]int8, map[string]int16, map[string]int32, map[string]int64, map[string]bool, map[string]time.Time, map[string]any:
		return SubstitutionTypeMap
	default:
		return fmt.Sprintf("illegal class passed as substitution, found %T at col: %d", value, col)
	}
}

func isEmpty(substType string, value any) bool {
	switch substType {
	case SubstitutionTypeString:
		switch value.(type) {
		case string:
			return value == nil
		default: // primitives are value types
			return false
		}
	case SubstitutionTypeList:
		return getListLength(value) == 0
	case SubstitutionTypeMap:
		return getMapLength(value) == 0
	default:
		return true
	}
}

func getListLength(value any) int {
	switch value.(type) {
	case []string:
		return len(value.([]string))
	case []float32:
		return len(value.([]float32))
	case []float64:
		return len(value.([]float64))
	case []int:
		return len(value.([]int))
	case []int8:
		return len(value.([]int8))
	case []int16:
		return len(value.([]int16))
	case []int32:
		return len(value.([]int32))
	case []int64:
		return len(value.([]int64))
	case []bool:
		return len(value.([]bool))
	case []time.Time:
		return len(value.([]time.Time))
	case []any:
		return len(value.([]any))
	}
	return 0
}

func getMapLength(value any) int {
	switch value.(type) {
	case map[string]string:
		return len(value.(map[string]string))
	case map[string]float32:
		return len(value.(map[string]float32))
	case map[string]float64:
		return len(value.(map[string]float64))
	case map[string]int:
		return len(value.(map[string]int))
	case map[string]int8:
		return len(value.(map[string]int8))
	case map[string]int16:
		return len(value.(map[string]int16))
	case map[string]int32:
		return len(value.(map[string]int32))
	case map[string]int64:
		return len(value.(map[string]int64))
	case map[string]bool:
		return len(value.(map[string]bool))
	case map[string]time.Time:
		return len(value.(map[string]time.Time))
	case map[string]any:
		return len(value.(map[string]any))
	}
	return 0
}

func convertNativeList(value any) ([]string, error) {
	var stringList = make([]string, getListLength(value))
	switch value.(type) {
	case []string:
		for index, val := range value.([]string) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringList[index] = str
		}
	case []float32:
		for index, val := range value.([]float32) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringList[index] = str
		}
	case []float64:
		for index, val := range value.([]float64) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringList[index] = str
		}
	case []int:
		for index, val := range value.([]int) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringList[index] = str
		}
	case []int8:
		for index, val := range value.([]int8) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringList[index] = str
		}
	case []int16:
		for index, val := range value.([]int16) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringList[index] = str
		}
	case []int32:
		for index, val := range value.([]int32) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringList[index] = str
		}
	case []int64:
		for index, val := range value.([]int64) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringList[index] = str
		}
	case []bool:
		for index, val := range value.([]bool) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringList[index] = str
		}
	case []time.Time:
		for index, val := range value.([]time.Time) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringList[index] = str
		}
	case []any:
		for index, val := range value.([]any) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringList[index] = str
		}
	default:
		return nil, fmt.Errorf("unrecognized type: %s", value)
	}
	return stringList, nil
}

func convertNativeMap(value any) (map[string]string, error) {
	var stringMap = make(map[string]string, getMapLength(value))
	switch value.(type) {
	case map[string]string:
		for key, val := range value.(map[string]string) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringMap[key] = str
		}
	case map[string]float32:
		for key, val := range value.(map[string]float32) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringMap[key] = str
		}
	case map[string]float64:
		for key, val := range value.(map[string]float64) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringMap[key] = str
		}
	case map[string]int:
		for key, val := range value.(map[string]int) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringMap[key] = str
		}
	case map[string]int8:
		for key, val := range value.(map[string]int8) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringMap[key] = str
		}
	case map[string]int16:
		for key, val := range value.(map[string]int16) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringMap[key] = str
		}
	case map[string]int32:
		for key, val := range value.(map[string]int32) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringMap[key] = str
		}
	case map[string]int64:
		for key, val := range value.(map[string]int64) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringMap[key] = str
		}
	case map[string]bool:
		for key, val := range value.(map[string]bool) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringMap[key] = str
		}
	case map[string]time.Time:
		for key, val := range value.(map[string]time.Time) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringMap[key] = str
		}
	case map[string]any:
		for key, val := range value.(map[string]any) {
			str, err := convertNativeTypes(val)
			if err != nil {
				return nil, err
			}
			stringMap[key] = str
		}
	default:
		return nil, fmt.Errorf("unrecognized type: %s", value)
	}
	return stringMap, nil
}

func convertNativeTypes(value any) (string, error) {
	switch value.(type) {
	case string, float32, float64, int, int8, int16, int32, int64, bool:
		return fmt.Sprintf("%v", value), nil
	case time.Time:
		return value.(time.Time).Format(time.RFC3339), nil
	default:
		return "", fmt.Errorf("unrecognized type: %s", value)
	}
}

func expandToken(
	operator Op,
	token string,
	composite bool,
	maxChar int,
	firstToken bool,
	substitutions Substitutions,
	result *strings.Builder,
	col int,
) (bool, error) {
	if len(token) == 0 {
		return false, fmt.Errorf("found an empty token at col: %d", col)
	}

	value, ok := substitutions[token]
	if !ok {
		return false, nil
	}

	substType := getSubstitutionType(value, col)
	if substType == SubstitutionTypeEmpty || isEmpty(substType, value) {
		return false, nil
	}

	if firstToken {
		addPrefix(operator, result)
	} else {
		addSeparator(operator, result)
	}

	switch substType {
	case SubstitutionTypeString:
		stringValue, err := convertNativeTypes(value)
		if err != nil {
			return false, err
		}
		addStringValue(operator, token, stringValue, result, maxChar)
	case SubstitutionTypeList:
		listValue, err := convertNativeList(value)
		if err != nil {
			return false, err
		}
		addListValue(operator, token, listValue, result, maxChar, composite)
	case SubstitutionTypeMap:
		mapValue, err := convertNativeMap(value)
		if err != nil {
			return false, err
		}
		err = addMapValue(operator, token, mapValue, result, maxChar, composite)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func addStringValue(operator Op, token string, value string, result *strings.Builder, maxChar int) {
	addValue(operator, token, value, result, maxChar)
}

func addListValue(operator Op, token string, value []string, result *strings.Builder, maxChar int, composite bool) {
	first := true
	for _, v := range value {
		if first {
			addValue(operator, token, v, result, maxChar)
			first = false
		} else {
			if composite {
				addSeparator(operator, result)
				addValue(operator, token, v, result, maxChar)
			} else {
				result.WriteString(",")
				addValueElement(operator, token, v, result, maxChar)
			}
		}
	}
}

func addMapValue(operator Op, token string, value map[string]string, result *strings.Builder, maxChar int, composite bool) error {
	if maxChar != -1 {
		return fmt.Errorf("value trimming is not allowed on Maps")
	}

	// workaround to make Map ordering not random
	// https://github.com/uri-templates/uritemplate-test/pull/58#issuecomment-1640029982
	keys := make([]string, 0, len(value))
	for k := range value {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i := range keys {
		k := keys[i]
		v := value[k]

		if composite {
			if i > 0 {
				addSeparator(operator, result)
			}
			addValueElement(operator, token, k, result, maxChar)
			result.WriteString("=")
		} else {
			if i == 0 {
				addValue(operator, token, k, result, maxChar)
			} else {
				result.WriteString(",")
				addValueElement(operator, token, k, result, maxChar)
			}
			result.WriteString(",")
		}
		addValueElement(operator, token, v, result, maxChar)
	}
	return nil
}
