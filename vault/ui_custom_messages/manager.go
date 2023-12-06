package uicustommessages

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// StoragePrefix is the prefix string to use to create the view of the
	// logical.Storage used by the Manager struct.
	StoragePrefix string = "ui/custom-messages/namespaces/"

	// MaximumMessageCountPerNamespace is the maximum number of custom messages
	// that can be stored for any namespace. This constraint is in place to
	// restrict the total number of custom messages in the system to make sure
	// that this doesn't become a performance drain.
	MaximumMessageCountPerNamespace int = 100
)

// Manager is a struct that provides methods to manage messages stored in a
// logical.Storage.
type Manager struct {
	view logical.Storage

	l sync.RWMutex
}

// NewManager creates a new Manager struct that has been fully initialized.
func NewManager(storage logical.Storage) *Manager {
	return &Manager{
		view: storage,
	}
}

// FindMessages handles getting a list of existing messages that match the
// criteria set in the provided FindFilter struct.
func (m *Manager) FindMessages(ctx context.Context, filters FindFilter) ([]Message, error) {
	nsList, err := getNamespacesToSearch(ctx, filters)
	if err != nil {
		return nil, err
	}

	results := make([]Message, 0)

	for _, ns := range nsList {
		entry, err := m.getEntryForNamespace(ctx, ns)
		if err != nil {
			return nil, err
		}

		results = append(results, entry.FindMessages(filters)...)
	}

	return results, nil
}

// CreateMessage handles adding the provided Message in the current namespace.
// A ID for the message is automatically generated and the ID field of the
// returned Message struct is set to this value. If the maximum number of
// messages already exists, then the message is not added.
func (m *Manager) CreateMessage(ctx context.Context, message Message) (*Message, error) {
	m.l.Lock()
	defer m.l.Unlock()

	entry, err := m.getEntry(ctx)
	if err != nil {
		return nil, err
	}

	err = entry.CreateMessage(&message)
	if err != nil {
		return nil, err
	}

	err = m.putEntry(ctx, entry)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

// ReadMessage handles retrieving the properties of the message in the current
// namespace specified by the id value from the logical.Storage.
func (m *Manager) ReadMessage(ctx context.Context, id string) (*Message, error) {
	entry, err := m.getEntry(ctx)
	if err != nil {
		return nil, err
	}

	message, ok := entry.Messages[id]
	if !ok {
		return nil, nil
	}

	return &message, nil
}

// UpdateMessage handles updating the message referenced by the provided
// Message struct in the current namespace with its content in the
// logical.Storage.
func (m *Manager) UpdateMessage(ctx context.Context, message Message) (*Message, error) {
	m.l.Lock()
	defer m.l.Unlock()

	entry, err := m.getEntry(ctx)
	if err != nil {
		return nil, err
	}

	found, err := entry.UpdateMessage(&message)
	if err != nil {
		return nil, err
	}

	if !found {
		// If no error occured but the specified message doesn't exist...
		return nil, nil
	}

	if err = m.putEntry(ctx, entry); err != nil {
		return nil, err
	}

	return &message, nil
}

// DeleteMessage handles deleting the message with the provided id value in the
// current namespace if it exists. The method updates the logical.Storage as
// well.
func (m *Manager) DeleteMessage(ctx context.Context, id string) error {
	m.l.Lock()
	defer m.l.Unlock()

	entry, err := m.getEntry(ctx)
	if err != nil {
		return err
	}

	delete(entry.Messages, id)

	return m.putEntry(ctx, entry)
}

// getEntry is a helper method that retrieves the current namespace from the
// context.Context and uses it to call getEntryForNamespace.
func (m *Manager) getEntry(ctx context.Context) (*Entry, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	return m.getEntryForNamespace(ctx, ns)
}

// getEntryForNamespace takes care of retrieving the logical.StorageEntry that
// corresponds to the provided namespace.Namespace. The logical.StorageEntry is
// then used to build an Entry struct.
func (m *Manager) getEntryForNamespace(ctx context.Context, ns *namespace.Namespace) (*Entry, error) {
	if ns == nil {
		return nil, errors.New("missing namespace")
	}

	storageEntry, err := m.view.Get(ctx, ns.ID)
	if err != nil {
		return nil, err
	}

	if storageEntry == nil {
		return &Entry{
			Messages: make(map[string]Message),
		}, nil
	}

	var entry *Entry = new(Entry)
	if err := storageEntry.DecodeJSON(entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// putEntry takes care of determining the current namespace from the
// context.Context and then marshalling the provided Entry pointer to a slice
// of byte to create the appropriate logical.StorageEntry and then storing it
// in the logical.Storage.
func (m *Manager) putEntry(ctx context.Context, entry *Entry) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	value, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	storageEntry := &logical.StorageEntry{
		Key:   ns.ID,
		Value: value,
	}

	return m.view.Put(ctx, storageEntry)
}

// getNamespacesToSearch builds a slice of pointers to namespace.Namespace
// struct that will be walked by the (*Manager).FindMessage method above.
// This function handles the complexity of gathering all of the applicable
// namespaces depending on the namespace set in the context and whether the
// IncludeAncestors criterion is set to true in the provided FindFilter struct.
func getNamespacesToSearch(ctx context.Context, filters FindFilter) ([]*namespace.Namespace, error) {
	var nsList []*namespace.Namespace

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Add the current namespace based on the context.Context to nsList.
	nsList = append(nsList, ns)

	//if filters.IncludeAncestors {
	// Add the parent, grand-parent, etc... namespaces all the way back up
	// to the root namespace to nsList.
	//}

	return nsList, nil
}
