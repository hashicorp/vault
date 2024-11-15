package models
type PrinterProcessingState int

const (
    UNKNOWN_PRINTERPROCESSINGSTATE PrinterProcessingState = iota
    IDLE_PRINTERPROCESSINGSTATE
    PROCESSING_PRINTERPROCESSINGSTATE
    STOPPED_PRINTERPROCESSINGSTATE
    UNKNOWNFUTUREVALUE_PRINTERPROCESSINGSTATE
)

func (i PrinterProcessingState) String() string {
    return []string{"unknown", "idle", "processing", "stopped", "unknownFutureValue"}[i]
}
func ParsePrinterProcessingState(v string) (any, error) {
    result := UNKNOWN_PRINTERPROCESSINGSTATE
    switch v {
        case "unknown":
            result = UNKNOWN_PRINTERPROCESSINGSTATE
        case "idle":
            result = IDLE_PRINTERPROCESSINGSTATE
        case "processing":
            result = PROCESSING_PRINTERPROCESSINGSTATE
        case "stopped":
            result = STOPPED_PRINTERPROCESSINGSTATE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRINTERPROCESSINGSTATE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrinterProcessingState(values []PrinterProcessingState) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrinterProcessingState) isMultiValue() bool {
    return false
}
