// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package mock

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/automatedrotationutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/rotation"
)

// pathConfig is used to test auto rotation.
func pathConfig(b *backend) *framework.Path {
	p := &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"fail_rotate": {
				Type: framework.TypeBool,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathConfigUpdate,
			logical.UpdateOperation: b.pathConfigUpdate,
			logical.ReadOperation:   b.pathConfigRead,
		},
		ExistenceCheck: b.pathConfigExistenceCheck,
	}
	automatedrotationutil.AddAutomatedRotationFields(p.Fields)
	return p
}

func (b *backend) pathConfigUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := b.configEntry(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if conf == nil {
		conf = &config{}
	}

	if failRotateRaw, ok := data.GetOk("fail_rotate"); ok {
		conf.FailRotate = failRotateRaw.(bool)
	}

	if err := conf.ParseAutomatedRotationFields(data); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if conf.ShouldDeregisterRotationJob() {
		deregisterReq := &rotation.RotationJobDeregisterRequest{
			MountPoint: req.MountPoint,
			ReqPath:    req.Path,
		}

		b.Logger().Debug("Deregistering rotation job", "mount", req.MountPoint+req.Path)
		if err := b.System().DeregisterRotationJob(ctx, deregisterReq); err != nil {
			return logical.ErrorResponse("error deregistering rotation job: %s", err), nil
		}
	} else if conf.ShouldRegisterRotationJob() {
		cfgReq := &rotation.RotationJobConfigureRequest{
			MountPoint:       req.MountPoint,
			ReqPath:          req.Path,
			RotationSchedule: conf.RotationSchedule,
			RotationWindow:   conf.RotationWindow,
			RotationPeriod:   conf.RotationPeriod,
			RotationPolicy:   conf.RotationPolicy,
		}

		b.Logger().Debug("Registering rotation job", "mount", req.MountPoint+req.Path)
		if _, err = b.System().RegisterRotationJob(ctx, cfgReq); err != nil {
			return logical.ErrorResponse("error registering rotation job: %s", err), nil
		}
	}

	entry, err := logical.StorageEntryJSON("config", conf)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	conf, err := b.configEntry(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if conf == nil {
		return nil, nil
	}

	configData := map[string]interface{}{}
	conf.PopulateAutomatedRotationData(configData)

	if conf.HasNonzeroRotationValues() {
		resp, err := b.System().GetRotationInformation(ctx, &rotation.RotationInfoRequest{ReqPath: req.Path})
		if err != nil {
			return nil, err
		}
		if resp != nil {
			configData["expire_time"] = resp.NextVaultRotation.Unix()
			configData["creation_time"] = resp.LastVaultRotation.Unix()
			configData["ttl"] = int64(resp.TTL)
		}
	}

	return &logical.Response{
		Data: configData,
	}, nil
}

func (b *backend) pathConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := req.Storage.Delete(ctx, "config"); err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *backend) pathConfigExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.configEntry(ctx, req.Storage)
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// Fetch the client configuration required to access the AWS API.
func (b *backend) configEntry(ctx context.Context, s logical.Storage) (*config, error) {
	entry, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result config
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

type config struct {
	FailRotate bool
	automatedrotationutil.AutomatedRotationParams
}
