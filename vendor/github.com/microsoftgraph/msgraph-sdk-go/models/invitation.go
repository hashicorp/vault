package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Invitation struct {
    Entity
}
// NewInvitation instantiates a new Invitation and sets the default values.
func NewInvitation()(*Invitation) {
    m := &Invitation{
        Entity: *NewEntity(),
    }
    return m
}
// CreateInvitationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateInvitationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewInvitation(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Invitation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["invitedUser"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInvitedUser(val.(Userable))
        }
        return nil
    }
    res["invitedUserDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInvitedUserDisplayName(val)
        }
        return nil
    }
    res["invitedUserEmailAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInvitedUserEmailAddress(val)
        }
        return nil
    }
    res["invitedUserMessageInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateInvitedUserMessageInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInvitedUserMessageInfo(val.(InvitedUserMessageInfoable))
        }
        return nil
    }
    res["invitedUserSponsors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetInvitedUserSponsors(res)
        }
        return nil
    }
    res["invitedUserType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInvitedUserType(val)
        }
        return nil
    }
    res["inviteRedeemUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInviteRedeemUrl(val)
        }
        return nil
    }
    res["inviteRedirectUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInviteRedirectUrl(val)
        }
        return nil
    }
    res["resetRedemption"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResetRedemption(val)
        }
        return nil
    }
    res["sendInvitationMessage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSendInvitationMessage(val)
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val)
        }
        return nil
    }
    return res
}
// GetInvitedUser gets the invitedUser property value. The user created as part of the invitation creation. Read-only. The id property is required in the request body to reset a redemption status.
// returns a Userable when successful
func (m *Invitation) GetInvitedUser()(Userable) {
    val, err := m.GetBackingStore().Get("invitedUser")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Userable)
    }
    return nil
}
// GetInvitedUserDisplayName gets the invitedUserDisplayName property value. The display name of the user being invited.
// returns a *string when successful
func (m *Invitation) GetInvitedUserDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("invitedUserDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetInvitedUserEmailAddress gets the invitedUserEmailAddress property value. The email address of the user being invited. Required. The following special characters aren't permitted in the email address:Tilde (~)Exclamation point (!)Number sign (#)Dollar sign ($)Percent (%)Circumflex (^)Ampersand (&)Asterisk (*)Parentheses (( ))Plus sign (+)Equal sign (=)Brackets ([ ])Braces ({ })Backslash (/)Slash mark (/)Pipe (/|)Semicolon (;)Colon (:)Quotation marks (')Angle brackets (< >)Question mark (?)Comma (,)However, the following exceptions apply:A period (.) or a hyphen (-) is permitted anywhere in the user name, except at the beginning or end of the name.An underscore (_) is permitted anywhere in the user name, including at the beginning or end of the name.
// returns a *string when successful
func (m *Invitation) GetInvitedUserEmailAddress()(*string) {
    val, err := m.GetBackingStore().Get("invitedUserEmailAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetInvitedUserMessageInfo gets the invitedUserMessageInfo property value. Contains configuration for the message being sent to the invited user, including customizing message text, language, and cc recipient list.
// returns a InvitedUserMessageInfoable when successful
func (m *Invitation) GetInvitedUserMessageInfo()(InvitedUserMessageInfoable) {
    val, err := m.GetBackingStore().Get("invitedUserMessageInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(InvitedUserMessageInfoable)
    }
    return nil
}
// GetInvitedUserSponsors gets the invitedUserSponsors property value. The users or groups who are sponsors of the invited user. Sponsors are users and groups that are responsible for guest users' privileges in the tenant and for keeping the guest users' information and access up to date.
// returns a []DirectoryObjectable when successful
func (m *Invitation) GetInvitedUserSponsors()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("invitedUserSponsors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetInvitedUserType gets the invitedUserType property value. The userType of the user being invited. By default, this is Guest. You can invite as Member if you're a company administrator.
// returns a *string when successful
func (m *Invitation) GetInvitedUserType()(*string) {
    val, err := m.GetBackingStore().Get("invitedUserType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetInviteRedeemUrl gets the inviteRedeemUrl property value. The URL the user can use to redeem their invitation. Read-only.
// returns a *string when successful
func (m *Invitation) GetInviteRedeemUrl()(*string) {
    val, err := m.GetBackingStore().Get("inviteRedeemUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetInviteRedirectUrl gets the inviteRedirectUrl property value. The URL the user should be redirected to after the invitation is redeemed. Required.
// returns a *string when successful
func (m *Invitation) GetInviteRedirectUrl()(*string) {
    val, err := m.GetBackingStore().Get("inviteRedirectUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResetRedemption gets the resetRedemption property value. Reset the user's redemption status and reinvite a user while retaining their user identifier, group memberships, and app assignments. This property allows you to enable a user to sign-in using a different email address from the one in the previous invitation. When true, the invitedUser/id relationship is required. For more information about using this property, see Reset redemption status for a guest user.
// returns a *bool when successful
func (m *Invitation) GetResetRedemption()(*bool) {
    val, err := m.GetBackingStore().Get("resetRedemption")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSendInvitationMessage gets the sendInvitationMessage property value. Indicates whether an email should be sent to the user being invited. The default is false.
// returns a *bool when successful
func (m *Invitation) GetSendInvitationMessage()(*bool) {
    val, err := m.GetBackingStore().Get("sendInvitationMessage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetStatus gets the status property value. The status of the invitation. Possible values are: PendingAcceptance, Completed, InProgress, and Error.
// returns a *string when successful
func (m *Invitation) GetStatus()(*string) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Invitation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("invitedUser", m.GetInvitedUser())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("invitedUserDisplayName", m.GetInvitedUserDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("invitedUserEmailAddress", m.GetInvitedUserEmailAddress())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("invitedUserMessageInfo", m.GetInvitedUserMessageInfo())
        if err != nil {
            return err
        }
    }
    if m.GetInvitedUserSponsors() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInvitedUserSponsors()))
        for i, v := range m.GetInvitedUserSponsors() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("invitedUserSponsors", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("invitedUserType", m.GetInvitedUserType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("inviteRedeemUrl", m.GetInviteRedeemUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("inviteRedirectUrl", m.GetInviteRedirectUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("resetRedemption", m.GetResetRedemption())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("sendInvitationMessage", m.GetSendInvitationMessage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInvitedUser sets the invitedUser property value. The user created as part of the invitation creation. Read-only. The id property is required in the request body to reset a redemption status.
func (m *Invitation) SetInvitedUser(value Userable)() {
    err := m.GetBackingStore().Set("invitedUser", value)
    if err != nil {
        panic(err)
    }
}
// SetInvitedUserDisplayName sets the invitedUserDisplayName property value. The display name of the user being invited.
func (m *Invitation) SetInvitedUserDisplayName(value *string)() {
    err := m.GetBackingStore().Set("invitedUserDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetInvitedUserEmailAddress sets the invitedUserEmailAddress property value. The email address of the user being invited. Required. The following special characters aren't permitted in the email address:Tilde (~)Exclamation point (!)Number sign (#)Dollar sign ($)Percent (%)Circumflex (^)Ampersand (&)Asterisk (*)Parentheses (( ))Plus sign (+)Equal sign (=)Brackets ([ ])Braces ({ })Backslash (/)Slash mark (/)Pipe (/|)Semicolon (;)Colon (:)Quotation marks (')Angle brackets (< >)Question mark (?)Comma (,)However, the following exceptions apply:A period (.) or a hyphen (-) is permitted anywhere in the user name, except at the beginning or end of the name.An underscore (_) is permitted anywhere in the user name, including at the beginning or end of the name.
func (m *Invitation) SetInvitedUserEmailAddress(value *string)() {
    err := m.GetBackingStore().Set("invitedUserEmailAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetInvitedUserMessageInfo sets the invitedUserMessageInfo property value. Contains configuration for the message being sent to the invited user, including customizing message text, language, and cc recipient list.
func (m *Invitation) SetInvitedUserMessageInfo(value InvitedUserMessageInfoable)() {
    err := m.GetBackingStore().Set("invitedUserMessageInfo", value)
    if err != nil {
        panic(err)
    }
}
// SetInvitedUserSponsors sets the invitedUserSponsors property value. The users or groups who are sponsors of the invited user. Sponsors are users and groups that are responsible for guest users' privileges in the tenant and for keeping the guest users' information and access up to date.
func (m *Invitation) SetInvitedUserSponsors(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("invitedUserSponsors", value)
    if err != nil {
        panic(err)
    }
}
// SetInvitedUserType sets the invitedUserType property value. The userType of the user being invited. By default, this is Guest. You can invite as Member if you're a company administrator.
func (m *Invitation) SetInvitedUserType(value *string)() {
    err := m.GetBackingStore().Set("invitedUserType", value)
    if err != nil {
        panic(err)
    }
}
// SetInviteRedeemUrl sets the inviteRedeemUrl property value. The URL the user can use to redeem their invitation. Read-only.
func (m *Invitation) SetInviteRedeemUrl(value *string)() {
    err := m.GetBackingStore().Set("inviteRedeemUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetInviteRedirectUrl sets the inviteRedirectUrl property value. The URL the user should be redirected to after the invitation is redeemed. Required.
func (m *Invitation) SetInviteRedirectUrl(value *string)() {
    err := m.GetBackingStore().Set("inviteRedirectUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetResetRedemption sets the resetRedemption property value. Reset the user's redemption status and reinvite a user while retaining their user identifier, group memberships, and app assignments. This property allows you to enable a user to sign-in using a different email address from the one in the previous invitation. When true, the invitedUser/id relationship is required. For more information about using this property, see Reset redemption status for a guest user.
func (m *Invitation) SetResetRedemption(value *bool)() {
    err := m.GetBackingStore().Set("resetRedemption", value)
    if err != nil {
        panic(err)
    }
}
// SetSendInvitationMessage sets the sendInvitationMessage property value. Indicates whether an email should be sent to the user being invited. The default is false.
func (m *Invitation) SetSendInvitationMessage(value *bool)() {
    err := m.GetBackingStore().Set("sendInvitationMessage", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The status of the invitation. Possible values are: PendingAcceptance, Completed, InProgress, and Error.
func (m *Invitation) SetStatus(value *string)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type Invitationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetInvitedUser()(Userable)
    GetInvitedUserDisplayName()(*string)
    GetInvitedUserEmailAddress()(*string)
    GetInvitedUserMessageInfo()(InvitedUserMessageInfoable)
    GetInvitedUserSponsors()([]DirectoryObjectable)
    GetInvitedUserType()(*string)
    GetInviteRedeemUrl()(*string)
    GetInviteRedirectUrl()(*string)
    GetResetRedemption()(*bool)
    GetSendInvitationMessage()(*bool)
    GetStatus()(*string)
    SetInvitedUser(value Userable)()
    SetInvitedUserDisplayName(value *string)()
    SetInvitedUserEmailAddress(value *string)()
    SetInvitedUserMessageInfo(value InvitedUserMessageInfoable)()
    SetInvitedUserSponsors(value []DirectoryObjectable)()
    SetInvitedUserType(value *string)()
    SetInviteRedeemUrl(value *string)()
    SetInviteRedirectUrl(value *string)()
    SetResetRedemption(value *bool)()
    SetSendInvitationMessage(value *bool)()
    SetStatus(value *string)()
}
