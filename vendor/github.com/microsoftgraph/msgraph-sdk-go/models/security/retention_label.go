package security

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242 "github.com/microsoftgraph/msgraph-sdk-go/models"
)

type RetentionLabel struct {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entity
}
// NewRetentionLabel instantiates a new RetentionLabel and sets the default values.
func NewRetentionLabel()(*RetentionLabel) {
    m := &RetentionLabel{
        Entity: *iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.NewEntity(),
    }
    return m
}
// CreateRetentionLabelFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateRetentionLabelFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewRetentionLabel(), nil
}
// GetActionAfterRetentionPeriod gets the actionAfterRetentionPeriod property value. Specifies the action to take on the labeled document after the period specified by the retentionDuration property expires. The possible values are: none, delete, startDispositionReview, unknownFutureValue.
// returns a *ActionAfterRetentionPeriod when successful
func (m *RetentionLabel) GetActionAfterRetentionPeriod()(*ActionAfterRetentionPeriod) {
    val, err := m.GetBackingStore().Get("actionAfterRetentionPeriod")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ActionAfterRetentionPeriod)
    }
    return nil
}
// GetBehaviorDuringRetentionPeriod gets the behaviorDuringRetentionPeriod property value. Specifies how the behavior of a document with this label should be during the retention period. The possible values are: doNotRetain, retain, retainAsRecord, retainAsRegulatoryRecord, unknownFutureValue.
// returns a *BehaviorDuringRetentionPeriod when successful
func (m *RetentionLabel) GetBehaviorDuringRetentionPeriod()(*BehaviorDuringRetentionPeriod) {
    val, err := m.GetBackingStore().Get("behaviorDuringRetentionPeriod")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*BehaviorDuringRetentionPeriod)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. Represents the user who created the retentionLabel.
// returns a IdentitySetable when successful
func (m *RetentionLabel) GetCreatedBy()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Represents the date and time in which the retentionLabel is created.
// returns a *Time when successful
func (m *RetentionLabel) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDefaultRecordBehavior gets the defaultRecordBehavior property value. Specifies the locked or unlocked state of a record label when it is created.The possible values are: startLocked, startUnlocked, unknownFutureValue.
// returns a *DefaultRecordBehavior when successful
func (m *RetentionLabel) GetDefaultRecordBehavior()(*DefaultRecordBehavior) {
    val, err := m.GetBackingStore().Get("defaultRecordBehavior")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*DefaultRecordBehavior)
    }
    return nil
}
// GetDescriptionForAdmins gets the descriptionForAdmins property value. Provides label information for the admin. Optional.
// returns a *string when successful
func (m *RetentionLabel) GetDescriptionForAdmins()(*string) {
    val, err := m.GetBackingStore().Get("descriptionForAdmins")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDescriptionForUsers gets the descriptionForUsers property value. Provides the label information for the user. Optional.
// returns a *string when successful
func (m *RetentionLabel) GetDescriptionForUsers()(*string) {
    val, err := m.GetBackingStore().Get("descriptionForUsers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDescriptors gets the descriptors property value. Represents out-of-the-box values that provide more options to improve the manageability and organization of the content you need to label.
// returns a FilePlanDescriptorable when successful
func (m *RetentionLabel) GetDescriptors()(FilePlanDescriptorable) {
    val, err := m.GetBackingStore().Get("descriptors")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FilePlanDescriptorable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Unique string that defines a label name.
// returns a *string when successful
func (m *RetentionLabel) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDispositionReviewStages gets the dispositionReviewStages property value. When action at the end of retention is chosen as 'dispositionReview', dispositionReviewStages specifies a sequential set of stages with at least one reviewer in each stage.
// returns a []DispositionReviewStageable when successful
func (m *RetentionLabel) GetDispositionReviewStages()([]DispositionReviewStageable) {
    val, err := m.GetBackingStore().Get("dispositionReviewStages")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DispositionReviewStageable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *RetentionLabel) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["actionAfterRetentionPeriod"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseActionAfterRetentionPeriod)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActionAfterRetentionPeriod(val.(*ActionAfterRetentionPeriod))
        }
        return nil
    }
    res["behaviorDuringRetentionPeriod"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseBehaviorDuringRetentionPeriod)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBehaviorDuringRetentionPeriod(val.(*BehaviorDuringRetentionPeriod))
        }
        return nil
    }
    res["createdBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedBy(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable))
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["defaultRecordBehavior"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseDefaultRecordBehavior)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDefaultRecordBehavior(val.(*DefaultRecordBehavior))
        }
        return nil
    }
    res["descriptionForAdmins"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescriptionForAdmins(val)
        }
        return nil
    }
    res["descriptionForUsers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescriptionForUsers(val)
        }
        return nil
    }
    res["descriptors"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFilePlanDescriptorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescriptors(val.(FilePlanDescriptorable))
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["dispositionReviewStages"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDispositionReviewStageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DispositionReviewStageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DispositionReviewStageable)
                }
            }
            m.SetDispositionReviewStages(res)
        }
        return nil
    }
    res["isInUse"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsInUse(val)
        }
        return nil
    }
    res["labelToBeApplied"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLabelToBeApplied(val)
        }
        return nil
    }
    res["lastModifiedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedBy(val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable))
        }
        return nil
    }
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
        }
        return nil
    }
    res["retentionDuration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRetentionDurationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRetentionDuration(val.(RetentionDurationable))
        }
        return nil
    }
    res["retentionEventType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRetentionEventTypeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRetentionEventType(val.(RetentionEventTypeable))
        }
        return nil
    }
    res["retentionTrigger"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRetentionTrigger)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRetentionTrigger(val.(*RetentionTrigger))
        }
        return nil
    }
    return res
}
// GetIsInUse gets the isInUse property value. Specifies whether the label is currently being used.
// returns a *bool when successful
func (m *RetentionLabel) GetIsInUse()(*bool) {
    val, err := m.GetBackingStore().Get("isInUse")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLabelToBeApplied gets the labelToBeApplied property value. Specifies the replacement label to be applied automatically after the retention period of the current label ends.
// returns a *string when successful
func (m *RetentionLabel) GetLabelToBeApplied()(*string) {
    val, err := m.GetBackingStore().Get("labelToBeApplied")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLastModifiedBy gets the lastModifiedBy property value. The user who last modified the retentionLabel.
// returns a IdentitySetable when successful
func (m *RetentionLabel) GetLastModifiedBy()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The latest date time when the retentionLabel was modified.
// returns a *Time when successful
func (m *RetentionLabel) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetRetentionDuration gets the retentionDuration property value. Specifies the number of days to retain the content.
// returns a RetentionDurationable when successful
func (m *RetentionLabel) GetRetentionDuration()(RetentionDurationable) {
    val, err := m.GetBackingStore().Get("retentionDuration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RetentionDurationable)
    }
    return nil
}
// GetRetentionEventType gets the retentionEventType property value. Represents the type associated with a retention event.
// returns a RetentionEventTypeable when successful
func (m *RetentionLabel) GetRetentionEventType()(RetentionEventTypeable) {
    val, err := m.GetBackingStore().Get("retentionEventType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RetentionEventTypeable)
    }
    return nil
}
// GetRetentionTrigger gets the retentionTrigger property value. Specifies whether the retention duration is calculated from the content creation date, labeled date, or last modification date. The possible values are: dateLabeled, dateCreated, dateModified, dateOfEvent, unknownFutureValue.
// returns a *RetentionTrigger when successful
func (m *RetentionLabel) GetRetentionTrigger()(*RetentionTrigger) {
    val, err := m.GetBackingStore().Get("retentionTrigger")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RetentionTrigger)
    }
    return nil
}
// Serialize serializes information the current object
func (m *RetentionLabel) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetActionAfterRetentionPeriod() != nil {
        cast := (*m.GetActionAfterRetentionPeriod()).String()
        err = writer.WriteStringValue("actionAfterRetentionPeriod", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetBehaviorDuringRetentionPeriod() != nil {
        cast := (*m.GetBehaviorDuringRetentionPeriod()).String()
        err = writer.WriteStringValue("behaviorDuringRetentionPeriod", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("createdBy", m.GetCreatedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetDefaultRecordBehavior() != nil {
        cast := (*m.GetDefaultRecordBehavior()).String()
        err = writer.WriteStringValue("defaultRecordBehavior", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("descriptionForAdmins", m.GetDescriptionForAdmins())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("descriptionForUsers", m.GetDescriptionForUsers())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("descriptors", m.GetDescriptors())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetDispositionReviewStages() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetDispositionReviewStages()))
        for i, v := range m.GetDispositionReviewStages() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("dispositionReviewStages", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isInUse", m.GetIsInUse())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("labelToBeApplied", m.GetLabelToBeApplied())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastModifiedBy", m.GetLastModifiedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("retentionDuration", m.GetRetentionDuration())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("retentionEventType", m.GetRetentionEventType())
        if err != nil {
            return err
        }
    }
    if m.GetRetentionTrigger() != nil {
        cast := (*m.GetRetentionTrigger()).String()
        err = writer.WriteStringValue("retentionTrigger", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActionAfterRetentionPeriod sets the actionAfterRetentionPeriod property value. Specifies the action to take on the labeled document after the period specified by the retentionDuration property expires. The possible values are: none, delete, startDispositionReview, unknownFutureValue.
func (m *RetentionLabel) SetActionAfterRetentionPeriod(value *ActionAfterRetentionPeriod)() {
    err := m.GetBackingStore().Set("actionAfterRetentionPeriod", value)
    if err != nil {
        panic(err)
    }
}
// SetBehaviorDuringRetentionPeriod sets the behaviorDuringRetentionPeriod property value. Specifies how the behavior of a document with this label should be during the retention period. The possible values are: doNotRetain, retain, retainAsRecord, retainAsRegulatoryRecord, unknownFutureValue.
func (m *RetentionLabel) SetBehaviorDuringRetentionPeriod(value *BehaviorDuringRetentionPeriod)() {
    err := m.GetBackingStore().Set("behaviorDuringRetentionPeriod", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. Represents the user who created the retentionLabel.
func (m *RetentionLabel) SetCreatedBy(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Represents the date and time in which the retentionLabel is created.
func (m *RetentionLabel) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDefaultRecordBehavior sets the defaultRecordBehavior property value. Specifies the locked or unlocked state of a record label when it is created.The possible values are: startLocked, startUnlocked, unknownFutureValue.
func (m *RetentionLabel) SetDefaultRecordBehavior(value *DefaultRecordBehavior)() {
    err := m.GetBackingStore().Set("defaultRecordBehavior", value)
    if err != nil {
        panic(err)
    }
}
// SetDescriptionForAdmins sets the descriptionForAdmins property value. Provides label information for the admin. Optional.
func (m *RetentionLabel) SetDescriptionForAdmins(value *string)() {
    err := m.GetBackingStore().Set("descriptionForAdmins", value)
    if err != nil {
        panic(err)
    }
}
// SetDescriptionForUsers sets the descriptionForUsers property value. Provides the label information for the user. Optional.
func (m *RetentionLabel) SetDescriptionForUsers(value *string)() {
    err := m.GetBackingStore().Set("descriptionForUsers", value)
    if err != nil {
        panic(err)
    }
}
// SetDescriptors sets the descriptors property value. Represents out-of-the-box values that provide more options to improve the manageability and organization of the content you need to label.
func (m *RetentionLabel) SetDescriptors(value FilePlanDescriptorable)() {
    err := m.GetBackingStore().Set("descriptors", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Unique string that defines a label name.
func (m *RetentionLabel) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDispositionReviewStages sets the dispositionReviewStages property value. When action at the end of retention is chosen as 'dispositionReview', dispositionReviewStages specifies a sequential set of stages with at least one reviewer in each stage.
func (m *RetentionLabel) SetDispositionReviewStages(value []DispositionReviewStageable)() {
    err := m.GetBackingStore().Set("dispositionReviewStages", value)
    if err != nil {
        panic(err)
    }
}
// SetIsInUse sets the isInUse property value. Specifies whether the label is currently being used.
func (m *RetentionLabel) SetIsInUse(value *bool)() {
    err := m.GetBackingStore().Set("isInUse", value)
    if err != nil {
        panic(err)
    }
}
// SetLabelToBeApplied sets the labelToBeApplied property value. Specifies the replacement label to be applied automatically after the retention period of the current label ends.
func (m *RetentionLabel) SetLabelToBeApplied(value *string)() {
    err := m.GetBackingStore().Set("labelToBeApplied", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. The user who last modified the retentionLabel.
func (m *RetentionLabel) SetLastModifiedBy(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The latest date time when the retentionLabel was modified.
func (m *RetentionLabel) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetRetentionDuration sets the retentionDuration property value. Specifies the number of days to retain the content.
func (m *RetentionLabel) SetRetentionDuration(value RetentionDurationable)() {
    err := m.GetBackingStore().Set("retentionDuration", value)
    if err != nil {
        panic(err)
    }
}
// SetRetentionEventType sets the retentionEventType property value. Represents the type associated with a retention event.
func (m *RetentionLabel) SetRetentionEventType(value RetentionEventTypeable)() {
    err := m.GetBackingStore().Set("retentionEventType", value)
    if err != nil {
        panic(err)
    }
}
// SetRetentionTrigger sets the retentionTrigger property value. Specifies whether the retention duration is calculated from the content creation date, labeled date, or last modification date. The possible values are: dateLabeled, dateCreated, dateModified, dateOfEvent, unknownFutureValue.
func (m *RetentionLabel) SetRetentionTrigger(value *RetentionTrigger)() {
    err := m.GetBackingStore().Set("retentionTrigger", value)
    if err != nil {
        panic(err)
    }
}
type RetentionLabelable interface {
    iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActionAfterRetentionPeriod()(*ActionAfterRetentionPeriod)
    GetBehaviorDuringRetentionPeriod()(*BehaviorDuringRetentionPeriod)
    GetCreatedBy()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDefaultRecordBehavior()(*DefaultRecordBehavior)
    GetDescriptionForAdmins()(*string)
    GetDescriptionForUsers()(*string)
    GetDescriptors()(FilePlanDescriptorable)
    GetDisplayName()(*string)
    GetDispositionReviewStages()([]DispositionReviewStageable)
    GetIsInUse()(*bool)
    GetLabelToBeApplied()(*string)
    GetLastModifiedBy()(iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetRetentionDuration()(RetentionDurationable)
    GetRetentionEventType()(RetentionEventTypeable)
    GetRetentionTrigger()(*RetentionTrigger)
    SetActionAfterRetentionPeriod(value *ActionAfterRetentionPeriod)()
    SetBehaviorDuringRetentionPeriod(value *BehaviorDuringRetentionPeriod)()
    SetCreatedBy(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDefaultRecordBehavior(value *DefaultRecordBehavior)()
    SetDescriptionForAdmins(value *string)()
    SetDescriptionForUsers(value *string)()
    SetDescriptors(value FilePlanDescriptorable)()
    SetDisplayName(value *string)()
    SetDispositionReviewStages(value []DispositionReviewStageable)()
    SetIsInUse(value *bool)()
    SetLabelToBeApplied(value *string)()
    SetLastModifiedBy(value iadcd81124412c61e647227ecfc4449d8bba17de0380ddda76f641a29edf2b242.IdentitySetable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetRetentionDuration(value RetentionDurationable)()
    SetRetentionEventType(value RetentionEventTypeable)()
    SetRetentionTrigger(value *RetentionTrigger)()
}
