package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type InternetExplorerMode struct {
    Entity
}
// NewInternetExplorerMode instantiates a new InternetExplorerMode and sets the default values.
func NewInternetExplorerMode()(*InternetExplorerMode) {
    m := &InternetExplorerMode{
        Entity: *NewEntity(),
    }
    return m
}
// CreateInternetExplorerModeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateInternetExplorerModeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewInternetExplorerMode(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *InternetExplorerMode) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["siteLists"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBrowserSiteListFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BrowserSiteListable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BrowserSiteListable)
                }
            }
            m.SetSiteLists(res)
        }
        return nil
    }
    return res
}
// GetSiteLists gets the siteLists property value. A collection of site lists to support Internet Explorer mode.
// returns a []BrowserSiteListable when successful
func (m *InternetExplorerMode) GetSiteLists()([]BrowserSiteListable) {
    val, err := m.GetBackingStore().Get("siteLists")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BrowserSiteListable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *InternetExplorerMode) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetSiteLists() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSiteLists()))
        for i, v := range m.GetSiteLists() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("siteLists", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetSiteLists sets the siteLists property value. A collection of site lists to support Internet Explorer mode.
func (m *InternetExplorerMode) SetSiteLists(value []BrowserSiteListable)() {
    err := m.GetBackingStore().Set("siteLists", value)
    if err != nil {
        panic(err)
    }
}
type InternetExplorerModeable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetSiteLists()([]BrowserSiteListable)
    SetSiteLists(value []BrowserSiteListable)()
}
