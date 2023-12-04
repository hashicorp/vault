package vault

import "context"

type UICustomMessage interface {
	// Create
	Create(context.Context, uicustommessages.Entry)
	// Delete
	// Get
	// Update
	// Find
}
