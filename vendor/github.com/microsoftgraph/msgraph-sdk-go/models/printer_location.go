package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type PrinterLocation struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewPrinterLocation instantiates a new PrinterLocation and sets the default values.
func NewPrinterLocation()(*PrinterLocation) {
    m := &PrinterLocation{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreatePrinterLocationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrinterLocationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrinterLocation(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *PrinterLocation) GetAdditionalData()(map[string]any) {
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
// GetAltitudeInMeters gets the altitudeInMeters property value. The altitude, in meters, that the printer is located at.
// returns a *int32 when successful
func (m *PrinterLocation) GetAltitudeInMeters()(*int32) {
    val, err := m.GetBackingStore().Get("altitudeInMeters")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *PrinterLocation) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBuilding gets the building property value. The building that the printer is located in.
// returns a *string when successful
func (m *PrinterLocation) GetBuilding()(*string) {
    val, err := m.GetBackingStore().Get("building")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCity gets the city property value. The city that the printer is located in.
// returns a *string when successful
func (m *PrinterLocation) GetCity()(*string) {
    val, err := m.GetBackingStore().Get("city")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCountryOrRegion gets the countryOrRegion property value. The country or region that the printer is located in.
// returns a *string when successful
func (m *PrinterLocation) GetCountryOrRegion()(*string) {
    val, err := m.GetBackingStore().Get("countryOrRegion")
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
func (m *PrinterLocation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["altitudeInMeters"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAltitudeInMeters(val)
        }
        return nil
    }
    res["building"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBuilding(val)
        }
        return nil
    }
    res["city"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCity(val)
        }
        return nil
    }
    res["countryOrRegion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCountryOrRegion(val)
        }
        return nil
    }
    res["floor"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFloor(val)
        }
        return nil
    }
    res["floorDescription"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFloorDescription(val)
        }
        return nil
    }
    res["latitude"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLatitude(val)
        }
        return nil
    }
    res["longitude"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLongitude(val)
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
    res["organization"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetOrganization(res)
        }
        return nil
    }
    res["postalCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPostalCode(val)
        }
        return nil
    }
    res["roomDescription"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRoomDescription(val)
        }
        return nil
    }
    res["roomName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRoomName(val)
        }
        return nil
    }
    res["site"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSite(val)
        }
        return nil
    }
    res["stateOrProvince"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStateOrProvince(val)
        }
        return nil
    }
    res["streetAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStreetAddress(val)
        }
        return nil
    }
    res["subdivision"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSubdivision(res)
        }
        return nil
    }
    res["subunit"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSubunit(res)
        }
        return nil
    }
    return res
}
// GetFloor gets the floor property value. The floor that the printer is located on. Only numerical values are supported right now.
// returns a *string when successful
func (m *PrinterLocation) GetFloor()(*string) {
    val, err := m.GetBackingStore().Get("floor")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFloorDescription gets the floorDescription property value. The description of the floor that the printer is located on.
// returns a *string when successful
func (m *PrinterLocation) GetFloorDescription()(*string) {
    val, err := m.GetBackingStore().Get("floorDescription")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLatitude gets the latitude property value. The latitude that the printer is located at.
// returns a *float64 when successful
func (m *PrinterLocation) GetLatitude()(*float64) {
    val, err := m.GetBackingStore().Get("latitude")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetLongitude gets the longitude property value. The longitude that the printer is located at.
// returns a *float64 when successful
func (m *PrinterLocation) GetLongitude()(*float64) {
    val, err := m.GetBackingStore().Get("longitude")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *PrinterLocation) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOrganization gets the organization property value. The organizational hierarchy that the printer belongs to. The elements should be in hierarchical order.
// returns a []string when successful
func (m *PrinterLocation) GetOrganization()([]string) {
    val, err := m.GetBackingStore().Get("organization")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetPostalCode gets the postalCode property value. The postal code that the printer is located in.
// returns a *string when successful
func (m *PrinterLocation) GetPostalCode()(*string) {
    val, err := m.GetBackingStore().Get("postalCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRoomDescription gets the roomDescription property value. The description of the room that the printer is located in.
// returns a *string when successful
func (m *PrinterLocation) GetRoomDescription()(*string) {
    val, err := m.GetBackingStore().Get("roomDescription")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRoomName gets the roomName property value. The room that the printer is located in. Only numerical values are supported right now.
// returns a *string when successful
func (m *PrinterLocation) GetRoomName()(*string) {
    val, err := m.GetBackingStore().Get("roomName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSite gets the site property value. The site that the printer is located in.
// returns a *string when successful
func (m *PrinterLocation) GetSite()(*string) {
    val, err := m.GetBackingStore().Get("site")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStateOrProvince gets the stateOrProvince property value. The state or province that the printer is located in.
// returns a *string when successful
func (m *PrinterLocation) GetStateOrProvince()(*string) {
    val, err := m.GetBackingStore().Get("stateOrProvince")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStreetAddress gets the streetAddress property value. The street address where the printer is located.
// returns a *string when successful
func (m *PrinterLocation) GetStreetAddress()(*string) {
    val, err := m.GetBackingStore().Get("streetAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSubdivision gets the subdivision property value. The subdivision that the printer is located in. The elements should be in hierarchical order.
// returns a []string when successful
func (m *PrinterLocation) GetSubdivision()([]string) {
    val, err := m.GetBackingStore().Get("subdivision")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSubunit gets the subunit property value. The subunit property
// returns a []string when successful
func (m *PrinterLocation) GetSubunit()([]string) {
    val, err := m.GetBackingStore().Get("subunit")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrinterLocation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("altitudeInMeters", m.GetAltitudeInMeters())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("building", m.GetBuilding())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("city", m.GetCity())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("countryOrRegion", m.GetCountryOrRegion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("floor", m.GetFloor())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("floorDescription", m.GetFloorDescription())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("latitude", m.GetLatitude())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("longitude", m.GetLongitude())
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
    if m.GetOrganization() != nil {
        err := writer.WriteCollectionOfStringValues("organization", m.GetOrganization())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("postalCode", m.GetPostalCode())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("roomDescription", m.GetRoomDescription())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("roomName", m.GetRoomName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("site", m.GetSite())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("stateOrProvince", m.GetStateOrProvince())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("streetAddress", m.GetStreetAddress())
        if err != nil {
            return err
        }
    }
    if m.GetSubdivision() != nil {
        err := writer.WriteCollectionOfStringValues("subdivision", m.GetSubdivision())
        if err != nil {
            return err
        }
    }
    if m.GetSubunit() != nil {
        err := writer.WriteCollectionOfStringValues("subunit", m.GetSubunit())
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
func (m *PrinterLocation) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAltitudeInMeters sets the altitudeInMeters property value. The altitude, in meters, that the printer is located at.
func (m *PrinterLocation) SetAltitudeInMeters(value *int32)() {
    err := m.GetBackingStore().Set("altitudeInMeters", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *PrinterLocation) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBuilding sets the building property value. The building that the printer is located in.
func (m *PrinterLocation) SetBuilding(value *string)() {
    err := m.GetBackingStore().Set("building", value)
    if err != nil {
        panic(err)
    }
}
// SetCity sets the city property value. The city that the printer is located in.
func (m *PrinterLocation) SetCity(value *string)() {
    err := m.GetBackingStore().Set("city", value)
    if err != nil {
        panic(err)
    }
}
// SetCountryOrRegion sets the countryOrRegion property value. The country or region that the printer is located in.
func (m *PrinterLocation) SetCountryOrRegion(value *string)() {
    err := m.GetBackingStore().Set("countryOrRegion", value)
    if err != nil {
        panic(err)
    }
}
// SetFloor sets the floor property value. The floor that the printer is located on. Only numerical values are supported right now.
func (m *PrinterLocation) SetFloor(value *string)() {
    err := m.GetBackingStore().Set("floor", value)
    if err != nil {
        panic(err)
    }
}
// SetFloorDescription sets the floorDescription property value. The description of the floor that the printer is located on.
func (m *PrinterLocation) SetFloorDescription(value *string)() {
    err := m.GetBackingStore().Set("floorDescription", value)
    if err != nil {
        panic(err)
    }
}
// SetLatitude sets the latitude property value. The latitude that the printer is located at.
func (m *PrinterLocation) SetLatitude(value *float64)() {
    err := m.GetBackingStore().Set("latitude", value)
    if err != nil {
        panic(err)
    }
}
// SetLongitude sets the longitude property value. The longitude that the printer is located at.
func (m *PrinterLocation) SetLongitude(value *float64)() {
    err := m.GetBackingStore().Set("longitude", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *PrinterLocation) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOrganization sets the organization property value. The organizational hierarchy that the printer belongs to. The elements should be in hierarchical order.
func (m *PrinterLocation) SetOrganization(value []string)() {
    err := m.GetBackingStore().Set("organization", value)
    if err != nil {
        panic(err)
    }
}
// SetPostalCode sets the postalCode property value. The postal code that the printer is located in.
func (m *PrinterLocation) SetPostalCode(value *string)() {
    err := m.GetBackingStore().Set("postalCode", value)
    if err != nil {
        panic(err)
    }
}
// SetRoomDescription sets the roomDescription property value. The description of the room that the printer is located in.
func (m *PrinterLocation) SetRoomDescription(value *string)() {
    err := m.GetBackingStore().Set("roomDescription", value)
    if err != nil {
        panic(err)
    }
}
// SetRoomName sets the roomName property value. The room that the printer is located in. Only numerical values are supported right now.
func (m *PrinterLocation) SetRoomName(value *string)() {
    err := m.GetBackingStore().Set("roomName", value)
    if err != nil {
        panic(err)
    }
}
// SetSite sets the site property value. The site that the printer is located in.
func (m *PrinterLocation) SetSite(value *string)() {
    err := m.GetBackingStore().Set("site", value)
    if err != nil {
        panic(err)
    }
}
// SetStateOrProvince sets the stateOrProvince property value. The state or province that the printer is located in.
func (m *PrinterLocation) SetStateOrProvince(value *string)() {
    err := m.GetBackingStore().Set("stateOrProvince", value)
    if err != nil {
        panic(err)
    }
}
// SetStreetAddress sets the streetAddress property value. The street address where the printer is located.
func (m *PrinterLocation) SetStreetAddress(value *string)() {
    err := m.GetBackingStore().Set("streetAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetSubdivision sets the subdivision property value. The subdivision that the printer is located in. The elements should be in hierarchical order.
func (m *PrinterLocation) SetSubdivision(value []string)() {
    err := m.GetBackingStore().Set("subdivision", value)
    if err != nil {
        panic(err)
    }
}
// SetSubunit sets the subunit property value. The subunit property
func (m *PrinterLocation) SetSubunit(value []string)() {
    err := m.GetBackingStore().Set("subunit", value)
    if err != nil {
        panic(err)
    }
}
type PrinterLocationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAltitudeInMeters()(*int32)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBuilding()(*string)
    GetCity()(*string)
    GetCountryOrRegion()(*string)
    GetFloor()(*string)
    GetFloorDescription()(*string)
    GetLatitude()(*float64)
    GetLongitude()(*float64)
    GetOdataType()(*string)
    GetOrganization()([]string)
    GetPostalCode()(*string)
    GetRoomDescription()(*string)
    GetRoomName()(*string)
    GetSite()(*string)
    GetStateOrProvince()(*string)
    GetStreetAddress()(*string)
    GetSubdivision()([]string)
    GetSubunit()([]string)
    SetAltitudeInMeters(value *int32)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBuilding(value *string)()
    SetCity(value *string)()
    SetCountryOrRegion(value *string)()
    SetFloor(value *string)()
    SetFloorDescription(value *string)()
    SetLatitude(value *float64)()
    SetLongitude(value *float64)()
    SetOdataType(value *string)()
    SetOrganization(value []string)()
    SetPostalCode(value *string)()
    SetRoomDescription(value *string)()
    SetRoomName(value *string)()
    SetSite(value *string)()
    SetStateOrProvince(value *string)()
    SetStreetAddress(value *string)()
    SetSubdivision(value []string)()
    SetSubunit(value []string)()
}
