package models
type Mutability int

const (
    READWRITE_MUTABILITY Mutability = iota
    READONLY_MUTABILITY
    IMMUTABLE_MUTABILITY
    WRITEONLY_MUTABILITY
)

func (i Mutability) String() string {
    return []string{"ReadWrite", "ReadOnly", "Immutable", "WriteOnly"}[i]
}
func ParseMutability(v string) (any, error) {
    result := READWRITE_MUTABILITY
    switch v {
        case "ReadWrite":
            result = READWRITE_MUTABILITY
        case "ReadOnly":
            result = READONLY_MUTABILITY
        case "Immutable":
            result = IMMUTABLE_MUTABILITY
        case "WriteOnly":
            result = WRITEONLY_MUTABILITY
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeMutability(values []Mutability) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Mutability) isMultiValue() bool {
    return false
}
