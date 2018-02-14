package physical

import "context"

// PhysicalAccess is a wrapper around physical.Backend that allows Core to
// expose its physical storage operations through PhysicalAccess() while
// restricting the ability to modify Core.physical itself.
type PhysicalAccess struct {
	physical Backend
}

var _ Backend = (*PhysicalAccess)(nil)

func NewPhysicalAccess(physical Backend) *PhysicalAccess {
	return &PhysicalAccess{physical: physical}
}

func (p *PhysicalAccess) Put(ctx context.Context, entry *Entry) error {
	return p.physical.Put(ctx, entry)
}

func (p *PhysicalAccess) Get(ctx context.Context, key string) (*Entry, error) {
	return p.physical.Get(ctx, key)
}

func (p *PhysicalAccess) Delete(ctx context.Context, key string) error {
	return p.physical.Delete(ctx, key)
}

func (p *PhysicalAccess) List(ctx context.Context, prefix string) ([]string, error) {
	return p.physical.List(ctx, prefix)
}

func (p *PhysicalAccess) Purge(ctx context.Context) {
	if purgeable, ok := p.physical.(ToggleablePurgemonster); ok {
		purgeable.Purge(ctx)
	}
}
