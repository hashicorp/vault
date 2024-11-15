/**
 * Copyright 2016-2024 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed
 * on an "AS IS" BASIS,WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

// AUTOMATICALLY GENERATED CODE - DO NOT MODIFY

package services

import (
	"fmt"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

// The SoftLayer_Dns_Domain data type represents a single DNS domain record hosted on the SoftLayer nameservers. Domains contain general information about the domain name such as name and serial. Individual records such as A, AAAA, CTYPE, and MX records are stored in the domain's associated [[SoftLayer_Dns_Domain_ResourceRecord (type)|SoftLayer_Dns_Domain_ResourceRecord]] records.
type Dns_Domain struct {
	Session session.SLSession
	Options sl.Options
}

// GetDnsDomainService returns an instance of the Dns_Domain SoftLayer service
func GetDnsDomainService(sess session.SLSession) Dns_Domain {
	return Dns_Domain{Session: sess}
}

func (r Dns_Domain) Id(id int) Dns_Domain {
	r.Options.Id = &id
	return r
}

func (r Dns_Domain) Mask(mask string) Dns_Domain {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Dns_Domain) Filter(filter string) Dns_Domain {
	r.Options.Filter = filter
	return r
}

func (r Dns_Domain) Limit(limit int) Dns_Domain {
	r.Options.Limit = &limit
	return r
}

func (r Dns_Domain) Offset(offset int) Dns_Domain {
	r.Options.Offset = &offset
	return r
}

// Create an A record on a SoftLayer domain. This is a shortcut method, meant to take the work out of creating a SoftLayer_Dns_Domain_ResourceRecord if you already have a domain record available. createARecord returns the newly created SoftLayer_Dns_Domain_ResourceRecord_AType.
func (r Dns_Domain) CreateARecord(host *string, data *string, ttl *int) (resp datatypes.Dns_Domain_ResourceRecord_AType, err error) {
	params := []interface{}{
		host,
		data,
		ttl,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "createARecord", params, &r.Options, &resp)
	return
}

// Create an AAAA record on a SoftLayer domain. This is a shortcut method, meant to take the work out of creating a SoftLayer_Dns_Domain_ResourceRecord if you already have a domain record available. createARecord returns the newly created SoftLayer_Dns_Domain_ResourceRecord_AaaaType.
func (r Dns_Domain) CreateAaaaRecord(host *string, data *string, ttl *int) (resp datatypes.Dns_Domain_ResourceRecord_AaaaType, err error) {
	params := []interface{}{
		host,
		data,
		ttl,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "createAaaaRecord", params, &r.Options, &resp)
	return
}

// Create a CNAME record on a SoftLayer domain. This is a shortcut method, meant to take the work out of creating a SoftLayer_Dns_Domain_ResourceRecord if you already have a domain record available. createCnameRecord returns the newly created SoftLayer_Dns_Domain_ResourceRecord_CnameType.
func (r Dns_Domain) CreateCnameRecord(host *string, data *string, ttl *int) (resp datatypes.Dns_Domain_ResourceRecord_CnameType, err error) {
	params := []interface{}{
		host,
		data,
		ttl,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "createCnameRecord", params, &r.Options, &resp)
	return
}

// Create an MX record on a SoftLayer domain. This is a shortcut method, meant to take the work out of creating a SoftLayer_Dns_Domain_ResourceRecord if you already have a domain record available. MX records are created with a default priority of 10. createMxRecord returns the newly created SoftLayer_Dns_Domain_ResourceRecord_MxType.
func (r Dns_Domain) CreateMxRecord(host *string, data *string, ttl *int, mxPriority *int) (resp datatypes.Dns_Domain_ResourceRecord_MxType, err error) {
	params := []interface{}{
		host,
		data,
		ttl,
		mxPriority,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "createMxRecord", params, &r.Options, &resp)
	return
}

// Create an NS record on a SoftLayer domain. This is a shortcut method, meant to take the work out of creating a SoftLayer_Dns_Domain_ResourceRecord if you already have a domain record available. createNsRecord returns the newly created SoftLayer_Dns_Domain_ResourceRecord_NsType.
func (r Dns_Domain) CreateNsRecord(host *string, data *string, ttl *int) (resp datatypes.Dns_Domain_ResourceRecord_NsType, err error) {
	params := []interface{}{
		host,
		data,
		ttl,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "createNsRecord", params, &r.Options, &resp)
	return
}

// Create a new domain on the SoftLayer name servers. The SoftLayer_Dns_Domain object passed to this function must have at least one A or AAAA resource record.
//
// createObject creates a default SOA record with the data:
// * ”'host”': "@"
// * ”'data”': "ns1.softlayer.com."
// * ”'responsible person”': "root.[your domain name]."
// * ”'expire”': 604800 seconds
// * ”'refresh”': 3600 seconds
// * ”'retry”': 300 seconds
// * ”'minimum”': 3600 seconds
//
// If your new domain uses the .de top-level domain then SOA refresh is set to 10000 seconds, retry is set to 1800 seconds, and minimum to 10000 seconds.
//
// If your domain doesn't contain NS resource records for ns1.softlayer.com or ns2.softlayer.com then ”createObject” will create them for you.
//
// ”createObject” returns a Boolean ”true” on successful object creation or ”false” if your domain was unable to be created..
func (r Dns_Domain) CreateObject(templateObject *datatypes.Dns_Domain) (resp datatypes.Dns_Domain, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "createObject", params, &r.Options, &resp)
	return
}

// Create multiple domains on the SoftLayer name servers. Each domain record passed to ”createObjects” follows the logic in the SoftLayer_Dns_Domain ”createObject” method.
func (r Dns_Domain) CreateObjects(templateObjects []datatypes.Dns_Domain) (resp []datatypes.Dns_Domain, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "createObjects", params, &r.Options, &resp)
	return
}

// setPtrRecordForIpAddress() sets a single reverse DNS record for a single IP address and returns the newly created or edited [[SoftLayer_Dns_Domain_ResourceRecord]] record. Currently this method only supports IPv4 addresses and performs no operation when given an IPv6 address.
func (r Dns_Domain) CreatePtrRecord(ipAddress *string, ptrRecord *string, ttl *int) (resp datatypes.Dns_Domain_ResourceRecord, err error) {
	params := []interface{}{
		ipAddress,
		ptrRecord,
		ttl,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "createPtrRecord", params, &r.Options, &resp)
	return
}

// Create an SPF record on a SoftLayer domain. This is a shortcut method, meant to take the work out of creating a SoftLayer_Dns_Domain_ResourceRecord if you already have a domain record available. createARecord returns the newly created SoftLayer_Dns_Domain_ResourceRecord_SpfType.
func (r Dns_Domain) CreateSpfRecord(host *string, data *string, ttl *int) (resp datatypes.Dns_Domain_ResourceRecord_SpfType, err error) {
	params := []interface{}{
		host,
		data,
		ttl,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "createSpfRecord", params, &r.Options, &resp)
	return
}

// Create a TXT record on a SoftLayer domain. This is a shortcut method, meant to take the work out of creating a SoftLayer_Dns_Domain_ResourceRecord if you already have a domain record available. createARecord returns the newly created SoftLayer_Dns_Domain_ResourceRecord_TxtType.
func (r Dns_Domain) CreateTxtRecord(host *string, data *string, ttl *int) (resp datatypes.Dns_Domain_ResourceRecord_TxtType, err error) {
	params := []interface{}{
		host,
		data,
		ttl,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "createTxtRecord", params, &r.Options, &resp)
	return
}

// deleteObject permanently removes a domain and all of it's associated resource records from the softlayer name servers. ”'This cannot be undone.”' Be wary of running this method. If you remove a domain in error you will need to re-create it by creating a new SoftLayer_Dns_Domain object.
func (r Dns_Domain) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "deleteObject", nil, &r.Options, &resp)
	return
}

// Retrieve The SoftLayer customer account that owns a domain.
func (r Dns_Domain) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "getAccount", nil, &r.Options, &resp)
	return
}

// Search for [[SoftLayer_Dns_Domain]] records by domain name. getByDomainName() performs an inclusive search for domain records, returning multiple records based on partial name matches. Use this method to locate domain records if you don't have access to their id numbers.
func (r Dns_Domain) GetByDomainName(name *string) (resp []datatypes.Dns_Domain, err error) {
	params := []interface{}{
		name,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "getByDomainName", params, &r.Options, &resp)
	return
}

// Retrieve A flag indicating that the dns domain record is a managed resource.
func (r Dns_Domain) GetManagedResourceFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "getManagedResourceFlag", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_Dns_Domain object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_Dns_Domain service. You can only retrieve domains that are assigned to your SoftLayer account.
func (r Dns_Domain) GetObject() (resp datatypes.Dns_Domain, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The individual records contained within a domain record. These include but are not limited to A, AAAA, MX, CTYPE, SPF and TXT records.
func (r Dns_Domain) GetResourceRecords() (resp []datatypes.Dns_Domain_ResourceRecord, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "getResourceRecords", nil, &r.Options, &resp)
	return
}

// Retrieve The secondary DNS record that defines this domain as being managed through zone transfers.
func (r Dns_Domain) GetSecondary() (resp datatypes.Dns_Secondary, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "getSecondary", nil, &r.Options, &resp)
	return
}

// Retrieve The start of authority (SOA) record contains authoritative and propagation details for a DNS zone. This property is not considered in requests to createObject and editObject.
func (r Dns_Domain) GetSoaResourceRecord() (resp datatypes.Dns_Domain_ResourceRecord_SoaType, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "getSoaResourceRecord", nil, &r.Options, &resp)
	return
}

// Return a SoftLayer hosted domain and resource records' data formatted as zone file.
func (r Dns_Domain) GetZoneFileContents() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain", "getZoneFileContents", nil, &r.Options, &resp)
	return
}

// The SoftLayer_Dns_Domain_ResourceRecord data type represents a single resource record entry in a SoftLayer hosted domain. Each resource record contains a ”host” and ”data” property, defining a resource's name and it's target data. Domains contain multiple types of resource records. The ”type” property separates out resource records by type. ”Type” can take one of the following values:
// * ”'"a"”' for [[SoftLayer_Dns_Domain_ResourceRecord_AType|address]] records
// * ”'"aaaa"”' for [[SoftLayer_Dns_Domain_ResourceRecord_AaaaType|address]] records
// * ”'"cname"”' for [[SoftLayer_Dns_Domain_ResourceRecord_CnameType|canonical name]] records
// * ”'"mx"”' for [[SoftLayer_Dns_Domain_ResourceRecord_MxType|mail exchanger]] records
// * ”'"ns"”' for [[SoftLayer_Dns_Domain_ResourceRecord_NsType|name server]] records
// * ”'"ptr"”' for [[SoftLayer_Dns_Domain_ResourceRecord_PtrType|pointer]] records in reverse domains
// * ”'"soa"”' for a domain's [[SoftLayer_Dns_Domain_ResourceRecord_SoaType|start of authority]] record
// * ”'"spf"”' for [[SoftLayer_Dns_Domain_ResourceRecord_SpfType|sender policy framework]] records
// * ”'"srv"”' for [[SoftLayer_Dns_Domain_ResourceRecord_SrvType|service]] records
// * ”'"txt"”' for [[SoftLayer_Dns_Domain_ResourceRecord_TxtType|text]] records
//
// As ”SoftLayer_Dns_Domain_ResourceRecord” objects are created and loaded, the API verifies the ”type” property and casts the object as the appropriate type.
type Dns_Domain_ResourceRecord struct {
	Session session.SLSession
	Options sl.Options
}

// GetDnsDomainResourceRecordService returns an instance of the Dns_Domain_ResourceRecord SoftLayer service
func GetDnsDomainResourceRecordService(sess session.SLSession) Dns_Domain_ResourceRecord {
	return Dns_Domain_ResourceRecord{Session: sess}
}

func (r Dns_Domain_ResourceRecord) Id(id int) Dns_Domain_ResourceRecord {
	r.Options.Id = &id
	return r
}

func (r Dns_Domain_ResourceRecord) Mask(mask string) Dns_Domain_ResourceRecord {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Dns_Domain_ResourceRecord) Filter(filter string) Dns_Domain_ResourceRecord {
	r.Options.Filter = filter
	return r
}

func (r Dns_Domain_ResourceRecord) Limit(limit int) Dns_Domain_ResourceRecord {
	r.Options.Limit = &limit
	return r
}

func (r Dns_Domain_ResourceRecord) Offset(offset int) Dns_Domain_ResourceRecord {
	r.Options.Offset = &offset
	return r
}

// createObject creates a new domain resource record. The ”host” property of the templateObject parameter is scrubbed to remove all non-alpha numeric characters except for "@", "_", ".", "*", and "-". The ”data” property of the templateObject parameter is scrubbed to remove all non-alphanumeric characters for "." and "-". Creating a resource record updates the serial number of the domain the resource record is associated with.
//
// ”createObject” returns Boolean ”true” on successful create or ”false” if it was unable to create a resource record.
func (r Dns_Domain_ResourceRecord) CreateObject(templateObject *datatypes.Dns_Domain_ResourceRecord) (resp datatypes.Dns_Domain_ResourceRecord, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord", "createObject", params, &r.Options, &resp)
	return
}

// Create multiple resource records on a domain. This follows the same logic as ”createObject'. The serial number of the domain associated with this resource record is updated upon creation.
//
// ”createObjects” returns Boolean ”true” on successful creation or ”false” if it was unable to create a resource record.
func (r Dns_Domain_ResourceRecord) CreateObjects(templateObjects []datatypes.Dns_Domain_ResourceRecord) (resp []datatypes.Dns_Domain_ResourceRecord, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord", "createObjects", params, &r.Options, &resp)
	return
}

// Delete a domain's resource record. ”'This cannot be undone.”' Be wary of running this method. If you remove a resource record in error you will need to re-create it by creating a new SoftLayer_Dns_Domain_ResourceRecord object. The serial number of the domain associated with this resource record is updated upon deletion. You may not delete SOA, NS, or PTR resource records.
//
// ”deleteObject” returns Boolean ”true” on successful deletion or ”false” if it was unable to remove a resource record.
func (r Dns_Domain_ResourceRecord) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord", "deleteObject", nil, &r.Options, &resp)
	return
}

// Remove multiple resource records from a domain. This follows the same logic as ”deleteObject” and ”'cannot be undone”'. The serial number of the domain associated with this resource record is updated upon deletion. You may not delete SOA records, PTR records, or NS resource records that point to ns1.softlayer.com or ns2.softlayer.com.
//
// ”deleteObjects” returns Boolean ”true” on successful deletion or ”false” if it was unable to remove a resource record.
func (r Dns_Domain_ResourceRecord) DeleteObjects(templateObjects []datatypes.Dns_Domain_ResourceRecord) (resp bool, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord", "deleteObjects", params, &r.Options, &resp)
	return
}

// editObject edits an existing domain resource record. The ”host” property of the templateObject parameter is scrubbed to remove all non-alpha numeric characters except for "@", "_", ".", "*", and "-". The ”data” property of the templateObject parameter is scrubbed to remove all non-alphanumeric characters for "." and "-". Editing a resource record updates the serial number of the domain the resource record is associated with.
//
// ”editObject” returns Boolean ”true” on a successful edit or ”false” if it was unable to edit the resource record.
func (r Dns_Domain_ResourceRecord) EditObject(templateObject *datatypes.Dns_Domain_ResourceRecord) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord", "editObject", params, &r.Options, &resp)
	return
}

// Edit multiple resource records on a domain. This follows the same logic as ”createObject'. The serial number of the domain associated with this resource record is updated upon creation.
//
// ”createObjects” returns Boolean ”true” on successful creation or ”false” if it was unable to create a resource record.
func (r Dns_Domain_ResourceRecord) EditObjects(templateObjects []datatypes.Dns_Domain_ResourceRecord) (resp bool, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord", "editObjects", params, &r.Options, &resp)
	return
}

// Retrieve The domain that a resource record belongs to.
func (r Dns_Domain_ResourceRecord) GetDomain() (resp datatypes.Dns_Domain, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord", "getDomain", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_Dns_Domain_ResourceRecord object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_Dns_Domain_ResourceRecord service. You can only retrieve resource records belonging to domains that are assigned to your SoftLayer account.
func (r Dns_Domain_ResourceRecord) GetObject() (resp datatypes.Dns_Domain_ResourceRecord, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord", "getObject", nil, &r.Options, &resp)
	return
}

// SoftLayer_Dns_Domain_ResourceRecord_MxType is a SoftLayer_Dns_Domain_ResourceRecord object whose ”type” property is set to "mx" and used to describe MX resource records. MX records control which hosts are responsible as mail exchangers for a domain. For instance, in the domain example.org, an MX record whose host is "@" and data is "mail" says that the host "mail.example.org" is responsible for handling mail for example.org. That means mail sent to users @example.org are delivered to mail.example.org.
//
// Domains can have more than one MX record if it uses more than one server to send mail through. Multiple MX records are denoted by their priority, defined by the mxPriority property.
//
// MX records must be defined for hosts with accompanying A or AAAA resource records. They may not point mail towards a host defined by a CNAME record.
type Dns_Domain_ResourceRecord_MxType struct {
	Session session.SLSession
	Options sl.Options
}

// GetDnsDomainResourceRecordMxTypeService returns an instance of the Dns_Domain_ResourceRecord_MxType SoftLayer service
func GetDnsDomainResourceRecordMxTypeService(sess session.SLSession) Dns_Domain_ResourceRecord_MxType {
	return Dns_Domain_ResourceRecord_MxType{Session: sess}
}

func (r Dns_Domain_ResourceRecord_MxType) Id(id int) Dns_Domain_ResourceRecord_MxType {
	r.Options.Id = &id
	return r
}

func (r Dns_Domain_ResourceRecord_MxType) Mask(mask string) Dns_Domain_ResourceRecord_MxType {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Dns_Domain_ResourceRecord_MxType) Filter(filter string) Dns_Domain_ResourceRecord_MxType {
	r.Options.Filter = filter
	return r
}

func (r Dns_Domain_ResourceRecord_MxType) Limit(limit int) Dns_Domain_ResourceRecord_MxType {
	r.Options.Limit = &limit
	return r
}

func (r Dns_Domain_ResourceRecord_MxType) Offset(offset int) Dns_Domain_ResourceRecord_MxType {
	r.Options.Offset = &offset
	return r
}

// createObject creates a new MX record. The ”host” property of the templateObject parameter is scrubbed to remove all non-alpha numeric characters except for "@", "_", ".", "*", and "-". The ”data” property of the templateObject parameter is scrubbed to remove all non-alphanumeric characters for "." and "-". Creating an MX record updates the serial number of the domain the resource record is associated with.
func (r Dns_Domain_ResourceRecord_MxType) CreateObject(templateObject *datatypes.Dns_Domain_ResourceRecord_MxType) (resp datatypes.Dns_Domain_ResourceRecord_MxType, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_MxType", "createObject", params, &r.Options, &resp)
	return
}

// Create multiple MX records on a domain. This follows the same logic as ”createObject'. The serial number of the domain associated with this MX record is updated upon creation.
//
// ”createObjects” returns Boolean ”true” on successful creation or ”false” if it was unable to create a resource record.
func (r Dns_Domain_ResourceRecord_MxType) CreateObjects(templateObjects []datatypes.Dns_Domain_ResourceRecord) (resp []datatypes.Dns_Domain_ResourceRecord, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_MxType", "createObjects", params, &r.Options, &resp)
	return
}

// Delete a domain's MX record. ”'This cannot be undone.”' Be wary of running this method. If you remove a resource record in error you will need to re-create it by creating a new SoftLayer_Dns_Domain_ResourceRecord_MxType object. The serial number of the domain associated with this MX record is updated upon deletion.
//
// ”deleteObject” returns Boolean ”true” on successful deletion or ”false” if it was unable to remove a resource record.
func (r Dns_Domain_ResourceRecord_MxType) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_MxType", "deleteObject", nil, &r.Options, &resp)
	return
}

// Remove multiple MX records from a domain. This follows the same logic as ”deleteObject” and ”'cannot be undone”'. The serial number of the domain associated with this MX record is updated upon deletion.
//
// ”deleteObjects” returns Boolean ”true” on successful deletion or ”false” if it was unable to remove a resource record.
func (r Dns_Domain_ResourceRecord_MxType) DeleteObjects(templateObjects []datatypes.Dns_Domain_ResourceRecord_MxType) (resp bool, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_MxType", "deleteObjects", params, &r.Options, &resp)
	return
}

// editObject edits an existing MX resource record. The ”host” property of the templateObject parameter is scrubbed to remove all non-alpha numeric characters except for "@", "_", ".", "*", and "-". The ”data” property of the templateObject parameter is scrubbed to remove all non-alphanumeric characters for "." and "-". Editing an MX record updates the serial number of the domain the record is associated with.
//
// ”editObject” returns Boolean ”true” on a successful edit or ”false” if it was unable to edit the resource record.
func (r Dns_Domain_ResourceRecord_MxType) EditObject(templateObject *datatypes.Dns_Domain_ResourceRecord_MxType) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_MxType", "editObject", params, &r.Options, &resp)
	return
}

// Edit multiple MX records on a domain. This follows the same logic as ”createObject'. The serial number of the domain associated with this MX record is updated upon creation.
//
// ”createObjects” returns Boolean ”true” on successful creation or ”false” if it was unable to create a resource record.
func (r Dns_Domain_ResourceRecord_MxType) EditObjects(templateObjects []datatypes.Dns_Domain_ResourceRecord_MxType) (resp bool, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_MxType", "editObjects", params, &r.Options, &resp)
	return
}

// Retrieve The domain that a resource record belongs to.
func (r Dns_Domain_ResourceRecord_MxType) GetDomain() (resp datatypes.Dns_Domain, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_MxType", "getDomain", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_Dns_Domain_ResourceRecord_MxType object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_Dns_Domain_ResourceRecord_MxType service. You can only retrieve resource records belonging to domains that are assigned to your SoftLayer account.
func (r Dns_Domain_ResourceRecord_MxType) GetObject() (resp datatypes.Dns_Domain_ResourceRecord_MxType, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_MxType", "getObject", nil, &r.Options, &resp)
	return
}

// SoftLayer_Dns_Domain_ResourceRecord_SrvType is a SoftLayer_Dns_Domain_ResourceRecord object whose ”type” property is set to "srv" and defines a DNS SRV record on a SoftLayer hosted domain.
type Dns_Domain_ResourceRecord_SrvType struct {
	Session session.SLSession
	Options sl.Options
}

// GetDnsDomainResourceRecordSrvTypeService returns an instance of the Dns_Domain_ResourceRecord_SrvType SoftLayer service
func GetDnsDomainResourceRecordSrvTypeService(sess session.SLSession) Dns_Domain_ResourceRecord_SrvType {
	return Dns_Domain_ResourceRecord_SrvType{Session: sess}
}

func (r Dns_Domain_ResourceRecord_SrvType) Id(id int) Dns_Domain_ResourceRecord_SrvType {
	r.Options.Id = &id
	return r
}

func (r Dns_Domain_ResourceRecord_SrvType) Mask(mask string) Dns_Domain_ResourceRecord_SrvType {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Dns_Domain_ResourceRecord_SrvType) Filter(filter string) Dns_Domain_ResourceRecord_SrvType {
	r.Options.Filter = filter
	return r
}

func (r Dns_Domain_ResourceRecord_SrvType) Limit(limit int) Dns_Domain_ResourceRecord_SrvType {
	r.Options.Limit = &limit
	return r
}

func (r Dns_Domain_ResourceRecord_SrvType) Offset(offset int) Dns_Domain_ResourceRecord_SrvType {
	r.Options.Offset = &offset
	return r
}

// createObject creates a new SRV record. The ”host” property of the templateObject parameter is scrubbed to remove all non-alpha numeric characters except for "@", "_", ".", "*", and "-". The ”data” property of the templateObject parameter is scrubbed to remove all non-alphanumeric characters for "." and "-". Creating an SRV record updates the serial number of the domain the resource record is associated with.
func (r Dns_Domain_ResourceRecord_SrvType) CreateObject(templateObject *datatypes.Dns_Domain_ResourceRecord_SrvType) (resp datatypes.Dns_Domain_ResourceRecord_SrvType, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_SrvType", "createObject", params, &r.Options, &resp)
	return
}

// Create multiple SRV records on a domain. This follows the same logic as ”createObject'. The serial number of the domain associated with this SRV record is updated upon creation.
//
// ”createObjects” returns Boolean ”true” on successful creation or ”false” if it was unable to create a resource record.
func (r Dns_Domain_ResourceRecord_SrvType) CreateObjects(templateObjects []datatypes.Dns_Domain_ResourceRecord) (resp []datatypes.Dns_Domain_ResourceRecord, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_SrvType", "createObjects", params, &r.Options, &resp)
	return
}

// Delete a domain's SRV record. ”'This cannot be undone.”' Be wary of running this method. If you remove a resource record in error you will need to re-create it by creating a new SoftLayer_Dns_Domain_ResourceRecord_SrvType object. The serial number of the domain associated with this SRV record is updated upon deletion.
//
// ”deleteObject” returns Boolean ”true” on successful deletion or ”false” if it was unable to remove a resource record.
func (r Dns_Domain_ResourceRecord_SrvType) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_SrvType", "deleteObject", nil, &r.Options, &resp)
	return
}

// Remove multiple SRV records from a domain. This follows the same logic as ”deleteObject” and ”'cannot be undone”'. The serial number of the domain associated with this SRV record is updated upon deletion.
//
// ”deleteObjects” returns Boolean ”true” on successful deletion or ”false” if it was unable to remove a resource record.
func (r Dns_Domain_ResourceRecord_SrvType) DeleteObjects(templateObjects []datatypes.Dns_Domain_ResourceRecord_SrvType) (resp bool, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_SrvType", "deleteObjects", params, &r.Options, &resp)
	return
}

// editObject edits an existing SRV resource record. The ”host” property of the templateObject parameter is scrubbed to remove all non-alpha numeric characters except for "@", "_", ".", "*", and "-". The ”data” property of the templateObject parameter is scrubbed to remove all non-alphanumeric characters for "." and "-". Editing an SRV record updates the serial number of the domain the record is associated with.
//
// ”editObject” returns Boolean ”true” on a successful edit or ”false” if it was unable to edit the resource record.
func (r Dns_Domain_ResourceRecord_SrvType) EditObject(templateObject *datatypes.Dns_Domain_ResourceRecord_SrvType) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_SrvType", "editObject", params, &r.Options, &resp)
	return
}

// Edit multiple SRV records on a domain. This follows the same logic as ”createObject'. The serial number of the domain associated with this SRV record is updated upon creation.
//
// ”createObjects” returns Boolean ”true” on successful creation or ”false” if it was unable to create a resource record.
func (r Dns_Domain_ResourceRecord_SrvType) EditObjects(templateObjects []datatypes.Dns_Domain_ResourceRecord_SrvType) (resp bool, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_SrvType", "editObjects", params, &r.Options, &resp)
	return
}

// Retrieve The domain that a resource record belongs to.
func (r Dns_Domain_ResourceRecord_SrvType) GetDomain() (resp datatypes.Dns_Domain, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_SrvType", "getDomain", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_Dns_Domain_ResourceRecord_SrvType object whose ID number corresponds to the ID number of the init parameter passed to the SoftLayer_Dns_Domain_ResourceRecord_SrvType service. You can only retrieve resource records belonging to domains that are assigned to your SoftLayer account.
func (r Dns_Domain_ResourceRecord_SrvType) GetObject() (resp datatypes.Dns_Domain_ResourceRecord_SrvType, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Domain_ResourceRecord_SrvType", "getObject", nil, &r.Options, &resp)
	return
}

// The SoftLayer_Dns_Secondary data type contains information on a single secondary DNS zone which is managed through SoftLayer's zone transfer service. Domains created via zone transfer may not be modified by the SoftLayer portal or API.
type Dns_Secondary struct {
	Session session.SLSession
	Options sl.Options
}

// GetDnsSecondaryService returns an instance of the Dns_Secondary SoftLayer service
func GetDnsSecondaryService(sess session.SLSession) Dns_Secondary {
	return Dns_Secondary{Session: sess}
}

func (r Dns_Secondary) Id(id int) Dns_Secondary {
	r.Options.Id = &id
	return r
}

func (r Dns_Secondary) Mask(mask string) Dns_Secondary {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Dns_Secondary) Filter(filter string) Dns_Secondary {
	r.Options.Filter = filter
	return r
}

func (r Dns_Secondary) Limit(limit int) Dns_Secondary {
	r.Options.Limit = &limit
	return r
}

func (r Dns_Secondary) Offset(offset int) Dns_Secondary {
	r.Options.Offset = &offset
	return r
}

// A secondary DNS record may be converted to a primary DNS record. By converting a secondary DNS record, the SoftLayer name servers will be the authoritative nameserver for this domain and will be directly editable in the SoftLayer API and Portal.
//
// Primary DNS record conversion performs the following steps:
// * The SOA record is updated with SoftLayer's primary name server.
// * All NS records are removed and replaced with SoftLayer's NS records.
// * The secondary DNS record is removed.
//
// After the DNS records are converted, the following restrictions will apply to the new domain record:
// * You will need to manage the zone record using the [[SoftLayer_Dns_Domain]] service.
// * You may not edit the SOA or NS records.
// * You may only edit the following resource records: A, AAAA, CNAME, MX, TX, SRV.
//
// This change can not be undone, and the record can not be converted back into a secondary DNS record once the conversion is complete.
func (r Dns_Secondary) ConvertToPrimary() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "convertToPrimary", nil, &r.Options, &resp)
	return
}

// Create a secondary DNS record. The ”zoneName”, ”masterIpAddress”, and ”transferFrequency” properties in the templateObject parameter are required parameters to create a secondary DNS record.
func (r Dns_Secondary) CreateObject(templateObject *datatypes.Dns_Secondary) (resp datatypes.Dns_Secondary, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "createObject", params, &r.Options, &resp)
	return
}

// Create multiple secondary DNS records. Each record passed to ”createObjects” follows the logic in the SoftLayer_Dns_Secondary [[SoftLayer_Dns_Secondary::createObject|createObject]] method.
func (r Dns_Secondary) CreateObjects(templateObjects []datatypes.Dns_Secondary) (resp []datatypes.Dns_Secondary, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "createObjects", params, &r.Options, &resp)
	return
}

// Delete a secondary DNS Record. This will also remove any associated domain records and resource records on the SoftLayer nameservers that were created as a result of the zone transfers. This action cannot be undone.
func (r Dns_Secondary) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "deleteObject", nil, &r.Options, &resp)
	return
}

// Edit the properties of a secondary DNS record by passing in a modified instance of a SoftLayer_Dns_Secondary object. You may only edit the ”masterIpAddress” and ”transferFrequency” properties of your secondary DNS record. ”ZoneName” may not be altered after a secondary DNS record has been created.  Please remove and re-create the record if you need to make changes to your zone name.
func (r Dns_Secondary) EditObject(templateObject *datatypes.Dns_Secondary) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "editObject", params, &r.Options, &resp)
	return
}

// Retrieve The SoftLayer account that owns a secondary DNS record.
func (r Dns_Secondary) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "getAccount", nil, &r.Options, &resp)
	return
}

// Search for [[SoftLayer_Dns_Secondary]] records by domain name. getByDomainName() performs an inclusive search for secondary domain records, returning multiple records based on partial name matches. Use this method to locate secondary domain records if you don't have access to their id numbers.
func (r Dns_Secondary) GetByDomainName(name *string) (resp []datatypes.Dns_Secondary, err error) {
	params := []interface{}{
		name,
	}
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "getByDomainName", params, &r.Options, &resp)
	return
}

// Retrieve The domain record created by zone transfer from a secondary DNS record.
func (r Dns_Secondary) GetDomain() (resp datatypes.Dns_Domain, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "getDomain", nil, &r.Options, &resp)
	return
}

// Retrieve The error messages created during secondary DNS record transfer.
func (r Dns_Secondary) GetErrorMessages() (resp []datatypes.Dns_Message, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "getErrorMessages", nil, &r.Options, &resp)
	return
}

// getObject retrieves the SoftLayer_Dns_Secondary object whose ID number corresponds to the ID number of the init paramater passed to the SoftLayer_Dns_Secondary service. You can only retrieve a secondary DNS record that is assigned to your SoftLayer customer account.
func (r Dns_Secondary) GetObject() (resp datatypes.Dns_Secondary, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The current status of the secondary DNS zone.
func (r Dns_Secondary) GetStatus() (resp datatypes.Dns_Status, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "getStatus", nil, &r.Options, &resp)
	return
}

// Force a secondary DNS zone transfer by setting it's status "Transfer Now".  A zone transfer will be initiated within a minute of receiving this API call.
func (r Dns_Secondary) TransferNow() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Dns_Secondary", "transferNow", nil, &r.Options, &resp)
	return
}
