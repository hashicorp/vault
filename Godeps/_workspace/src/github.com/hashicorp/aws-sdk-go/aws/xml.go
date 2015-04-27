package aws

import (
	"encoding/xml"
	"reflect"
	"strings"
)

// MarshalXML is a weird and stunted version of xml.Marshal which is used by the
// REST-XML request types to get around a bug in encoding/xml which doesn't
// allow us to marshal pointers to zero values:
//
// https://github.com/golang/go/issues/5452
func MarshalXML(v interface{}, e *xml.Encoder, start xml.StartElement) error {
	value := reflect.ValueOf(v)
	t := value.Type()
	switch value.Kind() {
	case reflect.Ptr:
		if !value.IsNil() {
			return MarshalXML(value.Elem().Interface(), e, start)
		}
	case reflect.Struct:
		var rootInfo xmlFieldInfo

		// detect xml.Name, if any
		for i := 0; i < value.NumField(); i++ {
			f := t.Field(i)
			v := value.Field(i)
			if f.Type == xmlName {
				rootInfo = parseXMLTag(f.Tag.Get("xml"))
				if rootInfo.name == "" {
					// name not in tag, try value
					name := v.Interface().(xml.Name)
					rootInfo = xmlFieldInfo{
						name: name.Local,
						ns:   name.Space,
					}
				}
			}
		}

		for _, start := range rootInfo.start(t.Name()) {
			if err := e.EncodeToken(start); err != nil {
				return err
			}
		}

		for i := 0; i < value.NumField(); i++ {
			ft := value.Type().Field(i)

			if ft.Type == xmlName {
				continue
			}

			fv := value.Field(i)
			fi := parseXMLTag(ft.Tag.Get("xml"))

			if fi.name == "-" {
				continue
			}

			if fi.omit {
				switch fv.Kind() {
				case reflect.Ptr:
					if fv.IsNil() {
						continue
					}
				case reflect.Slice, reflect.Map:
					if fv.Len() == 0 {
						continue
					}
				default:
					if !fv.IsValid() {
						continue
					}
				}
			}

			starts := fi.start(ft.Name)
			for _, start := range starts[:len(starts)-1] {
				if err := e.EncodeToken(start); err != nil {
					return err
				}
			}

			start := starts[len(starts)-1]
			if err := e.EncodeElement(fv.Interface(), start); err != nil {
				return err
			}

			for _, end := range fi.end(ft.Name)[1:] {
				if err := e.EncodeToken(end); err != nil {
					return err
				}
			}
		}

		for _, end := range rootInfo.end(t.Name()) {
			if err := e.EncodeToken(end); err != nil {
				return err
			}
		}
	default:
		return e.Encode(v)
	}
	return nil
}

var xmlName = reflect.TypeOf(xml.Name{})

type xmlFieldInfo struct {
	name string
	ns   string
	omit bool
}

func (fi xmlFieldInfo) start(name string) []xml.StartElement {
	if fi.name != "" {
		name = fi.name
	}

	var elements []xml.StartElement
	for _, part := range strings.Split(name, ">") {
		elements = append(elements, xml.StartElement{
			Name: xml.Name{
				Local: part,
				Space: fi.ns,
			},
		})
	}
	return elements
}

func (fi xmlFieldInfo) end(name string) []xml.EndElement {
	if fi.name != "" {
		name = fi.name
	}

	var elements []xml.EndElement
	parts := strings.Split(name, ">")
	for i := range parts {
		part := parts[len(parts)-i-1]
		elements = append(elements, xml.EndElement{
			Name: xml.Name{
				Local: part,
				Space: fi.ns,
			},
		})
	}
	return elements
}

func parseXMLTag(t string) xmlFieldInfo {
	parts := strings.Split(t, ",")

	var omit bool
	for _, p := range parts {
		omit = omit || p == "omitempty"
	}

	var name, ns string
	if len(parts) > 0 {
		nameParts := strings.Split(parts[0], " ")
		if len(nameParts) == 2 {
			name = nameParts[1]
			ns = nameParts[0]
		} else if len(nameParts) == 1 {
			name = nameParts[0]
		}

	}

	return xmlFieldInfo{
		name: name,
		ns:   ns,
		omit: omit,
	}
}
