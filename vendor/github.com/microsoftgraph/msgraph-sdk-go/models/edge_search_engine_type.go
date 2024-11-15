package models
// Allows IT admind to set a predefined default search engine for MDM-Controlled devices
type EdgeSearchEngineType int

const (
    // Uses factory settings of Edge to assign the default search engine as per the user market
    DEFAULT_EDGESEARCHENGINETYPE EdgeSearchEngineType = iota
    // Sets Bing as the default search engine
    BING_EDGESEARCHENGINETYPE
)

func (i EdgeSearchEngineType) String() string {
    return []string{"default", "bing"}[i]
}
func ParseEdgeSearchEngineType(v string) (any, error) {
    result := DEFAULT_EDGESEARCHENGINETYPE
    switch v {
        case "default":
            result = DEFAULT_EDGESEARCHENGINETYPE
        case "bing":
            result = BING_EDGESEARCHENGINETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeEdgeSearchEngineType(values []EdgeSearchEngineType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i EdgeSearchEngineType) isMultiValue() bool {
    return false
}
