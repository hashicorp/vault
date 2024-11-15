package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CloudPcManagementGroupAssignmentTarget struct {
    CloudPcManagementAssignmentTarget
}
// NewCloudPcManagementGroupAssignmentTarget instantiates a new CloudPcManagementGroupAssignmentTarget and sets the default values.
func NewCloudPcManagementGroupAssignmentTarget()(*CloudPcManagementGroupAssignmentTarget) {
    m := &CloudPcManagementGroupAssignmentTarget{
        CloudPcManagementAssignmentTarget: *NewCloudPcManagementAssignmentTarget(),
    }
    odataTypeValue := "#microsoft.graph.cloudPcManagementGroupAssignmentTarget"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCloudPcManagementGroupAssignmentTargetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCloudPcManagementGroupAssignmentTargetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCloudPcManagementGroupAssignmentTarget(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CloudPcManagementGroupAssignmentTarget) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.CloudPcManagementAssignmentTarget.GetFieldDeserializers()
    res["groupId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroupId(val)
        }
        return nil
    }
    res["servicePlanId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetServicePlanId(val)
        }
        return nil
    }
    return res
}
// GetGroupId gets the groupId property value. The ID of the target group for the assignment.
// returns a *string when successful
func (m *CloudPcManagementGroupAssignmentTarget) GetGroupId()(*string) {
    val, err := m.GetBackingStore().Get("groupId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetServicePlanId gets the servicePlanId property value. The unique identifier for the service plan that indicates which size of the Cloud PC to provision for the user. Use a null value, when the provisioningType is dedicated.
// returns a *string when successful
func (m *CloudPcManagementGroupAssignmentTarget) GetServicePlanId()(*string) {
    val, err := m.GetBackingStore().Get("servicePlanId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CloudPcManagementGroupAssignmentTarget) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.CloudPcManagementAssignmentTarget.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("groupId", m.GetGroupId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("servicePlanId", m.GetServicePlanId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetGroupId sets the groupId property value. The ID of the target group for the assignment.
func (m *CloudPcManagementGroupAssignmentTarget) SetGroupId(value *string)() {
    err := m.GetBackingStore().Set("groupId", value)
    if err != nil {
        panic(err)
    }
}
// SetServicePlanId sets the servicePlanId property value. The unique identifier for the service plan that indicates which size of the Cloud PC to provision for the user. Use a null value, when the provisioningType is dedicated.
func (m *CloudPcManagementGroupAssignmentTarget) SetServicePlanId(value *string)() {
    err := m.GetBackingStore().Set("servicePlanId", value)
    if err != nil {
        panic(err)
    }
}
type CloudPcManagementGroupAssignmentTargetable interface {
    CloudPcManagementAssignmentTargetable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetGroupId()(*string)
    GetServicePlanId()(*string)
    SetGroupId(value *string)()
    SetServicePlanId(value *string)()
}
