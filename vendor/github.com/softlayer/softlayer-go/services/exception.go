/**
 * Copyright 2016 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * AUTOMATICALLY GENERATED CODE - DO NOT MODIFY
 */

package services

import (
	"fmt"
	"strings"

	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

// Throw this exception if there are validation errors. The types are specified in SoftLayer_Brand_Creation_Input including: KEY_NAME, PREFIX, NAME, LONG_NAME, SUPPORT_POLICY, POLICY_ACKNOWLEDGEMENT_FLAG, etc.
type Exception_Brand_Creation struct {
	Session *session.Session
	Options sl.Options
}

// GetExceptionBrandCreationService returns an instance of the Exception_Brand_Creation SoftLayer service
func GetExceptionBrandCreationService(sess *session.Session) Exception_Brand_Creation {
	return Exception_Brand_Creation{Session: sess}
}

func (r Exception_Brand_Creation) Id(id int) Exception_Brand_Creation {
	r.Options.Id = &id
	return r
}

func (r Exception_Brand_Creation) Mask(mask string) Exception_Brand_Creation {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Exception_Brand_Creation) Filter(filter string) Exception_Brand_Creation {
	r.Options.Filter = filter
	return r
}

func (r Exception_Brand_Creation) Limit(limit int) Exception_Brand_Creation {
	r.Options.Limit = &limit
	return r
}

func (r Exception_Brand_Creation) Offset(offset int) Exception_Brand_Creation {
	r.Options.Offset = &offset
	return r
}
