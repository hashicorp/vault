package packngo

// API documentation https://www.packet.net/developers/api/paymentmethods/
const paymentMethodBasePath = "/payment-methods"

// ProjectService interface defines available project methods
type PaymentMethodService interface {
	List() ([]PaymentMethod, *Response, error)
	Get(string) (*PaymentMethod, *Response, error)
	Create(*PaymentMethodCreateRequest) (*PaymentMethod, *Response, error)
	Update(string, *PaymentMethodUpdateRequest) (*PaymentMethod, *Response, error)
	Delete(string) (*Response, error)
}

type paymentMethodsRoot struct {
	PaymentMethods []PaymentMethod `json:"payment_methods"`
}

// PaymentMethod represents a Packet payment method of an organization
type PaymentMethod struct {
	ID             string         `json:"id"`
	Name           string         `json:"name,omitempty"`
	Created        string         `json:"created_at,omitempty"`
	Updated        string         `json:"updated_at,omitempty"`
	Nonce          string         `json:"nonce,omitempty"`
	Default        bool           `json:"default,omitempty"`
	Organization   Organization   `json:"organization,omitempty"`
	Projects       []Project      `json:"projects,omitempty"`
	Type           string         `json:"type,omitempty"`
	CardholderName string         `json:"cardholder_name,omitempty"`
	ExpMonth       string         `json:"expiration_month,omitempty"`
	ExpYear        string         `json:"expiration_year,omitempty"`
	Last4          string         `json:"last_4,omitempty"`
	BillingAddress BillingAddress `json:"billing_address,omitempty"`
	URL            string         `json:"href,omitempty"`
}

func (pm PaymentMethod) String() string {
	return Stringify(pm)
}

// PaymentMethodCreateRequest type used to create a Packet payment method of an organization
type PaymentMethodCreateRequest struct {
	Name           string `json:"name"`
	Nonce          string `json:"name"`
	CardholderName string `json:"cardholder_name,omitempty"`
	ExpMonth       string `json:"expiration_month,omitempty"`
	ExpYear        string `json:"expiration_year,omitempty"`
	BillingAddress string `json:"billing_address,omitempty"`
}

func (pm PaymentMethodCreateRequest) String() string {
	return Stringify(pm)
}

// PaymentMethodUpdateRequest type used to update a Packet payment method of an organization
type PaymentMethodUpdateRequest struct {
	Name           *string `json:"name,omitempty"`
	CardholderName *string `json:"cardholder_name,omitempty"`
	ExpMonth       *string `json:"expiration_month,omitempty"`
	ExpYear        *string `json:"expiration_year,omitempty"`
	BillingAddress *string `json:"billing_address,omitempty"`
}

func (pm PaymentMethodUpdateRequest) String() string {
	return Stringify(pm)
}

// PaymentMethodServiceOp implements PaymentMethodService
type PaymentMethodServiceOp struct {
	client *Client
}
