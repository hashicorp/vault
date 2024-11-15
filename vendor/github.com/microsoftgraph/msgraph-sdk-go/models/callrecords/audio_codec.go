package callrecords
type AudioCodec int

const (
    UNKNOWN_AUDIOCODEC AudioCodec = iota
    INVALID_AUDIOCODEC
    CN_AUDIOCODEC
    PCMA_AUDIOCODEC
    PCMU_AUDIOCODEC
    AMRWIDE_AUDIOCODEC
    G722_AUDIOCODEC
    G7221_AUDIOCODEC
    G7221C_AUDIOCODEC
    G729_AUDIOCODEC
    MULTICHANNELAUDIO_AUDIOCODEC
    MUCHV2_AUDIOCODEC
    OPUS_AUDIOCODEC
    SATIN_AUDIOCODEC
    SATINFULLBAND_AUDIOCODEC
    RTAUDIO8_AUDIOCODEC
    RTAUDIO16_AUDIOCODEC
    SILK_AUDIOCODEC
    SILKNARROW_AUDIOCODEC
    SILKWIDE_AUDIOCODEC
    SIREN_AUDIOCODEC
    XMSRTA_AUDIOCODEC
    UNKNOWNFUTUREVALUE_AUDIOCODEC
)

func (i AudioCodec) String() string {
    return []string{"unknown", "invalid", "cn", "pcma", "pcmu", "amrWide", "g722", "g7221", "g7221c", "g729", "multiChannelAudio", "muchv2", "opus", "satin", "satinFullband", "rtAudio8", "rtAudio16", "silk", "silkNarrow", "silkWide", "siren", "xmsRta", "unknownFutureValue"}[i]
}
func ParseAudioCodec(v string) (any, error) {
    result := UNKNOWN_AUDIOCODEC
    switch v {
        case "unknown":
            result = UNKNOWN_AUDIOCODEC
        case "invalid":
            result = INVALID_AUDIOCODEC
        case "cn":
            result = CN_AUDIOCODEC
        case "pcma":
            result = PCMA_AUDIOCODEC
        case "pcmu":
            result = PCMU_AUDIOCODEC
        case "amrWide":
            result = AMRWIDE_AUDIOCODEC
        case "g722":
            result = G722_AUDIOCODEC
        case "g7221":
            result = G7221_AUDIOCODEC
        case "g7221c":
            result = G7221C_AUDIOCODEC
        case "g729":
            result = G729_AUDIOCODEC
        case "multiChannelAudio":
            result = MULTICHANNELAUDIO_AUDIOCODEC
        case "muchv2":
            result = MUCHV2_AUDIOCODEC
        case "opus":
            result = OPUS_AUDIOCODEC
        case "satin":
            result = SATIN_AUDIOCODEC
        case "satinFullband":
            result = SATINFULLBAND_AUDIOCODEC
        case "rtAudio8":
            result = RTAUDIO8_AUDIOCODEC
        case "rtAudio16":
            result = RTAUDIO16_AUDIOCODEC
        case "silk":
            result = SILK_AUDIOCODEC
        case "silkNarrow":
            result = SILKNARROW_AUDIOCODEC
        case "silkWide":
            result = SILKWIDE_AUDIOCODEC
        case "siren":
            result = SIREN_AUDIOCODEC
        case "xmsRta":
            result = XMSRTA_AUDIOCODEC
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_AUDIOCODEC
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeAudioCodec(values []AudioCodec) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i AudioCodec) isMultiValue() bool {
    return false
}
