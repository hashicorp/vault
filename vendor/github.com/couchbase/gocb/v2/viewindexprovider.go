package gocb

type viewIndexProvider interface {
	GetDesignDocument(name string, namespace DesignDocumentNamespace, opts *GetDesignDocumentOptions) (*DesignDocument, error)
	GetAllDesignDocuments(namespace DesignDocumentNamespace, opts *GetAllDesignDocumentsOptions) ([]DesignDocument, error)
	UpsertDesignDocument(ddoc DesignDocument, namespace DesignDocumentNamespace, opts *UpsertDesignDocumentOptions) error
	DropDesignDocument(name string, namespace DesignDocumentNamespace, opts *DropDesignDocumentOptions) error
	PublishDesignDocument(name string, opts *PublishDesignDocumentOptions) error
}
