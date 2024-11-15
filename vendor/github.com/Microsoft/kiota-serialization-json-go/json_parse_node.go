// Package jsonserialization is the default Kiota serialization implementation for JSON.
// It relies on the standard Go JSON library.
package jsonserialization

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// JsonParseNode is a ParseNode implementation for JSON.
type JsonParseNode struct {
	value                     interface{}
	onBeforeAssignFieldValues absser.ParsableAction
	onAfterAssignFieldValues  absser.ParsableAction
}

// NewJsonParseNode creates a new JsonParseNode.
func NewJsonParseNode(content []byte) (*JsonParseNode, error) {
	if len(content) == 0 {
		return nil, errors.New("content is empty")
	}
	if !json.Valid(content) {
		return nil, errors.New("invalid json type")
	}
	decoder := json.NewDecoder(bytes.NewReader(content))
	value, err := loadJsonTree(decoder)
	return value, err
}

func loadJsonTree(decoder *json.Decoder) (*JsonParseNode, error) {
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch token.(type) {
		case json.Delim:
			switch token.(json.Delim) {
			case '{':
				v := make(map[string]*JsonParseNode)
				for decoder.More() {
					key, err := decoder.Token()
					if err != nil {
						return nil, err
					}
					keyStr, ok := key.(string)
					if !ok {
						return nil, errors.New("key is not a string")
					}
					childNode, err := loadJsonTree(decoder)
					if err != nil {
						return nil, err
					}
					v[keyStr] = childNode
				}
				decoder.Token() // skip the closing curly
				result := &JsonParseNode{value: v}
				return result, nil
			case '[':
				v := make([]*JsonParseNode, 0)
				for decoder.More() {
					childNode, err := loadJsonTree(decoder)
					if err != nil {
						return nil, err
					}
					v = append(v, childNode)
				}
				decoder.Token() // skip the closing bracket
				result := &JsonParseNode{value: v}
				return result, nil
			case ']':
			case '}':
			}
		case json.Number:
			number := token.(json.Number)
			i, err := number.Int64()
			c := &JsonParseNode{}
			if err == nil {
				c.setValue(&i)
			} else {
				f, err := number.Float64()
				if err == nil {
					c.setValue(&f)
				} else {
					return nil, err
				}
			}
			return c, nil
		case string:
			v := token.(string)
			c := &JsonParseNode{}
			c.setValue(&v)
			return c, nil
		case bool:
			c := &JsonParseNode{}
			v := token.(bool)
			c.setValue(&v)
			return c, nil
		case int8:
			c := &JsonParseNode{}
			v := token.(int8)
			c.setValue(&v)
			return c, nil
		case byte:
			c := &JsonParseNode{}
			v := token.(byte)
			c.setValue(&v)
			return c, nil
		case float64:
			c := &JsonParseNode{}
			v := token.(float64)
			c.setValue(&v)
			return c, nil
		case float32:
			c := &JsonParseNode{}
			v := token.(float32)
			c.setValue(&v)
			return c, nil
		case int32:
			c := &JsonParseNode{}
			v := token.(int32)
			c.setValue(&v)
			return c, nil
		case int64:
			c := &JsonParseNode{}
			v := token.(int64)
			c.setValue(&v)
			return c, nil
		case nil:
			return nil, nil
		default:
		}
	}
	return nil, nil
}

// SetValue is obsolete, parse nodes are not meant to be settable externally
func (n *JsonParseNode) SetValue(value interface{}) {
	n.setValue(value)
}

// setValue sets the value represented by the node
func (n *JsonParseNode) setValue(value interface{}) {
	n.value = value
}

// GetChildNode returns a new parse node for the given identifier.
func (n *JsonParseNode) GetChildNode(index string) (absser.ParseNode, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	if index == "" {
		return nil, errors.New("index is empty")
	}
	childNodes, ok := n.value.(map[string]*JsonParseNode)
	if !ok || len(childNodes) == 0 {
		return nil, nil
	}

	childNode := childNodes[index]
	if childNode != nil {
		err := childNode.SetOnBeforeAssignFieldValues(n.GetOnBeforeAssignFieldValues())
		if err != nil {
			return nil, err
		}
		err = childNode.SetOnAfterAssignFieldValues(n.GetOnAfterAssignFieldValues())
		if err != nil {
			return nil, err
		}
	}

	return childNode, nil
}

// GetObjectValue returns the Parsable value from the node.
func (n *JsonParseNode) GetObjectValue(ctor absser.ParsableFactory) (absser.Parsable, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	if ctor == nil {
		return nil, errors.New("constructor is nil")
	}
	result, err := ctor(n)
	if err != nil {
		return nil, err
	}

	_, isUntypedNode := result.(absser.UntypedNodeable)
	if isUntypedNode {
		switch value := n.value.(type) {
		case *bool:
			return absser.NewUntypedBoolean(*value), nil
		case *string:
			return absser.NewUntypedString(*value), nil
		case *float32:
			return absser.NewUntypedFloat(*value), nil
		case *float64:
			return absser.NewUntypedDouble(*value), nil
		case *int32:
			return absser.NewUntypedInteger(*value), nil
		case *int64:
			return absser.NewUntypedLong(*value), nil
		case nil:
			return absser.NewUntypedNull(), nil
		case map[string]*JsonParseNode:
			properties := make(map[string]absser.UntypedNodeable)
			for key, value := range value {
				parsable, err := value.GetObjectValue(absser.CreateUntypedNodeFromDiscriminatorValue)
				if err != nil {
					return nil, errors.New("cannot parse object value")
				}
				if parsable == nil {
					parsable = absser.NewUntypedNull()
				}
				property, ok := parsable.(absser.UntypedNodeable)
				if ok {
					properties[key] = property
				}
			}
			return absser.NewUntypedObject(properties), nil
		case []*JsonParseNode:
			collection := make([]absser.UntypedNodeable, len(value))
			for index, node := range value {
				parsable, err := node.GetObjectValue(absser.CreateUntypedNodeFromDiscriminatorValue)
				if err != nil {
					return nil, errors.New("cannot parse object value")
				}
				if parsable == nil {
					parsable = absser.NewUntypedNull()
				}
				property, ok := parsable.(absser.UntypedNodeable)
				if ok {
					collection[index] = property
				}

			}
			return absser.NewUntypedArray(collection), nil
		default:
			return absser.NewUntypedNode(value), nil
		}
	}

	abstractions.InvokeParsableAction(n.GetOnBeforeAssignFieldValues(), result)
	properties, ok := n.value.(map[string]*JsonParseNode)
	fields := result.GetFieldDeserializers()
	if ok && len(properties) != 0 {
		itemAsHolder, isHolder := result.(absser.AdditionalDataHolder)
		var itemAdditionalData map[string]interface{}
		if isHolder {
			itemAdditionalData = itemAsHolder.GetAdditionalData()
			if itemAdditionalData == nil {
				itemAdditionalData = make(map[string]interface{})
				itemAsHolder.SetAdditionalData(itemAdditionalData)
			}
		}

		for key, value := range properties {
			field := fields[key]
			if value != nil {
				err := value.SetOnBeforeAssignFieldValues(n.GetOnBeforeAssignFieldValues())
				if err != nil {
					return nil, err
				}
				err = value.SetOnAfterAssignFieldValues(n.GetOnAfterAssignFieldValues())
				if err != nil {
					return nil, err
				}
			}
			if field == nil {
				if value != nil && isHolder {
					rawValue, err := value.GetRawValue()
					if err != nil {
						return nil, err
					}
					itemAdditionalData[key] = rawValue
				}
			} else {
				err := field(value)
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
func (n *JsonParseNode) GetCollectionOfObjectValues(ctor absser.ParsableFactory) ([]absser.Parsable, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	if ctor == nil {
		return nil, errors.New("ctor is nil")
	}
	nodes, ok := n.value.([]*JsonParseNode)
	if !ok {
		return nil, errors.New("value is not a collection")
	}
	result := make([]absser.Parsable, len(nodes))
	for i, v := range nodes {
		if v != nil {
			val, err := (*v).GetObjectValue(ctor)
			if err != nil {
				return nil, err
			}
			result[i] = val
		} else {
			result[i] = nil
		}
	}
	return result, nil
}

// GetCollectionOfPrimitiveValues returns the collection of primitive values from the node.
func (n *JsonParseNode) GetCollectionOfPrimitiveValues(targetType string) ([]interface{}, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	if targetType == "" {
		return nil, errors.New("targetType is empty")
	}
	nodes, ok := n.value.([]*JsonParseNode)
	if !ok {
		return nil, errors.New("value is not a collection")
	}
	result := make([]interface{}, len(nodes))
	for i, v := range nodes {
		if v != nil {
			val, err := v.getPrimitiveValue(targetType)
			if err != nil {
				return nil, err
			}
			result[i] = val
		} else {
			result[i] = nil
		}
	}
	return result, nil
}

func (n *JsonParseNode) getPrimitiveValue(targetType string) (interface{}, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
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
		return nil, fmt.Errorf("targetType %s is not supported", targetType)
	}
}

// GetCollectionOfEnumValues returns the collection of Enum values from the node.
func (n *JsonParseNode) GetCollectionOfEnumValues(parser absser.EnumFactory) ([]interface{}, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	if parser == nil {
		return nil, errors.New("parser is nil")
	}
	nodes, ok := n.value.([]*JsonParseNode)
	if !ok {
		return nil, errors.New("value is not a collection")
	}
	result := make([]interface{}, len(nodes))
	for i, v := range nodes {
		if v != nil {
			val, err := v.GetEnumValue(parser)
			if err != nil {
				return nil, err
			}
			result[i] = val
		} else {
			result[i] = nil
		}
	}
	return result, nil
}

// GetStringValue returns a String value from the nodes.
func (n *JsonParseNode) GetStringValue() (*string, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}

	val, ok := n.value.(*string)
	if !ok {
		return nil, fmt.Errorf("type '%T' is not compatible with type string", n.value)
	}
	return val, nil
}

// GetBoolValue returns a Bool value from the nodes.
func (n *JsonParseNode) GetBoolValue() (*bool, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}

	val, ok := n.value.(*bool)
	if !ok {
		return nil, fmt.Errorf("type '%T' is not compatible with type bool", n.value)
	}
	return val, nil
}

// GetInt8Value returns a int8 value from the nodes.
func (n *JsonParseNode) GetInt8Value() (*int8, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val int8

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetBoolValue returns a Bool value from the nodes.
func (n *JsonParseNode) GetByteValue() (*byte, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val byte

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetFloat32Value returns a Float32 value from the nodes.
func (n *JsonParseNode) GetFloat32Value() (*float32, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val float32

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetFloat64Value returns a Float64 value from the nodes.
func (n *JsonParseNode) GetFloat64Value() (*float64, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val float64

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetInt32Value returns a Int32 value from the nodes.
func (n *JsonParseNode) GetInt32Value() (*int32, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val int32

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetInt64Value returns a Int64 value from the nodes.
func (n *JsonParseNode) GetInt64Value() (*int64, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	var val int64

	if err := as(n.value, &val); err != nil {
		return nil, err
	}

	return &val, nil
}

// GetTimeValue returns a Time value from the nodes.
func (n *JsonParseNode) GetTimeValue() (*time.Time, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}

	// if string does not have timezone information, add local timezone
	if len(*v) == 19 {
		*v = *v + time.Now().Format("-07:00")
	}
	parsed, err := time.Parse(time.RFC3339, *v)
	return &parsed, err
}

// GetISODurationValue returns a ISODuration value from the nodes.
func (n *JsonParseNode) GetISODurationValue() (*absser.ISODuration, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
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
func (n *JsonParseNode) GetTimeOnlyValue() (*absser.TimeOnly, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
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
func (n *JsonParseNode) GetDateOnlyValue() (*absser.DateOnly, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
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
func (n *JsonParseNode) GetUUIDValue() (*uuid.UUID, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
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
func (n *JsonParseNode) GetEnumValue(parser absser.EnumFactory) (interface{}, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
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
func (n *JsonParseNode) GetByteArrayValue() ([]byte, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
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
func (n *JsonParseNode) GetRawValue() (interface{}, error) {
	if isNil(n) || isNil(n.value) {
		return nil, nil
	}
	switch v := n.value.(type) {
	case *JsonParseNode:
		return v.GetRawValue()
	case []*JsonParseNode:
		result := make([]interface{}, len(v))
		for i, x := range v {
			val, err := x.GetRawValue()
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return result, nil
	case map[string]*JsonParseNode:
		m := make(map[string]interface{})
		for key, element := range v {
			elementVal, err := element.GetRawValue()
			if err != nil {
				return nil, err
			}
			m[key] = elementVal
		}
		return m, nil
	default:
		return n.value, nil
	}
}

func (n *JsonParseNode) GetOnBeforeAssignFieldValues() absser.ParsableAction {
	return n.onBeforeAssignFieldValues
}

func (n *JsonParseNode) SetOnBeforeAssignFieldValues(action absser.ParsableAction) error {
	n.onBeforeAssignFieldValues = action
	return nil
}

func (n *JsonParseNode) GetOnAfterAssignFieldValues() absser.ParsableAction {
	return n.onAfterAssignFieldValues
}

func (n *JsonParseNode) SetOnAfterAssignFieldValues(action absser.ParsableAction) error {
	n.onAfterAssignFieldValues = action
	return nil
}
