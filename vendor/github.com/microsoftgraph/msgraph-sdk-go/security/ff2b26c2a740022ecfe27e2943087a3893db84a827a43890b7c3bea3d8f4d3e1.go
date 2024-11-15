package security

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
    idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae "github.com/microsoftgraph/msgraph-sdk-go/models/security"
)

type CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewCasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody instantiates a new CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody and sets the default values.
func NewCasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody()(*CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) {
    m := &CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateCasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBodyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBodyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) GetAdditionalData()(map[string]any) {
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
// GetAdditionalOptions gets the additionalOptions property value. The additionalOptions property
// returns a *AdditionalOptions when successful
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) GetAdditionalOptions()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.AdditionalOptions) {
    val, err := m.GetBackingStore().Get("additionalOptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.AdditionalOptions)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDescription gets the description property value. The description property
// returns a *string when successful
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The displayName property
// returns a *string when successful
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExportCriteria gets the exportCriteria property value. The exportCriteria property
// returns a *ExportCriteria when successful
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) GetExportCriteria()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportCriteria) {
    val, err := m.GetBackingStore().Get("exportCriteria")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportCriteria)
    }
    return nil
}
// GetExportLocation gets the exportLocation property value. The exportLocation property
// returns a *ExportLocation when successful
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) GetExportLocation()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportLocation) {
    val, err := m.GetBackingStore().Get("exportLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportLocation)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["additionalOptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ParseAdditionalOptions)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAdditionalOptions(val.(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.AdditionalOptions))
        }
        return nil
    }
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
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
    res["exportCriteria"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ParseExportCriteria)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExportCriteria(val.(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportCriteria))
        }
        return nil
    }
    res["exportLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ParseExportLocation)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExportLocation(val.(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportLocation))
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAdditionalOptions() != nil {
        cast := (*m.GetAdditionalOptions()).String()
        err := writer.WriteStringValue("additionalOptions", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("description", m.GetDescription())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetExportCriteria() != nil {
        cast := (*m.GetExportCriteria()).String()
        err := writer.WriteStringValue("exportCriteria", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetExportLocation() != nil {
        cast := (*m.GetExportLocation()).String()
        err := writer.WriteStringValue("exportLocation", &cast)
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
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalOptions sets the additionalOptions property value. The additionalOptions property
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) SetAdditionalOptions(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.AdditionalOptions)() {
    err := m.GetBackingStore().Set("additionalOptions", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDescription sets the description property value. The description property
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The displayName property
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetExportCriteria sets the exportCriteria property value. The exportCriteria property
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) SetExportCriteria(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportCriteria)() {
    err := m.GetBackingStore().Set("exportCriteria", value)
    if err != nil {
        panic(err)
    }
}
// SetExportLocation sets the exportLocation property value. The exportLocation property
func (m *CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBody) SetExportLocation(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportLocation)() {
    err := m.GetBackingStore().Set("exportLocation", value)
    if err != nil {
        panic(err)
    }
}
type CasesEdiscoveryCasesItemSearchesItemMicrosoftGraphSecurityExportReportExportReportPostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAdditionalOptions()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.AdditionalOptions)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetExportCriteria()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportCriteria)
    GetExportLocation()(*idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportLocation)
    SetAdditionalOptions(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.AdditionalOptions)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetExportCriteria(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportCriteria)()
    SetExportLocation(value *idd6d442c3cc83a389b8f0b8dd7ac355916e813c2882ff3aaa23331424ba827ae.ExportLocation)()
}
