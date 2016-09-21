package jsonx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"

	"github.com/Jeffail/gabs"
)

const (
	XMLHeader = `<?xml version="1.0" encoding="UTF-8"?>`
	Header    = `<json:object xsi:schemaLocation="http://www.datapower.com/schemas/json jsonx.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:json="http://www.ibm.com/xmlns/prod/2009/jsonx">`
	Footer    = `</json:object>`
)

// namedContainer wraps a gabs.Container to carry name information with it
type namedContainer struct {
	name string
	*gabs.Container
}

// Marshal marshals the input data into JSONx.
func Marshal(input interface{}) (string, error) {
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	xmlBytes, err := EncodeJSONBytes(jsonBytes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s%s", XMLHeader, Header, string(xmlBytes), Footer), nil
}

// EncodeJSONBytes encodes JSON-formatted bytes into JSONx. It is designed to
// be used for multiple entries so does not prepend the JSONx header tag or
// append the JSONx footer tag. You can use jsonx.Header and jsonx.Footer to
// easily add these when necessary.
func EncodeJSONBytes(input []byte) ([]byte, error) {
	o := bytes.NewBuffer(nil)
	reader := bytes.NewReader(input)
	dec := json.NewDecoder(reader)
	dec.UseNumber()

	cont, err := gabs.ParseJSONDecoder(dec)
	if err != nil {
		return nil, err
	}

	if err := sortAndTransformObject(o, &namedContainer{Container: cont}); err != nil {
		return nil, err
	}

	return o.Bytes(), nil
}

func transformContainer(o *bytes.Buffer, cont *namedContainer) error {
	var printName string

	if cont.name != "" {
		escapedNameBuf := bytes.NewBuffer(nil)
		err := xml.EscapeText(escapedNameBuf, []byte(cont.name))
		if err != nil {
			return err
		}
		printName = fmt.Sprintf(" name=\"%s\"", escapedNameBuf.String())
	}

	data := cont.Data()
	switch data.(type) {
	case nil:
		o.WriteString(fmt.Sprintf("<json:null%s />", printName))

	case bool:
		o.WriteString(fmt.Sprintf("<json:boolean%s>%t</json:boolean>", printName, data))

	case json.Number:
		o.WriteString(fmt.Sprintf("<json:number%s>%v</json:number>", printName, data))

	case string:
		o.WriteString(fmt.Sprintf("<json:string%s>%v</json:string>", printName, data))

	case []interface{}:
		o.WriteString(fmt.Sprintf("<json:array%s>", printName))
		arrayChildren, err := cont.Children()
		if err != nil {
			return err
		}
		for _, child := range arrayChildren {
			if err := transformContainer(o, &namedContainer{Container: child}); err != nil {
				return err
			}
		}
		o.WriteString("</json:array>")

	case map[string]interface{}:
		o.WriteString(fmt.Sprintf("<json:object%s>", printName))

		if err := sortAndTransformObject(o, cont); err != nil {
			return err
		}

		o.WriteString("</json:object>")
	}

	return nil
}

// sortAndTransformObject sorts object keys to make the output predictable so
// the package can be tested; logic is here to prevent code duplication
func sortAndTransformObject(o *bytes.Buffer, cont *namedContainer) error {
	objectChildren, err := cont.ChildrenMap()
	if err != nil {
		return err
	}

	sortedNames := make([]string, 0, len(objectChildren))
	for name, _ := range objectChildren {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)
	for _, name := range sortedNames {
		if err := transformContainer(o, &namedContainer{name: name, Container: objectChildren[name]}); err != nil {
			return err
		}
	}

	return nil
}
