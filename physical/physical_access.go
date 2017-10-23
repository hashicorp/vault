package physical

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

func (p *PhysicalAccess) Put(entry *Entry) error {
	return p.physical.Put(entry)
}

func (p *PhysicalAccess) Get(key string) (*Entry, error) {
	return p.physical.Get(key)
}

func (p *PhysicalAccess) Delete(key string) error {
	return p.physical.Delete(key)
}

func (p *PhysicalAccess) List(prefix string) ([]string, error) {
	return p.physical.List(prefix)
}

func (p *PhysicalAccess) Purge() {
	if purgeable, ok := p.physical.(Purgable); ok {
		purgeable.Purge()
	}
}
