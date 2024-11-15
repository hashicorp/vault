package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// Payment represents a Payment object
type Payment struct {
	// The unique ID of the Payment
	ID int `json:"id"`

	// The amount, in US dollars, of the Payment.
	USD json.Number `json:"usd"`

	// When the Payment was made.
	Date *time.Time `json:"-"`
}

// PaymentCreateOptions fields are those accepted by CreatePayment
type PaymentCreateOptions struct {
	// CVV (Card Verification Value) of the credit card to be used for the Payment
	CVV string `json:"cvv,omitempty"`

	// The amount, in US dollars, of the Payment
	USD json.Number `json:"usd"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *Payment) UnmarshalJSON(b []byte) error {
	type Mask Payment

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

// GetCreateOptions converts a Payment to PaymentCreateOptions for use in CreatePayment
func (i Payment) GetCreateOptions() (o PaymentCreateOptions) {
	o.USD = i.USD
	return
}

// ListPayments lists Payments
func (c *Client) ListPayments(ctx context.Context, opts *ListOptions) ([]Payment, error) {
	response, err := getPaginatedResults[Payment](ctx, c, "account/payments", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetPayment gets the payment with the provided ID
func (c *Client) GetPayment(ctx context.Context, paymentID int) (*Payment, error) {
	e := formatAPIPath("account/payments/%d", paymentID)
	response, err := doGETRequest[Payment](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreatePayment creates a Payment
func (c *Client) CreatePayment(ctx context.Context, opts PaymentCreateOptions) (*Payment, error) {
	e := "accounts/payments"
	response, err := doPOSTRequest[Payment](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}
