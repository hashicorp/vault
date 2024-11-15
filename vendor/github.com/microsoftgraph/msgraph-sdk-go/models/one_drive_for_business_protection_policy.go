package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type OneDriveForBusinessProtectionPolicy struct {
    ProtectionPolicyBase
}
// NewOneDriveForBusinessProtectionPolicy instantiates a new OneDriveForBusinessProtectionPolicy and sets the default values.
func NewOneDriveForBusinessProtectionPolicy()(*OneDriveForBusinessProtectionPolicy) {
    m := &OneDriveForBusinessProtectionPolicy{
        ProtectionPolicyBase: *NewProtectionPolicyBase(),
    }
    odataTypeValue := "#microsoft.graph.oneDriveForBusinessProtectionPolicy"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateOneDriveForBusinessProtectionPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOneDriveForBusinessProtectionPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOneDriveForBusinessProtectionPolicy(), nil
}
// GetDriveInclusionRules gets the driveInclusionRules property value. Contains the details of the Onedrive for Business protection rule.
// returns a []DriveProtectionRuleable when successful
func (m *OneDriveForBusinessProtectionPolicy) GetDriveInclusionRules()([]DriveProtectionRuleable) {
    val, err := m.GetBackingStore().Get("driveInclusionRules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveProtectionRuleable)
    }
    return nil
}
// GetDriveProtectionUnits gets the driveProtectionUnits property value. Contains the protection units associated with a  OneDrive for Business protection policy.
// returns a []DriveProtectionUnitable when successful
func (m *OneDriveForBusinessProtectionPolicy) GetDriveProtectionUnits()([]DriveProtectionUnitable) {
    val, err := m.GetBackingStore().Get("driveProtectionUnits")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DriveProtectionUnitable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OneDriveForBusinessProtectionPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ProtectionPolicyBase.GetFieldDeserializers()
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
    return res
}
// Serialize serializes information the current object
func (m *OneDriveForBusinessProtectionPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ProtectionPolicyBase.Serialize(writer)
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
    return nil
}
// SetDriveInclusionRules sets the driveInclusionRules property value. Contains the details of the Onedrive for Business protection rule.
func (m *OneDriveForBusinessProtectionPolicy) SetDriveInclusionRules(value []DriveProtectionRuleable)() {
    err := m.GetBackingStore().Set("driveInclusionRules", value)
    if err != nil {
        panic(err)
    }
}
// SetDriveProtectionUnits sets the driveProtectionUnits property value. Contains the protection units associated with a  OneDrive for Business protection policy.
func (m *OneDriveForBusinessProtectionPolicy) SetDriveProtectionUnits(value []DriveProtectionUnitable)() {
    err := m.GetBackingStore().Set("driveProtectionUnits", value)
    if err != nil {
        panic(err)
    }
}
type OneDriveForBusinessProtectionPolicyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ProtectionPolicyBaseable
    GetDriveInclusionRules()([]DriveProtectionRuleable)
    GetDriveProtectionUnits()([]DriveProtectionUnitable)
    SetDriveInclusionRules(value []DriveProtectionRuleable)()
    SetDriveProtectionUnits(value []DriveProtectionUnitable)()
}
