package creds

// TODO I don't think this actually needs to be exported,
// but if I don't I need to change the var name
type Cred struct {
	RoleName        string `json:"role_name"`
	Username        string `json:"username"`
	CurrentPassword string `json:"current_password"`
	LastPassword    string `json:"last_password,omitempty"`
}

func (c *Cred) Map() map[string]interface{} {
	m := map[string]interface{}{
		"role_name":        c.RoleName,
		"username":         c.Username,
		"current_password": c.CurrentPassword,
	}
	if c.LastPassword != "" {
		m["last_password"] = c.LastPassword
	}
	return m
}
