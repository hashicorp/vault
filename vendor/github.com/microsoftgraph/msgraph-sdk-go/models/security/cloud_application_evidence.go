package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CloudApplicationEvidence struct {
    AlertEvidence
}
// NewCloudApplicationEvidence instantiates a new CloudApplicationEvidence and sets the default values.
func NewCloudApplicationEvidence()(*CloudApplicationEvidence) {
    m := &CloudApplicationEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.cloudApplicationEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCloudApplicationEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCloudApplicationEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCloudApplicationEvidence(), nil
}
// GetAppId gets the appId property value. Unique identifier of the application.
// returns a *int64 when successful
func (m *CloudApplicationEvidence) GetAppId()(*int64) {
    val, err := m.GetBackingStore().Get("appId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Name of the application.
// returns a *string when successful
func (m *CloudApplicationEvidence) GetDisplayName()(*string) {
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
func (m *CloudApplicationEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["appId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppId(val)
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
    res["instanceId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstanceId(val)
        }
        return nil
    }
    res["instanceName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstanceName(val)
        }
        return nil
    }
    res["saasAppId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSaasAppId(val)
        }
        return nil
    }
    res["stream"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateStreamFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStream(val.(Streamable))
        }
        return nil
    }
    return res
}
// GetInstanceId gets the instanceId property value. Identifier of the instance of the Software as a Service (SaaS) application.
// returns a *int64 when successful
func (m *CloudApplicationEvidence) GetInstanceId()(*int64) {
    val, err := m.GetBackingStore().Get("instanceId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetInstanceName gets the instanceName property value. Name of the instance of the SaaS application.
// returns a *string when successful
func (m *CloudApplicationEvidence) GetInstanceName()(*string) {
    val, err := m.GetBackingStore().Get("instanceName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSaasAppId gets the saasAppId property value. The identifier of the SaaS application.
// returns a *int64 when successful
func (m *CloudApplicationEvidence) GetSaasAppId()(*int64) {
    val, err := m.GetBackingStore().Get("saasAppId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetStream gets the stream property value. The stream property
// returns a Streamable when successful
func (m *CloudApplicationEvidence) GetStream()(Streamable) {
    val, err := m.GetBackingStore().Get("stream")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Streamable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CloudApplicationEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteInt64Value("appId", m.GetAppId())
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
        err = writer.WriteInt64Value("instanceId", m.GetInstanceId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("instanceName", m.GetInstanceName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("saasAppId", m.GetSaasAppId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("stream", m.GetStream())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppId sets the appId property value. Unique identifier of the application.
func (m *CloudApplicationEvidence) SetAppId(value *int64)() {
    err := m.GetBackingStore().Set("appId", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Name of the application.
func (m *CloudApplicationEvidence) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetInstanceId sets the instanceId property value. Identifier of the instance of the Software as a Service (SaaS) application.
func (m *CloudApplicationEvidence) SetInstanceId(value *int64)() {
    err := m.GetBackingStore().Set("instanceId", value)
    if err != nil {
        panic(err)
    }
}
// SetInstanceName sets the instanceName property value. Name of the instance of the SaaS application.
func (m *CloudApplicationEvidence) SetInstanceName(value *string)() {
    err := m.GetBackingStore().Set("instanceName", value)
    if err != nil {
        panic(err)
    }
}
// SetSaasAppId sets the saasAppId property value. The identifier of the SaaS application.
func (m *CloudApplicationEvidence) SetSaasAppId(value *int64)() {
    err := m.GetBackingStore().Set("saasAppId", value)
    if err != nil {
        panic(err)
    }
}
// SetStream sets the stream property value. The stream property
func (m *CloudApplicationEvidence) SetStream(value Streamable)() {
    err := m.GetBackingStore().Set("stream", value)
    if err != nil {
        panic(err)
    }
}
type CloudApplicationEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppId()(*int64)
    GetDisplayName()(*string)
    GetInstanceId()(*int64)
    GetInstanceName()(*string)
    GetSaasAppId()(*int64)
    GetStream()(Streamable)
    SetAppId(value *int64)()
    SetDisplayName(value *string)()
    SetInstanceId(value *int64)()
    SetInstanceName(value *string)()
    SetSaasAppId(value *int64)()
    SetStream(value Streamable)()
}
