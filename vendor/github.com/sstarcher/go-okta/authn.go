package okta

import (
	"time"
)

type ErrorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorSummary string `json:"errorSummary"`
	ErrorLink    string `json:"errorLink"`
	ErrorID      string `json:"errorId"`
	ErrorCauses  []struct {
		ErrorSummary string `json:"errorSummary"`
	} `json:"errorCauses"`
}

type AuthnRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RelayState string `json:"relayState"`
	Options    struct {
		MultiOptionalFactorEnroll bool `json:"multiOptionalFactorEnroll"`
		WarnBeforePasswordExpired bool `json:"warnBeforePasswordExpired"`
	} `json:"options"`
}

type AuthnResponse struct {
	ExpiresAt    time.Time `json:"expiresAt"`
	Status       string    `json:"status"`
	RelayState   string    `json:"relayState"`
	SessionToken string    `json:"sessionToken"`
	Embedded     struct {
		User struct {
			ID              string    `json:"id"`
			PasswordChanged time.Time `json:"passwordChanged"`
			Profile         struct {
				Login     string `json:"login"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Locale    string `json:"locale"`
				TimeZone  string `json:"timeZone"`
			} `json:"profile"`
		} `json:"user"`
	} `json:"_embedded"`
}
