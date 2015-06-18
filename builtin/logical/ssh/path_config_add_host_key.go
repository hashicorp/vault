package ssh

import (
	"bytes"
	"fmt"
	"log"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigAddHostKey(b *backend) *framework.Path {
	log.Printf("Vishal: ssh.pathConfigAddHostKey\n")
	return &framework.Path{
		Pattern: "config/addhostkey",
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username in host.",
			},
			"ip": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "IP address of host.",
			},
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "SSH private key for host.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathAddHostKeyWrite,
		},
		HelpSynopsis:    pathConfigAddHostKeySyn,
		HelpDescription: pathConfigAddHostKeyDesc,
	}
}

func (b *backend) pathAddHostKeyWrite(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	log.Printf("Vishal: ssh.pathAddHostKeyWrite\n")
	username := d.Get("username").(string)
	ip := d.Get("ip").(string)
	key := d.Get("key").(string)
	log.Printf("Vishal: ssh.pathAddHostKeyWrite username:%#v ip:%#v key:%#v\n", username, ip, key)
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
	return nil, nil
}

const pathConfigAddHostKeySyn = `
pathConfigAddHostKeySyn
`

const pathConfigAddHostKeyDesc = `
pathConfigAddHostKeyDesc
`
