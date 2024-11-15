package callrecords
type ServiceRole int

const (
    UNKNOWN_SERVICEROLE ServiceRole = iota
    CUSTOMBOT_SERVICEROLE
    SKYPEFORBUSINESSMICROSOFTTEAMSGATEWAY_SERVICEROLE
    SKYPEFORBUSINESSAUDIOVIDEOMCU_SERVICEROLE
    SKYPEFORBUSINESSAPPLICATIONSHARINGMCU_SERVICEROLE
    SKYPEFORBUSINESSCALLQUEUES_SERVICEROLE
    SKYPEFORBUSINESSAUTOATTENDANT_SERVICEROLE
    MEDIATIONSERVER_SERVICEROLE
    MEDIATIONSERVERCLOUDCONNECTOREDITION_SERVICEROLE
    EXCHANGEUNIFIEDMESSAGINGSERVICE_SERVICEROLE
    MEDIACONTROLLER_SERVICEROLE
    CONFERENCINGANNOUNCEMENTSERVICE_SERVICEROLE
    CONFERENCINGATTENDANT_SERVICEROLE
    AUDIOTELECONFERENCERCONTROLLER_SERVICEROLE
    SKYPEFORBUSINESSUNIFIEDCOMMUNICATIONAPPLICATIONPLATFORM_SERVICEROLE
    RESPONSEGROUPSERVICEANNOUNCEMENTSERVICE_SERVICEROLE
    GATEWAY_SERVICEROLE
    SKYPETRANSLATOR_SERVICEROLE
    SKYPEFORBUSINESSATTENDANT_SERVICEROLE
    RESPONSEGROUPSERVICE_SERVICEROLE
    VOICEMAIL_SERVICEROLE
    UNKNOWNFUTUREVALUE_SERVICEROLE
)

func (i ServiceRole) String() string {
    return []string{"unknown", "customBot", "skypeForBusinessMicrosoftTeamsGateway", "skypeForBusinessAudioVideoMcu", "skypeForBusinessApplicationSharingMcu", "skypeForBusinessCallQueues", "skypeForBusinessAutoAttendant", "mediationServer", "mediationServerCloudConnectorEdition", "exchangeUnifiedMessagingService", "mediaController", "conferencingAnnouncementService", "conferencingAttendant", "audioTeleconferencerController", "skypeForBusinessUnifiedCommunicationApplicationPlatform", "responseGroupServiceAnnouncementService", "gateway", "skypeTranslator", "skypeForBusinessAttendant", "responseGroupService", "voicemail", "unknownFutureValue"}[i]
}
func ParseServiceRole(v string) (any, error) {
    result := UNKNOWN_SERVICEROLE
    switch v {
        case "unknown":
            result = UNKNOWN_SERVICEROLE
        case "customBot":
            result = CUSTOMBOT_SERVICEROLE
        case "skypeForBusinessMicrosoftTeamsGateway":
            result = SKYPEFORBUSINESSMICROSOFTTEAMSGATEWAY_SERVICEROLE
        case "skypeForBusinessAudioVideoMcu":
            result = SKYPEFORBUSINESSAUDIOVIDEOMCU_SERVICEROLE
        case "skypeForBusinessApplicationSharingMcu":
            result = SKYPEFORBUSINESSAPPLICATIONSHARINGMCU_SERVICEROLE
        case "skypeForBusinessCallQueues":
            result = SKYPEFORBUSINESSCALLQUEUES_SERVICEROLE
        case "skypeForBusinessAutoAttendant":
            result = SKYPEFORBUSINESSAUTOATTENDANT_SERVICEROLE
        case "mediationServer":
            result = MEDIATIONSERVER_SERVICEROLE
        case "mediationServerCloudConnectorEdition":
            result = MEDIATIONSERVERCLOUDCONNECTOREDITION_SERVICEROLE
        case "exchangeUnifiedMessagingService":
            result = EXCHANGEUNIFIEDMESSAGINGSERVICE_SERVICEROLE
        case "mediaController":
            result = MEDIACONTROLLER_SERVICEROLE
        case "conferencingAnnouncementService":
            result = CONFERENCINGANNOUNCEMENTSERVICE_SERVICEROLE
        case "conferencingAttendant":
            result = CONFERENCINGATTENDANT_SERVICEROLE
        case "audioTeleconferencerController":
            result = AUDIOTELECONFERENCERCONTROLLER_SERVICEROLE
        case "skypeForBusinessUnifiedCommunicationApplicationPlatform":
            result = SKYPEFORBUSINESSUNIFIEDCOMMUNICATIONAPPLICATIONPLATFORM_SERVICEROLE
        case "responseGroupServiceAnnouncementService":
            result = RESPONSEGROUPSERVICEANNOUNCEMENTSERVICE_SERVICEROLE
        case "gateway":
            result = GATEWAY_SERVICEROLE
        case "skypeTranslator":
            result = SKYPETRANSLATOR_SERVICEROLE
        case "skypeForBusinessAttendant":
            result = SKYPEFORBUSINESSATTENDANT_SERVICEROLE
        case "responseGroupService":
            result = RESPONSEGROUPSERVICE_SERVICEROLE
        case "voicemail":
            result = VOICEMAIL_SERVICEROLE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_SERVICEROLE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeServiceRole(values []ServiceRole) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i ServiceRole) isMultiValue() bool {
    return false
}
