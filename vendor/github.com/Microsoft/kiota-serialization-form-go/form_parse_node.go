// Package formserialization is the default Kiota serialization implementation for URI form encoded.
package formserialization

import (
	"encoding/base64"
	"errors"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// FormParseNode is a ParseNode implementation for JSON.
type FormParseNode struct {
	value                     string
	fields                    map[string]string
	onBeforeAssignFieldValues absser.ParsableAction
	onAfterAssignFieldValues  absser.ParsableAction
}

// NewFormParseNode creates a new FormParseNode.
func NewFormParseNode(content []byte) (*FormParseNode, error) {
	if len(content) == 0 {
		return nil, errors.New("content is empty")
	}
	rawValue := string(content)
	fields, err := loadFields(rawValue)
	if err != nil {
		return nil, err
	}
	return &FormParseNode{
		value:  rawValue,
		fields: fields,
	}, nil
}
func loadFields(value string) (map[string]string, error) {
	result := make(map[string]string)
	if len(value) == 0 {
		return result, nil
	}
	parts := strings.Split(value, "&")
	for _, part := range parts {
		keyValue := strings.Split(part, "=")
		if len(keyValue) == 2 {
			key, err := sanitizeKey(keyValue[0])
			if err != nil {
				return nil, err
			}
			if result[key] == "" {
				result[key] = keyValue[1]
			} else {
				result[key] += "," + keyValue[1]
			}
		}
	}
	return result, nil
}
func sanitizeKey(key string) (string, error) {
	if key == "" {
		return "", nil
	}
	res, err := url.QueryUnescape(key)
	if err != nil {
		return "", err
	}
	return strings.Trim(res, " "), nil
}

// GetChildNode returns a new parse node for the given identifier.
func (n *FormParseNode) GetChildNode(index string) (absser.ParseNode, error) {
	if index == "" {
		return nil, errors.New("index is empty")
	}
	key, err := sanitizeKey(index)
	if err != nil {
		return nil, err
	}
	fieldValue := n.fields[key]
	if fieldValue == "" {
		return nil, nil
	}

	node := &FormParseNode{
		value:                     fieldValue,
		onBeforeAssignFieldValues: n.GetOnBeforeAssignFieldValues(),
		onAfterAssignFieldValues:  n.GetOnAfterAssignFieldValues(),
	}
	return node, nil
}

// GetObjectValue returns the Parsable value from the node.
func (n *FormParseNode) GetObjectValue(ctor absser.ParsableFactory) (absser.Parsable, error) {
	if ctor == nil {
		return nil, errors.New("constructor is nil")
	}
	if n == nil || n.value == "" {
		return nil, nil
	}
	result, err := ctor(n)
	if err != nil {
		return nil, err
	}
	abstractions.InvokeParsableAction(n.GetOnBeforeAssignFieldValues(), result)
	fields := result.GetFieldDeserializers()
	if len(n.fields) != 0 {
		itemAsHolder, isHolder := result.(absser.AdditionalDataHolder)
		var itemAdditionalData map[string]interface{}
		if isHolder {
			itemAdditionalData = itemAsHolder.GetAdditionalData()
			if itemAdditionalData == nil {
				itemAdditionalData = make(map[string]interface{})
				itemAsHolder.SetAdditionalData(itemAdditionalData)
			}
		}

		for key, value := range n.fields {
			field := fields[key]
			if field == nil {
				if value != "" && isHolder {
					if err != nil {
						return nil, err
					}
					itemAdditionalData[key] = value
				}
			} else {
				childNode, err := n.GetChildNode(key)
				if err != nil {
					return nil, err
				}
				err = field(childNode)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	abstractions.InvokeParsableAction(n.GetOnAfterAssignFieldValues(), result)
	return result, nil
}

// GetCollectionOfObjectValues returns the collection of Parsable values from the node.
func (n *FormParseNode) GetCollectionOfObjectValues(ctor absser.ParsableFactory) ([]absser.Parsable, error) {
	return nil, errors.New("collections are not supported in form serialization")
}

// GetCollectionOfPrimitiveValues returns the collection of primitive values from the node.
func (n *FormParseNode) GetCollectionOfPrimitiveValues(targetType string) ([]interface{}, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	if targetType == "" {
		return nil, errors.New("targetType is empty")
	}
	valueList := strings.Split(n.value, ",")

	result := make([]interface{}, len(valueList))
	for i, element := range valueList {
		parseNode, err := NewFormParseNode([]byte(element))
		if err != nil {
			return nil, err
		}

		val, err := parseNode.getPrimitiveValue(targetType)
		if err != nil {
			return nil, err
		}
		result[i] = val
	}
	return result, nil
}

func (n *FormParseNode) getPrimitiveValue(targetType string) (interface{}, error) {
	switch targetType {
	case "string":
		return n.GetStringValue()
	case "bool":
		return n.GetBoolValue()
	case "uint8":
		return n.GetInt8Value()
	case "byte":
		return n.GetByteValue()
	case "float32":
		return n.GetFloat32Value()
	case "float64":
		return n.GetFloat64Value()
	case "int32":
		return n.GetInt32Value()
	case "int64":
		return n.GetInt64Value()
	case "time":
		return n.GetTimeValue()
	case "timeonly":
		return n.GetTimeOnlyValue()
	case "dateonly":
		return n.GetDateOnlyValue()
	case "isoduration":
		return n.GetISODurationValue()
	case "uuid":
		return n.GetUUIDValue()
	case "base64":
		return n.GetByteArrayValue()
	default:
		return nil, errors.New("targetType is not supported")
	}
}

// GetCollectionOfEnumValues returns the collection of Enum values from the node.
func (n *FormParseNode) GetCollectionOfEnumValues(parser absser.EnumFactory) ([]interface{}, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	if parser == nil {
		return nil, errors.New("parser is nil")
	}
	rawValues := strings.Split(n.value, ",")
	result := make([]interface{}, len(rawValues))
	for i, rawValue := range rawValues {
		node := &FormParseNode{
			value: rawValue,
		}
		val, err := node.GetEnumValue(parser)
		if err != nil {
			return nil, err
		}
		result[i] = val
	}
	return result, nil
}

// GetStringValue returns a String value from the nodes.
func (n *FormParseNode) GetStringValue() (*string, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	res, err := url.QueryUnescape(n.value)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetBoolValue returns a Bool value from the nodes.
func (n *FormParseNode) GetBoolValue() (*bool, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	res, err := strconv.ParseBool(n.value)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetInt8Value returns a int8 value from the nodes.
func (n *FormParseNode) GetInt8Value() (*int8, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	res, err := strconv.ParseInt(n.value, 10, 8)
	if err != nil {
		return nil, err
	}
	cast := int8(res)
	return &cast, nil
}

// GetBoolValue returns a Bool value from the nodes.
func (n *FormParseNode) GetByteValue() (*byte, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	res, err := strconv.ParseInt(n.value, 10, 8)
	if err != nil {
		return nil, err
	}
	cast := byte(res)
	return &cast, nil
}

// GetFloat32Value returns a Float32 value from the nodes.
func (n *FormParseNode) GetFloat32Value() (*float32, error) {
	v, err := n.GetFloat64Value()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	cast := float32(*v)
	return &cast, nil
}

// GetFloat64Value returns a Float64 value from the nodes.
func (n *FormParseNode) GetFloat64Value() (*float64, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	res, err := strconv.ParseFloat(n.value, 64)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetInt32Value returns a Int32 value from the nodes.
func (n *FormParseNode) GetInt32Value() (*int32, error) {
	v, err := n.GetFloat64Value()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	cast := int32(*v)
	return &cast, nil
}

// GetInt64Value returns a Int64 value from the nodes.
func (n *FormParseNode) GetInt64Value() (*int64, error) {
	v, err := n.GetFloat64Value()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	cast := int64(*v)
	return &cast, nil
}

// GetTimeValue returns a Time value from the nodes.
func (n *FormParseNode) GetTimeValue() (*time.Time, error) {
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, *v)
	return &parsed, err
}

// GetISODurationValue returns a ISODuration value from the nodes.
func (n *FormParseNode) GetISODurationValue() (*absser.ISODuration, error) {
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return absser.ParseISODuration(*v)
}

// GetTimeOnlyValue returns a TimeOnly value from the nodes.
func (n *FormParseNode) GetTimeOnlyValue() (*absser.TimeOnly, error) {
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return absser.ParseTimeOnly(*v)
}

// GetDateOnlyValue returns a DateOnly value from the nodes.
func (n *FormParseNode) GetDateOnlyValue() (*absser.DateOnly, error) {
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return absser.ParseDateOnly(*v)
}

// GetUUIDValue returns a UUID value from the nodes.
func (n *FormParseNode) GetUUIDValue() (*uuid.UUID, error) {
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	parsed, err := uuid.Parse(*v)
	return &parsed, err
}

// GetEnumValue returns a Enum value from the nodes.
func (n *FormParseNode) GetEnumValue(parser absser.EnumFactory) (interface{}, error) {
	if parser == nil {
		return nil, errors.New("parser is nil")
	}
	s, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}
	return parser(*s)
}

// GetByteArrayValue returns a ByteArray value from the nodes.
func (n *FormParseNode) GetByteArrayValue() ([]byte, error) {
	s, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}
	return base64.StdEncoding.DecodeString(*s)
}

// GetRawValue returns a ByteArray value from the nodes.
func (n *FormParseNode) GetRawValue() (interface{}, error) {
	if n == nil || n.value == "" {
		return nil, nil
	}
	res := n.value
	return &res, nil
}

func (n *FormParseNode) GetOnBeforeAssignFieldValues() absser.ParsableAction {
	return n.onBeforeAssignFieldValues
}

func (n *FormParseNode) SetOnBeforeAssignFieldValues(action absser.ParsableAction) error {
	n.onBeforeAssignFieldValues = action
	return nil
}

func (n *FormParseNode) GetOnAfterAssignFieldValues() absser.ParsableAction {
	return n.onAfterAssignFieldValues
}

func (n *FormParseNode) SetOnAfterAssignFieldValues(action absser.ParsableAction) error {
	n.onAfterAssignFieldValues = action
	return nil
}
