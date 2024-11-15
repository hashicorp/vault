package security
import (
    "math"
    "strings"
)
type AdditionalOptions int

const (
    NONE_ADDITIONALOPTIONS = 1
    TEAMSANDYAMMERCONVERSATIONS_ADDITIONALOPTIONS = 2
    CLOUDATTACHMENTS_ADDITIONALOPTIONS = 4
    ALLDOCUMENTVERSIONS_ADDITIONALOPTIONS = 8
    SUBFOLDERCONTENTS_ADDITIONALOPTIONS = 16
    LISTATTACHMENTS_ADDITIONALOPTIONS = 32
    UNKNOWNFUTUREVALUE_ADDITIONALOPTIONS = 64
)

func (i AdditionalOptions) String() string {
    var values []string
    options := []string{"none", "teamsAndYammerConversations", "cloudAttachments", "allDocumentVersions", "subfolderContents", "listAttachments", "unknownFutureValue"}
    for p := 0; p < 7; p++ {
        mantis := AdditionalOptions(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseAdditionalOptions(v string) (any, error) {
    var result AdditionalOptions
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "none":
                result |= NONE_ADDITIONALOPTIONS
            case "teamsAndYammerConversations":
                result |= TEAMSANDYAMMERCONVERSATIONS_ADDITIONALOPTIONS
            case "cloudAttachments":
                result |= CLOUDATTACHMENTS_ADDITIONALOPTIONS
            case "allDocumentVersions":
                result |= ALLDOCUMENTVERSIONS_ADDITIONALOPTIONS
            case "subfolderContents":
                result |= SUBFOLDERCONTENTS_ADDITIONALOPTIONS
            case "listAttachments":
                result |= LISTATTACHMENTS_ADDITIONALOPTIONS
            case "unknownFutureValue":
                result |= UNKNOWNFUTUREVALUE_ADDITIONALOPTIONS
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeAdditionalOptions(values []AdditionalOptions) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AdditionalOptions) isMultiValue() bool {
    return true
}
