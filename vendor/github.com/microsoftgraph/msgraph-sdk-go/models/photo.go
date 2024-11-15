package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type Photo struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewPhoto instantiates a new Photo and sets the default values.
func NewPhoto()(*Photo) {
    m := &Photo{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreatePhotoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePhotoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPhoto(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Photo) GetAdditionalData()(map[string]any) {
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
func (m *Photo) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCameraMake gets the cameraMake property value. Camera manufacturer. Read-only.
// returns a *string when successful
func (m *Photo) GetCameraMake()(*string) {
    val, err := m.GetBackingStore().Get("cameraMake")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCameraModel gets the cameraModel property value. Camera model. Read-only.
// returns a *string when successful
func (m *Photo) GetCameraModel()(*string) {
    val, err := m.GetBackingStore().Get("cameraModel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExposureDenominator gets the exposureDenominator property value. The denominator for the exposure time fraction from the camera. Read-only.
// returns a *float64 when successful
func (m *Photo) GetExposureDenominator()(*float64) {
    val, err := m.GetBackingStore().Get("exposureDenominator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetExposureNumerator gets the exposureNumerator property value. The numerator for the exposure time fraction from the camera. Read-only.
// returns a *float64 when successful
func (m *Photo) GetExposureNumerator()(*float64) {
    val, err := m.GetBackingStore().Get("exposureNumerator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Photo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["cameraMake"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCameraMake(val)
        }
        return nil
    }
    res["cameraModel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCameraModel(val)
        }
        return nil
    }
    res["exposureDenominator"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExposureDenominator(val)
        }
        return nil
    }
    res["exposureNumerator"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExposureNumerator(val)
        }
        return nil
    }
    res["fNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFNumber(val)
        }
        return nil
    }
    res["focalLength"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFocalLength(val)
        }
        return nil
    }
    res["iso"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIso(val)
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
    res["orientation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOrientation(val)
        }
        return nil
    }
    res["takenDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTakenDateTime(val)
        }
        return nil
    }
    return res
}
// GetFNumber gets the fNumber property value. The F-stop value from the camera. Read-only.
// returns a *float64 when successful
func (m *Photo) GetFNumber()(*float64) {
    val, err := m.GetBackingStore().Get("fNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetFocalLength gets the focalLength property value. The focal length from the camera. Read-only.
// returns a *float64 when successful
func (m *Photo) GetFocalLength()(*float64) {
    val, err := m.GetBackingStore().Get("focalLength")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetIso gets the iso property value. The ISO value from the camera. Read-only.
// returns a *int32 when successful
func (m *Photo) GetIso()(*int32) {
    val, err := m.GetBackingStore().Get("iso")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *Photo) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOrientation gets the orientation property value. The orientation value from the camera. Writable on OneDrive Personal.
// returns a *int32 when successful
func (m *Photo) GetOrientation()(*int32) {
    val, err := m.GetBackingStore().Get("orientation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTakenDateTime gets the takenDateTime property value. Represents the date and time the photo was taken. Read-only.
// returns a *Time when successful
func (m *Photo) GetTakenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("takenDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Photo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("cameraMake", m.GetCameraMake())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("cameraModel", m.GetCameraModel())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("exposureDenominator", m.GetExposureDenominator())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("exposureNumerator", m.GetExposureNumerator())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("fNumber", m.GetFNumber())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteFloat64Value("focalLength", m.GetFocalLength())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("iso", m.GetIso())
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
        err := writer.WriteInt32Value("orientation", m.GetOrientation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteTimeValue("takenDateTime", m.GetTakenDateTime())
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
func (m *Photo) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Photo) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCameraMake sets the cameraMake property value. Camera manufacturer. Read-only.
func (m *Photo) SetCameraMake(value *string)() {
    err := m.GetBackingStore().Set("cameraMake", value)
    if err != nil {
        panic(err)
    }
}
// SetCameraModel sets the cameraModel property value. Camera model. Read-only.
func (m *Photo) SetCameraModel(value *string)() {
    err := m.GetBackingStore().Set("cameraModel", value)
    if err != nil {
        panic(err)
    }
}
// SetExposureDenominator sets the exposureDenominator property value. The denominator for the exposure time fraction from the camera. Read-only.
func (m *Photo) SetExposureDenominator(value *float64)() {
    err := m.GetBackingStore().Set("exposureDenominator", value)
    if err != nil {
        panic(err)
    }
}
// SetExposureNumerator sets the exposureNumerator property value. The numerator for the exposure time fraction from the camera. Read-only.
func (m *Photo) SetExposureNumerator(value *float64)() {
    err := m.GetBackingStore().Set("exposureNumerator", value)
    if err != nil {
        panic(err)
    }
}
// SetFNumber sets the fNumber property value. The F-stop value from the camera. Read-only.
func (m *Photo) SetFNumber(value *float64)() {
    err := m.GetBackingStore().Set("fNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetFocalLength sets the focalLength property value. The focal length from the camera. Read-only.
func (m *Photo) SetFocalLength(value *float64)() {
    err := m.GetBackingStore().Set("focalLength", value)
    if err != nil {
        panic(err)
    }
}
// SetIso sets the iso property value. The ISO value from the camera. Read-only.
func (m *Photo) SetIso(value *int32)() {
    err := m.GetBackingStore().Set("iso", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Photo) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetOrientation sets the orientation property value. The orientation value from the camera. Writable on OneDrive Personal.
func (m *Photo) SetOrientation(value *int32)() {
    err := m.GetBackingStore().Set("orientation", value)
    if err != nil {
        panic(err)
    }
}
// SetTakenDateTime sets the takenDateTime property value. Represents the date and time the photo was taken. Read-only.
func (m *Photo) SetTakenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("takenDateTime", value)
    if err != nil {
        panic(err)
    }
}
type Photoable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCameraMake()(*string)
    GetCameraModel()(*string)
    GetExposureDenominator()(*float64)
    GetExposureNumerator()(*float64)
    GetFNumber()(*float64)
    GetFocalLength()(*float64)
    GetIso()(*int32)
    GetOdataType()(*string)
    GetOrientation()(*int32)
    GetTakenDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCameraMake(value *string)()
    SetCameraModel(value *string)()
    SetExposureDenominator(value *float64)()
    SetExposureNumerator(value *float64)()
    SetFNumber(value *float64)()
    SetFocalLength(value *float64)()
    SetIso(value *int32)()
    SetOdataType(value *string)()
    SetOrientation(value *int32)()
    SetTakenDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
}
