package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type BackupRestoreRoot struct {
    Entity
}
// NewBackupRestoreRoot instantiates a new BackupRestoreRoot and sets the default values.
func NewBackupRestoreRoot()(*BackupRestoreRoot) {
    m := &BackupRestoreRoot{
        Entity: *NewEntity(),
    }
    return m
}
// CreateBackupRestoreRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateBackupRestoreRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewBackupRestoreRoot(), nil
}
// GetDriveInclusionRules gets the driveInclusionRules property value. The list of drive inclusion rules applied to the tenant.
// returns a []DriveProtectionRuleable when successful
func (m *BackupRestoreRoot) GetDriveInclusionRules()([]DriveProtectionRuleable) {
    val, err := m.GetBackingStore().Get("driveInclusionRules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveProtectionRuleable)
    }
    return nil
}
// GetDriveProtectionUnits gets the driveProtectionUnits property value. The list of drive protection units in the tenant.
// returns a []DriveProtectionUnitable when successful
func (m *BackupRestoreRoot) GetDriveProtectionUnits()([]DriveProtectionUnitable) {
    val, err := m.GetBackingStore().Get("driveProtectionUnits")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveProtectionUnitable)
    }
    return nil
}
// GetExchangeProtectionPolicies gets the exchangeProtectionPolicies property value. The list of Exchange protection policies in the tenant.
// returns a []ExchangeProtectionPolicyable when successful
func (m *BackupRestoreRoot) GetExchangeProtectionPolicies()([]ExchangeProtectionPolicyable) {
    val, err := m.GetBackingStore().Get("exchangeProtectionPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ExchangeProtectionPolicyable)
    }
    return nil
}
// GetExchangeRestoreSessions gets the exchangeRestoreSessions property value. The list of Exchange restore sessions available in the tenant.
// returns a []ExchangeRestoreSessionable when successful
func (m *BackupRestoreRoot) GetExchangeRestoreSessions()([]ExchangeRestoreSessionable) {
    val, err := m.GetBackingStore().Get("exchangeRestoreSessions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ExchangeRestoreSessionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *BackupRestoreRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["driveInclusionRules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDriveProtectionRuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DriveProtectionRuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DriveProtectionRuleable)
                }
            }
            m.SetDriveInclusionRules(res)
        }
        return nil
    }
    res["driveProtectionUnits"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDriveProtectionUnitFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DriveProtectionUnitable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DriveProtectionUnitable)
                }
            }
            m.SetDriveProtectionUnits(res)
        }
        return nil
    }
    res["exchangeProtectionPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateExchangeProtectionPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ExchangeProtectionPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ExchangeProtectionPolicyable)
                }
            }
            m.SetExchangeProtectionPolicies(res)
        }
        return nil
    }
    res["exchangeRestoreSessions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateExchangeRestoreSessionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ExchangeRestoreSessionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ExchangeRestoreSessionable)
                }
            }
            m.SetExchangeRestoreSessions(res)
        }
        return nil
    }
    res["mailboxInclusionRules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMailboxProtectionRuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MailboxProtectionRuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MailboxProtectionRuleable)
                }
            }
            m.SetMailboxInclusionRules(res)
        }
        return nil
    }
    res["mailboxProtectionUnits"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMailboxProtectionUnitFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MailboxProtectionUnitable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MailboxProtectionUnitable)
                }
            }
            m.SetMailboxProtectionUnits(res)
        }
        return nil
    }
    res["oneDriveForBusinessProtectionPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOneDriveForBusinessProtectionPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OneDriveForBusinessProtectionPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OneDriveForBusinessProtectionPolicyable)
                }
            }
            m.SetOneDriveForBusinessProtectionPolicies(res)
        }
        return nil
    }
    res["oneDriveForBusinessRestoreSessions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateOneDriveForBusinessRestoreSessionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]OneDriveForBusinessRestoreSessionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(OneDriveForBusinessRestoreSessionable)
                }
            }
            m.SetOneDriveForBusinessRestoreSessions(res)
        }
        return nil
    }
    res["protectionPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateProtectionPolicyBaseFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ProtectionPolicyBaseable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ProtectionPolicyBaseable)
                }
            }
            m.SetProtectionPolicies(res)
        }
        return nil
    }
    res["protectionUnits"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateProtectionUnitBaseFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ProtectionUnitBaseable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ProtectionUnitBaseable)
                }
            }
            m.SetProtectionUnits(res)
        }
        return nil
    }
    res["restorePoints"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRestorePointFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RestorePointable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RestorePointable)
                }
            }
            m.SetRestorePoints(res)
        }
        return nil
    }
    res["restoreSessions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRestoreSessionBaseFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RestoreSessionBaseable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RestoreSessionBaseable)
                }
            }
            m.SetRestoreSessions(res)
        }
        return nil
    }
    res["serviceApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateServiceAppFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ServiceAppable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ServiceAppable)
                }
            }
            m.SetServiceApps(res)
        }
        return nil
    }
    res["serviceStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateServiceStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServiceStatus(val.(ServiceStatusable))
        }
        return nil
    }
    res["sharePointProtectionPolicies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSharePointProtectionPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SharePointProtectionPolicyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SharePointProtectionPolicyable)
                }
            }
            m.SetSharePointProtectionPolicies(res)
        }
        return nil
    }
    res["sharePointRestoreSessions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSharePointRestoreSessionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SharePointRestoreSessionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SharePointRestoreSessionable)
                }
            }
            m.SetSharePointRestoreSessions(res)
        }
        return nil
    }
    res["siteInclusionRules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSiteProtectionRuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SiteProtectionRuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SiteProtectionRuleable)
                }
            }
            m.SetSiteInclusionRules(res)
        }
        return nil
    }
    res["siteProtectionUnits"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSiteProtectionUnitFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SiteProtectionUnitable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SiteProtectionUnitable)
                }
            }
            m.SetSiteProtectionUnits(res)
        }
        return nil
    }
    return res
}
// GetMailboxInclusionRules gets the mailboxInclusionRules property value. The list of mailbox inclusion rules applied to the tenant.
// returns a []MailboxProtectionRuleable when successful
func (m *BackupRestoreRoot) GetMailboxInclusionRules()([]MailboxProtectionRuleable) {
    val, err := m.GetBackingStore().Get("mailboxInclusionRules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MailboxProtectionRuleable)
    }
    return nil
}
// GetMailboxProtectionUnits gets the mailboxProtectionUnits property value. The list of mailbox protection units in the tenant.
// returns a []MailboxProtectionUnitable when successful
func (m *BackupRestoreRoot) GetMailboxProtectionUnits()([]MailboxProtectionUnitable) {
    val, err := m.GetBackingStore().Get("mailboxProtectionUnits")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MailboxProtectionUnitable)
    }
    return nil
}
// GetOneDriveForBusinessProtectionPolicies gets the oneDriveForBusinessProtectionPolicies property value. The list of OneDrive for Business protection policies in the tenant.
// returns a []OneDriveForBusinessProtectionPolicyable when successful
func (m *BackupRestoreRoot) GetOneDriveForBusinessProtectionPolicies()([]OneDriveForBusinessProtectionPolicyable) {
    val, err := m.GetBackingStore().Get("oneDriveForBusinessProtectionPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OneDriveForBusinessProtectionPolicyable)
    }
    return nil
}
// GetOneDriveForBusinessRestoreSessions gets the oneDriveForBusinessRestoreSessions property value. The list of OneDrive for Business restore sessions available in the tenant.
// returns a []OneDriveForBusinessRestoreSessionable when successful
func (m *BackupRestoreRoot) GetOneDriveForBusinessRestoreSessions()([]OneDriveForBusinessRestoreSessionable) {
    val, err := m.GetBackingStore().Get("oneDriveForBusinessRestoreSessions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]OneDriveForBusinessRestoreSessionable)
    }
    return nil
}
// GetProtectionPolicies gets the protectionPolicies property value. List of protection policies in the tenant.
// returns a []ProtectionPolicyBaseable when successful
func (m *BackupRestoreRoot) GetProtectionPolicies()([]ProtectionPolicyBaseable) {
    val, err := m.GetBackingStore().Get("protectionPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ProtectionPolicyBaseable)
    }
    return nil
}
// GetProtectionUnits gets the protectionUnits property value. List of protection units in the tenant.
// returns a []ProtectionUnitBaseable when successful
func (m *BackupRestoreRoot) GetProtectionUnits()([]ProtectionUnitBaseable) {
    val, err := m.GetBackingStore().Get("protectionUnits")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ProtectionUnitBaseable)
    }
    return nil
}
// GetRestorePoints gets the restorePoints property value. List of restore points in the tenant.
// returns a []RestorePointable when successful
func (m *BackupRestoreRoot) GetRestorePoints()([]RestorePointable) {
    val, err := m.GetBackingStore().Get("restorePoints")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RestorePointable)
    }
    return nil
}
// GetRestoreSessions gets the restoreSessions property value. List of restore sessions in the tenant.
// returns a []RestoreSessionBaseable when successful
func (m *BackupRestoreRoot) GetRestoreSessions()([]RestoreSessionBaseable) {
    val, err := m.GetBackingStore().Get("restoreSessions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RestoreSessionBaseable)
    }
    return nil
}
// GetServiceApps gets the serviceApps property value. List of Backup Storage apps in the tenant.
// returns a []ServiceAppable when successful
func (m *BackupRestoreRoot) GetServiceApps()([]ServiceAppable) {
    val, err := m.GetBackingStore().Get("serviceApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ServiceAppable)
    }
    return nil
}
// GetServiceStatus gets the serviceStatus property value. Represents the tenant-level status of the Backup Storage service.
// returns a ServiceStatusable when successful
func (m *BackupRestoreRoot) GetServiceStatus()(ServiceStatusable) {
    val, err := m.GetBackingStore().Get("serviceStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ServiceStatusable)
    }
    return nil
}
// GetSharePointProtectionPolicies gets the sharePointProtectionPolicies property value. The list of SharePoint protection policies in the tenant.
// returns a []SharePointProtectionPolicyable when successful
func (m *BackupRestoreRoot) GetSharePointProtectionPolicies()([]SharePointProtectionPolicyable) {
    val, err := m.GetBackingStore().Get("sharePointProtectionPolicies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SharePointProtectionPolicyable)
    }
    return nil
}
// GetSharePointRestoreSessions gets the sharePointRestoreSessions property value. The list of SharePoint restore sessions available in the tenant.
// returns a []SharePointRestoreSessionable when successful
func (m *BackupRestoreRoot) GetSharePointRestoreSessions()([]SharePointRestoreSessionable) {
    val, err := m.GetBackingStore().Get("sharePointRestoreSessions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SharePointRestoreSessionable)
    }
    return nil
}
// GetSiteInclusionRules gets the siteInclusionRules property value. The list of site inclusion rules applied to the tenant.
// returns a []SiteProtectionRuleable when successful
func (m *BackupRestoreRoot) GetSiteInclusionRules()([]SiteProtectionRuleable) {
    val, err := m.GetBackingStore().Get("siteInclusionRules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SiteProtectionRuleable)
    }
    return nil
}
// GetSiteProtectionUnits gets the siteProtectionUnits property value. The list of site protection units in the tenant.
// returns a []SiteProtectionUnitable when successful
func (m *BackupRestoreRoot) GetSiteProtectionUnits()([]SiteProtectionUnitable) {
    val, err := m.GetBackingStore().Get("siteProtectionUnits")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SiteProtectionUnitable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *BackupRestoreRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetDriveInclusionRules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDriveInclusionRules()))
        for i, v := range m.GetDriveInclusionRules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("driveInclusionRules", cast)
        if err != nil {
            return err
        }
    }
    if m.GetDriveProtectionUnits() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDriveProtectionUnits()))
        for i, v := range m.GetDriveProtectionUnits() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("driveProtectionUnits", cast)
        if err != nil {
            return err
        }
    }
    if m.GetExchangeProtectionPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExchangeProtectionPolicies()))
        for i, v := range m.GetExchangeProtectionPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("exchangeProtectionPolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetExchangeRestoreSessions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExchangeRestoreSessions()))
        for i, v := range m.GetExchangeRestoreSessions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("exchangeRestoreSessions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMailboxInclusionRules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMailboxInclusionRules()))
        for i, v := range m.GetMailboxInclusionRules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("mailboxInclusionRules", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMailboxProtectionUnits() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMailboxProtectionUnits()))
        for i, v := range m.GetMailboxProtectionUnits() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("mailboxProtectionUnits", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOneDriveForBusinessProtectionPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOneDriveForBusinessProtectionPolicies()))
        for i, v := range m.GetOneDriveForBusinessProtectionPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("oneDriveForBusinessProtectionPolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOneDriveForBusinessRestoreSessions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOneDriveForBusinessRestoreSessions()))
        for i, v := range m.GetOneDriveForBusinessRestoreSessions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("oneDriveForBusinessRestoreSessions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetProtectionPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetProtectionPolicies()))
        for i, v := range m.GetProtectionPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("protectionPolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetProtectionUnits() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetProtectionUnits()))
        for i, v := range m.GetProtectionUnits() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("protectionUnits", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRestorePoints() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRestorePoints()))
        for i, v := range m.GetRestorePoints() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("restorePoints", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRestoreSessions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRestoreSessions()))
        for i, v := range m.GetRestoreSessions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("restoreSessions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetServiceApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetServiceApps()))
        for i, v := range m.GetServiceApps() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("serviceApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("serviceStatus", m.GetServiceStatus())
        if err != nil {
            return err
        }
    }
    if m.GetSharePointProtectionPolicies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSharePointProtectionPolicies()))
        for i, v := range m.GetSharePointProtectionPolicies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("sharePointProtectionPolicies", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSharePointRestoreSessions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSharePointRestoreSessions()))
        for i, v := range m.GetSharePointRestoreSessions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("sharePointRestoreSessions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSiteInclusionRules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSiteInclusionRules()))
        for i, v := range m.GetSiteInclusionRules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("siteInclusionRules", cast)
        if err != nil {
            return err
        }
    }
    if m.GetSiteProtectionUnits() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSiteProtectionUnits()))
        for i, v := range m.GetSiteProtectionUnits() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("siteProtectionUnits", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDriveInclusionRules sets the driveInclusionRules property value. The list of drive inclusion rules applied to the tenant.
func (m *BackupRestoreRoot) SetDriveInclusionRules(value []DriveProtectionRuleable)() {
    err := m.GetBackingStore().Set("driveInclusionRules", value)
    if err != nil {
        panic(err)
    }
}
// SetDriveProtectionUnits sets the driveProtectionUnits property value. The list of drive protection units in the tenant.
func (m *BackupRestoreRoot) SetDriveProtectionUnits(value []DriveProtectionUnitable)() {
    err := m.GetBackingStore().Set("driveProtectionUnits", value)
    if err != nil {
        panic(err)
    }
}
// SetExchangeProtectionPolicies sets the exchangeProtectionPolicies property value. The list of Exchange protection policies in the tenant.
func (m *BackupRestoreRoot) SetExchangeProtectionPolicies(value []ExchangeProtectionPolicyable)() {
    err := m.GetBackingStore().Set("exchangeProtectionPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetExchangeRestoreSessions sets the exchangeRestoreSessions property value. The list of Exchange restore sessions available in the tenant.
func (m *BackupRestoreRoot) SetExchangeRestoreSessions(value []ExchangeRestoreSessionable)() {
    err := m.GetBackingStore().Set("exchangeRestoreSessions", value)
    if err != nil {
        panic(err)
    }
}
// SetMailboxInclusionRules sets the mailboxInclusionRules property value. The list of mailbox inclusion rules applied to the tenant.
func (m *BackupRestoreRoot) SetMailboxInclusionRules(value []MailboxProtectionRuleable)() {
    err := m.GetBackingStore().Set("mailboxInclusionRules", value)
    if err != nil {
        panic(err)
    }
}
// SetMailboxProtectionUnits sets the mailboxProtectionUnits property value. The list of mailbox protection units in the tenant.
func (m *BackupRestoreRoot) SetMailboxProtectionUnits(value []MailboxProtectionUnitable)() {
    err := m.GetBackingStore().Set("mailboxProtectionUnits", value)
    if err != nil {
        panic(err)
    }
}
// SetOneDriveForBusinessProtectionPolicies sets the oneDriveForBusinessProtectionPolicies property value. The list of OneDrive for Business protection policies in the tenant.
func (m *BackupRestoreRoot) SetOneDriveForBusinessProtectionPolicies(value []OneDriveForBusinessProtectionPolicyable)() {
    err := m.GetBackingStore().Set("oneDriveForBusinessProtectionPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetOneDriveForBusinessRestoreSessions sets the oneDriveForBusinessRestoreSessions property value. The list of OneDrive for Business restore sessions available in the tenant.
func (m *BackupRestoreRoot) SetOneDriveForBusinessRestoreSessions(value []OneDriveForBusinessRestoreSessionable)() {
    err := m.GetBackingStore().Set("oneDriveForBusinessRestoreSessions", value)
    if err != nil {
        panic(err)
    }
}
// SetProtectionPolicies sets the protectionPolicies property value. List of protection policies in the tenant.
func (m *BackupRestoreRoot) SetProtectionPolicies(value []ProtectionPolicyBaseable)() {
    err := m.GetBackingStore().Set("protectionPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetProtectionUnits sets the protectionUnits property value. List of protection units in the tenant.
func (m *BackupRestoreRoot) SetProtectionUnits(value []ProtectionUnitBaseable)() {
    err := m.GetBackingStore().Set("protectionUnits", value)
    if err != nil {
        panic(err)
    }
}
// SetRestorePoints sets the restorePoints property value. List of restore points in the tenant.
func (m *BackupRestoreRoot) SetRestorePoints(value []RestorePointable)() {
    err := m.GetBackingStore().Set("restorePoints", value)
    if err != nil {
        panic(err)
    }
}
// SetRestoreSessions sets the restoreSessions property value. List of restore sessions in the tenant.
func (m *BackupRestoreRoot) SetRestoreSessions(value []RestoreSessionBaseable)() {
    err := m.GetBackingStore().Set("restoreSessions", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceApps sets the serviceApps property value. List of Backup Storage apps in the tenant.
func (m *BackupRestoreRoot) SetServiceApps(value []ServiceAppable)() {
    err := m.GetBackingStore().Set("serviceApps", value)
    if err != nil {
        panic(err)
    }
}
// SetServiceStatus sets the serviceStatus property value. Represents the tenant-level status of the Backup Storage service.
func (m *BackupRestoreRoot) SetServiceStatus(value ServiceStatusable)() {
    err := m.GetBackingStore().Set("serviceStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetSharePointProtectionPolicies sets the sharePointProtectionPolicies property value. The list of SharePoint protection policies in the tenant.
func (m *BackupRestoreRoot) SetSharePointProtectionPolicies(value []SharePointProtectionPolicyable)() {
    err := m.GetBackingStore().Set("sharePointProtectionPolicies", value)
    if err != nil {
        panic(err)
    }
}
// SetSharePointRestoreSessions sets the sharePointRestoreSessions property value. The list of SharePoint restore sessions available in the tenant.
func (m *BackupRestoreRoot) SetSharePointRestoreSessions(value []SharePointRestoreSessionable)() {
    err := m.GetBackingStore().Set("sharePointRestoreSessions", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteInclusionRules sets the siteInclusionRules property value. The list of site inclusion rules applied to the tenant.
func (m *BackupRestoreRoot) SetSiteInclusionRules(value []SiteProtectionRuleable)() {
    err := m.GetBackingStore().Set("siteInclusionRules", value)
    if err != nil {
        panic(err)
    }
}
// SetSiteProtectionUnits sets the siteProtectionUnits property value. The list of site protection units in the tenant.
func (m *BackupRestoreRoot) SetSiteProtectionUnits(value []SiteProtectionUnitable)() {
    err := m.GetBackingStore().Set("siteProtectionUnits", value)
    if err != nil {
        panic(err)
    }
}
type BackupRestoreRootable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDriveInclusionRules()([]DriveProtectionRuleable)
    GetDriveProtectionUnits()([]DriveProtectionUnitable)
    GetExchangeProtectionPolicies()([]ExchangeProtectionPolicyable)
    GetExchangeRestoreSessions()([]ExchangeRestoreSessionable)
    GetMailboxInclusionRules()([]MailboxProtectionRuleable)
    GetMailboxProtectionUnits()([]MailboxProtectionUnitable)
    GetOneDriveForBusinessProtectionPolicies()([]OneDriveForBusinessProtectionPolicyable)
    GetOneDriveForBusinessRestoreSessions()([]OneDriveForBusinessRestoreSessionable)
    GetProtectionPolicies()([]ProtectionPolicyBaseable)
    GetProtectionUnits()([]ProtectionUnitBaseable)
    GetRestorePoints()([]RestorePointable)
    GetRestoreSessions()([]RestoreSessionBaseable)
    GetServiceApps()([]ServiceAppable)
    GetServiceStatus()(ServiceStatusable)
    GetSharePointProtectionPolicies()([]SharePointProtectionPolicyable)
    GetSharePointRestoreSessions()([]SharePointRestoreSessionable)
    GetSiteInclusionRules()([]SiteProtectionRuleable)
    GetSiteProtectionUnits()([]SiteProtectionUnitable)
    SetDriveInclusionRules(value []DriveProtectionRuleable)()
    SetDriveProtectionUnits(value []DriveProtectionUnitable)()
    SetExchangeProtectionPolicies(value []ExchangeProtectionPolicyable)()
    SetExchangeRestoreSessions(value []ExchangeRestoreSessionable)()
    SetMailboxInclusionRules(value []MailboxProtectionRuleable)()
    SetMailboxProtectionUnits(value []MailboxProtectionUnitable)()
    SetOneDriveForBusinessProtectionPolicies(value []OneDriveForBusinessProtectionPolicyable)()
    SetOneDriveForBusinessRestoreSessions(value []OneDriveForBusinessRestoreSessionable)()
    SetProtectionPolicies(value []ProtectionPolicyBaseable)()
    SetProtectionUnits(value []ProtectionUnitBaseable)()
    SetRestorePoints(value []RestorePointable)()
    SetRestoreSessions(value []RestoreSessionBaseable)()
    SetServiceApps(value []ServiceAppable)()
    SetServiceStatus(value ServiceStatusable)()
    SetSharePointProtectionPolicies(value []SharePointProtectionPolicyable)()
    SetSharePointRestoreSessions(value []SharePointRestoreSessionable)()
    SetSiteInclusionRules(value []SiteProtectionRuleable)()
    SetSiteProtectionUnits(value []SiteProtectionUnitable)()
}
