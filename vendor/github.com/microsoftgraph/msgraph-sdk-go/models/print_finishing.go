package models
type PrintFinishing int

const (
    NONE_PRINTFINISHING PrintFinishing = iota
    STAPLE_PRINTFINISHING
    PUNCH_PRINTFINISHING
    COVER_PRINTFINISHING
    BIND_PRINTFINISHING
    SADDLESTITCH_PRINTFINISHING
    STITCHEDGE_PRINTFINISHING
    STAPLETOPLEFT_PRINTFINISHING
    STAPLEBOTTOMLEFT_PRINTFINISHING
    STAPLETOPRIGHT_PRINTFINISHING
    STAPLEBOTTOMRIGHT_PRINTFINISHING
    STITCHLEFTEDGE_PRINTFINISHING
    STITCHTOPEDGE_PRINTFINISHING
    STITCHRIGHTEDGE_PRINTFINISHING
    STITCHBOTTOMEDGE_PRINTFINISHING
    STAPLEDUALLEFT_PRINTFINISHING
    STAPLEDUALTOP_PRINTFINISHING
    STAPLEDUALRIGHT_PRINTFINISHING
    STAPLEDUALBOTTOM_PRINTFINISHING
    UNKNOWNFUTUREVALUE_PRINTFINISHING
    STAPLETRIPLELEFT_PRINTFINISHING
    STAPLETRIPLETOP_PRINTFINISHING
    STAPLETRIPLERIGHT_PRINTFINISHING
    STAPLETRIPLEBOTTOM_PRINTFINISHING
    BINDLEFT_PRINTFINISHING
    BINDTOP_PRINTFINISHING
    BINDRIGHT_PRINTFINISHING
    BINDBOTTOM_PRINTFINISHING
    FOLDACCORDION_PRINTFINISHING
    FOLDDOUBLEGATE_PRINTFINISHING
    FOLDGATE_PRINTFINISHING
    FOLDHALF_PRINTFINISHING
    FOLDHALFZ_PRINTFINISHING
    FOLDLEFTGATE_PRINTFINISHING
    FOLDLETTER_PRINTFINISHING
    FOLDPARALLEL_PRINTFINISHING
    FOLDPOSTER_PRINTFINISHING
    FOLDRIGHTGATE_PRINTFINISHING
    FOLDZ_PRINTFINISHING
    FOLDENGINEERINGZ_PRINTFINISHING
    PUNCHTOPLEFT_PRINTFINISHING
    PUNCHBOTTOMLEFT_PRINTFINISHING
    PUNCHTOPRIGHT_PRINTFINISHING
    PUNCHBOTTOMRIGHT_PRINTFINISHING
    PUNCHDUALLEFT_PRINTFINISHING
    PUNCHDUALTOP_PRINTFINISHING
    PUNCHDUALRIGHT_PRINTFINISHING
    PUNCHDUALBOTTOM_PRINTFINISHING
    PUNCHTRIPLELEFT_PRINTFINISHING
    PUNCHTRIPLETOP_PRINTFINISHING
    PUNCHTRIPLERIGHT_PRINTFINISHING
    PUNCHTRIPLEBOTTOM_PRINTFINISHING
    PUNCHQUADLEFT_PRINTFINISHING
    PUNCHQUADTOP_PRINTFINISHING
    PUNCHQUADRIGHT_PRINTFINISHING
    PUNCHQUADBOTTOM_PRINTFINISHING
    FOLD_PRINTFINISHING
    TRIM_PRINTFINISHING
    BALE_PRINTFINISHING
    BOOKLETMAKER_PRINTFINISHING
    COAT_PRINTFINISHING
    LAMINATE_PRINTFINISHING
    TRIMAFTERPAGES_PRINTFINISHING
    TRIMAFTERDOCUMENTS_PRINTFINISHING
    TRIMAFTERCOPIES_PRINTFINISHING
    TRIMAFTERJOB_PRINTFINISHING
)

func (i PrintFinishing) String() string {
    return []string{"none", "staple", "punch", "cover", "bind", "saddleStitch", "stitchEdge", "stapleTopLeft", "stapleBottomLeft", "stapleTopRight", "stapleBottomRight", "stitchLeftEdge", "stitchTopEdge", "stitchRightEdge", "stitchBottomEdge", "stapleDualLeft", "stapleDualTop", "stapleDualRight", "stapleDualBottom", "unknownFutureValue", "stapleTripleLeft", "stapleTripleTop", "stapleTripleRight", "stapleTripleBottom", "bindLeft", "bindTop", "bindRight", "bindBottom", "foldAccordion", "foldDoubleGate", "foldGate", "foldHalf", "foldHalfZ", "foldLeftGate", "foldLetter", "foldParallel", "foldPoster", "foldRightGate", "foldZ", "foldEngineeringZ", "punchTopLeft", "punchBottomLeft", "punchTopRight", "punchBottomRight", "punchDualLeft", "punchDualTop", "punchDualRight", "punchDualBottom", "punchTripleLeft", "punchTripleTop", "punchTripleRight", "punchTripleBottom", "punchQuadLeft", "punchQuadTop", "punchQuadRight", "punchQuadBottom", "fold", "trim", "bale", "bookletMaker", "coat", "laminate", "trimAfterPages", "trimAfterDocuments", "trimAfterCopies", "trimAfterJob"}[i]
}
func ParsePrintFinishing(v string) (any, error) {
    result := NONE_PRINTFINISHING
    switch v {
        case "none":
            result = NONE_PRINTFINISHING
        case "staple":
            result = STAPLE_PRINTFINISHING
        case "punch":
            result = PUNCH_PRINTFINISHING
        case "cover":
            result = COVER_PRINTFINISHING
        case "bind":
            result = BIND_PRINTFINISHING
        case "saddleStitch":
            result = SADDLESTITCH_PRINTFINISHING
        case "stitchEdge":
            result = STITCHEDGE_PRINTFINISHING
        case "stapleTopLeft":
            result = STAPLETOPLEFT_PRINTFINISHING
        case "stapleBottomLeft":
            result = STAPLEBOTTOMLEFT_PRINTFINISHING
        case "stapleTopRight":
            result = STAPLETOPRIGHT_PRINTFINISHING
        case "stapleBottomRight":
            result = STAPLEBOTTOMRIGHT_PRINTFINISHING
        case "stitchLeftEdge":
            result = STITCHLEFTEDGE_PRINTFINISHING
        case "stitchTopEdge":
            result = STITCHTOPEDGE_PRINTFINISHING
        case "stitchRightEdge":
            result = STITCHRIGHTEDGE_PRINTFINISHING
        case "stitchBottomEdge":
            result = STITCHBOTTOMEDGE_PRINTFINISHING
        case "stapleDualLeft":
            result = STAPLEDUALLEFT_PRINTFINISHING
        case "stapleDualTop":
            result = STAPLEDUALTOP_PRINTFINISHING
        case "stapleDualRight":
            result = STAPLEDUALRIGHT_PRINTFINISHING
        case "stapleDualBottom":
            result = STAPLEDUALBOTTOM_PRINTFINISHING
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_PRINTFINISHING
        case "stapleTripleLeft":
            result = STAPLETRIPLELEFT_PRINTFINISHING
        case "stapleTripleTop":
            result = STAPLETRIPLETOP_PRINTFINISHING
        case "stapleTripleRight":
            result = STAPLETRIPLERIGHT_PRINTFINISHING
        case "stapleTripleBottom":
            result = STAPLETRIPLEBOTTOM_PRINTFINISHING
        case "bindLeft":
            result = BINDLEFT_PRINTFINISHING
        case "bindTop":
            result = BINDTOP_PRINTFINISHING
        case "bindRight":
            result = BINDRIGHT_PRINTFINISHING
        case "bindBottom":
            result = BINDBOTTOM_PRINTFINISHING
        case "foldAccordion":
            result = FOLDACCORDION_PRINTFINISHING
        case "foldDoubleGate":
            result = FOLDDOUBLEGATE_PRINTFINISHING
        case "foldGate":
            result = FOLDGATE_PRINTFINISHING
        case "foldHalf":
            result = FOLDHALF_PRINTFINISHING
        case "foldHalfZ":
            result = FOLDHALFZ_PRINTFINISHING
        case "foldLeftGate":
            result = FOLDLEFTGATE_PRINTFINISHING
        case "foldLetter":
            result = FOLDLETTER_PRINTFINISHING
        case "foldParallel":
            result = FOLDPARALLEL_PRINTFINISHING
        case "foldPoster":
            result = FOLDPOSTER_PRINTFINISHING
        case "foldRightGate":
            result = FOLDRIGHTGATE_PRINTFINISHING
        case "foldZ":
            result = FOLDZ_PRINTFINISHING
        case "foldEngineeringZ":
            result = FOLDENGINEERINGZ_PRINTFINISHING
        case "punchTopLeft":
            result = PUNCHTOPLEFT_PRINTFINISHING
        case "punchBottomLeft":
            result = PUNCHBOTTOMLEFT_PRINTFINISHING
        case "punchTopRight":
            result = PUNCHTOPRIGHT_PRINTFINISHING
        case "punchBottomRight":
            result = PUNCHBOTTOMRIGHT_PRINTFINISHING
        case "punchDualLeft":
            result = PUNCHDUALLEFT_PRINTFINISHING
        case "punchDualTop":
            result = PUNCHDUALTOP_PRINTFINISHING
        case "punchDualRight":
            result = PUNCHDUALRIGHT_PRINTFINISHING
        case "punchDualBottom":
            result = PUNCHDUALBOTTOM_PRINTFINISHING
        case "punchTripleLeft":
            result = PUNCHTRIPLELEFT_PRINTFINISHING
        case "punchTripleTop":
            result = PUNCHTRIPLETOP_PRINTFINISHING
        case "punchTripleRight":
            result = PUNCHTRIPLERIGHT_PRINTFINISHING
        case "punchTripleBottom":
            result = PUNCHTRIPLEBOTTOM_PRINTFINISHING
        case "punchQuadLeft":
            result = PUNCHQUADLEFT_PRINTFINISHING
        case "punchQuadTop":
            result = PUNCHQUADTOP_PRINTFINISHING
        case "punchQuadRight":
            result = PUNCHQUADRIGHT_PRINTFINISHING
        case "punchQuadBottom":
            result = PUNCHQUADBOTTOM_PRINTFINISHING
        case "fold":
            result = FOLD_PRINTFINISHING
        case "trim":
            result = TRIM_PRINTFINISHING
        case "bale":
            result = BALE_PRINTFINISHING
        case "bookletMaker":
            result = BOOKLETMAKER_PRINTFINISHING
        case "coat":
            result = COAT_PRINTFINISHING
        case "laminate":
            result = LAMINATE_PRINTFINISHING
        case "trimAfterPages":
            result = TRIMAFTERPAGES_PRINTFINISHING
        case "trimAfterDocuments":
            result = TRIMAFTERDOCUMENTS_PRINTFINISHING
        case "trimAfterCopies":
            result = TRIMAFTERCOPIES_PRINTFINISHING
        case "trimAfterJob":
            result = TRIMAFTERJOB_PRINTFINISHING
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializePrintFinishing(values []PrintFinishing) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i PrintFinishing) isMultiValue() bool {
    return false
}
