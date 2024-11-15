package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosHomeScreenFolder a folder containing pages of apps and web clips on the Home Screen.
type IosHomeScreenFolder struct {
    IosHomeScreenItem
}
// NewIosHomeScreenFolder instantiates a new IosHomeScreenFolder and sets the default values.
func NewIosHomeScreenFolder()(*IosHomeScreenFolder) {
    m := &IosHomeScreenFolder{
        IosHomeScreenItem: *NewIosHomeScreenItem(),
    }
    odataTypeValue := "#microsoft.graph.iosHomeScreenFolder"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIosHomeScreenFolderFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosHomeScreenFolderFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosHomeScreenFolder(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IosHomeScreenFolder) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.IosHomeScreenItem.GetFieldDeserializers()
    res["pages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateIosHomeScreenFolderPageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]IosHomeScreenFolderPageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(IosHomeScreenFolderPageable)
                }
            }
            m.SetPages(res)
        }
        return nil
    }
    return res
}
// GetPages gets the pages property value. Pages of Home Screen Layout Icons which must be applications or web clips. This collection can contain a maximum of 500 elements.
// returns a []IosHomeScreenFolderPageable when successful
func (m *IosHomeScreenFolder) GetPages()([]IosHomeScreenFolderPageable) {
    val, err := m.GetBackingStore().Get("pages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]IosHomeScreenFolderPageable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IosHomeScreenFolder) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.IosHomeScreenItem.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetPages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPages()))
        for i, v := range m.GetPages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("pages", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetPages sets the pages property value. Pages of Home Screen Layout Icons which must be applications or web clips. This collection can contain a maximum of 500 elements.
func (m *IosHomeScreenFolder) SetPages(value []IosHomeScreenFolderPageable)() {
    err := m.GetBackingStore().Set("pages", value)
    if err != nil {
        panic(err)
    }
}
type IosHomeScreenFolderable interface {
    IosHomeScreenItemable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetPages()([]IosHomeScreenFolderPageable)
    SetPages(value []IosHomeScreenFolderPageable)()
}
