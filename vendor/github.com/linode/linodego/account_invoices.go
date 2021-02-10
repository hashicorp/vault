package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
	"github.com/linode/linodego/pkg/errors"
)

// Invoice structs reflect an invoice for billable activity on the account.
type Invoice struct {
	ID    int        `json:"id"`
	Label string     `json:"label"`
	Total float32    `json:"total"`
	Date  *time.Time `json:"-"`
}

// InvoiceItem structs reflect an single billable activity associate with an Invoice
type InvoiceItem struct {
	Label     string     `json:"label"`
	Type      string     `json:"type"`
	UnitPrice int        `json:"unitprice"`
	Quantity  int        `json:"quantity"`
	Amount    float32    `json:"amount"`
	From      *time.Time `json:"-"`
	To        *time.Time `json:"-"`
}

// InvoicesPagedResponse represents a paginated Invoice API response
type InvoicesPagedResponse struct {
	*PageOptions
	Data []Invoice `json:"data"`
}

// endpoint gets the endpoint URL for Invoice
func (InvoicesPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.Invoices.Endpoint()
	if err != nil {
		panic(err)
	}

	return endpoint
}

// appendData appends Invoices when processing paginated Invoice responses
func (resp *InvoicesPagedResponse) appendData(r *InvoicesPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListInvoices gets a paginated list of Invoices against the Account
func (c *Client) ListInvoices(ctx context.Context, opts *ListOptions) ([]Invoice, error) {
	response := InvoicesPagedResponse{}
	err := c.listHelper(ctx, &response, opts)

	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *Invoice) UnmarshalJSON(b []byte) error {
	type Mask Invoice

	p := struct {
		*Mask
		Date *parseabletime.ParseableTime `json:"date"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Date = (*time.Time)(p.Date)

	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *InvoiceItem) UnmarshalJSON(b []byte) error {
	type Mask InvoiceItem

	p := struct {
		*Mask
		From *parseabletime.ParseableTime `json:"from"`
		To   *parseabletime.ParseableTime `json:"to"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.From = (*time.Time)(p.From)
	i.To = (*time.Time)(p.To)

	return nil
}

// GetInvoice gets the a single Invoice matching the provided ID
func (c *Client) GetInvoice(ctx context.Context, id int) (*Invoice, error) {
	e, err := c.Invoices.Endpoint()
	if err != nil {
		return nil, err
	}

	e = fmt.Sprintf("%s/%d", e, id)
	r, err := errors.CoupleAPIErrors(c.R(ctx).SetResult(&Invoice{}).Get(e))

	if err != nil {
		return nil, err
	}

	return r.Result().(*Invoice), nil
}

// InvoiceItemsPagedResponse represents a paginated Invoice Item API response
type InvoiceItemsPagedResponse struct {
	*PageOptions
	Data []InvoiceItem `json:"data"`
}

// endpointWithID gets the endpoint URL for InvoiceItems associated with a specific Invoice
func (InvoiceItemsPagedResponse) endpointWithID(c *Client, id int) string {
	endpoint, err := c.InvoiceItems.endpointWithID(id)
	if err != nil {
		panic(err)
	}

	return endpoint
}

// appendData appends InvoiceItems when processing paginated Invoice Item responses
func (resp *InvoiceItemsPagedResponse) appendData(r *InvoiceItemsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListInvoiceItems gets the invoice items associated with a specific Invoice
func (c *Client) ListInvoiceItems(ctx context.Context, id int, opts *ListOptions) ([]InvoiceItem, error) {
	response := InvoiceItemsPagedResponse{}
	err := c.listHelperWithID(ctx, &response, id, opts)

	if err != nil {
		return nil, err
	}

	return response.Data, nil
}
