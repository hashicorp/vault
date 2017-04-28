package okta

import (
	"time"
)

type SessionRequest struct {
	SessionToken string `json:"sessionToken"`
}

type SessionResponse struct {
	ID                       string      `json:"id"`
	Login                    string      `json:"login"`
	UserID                   string      `json:"userId"`
	ExpiresAt                time.Time   `json:"expiresAt"`
	Status                   string      `json:"status"`
	LastPasswordVerification time.Time   `json:"lastPasswordVerification"`
	LastFactorVerification   interface{} `json:"lastFactorVerification"`
	Amr                      []string    `json:"amr"`
	Idp                      struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"idp"`
	MfaActive bool `json:"mfaActive"`
	Links     struct {
		Self struct {
			Href  string `json:"href"`
			Hints struct {
				Allow []string `json:"allow"`
			} `json:"hints"`
		} `json:"self"`
		Refresh struct {
			Href  string `json:"href"`
			Hints struct {
				Allow []string `json:"allow"`
			} `json:"hints"`
		} `json:"refresh"`
		User struct {
			Name  string `json:"name"`
			Href  string `json:"href"`
			Hints struct {
				Allow []string `json:"allow"`
			} `json:"hints"`
		} `json:"user"`
	} `json:"_links"`
}
