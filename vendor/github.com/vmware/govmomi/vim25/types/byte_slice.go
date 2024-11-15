/*
Copyright (c) 2024-2024 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/vmware/govmomi/vim25/xml"
)

// ByteSlice implements vCenter compatibile xml encoding and decoding for a byte slice.
// vCenter encodes each byte of the array in its own xml element, whereas
// Go encodes the entire byte array in a single xml element.
type ByteSlice []byte

// MarshalXML implements xml.Marshaler
func (b ByteSlice) MarshalXML(e *xml.Encoder, field xml.StartElement) error {
	start := xml.StartElement{
		Name: field.Name,
	}
	for i := range b {
		if err := e.EncodeElement(b[i], start); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler
func (b *ByteSlice) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}

		if c, ok := t.(xml.CharData); ok {
			n, err := strconv.ParseInt(string(c), 10, 16)
			if err != nil {
				return err
			}
			if n > math.MaxUint8 {
				return fmt.Errorf("parsing %q: uint8 overflow", start.Name.Local)
			}
			*b = append(*b, byte(n))
		}
	}

	return nil
}
