package models
type OAuthAppScope int

const (
    UNKNOWN_OAUTHAPPSCOPE OAuthAppScope = iota
    READCALENDAR_OAUTHAPPSCOPE
    READCONTACT_OAUTHAPPSCOPE
    READMAIL_OAUTHAPPSCOPE
    READALLCHAT_OAUTHAPPSCOPE
    READALLFILE_OAUTHAPPSCOPE
    READANDWRITEMAIL_OAUTHAPPSCOPE
    SENDMAIL_OAUTHAPPSCOPE
    UNKNOWNFUTUREVALUE_OAUTHAPPSCOPE
)

func (i OAuthAppScope) String() string {
    return []string{"unknown", "readCalendar", "readContact", "readMail", "readAllChat", "readAllFile", "readAndWriteMail", "sendMail", "unknownFutureValue"}[i]
}
func ParseOAuthAppScope(v string) (any, error) {
    result := UNKNOWN_OAUTHAPPSCOPE
    switch v {
        case "unknown":
            result = UNKNOWN_OAUTHAPPSCOPE
        case "readCalendar":
            result = READCALENDAR_OAUTHAPPSCOPE
        case "readContact":
            result = READCONTACT_OAUTHAPPSCOPE
        case "readMail":
            result = READMAIL_OAUTHAPPSCOPE
        case "readAllChat":
            result = READALLCHAT_OAUTHAPPSCOPE
        case "readAllFile":
            result = READALLFILE_OAUTHAPPSCOPE
        case "readAndWriteMail":
            result = READANDWRITEMAIL_OAUTHAPPSCOPE
        case "sendMail":
            result = SENDMAIL_OAUTHAPPSCOPE
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_OAUTHAPPSCOPE
        default:
            return nil, nil
    }
    return &result, nil
}
func SerializeOAuthAppScope(values []OAuthAppScope) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
func (i OAuthAppScope) isMultiValue() bool {
    return false
}
