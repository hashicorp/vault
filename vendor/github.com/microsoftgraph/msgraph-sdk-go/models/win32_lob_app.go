package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Win32LobApp contains properties and inherited properties for Win32 apps.
type Win32LobApp struct {
    MobileLobApp
}
// NewWin32LobApp instantiates a new Win32LobApp and sets the default values.
func NewWin32LobApp()(*Win32LobApp) {
    m := &Win32LobApp{
        MobileLobApp: *NewMobileLobApp(),
    }
    odataTypeValue := "#microsoft.graph.win32LobApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWin32LobAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWin32LobAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWin32LobApp(), nil
}
// GetApplicableArchitectures gets the applicableArchitectures property value. Contains properties for Windows architecture.
// returns a *WindowsArchitecture when successful
func (m *Win32LobApp) GetApplicableArchitectures()(*WindowsArchitecture) {
    val, err := m.GetBackingStore().Get("applicableArchitectures")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsArchitecture)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Win32LobApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileLobApp.GetFieldDeserializers()
    res["applicableArchitectures"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsArchitecture)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicableArchitectures(val.(*WindowsArchitecture))
        }
        return nil
    }
    res["installCommandLine"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstallCommandLine(val)
        }
        return nil
    }
    res["installExperience"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWin32LobAppInstallExperienceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstallExperience(val.(Win32LobAppInstallExperienceable))
        }
        return nil
    }
    res["minimumCpuSpeedInMHz"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumCpuSpeedInMHz(val)
        }
        return nil
    }
    res["minimumFreeDiskSpaceInMB"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumFreeDiskSpaceInMB(val)
        }
        return nil
    }
    res["minimumMemoryInMB"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumMemoryInMB(val)
        }
        return nil
    }
    res["minimumNumberOfProcessors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumNumberOfProcessors(val)
        }
        return nil
    }
    res["minimumSupportedWindowsRelease"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumSupportedWindowsRelease(val)
        }
        return nil
    }
    res["msiInformation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWin32LobAppMsiInformationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMsiInformation(val.(Win32LobAppMsiInformationable))
        }
        return nil
    }
    res["returnCodes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWin32LobAppReturnCodeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Win32LobAppReturnCodeable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Win32LobAppReturnCodeable)
                }
            }
            m.SetReturnCodes(res)
        }
        return nil
    }
    res["rules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWin32LobAppRuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Win32LobAppRuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Win32LobAppRuleable)
                }
            }
            m.SetRules(res)
        }
        return nil
    }
    res["setupFilePath"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSetupFilePath(val)
        }
        return nil
    }
    res["uninstallCommandLine"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUninstallCommandLine(val)
        }
        return nil
    }
    return res
}
// GetInstallCommandLine gets the installCommandLine property value. The command line to install this app
// returns a *string when successful
func (m *Win32LobApp) GetInstallCommandLine()(*string) {
    val, err := m.GetBackingStore().Get("installCommandLine")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetInstallExperience gets the installExperience property value. The install experience for this app.
// returns a Win32LobAppInstallExperienceable when successful
func (m *Win32LobApp) GetInstallExperience()(Win32LobAppInstallExperienceable) {
    val, err := m.GetBackingStore().Get("installExperience")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Win32LobAppInstallExperienceable)
    }
    return nil
}
// GetMinimumCpuSpeedInMHz gets the minimumCpuSpeedInMHz property value. The value for the minimum CPU speed which is required to install this app.
// returns a *int32 when successful
func (m *Win32LobApp) GetMinimumCpuSpeedInMHz()(*int32) {
    val, err := m.GetBackingStore().Get("minimumCpuSpeedInMHz")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMinimumFreeDiskSpaceInMB gets the minimumFreeDiskSpaceInMB property value. The value for the minimum free disk space which is required to install this app.
// returns a *int32 when successful
func (m *Win32LobApp) GetMinimumFreeDiskSpaceInMB()(*int32) {
    val, err := m.GetBackingStore().Get("minimumFreeDiskSpaceInMB")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMinimumMemoryInMB gets the minimumMemoryInMB property value. The value for the minimum physical memory which is required to install this app.
// returns a *int32 when successful
func (m *Win32LobApp) GetMinimumMemoryInMB()(*int32) {
    val, err := m.GetBackingStore().Get("minimumMemoryInMB")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMinimumNumberOfProcessors gets the minimumNumberOfProcessors property value. The value for the minimum number of processors which is required to install this app.
// returns a *int32 when successful
func (m *Win32LobApp) GetMinimumNumberOfProcessors()(*int32) {
    val, err := m.GetBackingStore().Get("minimumNumberOfProcessors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMinimumSupportedWindowsRelease gets the minimumSupportedWindowsRelease property value. The value for the minimum supported windows release.
// returns a *string when successful
func (m *Win32LobApp) GetMinimumSupportedWindowsRelease()(*string) {
    val, err := m.GetBackingStore().Get("minimumSupportedWindowsRelease")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMsiInformation gets the msiInformation property value. The MSI details if this Win32 app is an MSI app.
// returns a Win32LobAppMsiInformationable when successful
func (m *Win32LobApp) GetMsiInformation()(Win32LobAppMsiInformationable) {
    val, err := m.GetBackingStore().Get("msiInformation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Win32LobAppMsiInformationable)
    }
    return nil
}
// GetReturnCodes gets the returnCodes property value. The return codes for post installation behavior.
// returns a []Win32LobAppReturnCodeable when successful
func (m *Win32LobApp) GetReturnCodes()([]Win32LobAppReturnCodeable) {
    val, err := m.GetBackingStore().Get("returnCodes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Win32LobAppReturnCodeable)
    }
    return nil
}
// GetRules gets the rules property value. The detection and requirement rules for this app.
// returns a []Win32LobAppRuleable when successful
func (m *Win32LobApp) GetRules()([]Win32LobAppRuleable) {
    val, err := m.GetBackingStore().Get("rules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Win32LobAppRuleable)
    }
    return nil
}
// GetSetupFilePath gets the setupFilePath property value. The relative path of the setup file in the encrypted Win32LobApp package.
// returns a *string when successful
func (m *Win32LobApp) GetSetupFilePath()(*string) {
    val, err := m.GetBackingStore().Get("setupFilePath")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUninstallCommandLine gets the uninstallCommandLine property value. The command line to uninstall this app
// returns a *string when successful
func (m *Win32LobApp) GetUninstallCommandLine()(*string) {
    val, err := m.GetBackingStore().Get("uninstallCommandLine")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Win32LobApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileLobApp.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetApplicableArchitectures() != nil {
        cast := (*m.GetApplicableArchitectures()).String()
        err = writer.WriteStringValue("applicableArchitectures", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("installCommandLine", m.GetInstallCommandLine())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("installExperience", m.GetInstallExperience())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("minimumCpuSpeedInMHz", m.GetMinimumCpuSpeedInMHz())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("minimumFreeDiskSpaceInMB", m.GetMinimumFreeDiskSpaceInMB())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("minimumMemoryInMB", m.GetMinimumMemoryInMB())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("minimumNumberOfProcessors", m.GetMinimumNumberOfProcessors())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("minimumSupportedWindowsRelease", m.GetMinimumSupportedWindowsRelease())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("msiInformation", m.GetMsiInformation())
        if err != nil {
            return err
        }
    }
    if m.GetReturnCodes() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetReturnCodes()))
        for i, v := range m.GetReturnCodes() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("returnCodes", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRules()))
        for i, v := range m.GetRules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("rules", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("setupFilePath", m.GetSetupFilePath())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("uninstallCommandLine", m.GetUninstallCommandLine())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicableArchitectures sets the applicableArchitectures property value. Contains properties for Windows architecture.
func (m *Win32LobApp) SetApplicableArchitectures(value *WindowsArchitecture)() {
    err := m.GetBackingStore().Set("applicableArchitectures", value)
    if err != nil {
        panic(err)
    }
}
// SetInstallCommandLine sets the installCommandLine property value. The command line to install this app
func (m *Win32LobApp) SetInstallCommandLine(value *string)() {
    err := m.GetBackingStore().Set("installCommandLine", value)
    if err != nil {
        panic(err)
    }
}
// SetInstallExperience sets the installExperience property value. The install experience for this app.
func (m *Win32LobApp) SetInstallExperience(value Win32LobAppInstallExperienceable)() {
    err := m.GetBackingStore().Set("installExperience", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumCpuSpeedInMHz sets the minimumCpuSpeedInMHz property value. The value for the minimum CPU speed which is required to install this app.
func (m *Win32LobApp) SetMinimumCpuSpeedInMHz(value *int32)() {
    err := m.GetBackingStore().Set("minimumCpuSpeedInMHz", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumFreeDiskSpaceInMB sets the minimumFreeDiskSpaceInMB property value. The value for the minimum free disk space which is required to install this app.
func (m *Win32LobApp) SetMinimumFreeDiskSpaceInMB(value *int32)() {
    err := m.GetBackingStore().Set("minimumFreeDiskSpaceInMB", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumMemoryInMB sets the minimumMemoryInMB property value. The value for the minimum physical memory which is required to install this app.
func (m *Win32LobApp) SetMinimumMemoryInMB(value *int32)() {
    err := m.GetBackingStore().Set("minimumMemoryInMB", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumNumberOfProcessors sets the minimumNumberOfProcessors property value. The value for the minimum number of processors which is required to install this app.
func (m *Win32LobApp) SetMinimumNumberOfProcessors(value *int32)() {
    err := m.GetBackingStore().Set("minimumNumberOfProcessors", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumSupportedWindowsRelease sets the minimumSupportedWindowsRelease property value. The value for the minimum supported windows release.
func (m *Win32LobApp) SetMinimumSupportedWindowsRelease(value *string)() {
    err := m.GetBackingStore().Set("minimumSupportedWindowsRelease", value)
    if err != nil {
        panic(err)
    }
}
// SetMsiInformation sets the msiInformation property value. The MSI details if this Win32 app is an MSI app.
func (m *Win32LobApp) SetMsiInformation(value Win32LobAppMsiInformationable)() {
    err := m.GetBackingStore().Set("msiInformation", value)
    if err != nil {
        panic(err)
    }
}
// SetReturnCodes sets the returnCodes property value. The return codes for post installation behavior.
func (m *Win32LobApp) SetReturnCodes(value []Win32LobAppReturnCodeable)() {
    err := m.GetBackingStore().Set("returnCodes", value)
    if err != nil {
        panic(err)
    }
}
// SetRules sets the rules property value. The detection and requirement rules for this app.
func (m *Win32LobApp) SetRules(value []Win32LobAppRuleable)() {
    err := m.GetBackingStore().Set("rules", value)
    if err != nil {
        panic(err)
    }
}
// SetSetupFilePath sets the setupFilePath property value. The relative path of the setup file in the encrypted Win32LobApp package.
func (m *Win32LobApp) SetSetupFilePath(value *string)() {
    err := m.GetBackingStore().Set("setupFilePath", value)
    if err != nil {
        panic(err)
    }
}
// SetUninstallCommandLine sets the uninstallCommandLine property value. The command line to uninstall this app
func (m *Win32LobApp) SetUninstallCommandLine(value *string)() {
    err := m.GetBackingStore().Set("uninstallCommandLine", value)
    if err != nil {
        panic(err)
    }
}
type Win32LobAppable interface {
    MobileLobAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicableArchitectures()(*WindowsArchitecture)
    GetInstallCommandLine()(*string)
    GetInstallExperience()(Win32LobAppInstallExperienceable)
    GetMinimumCpuSpeedInMHz()(*int32)
    GetMinimumFreeDiskSpaceInMB()(*int32)
    GetMinimumMemoryInMB()(*int32)
    GetMinimumNumberOfProcessors()(*int32)
    GetMinimumSupportedWindowsRelease()(*string)
    GetMsiInformation()(Win32LobAppMsiInformationable)
    GetReturnCodes()([]Win32LobAppReturnCodeable)
    GetRules()([]Win32LobAppRuleable)
    GetSetupFilePath()(*string)
    GetUninstallCommandLine()(*string)
    SetApplicableArchitectures(value *WindowsArchitecture)()
    SetInstallCommandLine(value *string)()
    SetInstallExperience(value Win32LobAppInstallExperienceable)()
    SetMinimumCpuSpeedInMHz(value *int32)()
    SetMinimumFreeDiskSpaceInMB(value *int32)()
    SetMinimumMemoryInMB(value *int32)()
    SetMinimumNumberOfProcessors(value *int32)()
    SetMinimumSupportedWindowsRelease(value *string)()
    SetMsiInformation(value Win32LobAppMsiInformationable)()
    SetReturnCodes(value []Win32LobAppReturnCodeable)()
    SetRules(value []Win32LobAppRuleable)()
    SetSetupFilePath(value *string)()
    SetUninstallCommandLine(value *string)()
}
