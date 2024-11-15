package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type VirtualEventPresenterDetails struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewVirtualEventPresenterDetails instantiates a new VirtualEventPresenterDetails and sets the default values.
func NewVirtualEventPresenterDetails()(*VirtualEventPresenterDetails) {
    m := &VirtualEventPresenterDetails{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateVirtualEventPresenterDetailsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateVirtualEventPresenterDetailsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewVirtualEventPresenterDetails(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *VirtualEventPresenterDetails) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *VirtualEventPresenterDetails) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBio gets the bio property value. Bio of the presenter.
// returns a ItemBodyable when successful
func (m *VirtualEventPresenterDetails) GetBio()(ItemBodyable) {
    val, err := m.GetBackingStore().Get("bio")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemBodyable)
    }
    return nil
}
// GetCompany gets the company property value. The presenter's company name.
// returns a *string when successful
func (m *VirtualEventPresenterDetails) GetCompany()(*string) {
    val, err := m.GetBackingStore().Get("company")
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
func (m *VirtualEventPresenterDetails) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["bio"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemBodyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBio(val.(ItemBodyable))
        }
        return nil
    }
    res["company"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompany(val)
        }
        return nil
    }
    res["jobTitle"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetJobTitle(val)
        }
        return nil
    }
    res["linkedInProfileWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLinkedInProfileWebUrl(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["personalSiteWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPersonalSiteWebUrl(val)
        }
        return nil
    }
    res["photo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPhoto(val)
        }
        return nil
    }
    res["twitterProfileWebUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTwitterProfileWebUrl(val)
        }
        return nil
    }
    return res
}
// GetJobTitle gets the jobTitle property value. The presenter's job title.
// returns a *string when successful
func (m *VirtualEventPresenterDetails) GetJobTitle()(*string) {
    val, err := m.GetBackingStore().Get("jobTitle")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLinkedInProfileWebUrl gets the linkedInProfileWebUrl property value. The presenter's LinkedIn profile URL.
// returns a *string when successful
func (m *VirtualEventPresenterDetails) GetLinkedInProfileWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("linkedInProfileWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *VirtualEventPresenterDetails) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPersonalSiteWebUrl gets the personalSiteWebUrl property value. The presenter's personal website URL.
// returns a *string when successful
func (m *VirtualEventPresenterDetails) GetPersonalSiteWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("personalSiteWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPhoto gets the photo property value. The content stream of the presenter's photo.
// returns a []byte when successful
func (m *VirtualEventPresenterDetails) GetPhoto()([]byte) {
    val, err := m.GetBackingStore().Get("photo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetTwitterProfileWebUrl gets the twitterProfileWebUrl property value. The presenter's Twitter profile URL.
// returns a *string when successful
func (m *VirtualEventPresenterDetails) GetTwitterProfileWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("twitterProfileWebUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *VirtualEventPresenterDetails) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("bio", m.GetBio())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("company", m.GetCompany())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("jobTitle", m.GetJobTitle())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("linkedInProfileWebUrl", m.GetLinkedInProfileWebUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("personalSiteWebUrl", m.GetPersonalSiteWebUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteByteArrayValue("photo", m.GetPhoto())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("twitterProfileWebUrl", m.GetTwitterProfileWebUrl())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *VirtualEventPresenterDetails) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *VirtualEventPresenterDetails) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBio sets the bio property value. Bio of the presenter.
func (m *VirtualEventPresenterDetails) SetBio(value ItemBodyable)() {
    err := m.GetBackingStore().Set("bio", value)
    if err != nil {
        panic(err)
    }
}
// SetCompany sets the company property value. The presenter's company name.
func (m *VirtualEventPresenterDetails) SetCompany(value *string)() {
    err := m.GetBackingStore().Set("company", value)
    if err != nil {
        panic(err)
    }
}
// SetJobTitle sets the jobTitle property value. The presenter's job title.
func (m *VirtualEventPresenterDetails) SetJobTitle(value *string)() {
    err := m.GetBackingStore().Set("jobTitle", value)
    if err != nil {
        panic(err)
    }
}
// SetLinkedInProfileWebUrl sets the linkedInProfileWebUrl property value. The presenter's LinkedIn profile URL.
func (m *VirtualEventPresenterDetails) SetLinkedInProfileWebUrl(value *string)() {
    err := m.GetBackingStore().Set("linkedInProfileWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *VirtualEventPresenterDetails) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPersonalSiteWebUrl sets the personalSiteWebUrl property value. The presenter's personal website URL.
func (m *VirtualEventPresenterDetails) SetPersonalSiteWebUrl(value *string)() {
    err := m.GetBackingStore().Set("personalSiteWebUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetPhoto sets the photo property value. The content stream of the presenter's photo.
func (m *VirtualEventPresenterDetails) SetPhoto(value []byte)() {
    err := m.GetBackingStore().Set("photo", value)
    if err != nil {
        panic(err)
    }
}
// SetTwitterProfileWebUrl sets the twitterProfileWebUrl property value. The presenter's Twitter profile URL.
func (m *VirtualEventPresenterDetails) SetTwitterProfileWebUrl(value *string)() {
    err := m.GetBackingStore().Set("twitterProfileWebUrl", value)
    if err != nil {
        panic(err)
    }
}
type VirtualEventPresenterDetailsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBio()(ItemBodyable)
    GetCompany()(*string)
    GetJobTitle()(*string)
    GetLinkedInProfileWebUrl()(*string)
    GetOdataType()(*string)
    GetPersonalSiteWebUrl()(*string)
    GetPhoto()([]byte)
    GetTwitterProfileWebUrl()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBio(value ItemBodyable)()
    SetCompany(value *string)()
    SetJobTitle(value *string)()
    SetLinkedInProfileWebUrl(value *string)()
    SetOdataType(value *string)()
    SetPersonalSiteWebUrl(value *string)()
    SetPhoto(value []byte)()
    SetTwitterProfileWebUrl(value *string)()
}
