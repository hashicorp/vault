package ssh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func sshConnect(b *backend) *framework.Path {
	log.Printf("Vishal: ssh.sshConnect\n")
	return &framework.Path{
		Pattern: "connect",
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "username at SSH host",
			},
			"address": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "IPv4 address of SSH host",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.sshConnectWrite,
		},
		HelpSynopsis:    sshConnectHelpSyn,
		HelpDescription: sshConnectHelpDesc,
	}
}

func (b *backend) sshConnectWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	ipAddr := d.Get("address").(string)
	log.Printf("Vishal: ssh.sshConnectWrite username:%#v address:%#v\n", username, ipAddr)

	hostKeyPath := "hosts/" + ipAddr + "/" + username
	entry, err := req.Storage.Get(hostKeyPath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, fmt.Errorf("Host key is not configured. Please configure them at the config/addhostkey endpoint")
	}
	var hostKey sshHostKey
	if err := entry.DecodeJSON(&hostKey); err != nil {
		return nil, fmt.Errorf("Error reading the host key: %s", err)
	}
	log.Printf("Vishal: host key previously configured: \n---------------\n%#v\n--------------\n", hostKey.Key)

	//TODO: save th entry in a file
	//TODO: read the hosts path and get the key
	//TODO: Input validation for the commands below
	hostKeyFileName := "./vault_ssh_" + username + "_" + ipAddr + "_shared.pem"
	err = ioutil.WriteFile(hostKeyFileName, []byte(hostKey.Key), 0400)

	otkPrivateKeyFileName := "vault_ssh_" + username + "_" + ipAddr + "_otk.pem"
	otkPublicKeyFileName := otkPrivateKeyFileName + ".pub"
	rmCmd := "rm -f " + otkPrivateKeyFileName + " " + otkPublicKeyFileName + ";"
	sshKeygenCmd := "ssh-keygen -f " + otkPrivateKeyFileName + " -t rsa -N ''" + ";"
	chmodCmd := "chmod 400 " + otkPrivateKeyFileName + ";"
	scpCmd := "scp -i " + hostKeyFileName + " " + otkPublicKeyFileName + " " + username + "@" + ipAddr + ":~;"

	log.Printf("Vishal: scpCmd: \n", scpCmd)

	localCmdString := strings.Join([]string{
		rmCmd,
		sshKeygenCmd,
		chmodCmd,
		scpCmd,
	}, "")
	err = exec_command(localCmdString)
	if err != nil {
		fmt.Errorf("Running command failed " + err.Error())
	}
	log.Printf("Vishal: Creating session\n")
	session := createSSHPublicKeysSession(username, ipAddr)
	var buf bytes.Buffer
	session.Stdout = &buf
	log.Printf("Vishal: Installing keys\n")
	if err := installSshOtkInTarget(session, username, ipAddr); err != nil {
		fmt.Errorf("Failed to install one-time-key at target machine: " + err.Error())
	}
	session.Close()
	fmt.Println(buf.String())
	keyBytes, err := ioutil.ReadFile(otkPrivateKeyFileName)
	oneTimeKey := string(keyBytes)
	log.Printf("Vishal: Returning:[%s]\n", oneTimeKey)
	return b.Secret(SecretOneTimeKeyType).Response(map[string]interface{}{
		"key": oneTimeKey,
	}, nil), nil
	/*return &logical.Response{
		Data: map[string]interface{}{
			"key": oneTimeKey,
		},
	}, nil
	*/
}

const sshConnectHelpSyn = `
sshConnectionHelpSyn
`

const sshConnectHelpDesc = `
rshConnectionHelpDesc
`
