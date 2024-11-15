package models
type X509CertificateRuleType int

const (
    ISSUERSUBJECT_X509CERTIFICATERULETYPE X509CertificateRuleType = iota
    POLICYOID_X509CERTIFICATERULETYPE
    UNKNOWNFUTUREVALUE_X509CERTIFICATERULETYPE
    ISSUERSUBJECTANDPOLICYOID_X509CERTIFICATERULETYPE
)

func (i X509CertificateRuleType) String() string {
    return []string{"issuerSubject", "policyOID", "unknownFutureValue", "issuerSubjectAndPolicyOID"}[i]
}
func ParseX509CertificateRuleType(v string) (any, error) {
    result := ISSUERSUBJECT_X509CERTIFICATERULETYPE
    switch v {
        case "issuerSubject":
            result = ISSUERSUBJECT_X509CERTIFICATERULETYPE
        case "policyOID":
            result = POLICYOID_X509CERTIFICATERULETYPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_X509CERTIFICATERULETYPE
        case "issuerSubjectAndPolicyOID":
            result = ISSUERSUBJECTANDPOLICYOID_X509CERTIFICATERULETYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeX509CertificateRuleType(values []X509CertificateRuleType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i X509CertificateRuleType) isMultiValue() bool {
    return false
}
