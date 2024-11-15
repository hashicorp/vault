package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ContainerImageEvidence struct {
    AlertEvidence
}
// NewContainerImageEvidence instantiates a new ContainerImageEvidence and sets the default values.
func NewContainerImageEvidence()(*ContainerImageEvidence) {
    m := &ContainerImageEvidence{
        AlertEvidence: *NewAlertEvidence(),
    }
    odataTypeValue := "#microsoft.graph.security.containerImageEvidence"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateContainerImageEvidenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateContainerImageEvidenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewContainerImageEvidence(), nil
}
// GetDigestImage gets the digestImage property value. The digest image entity, in case this is a tag image.
// returns a ContainerImageEvidenceable when successful
func (m *ContainerImageEvidence) GetDigestImage()(ContainerImageEvidenceable) {
    val, err := m.GetBackingStore().Get("digestImage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ContainerImageEvidenceable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ContainerImageEvidence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AlertEvidence.GetFieldDeserializers()
    res["digestImage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateContainerImageEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDigestImage(val.(ContainerImageEvidenceable))
        }
        return nil
    }
    res["imageId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImageId(val)
        }
        return nil
    }
    res["registry"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateContainerRegistryEvidenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRegistry(val.(ContainerRegistryEvidenceable))
        }
        return nil
    }
    return res
}
// GetImageId gets the imageId property value. The unique identifier for the container image entity.
// returns a *string when successful
func (m *ContainerImageEvidence) GetImageId()(*string) {
    val, err := m.GetBackingStore().Get("imageId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRegistry gets the registry property value. The container registry for this image.
// returns a ContainerRegistryEvidenceable when successful
func (m *ContainerImageEvidence) GetRegistry()(ContainerRegistryEvidenceable) {
    val, err := m.GetBackingStore().Get("registry")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ContainerRegistryEvidenceable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ContainerImageEvidence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AlertEvidence.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("digestImage", m.GetDigestImage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("imageId", m.GetImageId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("registry", m.GetRegistry())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDigestImage sets the digestImage property value. The digest image entity, in case this is a tag image.
func (m *ContainerImageEvidence) SetDigestImage(value ContainerImageEvidenceable)() {
    err := m.GetBackingStore().Set("digestImage", value)
    if err != nil {
        panic(err)
    }
}
// SetImageId sets the imageId property value. The unique identifier for the container image entity.
func (m *ContainerImageEvidence) SetImageId(value *string)() {
    err := m.GetBackingStore().Set("imageId", value)
    if err != nil {
        panic(err)
    }
}
// SetRegistry sets the registry property value. The container registry for this image.
func (m *ContainerImageEvidence) SetRegistry(value ContainerRegistryEvidenceable)() {
    err := m.GetBackingStore().Set("registry", value)
    if err != nil {
        panic(err)
    }
}
type ContainerImageEvidenceable interface {
    AlertEvidenceable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDigestImage()(ContainerImageEvidenceable)
    GetImageId()(*string)
    GetRegistry()(ContainerRegistryEvidenceable)
    SetDigestImage(value ContainerImageEvidenceable)()
    SetImageId(value *string)()
    SetRegistry(value ContainerRegistryEvidenceable)()
}
