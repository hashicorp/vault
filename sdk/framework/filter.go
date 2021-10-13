package framework

import (
	"context"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/ryanuber/go-glob"
)

// GlobListFilter wraps an OperationFunc with an optional filter which excludes listed entries
// which don't match a glob style pattern
func GlobListFilter(fieldName string, callback OperationFunc) OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
		resp, err := callback(ctx, req, data)
		if err != nil {
			return nil, err
		}

		if keys, ok := resp.Data["keys"]; ok {
			if entries, ok := keys.([]string); ok {
				filter, ok := data.GetOk(fieldName)
				if ok && filter != "" && filter != "*" {
					var filteredEntries []string
					for _, e := range entries {
						if glob.Glob(filter.(string), e) {
							filteredEntries = append(filteredEntries, e)
						}
					}
					resp.Data["keys"] = filteredEntries
				}
			}
		}
		return resp, nil
	}
}
