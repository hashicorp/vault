package oidc

// Prompt is a string values that specifies whether the Authorization Server
// prompts the End-User for reauthentication and consent.
//
// See: https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
type Prompt string

const (
	// Defined the Prompt values that specifies whether the Authorization Server
	// prompts the End-User for reauthentication and consent.
	//
	// See: https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	None          Prompt = "none"
	Login         Prompt = "login"
	Consent       Prompt = "consent"
	SelectAccount Prompt = "select_account"
)
