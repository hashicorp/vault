// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.

package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// String returns a pointer to the provided string
func String(value string) *string {
	return &value
}

// Int returns a pointer to the provided int
func Int(value int) *int {
	return &value
}

// Int64 returns a pointer to the provided int64
func Int64(value int64) *int64 {
	return &value
}

// Uint returns a pointer to the provided uint
func Uint(value uint) *uint {
	return &value
}

//Float32 returns a pointer to the provided float32
func Float32(value float32) *float32 {
	return &value
}

//Float64 returns a pointer to the provided float64
func Float64(value float64) *float64 {
	return &value
}

//Bool returns a pointer to the provided bool
func Bool(value bool) *bool {
	return &value
}

//PointerString prints the values of pointers in a struct
//Producing a human friendly string for an struct with pointers.
//useful when debugging the values of a struct
func PointerString(datastruct interface{}) (representation string) {
	val := reflect.ValueOf(datastruct)
	typ := reflect.TypeOf(datastruct)
	all := make([]string, 2)
	all = append(all, "{")
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)

		//unexported
		if sf.PkgPath != "" && !sf.Anonymous {
			continue
		}

		sv := val.Field(i)
		stringValue := ""
		if isNil(sv) {
			stringValue = fmt.Sprintf("%s=<nil>", sf.Name)
		} else {
			if sv.Type().Kind() == reflect.Ptr {
				sv = sv.Elem()
			}
			stringValue = fmt.Sprintf("%s=%v", sf.Name, sv)
		}
		all = append(all, stringValue)
	}
	all = append(all, "}")
	representation = strings.TrimSpace(strings.Join(all, " "))
	return
}

// SDKTime a struct that parses/renders to/from json using RFC339 date-time information
type SDKTime struct {
	time.Time
}

// SDKDate a struct that parses/renders to/from json using only date information
type SDKDate struct {
	//Date date information
	Date time.Time
}

func sdkTimeFromTime(t time.Time) SDKTime {
	return SDKTime{t}
}

func sdkDateFromTime(t time.Time) SDKDate {
	return SDKDate{Date: t}
}

func formatTime(t SDKTime) string {
	return t.Format(sdkTimeFormat)
}

func formatDate(t SDKDate) string {
	return t.Date.Format(sdkDateFormat)
}

func now() *SDKTime {
	t := SDKTime{time.Now()}
	return &t
}

var timeType = reflect.TypeOf(SDKTime{})
var timeTypePtr = reflect.TypeOf(&SDKTime{})

var sdkDateType = reflect.TypeOf(SDKDate{})
var sdkDateTypePtr = reflect.TypeOf(&SDKDate{})

//Formats for sdk supported time representations
const sdkTimeFormat = time.RFC3339Nano
const rfc1123OptionalLeadingDigitsInDay = "Mon, _2 Jan 2006 15:04:05 MST"
const sdkDateFormat = "2006-01-02"

func tryParsingTimeWithValidFormatsForHeaders(data []byte, headerName string) (t time.Time, err error) {
	header := strings.ToLower(headerName)
	switch header {
	case "lastmodified", "date":
		t, err = tryParsing(data, time.RFC3339Nano, time.RFC3339, time.RFC1123, rfc1123OptionalLeadingDigitsInDay, time.RFC850, time.ANSIC)
		return
	default: //By default we parse with RFC3339
		t, err = time.Parse(sdkTimeFormat, string(data))
		return
	}
}

func tryParsing(data []byte, layouts ...string) (tm time.Time, err error) {
	datestring := string(data)
	for _, l := range layouts {
		tm, err = time.Parse(l, datestring)
		if err == nil {
			return
		}
	}
	err = fmt.Errorf("Could not parse time: %s with formats: %s", datestring, layouts[:])
	return
}

// String returns string representation of SDKDate
func (t *SDKDate) String() string {
	return t.Date.Format(sdkDateFormat)
}

// NewSDKDateFromString parses the dateString into SDKDate
func NewSDKDateFromString(dateString string) (*SDKDate, error) {
	parsedTime, err := time.Parse(sdkDateFormat, dateString)
	if err != nil {
		return nil, err
	}

	return &SDKDate{Date: parsedTime}, nil
}

// UnmarshalJSON unmarshals from json
func (t *SDKTime) UnmarshalJSON(data []byte) (e error) {
	s := string(data)
	if s == "null" {
		t.Time = time.Time{}
	} else {
		//Try parsing with RFC3339
		t.Time, e = time.Parse(`"`+sdkTimeFormat+`"`, string(data))
	}
	return
}

// MarshalJSON marshals to JSON
func (t *SDKTime) MarshalJSON() (buff []byte, e error) {
	s := t.Format(sdkTimeFormat)
	buff = []byte(`"` + s + `"`)
	return
}

// UnmarshalJSON unmarshals from json
func (t *SDKDate) UnmarshalJSON(data []byte) (e error) {
	if string(data) == `"null"` {
		t.Date = time.Time{}
		return
	}

	t.Date, e = tryParsing(data,
		strconv.Quote(sdkDateFormat),
	)
	return
}

// MarshalJSON marshals to JSON
func (t *SDKDate) MarshalJSON() (buff []byte, e error) {
	s := t.Date.Format(sdkDateFormat)
	buff = []byte(strconv.Quote(s))
	return
}

// PrivateKeyFromBytes is a helper function that will produce a RSA private
// key from bytes. This function is deprecated in favour of PrivateKeyFromBytesWithPassword
// Deprecated
func PrivateKeyFromBytes(pemData []byte, password *string) (key *rsa.PrivateKey, e error) {
	if password == nil {
		return PrivateKeyFromBytesWithPassword(pemData, nil)
	}

	return PrivateKeyFromBytesWithPassword(pemData, []byte(*password))
}

// PrivateKeyFromBytesWithPassword is a helper function that will produce a RSA private
// key from bytes and a password.
func PrivateKeyFromBytesWithPassword(pemData, password []byte) (key *rsa.PrivateKey, e error) {
	if pemBlock, _ := pem.Decode(pemData); pemBlock != nil {
		decrypted := pemBlock.Bytes
		if x509.IsEncryptedPEMBlock(pemBlock) {
			if password == nil {
				e = fmt.Errorf("private key password is required for encrypted private keys")
				return
			}
			if decrypted, e = x509.DecryptPEMBlock(pemBlock, password); e != nil {
				return
			}
		}

		key, e = parsePKCSPrivateKey(decrypted)

	} else {
		e = fmt.Errorf("PEM data was not found in buffer")
		return
	}
	return
}

// ParsePrivateKey using PKCS1 or PKCS8
func parsePKCSPrivateKey(decryptedKey []byte) (*rsa.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(decryptedKey); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(decryptedKey); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey:
			return key, nil
		default:
			return nil, fmt.Errorf("unsupportesd private key type in PKCS8 wrapping")
		}
	}
	return nil, fmt.Errorf("failed to parse private key")
}

func generateRandUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid, nil
}

func makeACopy(original []string) []string {
	tmp := make([]string, len(original))
	copy(tmp, original)
	return tmp
}
