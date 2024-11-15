package models
// Windows 10 Edition type.
type Windows10EditionType int

const (
    // Windows 10 Enterprise
    WINDOWS10ENTERPRISE_WINDOWS10EDITIONTYPE Windows10EditionType = iota
    // Windows 10 EnterpriseN
    WINDOWS10ENTERPRISEN_WINDOWS10EDITIONTYPE
    // Windows 10 Education
    WINDOWS10EDUCATION_WINDOWS10EDITIONTYPE
    // Windows 10 EducationN
    WINDOWS10EDUCATIONN_WINDOWS10EDITIONTYPE
    // Windows 10 Mobile Enterprise
    WINDOWS10MOBILEENTERPRISE_WINDOWS10EDITIONTYPE
    // Windows 10 Holographic Enterprise
    WINDOWS10HOLOGRAPHICENTERPRISE_WINDOWS10EDITIONTYPE
    // Windows 10 Professional
    WINDOWS10PROFESSIONAL_WINDOWS10EDITIONTYPE
    // Windows 10 ProfessionalN
    WINDOWS10PROFESSIONALN_WINDOWS10EDITIONTYPE
    // Windows 10 Professional Education
    WINDOWS10PROFESSIONALEDUCATION_WINDOWS10EDITIONTYPE
    // Windows 10 Professional EducationN
    WINDOWS10PROFESSIONALEDUCATIONN_WINDOWS10EDITIONTYPE
    // Windows 10 Professional for Workstations
    WINDOWS10PROFESSIONALWORKSTATION_WINDOWS10EDITIONTYPE
    // Windows 10 Professional for Workstations N
    WINDOWS10PROFESSIONALWORKSTATIONN_WINDOWS10EDITIONTYPE
)

func (i Windows10EditionType) String() string {
    return []string{"windows10Enterprise", "windows10EnterpriseN", "windows10Education", "windows10EducationN", "windows10MobileEnterprise", "windows10HolographicEnterprise", "windows10Professional", "windows10ProfessionalN", "windows10ProfessionalEducation", "windows10ProfessionalEducationN", "windows10ProfessionalWorkstation", "windows10ProfessionalWorkstationN"}[i]
}
func ParseWindows10EditionType(v string) (any, error) {
    result := WINDOWS10ENTERPRISE_WINDOWS10EDITIONTYPE
    switch v {
        case "windows10Enterprise":
            result = WINDOWS10ENTERPRISE_WINDOWS10EDITIONTYPE
        case "windows10EnterpriseN":
            result = WINDOWS10ENTERPRISEN_WINDOWS10EDITIONTYPE
        case "windows10Education":
            result = WINDOWS10EDUCATION_WINDOWS10EDITIONTYPE
        case "windows10EducationN":
            result = WINDOWS10EDUCATIONN_WINDOWS10EDITIONTYPE
        case "windows10MobileEnterprise":
            result = WINDOWS10MOBILEENTERPRISE_WINDOWS10EDITIONTYPE
        case "windows10HolographicEnterprise":
            result = WINDOWS10HOLOGRAPHICENTERPRISE_WINDOWS10EDITIONTYPE
        case "windows10Professional":
            result = WINDOWS10PROFESSIONAL_WINDOWS10EDITIONTYPE
        case "windows10ProfessionalN":
            result = WINDOWS10PROFESSIONALN_WINDOWS10EDITIONTYPE
        case "windows10ProfessionalEducation":
            result = WINDOWS10PROFESSIONALEDUCATION_WINDOWS10EDITIONTYPE
        case "windows10ProfessionalEducationN":
            result = WINDOWS10PROFESSIONALEDUCATIONN_WINDOWS10EDITIONTYPE
        case "windows10ProfessionalWorkstation":
            result = WINDOWS10PROFESSIONALWORKSTATION_WINDOWS10EDITIONTYPE
        case "windows10ProfessionalWorkstationN":
            result = WINDOWS10PROFESSIONALWORKSTATIONN_WINDOWS10EDITIONTYPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeWindows10EditionType(values []Windows10EditionType) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i Windows10EditionType) isMultiValue() bool {
    return false
}
