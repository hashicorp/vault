package linodego

import (
	"encoding/json"
)

type FilterOperator string

const (
	Eq         FilterOperator = "+eq"
	Neq        FilterOperator = "+neq"
	Gt         FilterOperator = "+gt"
	Gte        FilterOperator = "+gte"
	Lt         FilterOperator = "+lt"
	Lte        FilterOperator = "+lte"
	Contains   FilterOperator = "+contains"
	Ascending                 = "asc"
	Descending                = "desc"
)

type FilterNode interface {
	Key() string
	JSONValueSegment() any
}

type Filter struct {
	// Operator is the logic for all Children nodes ("+and"/"+or")
	Operator string
	Children []FilterNode
	// OrderBy is the field you want to order your results by (ex: "+order_by": "class")
	OrderBy string
	// Order is the direction in which to order the results ("+order": "asc"/"desc")
	Order string
}

func (f *Filter) AddField(op FilterOperator, key string, value any) {
	f.Children = append(f.Children, &Comp{key, op, value})
}

func (f *Filter) MarshalJSON() ([]byte, error) {
	result := make(map[string]any)

	if f.OrderBy != "" {
		result["+order_by"] = f.OrderBy
	}

	if f.Order != "" {
		result["+order"] = f.Order
	}

	if f.Operator == "" {
		for _, c := range f.Children {
			result[c.Key()] = c.JSONValueSegment()
		}

		return json.Marshal(result)
	}

	fields := make([]map[string]any, len(f.Children))
	for i, c := range f.Children {
		fields[i] = map[string]any{
			c.Key(): c.JSONValueSegment(),
		}
	}

	result[f.Operator] = fields

	return json.Marshal(result)
}

type Comp struct {
	Column   string
	Operator FilterOperator
	Value    any
}

func (c *Comp) Key() string {
	return c.Column
}

func (c *Comp) JSONValueSegment() any {
	if c.Operator == Eq {
		return c.Value
	}

	return map[string]any{
		string(c.Operator): c.Value,
	}
}

func Or(order string, orderBy string, nodes ...FilterNode) *Filter {
	return &Filter{"+or", nodes, orderBy, order}
}

func And(order string, orderBy string, nodes ...FilterNode) *Filter {
	return &Filter{"+and", nodes, orderBy, order}
}
