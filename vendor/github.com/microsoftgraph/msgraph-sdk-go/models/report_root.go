package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ReportRoot struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewReportRoot instantiates a new ReportRoot and sets the default values.
func NewReportRoot()(*ReportRoot) {
    m := &ReportRoot{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateReportRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateReportRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewReportRoot(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ReportRoot) GetAdditionalData()(map[string]any) {
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
// GetAuthenticationMethods gets the authenticationMethods property value. Container for navigation properties for Microsoft Entra authentication methods resources.
// returns a AuthenticationMethodsRootable when successful
func (m *ReportRoot) GetAuthenticationMethods()(AuthenticationMethodsRootable) {
    val, err := m.GetBackingStore().Get("authenticationMethods")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AuthenticationMethodsRootable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *ReportRoot) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDailyPrintUsageByPrinter gets the dailyPrintUsageByPrinter property value. Retrieve a list of daily print usage summaries, grouped by printer.
// returns a []PrintUsageByPrinterable when successful
func (m *ReportRoot) GetDailyPrintUsageByPrinter()([]PrintUsageByPrinterable) {
    val, err := m.GetBackingStore().Get("dailyPrintUsageByPrinter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintUsageByPrinterable)
    }
    return nil
}
// GetDailyPrintUsageByUser gets the dailyPrintUsageByUser property value. Retrieve a list of daily print usage summaries, grouped by user.
// returns a []PrintUsageByUserable when successful
func (m *ReportRoot) GetDailyPrintUsageByUser()([]PrintUsageByUserable) {
    val, err := m.GetBackingStore().Get("dailyPrintUsageByUser")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintUsageByUserable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ReportRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["authenticationMethods"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAuthenticationMethodsRootFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAuthenticationMethods(val.(AuthenticationMethodsRootable))
        }
        return nil
    }
    res["dailyPrintUsageByPrinter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrintUsageByPrinterFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintUsageByPrinterable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrintUsageByPrinterable)
                }
            }
            m.SetDailyPrintUsageByPrinter(res)
        }
        return nil
    }
    res["dailyPrintUsageByUser"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrintUsageByUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintUsageByUserable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrintUsageByUserable)
                }
            }
            m.SetDailyPrintUsageByUser(res)
        }
        return nil
    }
    res["monthlyPrintUsageByPrinter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrintUsageByPrinterFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintUsageByPrinterable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrintUsageByPrinterable)
                }
            }
            m.SetMonthlyPrintUsageByPrinter(res)
        }
        return nil
    }
    res["monthlyPrintUsageByUser"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePrintUsageByUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]PrintUsageByUserable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(PrintUsageByUserable)
                }
            }
            m.SetMonthlyPrintUsageByUser(res)
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
    res["partners"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePartnersFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPartners(val.(Partnersable))
        }
        return nil
    }
    res["security"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSecurityReportsRootFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSecurity(val.(SecurityReportsRootable))
        }
        return nil
    }
    return res
}
// GetMonthlyPrintUsageByPrinter gets the monthlyPrintUsageByPrinter property value. Retrieve a list of monthly print usage summaries, grouped by printer.
// returns a []PrintUsageByPrinterable when successful
func (m *ReportRoot) GetMonthlyPrintUsageByPrinter()([]PrintUsageByPrinterable) {
    val, err := m.GetBackingStore().Get("monthlyPrintUsageByPrinter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintUsageByPrinterable)
    }
    return nil
}
// GetMonthlyPrintUsageByUser gets the monthlyPrintUsageByUser property value. Retrieve a list of monthly print usage summaries, grouped by user.
// returns a []PrintUsageByUserable when successful
func (m *ReportRoot) GetMonthlyPrintUsageByUser()([]PrintUsageByUserable) {
    val, err := m.GetBackingStore().Get("monthlyPrintUsageByUser")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]PrintUsageByUserable)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ReportRoot) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPartners gets the partners property value. Represents billing details for a Microsoft direct partner.
// returns a Partnersable when successful
func (m *ReportRoot) GetPartners()(Partnersable) {
    val, err := m.GetBackingStore().Get("partners")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Partnersable)
    }
    return nil
}
// GetSecurity gets the security property value. Represents an abstract type that contains resources for attack simulation and training reports.
// returns a SecurityReportsRootable when successful
func (m *ReportRoot) GetSecurity()(SecurityReportsRootable) {
    val, err := m.GetBackingStore().Get("security")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SecurityReportsRootable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ReportRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("authenticationMethods", m.GetAuthenticationMethods())
        if err != nil {
            return err
        }
    }
    if m.GetDailyPrintUsageByPrinter() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDailyPrintUsageByPrinter()))
        for i, v := range m.GetDailyPrintUsageByPrinter() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("dailyPrintUsageByPrinter", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDailyPrintUsageByUser() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDailyPrintUsageByUser()))
        for i, v := range m.GetDailyPrintUsageByUser() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("dailyPrintUsageByUser", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMonthlyPrintUsageByPrinter() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMonthlyPrintUsageByPrinter()))
        for i, v := range m.GetMonthlyPrintUsageByPrinter() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("monthlyPrintUsageByPrinter", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMonthlyPrintUsageByUser() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMonthlyPrintUsageByUser()))
        for i, v := range m.GetMonthlyPrintUsageByUser() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("monthlyPrintUsageByUser", cast)
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
        err := writer.WriteObjectValue("partners", m.GetPartners())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("security", m.GetSecurity())
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
func (m *ReportRoot) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAuthenticationMethods sets the authenticationMethods property value. Container for navigation properties for Microsoft Entra authentication methods resources.
func (m *ReportRoot) SetAuthenticationMethods(value AuthenticationMethodsRootable)() {
    err := m.GetBackingStore().Set("authenticationMethods", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ReportRoot) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDailyPrintUsageByPrinter sets the dailyPrintUsageByPrinter property value. Retrieve a list of daily print usage summaries, grouped by printer.
func (m *ReportRoot) SetDailyPrintUsageByPrinter(value []PrintUsageByPrinterable)() {
    err := m.GetBackingStore().Set("dailyPrintUsageByPrinter", value)
    if err != nil {
        panic(err)
    }
}
// SetDailyPrintUsageByUser sets the dailyPrintUsageByUser property value. Retrieve a list of daily print usage summaries, grouped by user.
func (m *ReportRoot) SetDailyPrintUsageByUser(value []PrintUsageByUserable)() {
    err := m.GetBackingStore().Set("dailyPrintUsageByUser", value)
    if err != nil {
        panic(err)
    }
}
// SetMonthlyPrintUsageByPrinter sets the monthlyPrintUsageByPrinter property value. Retrieve a list of monthly print usage summaries, grouped by printer.
func (m *ReportRoot) SetMonthlyPrintUsageByPrinter(value []PrintUsageByPrinterable)() {
    err := m.GetBackingStore().Set("monthlyPrintUsageByPrinter", value)
    if err != nil {
        panic(err)
    }
}
// SetMonthlyPrintUsageByUser sets the monthlyPrintUsageByUser property value. Retrieve a list of monthly print usage summaries, grouped by user.
func (m *ReportRoot) SetMonthlyPrintUsageByUser(value []PrintUsageByUserable)() {
    err := m.GetBackingStore().Set("monthlyPrintUsageByUser", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ReportRoot) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPartners sets the partners property value. Represents billing details for a Microsoft direct partner.
func (m *ReportRoot) SetPartners(value Partnersable)() {
    err := m.GetBackingStore().Set("partners", value)
    if err != nil {
        panic(err)
    }
}
// SetSecurity sets the security property value. Represents an abstract type that contains resources for attack simulation and training reports.
func (m *ReportRoot) SetSecurity(value SecurityReportsRootable)() {
    err := m.GetBackingStore().Set("security", value)
    if err != nil {
        panic(err)
    }
}
type ReportRootable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAuthenticationMethods()(AuthenticationMethodsRootable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDailyPrintUsageByPrinter()([]PrintUsageByPrinterable)
    GetDailyPrintUsageByUser()([]PrintUsageByUserable)
    GetMonthlyPrintUsageByPrinter()([]PrintUsageByPrinterable)
    GetMonthlyPrintUsageByUser()([]PrintUsageByUserable)
    GetOdataType()(*string)
    GetPartners()(Partnersable)
    GetSecurity()(SecurityReportsRootable)
    SetAuthenticationMethods(value AuthenticationMethodsRootable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDailyPrintUsageByPrinter(value []PrintUsageByPrinterable)()
    SetDailyPrintUsageByUser(value []PrintUsageByUserable)()
    SetMonthlyPrintUsageByPrinter(value []PrintUsageByPrinterable)()
    SetMonthlyPrintUsageByUser(value []PrintUsageByUserable)()
    SetOdataType(value *string)()
    SetPartners(value Partnersable)()
    SetSecurity(value SecurityReportsRootable)()
}
