package models

import (
    i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosVppEBook a class containing the properties for iOS Vpp eBook.
type IosVppEBook struct {
    ManagedEBook
}
// NewIosVppEBook instantiates a new IosVppEBook and sets the default values.
func NewIosVppEBook()(*IosVppEBook) {
    m := &IosVppEBook{
        ManagedEBook: *NewManagedEBook(),
    }
    odataTypeValue := "#microsoft.graph.iosVppEBook"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIosVppEBookFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosVppEBookFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosVppEBook(), nil
}
// GetAppleId gets the appleId property value. The Apple ID associated with Vpp token.
// returns a *string when successful
func (m *IosVppEBook) GetAppleId()(*string) {
    val, err := m.GetBackingStore().Get("appleId")
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
func (m *IosVppEBook) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedEBook.GetFieldDeserializers()
    res["appleId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppleId(val)
        }
        return nil
    }
    res["genres"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetGenres(res)
        }
        return nil
    }
    res["language"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLanguage(val)
        }
        return nil
    }
    res["seller"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSeller(val)
        }
        return nil
    }
    res["totalLicenseCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalLicenseCount(val)
        }
        return nil
    }
    res["usedLicenseCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUsedLicenseCount(val)
        }
        return nil
    }
    res["vppOrganizationName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVppOrganizationName(val)
        }
        return nil
    }
    res["vppTokenId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetUUIDValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVppTokenId(val)
        }
        return nil
    }
    return res
}
// GetGenres gets the genres property value. Genres.
// returns a []string when successful
func (m *IosVppEBook) GetGenres()([]string) {
    val, err := m.GetBackingStore().Get("genres")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetLanguage gets the language property value. Language.
// returns a *string when successful
func (m *IosVppEBook) GetLanguage()(*string) {
    val, err := m.GetBackingStore().Get("language")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSeller gets the seller property value. Seller.
// returns a *string when successful
func (m *IosVppEBook) GetSeller()(*string) {
    val, err := m.GetBackingStore().Get("seller")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTotalLicenseCount gets the totalLicenseCount property value. Total license count.
// returns a *int32 when successful
func (m *IosVppEBook) GetTotalLicenseCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalLicenseCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUsedLicenseCount gets the usedLicenseCount property value. Used license count.
// returns a *int32 when successful
func (m *IosVppEBook) GetUsedLicenseCount()(*int32) {
    val, err := m.GetBackingStore().Get("usedLicenseCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetVppOrganizationName gets the vppOrganizationName property value. The Vpp token's organization name.
// returns a *string when successful
func (m *IosVppEBook) GetVppOrganizationName()(*string) {
    val, err := m.GetBackingStore().Get("vppOrganizationName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVppTokenId gets the vppTokenId property value. The Vpp token ID.
// returns a *UUID when successful
func (m *IosVppEBook) GetVppTokenId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
    val, err := m.GetBackingStore().Get("vppTokenId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IosVppEBook) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedEBook.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appleId", m.GetAppleId())
        if err != nil {
            return err
        }
    }
    if m.GetGenres() != nil {
        err = writer.WriteCollectionOfStringValues("genres", m.GetGenres())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("language", m.GetLanguage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("seller", m.GetSeller())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("totalLicenseCount", m.GetTotalLicenseCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("usedLicenseCount", m.GetUsedLicenseCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("vppOrganizationName", m.GetVppOrganizationName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteUUIDValue("vppTokenId", m.GetVppTokenId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppleId sets the appleId property value. The Apple ID associated with Vpp token.
func (m *IosVppEBook) SetAppleId(value *string)() {
    err := m.GetBackingStore().Set("appleId", value)
    if err != nil {
        panic(err)
    }
}
// SetGenres sets the genres property value. Genres.
func (m *IosVppEBook) SetGenres(value []string)() {
    err := m.GetBackingStore().Set("genres", value)
    if err != nil {
        panic(err)
    }
}
// SetLanguage sets the language property value. Language.
func (m *IosVppEBook) SetLanguage(value *string)() {
    err := m.GetBackingStore().Set("language", value)
    if err != nil {
        panic(err)
    }
}
// SetSeller sets the seller property value. Seller.
func (m *IosVppEBook) SetSeller(value *string)() {
    err := m.GetBackingStore().Set("seller", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalLicenseCount sets the totalLicenseCount property value. Total license count.
func (m *IosVppEBook) SetTotalLicenseCount(value *int32)() {
    err := m.GetBackingStore().Set("totalLicenseCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUsedLicenseCount sets the usedLicenseCount property value. Used license count.
func (m *IosVppEBook) SetUsedLicenseCount(value *int32)() {
    err := m.GetBackingStore().Set("usedLicenseCount", value)
    if err != nil {
        panic(err)
    }
}
// SetVppOrganizationName sets the vppOrganizationName property value. The Vpp token's organization name.
func (m *IosVppEBook) SetVppOrganizationName(value *string)() {
    err := m.GetBackingStore().Set("vppOrganizationName", value)
    if err != nil {
        panic(err)
    }
}
// SetVppTokenId sets the vppTokenId property value. The Vpp token ID.
func (m *IosVppEBook) SetVppTokenId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)() {
    err := m.GetBackingStore().Set("vppTokenId", value)
    if err != nil {
        panic(err)
    }
}
type IosVppEBookable interface {
    ManagedEBookable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppleId()(*string)
    GetGenres()([]string)
    GetLanguage()(*string)
    GetSeller()(*string)
    GetTotalLicenseCount()(*int32)
    GetUsedLicenseCount()(*int32)
    GetVppOrganizationName()(*string)
    GetVppTokenId()(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
    SetAppleId(value *string)()
    SetGenres(value []string)()
    SetLanguage(value *string)()
    SetSeller(value *string)()
    SetTotalLicenseCount(value *int32)()
    SetUsedLicenseCount(value *int32)()
    SetVppOrganizationName(value *string)()
    SetVppTokenId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)()
}
