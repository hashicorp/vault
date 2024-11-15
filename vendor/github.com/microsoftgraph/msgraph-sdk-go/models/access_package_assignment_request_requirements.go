package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AccessPackageAssignmentRequestRequirements struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAccessPackageAssignmentRequestRequirements instantiates a new AccessPackageAssignmentRequestRequirements and sets the default values.
func NewAccessPackageAssignmentRequestRequirements()(*AccessPackageAssignmentRequestRequirements) {
    m := &AccessPackageAssignmentRequestRequirements{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAccessPackageAssignmentRequestRequirementsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageAssignmentRequestRequirementsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageAssignmentRequestRequirements(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AccessPackageAssignmentRequestRequirements) GetAdditionalData()(map[string]any) {
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
// GetAllowCustomAssignmentSchedule gets the allowCustomAssignmentSchedule property value. Indicates whether the requestor is allowed to set a custom schedule.
// returns a *bool when successful
func (m *AccessPackageAssignmentRequestRequirements) GetAllowCustomAssignmentSchedule()(*bool) {
    val, err := m.GetBackingStore().Get("allowCustomAssignmentSchedule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *AccessPackageAssignmentRequestRequirements) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessPackageAssignmentRequestRequirements) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowCustomAssignmentSchedule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowCustomAssignmentSchedule(val)
        }
        return nil
    }
    res["isApprovalRequiredForAdd"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsApprovalRequiredForAdd(val)
        }
        return nil
    }
    res["isApprovalRequiredForUpdate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsApprovalRequiredForUpdate(val)
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
    res["policyDescription"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicyDescription(val)
        }
        return nil
    }
    res["policyDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicyDisplayName(val)
        }
        return nil
    }
    res["policyId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicyId(val)
        }
        return nil
    }
    res["questions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAccessPackageQuestionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]AccessPackageQuestionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(AccessPackageQuestionable)
                }
            }
            m.SetQuestions(res)
        }
        return nil
    }
    res["schedule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEntitlementManagementScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSchedule(val.(EntitlementManagementScheduleable))
        }
        return nil
    }
    return res
}
// GetIsApprovalRequiredForAdd gets the isApprovalRequiredForAdd property value. Indicates whether a request to add must be approved by an approver.
// returns a *bool when successful
func (m *AccessPackageAssignmentRequestRequirements) GetIsApprovalRequiredForAdd()(*bool) {
    val, err := m.GetBackingStore().Get("isApprovalRequiredForAdd")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsApprovalRequiredForUpdate gets the isApprovalRequiredForUpdate property value. Indicates whether a request to update must be approved by an approver.
// returns a *bool when successful
func (m *AccessPackageAssignmentRequestRequirements) GetIsApprovalRequiredForUpdate()(*bool) {
    val, err := m.GetBackingStore().Get("isApprovalRequiredForUpdate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *AccessPackageAssignmentRequestRequirements) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPolicyDescription gets the policyDescription property value. The description of the policy that the user is trying to request access using.
// returns a *string when successful
func (m *AccessPackageAssignmentRequestRequirements) GetPolicyDescription()(*string) {
    val, err := m.GetBackingStore().Get("policyDescription")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPolicyDisplayName gets the policyDisplayName property value. The display name of the policy that the user is trying to request access using.
// returns a *string when successful
func (m *AccessPackageAssignmentRequestRequirements) GetPolicyDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("policyDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPolicyId gets the policyId property value. The identifier of the policy that these requirements are associated with. This identifier can be used when creating a new assignment request.
// returns a *string when successful
func (m *AccessPackageAssignmentRequestRequirements) GetPolicyId()(*string) {
    val, err := m.GetBackingStore().Get("policyId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQuestions gets the questions property value. The questions property
// returns a []AccessPackageQuestionable when successful
func (m *AccessPackageAssignmentRequestRequirements) GetQuestions()([]AccessPackageQuestionable) {
    val, err := m.GetBackingStore().Get("questions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]AccessPackageQuestionable)
    }
    return nil
}
// GetSchedule gets the schedule property value. Schedule restrictions enforced, if any.
// returns a EntitlementManagementScheduleable when successful
func (m *AccessPackageAssignmentRequestRequirements) GetSchedule()(EntitlementManagementScheduleable) {
    val, err := m.GetBackingStore().Get("schedule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EntitlementManagementScheduleable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageAssignmentRequestRequirements) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("allowCustomAssignmentSchedule", m.GetAllowCustomAssignmentSchedule())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isApprovalRequiredForAdd", m.GetIsApprovalRequiredForAdd())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isApprovalRequiredForUpdate", m.GetIsApprovalRequiredForUpdate())
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
        err := writer.WriteStringValue("policyDescription", m.GetPolicyDescription())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("policyDisplayName", m.GetPolicyDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("policyId", m.GetPolicyId())
        if err != nil {
            return err
        }
    }
    if m.GetQuestions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetQuestions()))
        for i, v := range m.GetQuestions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("questions", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("schedule", m.GetSchedule())
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
func (m *AccessPackageAssignmentRequestRequirements) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowCustomAssignmentSchedule sets the allowCustomAssignmentSchedule property value. Indicates whether the requestor is allowed to set a custom schedule.
func (m *AccessPackageAssignmentRequestRequirements) SetAllowCustomAssignmentSchedule(value *bool)() {
    err := m.GetBackingStore().Set("allowCustomAssignmentSchedule", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AccessPackageAssignmentRequestRequirements) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetIsApprovalRequiredForAdd sets the isApprovalRequiredForAdd property value. Indicates whether a request to add must be approved by an approver.
func (m *AccessPackageAssignmentRequestRequirements) SetIsApprovalRequiredForAdd(value *bool)() {
    err := m.GetBackingStore().Set("isApprovalRequiredForAdd", value)
    if err != nil {
        panic(err)
    }
}
// SetIsApprovalRequiredForUpdate sets the isApprovalRequiredForUpdate property value. Indicates whether a request to update must be approved by an approver.
func (m *AccessPackageAssignmentRequestRequirements) SetIsApprovalRequiredForUpdate(value *bool)() {
    err := m.GetBackingStore().Set("isApprovalRequiredForUpdate", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AccessPackageAssignmentRequestRequirements) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicyDescription sets the policyDescription property value. The description of the policy that the user is trying to request access using.
func (m *AccessPackageAssignmentRequestRequirements) SetPolicyDescription(value *string)() {
    err := m.GetBackingStore().Set("policyDescription", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicyDisplayName sets the policyDisplayName property value. The display name of the policy that the user is trying to request access using.
func (m *AccessPackageAssignmentRequestRequirements) SetPolicyDisplayName(value *string)() {
    err := m.GetBackingStore().Set("policyDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicyId sets the policyId property value. The identifier of the policy that these requirements are associated with. This identifier can be used when creating a new assignment request.
func (m *AccessPackageAssignmentRequestRequirements) SetPolicyId(value *string)() {
    err := m.GetBackingStore().Set("policyId", value)
    if err != nil {
        panic(err)
    }
}
// SetQuestions sets the questions property value. The questions property
func (m *AccessPackageAssignmentRequestRequirements) SetQuestions(value []AccessPackageQuestionable)() {
    err := m.GetBackingStore().Set("questions", value)
    if err != nil {
        panic(err)
    }
}
// SetSchedule sets the schedule property value. Schedule restrictions enforced, if any.
func (m *AccessPackageAssignmentRequestRequirements) SetSchedule(value EntitlementManagementScheduleable)() {
    err := m.GetBackingStore().Set("schedule", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageAssignmentRequestRequirementsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowCustomAssignmentSchedule()(*bool)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetIsApprovalRequiredForAdd()(*bool)
    GetIsApprovalRequiredForUpdate()(*bool)
    GetOdataType()(*string)
    GetPolicyDescription()(*string)
    GetPolicyDisplayName()(*string)
    GetPolicyId()(*string)
    GetQuestions()([]AccessPackageQuestionable)
    GetSchedule()(EntitlementManagementScheduleable)
    SetAllowCustomAssignmentSchedule(value *bool)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetIsApprovalRequiredForAdd(value *bool)()
    SetIsApprovalRequiredForUpdate(value *bool)()
    SetOdataType(value *string)()
    SetPolicyDescription(value *string)()
    SetPolicyDisplayName(value *string)()
    SetPolicyId(value *string)()
    SetQuestions(value []AccessPackageQuestionable)()
    SetSchedule(value EntitlementManagementScheduleable)()
}
