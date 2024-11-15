package models
import (
    "math"
    "strings"
)
type SynchronizationJobRestartScope int

const (
    NONE_SYNCHRONIZATIONJOBRESTARTSCOPE = 1
    CONNECTORDATASTORE_SYNCHRONIZATIONJOBRESTARTSCOPE = 2
    ESCROWS_SYNCHRONIZATIONJOBRESTARTSCOPE = 4
    WATERMARK_SYNCHRONIZATIONJOBRESTARTSCOPE = 8
    QUARANTINESTATE_SYNCHRONIZATIONJOBRESTARTSCOPE = 16
    FULL_SYNCHRONIZATIONJOBRESTARTSCOPE = 32
    FORCEDELETES_SYNCHRONIZATIONJOBRESTARTSCOPE = 64
)

func (i SynchronizationJobRestartScope) String() string {
    var values []string
    options := []string{"None", "ConnectorDataStore", "Escrows", "Watermark", "QuarantineState", "Full", "ForceDeletes"}
    for p := 0; p < 7; p++ {
        mantis := SynchronizationJobRestartScope(int(math.Pow(2, float64(p))))
        if i&mantis == mantis {
            values = append(values, options[p])
        }
    }
    return strings.Join(values, ",")
}
func ParseSynchronizationJobRestartScope(v string) (any, error) {
    var result SynchronizationJobRestartScope
    values := strings.Split(v, ",")
    for _, str := range values {
        switch str {
            case "None":
                result |= NONE_SYNCHRONIZATIONJOBRESTARTSCOPE
            case "ConnectorDataStore":
                result |= CONNECTORDATASTORE_SYNCHRONIZATIONJOBRESTARTSCOPE
            case "Escrows":
                result |= ESCROWS_SYNCHRONIZATIONJOBRESTARTSCOPE
            case "Watermark":
                result |= WATERMARK_SYNCHRONIZATIONJOBRESTARTSCOPE
            case "QuarantineState":
                result |= QUARANTINESTATE_SYNCHRONIZATIONJOBRESTARTSCOPE
            case "Full":
                result |= FULL_SYNCHRONIZATIONJOBRESTARTSCOPE
            case "ForceDeletes":
                result |= FORCEDELETES_SYNCHRONIZATIONJOBRESTARTSCOPE
            default:
                return nil, nil
        }
    }
    return &result, nil
}
func SerializeSynchronizationJobRestartScope(values []SynchronizationJobRestartScope) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SynchronizationJobRestartScope) isMultiValue() bool {
    return true
}
