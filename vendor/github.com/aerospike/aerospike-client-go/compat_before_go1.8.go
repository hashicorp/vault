// +build !go1.8

// Copyright 2013-2020 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"crypto/tls"
)

func cloneTLSConfig(c *tls.Config) *tls.Config {
	// .Clone() method is not available in go versions before 1.8
	return &tls.Config{
		Certificates:             c.Certificates,
		CipherSuites:             c.CipherSuites,
		ClientAuth:               c.ClientAuth,
		ClientCAs:                c.ClientCAs,
		ClientSessionCache:       c.ClientSessionCache,
		CurvePreferences:         c.CurvePreferences,
		GetCertificate:           c.GetCertificate,
		InsecureSkipVerify:       c.InsecureSkipVerify,
		MaxVersion:               c.MaxVersion,
		MinVersion:               c.MinVersion,
		NameToCertificate:        c.NameToCertificate,
		NextProtos:               c.NextProtos,
		PreferServerCipherSuites: c.PreferServerCipherSuites,
		Rand:                     c.Rand,
		RootCAs:                  c.RootCAs,
		ServerName:               c.ServerName,
		Time:                     c.Time,
	}
}
