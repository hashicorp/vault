// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cloudsql

import (
	"context"
	"fmt"
	"net"
	"sort"

	"cloud.google.com/go/cloudsqlconn/instance"
)

// DNSResolver uses the default net.Resolver to find
// TXT records containing an instance name for a DNS record.
var DNSResolver = &DNSInstanceConnectionNameResolver{
	dnsResolver: net.DefaultResolver,
}

// DefaultResolver simply parses instance names.
var DefaultResolver = &ConnNameResolver{}

// ConnNameResolver simply parses instance names. Implements
// InstanceConnectionNameResolver
type ConnNameResolver struct {
}

// Resolve returns the instance name, possibly using DNS. This will return an
// instance.ConnName or an error if it was unable to resolve an instance name.
func (r *ConnNameResolver) Resolve(_ context.Context, icn string) (instanceName instance.ConnName, err error) {
	return instance.ParseConnName(icn)
}

// netResolver groups the methods on net.Resolver that are used by the DNS
// resolver implementation. This allows an application to replace the default
// net.DefaultResolver with a custom implementation. For example: the
// application may need to connect to a specific DNS server using a specially
// configured instance of net.Resolver.
type netResolver interface {
	LookupTXT(ctx context.Context, name string) ([]string, error)
}

// DNSInstanceConnectionNameResolver can resolve domain names into instance names using
// TXT records in DNS. Implements InstanceConnectionNameResolver
type DNSInstanceConnectionNameResolver struct {
	dnsResolver netResolver
}

// Resolve returns the instance name, possibly using DNS. This will return an
// instance.ConnName or an error if it was unable to resolve an instance name.
func (r *DNSInstanceConnectionNameResolver) Resolve(ctx context.Context, icn string) (instanceName instance.ConnName, err error) {
	cn, err := instance.ParseConnName(icn)
	if err != nil {
		// The connection name was not project:region:instance
		// Attempt to query a TXT record and see if it works instead.
		cn, err = r.queryDNS(ctx, icn)
		if err != nil {
			return instance.ConnName{}, err
		}
	}

	return cn, nil
}

// queryDNS attempts to resolve a TXT record for the domain name.
// The DNS TXT record's target field is used as instance name.
//
// This handles several conditions where the DNS records may be missing or
// invalid:
//   - The domain name resolves to 0 DNS records - return an error
//   - Some DNS records to not contain a well-formed instance name - return the
//     first well-formed instance name. If none found return an error.
//   - The domain name resolves to 2 or more DNS record - return first valid
//     record when sorted by priority: lowest value first, then by target:
//     alphabetically.
func (r *DNSInstanceConnectionNameResolver) queryDNS(ctx context.Context, domainName string) (instance.ConnName, error) {
	// Attempt to query the TXT records.
	// This could return a partial error where both err != nil && len(records) > 0.
	records, err := r.dnsResolver.LookupTXT(ctx, domainName)
	// If resolve failed and no records were found, return the error.
	if err != nil {
		return instance.ConnName{}, fmt.Errorf("unable to resolve TXT record for %q: %v", domainName, err)
	}

	// Process the records returning the first valid TXT record.

	// Sort the TXT record values alphabetically by instance name
	sort.Slice(records, func(i, j int) bool {
		return records[i] < records[j]
	})

	var perr error
	// Attempt to parse records, returning the first valid record.
	for _, record := range records {
		// Parse the target as a CN
		cn, parseErr := instance.ParseConnName(record)
		if parseErr != nil {
			perr = fmt.Errorf("unable to parse TXT for %q -> %q : %v", domainName, record, parseErr)
			continue
		}
		return cn, nil
	}

	// If all the records failed to parse, return one of the parse errors
	if perr != nil {
		return instance.ConnName{}, perr
	}

	// No records were found, return an error.
	return instance.ConnName{}, fmt.Errorf("no valid TXT records found for %q", domainName)
}
