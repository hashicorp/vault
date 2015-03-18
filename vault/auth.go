package vault

// setupCredentials is invoked after we've loaded the mount table to
// initialize the credential backends and setup the router
func (c *Core) setupCredentials() error {
	return nil
}

// teardownCredentials is used before we seal the vault to reset the credential
// backends to their unloaded state. This is reversed by loadCredentials.
func (c *Core) teardownCredentials() error {
	return nil
}
