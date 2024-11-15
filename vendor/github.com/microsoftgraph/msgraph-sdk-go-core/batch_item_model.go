package msgraphgocore

import (
	"errors"
	abs "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	"reflect"
)

// BatchItem is an instance of the BatchRequest payload to be later serialized to a json payload
type BatchItem interface {
	serialization.Parsable
	GetId() *string
	SetId(value *string)
	GetMethod() *string
	SetMethod(value *string)
	GetUrl() *string
	SetUrl(value *string)
	GetHeaders() RequestHeader
	SetHeaders(value RequestHeader)
	GetBody() RequestBody
	SetBody(value RequestBody)
	GetDependsOn() []string
	SetDependsOn(value []string)
	GetStatus() *int32
	SetStatus(value *int32)
	DependsOnItem(item BatchItem)
}

type batchItem struct {
	Id        *string
	method    *string
	Url       *string
	Headers   RequestHeader
	Body      RequestBody
	DependsOn []string
	Status    *int32
}

// DependsOnItem creates a dependency chain between BatchItems.If A depends on B, then B will be sent before B
// A batchItem can only depend on one other batchItem
// see: https://docs.microsoft.com/en-us/graph/known-issues#request-dependencies-are-limited
func (bi *batchItem) DependsOnItem(item BatchItem) {
	dependsOn := append(item.GetDependsOn(), *item.GetId())
	bi.SetDependsOn(dependsOn)
}

// NewBatchItem creates an instance of BatchItem
func NewBatchItem() BatchItem {
	return &batchItem{
		DependsOn: make([]string, 0),
	}
}

// GetId returns batch item `id` property
func (bi *batchItem) GetId() *string {
	return bi.Id
}

// SetId sets string value as batch item `id` property
func (bi *batchItem) SetId(value *string) {
	bi.Id = value
}

// GetMethod returns batch item `Method` property
func (bi *batchItem) GetMethod() *string {
	return bi.method
}

// SetMethod sets string value as batch item `Method` property
func (bi *batchItem) SetMethod(value *string) {
	bi.method = value
}

// GetUrl returns batch item `Url` property
func (bi *batchItem) GetUrl() *string {
	return bi.Url
}

// SetUrl sets string value as batch item `Url` property
func (bi *batchItem) SetUrl(value *string) {
	bi.Url = value
}

// GetHeaders returns batch item `Header` as a map[string]string
func (bi *batchItem) GetHeaders() RequestHeader {
	return bi.Headers
}

// SetHeaders sets map[string]string value as batch item `Header` property
func (bi *batchItem) SetHeaders(value RequestHeader) {
	bi.Headers = value
}

// GetBody returns batch item `RequestBody` property
func (bi *batchItem) GetBody() RequestBody {
	return bi.Body
}

// SetBody sets map[string]string value as batch item `RequestBody` property
func (bi *batchItem) SetBody(value RequestBody) {
	bi.Body = value
}

// GetDependsOn returns batch item `dependsOn` property as a string array
func (bi *batchItem) GetDependsOn() []string {
	return bi.DependsOn
}

// SetDependsOn sets []string value as batch item `dependsOn` property
func (bi *batchItem) SetDependsOn(value []string) {
	bi.DependsOn = value
}

// GetStatus returns batch item `status` property
func (bi *batchItem) GetStatus() *int32 {
	return bi.Status
}

// SetStatus sets int32 value as batch item `int` property
func (bi *batchItem) SetStatus(value *int32) {
	bi.Status = value
}

// Serialize serializes information the current object
func (bi *batchItem) Serialize(writer serialization.SerializationWriter) error {
	{
		err := writer.WriteStringValue("id", bi.GetId())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteStringValue("method", bi.GetMethod())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteStringValue("url", bi.GetUrl())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteAnyValue("headers", bi.GetHeaders())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteAnyValue("body", bi.GetBody())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteCollectionOfStringValues("dependsOn", bi.GetDependsOn())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteInt32Value("status", bi.GetStatus())
		if err != nil {
			return err
		}
	}
	return nil
}

// GetFieldDeserializers the deserialization information for the current model
func (bi *batchItem) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	res := make(map[string]func(serialization.ParseNode) error)
	res["id"] = abs.SetStringValue(bi.SetId)
	res["method"] = abs.SetStringValue(bi.SetMethod)
	res["url"] = abs.SetStringValue(bi.SetUrl)
	res["headers"] = func(n serialization.ParseNode) error {
		rawVal, err := n.GetRawValue()
		if err != nil {
			return err
		}

		if rawVal == nil {
			return nil
		}

		result, err := castMapOfStrings(rawVal)
		if err != nil {
			return err
		}

		bi.SetHeaders(result)
		return nil
	}
	res["body"] = func(n serialization.ParseNode) error {
		rawVal, err := n.GetRawValue()
		if err != nil {
			return err
		}

		if rawVal == nil {
			return nil
		}

		result, err := convertToMap(rawVal)
		if err != nil {
			return err
		}

		bi.SetBody(result)
		return nil
	}
	res["dependsOn"] = abs.SetCollectionOfPrimitiveValues("string", bi.SetDependsOn)
	res["status"] = abs.SetInt32Value(bi.SetStatus)
	return res
}

func convertToMap(rawVal interface{}) (map[string]interface{}, error) {
	kind := reflect.ValueOf(rawVal)
	if kind.Kind() == reflect.Map {
		result := make(map[string]interface{})
		err := deserializeMapped(kind, result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
	return nil, errors.New("interface was not a map")
}

func deserializeNode(value serialization.ParseNode) (interface{}, error) {
	rawVal, err := value.GetRawValue()
	if err != nil {
		return nil, err
	} else {
		kind := reflect.ValueOf(rawVal)
		if kind.Kind() == reflect.Map {

			result := make(map[string]interface{})
			err := deserializeMapped(kind, result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return deserializeValue(rawVal)
		}
	}
}

func deserializeMapped(v reflect.Value, result map[string]interface{}) error {
	for _, key := range v.MapKeys() {
		value, err := deserializeValue(v.MapIndex(key).Interface())
		if err != nil {
			return err
		} else {
			result[key.String()] = value
		}
	}
	return nil
}

func deserializeNodes(value []*jsonserialization.JsonParseNode) (interface{}, error) {
	slice := make([]interface{}, len(value))
	for index, element := range value {
		res, err := deserializeNode(element)
		if err != nil {
			return nil, err
		}
		slice[index] = res
	}
	return slice, nil
}

func deserializeValue(value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case int:
	case float64:
	case string:
		return value, nil
	case *int:
	case *float64:
	case *string:
		return value, nil
	case jsonserialization.JsonParseNode:
	case *jsonserialization.JsonParseNode:
		return deserializeNode(v)
	case []*jsonserialization.JsonParseNode:
		return deserializeNodes(v)
	case []jsonserialization.JsonParseNode:
		return deserializeNodes(abs.CollectionApply(v, func(x jsonserialization.JsonParseNode) *jsonserialization.JsonParseNode {
			return &x
		}))
	default:
		return value, nil
	}
	return nil, nil
}

func castMapOfStrings(rawVal interface{}) (map[string]string, error) {
	result := make(map[string]string)
	v := reflect.ValueOf(rawVal)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			val, err := deserializeValue(v.MapIndex(key).Interface())
			if err != nil {
				return nil, err
			}
			result[key.String()] = *(val.(*string))
		}
	}
	return result, nil
}

// CreateBatchRequestItemDiscriminator creates a new instance of the appropriate class based on discriminator value
func CreateBatchRequestItemDiscriminator(serialization.ParseNode) (serialization.Parsable, error) {
	var res batchItem
	return &res, nil
}
