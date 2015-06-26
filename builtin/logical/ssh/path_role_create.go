package ssh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
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

	//fetch the parameters
	roleName := d.Get("name").(string)
	username := d.Get("username").(string)
	ipRaw := d.Get("ip").(string)

	//find the role to be used for installing dynamic key
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

	//validate the IP address
	ipAddr := net.ParseIP(ipRaw)
	if ipAddr == nil {
		return logical.ErrorResponse(fmt.Sprintf("Invalid IP '%s'", ipRaw)), nil
	}
	ip := ipAddr.String()

	var cidrEntry sshCIDR
	ipMatched := false
	json.Unmarshal([]byte(role.CIDR), &cidrEntry)
	for _, item := range cidrEntry.CIDR {
		log.Println(item)
		_, cidrIPNet, _ := net.ParseCIDR(item)
		ipMatched = cidrIPNet.Contains(ipAddr)
		if ipMatched {
			break
		}
	}
	if !ipMatched {
		return logical.ErrorResponse(fmt.Sprintf("IP[%s] does not belong to role[%s]", ip, roleName)), nil
	}

	//fetch the host key to be used for installation
	keyPath := "keys/" + role.KeyName
	keyEntry, err := req.Storage.Get(keyPath)
	if err != nil {
		return nil, fmt.Errorf("Key '%s' not found error:%s", role.KeyName, err)
	}
	var hostKey sshHostKey
	if err := keyEntry.DecodeJSON(&hostKey); err != nil {
		return nil, fmt.Errorf("Error reading the host key: %s", err)
	}

	//store the host key to file. Use it as parameter for scp command
	hostKeyFileName := "./vault_ssh_" + username + "_" + ip + "_shared.pem"
	err = ioutil.WriteFile(hostKeyFileName, []byte(hostKey.Key), 0400)

	otkPrivateKeyFileName := "vault_ssh_" + username + "_" + ip + "_otk.pem"
	otkPublicKeyFileName := otkPrivateKeyFileName + ".pub"

	//commands to be run on vault server
	removeFile(otkPrivateKeyFileName)
	removeFile(otkPublicKeyFileName)
	dynamicPublicKey, dynamicPrivateKey, _ := generateRSAKeys()
	ioutil.WriteFile(otkPrivateKeyFileName, []byte(dynamicPrivateKey), 0600)
	ioutil.WriteFile(otkPublicKeyFileName, []byte(dynamicPublicKey), 0600)
	//ioutil.WriteFile("testkey.pub", []byte(publicKeyRsa), 0600)
	//sshKeygenCmd := "ssh-keygen -f " + otkPrivateKeyFileName + " -t rsa -N ''" + ";"
	//chmodCmd := "chmod 600 " + otkPrivateKeyFileName + ";"
	scpCmd := "scp -i " + hostKeyFileName + " " + otkPublicKeyFileName + " " + username + "@" + ip + ":~;"
	localCmdString := strings.Join([]string{
		scpCmd,
	}, "")
	//run the commands on vault server
	err = exec_command(localCmdString)
	if err != nil {
		fmt.Errorf("Running command failed " + err.Error())
	}

	//connect to target machine
	session := createSSHPublicKeysSession(username, ip, hostKey.Key)
	var buf bytes.Buffer
	session.Stdout = &buf

	authKeysFileName := "~/.ssh/authorized_keys"
	tempKeysFileName := "~/temp_authorized_keys"

	//commands to be run on target machine
	grepCmd := "grep -vFf " + otkPublicKeyFileName + " " + authKeysFileName + " > " + tempKeysFileName + ";"
	catCmdRemoveDuplicate := "cat " + tempKeysFileName + " > " + authKeysFileName + ";"
	catCmdAppendNew := "cat " + otkPublicKeyFileName + " >> " + authKeysFileName + ";"
	removeCmd := "rm -f " + tempKeysFileName + " " + otkPublicKeyFileName + ";"
	remoteCmdString := strings.Join([]string{
		grepCmd,
		catCmdRemoveDuplicate,
		catCmdAppendNew,
		removeCmd,
	}, "")

	//run the commands on target machine
	if err := session.Run(remoteCmdString); err != nil {
		return nil, err
	}
	session.Close()
	fmt.Println(buf.String())

	//preparing the secret
	dynamicPrivateKeyBytes, err := ioutil.ReadFile(otkPrivateKeyFileName)
	if err != nil {
		fmt.Errorf("Failed to open '%s':%s", otkPrivateKeyFileName, err)
	}
	dynamicPrivateKey = string(dynamicPrivateKeyBytes)

	dynamicPublicKeyBytes, err := ioutil.ReadFile(otkPublicKeyFileName)
	if err != nil {
		fmt.Errorf("Failed to open '%s':%s", otkPublicKeyFileName, err)
	}
	dynamicPublicKey = string(dynamicPublicKeyBytes)
	return b.Secret(SecretOneTimeKeyType).Response(map[string]interface{}{
		"key": dynamicPrivateKey,
	}, map[string]interface{}{
		"username":           username,
		"ip":                 ip,
		"host_key_name":      role.KeyName,
		"dynamic_public_key": dynamicPublicKey,
	}), nil
}

type sshCIDR struct {
	CIDR []string
}

const sshConnectHelpSyn = `
sshConnectionHelpSyn
`

const sshConnectHelpDesc = `
rshConnectionHelpDesc
`
