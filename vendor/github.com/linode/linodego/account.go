package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// Account associated with the token in use.
type Account struct {
	FirstName         string      `json:"first_name"`
	LastName          string      `json:"last_name"`
	Email             string      `json:"email"`
	Company           string      `json:"company"`
	Address1          string      `json:"address_1"`
	Address2          string      `json:"address_2"`
	Balance           float32     `json:"balance"`
	BalanceUninvoiced float32     `json:"balance_uninvoiced"`
	City              string      `json:"city"`
	State             string      `json:"state"`
	Zip               string      `json:"zip"`
	Country           string      `json:"country"`
	TaxID             string      `json:"tax_id"`
	Phone             string      `json:"phone"`
	CreditCard        *CreditCard `json:"credit_card"`
	EUUID             string      `json:"euuid"`
	BillingSource     string      `json:"billing_source"`
	Capabilities      []string    `json:"capabilities"`
	ActiveSince       *time.Time  `json:"-"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (account *Account) UnmarshalJSON(b []byte) error {
	type Mask Account

	p := struct {
		*Mask
		ActiveSince *parseabletime.ParseableTime `json:"active_since"`
	}{
		Mask: (*Mask)(account),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	account.ActiveSince = (*time.Time)(p.ActiveSince)

	return nil
}

// CreditCard information associated with the Account.
type CreditCard struct {
	LastFour string `json:"last_four"`
	Expiry   string `json:"expiry"`
}

// GetAccount gets the contact and billing information related to the Account.
func (c *Client) GetAccount(ctx context.Context) (*Account, error) {
	e := "account"
	response, err := doGETRequest[Account](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}
