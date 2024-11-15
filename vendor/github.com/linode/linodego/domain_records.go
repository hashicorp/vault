package linodego

import (
	"context"
)

// DomainRecord represents a DomainRecord object
type DomainRecord struct {
	ID       int              `json:"id"`
	Type     DomainRecordType `json:"type"`
	Name     string           `json:"name"`
	Target   string           `json:"target"`
	Priority int              `json:"priority"`
	Weight   int              `json:"weight"`
	Port     int              `json:"port"`
	Service  *string          `json:"service"`
	Protocol *string          `json:"protocol"`
	TTLSec   int              `json:"ttl_sec"`
	Tag      *string          `json:"tag"`
}

// DomainRecordCreateOptions fields are those accepted by CreateDomainRecord
type DomainRecordCreateOptions struct {
	Type     DomainRecordType `json:"type"`
	Name     string           `json:"name"`
	Target   string           `json:"target"`
	Priority *int             `json:"priority,omitempty"`
	Weight   *int             `json:"weight,omitempty"`
	Port     *int             `json:"port,omitempty"`
	Service  *string          `json:"service,omitempty"`
	Protocol *string          `json:"protocol,omitempty"`
	TTLSec   int              `json:"ttl_sec,omitempty"` // 0 is not accepted by Linode, so can be omitted
	Tag      *string          `json:"tag,omitempty"`
}

// DomainRecordUpdateOptions fields are those accepted by UpdateDomainRecord
type DomainRecordUpdateOptions struct {
	Type     DomainRecordType `json:"type,omitempty"`
	Name     string           `json:"name,omitempty"`
	Target   string           `json:"target,omitempty"`
	Priority *int             `json:"priority,omitempty"` // 0 is valid, so omit only nil values
	Weight   *int             `json:"weight,omitempty"`   // 0 is valid, so omit only nil values
	Port     *int             `json:"port,omitempty"`     // 0 is valid to spec, so omit only nil values
	Service  *string          `json:"service,omitempty"`
	Protocol *string          `json:"protocol,omitempty"`
	TTLSec   int              `json:"ttl_sec,omitempty"` // 0 is not accepted by Linode, so can be omitted
	Tag      *string          `json:"tag,omitempty"`
}

// DomainRecordType constants start with RecordType and include Linode API Domain Record Types
type DomainRecordType string

// DomainRecordType contants are the DNS record types a DomainRecord can assign
const (
	RecordTypeA     DomainRecordType = "A"
	RecordTypeAAAA  DomainRecordType = "AAAA"
	RecordTypeNS    DomainRecordType = "NS"
	RecordTypeMX    DomainRecordType = "MX"
	RecordTypeCNAME DomainRecordType = "CNAME"
	RecordTypeTXT   DomainRecordType = "TXT"
	RecordTypeSRV   DomainRecordType = "SRV"
	RecordTypePTR   DomainRecordType = "PTR"
	RecordTypeCAA   DomainRecordType = "CAA"
)

// GetUpdateOptions converts a DomainRecord to DomainRecordUpdateOptions for use in UpdateDomainRecord
func (d DomainRecord) GetUpdateOptions() (du DomainRecordUpdateOptions) {
	du.Type = d.Type
	du.Name = d.Name
	du.Target = d.Target
	du.Priority = copyInt(&d.Priority)
	du.Weight = copyInt(&d.Weight)
	du.Port = copyInt(&d.Port)
	du.Service = copyString(d.Service)
	du.Protocol = copyString(d.Protocol)
	du.TTLSec = d.TTLSec
	du.Tag = copyString(d.Tag)

	return
}

// ListDomainRecords lists DomainRecords
func (c *Client) ListDomainRecords(ctx context.Context, domainID int, opts *ListOptions) ([]DomainRecord, error) {
	response, err := getPaginatedResults[DomainRecord](ctx, c, formatAPIPath("domains/%d/records", domainID), opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetDomainRecord gets the domainrecord with the provided ID
func (c *Client) GetDomainRecord(ctx context.Context, domainID int, recordID int) (*DomainRecord, error) {
	e := formatAPIPath("domains/%d/records/%d", domainID, recordID)
	response, err := doGETRequest[DomainRecord](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateDomainRecord creates a DomainRecord
func (c *Client) CreateDomainRecord(ctx context.Context, domainID int, opts DomainRecordCreateOptions) (*DomainRecord, error) {
	e := formatAPIPath("domains/%d/records", domainID)
	response, err := doPOSTRequest[DomainRecord](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateDomainRecord updates the DomainRecord with the specified id
func (c *Client) UpdateDomainRecord(ctx context.Context, domainID int, recordID int, opts DomainRecordUpdateOptions) (*DomainRecord, error) {
	e := formatAPIPath("domains/%d/records/%d", domainID, recordID)
	response, err := doPUTRequest[DomainRecord](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteDomainRecord deletes the DomainRecord with the specified id
func (c *Client) DeleteDomainRecord(ctx context.Context, domainID int, recordID int) error {
	e := formatAPIPath("domains/%d/records/%d", domainID, recordID)
	err := doDELETERequest(ctx, c, e)
	return err
}
