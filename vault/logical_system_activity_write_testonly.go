// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build testonly

package vault

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/activity"
	"github.com/hashicorp/vault/vault/activity/generation"
	"google.golang.org/protobuf/encoding/protojson"
)

const helpText = "Create activity log data for testing purposes"

func (b *SystemBackend) activityWritePath() *framework.Path {
	return &framework.Path{
		Pattern:         "internal/counters/activity/write$",
		HelpDescription: helpText,
		HelpSynopsis:    helpText,
		Fields: map[string]*framework.FieldSchema{
			"input": {
				Type:        framework.TypeString,
				Description: "JSON input for generating mock data",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.handleActivityWriteData,
				Summary:  "Write activity log data",
			},
		},
	}
}

func (b *SystemBackend) handleActivityWriteData(ctx context.Context, request *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	json := data.Get("input")
	input := &generation.ActivityLogMockInput{}
	err := protojson.Unmarshal([]byte(json.(string)), input)
	if err != nil {
		return logical.ErrorResponse("Invalid input data: %s", err), logical.ErrInvalidRequest
	}
	if len(input.Write) == 0 {
		return logical.ErrorResponse("Missing required \"write\" values"), logical.ErrInvalidRequest
	}
	if len(input.Data) == 0 {
		return logical.ErrorResponse("Missing required \"data\" values"), logical.ErrInvalidRequest
	}
	return nil, nil
}

// singleMonthActivityClients holds a single month's client IDs, in the order they were seen
type singleMonthActivityClients struct {
	// clients are indexed by ID
	clients []*activity.EntityRecord
}

// multipleMonthsActivityClients holds multiple month's data
type multipleMonthsActivityClients struct {
	// months are in order, with month 0 being the current month and index 1 being 1 month ago
	months []*singleMonthActivityClients
}

// addNewClients generates clients according to the given parameters, and adds them to the month
// the client will always have the mountAccessor as its mount accessor
func (s *singleMonthActivityClients) addNewClients(c *generation.Client, mountAccessor string) error {
	count := 1
	if c.Count > 1 {
		count = int(c.Count)
	}
	for i := 0; i < count; i++ {
		record := &activity.EntityRecord{
			ClientID:      c.Id,
			NamespaceID:   c.Namespace,
			NonEntity:     c.NonEntity,
			MountAccessor: mountAccessor,
		}
		if record.ClientID == "" {
			var err error
			record.ClientID, err = uuid.GenerateUUID()
			if err != nil {
				return err
			}
		}
		s.clients = append(s.clients, record)
	}
	return nil
}

// processMonth populates a month of client data
func (m *multipleMonthsActivityClients) processMonth(ctx context.Context, core *Core, month *generation.Data) error {
	if month.GetAll() == nil {
		return errors.New("segmented monthly data is not yet supported")
	}

	// default to using the root namespace and the first mount on the root namespace
	mounts, err := core.ListMounts()
	if err != nil {
		return err
	}
	defaultMountAccessorRootNS := ""
	for _, mount := range mounts {
		if mount.NamespaceID == namespace.RootNamespaceID {
			defaultMountAccessorRootNS = mount.Accessor
			break
		}
	}
	addingTo := m.months[month.GetMonthsAgo()]

	for _, clients := range month.GetAll().Clients {
		if clients.Repeated || clients.RepeatedFromMonth > 0 {
			return errors.New("repeated clients are not yet supported")
		}

		if clients.Namespace == "" {
			clients.Namespace = namespace.RootNamespaceID
		}

		// verify that the namespace exists
		ns, err := core.NamespaceByID(ctx, clients.Namespace)
		if err != nil {
			return err
		}

		// verify that the mount exists
		if clients.Mount != "" {
			nctx := namespace.ContextWithNamespace(ctx, ns)
			mountEntry := core.router.MatchingMountEntry(nctx, clients.Mount)
			if mountEntry == nil {
				return fmt.Errorf("unable to find matching mount in namespace %s", clients.Namespace)
			}
		}

		mountAccessor := defaultMountAccessorRootNS
		if clients.Namespace != namespace.RootNamespaceID && clients.Mount == "" {
			// if we're not using the root namespace, find a mount on the namespace that we are using
			found := false
			for _, mount := range mounts {
				if mount.NamespaceID == clients.Namespace {
					mountAccessor = mount.Accessor
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("unable to find matching mount in namespace %s", clients.Namespace)
			}
		}
		err = addingTo.addNewClients(clients, mountAccessor)
		if err != nil {
			return err
		}
	}
	return nil
}

func newMultipleMonthsActivityClients(numberOfMonths int) *multipleMonthsActivityClients {
	m := &multipleMonthsActivityClients{
		months: make([]*singleMonthActivityClients, numberOfMonths),
	}
	for i := 0; i < numberOfMonths; i++ {
		m.months[i] = new(singleMonthActivityClients)
	}
	return m
}
