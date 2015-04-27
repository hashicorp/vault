package aws

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// QueryClient is the underlying client for Query APIs.
type QueryClient struct {
	Context    Context
	Client     *http.Client
	Endpoint   string
	APIVersion string
}

// Do sends an HTTP request and returns an HTTP response, following policy
// (e.g. redirects, cookies, auth) as configured on the client.
func (c *QueryClient) Do(op, method, uri string, req, resp interface{}) error {
	body := url.Values{"Action": {op}, "Version": {c.APIVersion}}
	if err := c.loadValues(body, req, ""); err != nil {
		return err
	}

	httpReq, err := http.NewRequest(method, c.Endpoint+uri, strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("User-Agent", "aws-go")
	if err := c.Context.sign(httpReq); err != nil {
		return err
	}

	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return err
	}
	defer func() {
		_ = httpResp.Body.Close()
	}()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}
		if len(bodyBytes) == 0 {
			return APIError{
				StatusCode: httpResp.StatusCode,
				Message:    httpResp.Status,
			}
		}
		var queryErr queryErrorResponse
		if err := xml.Unmarshal(bodyBytes, &queryErr); err != nil {
			return err
		}
		return queryErr.Err(httpResp.StatusCode)
	}

	if resp != nil {
		return xml.NewDecoder(httpResp.Body).Decode(resp)
	}
	return nil
}

type queryErrorResponse struct {
	XMLName   xml.Name `xml:"ErrorResponse"`
	Type      string   `xml:"Error>Type"`
	Code      string   `xml:"Error>Code"`
	Message   string   `xml:"Error>Message"`
	RequestID string   `xml:"RequestId"`
}

func (e queryErrorResponse) Err(StatusCode int) error {
	return APIError{
		StatusCode: StatusCode,
		Type:       e.Type,
		Code:       e.Code,
		Message:    e.Message,
		RequestID:  e.RequestID,
	}
}

func (c *QueryClient) loadValues(v url.Values, i interface{}, prefix string) error {
	value := reflect.ValueOf(i)

	// follow any pointers
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// no need to handle zero values
	if !value.IsValid() {
		return nil
	}

	switch value.Kind() {
	case reflect.Struct:
		return c.loadStruct(v, value, prefix)
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			slicePrefix := prefix
			if slicePrefix == "" {
				slicePrefix = strconv.Itoa(i + 1)
			} else {
				slicePrefix = slicePrefix + "." + strconv.Itoa(i+1)
			}
			if err := c.loadValues(v, value.Index(i).Interface(), slicePrefix); err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		sortedKeys := []string{}
		keysByString := map[string]reflect.Value{}
		for _, k := range value.MapKeys() {
			s := fmt.Sprintf("%v", k.Interface())
			sortedKeys = append(sortedKeys, s)
			keysByString[s] = k
		}
		sort.Strings(sortedKeys)

		for i, sortKey := range sortedKeys {
			mapKey := keysByString[sortKey]

			var keyName string
			if prefix == "" {
				keyName = strconv.Itoa(i+1) + ".Name"
			} else {
				keyName = prefix + "." + strconv.Itoa(i+1) + ".Name"
			}

			if err := c.loadValue(v, mapKey, keyName); err != nil {
				return err
			}

			mapValue := value.MapIndex(mapKey)

			var valueName string
			if prefix == "" {
				valueName = strconv.Itoa(i+1) + ".Value"
			} else {
				valueName = prefix + "." + strconv.Itoa(i+1) + ".Value"
			}

			if err := c.loadValue(v, mapValue, valueName); err != nil {
				return err
			}
		}

		return nil
	default:
		panic("unknown request member type: " + value.String())
	}
}

func (c *QueryClient) loadStruct(v url.Values, value reflect.Value, prefix string) error {
	if !value.IsValid() {
		return nil
	}

	t := value.Type()
	for i := 0; i < value.NumField(); i++ {
		value := value.Field(i)
		name := t.Field(i).Tag.Get("query")
		if name == "" {
			name = t.Field(i).Name
		}
		if prefix != "" {
			name = prefix + "." + name
		}
		if err := c.loadValue(v, value, name); err != nil {
			return err
		}
	}
	return nil
}

func (c *QueryClient) loadValue(v url.Values, value reflect.Value, name string) error {
	switch casted := value.Interface().(type) {
	case string:
		if casted != "" {
			v.Set(name, casted)
		}
	case StringValue:
		if casted != nil {
			v.Set(name, *casted)
		}
	case BooleanValue:
		if casted != nil {
			v.Set(name, strconv.FormatBool(*casted))
		}
	case LongValue:
		if casted != nil {
			v.Set(name, strconv.FormatInt(*casted, 10))
		}
	case IntegerValue:
		if casted != nil {
			v.Set(name, strconv.Itoa(*casted))
		}
	case DoubleValue:
		if casted != nil {
			v.Set(name, strconv.FormatFloat(*casted, 'f', -1, 64))
		}
	case FloatValue:
		if casted != nil {
			v.Set(name, strconv.FormatFloat(float64(*casted), 'f', -1, 32))
		}
	case time.Time:
		if !casted.IsZero() {
			const ISO8601UTC = "2006-01-02T15:04:05Z"
			v.Set(name, casted.UTC().Format(ISO8601UTC))
		}
	case []string:
		if len(casted) != 0 {
			for i, val := range casted {
				v.Set(fmt.Sprintf("%s.%d", name, i+1), val)
			}
		}
	default:
		if err := c.loadValues(v, value.Interface(), name); err != nil {
			return err
		}
	}
	return nil
}
