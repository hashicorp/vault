package api

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// Default path at which SSH backend will be mounted
const SSHAgentDefaultMountPoint = "ssh"

// This is a structure representing an SSH agent which can talk to vault server
// in order to verify the OTP entered by the user. It contains the path at which
// SSH backend is mounted at the server.
type SSHAgent struct {
	c          *Client
	MountPoint string
}

// SSHVerifyResp is a structure representing the fields in Vault server's
// response.
type SSHVerifyResponse struct {
	Message  string `mapstructure:"message"`
	Username string `mapstructure:"username"`
	IP       string `mapstructure:"ip"`
}

// Creates an SSHAgent object which can talk to Vault server with SSH backend
// mounted at default path ("ssh").
func (c *Client) SSHAgent() *SSHAgent {
	return c.SSHAgentWithMountPoint(SSHAgentDefaultMountPoint)
}

// Creates an SSHAgent object which can talk to Vault server with SSH backend
// mounted at a specific mount point.
func (c *Client) SSHAgentWithMountPoint(mountPoint string) *SSHAgent {
	return &SSHAgent{
		c:          c,
		MountPoint: mountPoint,
	}
}

// Verifies if the key provided by user is present in Vault server. If yes,
// the response will contain the IP address and username associated with the
// key.
func (c *SSHAgent) Verify(otp string) (*SSHVerifyResponse, error) {
	data := map[string]interface{}{
		"otp": otp,
	}
	verifyPath := fmt.Sprintf("/v1/%s/verify", c.MountPoint)
	r := c.c.NewRequest("PUT", verifyPath)
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	secret, err := ParseSecret(resp.Body)
	if err != nil {
		return nil, err
	}

	if secret.Data == nil {
		return nil, nil
	}

	var verifyResp SSHVerifyResponse
	err = mapstructure.Decode(secret.Data, &verifyResp)
	if err != nil {
		return nil, err
	}
	return &verifyResp, nil
}
