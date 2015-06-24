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

func pathRoleCreate(b *backend) *framework.Path {
	log.Printf("Vishal: ssh.sshConnect\n")
	return &framework.Path{
		Pattern: "creds/(?P<name>\\w+)",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "name of the policy",
			},
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "username in target",
			},
			"ip": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "IP of the target machine",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathRoleCreateWrite,
		},
		HelpSynopsis:    sshConnectHelpSyn,
		HelpDescription: sshConnectHelpDesc,
	}
}

func (b *backend) pathRoleCreateWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.pathRoleCreateWrite\n")

	roleName := d.Get("name").(string)
	username := d.Get("username").(string)
	ipAddr := d.Get("ip").(string)

	rolePath := "policy/" + roleName
	roleEntry, err := req.Storage.Get(rolePath)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %s", err)
	}
	if roleEntry == nil {
		return logical.ErrorResponse(fmt.Sprintf("Role '%s' not found", roleName)), nil
	}

	var role sshRole
	if err := roleEntry.DecodeJSON(&role); err != nil {
		return nil, err
	}

	log.Printf("Vishal: ssh.pathRoleCreateWrite username:%#v address:%#v name:%#v result:%s\n", username, ipAddr, roleName, role)
	//TODO: do the role verification here

	keyPath := "keys/" + role.KeyName
	keyEntry, err := req.Storage.Get(keyPath)
	if err != nil {
		return nil, fmt.Errorf("Key '%s' not found error:%s", role.KeyName, err)
	}

	log.Printf("Vishal: KeyName:%s keyPath:%s\n", role.KeyName, keyPath)
	var hostKey sshHostKey
	if err := keyEntry.DecodeJSON(&hostKey); err != nil {
		return nil, fmt.Errorf("Error reading the host key: %s", err)
	}
	log.Printf("Vishal: host key previously configured: \n---------------\n%#v\n--------------\n", hostKey.Key)

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
}

const sshConnectHelpSyn = `
sshConnectionHelpSyn
`

const sshConnectHelpDesc = `
rshConnectionHelpDesc
`
