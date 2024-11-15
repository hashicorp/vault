package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosVppApp contains properties and inherited properties for iOS Volume-Purchased Program (VPP) Apps.
type IosVppApp struct {
    MobileApp
}
// NewIosVppApp instantiates a new IosVppApp and sets the default values.
func NewIosVppApp()(*IosVppApp) {
    m := &IosVppApp{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.iosVppApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIosVppAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosVppAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosVppApp(), nil
}
// GetApplicableDeviceType gets the applicableDeviceType property value. The applicable iOS Device Type.
// returns a IosDeviceTypeable when successful
func (m *IosVppApp) GetApplicableDeviceType()(IosDeviceTypeable) {
    val, err := m.GetBackingStore().Get("applicableDeviceType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IosDeviceTypeable)
    }
    return nil
}
// GetAppStoreUrl gets the appStoreUrl property value. The store URL.
// returns a *string when successful
func (m *IosVppApp) GetAppStoreUrl()(*string) {
    val, err := m.GetBackingStore().Get("appStoreUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBundleId gets the bundleId property value. The Identity Name.
// returns a *string when successful
func (m *IosVppApp) GetBundleId()(*string) {
    val, err := m.GetBackingStore().Get("bundleId")
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
func (m *IosVppApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileApp.GetFieldDeserializers()
    res["applicableDeviceType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIosDeviceTypeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicableDeviceType(val.(IosDeviceTypeable))
        }
        return nil
    }
    res["appStoreUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppStoreUrl(val)
        }
        return nil
    }
    res["bundleId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBundleId(val)
        }
        return nil
    }
    res["licensingType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateVppLicensingTypeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLicensingType(val.(VppLicensingTypeable))
        }
        return nil
    }
    res["releaseDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReleaseDateTime(val)
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
    res["vppTokenAccountType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseVppTokenAccountType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVppTokenAccountType(val.(*VppTokenAccountType))
        }
        return nil
    }
    res["vppTokenAppleId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVppTokenAppleId(val)
        }
        return nil
    }
    res["vppTokenOrganizationName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVppTokenOrganizationName(val)
        }
        return nil
    }
    return res
}
// GetLicensingType gets the licensingType property value. The supported License Type.
// returns a VppLicensingTypeable when successful
func (m *IosVppApp) GetLicensingType()(VppLicensingTypeable) {
    val, err := m.GetBackingStore().Get("licensingType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(VppLicensingTypeable)
    }
    return nil
}
// GetReleaseDateTime gets the releaseDateTime property value. The VPP application release date and time.
// returns a *Time when successful
func (m *IosVppApp) GetReleaseDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("releaseDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetTotalLicenseCount gets the totalLicenseCount property value. The total number of VPP licenses.
// returns a *int32 when successful
func (m *IosVppApp) GetTotalLicenseCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalLicenseCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUsedLicenseCount gets the usedLicenseCount property value. The number of VPP licenses in use.
// returns a *int32 when successful
func (m *IosVppApp) GetUsedLicenseCount()(*int32) {
    val, err := m.GetBackingStore().Get("usedLicenseCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetVppTokenAccountType gets the vppTokenAccountType property value. Possible types of an Apple Volume Purchase Program token.
// returns a *VppTokenAccountType when successful
func (m *IosVppApp) GetVppTokenAccountType()(*VppTokenAccountType) {
    val, err := m.GetBackingStore().Get("vppTokenAccountType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*VppTokenAccountType)
    }
    return nil
}
// GetVppTokenAppleId gets the vppTokenAppleId property value. The Apple Id associated with the given Apple Volume Purchase Program Token.
// returns a *string when successful
func (m *IosVppApp) GetVppTokenAppleId()(*string) {
    val, err := m.GetBackingStore().Get("vppTokenAppleId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVppTokenOrganizationName gets the vppTokenOrganizationName property value. The organization associated with the Apple Volume Purchase Program Token
// returns a *string when successful
func (m *IosVppApp) GetVppTokenOrganizationName()(*string) {
    val, err := m.GetBackingStore().Get("vppTokenOrganizationName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IosVppApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileApp.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("applicableDeviceType", m.GetApplicableDeviceType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appStoreUrl", m.GetAppStoreUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("bundleId", m.GetBundleId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("licensingType", m.GetLicensingType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("releaseDateTime", m.GetReleaseDateTime())
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
    if m.GetVppTokenAccountType() != nil {
        cast := (*m.GetVppTokenAccountType()).String()
        err = writer.WriteStringValue("vppTokenAccountType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("vppTokenAppleId", m.GetVppTokenAppleId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("vppTokenOrganizationName", m.GetVppTokenOrganizationName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicableDeviceType sets the applicableDeviceType property value. The applicable iOS Device Type.
func (m *IosVppApp) SetApplicableDeviceType(value IosDeviceTypeable)() {
    err := m.GetBackingStore().Set("applicableDeviceType", value)
    if err != nil {
        panic(err)
    }
}
// SetAppStoreUrl sets the appStoreUrl property value. The store URL.
func (m *IosVppApp) SetAppStoreUrl(value *string)() {
    err := m.GetBackingStore().Set("appStoreUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetBundleId sets the bundleId property value. The Identity Name.
func (m *IosVppApp) SetBundleId(value *string)() {
    err := m.GetBackingStore().Set("bundleId", value)
    if err != nil {
        panic(err)
    }
}
// SetLicensingType sets the licensingType property value. The supported License Type.
func (m *IosVppApp) SetLicensingType(value VppLicensingTypeable)() {
    err := m.GetBackingStore().Set("licensingType", value)
    if err != nil {
        panic(err)
    }
}
// SetReleaseDateTime sets the releaseDateTime property value. The VPP application release date and time.
func (m *IosVppApp) SetReleaseDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("releaseDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalLicenseCount sets the totalLicenseCount property value. The total number of VPP licenses.
func (m *IosVppApp) SetTotalLicenseCount(value *int32)() {
    err := m.GetBackingStore().Set("totalLicenseCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUsedLicenseCount sets the usedLicenseCount property value. The number of VPP licenses in use.
func (m *IosVppApp) SetUsedLicenseCount(value *int32)() {
    err := m.GetBackingStore().Set("usedLicenseCount", value)
    if err != nil {
        panic(err)
    }
}
// SetVppTokenAccountType sets the vppTokenAccountType property value. Possible types of an Apple Volume Purchase Program token.
func (m *IosVppApp) SetVppTokenAccountType(value *VppTokenAccountType)() {
    err := m.GetBackingStore().Set("vppTokenAccountType", value)
    if err != nil {
        panic(err)
    }
}
// SetVppTokenAppleId sets the vppTokenAppleId property value. The Apple Id associated with the given Apple Volume Purchase Program Token.
func (m *IosVppApp) SetVppTokenAppleId(value *string)() {
    err := m.GetBackingStore().Set("vppTokenAppleId", value)
    if err != nil {
        panic(err)
    }
}
// SetVppTokenOrganizationName sets the vppTokenOrganizationName property value. The organization associated with the Apple Volume Purchase Program Token
func (m *IosVppApp) SetVppTokenOrganizationName(value *string)() {
    err := m.GetBackingStore().Set("vppTokenOrganizationName", value)
    if err != nil {
        panic(err)
    }
}
type IosVppAppable interface {
    MobileAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicableDeviceType()(IosDeviceTypeable)
    GetAppStoreUrl()(*string)
    GetBundleId()(*string)
    GetLicensingType()(VppLicensingTypeable)
    GetReleaseDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetTotalLicenseCount()(*int32)
    GetUsedLicenseCount()(*int32)
    GetVppTokenAccountType()(*VppTokenAccountType)
    GetVppTokenAppleId()(*string)
    GetVppTokenOrganizationName()(*string)
    SetApplicableDeviceType(value IosDeviceTypeable)()
    SetAppStoreUrl(value *string)()
    SetBundleId(value *string)()
    SetLicensingType(value VppLicensingTypeable)()
    SetReleaseDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetTotalLicenseCount(value *int32)()
    SetUsedLicenseCount(value *int32)()
    SetVppTokenAccountType(value *VppTokenAccountType)()
    SetVppTokenAppleId(value *string)()
    SetVppTokenOrganizationName(value *string)()
}
