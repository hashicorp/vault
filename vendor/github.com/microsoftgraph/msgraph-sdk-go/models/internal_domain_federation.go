package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type InternalDomainFederation struct {
    SamlOrWsFedProvider
}
// NewInternalDomainFederation instantiates a new InternalDomainFederation and sets the default values.
func NewInternalDomainFederation()(*InternalDomainFederation) {
    m := &InternalDomainFederation{
        SamlOrWsFedProvider: *NewSamlOrWsFedProvider(),
    }
    odataTypeValue := "#microsoft.graph.internalDomainFederation"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateInternalDomainFederationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateInternalDomainFederationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewInternalDomainFederation(), nil
}
// GetActiveSignInUri gets the activeSignInUri property value. URL of the endpoint used by active clients when authenticating with federated domains set up for single sign-on in Microsoft Entra ID. Corresponds to the ActiveLogOnUri property of the Set-MsolDomainFederationSettings MSOnline v1 PowerShell cmdlet.
// returns a *string when successful
func (m *InternalDomainFederation) GetActiveSignInUri()(*string) {
    val, err := m.GetBackingStore().Get("activeSignInUri")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFederatedIdpMfaBehavior gets the federatedIdpMfaBehavior property value. Determines whether Microsoft Entra ID accepts the MFA performed by the federated IdP when a federated user accesses an application that is governed by a conditional access policy that requires MFA. The possible values are: acceptIfMfaDoneByFederatedIdp, enforceMfaByFederatedIdp, rejectMfaByFederatedIdp, unknownFutureValue. For more information, see federatedIdpMfaBehavior values.
// returns a *FederatedIdpMfaBehavior when successful
func (m *InternalDomainFederation) GetFederatedIdpMfaBehavior()(*FederatedIdpMfaBehavior) {
    val, err := m.GetBackingStore().Get("federatedIdpMfaBehavior")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*FederatedIdpMfaBehavior)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *InternalDomainFederation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.SamlOrWsFedProvider.GetFieldDeserializers()
    res["activeSignInUri"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActiveSignInUri(val)
        }
        return nil
    }
    res["federatedIdpMfaBehavior"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseFederatedIdpMfaBehavior)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFederatedIdpMfaBehavior(val.(*FederatedIdpMfaBehavior))
        }
        return nil
    }
    res["isSignedAuthenticationRequestRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSignedAuthenticationRequestRequired(val)
        }
        return nil
    }
    res["nextSigningCertificate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNextSigningCertificate(val)
        }
        return nil
    }
    res["promptLoginBehavior"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePromptLoginBehavior)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPromptLoginBehavior(val.(*PromptLoginBehavior))
        }
        return nil
    }
    res["signingCertificateUpdateStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSigningCertificateUpdateStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSigningCertificateUpdateStatus(val.(SigningCertificateUpdateStatusable))
        }
        return nil
    }
    res["signOutUri"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSignOutUri(val)
        }
        return nil
    }
    return res
}
// GetIsSignedAuthenticationRequestRequired gets the isSignedAuthenticationRequestRequired property value. If true, when SAML authentication requests are sent to the federated SAML IdP, Microsoft Entra ID will sign those requests using the OrgID signing key. If false (default), the SAML authentication requests sent to the federated IdP aren't signed.
// returns a *bool when successful
func (m *InternalDomainFederation) GetIsSignedAuthenticationRequestRequired()(*bool) {
    val, err := m.GetBackingStore().Get("isSignedAuthenticationRequestRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNextSigningCertificate gets the nextSigningCertificate property value. Fallback token signing certificate that can also be used to sign tokens, for example when the primary signing certificate expires. Formatted as Base64 encoded strings of the public portion of the federated IdP's token signing certificate. Needs to be compatible with the X509Certificate2 class. Much like the signingCertificate, the nextSigningCertificate property is used if a rollover is required outside of the auto-rollover update, a new federation service is being set up, or if the new token signing certificate isn't present in the federation properties after the federation service certificate has been updated.
// returns a *string when successful
func (m *InternalDomainFederation) GetNextSigningCertificate()(*string) {
    val, err := m.GetBackingStore().Get("nextSigningCertificate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPromptLoginBehavior gets the promptLoginBehavior property value. Sets the preferred behavior for the sign-in prompt. The possible values are: translateToFreshPasswordAuthentication, nativeSupport, disabled, unknownFutureValue.
// returns a *PromptLoginBehavior when successful
func (m *InternalDomainFederation) GetPromptLoginBehavior()(*PromptLoginBehavior) {
    val, err := m.GetBackingStore().Get("promptLoginBehavior")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PromptLoginBehavior)
    }
    return nil
}
// GetSigningCertificateUpdateStatus gets the signingCertificateUpdateStatus property value. Provides status and timestamp of the last update of the signing certificate.
// returns a SigningCertificateUpdateStatusable when successful
func (m *InternalDomainFederation) GetSigningCertificateUpdateStatus()(SigningCertificateUpdateStatusable) {
    val, err := m.GetBackingStore().Get("signingCertificateUpdateStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SigningCertificateUpdateStatusable)
    }
    return nil
}
// GetSignOutUri gets the signOutUri property value. URI that clients are redirected to when they sign out of Microsoft Entra services. Corresponds to the LogOffUri property of the Set-MsolDomainFederationSettings MSOnline v1 PowerShell cmdlet.
// returns a *string when successful
func (m *InternalDomainFederation) GetSignOutUri()(*string) {
    val, err := m.GetBackingStore().Get("signOutUri")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *InternalDomainFederation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.SamlOrWsFedProvider.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("activeSignInUri", m.GetActiveSignInUri())
        if err != nil {
            return err
        }
    }
    if m.GetFederatedIdpMfaBehavior() != nil {
        cast := (*m.GetFederatedIdpMfaBehavior()).String()
        err = writer.WriteStringValue("federatedIdpMfaBehavior", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isSignedAuthenticationRequestRequired", m.GetIsSignedAuthenticationRequestRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("nextSigningCertificate", m.GetNextSigningCertificate())
        if err != nil {
            return err
        }
    }
    if m.GetPromptLoginBehavior() != nil {
        cast := (*m.GetPromptLoginBehavior()).String()
        err = writer.WriteStringValue("promptLoginBehavior", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("signingCertificateUpdateStatus", m.GetSigningCertificateUpdateStatus())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("signOutUri", m.GetSignOutUri())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActiveSignInUri sets the activeSignInUri property value. URL of the endpoint used by active clients when authenticating with federated domains set up for single sign-on in Microsoft Entra ID. Corresponds to the ActiveLogOnUri property of the Set-MsolDomainFederationSettings MSOnline v1 PowerShell cmdlet.
func (m *InternalDomainFederation) SetActiveSignInUri(value *string)() {
    err := m.GetBackingStore().Set("activeSignInUri", value)
    if err != nil {
        panic(err)
    }
}
// SetFederatedIdpMfaBehavior sets the federatedIdpMfaBehavior property value. Determines whether Microsoft Entra ID accepts the MFA performed by the federated IdP when a federated user accesses an application that is governed by a conditional access policy that requires MFA. The possible values are: acceptIfMfaDoneByFederatedIdp, enforceMfaByFederatedIdp, rejectMfaByFederatedIdp, unknownFutureValue. For more information, see federatedIdpMfaBehavior values.
func (m *InternalDomainFederation) SetFederatedIdpMfaBehavior(value *FederatedIdpMfaBehavior)() {
    err := m.GetBackingStore().Set("federatedIdpMfaBehavior", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSignedAuthenticationRequestRequired sets the isSignedAuthenticationRequestRequired property value. If true, when SAML authentication requests are sent to the federated SAML IdP, Microsoft Entra ID will sign those requests using the OrgID signing key. If false (default), the SAML authentication requests sent to the federated IdP aren't signed.
func (m *InternalDomainFederation) SetIsSignedAuthenticationRequestRequired(value *bool)() {
    err := m.GetBackingStore().Set("isSignedAuthenticationRequestRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetNextSigningCertificate sets the nextSigningCertificate property value. Fallback token signing certificate that can also be used to sign tokens, for example when the primary signing certificate expires. Formatted as Base64 encoded strings of the public portion of the federated IdP's token signing certificate. Needs to be compatible with the X509Certificate2 class. Much like the signingCertificate, the nextSigningCertificate property is used if a rollover is required outside of the auto-rollover update, a new federation service is being set up, or if the new token signing certificate isn't present in the federation properties after the federation service certificate has been updated.
func (m *InternalDomainFederation) SetNextSigningCertificate(value *string)() {
    err := m.GetBackingStore().Set("nextSigningCertificate", value)
    if err != nil {
        panic(err)
    }
}
// SetPromptLoginBehavior sets the promptLoginBehavior property value. Sets the preferred behavior for the sign-in prompt. The possible values are: translateToFreshPasswordAuthentication, nativeSupport, disabled, unknownFutureValue.
func (m *InternalDomainFederation) SetPromptLoginBehavior(value *PromptLoginBehavior)() {
    err := m.GetBackingStore().Set("promptLoginBehavior", value)
    if err != nil {
        panic(err)
    }
}
// SetSigningCertificateUpdateStatus sets the signingCertificateUpdateStatus property value. Provides status and timestamp of the last update of the signing certificate.
func (m *InternalDomainFederation) SetSigningCertificateUpdateStatus(value SigningCertificateUpdateStatusable)() {
    err := m.GetBackingStore().Set("signingCertificateUpdateStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetSignOutUri sets the signOutUri property value. URI that clients are redirected to when they sign out of Microsoft Entra services. Corresponds to the LogOffUri property of the Set-MsolDomainFederationSettings MSOnline v1 PowerShell cmdlet.
func (m *InternalDomainFederation) SetSignOutUri(value *string)() {
    err := m.GetBackingStore().Set("signOutUri", value)
    if err != nil {
        panic(err)
    }
}
type InternalDomainFederationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SamlOrWsFedProviderable
    GetActiveSignInUri()(*string)
    GetFederatedIdpMfaBehavior()(*FederatedIdpMfaBehavior)
    GetIsSignedAuthenticationRequestRequired()(*bool)
    GetNextSigningCertificate()(*string)
    GetPromptLoginBehavior()(*PromptLoginBehavior)
    GetSigningCertificateUpdateStatus()(SigningCertificateUpdateStatusable)
    GetSignOutUri()(*string)
    SetActiveSignInUri(value *string)()
    SetFederatedIdpMfaBehavior(value *FederatedIdpMfaBehavior)()
    SetIsSignedAuthenticationRequestRequired(value *bool)()
    SetNextSigningCertificate(value *string)()
    SetPromptLoginBehavior(value *PromptLoginBehavior)()
    SetSigningCertificateUpdateStatus(value SigningCertificateUpdateStatusable)()
    SetSignOutUri(value *string)()
}
