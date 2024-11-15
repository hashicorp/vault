package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type FilePlanAppliedCategory struct {
    FilePlanDescriptorBase
}
// NewFilePlanAppliedCategory instantiates a new FilePlanAppliedCategory and sets the default values.
func NewFilePlanAppliedCategory()(*FilePlanAppliedCategory) {
    m := &FilePlanAppliedCategory{
        FilePlanDescriptorBase: *NewFilePlanDescriptorBase(),
    }
    return m
}
// CreateFilePlanAppliedCategoryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFilePlanAppliedCategoryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFilePlanAppliedCategory(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *FilePlanAppliedCategory) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.FilePlanDescriptorBase.GetFieldDeserializers()
    res["subcategory"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFilePlanSubcategoryFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubcategory(val.(FilePlanSubcategoryable))
        }
        return nil
    }
    return res
}
// GetSubcategory gets the subcategory property value. Represents the file plan descriptor for a subcategory under a specific category, which has been assigned to a particular retention label.
// returns a FilePlanSubcategoryable when successful
func (m *FilePlanAppliedCategory) GetSubcategory()(FilePlanSubcategoryable) {
    val, err := m.GetBackingStore().Get("subcategory")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FilePlanSubcategoryable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *FilePlanAppliedCategory) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.FilePlanDescriptorBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("subcategory", m.GetSubcategory())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSubcategory sets the subcategory property value. Represents the file plan descriptor for a subcategory under a specific category, which has been assigned to a particular retention label.
func (m *FilePlanAppliedCategory) SetSubcategory(value FilePlanSubcategoryable)() {
    err := m.GetBackingStore().Set("subcategory", value)
    if err != nil {
        panic(err)
    }
}
type FilePlanAppliedCategoryable interface {
    FilePlanDescriptorBaseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetSubcategory()(FilePlanSubcategoryable)
    SetSubcategory(value FilePlanSubcategoryable)()
}
