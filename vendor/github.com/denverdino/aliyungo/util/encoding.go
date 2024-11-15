package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// change instance=["a", "b"]
// to instance.1="a" instance.2="b"
func FlattenFn(fieldName string, field reflect.Value, values *url.Values) {
	l := field.Len()
	if l > 0 {
		for i := 0; i < l; i++ {
			str := field.Index(i).String()
			values.Set(fieldName+"."+strconv.Itoa(i+1), str)
		}
	}
}

func Underline2Dot(name string) string {
	return strings.Replace(name, "_", ".", -1)
}

//ConvertToQueryValues converts the struct to url.Values
func ConvertToQueryValues(ifc interface{}) url.Values {
	values := url.Values{}
	SetQueryValues(ifc, &values)
	return values
}

//SetQueryValues sets the struct to existing url.Values following ECS encoding rules
func SetQueryValues(ifc interface{}, values *url.Values) {
	setQueryValues(ifc, values, "")
}

func SetQueryValueByFlattenMethod(ifc interface{}, values *url.Values) {
	setQueryValuesByFlattenMethod(ifc, values, "")
}

func setQueryValues(i interface{}, values *url.Values, prefix string) {
	// add to support url.Values
	mapValues, ok := i.(url.Values)
	if ok {
		for k := range mapValues {
			values.Set(k, mapValues.Get(k))
		}
		return
	}

	elem := reflect.ValueOf(i)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	elemType := elem.Type()
	for i := 0; i < elem.NumField(); i++ {

		fieldName := elemType.Field(i).Name
		anonymous := elemType.Field(i).Anonymous
		tag := elemType.Field(i).Tag.Get("query")
		argName := elemType.Field(i).Tag.Get("ArgName")
		field := elem.Field(i)
		// TODO Use Tag for validation
		// tag := typ.Field(i).Tag.Get("tagname")
		kind := field.Kind()
		isPtr := false
		if (kind == reflect.Ptr || kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map || kind == reflect.Chan) && field.IsNil() {
			continue
		}
		if kind == reflect.Ptr {
			field = field.Elem()
			kind = field.Kind()
			isPtr = true
		}
		var value string
		//switch field.Interface().(type) {
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i := field.Int()
			if i != 0 || isPtr {
				value = strconv.FormatInt(i, 10)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i := field.Uint()
			if i != 0 || isPtr {
				value = strconv.FormatUint(i, 10)
			}
		case reflect.Float32:
			value = strconv.FormatFloat(field.Float(), 'f', 4, 32)
		case reflect.Float64:
			value = strconv.FormatFloat(field.Float(), 'f', 4, 64)
		case reflect.Bool:
			value = strconv.FormatBool(field.Bool())
		case reflect.String:
			value = field.String()
		case reflect.Map:
			ifc := field.Interface()
			m := ifc.(map[string]string)
			if m != nil {
				j := 0
				for k, v := range m {
					j++
					keyName := fmt.Sprintf("%s.%d.Key", fieldName, j)
					values.Set(keyName, k)
					valueName := fmt.Sprintf("%s.%d.Value", fieldName, j)
					values.Set(valueName, v)
				}
			}
		case reflect.Slice:
			switch field.Type().Elem().Kind() {
			case reflect.Uint8:
				value = string(field.Bytes())
			case reflect.String:
				l := field.Len()
				if l > 0 {
					if tag == "list" {
						name := argName
						if argName == "" {
							name = fieldName
						}
						for i := 0; i < l; i++ {
							valueName := fmt.Sprintf("%s.%d", name, (i + 1))
							values.Set(valueName, field.Index(i).String())
						}
					} else {
						strArray := make([]string, l)
						for i := 0; i < l; i++ {
							strArray[i] = field.Index(i).String()
						}
						bytes, err := json.Marshal(strArray)
						if err == nil {
							value = string(bytes)
						} else {
							log.Printf("Failed to convert JSON: %v", err)
						}
					}
				}
			default:
				l := field.Len()
				for j := 0; j < l; j++ {
					prefixName := fmt.Sprintf("%s.%d.", fieldName, (j + 1))
					ifc := field.Index(j).Interface()
					//log.Printf("%s : %v", prefixName, ifc)
					if ifc != nil {
						setQueryValues(ifc, values, prefixName)
					}
				}
				continue
			}

		default:
			switch field.Interface().(type) {
			case ISO6801Time:
				t := field.Interface().(ISO6801Time)
				value = t.String()
			case time.Time:
				t := field.Interface().(time.Time)
				value = GetISO8601TimeStamp(t)
			default:
				ifc := field.Interface()
				if ifc != nil {
					if anonymous {
						SetQueryValues(ifc, values)
					} else {
						prefixName := fieldName + "."
						setQueryValues(ifc, values, prefixName)
					}
					continue
				}
			}
		}
		if value != "" {
			name := argName
			if argName == "" {
				name = fieldName
			}
			if prefix != "" {
				name = prefix + name
			}
			values.Set(name, value)
		}
	}
}

func setQueryValuesByFlattenMethod(i interface{}, values *url.Values, prefix string) {
	// add to support url.Values
	mapValues, ok := i.(url.Values)
	if ok {
		for k := range mapValues {
			values.Set(k, mapValues.Get(k))
		}
		return
	}

	elem := reflect.ValueOf(i)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	elemType := elem.Type()
	for i := 0; i < elem.NumField(); i++ {

		fieldName := elemType.Field(i).Name
		anonymous := elemType.Field(i).Anonymous
		field := elem.Field(i)

		// TODO Use Tag for validation
		// tag := typ.Field(i).Tag.Get("tagname")
		kind := field.Kind()

		isPtr := false
		if (kind == reflect.Ptr || kind == reflect.Array || kind == reflect.Slice || kind == reflect.Map || kind == reflect.Chan) && field.IsNil() {
			continue
		}
		if kind == reflect.Ptr {
			field = field.Elem()
			kind = field.Kind()
			isPtr = true
		}

		var value string
		//switch field.Interface().(type) {
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i := field.Int()
			if i != 0 || isPtr {
				value = strconv.FormatInt(i, 10)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i := field.Uint()
			if i != 0 || isPtr {
				value = strconv.FormatUint(i, 10)
			}
		case reflect.Float32:
			value = strconv.FormatFloat(field.Float(), 'f', 4, 32)
		case reflect.Float64:
			value = strconv.FormatFloat(field.Float(), 'f', 4, 64)
		case reflect.Bool:
			value = strconv.FormatBool(field.Bool())
		case reflect.String:
			value = field.String()
		case reflect.Map:
			ifc := field.Interface()
			m := ifc.(map[string]string)
			if m != nil {
				j := 0
				for k, v := range m {
					j++
					keyName := fmt.Sprintf("%s.%d.Key", fieldName, j)
					values.Set(keyName, k)
					valueName := fmt.Sprintf("%s.%d.Value", fieldName, j)
					values.Set(valueName, v)
				}
			}
		case reflect.Slice:
			if field.Type().Name() == "FlattenArray" {
				FlattenFn(fieldName, field, values)
			} else {
				switch field.Type().Elem().Kind() {
				case reflect.Uint8:
					value = string(field.Bytes())
				case reflect.String:
					l := field.Len()
					if l > 0 {
						strArray := make([]string, l)
						for i := 0; i < l; i++ {
							strArray[i] = field.Index(i).String()
						}
						bytes, err := json.Marshal(strArray)
						if err == nil {
							value = string(bytes)
						} else {
							log.Printf("Failed to convert JSON: %v", err)
						}
					}
				default:
					l := field.Len()
					for j := 0; j < l; j++ {
						prefixName := fmt.Sprintf("%s.%d.", fieldName, (j + 1))
						ifc := field.Index(j).Interface()
						//log.Printf("%s : %v", prefixName, ifc)
						if ifc != nil {
							setQueryValuesByFlattenMethod(ifc, values, prefixName)
						}
					}
					continue
				}
			}

		default:
			switch field.Interface().(type) {
			case ISO6801Time:
				t := field.Interface().(ISO6801Time)
				value = t.String()
			case time.Time:
				t := field.Interface().(time.Time)
				value = GetISO8601TimeStamp(t)
			default:

				ifc := field.Interface()
				if ifc != nil {
					if anonymous {
						SetQueryValues(ifc, values)
					} else {
						prefixName := fieldName + "."
						setQueryValuesByFlattenMethod(ifc, values, prefixName)
					}
					continue
				}
			}
		}
		if value != "" {
			name := elemType.Field(i).Tag.Get("ArgName")
			if name == "" {
				name = fieldName
			}
			if prefix != "" {
				name = prefix + name
			}
			// NOTE: here we will change name to underline style when the type is UnderlineString
			if field.Type().Name() == "UnderlineString" {
				name = Underline2Dot(name)
			}
			values.Set(name, value)
		}
	}
}
