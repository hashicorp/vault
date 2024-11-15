package models
type SynchronizationSecret int

const (
    NONE_SYNCHRONIZATIONSECRET SynchronizationSecret = iota
    USERNAME_SYNCHRONIZATIONSECRET
    PASSWORD_SYNCHRONIZATIONSECRET
    SECRETTOKEN_SYNCHRONIZATIONSECRET
    APPKEY_SYNCHRONIZATIONSECRET
    BASEADDRESS_SYNCHRONIZATIONSECRET
    CLIENTIDENTIFIER_SYNCHRONIZATIONSECRET
    CLIENTSECRET_SYNCHRONIZATIONSECRET
    SINGLESIGNONTYPE_SYNCHRONIZATIONSECRET
    SANDBOX_SYNCHRONIZATIONSECRET
    URL_SYNCHRONIZATIONSECRET
    DOMAIN_SYNCHRONIZATIONSECRET
    CONSUMERKEY_SYNCHRONIZATIONSECRET
    CONSUMERSECRET_SYNCHRONIZATIONSECRET
    TOKENKEY_SYNCHRONIZATIONSECRET
    TOKENEXPIRATION_SYNCHRONIZATIONSECRET
    OAUTH2ACCESSTOKEN_SYNCHRONIZATIONSECRET
    OAUTH2ACCESSTOKENCREATIONTIME_SYNCHRONIZATIONSECRET
    OAUTH2REFRESHTOKEN_SYNCHRONIZATIONSECRET
    SYNCALL_SYNCHRONIZATIONSECRET
    INSTANCENAME_SYNCHRONIZATIONSECRET
    OAUTH2CLIENTID_SYNCHRONIZATIONSECRET
    OAUTH2CLIENTSECRET_SYNCHRONIZATIONSECRET
    COMPANYID_SYNCHRONIZATIONSECRET
    UPDATEKEYONSOFTDELETE_SYNCHRONIZATIONSECRET
    SYNCHRONIZATIONSCHEDULE_SYNCHRONIZATIONSECRET
    SYSTEMOFRECORD_SYNCHRONIZATIONSECRET
    SANDBOXNAME_SYNCHRONIZATIONSECRET
    ENFORCEDOMAIN_SYNCHRONIZATIONSECRET
    SYNCNOTIFICATIONSETTINGS_SYNCHRONIZATIONSECRET
    SKIPOUTOFSCOPEDELETIONS_SYNCHRONIZATIONSECRET
    OAUTH2AUTHORIZATIONCODE_SYNCHRONIZATIONSECRET
    OAUTH2REDIRECTURI_SYNCHRONIZATIONSECRET
    APPLICATIONTEMPLATEIDENTIFIER_SYNCHRONIZATIONSECRET
    OAUTH2TOKENEXCHANGEURI_SYNCHRONIZATIONSECRET
    OAUTH2AUTHORIZATIONURI_SYNCHRONIZATIONSECRET
    AUTHENTICATIONTYPE_SYNCHRONIZATIONSECRET
    SERVER_SYNCHRONIZATIONSECRET
    PERFORMINBOUNDENTITLEMENTGRANTS_SYNCHRONIZATIONSECRET
    HARDDELETESENABLED_SYNCHRONIZATIONSECRET
    SYNCAGENTCOMPATIBILITYKEY_SYNCHRONIZATIONSECRET
    SYNCAGENTADCONTAINER_SYNCHRONIZATIONSECRET
    VALIDATEDOMAIN_SYNCHRONIZATIONSECRET
    TESTREFERENCES_SYNCHRONIZATIONSECRET
    CONNECTIONSTRING_SYNCHRONIZATIONSECRET
)

func (i SynchronizationSecret) String() string {
    return []string{"None", "UserName", "Password", "SecretToken", "AppKey", "BaseAddress", "ClientIdentifier", "ClientSecret", "SingleSignOnType", "Sandbox", "Url", "Domain", "ConsumerKey", "ConsumerSecret", "TokenKey", "TokenExpiration", "Oauth2AccessToken", "Oauth2AccessTokenCreationTime", "Oauth2RefreshToken", "SyncAll", "InstanceName", "Oauth2ClientId", "Oauth2ClientSecret", "CompanyId", "UpdateKeyOnSoftDelete", "SynchronizationSchedule", "SystemOfRecord", "SandboxName", "EnforceDomain", "SyncNotificationSettings", "SkipOutOfScopeDeletions", "Oauth2AuthorizationCode", "Oauth2RedirectUri", "ApplicationTemplateIdentifier", "Oauth2TokenExchangeUri", "Oauth2AuthorizationUri", "AuthenticationType", "Server", "PerformInboundEntitlementGrants", "HardDeletesEnabled", "SyncAgentCompatibilityKey", "SyncAgentADContainer", "ValidateDomain", "TestReferences", "ConnectionString"}[i]
}
func ParseSynchronizationSecret(v string) (any, error) {
    result := NONE_SYNCHRONIZATIONSECRET
    switch v {
        case "None":
            result = NONE_SYNCHRONIZATIONSECRET
        case "UserName":
            result = USERNAME_SYNCHRONIZATIONSECRET
        case "Password":
            result = PASSWORD_SYNCHRONIZATIONSECRET
        case "SecretToken":
            result = SECRETTOKEN_SYNCHRONIZATIONSECRET
        case "AppKey":
            result = APPKEY_SYNCHRONIZATIONSECRET
        case "BaseAddress":
            result = BASEADDRESS_SYNCHRONIZATIONSECRET
        case "ClientIdentifier":
            result = CLIENTIDENTIFIER_SYNCHRONIZATIONSECRET
        case "ClientSecret":
            result = CLIENTSECRET_SYNCHRONIZATIONSECRET
        case "SingleSignOnType":
            result = SINGLESIGNONTYPE_SYNCHRONIZATIONSECRET
        case "Sandbox":
            result = SANDBOX_SYNCHRONIZATIONSECRET
        case "Url":
            result = URL_SYNCHRONIZATIONSECRET
        case "Domain":
            result = DOMAIN_SYNCHRONIZATIONSECRET
        case "ConsumerKey":
            result = CONSUMERKEY_SYNCHRONIZATIONSECRET
        case "ConsumerSecret":
            result = CONSUMERSECRET_SYNCHRONIZATIONSECRET
        case "TokenKey":
            result = TOKENKEY_SYNCHRONIZATIONSECRET
        case "TokenExpiration":
            result = TOKENEXPIRATION_SYNCHRONIZATIONSECRET
        case "Oauth2AccessToken":
            result = OAUTH2ACCESSTOKEN_SYNCHRONIZATIONSECRET
        case "Oauth2AccessTokenCreationTime":
            result = OAUTH2ACCESSTOKENCREATIONTIME_SYNCHRONIZATIONSECRET
        case "Oauth2RefreshToken":
            result = OAUTH2REFRESHTOKEN_SYNCHRONIZATIONSECRET
        case "SyncAll":
            result = SYNCALL_SYNCHRONIZATIONSECRET
        case "InstanceName":
            result = INSTANCENAME_SYNCHRONIZATIONSECRET
        case "Oauth2ClientId":
            result = OAUTH2CLIENTID_SYNCHRONIZATIONSECRET
        case "Oauth2ClientSecret":
            result = OAUTH2CLIENTSECRET_SYNCHRONIZATIONSECRET
        case "CompanyId":
            result = COMPANYID_SYNCHRONIZATIONSECRET
        case "UpdateKeyOnSoftDelete":
            result = UPDATEKEYONSOFTDELETE_SYNCHRONIZATIONSECRET
        case "SynchronizationSchedule":
            result = SYNCHRONIZATIONSCHEDULE_SYNCHRONIZATIONSECRET
        case "SystemOfRecord":
            result = SYSTEMOFRECORD_SYNCHRONIZATIONSECRET
        case "SandboxName":
            result = SANDBOXNAME_SYNCHRONIZATIONSECRET
        case "EnforceDomain":
            result = ENFORCEDOMAIN_SYNCHRONIZATIONSECRET
        case "SyncNotificationSettings":
            result = SYNCNOTIFICATIONSETTINGS_SYNCHRONIZATIONSECRET
        case "SkipOutOfScopeDeletions":
            result = SKIPOUTOFSCOPEDELETIONS_SYNCHRONIZATIONSECRET
        case "Oauth2AuthorizationCode":
            result = OAUTH2AUTHORIZATIONCODE_SYNCHRONIZATIONSECRET
        case "Oauth2RedirectUri":
            result = OAUTH2REDIRECTURI_SYNCHRONIZATIONSECRET
        case "ApplicationTemplateIdentifier":
            result = APPLICATIONTEMPLATEIDENTIFIER_SYNCHRONIZATIONSECRET
        case "Oauth2TokenExchangeUri":
            result = OAUTH2TOKENEXCHANGEURI_SYNCHRONIZATIONSECRET
        case "Oauth2AuthorizationUri":
            result = OAUTH2AUTHORIZATIONURI_SYNCHRONIZATIONSECRET
        case "AuthenticationType":
            result = AUTHENTICATIONTYPE_SYNCHRONIZATIONSECRET
        case "Server":
            result = SERVER_SYNCHRONIZATIONSECRET
        case "PerformInboundEntitlementGrants":
            result = PERFORMINBOUNDENTITLEMENTGRANTS_SYNCHRONIZATIONSECRET
        case "HardDeletesEnabled":
            result = HARDDELETESENABLED_SYNCHRONIZATIONSECRET
        case "SyncAgentCompatibilityKey":
            result = SYNCAGENTCOMPATIBILITYKEY_SYNCHRONIZATIONSECRET
        case "SyncAgentADContainer":
            result = SYNCAGENTADCONTAINER_SYNCHRONIZATIONSECRET
        case "ValidateDomain":
            result = VALIDATEDOMAIN_SYNCHRONIZATIONSECRET
        case "TestReferences":
            result = TESTREFERENCES_SYNCHRONIZATIONSECRET
        case "ConnectionString":
            result = CONNECTIONSTRING_SYNCHRONIZATIONSECRET
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeSynchronizationSecret(values []SynchronizationSecret) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i SynchronizationSecret) isMultiValue() bool {
    return false
}
