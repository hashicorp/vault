package callrecords
type VideoCodec int

const (
    UNKNOWN_VIDEOCODEC VideoCodec = iota
    INVALID_VIDEOCODEC
    AV1_VIDEOCODEC
    H263_VIDEOCODEC
    H264_VIDEOCODEC
    H264S_VIDEOCODEC
    H264UC_VIDEOCODEC
    H265_VIDEOCODEC
    RTVC1_VIDEOCODEC
    RTVIDEO_VIDEOCODEC
    XRTVC1_VIDEOCODEC
    UNKNOWNFUTUREVALUE_VIDEOCODEC
)

func (i VideoCodec) String() string {
    return []string{"unknown", "invalid", "av1", "h263", "h264", "h264s", "h264uc", "h265", "rtvc1", "rtVideo", "xrtvc1", "unknownFutureValue"}[i]
}
func ParseVideoCodec(v string) (any, error) {
    result := UNKNOWN_VIDEOCODEC
    switch v {
        case "unknown":
            result = UNKNOWN_VIDEOCODEC
        case "invalid":
            result = INVALID_VIDEOCODEC
        case "av1":
            result = AV1_VIDEOCODEC
        case "h263":
            result = H263_VIDEOCODEC
        case "h264":
            result = H264_VIDEOCODEC
        case "h264s":
            result = H264S_VIDEOCODEC
        case "h264uc":
            result = H264UC_VIDEOCODEC
        case "h265":
            result = H265_VIDEOCODEC
        case "rtvc1":
            result = RTVC1_VIDEOCODEC
        case "rtVideo":
            result = RTVIDEO_VIDEOCODEC
        case "xrtvc1":
            result = XRTVC1_VIDEOCODEC
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_VIDEOCODEC
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeVideoCodec(values []VideoCodec) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i VideoCodec) isMultiValue() bool {
    return false
}
