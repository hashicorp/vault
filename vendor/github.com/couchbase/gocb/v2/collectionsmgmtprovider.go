package gocb

type collectionsManagementProvider interface {
	GetAllScopes(opts *GetAllScopesOptions) ([]ScopeSpec, error)
	CreateCollection(scopeName string, collectionName string, settings *CreateCollectionSettings, opts *CreateCollectionOptions) error
	UpdateCollection(scopeName string, collectionName string, settings UpdateCollectionSettings, opts *UpdateCollectionOptions) error
	DropCollection(scopeName string, collectionName string, opts *DropCollectionOptions) error
	CreateScope(scopeName string, opts *CreateScopeOptions) error
	DropScope(scopeName string, opts *DropScopeOptions) error
}
