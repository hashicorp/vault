// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package monday

import "encoding/json"

// statusColumnID is the ID of the status column that an initial status is
// written to.
const statusColumnID = "status"

// columnValuesJSON is the GraphQL JSON scalar used for the column_values
// argument of the create_item and create_subitem mutations. Implements
// graphql.GraphQLType so that hasura emits "JSON!" in the operation variable
// declaration rather than inferring "String!" from the underlying Go string
// value.
type columnValuesJSON string

// GetGraphQLType implements graphql.GraphQLType.
func (columnValuesJSON) GetGraphQLType() string { return "JSON" }

// statusColumnValues builds the column_values payload that sets the initial
// status label on a new item. An empty status yields an empty JSON object so
// that no columns are set. The label is passed through exactly as given;
// Monday.com labels are case-sensitive.
func statusColumnValues(status string) columnValuesJSON {
	if status == "" {
		return columnValuesJSON("{}")
	}

	// Marshaling a map of strings cannot fail.
	b, _ := json.Marshal(map[string]map[string]string{
		statusColumnID: {"label": status},
	})
	return columnValuesJSON(b)
}
