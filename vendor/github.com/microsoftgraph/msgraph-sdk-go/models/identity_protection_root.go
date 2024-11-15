package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type IdentityProtectionRoot struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewIdentityProtectionRoot instantiates a new IdentityProtectionRoot and sets the default values.
func NewIdentityProtectionRoot()(*IdentityProtectionRoot) {
    m := &IdentityProtectionRoot{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateIdentityProtectionRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIdentityProtectionRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIdentityProtectionRoot(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *IdentityProtectionRoot) GetAdditionalData()(map[string]any) {
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
func (m *IdentityProtectionRoot) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IdentityProtectionRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
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
    res["riskDetections"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRiskDetectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RiskDetectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RiskDetectionable)
                }
            }
            m.SetRiskDetections(res)
        }
        return nil
    }
    res["riskyServicePrincipals"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRiskyServicePrincipalFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RiskyServicePrincipalable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RiskyServicePrincipalable)
                }
            }
            m.SetRiskyServicePrincipals(res)
        }
        return nil
    }
    res["riskyUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRiskyUserFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]RiskyUserable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(RiskyUserable)
                }
            }
            m.SetRiskyUsers(res)
        }
        return nil
    }
    res["servicePrincipalRiskDetections"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateServicePrincipalRiskDetectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ServicePrincipalRiskDetectionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ServicePrincipalRiskDetectionable)
                }
            }
            m.SetServicePrincipalRiskDetections(res)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *IdentityProtectionRoot) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRiskDetections gets the riskDetections property value. Risk detection in Microsoft Entra ID Protection and the associated information about the detection.
// returns a []RiskDetectionable when successful
func (m *IdentityProtectionRoot) GetRiskDetections()([]RiskDetectionable) {
    val, err := m.GetBackingStore().Get("riskDetections")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RiskDetectionable)
    }
    return nil
}
// GetRiskyServicePrincipals gets the riskyServicePrincipals property value. Microsoft Entra service principals that are at risk.
// returns a []RiskyServicePrincipalable when successful
func (m *IdentityProtectionRoot) GetRiskyServicePrincipals()([]RiskyServicePrincipalable) {
    val, err := m.GetBackingStore().Get("riskyServicePrincipals")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RiskyServicePrincipalable)
    }
    return nil
}
// GetRiskyUsers gets the riskyUsers property value. Users that are flagged as at-risk by Microsoft Entra ID Protection.
// returns a []RiskyUserable when successful
func (m *IdentityProtectionRoot) GetRiskyUsers()([]RiskyUserable) {
    val, err := m.GetBackingStore().Get("riskyUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]RiskyUserable)
    }
    return nil
}
// GetServicePrincipalRiskDetections gets the servicePrincipalRiskDetections property value. Represents information about detected at-risk service principals in a Microsoft Entra tenant.
// returns a []ServicePrincipalRiskDetectionable when successful
func (m *IdentityProtectionRoot) GetServicePrincipalRiskDetections()([]ServicePrincipalRiskDetectionable) {
    val, err := m.GetBackingStore().Get("servicePrincipalRiskDetections")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ServicePrincipalRiskDetectionable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IdentityProtectionRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    if m.GetRiskDetections() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRiskDetections()))
        for i, v := range m.GetRiskDetections() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("riskDetections", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRiskyServicePrincipals() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRiskyServicePrincipals()))
        for i, v := range m.GetRiskyServicePrincipals() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("riskyServicePrincipals", cast)
        if err != nil {
            return err
        }
    }
    if m.GetRiskyUsers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRiskyUsers()))
        for i, v := range m.GetRiskyUsers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("riskyUsers", cast)
        if err != nil {
            return err
        }
    }
    if m.GetServicePrincipalRiskDetections() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetServicePrincipalRiskDetections()))
        for i, v := range m.GetServicePrincipalRiskDetections() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("servicePrincipalRiskDetections", cast)
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
func (m *IdentityProtectionRoot) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *IdentityProtectionRoot) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *IdentityProtectionRoot) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskDetections sets the riskDetections property value. Risk detection in Microsoft Entra ID Protection and the associated information about the detection.
func (m *IdentityProtectionRoot) SetRiskDetections(value []RiskDetectionable)() {
    err := m.GetBackingStore().Set("riskDetections", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskyServicePrincipals sets the riskyServicePrincipals property value. Microsoft Entra service principals that are at risk.
func (m *IdentityProtectionRoot) SetRiskyServicePrincipals(value []RiskyServicePrincipalable)() {
    err := m.GetBackingStore().Set("riskyServicePrincipals", value)
    if err != nil {
        panic(err)
    }
}
// SetRiskyUsers sets the riskyUsers property value. Users that are flagged as at-risk by Microsoft Entra ID Protection.
func (m *IdentityProtectionRoot) SetRiskyUsers(value []RiskyUserable)() {
    err := m.GetBackingStore().Set("riskyUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePrincipalRiskDetections sets the servicePrincipalRiskDetections property value. Represents information about detected at-risk service principals in a Microsoft Entra tenant.
func (m *IdentityProtectionRoot) SetServicePrincipalRiskDetections(value []ServicePrincipalRiskDetectionable)() {
    err := m.GetBackingStore().Set("servicePrincipalRiskDetections", value)
    if err != nil {
        panic(err)
    }
}
type IdentityProtectionRootable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    GetRiskDetections()([]RiskDetectionable)
    GetRiskyServicePrincipals()([]RiskyServicePrincipalable)
    GetRiskyUsers()([]RiskyUserable)
    GetServicePrincipalRiskDetections()([]ServicePrincipalRiskDetectionable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
    SetRiskDetections(value []RiskDetectionable)()
    SetRiskyServicePrincipals(value []RiskyServicePrincipalable)()
    SetRiskyUsers(value []RiskyUserable)()
    SetServicePrincipalRiskDetections(value []ServicePrincipalRiskDetectionable)()
}
