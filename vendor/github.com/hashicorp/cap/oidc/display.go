package oidc

// Display is a string value that specifies how the Authorization Server
// displays the authentication and consent user interface pages to the End-User.
//
// See: https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
type Display string

const (
	// Defined the Display values that specifies how the Authorization Server
	// displays the authentication and consent user interface pages to the End-User.
	//
	// See: https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
	Page  Display = "page"
	Popup Display = "popup"
	Touch Display = "touch"
	WAP   Display = "wap"
)
