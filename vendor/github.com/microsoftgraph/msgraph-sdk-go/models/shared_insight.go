package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type SharedInsight struct {
    Entity
}
// NewSharedInsight instantiates a new SharedInsight and sets the default values.
func NewSharedInsight()(*SharedInsight) {
    m := &SharedInsight{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSharedInsightFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSharedInsightFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSharedInsight(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SharedInsight) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["lastShared"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSharingDetailFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastShared(val.(SharingDetailable))
        }
        return nil
    }
    res["lastSharedMethod"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEntityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastSharedMethod(val.(Entityable))
        }
        return nil
    }
    res["resource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEntityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResource(val.(Entityable))
        }
        return nil
    }
    res["resourceReference"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateResourceReferenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceReference(val.(ResourceReferenceable))
        }
        return nil
    }
    res["resourceVisualization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateResourceVisualizationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceVisualization(val.(ResourceVisualizationable))
        }
        return nil
    }
    res["sharingHistory"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSharingDetailFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SharingDetailable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SharingDetailable)
                }
            }
            m.SetSharingHistory(res)
        }
        return nil
    }
    return res
}
// GetLastShared gets the lastShared property value. Details about the shared item. Read only.
// returns a SharingDetailable when successful
func (m *SharedInsight) GetLastShared()(SharingDetailable) {
    val, err := m.GetBackingStore().Get("lastShared")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SharingDetailable)
    }
    return nil
}
// GetLastSharedMethod gets the lastSharedMethod property value. The lastSharedMethod property
// returns a Entityable when successful
func (m *SharedInsight) GetLastSharedMethod()(Entityable) {
    val, err := m.GetBackingStore().Get("lastSharedMethod")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Entityable)
    }
    return nil
}
// GetResource gets the resource property value. Used for navigating to the item that was shared. For file attachments, the type is fileAttachment. For linked attachments, the type is driveItem.
// returns a Entityable when successful
func (m *SharedInsight) GetResource()(Entityable) {
    val, err := m.GetBackingStore().Get("resource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Entityable)
    }
    return nil
}
// GetResourceReference gets the resourceReference property value. Reference properties of the shared document, such as the url and type of the document. Read-only
// returns a ResourceReferenceable when successful
func (m *SharedInsight) GetResourceReference()(ResourceReferenceable) {
    val, err := m.GetBackingStore().Get("resourceReference")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ResourceReferenceable)
    }
    return nil
}
// GetResourceVisualization gets the resourceVisualization property value. Properties that you can use to visualize the document in your experience. Read-only
// returns a ResourceVisualizationable when successful
func (m *SharedInsight) GetResourceVisualization()(ResourceVisualizationable) {
    val, err := m.GetBackingStore().Get("resourceVisualization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ResourceVisualizationable)
    }
    return nil
}
// GetSharingHistory gets the sharingHistory property value. The sharingHistory property
// returns a []SharingDetailable when successful
func (m *SharedInsight) GetSharingHistory()([]SharingDetailable) {
    val, err := m.GetBackingStore().Get("sharingHistory")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SharingDetailable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SharedInsight) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("lastShared", m.GetLastShared())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastSharedMethod", m.GetLastSharedMethod())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("resource", m.GetResource())
        if err != nil {
            return err
        }
    }
    if m.GetSharingHistory() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSharingHistory()))
        for i, v := range m.GetSharingHistory() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("sharingHistory", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLastShared sets the lastShared property value. Details about the shared item. Read only.
func (m *SharedInsight) SetLastShared(value SharingDetailable)() {
    err := m.GetBackingStore().Set("lastShared", value)
    if err != nil {
        panic(err)
    }
}
// SetLastSharedMethod sets the lastSharedMethod property value. The lastSharedMethod property
func (m *SharedInsight) SetLastSharedMethod(value Entityable)() {
    err := m.GetBackingStore().Set("lastSharedMethod", value)
    if err != nil {
        panic(err)
    }
}
// SetResource sets the resource property value. Used for navigating to the item that was shared. For file attachments, the type is fileAttachment. For linked attachments, the type is driveItem.
func (m *SharedInsight) SetResource(value Entityable)() {
    err := m.GetBackingStore().Set("resource", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceReference sets the resourceReference property value. Reference properties of the shared document, such as the url and type of the document. Read-only
func (m *SharedInsight) SetResourceReference(value ResourceReferenceable)() {
    err := m.GetBackingStore().Set("resourceReference", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceVisualization sets the resourceVisualization property value. Properties that you can use to visualize the document in your experience. Read-only
func (m *SharedInsight) SetResourceVisualization(value ResourceVisualizationable)() {
    err := m.GetBackingStore().Set("resourceVisualization", value)
    if err != nil {
        panic(err)
    }
}
// SetSharingHistory sets the sharingHistory property value. The sharingHistory property
func (m *SharedInsight) SetSharingHistory(value []SharingDetailable)() {
    err := m.GetBackingStore().Set("sharingHistory", value)
    if err != nil {
        panic(err)
    }
}
type SharedInsightable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetLastShared()(SharingDetailable)
    GetLastSharedMethod()(Entityable)
    GetResource()(Entityable)
    GetResourceReference()(ResourceReferenceable)
    GetResourceVisualization()(ResourceVisualizationable)
    GetSharingHistory()([]SharingDetailable)
    SetLastShared(value SharingDetailable)()
    SetLastSharedMethod(value Entityable)()
    SetResource(value Entityable)()
    SetResourceReference(value ResourceReferenceable)()
    SetResourceVisualization(value ResourceVisualizationable)()
    SetSharingHistory(value []SharingDetailable)()
}
