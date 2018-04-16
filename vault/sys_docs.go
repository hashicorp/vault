package vault

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

var documentationPaths = []*framework.Path{
	&framework.Path{
		Pattern: "generate-root/attempt",

		Fields: map[string]*framework.FieldSchema{
			"otp": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Specifies a base64-encoded 16-byte value.",
			},
			"pgp_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Specifies a base64-encoded PGP public key.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   framework.NullCallback,
			logical.CreateOperation: framework.NullCallback,
			logical.DeleteOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Cancels any in-progress root generation attempt.",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "generate-root/update",

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Specifies a single master key share.",
			},
			"nonce": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Specifies the nonce of the attempt.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Enter a single master key share to progress the root generation attempt.",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "leader",

		Fields: map[string]*framework.FieldSchema{},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Check the high availability status and current leader of Vault",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "init",

		Fields: map[string]*framework.FieldSchema{
			"pgp_keys": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares.",
			},
			"root_token_pgp_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Specifies a PGP public key used to encrypt the initial root token. The key must be base64-encoded from its original binary representation.",
			},
			"secret_shares": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Specifies the number of shares to split the master key into.",
			},
			"secret_threshold": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Specifies the number of shares required to reconstruct the master key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares.",
			},
			"stored_shares": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Specifies the number of shares that should be encrypted by the HSM and stored for auto-unsealing. Currently must be the same as secret_shares.",
			},
			"recovery_pgp_keys": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Specifies an array of PGP public keys used to encrypt the output recovery keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as recovery_shares.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   framework.NullCallback,
			logical.CreateOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Initializes or returns the initialization status of the Vault.",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "health",

		Fields: map[string]*framework.FieldSchema{},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Returns the health status of Vault.",
		HelpDescription: "",
		Responses: map[logical.Operation]map[int]*framework.ResponseSchema{
			logical.ReadOperation: map[int]*framework.ResponseSchema{
				200: &framework.ResponseSchema{Description: "initialized, unsealed, and active"},
				429: &framework.ResponseSchema{Description: "unsealed and standby"},
				472: &framework.ResponseSchema{Description: "data recovery mode replication secondary and active"},
				501: &framework.ResponseSchema{Description: "not initialized"},
				503: &framework.ResponseSchema{Description: "sealed"},
			},
		},
	},

	&framework.Path{
		Pattern: "rekey/init",

		Fields: map[string]*framework.FieldSchema{
			"secret_shares": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Specifies the number of shares to split the master key into.",
			},
			"secret_threshold": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Specifies the number of shares required to reconstruct the master key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares.",
			},
			"pgp_keys": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares.",
			},
			"backup": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Specifies if using PGP-encrypted keys, whether Vault should also store a plaintext backup of the PGP-encrypted keys.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   framework.NullCallback,
			logical.CreateOperation: framework.NullCallback,
			logical.DeleteOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Cancels any in-progress rekey.",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "rekey/update",

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Specifies a single master key share.",
			},
			"nonce": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Specifies the nonce of the rekey attempt.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Enter a single master key share to progress the rekey of the Vault.",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "rekey-recovery-key/init",

		Fields: map[string]*framework.FieldSchema{
			"secret_shares": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Specifies the number of shares to split the recovery key into.",
			},
			"secret_threshold": &framework.FieldSchema{
				Type:        framework.TypeInt,
				Description: "Specifies the number of shares required to reconstruct the recovery key. This must be less than or equal secret_shares. If using Vault HSM with auto-unsealing, this value must be the same as secret_shares.",
			},
			"pgp_keys": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Specifies an array of PGP public keys used to encrypt the output unseal keys. Ordering is preserved. The keys must be base64-encoded from their original binary representation. The size of this array must be the same as secret_shares.",
			},
			"backup": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Specifies if using PGP-encrypted keys, whether Vault should also store a plaintext backup of the PGP-encrypted keys.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   framework.NullCallback,
			logical.CreateOperation: framework.NullCallback,
			logical.DeleteOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Cancels any in-progress rekey.",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "rekey-recovery-key/update",

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Specifies a single master key share.",
			},
			"nonce": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Specifies the nonce of the rekey attempt.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Enter a single master key share to progress the rekey of the Vault.",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "rekey-recovery-key/backup",

		Fields: map[string]*framework.FieldSchema{},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   framework.NullCallback,
			logical.DeleteOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Deletes the backup copy of PGP-encrypted recovery key shares.",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "seal-status",

		Fields: map[string]*framework.FieldSchema{},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Returns the seal status of the Vault.",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "seal",

		Fields: map[string]*framework.FieldSchema{},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Seals the Vault.",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "step-down",

		Fields: map[string]*framework.FieldSchema{},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Causes the node to give up active status.",
		HelpDescription: "",
	},

	&framework.Path{
		Pattern: "unseal",

		Fields: map[string]*framework.FieldSchema{
			"key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Specifies a single master key share. This is required unless reset is true.",
			},
			"reset": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Specifies if previously-provided unseal keys are discarded and the unseal process is reset.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: framework.NullCallback,
		},

		HelpSynopsis:    "Unseals the Vault.",
		HelpDescription: "",
	},
}
