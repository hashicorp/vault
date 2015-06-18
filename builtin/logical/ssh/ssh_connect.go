package ssh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

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
	log.Printf("Vishal: ssh.sshConnectWrite username:%#v address:%#v\n", d.Get("username").(string), d.Get("address").(string))

	//username := d.Get("username").(string)
	//ip := d.Get("ip").(string)
	//key := d.Get("key").(string)
	//log.Printf("Vishal: ssh.pathAddHostKeyWrite username:%#v ip:%#v key:%#v\n", username, ip, key)
	localCmdString := `
	rm -f vault_ssh_otk.pem vault_ssh_otk.pem.pub;
	ssh-keygen -f vault_ssh_otk.pem -t rsa -N '';
	chmod 400 vault_ssh_otk.pem;
	scp -i vault_ssh_shared.pem vault_ssh_otk.pem.pub vishal@localhost:/home/vishal
	echo done!
	`
	err := exec_command(localCmdString)
	if err != nil {
		fmt.Errorf("Running command failed " + err.Error())
	}
	session := createSSHPublicKeysSession("vishal", "127.0.0.1")
	var buf bytes.Buffer
	session.Stdout = &buf
	if err := installSshOtkInTarget(session); err != nil {
		fmt.Errorf("Failed to install one-time-key at target machine: " + err.Error())
	}
	session.Close()
	fmt.Println(buf.String())
	keyBytes, err := ioutil.ReadFile("vault_ssh_otk.pem")
	oneTimeKey := string(keyBytes)
	return &logical.Response{
		Data: map[string]interface{}{
			"key": oneTimeKey,
		},
	}, nil
}

const sshConnectHelpSyn = `
sshConnectionHelpSyn
`

const sshConnectHelpDesc = `
sshConnectionHelpDesc
`
