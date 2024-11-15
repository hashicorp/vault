package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Fido2AuthenticationMethod struct {
    AuthenticationMethod
}
// NewFido2AuthenticationMethod instantiates a new Fido2AuthenticationMethod and sets the default values.
func NewFido2AuthenticationMethod()(*Fido2AuthenticationMethod) {
    m := &Fido2AuthenticationMethod{
        AuthenticationMethod: *NewAuthenticationMethod(),
    }
    odataTypeValue := "#microsoft.graph.fido2AuthenticationMethod"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateFido2AuthenticationMethodFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFido2AuthenticationMethodFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFido2AuthenticationMethod(), nil
}
// GetAaGuid gets the aaGuid property value. Authenticator Attestation GUID, an identifier that indicates the type (e.g. make and model) of the authenticator.
// returns a *string when successful
func (m *Fido2AuthenticationMethod) GetAaGuid()(*string) {
    val, err := m.GetBackingStore().Get("aaGuid")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAttestationCertificates gets the attestationCertificates property value. The attestation certificate(s) attached to this security key.
// returns a []string when successful
func (m *Fido2AuthenticationMethod) GetAttestationCertificates()([]string) {
    val, err := m.GetBackingStore().Get("attestationCertificates")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetAttestationLevel gets the attestationLevel property value. The attestation level of this FIDO2 security key. Possible values are: attested, or notAttested.
// returns a *AttestationLevel when successful
func (m *Fido2AuthenticationMethod) GetAttestationLevel()(*AttestationLevel) {
    val, err := m.GetBackingStore().Get("attestationLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*AttestationLevel)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. The timestamp when this key was registered to the user.
// returns a *Time when successful
func (m *Fido2AuthenticationMethod) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name of the key as given by the user.
// returns a *string when successful
func (m *Fido2AuthenticationMethod) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Fido2AuthenticationMethod) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AuthenticationMethod.GetFieldDeserializers()
    res["aaGuid"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAaGuid(val)
        }
        return nil
    }
    res["attestationCertificates"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetAttestationCertificates(res)
        }
        return nil
    }
    res["attestationLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseAttestationLevel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttestationLevel(val.(*AttestationLevel))
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["model"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModel(val)
        }
        return nil
    }
    return res
}
// GetModel gets the model property value. The manufacturer-assigned model of the FIDO2 security key.
// returns a *string when successful
func (m *Fido2AuthenticationMethod) GetModel()(*string) {
    val, err := m.GetBackingStore().Get("model")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Fido2AuthenticationMethod) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AuthenticationMethod.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("aaGuid", m.GetAaGuid())
        if err != nil {
            return err
        }
    }
    if m.GetAttestationCertificates() != nil {
        err = writer.WriteCollectionOfStringValues("attestationCertificates", m.GetAttestationCertificates())
        if err != nil {
            return err
        }
    }
    if m.GetAttestationLevel() != nil {
        cast := (*m.GetAttestationLevel()).String()
        err = writer.WriteStringValue("attestationLevel", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("model", m.GetModel())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAaGuid sets the aaGuid property value. Authenticator Attestation GUID, an identifier that indicates the type (e.g. make and model) of the authenticator.
func (m *Fido2AuthenticationMethod) SetAaGuid(value *string)() {
    err := m.GetBackingStore().Set("aaGuid", value)
    if err != nil {
        panic(err)
    }
}
// SetAttestationCertificates sets the attestationCertificates property value. The attestation certificate(s) attached to this security key.
func (m *Fido2AuthenticationMethod) SetAttestationCertificates(value []string)() {
    err := m.GetBackingStore().Set("attestationCertificates", value)
    if err != nil {
        panic(err)
    }
}
// SetAttestationLevel sets the attestationLevel property value. The attestation level of this FIDO2 security key. Possible values are: attested, or notAttested.
func (m *Fido2AuthenticationMethod) SetAttestationLevel(value *AttestationLevel)() {
    err := m.GetBackingStore().Set("attestationLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. The timestamp when this key was registered to the user.
func (m *Fido2AuthenticationMethod) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name of the key as given by the user.
func (m *Fido2AuthenticationMethod) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetModel sets the model property value. The manufacturer-assigned model of the FIDO2 security key.
func (m *Fido2AuthenticationMethod) SetModel(value *string)() {
    err := m.GetBackingStore().Set("model", value)
    if err != nil {
        panic(err)
    }
}
type Fido2AuthenticationMethodable interface {
    AuthenticationMethodable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAaGuid()(*string)
    GetAttestationCertificates()([]string)
    GetAttestationLevel()(*AttestationLevel)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDisplayName()(*string)
    GetModel()(*string)
    SetAaGuid(value *string)()
    SetAttestationCertificates(value []string)()
    SetAttestationLevel(value *AttestationLevel)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDisplayName(value *string)()
    SetModel(value *string)()
}
