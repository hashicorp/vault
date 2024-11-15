package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AccessReviewInstanceDecisionItem struct {
    Entity
}
// NewAccessReviewInstanceDecisionItem instantiates a new AccessReviewInstanceDecisionItem and sets the default values.
func NewAccessReviewInstanceDecisionItem()(*AccessReviewInstanceDecisionItem) {
    m := &AccessReviewInstanceDecisionItem{
        Entity: *NewEntity(),
    }
    return m
}
// CreateAccessReviewInstanceDecisionItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessReviewInstanceDecisionItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessReviewInstanceDecisionItem(), nil
}
// GetAccessReviewId gets the accessReviewId property value. The identifier of the accessReviewInstance parent. Supports $select. Read-only.
// returns a *string when successful
func (m *AccessReviewInstanceDecisionItem) GetAccessReviewId()(*string) {
    val, err := m.GetBackingStore().Get("accessReviewId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppliedBy gets the appliedBy property value. The identifier of the user who applied the decision. Read-only.
// returns a UserIdentityable when successful
func (m *AccessReviewInstanceDecisionItem) GetAppliedBy()(UserIdentityable) {
    val, err := m.GetBackingStore().Get("appliedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserIdentityable)
    }
    return nil
}
// GetAppliedDateTime gets the appliedDateTime property value. The timestamp when the approval decision was applied.00000000-0000-0000-0000-000000000000 if the assigned reviewer hasn't applied the decision or it was automatically applied. The DatetimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.  Supports $select. Read-only.
// returns a *Time when successful
func (m *AccessReviewInstanceDecisionItem) GetAppliedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("appliedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetApplyResult gets the applyResult property value. The result of applying the decision. Possible values: New, AppliedSuccessfully, AppliedWithUnknownFailure, AppliedSuccessfullyButObjectNotFound and ApplyNotSupported. Supports $select, $orderby, and $filter (eq only). Read-only.
// returns a *string when successful
func (m *AccessReviewInstanceDecisionItem) GetApplyResult()(*string) {
    val, err := m.GetBackingStore().Get("applyResult")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDecision gets the decision property value. Result of the review. Possible values: Approve, Deny, NotReviewed, or DontKnow. Supports $select, $orderby, and $filter (eq only).
// returns a *string when successful
func (m *AccessReviewInstanceDecisionItem) GetDecision()(*string) {
    val, err := m.GetBackingStore().Get("decision")
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
func (m *AccessReviewInstanceDecisionItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["accessReviewId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccessReviewId(val)
        }
        return nil
    }
    res["appliedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppliedBy(val.(UserIdentityable))
        }
        return nil
    }
    res["appliedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppliedDateTime(val)
        }
        return nil
    }
    res["applyResult"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplyResult(val)
        }
        return nil
    }
    res["decision"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDecision(val)
        }
        return nil
    }
    res["insights"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateGovernanceInsightFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]GovernanceInsightable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(GovernanceInsightable)
                }
            }
            m.SetInsights(res)
        }
        return nil
    }
    res["justification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetJustification(val)
        }
        return nil
    }
    res["principal"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrincipal(val.(Identityable))
        }
        return nil
    }
    res["principalLink"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrincipalLink(val)
        }
        return nil
    }
    res["recommendation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecommendation(val)
        }
        return nil
    }
    res["resource"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccessReviewInstanceDecisionItemResourceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResource(val.(AccessReviewInstanceDecisionItemResourceable))
        }
        return nil
    }
    res["resourceLink"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceLink(val)
        }
        return nil
    }
    res["reviewedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUserIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReviewedBy(val.(UserIdentityable))
        }
        return nil
    }
    res["reviewedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReviewedDateTime(val)
        }
        return nil
    }
    return res
}
// GetInsights gets the insights property value. Insights are recommendations to reviewers on whether to approve or deny a decision. There can be multiple insights associated with an accessReviewInstanceDecisionItem.
// returns a []GovernanceInsightable when successful
func (m *AccessReviewInstanceDecisionItem) GetInsights()([]GovernanceInsightable) {
    val, err := m.GetBackingStore().Get("insights")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]GovernanceInsightable)
    }
    return nil
}
// GetJustification gets the justification property value. Justification left by the reviewer when they made the decision.
// returns a *string when successful
func (m *AccessReviewInstanceDecisionItem) GetJustification()(*string) {
    val, err := m.GetBackingStore().Get("justification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPrincipal gets the principal property value. Every decision item in an access review represents a principal's access to a resource. This property represents details of the principal. For example, if a decision item represents access of User 'Bob' to Group 'Sales' - The principal is 'Bob' and the resource is 'Sales'. Principals can be of two types - userIdentity and servicePrincipalIdentity. Supports $select. Read-only.
// returns a Identityable when successful
func (m *AccessReviewInstanceDecisionItem) GetPrincipal()(Identityable) {
    val, err := m.GetBackingStore().Get("principal")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Identityable)
    }
    return nil
}
// GetPrincipalLink gets the principalLink property value. A link to the principal object. For example, https://graph.microsoft.com/v1.0/users/a6c7aecb-cbfd-4763-87ef-e91b4bd509d9. Read-only.
// returns a *string when successful
func (m *AccessReviewInstanceDecisionItem) GetPrincipalLink()(*string) {
    val, err := m.GetBackingStore().Get("principalLink")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRecommendation gets the recommendation property value. A system-generated recommendation for the approval decision based off last interactive sign-in to tenant. The value is Approve if the sign-in is fewer than 30 days after the start of review, Deny if the sign-in is greater than 30 days after, or NoInfoAvailable. Possible values: Approve, Deny, or NoInfoAvailable. Supports $select, $orderby, and $filter (eq only). Read-only.
// returns a *string when successful
func (m *AccessReviewInstanceDecisionItem) GetRecommendation()(*string) {
    val, err := m.GetBackingStore().Get("recommendation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResource gets the resource property value. Every decision item in an access review represents a principal's access to a resource. This property represents details of the resource. For example, if a decision item represents access of User 'Bob' to Group 'Sales' - The principal is Bob and the resource is 'Sales'. Resources can be of multiple types. See accessReviewInstanceDecisionItemResource. Read-only.
// returns a AccessReviewInstanceDecisionItemResourceable when successful
func (m *AccessReviewInstanceDecisionItem) GetResource()(AccessReviewInstanceDecisionItemResourceable) {
    val, err := m.GetBackingStore().Get("resource")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccessReviewInstanceDecisionItemResourceable)
    }
    return nil
}
// GetResourceLink gets the resourceLink property value. A link to the resource. For example, https://graph.microsoft.com/v1.0/servicePrincipals/c86300f3-8695-4320-9f6e-32a2555f5ff8. Supports $select. Read-only.
// returns a *string when successful
func (m *AccessReviewInstanceDecisionItem) GetResourceLink()(*string) {
    val, err := m.GetBackingStore().Get("resourceLink")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetReviewedBy gets the reviewedBy property value. The identifier of the reviewer.00000000-0000-0000-0000-000000000000 if the assigned reviewer hasn't reviewed. Supports $select. Read-only.
// returns a UserIdentityable when successful
func (m *AccessReviewInstanceDecisionItem) GetReviewedBy()(UserIdentityable) {
    val, err := m.GetBackingStore().Get("reviewedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UserIdentityable)
    }
    return nil
}
// GetReviewedDateTime gets the reviewedDateTime property value. The timestamp when the review decision occurred. Supports $select. Read-only.
// returns a *Time when successful
func (m *AccessReviewInstanceDecisionItem) GetReviewedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("reviewedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessReviewInstanceDecisionItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("accessReviewId", m.GetAccessReviewId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("appliedBy", m.GetAppliedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("appliedDateTime", m.GetAppliedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("applyResult", m.GetApplyResult())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("decision", m.GetDecision())
        if err != nil {
            return err
        }
    }
    if m.GetInsights() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInsights()))
        for i, v := range m.GetInsights() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("insights", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("justification", m.GetJustification())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("principal", m.GetPrincipal())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("principalLink", m.GetPrincipalLink())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("recommendation", m.GetRecommendation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("resource", m.GetResource())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("resourceLink", m.GetResourceLink())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("reviewedBy", m.GetReviewedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("reviewedDateTime", m.GetReviewedDateTime())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccessReviewId sets the accessReviewId property value. The identifier of the accessReviewInstance parent. Supports $select. Read-only.
func (m *AccessReviewInstanceDecisionItem) SetAccessReviewId(value *string)() {
    err := m.GetBackingStore().Set("accessReviewId", value)
    if err != nil {
        panic(err)
    }
}
// SetAppliedBy sets the appliedBy property value. The identifier of the user who applied the decision. Read-only.
func (m *AccessReviewInstanceDecisionItem) SetAppliedBy(value UserIdentityable)() {
    err := m.GetBackingStore().Set("appliedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetAppliedDateTime sets the appliedDateTime property value. The timestamp when the approval decision was applied.00000000-0000-0000-0000-000000000000 if the assigned reviewer hasn't applied the decision or it was automatically applied. The DatetimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.  Supports $select. Read-only.
func (m *AccessReviewInstanceDecisionItem) SetAppliedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("appliedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetApplyResult sets the applyResult property value. The result of applying the decision. Possible values: New, AppliedSuccessfully, AppliedWithUnknownFailure, AppliedSuccessfullyButObjectNotFound and ApplyNotSupported. Supports $select, $orderby, and $filter (eq only). Read-only.
func (m *AccessReviewInstanceDecisionItem) SetApplyResult(value *string)() {
    err := m.GetBackingStore().Set("applyResult", value)
    if err != nil {
        panic(err)
    }
}
// SetDecision sets the decision property value. Result of the review. Possible values: Approve, Deny, NotReviewed, or DontKnow. Supports $select, $orderby, and $filter (eq only).
func (m *AccessReviewInstanceDecisionItem) SetDecision(value *string)() {
    err := m.GetBackingStore().Set("decision", value)
    if err != nil {
        panic(err)
    }
}
// SetInsights sets the insights property value. Insights are recommendations to reviewers on whether to approve or deny a decision. There can be multiple insights associated with an accessReviewInstanceDecisionItem.
func (m *AccessReviewInstanceDecisionItem) SetInsights(value []GovernanceInsightable)() {
    err := m.GetBackingStore().Set("insights", value)
    if err != nil {
        panic(err)
    }
}
// SetJustification sets the justification property value. Justification left by the reviewer when they made the decision.
func (m *AccessReviewInstanceDecisionItem) SetJustification(value *string)() {
    err := m.GetBackingStore().Set("justification", value)
    if err != nil {
        panic(err)
    }
}
// SetPrincipal sets the principal property value. Every decision item in an access review represents a principal's access to a resource. This property represents details of the principal. For example, if a decision item represents access of User 'Bob' to Group 'Sales' - The principal is 'Bob' and the resource is 'Sales'. Principals can be of two types - userIdentity and servicePrincipalIdentity. Supports $select. Read-only.
func (m *AccessReviewInstanceDecisionItem) SetPrincipal(value Identityable)() {
    err := m.GetBackingStore().Set("principal", value)
    if err != nil {
        panic(err)
    }
}
// SetPrincipalLink sets the principalLink property value. A link to the principal object. For example, https://graph.microsoft.com/v1.0/users/a6c7aecb-cbfd-4763-87ef-e91b4bd509d9. Read-only.
func (m *AccessReviewInstanceDecisionItem) SetPrincipalLink(value *string)() {
    err := m.GetBackingStore().Set("principalLink", value)
    if err != nil {
        panic(err)
    }
}
// SetRecommendation sets the recommendation property value. A system-generated recommendation for the approval decision based off last interactive sign-in to tenant. The value is Approve if the sign-in is fewer than 30 days after the start of review, Deny if the sign-in is greater than 30 days after, or NoInfoAvailable. Possible values: Approve, Deny, or NoInfoAvailable. Supports $select, $orderby, and $filter (eq only). Read-only.
func (m *AccessReviewInstanceDecisionItem) SetRecommendation(value *string)() {
    err := m.GetBackingStore().Set("recommendation", value)
    if err != nil {
        panic(err)
    }
}
// SetResource sets the resource property value. Every decision item in an access review represents a principal's access to a resource. This property represents details of the resource. For example, if a decision item represents access of User 'Bob' to Group 'Sales' - The principal is Bob and the resource is 'Sales'. Resources can be of multiple types. See accessReviewInstanceDecisionItemResource. Read-only.
func (m *AccessReviewInstanceDecisionItem) SetResource(value AccessReviewInstanceDecisionItemResourceable)() {
    err := m.GetBackingStore().Set("resource", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceLink sets the resourceLink property value. A link to the resource. For example, https://graph.microsoft.com/v1.0/servicePrincipals/c86300f3-8695-4320-9f6e-32a2555f5ff8. Supports $select. Read-only.
func (m *AccessReviewInstanceDecisionItem) SetResourceLink(value *string)() {
    err := m.GetBackingStore().Set("resourceLink", value)
    if err != nil {
        panic(err)
    }
}
// SetReviewedBy sets the reviewedBy property value. The identifier of the reviewer.00000000-0000-0000-0000-000000000000 if the assigned reviewer hasn't reviewed. Supports $select. Read-only.
func (m *AccessReviewInstanceDecisionItem) SetReviewedBy(value UserIdentityable)() {
    err := m.GetBackingStore().Set("reviewedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetReviewedDateTime sets the reviewedDateTime property value. The timestamp when the review decision occurred. Supports $select. Read-only.
func (m *AccessReviewInstanceDecisionItem) SetReviewedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("reviewedDateTime", value)
    if err != nil {
        panic(err)
    }
}
type AccessReviewInstanceDecisionItemable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccessReviewId()(*string)
    GetAppliedBy()(UserIdentityable)
    GetAppliedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetApplyResult()(*string)
    GetDecision()(*string)
    GetInsights()([]GovernanceInsightable)
    GetJustification()(*string)
    GetPrincipal()(Identityable)
    GetPrincipalLink()(*string)
    GetRecommendation()(*string)
    GetResource()(AccessReviewInstanceDecisionItemResourceable)
    GetResourceLink()(*string)
    GetReviewedBy()(UserIdentityable)
    GetReviewedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    SetAccessReviewId(value *string)()
    SetAppliedBy(value UserIdentityable)()
    SetAppliedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetApplyResult(value *string)()
    SetDecision(value *string)()
    SetInsights(value []GovernanceInsightable)()
    SetJustification(value *string)()
    SetPrincipal(value Identityable)()
    SetPrincipalLink(value *string)()
    SetRecommendation(value *string)()
    SetResource(value AccessReviewInstanceDecisionItemResourceable)()
    SetResourceLink(value *string)()
    SetReviewedBy(value UserIdentityable)()
    SetReviewedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
}
