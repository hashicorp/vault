package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const configPath = "config"

const domainConfigPropertyName = "domain"
const applicationIdConfigPropertyName = "applicationId"
const applicationSecretConfigPropertyName = "applicationSecret"
const TTLConfigPropertyName = "ttl"
const maxTTLConfigPropertyName = "max_ttl"

const writeConfigPathHelp = `configure the google credential backend with applicationId and applicationSecret first:
vault write auth/google/config applicationId=$GOOGLE_APPLICATION_ID applicationSecret=$GOOGLE_APPLICATION_SECRET domain=example.com`


func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: configPath,
		Fields: map[string]*framework.FieldSchema{
			domainConfigPropertyName: &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "The domain users must be part of",
			},
			applicationIdConfigPropertyName: &framework.FieldSchema{
				Type:	     framework.TypeString,
				Description: "google application id",
			},
			applicationSecretConfigPropertyName: &framework.FieldSchema{
				Type:	     framework.TypeString,
				Description: "google application secret",
			},
			TTLConfigPropertyName: &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Duration after which authentication will be expired`,
			},
			maxTTLConfigPropertyName: &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Maximum duration after which authentication will be expired`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigWrite,
		},
	}
}

const configEntry = "config"

func readDurationFromData(data *framework.FieldData, property string) (time.Duration, *logical.Response) {
	ttlRaw, ok := data.GetOk(property)
	var ttl time.Duration
	var err error
	var rsp *logical.Response
	if !ok || len(ttlRaw.(string)) == 0 {
		ttl = 0
		rsp = nil
	} else {
		ttl, err = time.ParseDuration(ttlRaw.(string))
		if err != nil {
			rsp = logical.ErrorResponse(fmt.Sprintf("Invalid '%s':%s", property , err))
		}
	}
	return ttl, rsp
}

func (b *backend) pathConfigWrite(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	domain := data.Get(domainConfigPropertyName).(string)
	applicationID := data.Get(applicationIdConfigPropertyName).(string)
	applicationSecret := data.Get(applicationSecretConfigPropertyName).(string)

	var rsp *logical.Response
	var ttl time.Duration

	ttl, rsp = readDurationFromData(data, TTLConfigPropertyName)
	if (rsp != nil) {
		return rsp, nil
	}

	var maxTTL time.Duration
	maxTTL, rsp = readDurationFromData(data, maxTTLConfigPropertyName)
	if (rsp != nil) {
		return rsp, nil
	}

	entry, err := logical.StorageEntryJSON(configEntry, config{
		Domain:     domain,
		TTL:     ttl,
		MaxTTL:  maxTTL,
		ApplicationID: applicationID,
		ApplicationSecret: applicationSecret,
	})

	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}



// Config returns the configuration for this backend.
func (b *backend) Config(s logical.Storage) (*config, error) {
	entry, err := s.Get(configEntry)
	if err != nil {
		return nil, err
	}

	var result config
	if entry != nil {
		if err := entry.DecodeJSON(&result); err != nil {
			return nil, fmt.Errorf("error reading configuration: %s", err)
		}
	}

	return &result, nil
}

type config struct {
	Domain            string        `json:"domain"`
	ApplicationID     string `json:"applicationId"`
	ApplicationSecret string `json:"applicationSecret"`
	TTL               time.Duration `json:"ttl"`
	MaxTTL            time.Duration `json:"max_ttl"`
}
