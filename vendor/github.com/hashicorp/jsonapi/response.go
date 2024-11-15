package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrBadJSONAPIStructTag is returned when the Struct field's JSON API
	// annotation is invalid.
	ErrBadJSONAPIStructTag = errors.New("Bad jsonapi struct tag format")
	// ErrBadJSONAPIID is returned when the Struct JSON API annotated "id" field
	// was not a valid numeric type.
	ErrBadJSONAPIID = errors.New(
		"id should be either string, int(8,16,32,64) or uint(8,16,32,64)")
	// ErrExpectedSlice is returned when a variable or argument was expected to
	// be a slice of *Structs; MarshalMany will return this error when its
	// interface{} argument is invalid.
	ErrExpectedSlice = errors.New("models should be a slice of struct pointers")
	// ErrUnexpectedType is returned when marshalling an interface; the interface
	// had to be a pointer or a slice; otherwise this error is returned.
	ErrUnexpectedType = errors.New("models should be a struct pointer or slice of struct pointers")
	// ErrUnexpectedNil is returned when a slice of relation structs contains nil values
	ErrUnexpectedNil = errors.New("slice of struct pointers cannot contain nil")
)

// MarshalPayload writes a jsonapi response for one or many records. The
// related records are sideloaded into the "included" array. If this method is
// given a struct pointer as an argument it will serialize in the form
// "data": {...}. If this method is given a slice of pointers, this method will
// serialize in the form "data": [...]
//
// One Example: you could pass it, w, your http.ResponseWriter, and, models, a
// ptr to a Blog to be written to the response body:
//
//	 func ShowBlog(w http.ResponseWriter, r *http.Request) {
//		 blog := &Blog{}
//
//		 w.Header().Set("Content-Type", jsonapi.MediaType)
//		 w.WriteHeader(http.StatusOK)
//
//		 if err := jsonapi.MarshalPayload(w, blog); err != nil {
//			 http.Error(w, err.Error(), http.StatusInternalServerError)
//		 }
//	 }
//
// Many Example: you could pass it, w, your http.ResponseWriter, and, models, a
// slice of Blog struct instance pointers to be written to the response body:
//
//		 func ListBlogs(w http.ResponseWriter, r *http.Request) {
//	    blogs := []*Blog{}
//
//			 w.Header().Set("Content-Type", jsonapi.MediaType)
//			 w.WriteHeader(http.StatusOK)
//
//			 if err := jsonapi.MarshalPayload(w, blogs); err != nil {
//				 http.Error(w, err.Error(), http.StatusInternalServerError)
//			 }
//		 }
func MarshalPayload(w io.Writer, models interface{}) error {
	payload, err := Marshal(models)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(payload)
}

// Marshal does the same as MarshalPayload except it just returns the payload
// and doesn't write out results. Useful if you use your own JSON rendering
// library.
func Marshal(models interface{}) (Payloader, error) {
	switch vals := reflect.ValueOf(models); vals.Kind() {
	case reflect.Slice:
		m, err := convertToSliceInterface(&models)
		if err != nil {
			return nil, err
		}

		payload, err := marshalMany(m)
		if err != nil {
			return nil, err
		}

		if linkableModels, isLinkable := models.(Linkable); isLinkable {
			jl := linkableModels.JSONAPILinks()
			if er := jl.validate(); er != nil {
				return nil, er
			}
			payload.Links = linkableModels.JSONAPILinks()
		}

		if metableModels, ok := models.(Metable); ok {
			payload.Meta = metableModels.JSONAPIMeta()
		}

		return payload, nil
	case reflect.Ptr:
		// Check that the pointer was to a struct
		if reflect.Indirect(vals).Kind() != reflect.Struct {
			return nil, ErrUnexpectedType
		}
		return marshalOne(models)
	default:
		return nil, ErrUnexpectedType
	}
}

// MarshalPayloadWithoutIncluded writes a jsonapi response with one or many
// records, without the related records sideloaded into "included" array.
// If you want to serialize the relations into the "included" array see
// MarshalPayload.
//
// models interface{} should be either a struct pointer or a slice of struct
// pointers.
func MarshalPayloadWithoutIncluded(w io.Writer, model interface{}) error {
	payload, err := Marshal(model)
	if err != nil {
		return err
	}
	payload.clearIncluded()

	return json.NewEncoder(w).Encode(payload)
}

// marshalOne does the same as MarshalOnePayload except it just returns the
// payload and doesn't write out results. Useful is you use your JSON rendering
// library.
func marshalOne(model interface{}) (*OnePayload, error) {
	included := make(map[string]*Node)

	rootNode, err := visitModelNode(model, &included, true)
	if err != nil {
		return nil, err
	}
	payload := &OnePayload{Data: rootNode}

	payload.Included = nodeMapValues(&included)

	return payload, nil
}

// marshalMany does the same as MarshalManyPayload except it just returns the
// payload and doesn't write out results. Useful is you use your JSON rendering
// library.
func marshalMany(models []interface{}) (*ManyPayload, error) {
	payload := &ManyPayload{
		Data: []*Node{},
	}
	included := map[string]*Node{}

	for _, model := range models {
		node, err := visitModelNode(model, &included, true)
		if err != nil {
			return nil, err
		}
		payload.Data = append(payload.Data, node)
	}
	payload.Included = nodeMapValues(&included)

	return payload, nil
}

// MarshalOnePayloadEmbedded - This method not meant to for use in
// implementation code, although feel free.  The purpose of this
// method is for use in tests.  In most cases, your request
// payloads for create will be embedded rather than sideloaded for
// related records. This method will serialize a single struct
// pointer into an embedded json response. In other words, there
// will be no, "included", array in the json all relationships will
// be serailized inline in the data.
//
// However, in tests, you may want to construct payloads to post
// to create methods that are embedded to most closely resemble
// the payloads that will be produced by the client. This is what
// this method is intended for.
//
// model interface{} should be a pointer to a struct.
func MarshalOnePayloadEmbedded(w io.Writer, model interface{}) error {
	rootNode, err := visitModelNode(model, nil, false)
	if err != nil {
		return err
	}

	payload := &OnePayload{Data: rootNode}

	return json.NewEncoder(w).Encode(payload)
}

// selectChoiceTypeStructField returns the first non-nil struct pointer field in the
// specified struct value that has a jsonapi type field defined within it.
// An error is returned if there are no fields matching that definition.
func selectChoiceTypeStructField(structValue reflect.Value) (reflect.Value, error) {
	for i := 0; i < structValue.NumField(); i++ {
		choiceFieldValue := structValue.Field(i)
		choiceTypeField := choiceFieldValue.Type()

		// Must be a pointer
		if choiceTypeField.Kind() != reflect.Ptr {
			continue
		}

		// Must not be nil
		if choiceFieldValue.IsNil() {
			continue
		}

		subtype := choiceTypeField.Elem()
		_, err := jsonapiTypeOfModel(subtype)
		if err == nil {
			return choiceFieldValue, nil
		}
	}

	return reflect.Value{}, errors.New("no non-nil choice field was found in the specified struct")
}

func visitModelNode(model interface{}, included *map[string]*Node,
	sideload bool) (*Node, error) {
	node := new(Node)

	var er error
	value := reflect.ValueOf(model)
	if value.IsNil() {
		return nil, nil
	}

	modelValue := value.Elem()
	modelType := value.Type().Elem()

	for i := 0; i < modelValue.NumField(); i++ {
		fieldValue := modelValue.Field(i)
		structField := modelValue.Type().Field(i)
		tag := structField.Tag.Get(annotationJSONAPI)
		if tag == "" {
			continue
		}

		fieldType := modelType.Field(i)

		args := strings.Split(tag, annotationSeparator)

		if len(args) < 1 {
			er = ErrBadJSONAPIStructTag
			break
		}

		annotation := args[0]

		if (annotation == annotationClientID && len(args) != 1) ||
			(annotation != annotationClientID && len(args) < 2) {
			er = ErrBadJSONAPIStructTag
			break
		}

		if annotation == annotationPrimary {
			v := fieldValue

			// Deal with PTRS
			var kind reflect.Kind
			if fieldValue.Kind() == reflect.Ptr {
				kind = fieldType.Type.Elem().Kind()
				v = reflect.Indirect(fieldValue)
			} else {
				kind = fieldType.Type.Kind()
			}

			// Handle allowed types
			switch kind {
			case reflect.String:
				node.ID = v.Interface().(string)
			case reflect.Int:
				node.ID = strconv.FormatInt(int64(v.Interface().(int)), 10)
			case reflect.Int8:
				node.ID = strconv.FormatInt(int64(v.Interface().(int8)), 10)
			case reflect.Int16:
				node.ID = strconv.FormatInt(int64(v.Interface().(int16)), 10)
			case reflect.Int32:
				node.ID = strconv.FormatInt(int64(v.Interface().(int32)), 10)
			case reflect.Int64:
				node.ID = strconv.FormatInt(v.Interface().(int64), 10)
			case reflect.Uint:
				node.ID = strconv.FormatUint(uint64(v.Interface().(uint)), 10)
			case reflect.Uint8:
				node.ID = strconv.FormatUint(uint64(v.Interface().(uint8)), 10)
			case reflect.Uint16:
				node.ID = strconv.FormatUint(uint64(v.Interface().(uint16)), 10)
			case reflect.Uint32:
				node.ID = strconv.FormatUint(uint64(v.Interface().(uint32)), 10)
			case reflect.Uint64:
				node.ID = strconv.FormatUint(v.Interface().(uint64), 10)
			default:
				// We had a JSON float (numeric), but our field was not one of the
				// allowed numeric types
				er = ErrBadJSONAPIID
			}

			if er != nil {
				break
			}

			node.Type = args[1]
		} else if annotation == annotationClientID {
			clientID := fieldValue.String()
			if clientID != "" {
				node.ClientID = clientID
			}
		} else if annotation == annotationAttribute {
			var omitEmpty, iso8601, rfc3339 bool

			if len(args) > 2 {
				for _, arg := range args[2:] {
					switch arg {
					case annotationOmitEmpty:
						omitEmpty = true
					case annotationISO8601:
						iso8601 = true
					case annotationRFC3339:
						rfc3339 = true
					}
				}
			}

			if node.Attributes == nil {
				node.Attributes = make(map[string]interface{})
			}

			// Handle Nullable[T]
			if strings.HasPrefix(fieldValue.Type().Name(), "NullableAttr[") {
				// handle unspecified
				if fieldValue.IsNil() {
					continue
				}

				// handle null
				if fieldValue.MapIndex(reflect.ValueOf(false)).IsValid() {
					node.Attributes[args[1]] = json.RawMessage("null")
					continue
				} else {

					// handle value
					fieldValue = fieldValue.MapIndex(reflect.ValueOf(true))
				}
			}

			if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				t := fieldValue.Interface().(time.Time)

				if t.IsZero() {
					continue
				}

				if iso8601 {
					node.Attributes[args[1]] = t.UTC().Format(iso8601TimeFormat)
				} else if rfc3339 {
					node.Attributes[args[1]] = t.UTC().Format(time.RFC3339)
				} else {
					node.Attributes[args[1]] = t.Unix()
				}
			} else if fieldValue.Type() == reflect.TypeOf(new(time.Time)) {
				// A time pointer may be nil
				if fieldValue.IsNil() {
					if omitEmpty {
						continue
					}

					node.Attributes[args[1]] = nil
				} else {
					tm := fieldValue.Interface().(*time.Time)

					if tm.IsZero() && omitEmpty {
						continue
					}

					if iso8601 {
						node.Attributes[args[1]] = tm.UTC().Format(iso8601TimeFormat)
					} else if rfc3339 {
						node.Attributes[args[1]] = tm.UTC().Format(time.RFC3339)
					} else {
						node.Attributes[args[1]] = tm.Unix()
					}
				}
			} else {
				// Dealing with a fieldValue that is not a time
				emptyValue := reflect.Zero(fieldValue.Type())

				// See if we need to omit this field
				if omitEmpty && reflect.DeepEqual(fieldValue.Interface(), emptyValue.Interface()) {
					continue
				}

				strAttr, ok := fieldValue.Interface().(string)
				if ok {
					node.Attributes[args[1]] = strAttr
				} else {
					node.Attributes[args[1]] = fieldValue.Interface()
				}
			}
		} else if annotation == annotationRelation || annotation == annotationPolyRelation {
			var omitEmpty bool

			//add support for 'omitempty' struct tag for marshaling as absent
			if len(args) > 2 {
				omitEmpty = args[2] == annotationOmitEmpty
			}

			isSlice := fieldValue.Type().Kind() == reflect.Slice
			if omitEmpty &&
				(isSlice && fieldValue.Len() < 1 ||
					(!isSlice && fieldValue.IsNil())) {
				continue
			}

			if annotation == annotationPolyRelation {
				// for polyrelation, we'll snoop out the actual relation model
				// through the choice type value by choosing the first non-nil
				// field that has a jsonapi type annotation and overwriting
				// `fieldValue` so normal annotation-assisted marshaling
				// can continue
				if !isSlice {
					choiceValue := fieldValue

					// must be a pointer type
					if choiceValue.Type().Kind() != reflect.Ptr {
						er = ErrUnexpectedType
						break
					}

					if choiceValue.IsNil() {
						fieldValue = reflect.ValueOf(nil)
					}
					structValue := choiceValue.Elem()

					// Short circuit if field is omitted from model
					if !structValue.IsValid() {
						break
					}

					if found, err := selectChoiceTypeStructField(structValue); err == nil {
						fieldValue = found
					}
				} else {
					// A slice polyrelation field can be... polymorphic... meaning
					// that we might snoop different types within each slice element.
					// Each snooped value will added to this collection and then
					// the recursion will take care of the rest. The only special case
					// is nil. For that, we'll just choose the first
					collection := make([]interface{}, 0)

					for i := 0; i < fieldValue.Len(); i++ {
						itemValue := fieldValue.Index(i)
						// Once again, must be a pointer type
						if itemValue.Type().Kind() != reflect.Ptr {
							er = ErrUnexpectedType
							break
						}

						if itemValue.IsNil() {
							er = ErrUnexpectedNil
							break
						}

						structValue := itemValue.Elem()

						if found, err := selectChoiceTypeStructField(structValue); err == nil {
							collection = append(collection, found.Interface())
						}
					}

					if er != nil {
						break
					}

					fieldValue = reflect.ValueOf(collection)
				}
			}

			if node.Relationships == nil {
				node.Relationships = make(map[string]interface{})
			}

			var relLinks *Links
			if linkableModel, ok := model.(RelationshipLinkable); ok {
				relLinks = linkableModel.JSONAPIRelationshipLinks(args[1])
			}

			var relMeta *Meta
			if metableModel, ok := model.(RelationshipMetable); ok {
				relMeta = metableModel.JSONAPIRelationshipMeta(args[1])
			}

			if isSlice {
				// to-many relationship
				relationship, err := visitModelNodeRelationships(
					fieldValue,
					included,
					sideload,
				)
				if err != nil {
					er = err
					break
				}
				relationship.Links = relLinks
				relationship.Meta = relMeta

				if sideload {
					shallowNodes := []*Node{}
					for _, n := range relationship.Data {
						appendIncluded(included, n)
						shallowNodes = append(shallowNodes, toShallowNode(n))
					}

					node.Relationships[args[1]] = &RelationshipManyNode{
						Data:  shallowNodes,
						Links: relationship.Links,
						Meta:  relationship.Meta,
					}
				} else {
					node.Relationships[args[1]] = relationship
				}
			} else {
				// to-one relationships

				// Handle null relationship case
				if fieldValue.IsNil() {
					node.Relationships[args[1]] = &RelationshipOneNode{Data: nil}
					continue
				}

				relationship, err := visitModelNode(
					fieldValue.Interface(),
					included,
					sideload,
				)
				if err != nil {
					er = err
					break
				}

				if sideload {
					appendIncluded(included, relationship)
					node.Relationships[args[1]] = &RelationshipOneNode{
						Data:  toShallowNode(relationship),
						Links: relLinks,
						Meta:  relMeta,
					}
				} else {
					node.Relationships[args[1]] = &RelationshipOneNode{
						Data:  relationship,
						Links: relLinks,
						Meta:  relMeta,
					}
				}
			}
		} else if annotation == annotationLinks {
			// Nothing. Ignore this field, as Links fields are only for unmarshaling requests.
			// The Linkable interface methods are used for marshaling data in a response.
		} else {
			er = ErrBadJSONAPIStructTag
			break
		}
	}

	if er != nil {
		return nil, er
	}

	if linkableModel, isLinkable := model.(Linkable); isLinkable {
		jl := linkableModel.JSONAPILinks()
		if er := jl.validate(); er != nil {
			return nil, er
		}
		node.Links = linkableModel.JSONAPILinks()
	}

	if metableModel, ok := model.(Metable); ok {
		node.Meta = metableModel.JSONAPIMeta()
	}

	return node, nil
}

// toShallowNode takes a node and returns a shallow version of the node.
// If the ID is empty, we include attributes into the shallow version.
//
// An example of where this is useful would be if an object
// within a relationship can be created at the same time as
// the root node.
//
// This is not 1.0 jsonapi spec compliant--it's a bespoke variation on
// resource object identifiers discussed in the pending 1.1 spec.
func toShallowNode(node *Node) *Node {
	ret := &Node{Type: node.Type}
	if node.ID == "" {
		ret.Attributes = node.Attributes
	} else {
		ret.ID = node.ID
	}
	return ret
}

func visitModelNodeRelationships(models reflect.Value, included *map[string]*Node,
	sideload bool) (*RelationshipManyNode, error) {
	nodes := []*Node{}

	for i := 0; i < models.Len(); i++ {
		model := models.Index(i)
		if !model.IsValid() || model.IsNil() {
			return nil, ErrUnexpectedNil
		}

		n := model.Interface()

		node, err := visitModelNode(n, included, sideload)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return &RelationshipManyNode{Data: nodes}, nil
}

func appendIncluded(m *map[string]*Node, nodes ...*Node) {
	included := *m

	for _, n := range nodes {
		k := fmt.Sprintf("%s,%s", n.Type, n.ID)

		if _, hasNode := included[k]; hasNode {
			continue
		}

		included[k] = n
	}
}

func nodeMapValues(m *map[string]*Node) []*Node {
	mp := *m
	nodes := make([]*Node, len(mp))

	i := 0
	for _, n := range mp {
		nodes[i] = n
		i++
	}

	return nodes
}

func convertToSliceInterface(i *interface{}) ([]interface{}, error) {
	vals := reflect.ValueOf(*i)
	if vals.Kind() != reflect.Slice {
		return nil, ErrExpectedSlice
	}
	var response []interface{}
	for x := 0; x < vals.Len(); x++ {
		response = append(response, vals.Index(x).Interface())
	}
	return response, nil
}
