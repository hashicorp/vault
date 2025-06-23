// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package json

import (
	"fmt"
	"strings"
)

func GenerateComplexJSON(depth int) string {
	var innerBuilder strings.Builder
	innerBuilder.WriteString("{")
	i := 0
	prefixes := []string{"", "_", ".", ","}
	for _, prefix := range prefixes {
		for x := 32; x < 126; x++ {
			if x == 34 || x == 92 {
				continue
			}
			innerBuilder.WriteString(fmt.Sprintf(`"%s%c":%d,`, prefix, rune(x), i))
			i++
		}
	}
	innerBuilder.WriteString(`"~":`)

	inner := innerBuilder.String()

	var jsonBuilder strings.Builder
	jsonBuilder.WriteString(`{"data":{"data":`)

	for k := 0; k < depth; k++ {
		jsonBuilder.WriteString(inner)
	}
	jsonBuilder.WriteString(`0.1`)
	for k := 0; k < depth; k++ {
		jsonBuilder.WriteString(`}`)
	}
	jsonBuilder.WriteString(`}}`)

	return jsonBuilder.String()
}
