package tfe

type TagList struct {
	*Pagination
	Items []*Tag
}

// Tag is owned by an organization and applied to workspaces. Used for grouping and search.
type Tag struct {
	ID   string `jsonapi:"primary,tags"`
	Name string `jsonapi:"attr,name,omitempty"`
}
