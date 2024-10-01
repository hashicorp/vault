// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"

	uicustommessages "github.com/hashicorp/vault/vault/ui_custom_messages"
)

// CustomMessagesManager is the interface used by the vault package when
// interacting with a uicustommessages.Manager instance.
type CustomMessagesManager interface {
	FindMessages(context.Context, uicustommessages.FindFilter) ([]uicustommessages.Message, error)
	AddMessage(context.Context, uicustommessages.Message) (*uicustommessages.Message, error)
	ReadMessage(context.Context, string) (*uicustommessages.Message, error)
	UpdateMessage(context.Context, uicustommessages.Message) (*uicustommessages.Message, error)
	DeleteMessage(context.Context, string) error
}
