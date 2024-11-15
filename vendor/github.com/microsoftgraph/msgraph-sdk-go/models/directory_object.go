package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DirectoryObject struct {
    Entity
}
// NewDirectoryObject instantiates a new DirectoryObject and sets the default values.
func NewDirectoryObject()(*DirectoryObject) {
    m := &DirectoryObject{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDirectoryObjectFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDirectoryObjectFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.activityBasedTimeoutPolicy":
                        return NewActivityBasedTimeoutPolicy(), nil
                    case "#microsoft.graph.administrativeUnit":
                        return NewAdministrativeUnit(), nil
                    case "#microsoft.graph.application":
                        return NewApplication(), nil
                    case "#microsoft.graph.appManagementPolicy":
                        return NewAppManagementPolicy(), nil
                    case "#microsoft.graph.appRoleAssignment":
                        return NewAppRoleAssignment(), nil
                    case "#microsoft.graph.authorizationPolicy":
                        return NewAuthorizationPolicy(), nil
                    case "#microsoft.graph.claimsMappingPolicy":
                        return NewClaimsMappingPolicy(), nil
                    case "#microsoft.graph.contract":
                        return NewContract(), nil
                    case "#microsoft.graph.crossTenantAccessPolicy":
                        return NewCrossTenantAccessPolicy(), nil
                    case "#microsoft.graph.device":
                        return NewDevice(), nil
                    case "#microsoft.graph.directoryObjectPartnerReference":
                        return NewDirectoryObjectPartnerReference(), nil
                    case "#microsoft.graph.directoryRole":
                        return NewDirectoryRole(), nil
                    case "#microsoft.graph.directoryRoleTemplate":
                        return NewDirectoryRoleTemplate(), nil
                    case "#microsoft.graph.endpoint":
                        return NewEndpoint(), nil
                    case "#microsoft.graph.extensionProperty":
                        return NewExtensionProperty(), nil
                    case "#microsoft.graph.group":
                        return NewGroup(), nil
                    case "#microsoft.graph.groupSettingTemplate":
                        return NewGroupSettingTemplate(), nil
                    case "#microsoft.graph.homeRealmDiscoveryPolicy":
                        return NewHomeRealmDiscoveryPolicy(), nil
                    case "#microsoft.graph.identitySecurityDefaultsEnforcementPolicy":
                        return NewIdentitySecurityDefaultsEnforcementPolicy(), nil
                    case "#microsoft.graph.multiTenantOrganizationMember":
                        return NewMultiTenantOrganizationMember(), nil
                    case "#microsoft.graph.organization":
                        return NewOrganization(), nil
                    case "#microsoft.graph.orgContact":
                        return NewOrgContact(), nil
                    case "#microsoft.graph.permissionGrantPolicy":
                        return NewPermissionGrantPolicy(), nil
                    case "#microsoft.graph.policyBase":
                        return NewPolicyBase(), nil
                    case "#microsoft.graph.resourceSpecificPermissionGrant":
                        return NewResourceSpecificPermissionGrant(), nil
                    case "#microsoft.graph.servicePrincipal":
                        return NewServicePrincipal(), nil
                    case "#microsoft.graph.stsPolicy":
                        return NewStsPolicy(), nil
                    case "#microsoft.graph.tenantAppManagementPolicy":
                        return NewTenantAppManagementPolicy(), nil
                    case "#microsoft.graph.tokenIssuancePolicy":
                        return NewTokenIssuancePolicy(), nil
                    case "#microsoft.graph.tokenLifetimePolicy":
                        return NewTokenLifetimePolicy(), nil
                    case "#microsoft.graph.user":
                        return NewUser(), nil
                }
            }
        }
    }
    return NewDirectoryObject(), nil
}
// GetDeletedDateTime gets the deletedDateTime property value. Date and time when this object was deleted. Always null when the object hasn't been deleted.
// returns a *Time when successful
func (m *DirectoryObject) GetDeletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("deletedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DirectoryObject) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["deletedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeletedDateTime(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *DirectoryObject) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteTimeValue("deletedDateTime", m.GetDeletedDateTime())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDeletedDateTime sets the deletedDateTime property value. Date and time when this object was deleted. Always null when the object hasn't been deleted.
func (m *DirectoryObject) SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("deletedDateTime", value)
    if err != nil {
        panic(err)
    }
}
type DirectoryObjectable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDeletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
}
