package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Simulation struct {
    Entity
}
// NewSimulation instantiates a new Simulation and sets the default values.
func NewSimulation()(*Simulation) {
    m := &Simulation{
        Entity: *NewEntity(),
    }
    return m
}
// CreateSimulationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSimulationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSimulation(), nil
}
// GetAttackTechnique gets the attackTechnique property value. The social engineering technique used in the attack simulation and training campaign. Supports $filter and $orderby. Possible values are: unknown, credentialHarvesting, attachmentMalware, driveByUrl, linkInAttachment, linkToMalwareFile, unknownFutureValue, oAuthConsentGrant. Note that you must use the Prefer: include-unknown-enum-members request header to get the following values from this evolvable enum: oAuthConsentGrant. For more information on the types of social engineering attack techniques, see simulations.
// returns a *SimulationAttackTechnique when successful
func (m *Simulation) GetAttackTechnique()(*SimulationAttackTechnique) {
    val, err := m.GetBackingStore().Get("attackTechnique")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SimulationAttackTechnique)
    }
    return nil
}
// GetAttackType gets the attackType property value. Attack type of the attack simulation and training campaign. Supports $filter and $orderby. Possible values are: unknown, social, cloud, endpoint, unknownFutureValue.
// returns a *SimulationAttackType when successful
func (m *Simulation) GetAttackType()(*SimulationAttackType) {
    val, err := m.GetBackingStore().Get("attackType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SimulationAttackType)
    }
    return nil
}
// GetAutomationId gets the automationId property value. Unique identifier for the attack simulation automation.
// returns a *string when successful
func (m *Simulation) GetAutomationId()(*string) {
    val, err := m.GetBackingStore().Get("automationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCompletionDateTime gets the completionDateTime property value. Date and time of completion of the attack simulation and training campaign. Supports $filter and $orderby.
// returns a *Time when successful
func (m *Simulation) GetCompletionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("completionDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetCreatedBy gets the createdBy property value. Identity of the user who created the attack simulation and training campaign.
// returns a EmailIdentityable when successful
func (m *Simulation) GetCreatedBy()(EmailIdentityable) {
    val, err := m.GetBackingStore().Get("createdBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EmailIdentityable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Date and time of creation of the attack simulation and training campaign.
// returns a *Time when successful
func (m *Simulation) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDescription gets the description property value. Description of the attack simulation and training campaign.
// returns a *string when successful
func (m *Simulation) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Display name of the attack simulation and training campaign. Supports $filter and $orderby.
// returns a *string when successful
func (m *Simulation) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDurationInDays gets the durationInDays property value. Simulation duration in days.
// returns a *int32 when successful
func (m *Simulation) GetDurationInDays()(*int32) {
    val, err := m.GetBackingStore().Get("durationInDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetEndUserNotificationSetting gets the endUserNotificationSetting property value. Details about the end user notification setting.
// returns a EndUserNotificationSettingable when successful
func (m *Simulation) GetEndUserNotificationSetting()(EndUserNotificationSettingable) {
    val, err := m.GetBackingStore().Get("endUserNotificationSetting")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EndUserNotificationSettingable)
    }
    return nil
}
// GetExcludedAccountTarget gets the excludedAccountTarget property value. Users excluded from the simulation.
// returns a AccountTargetContentable when successful
func (m *Simulation) GetExcludedAccountTarget()(AccountTargetContentable) {
    val, err := m.GetBackingStore().Get("excludedAccountTarget")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccountTargetContentable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Simulation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["attackTechnique"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSimulationAttackTechnique)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttackTechnique(val.(*SimulationAttackTechnique))
        }
        return nil
    }
    res["attackType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSimulationAttackType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAttackType(val.(*SimulationAttackType))
        }
        return nil
    }
    res["automationId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAutomationId(val)
        }
        return nil
    }
    res["completionDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompletionDateTime(val)
        }
        return nil
    }
    res["createdBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEmailIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedBy(val.(EmailIdentityable))
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
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
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
    res["durationInDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDurationInDays(val)
        }
        return nil
    }
    res["endUserNotificationSetting"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEndUserNotificationSettingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEndUserNotificationSetting(val.(EndUserNotificationSettingable))
        }
        return nil
    }
    res["excludedAccountTarget"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccountTargetContentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExcludedAccountTarget(val.(AccountTargetContentable))
        }
        return nil
    }
    res["includedAccountTarget"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAccountTargetContentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIncludedAccountTarget(val.(AccountTargetContentable))
        }
        return nil
    }
    res["isAutomated"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAutomated(val)
        }
        return nil
    }
    res["landingPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateLandingPageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLandingPage(val.(LandingPageable))
        }
        return nil
    }
    res["lastModifiedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEmailIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedBy(val.(EmailIdentityable))
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
    res["launchDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLaunchDateTime(val)
        }
        return nil
    }
    res["loginPage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateLoginPageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLoginPage(val.(LoginPageable))
        }
        return nil
    }
    res["oAuthConsentAppDetail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOAuthConsentAppDetailFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOAuthConsentAppDetail(val.(OAuthConsentAppDetailable))
        }
        return nil
    }
    res["payload"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePayloadFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPayload(val.(Payloadable))
        }
        return nil
    }
    res["payloadDeliveryPlatform"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePayloadDeliveryPlatform)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPayloadDeliveryPlatform(val.(*PayloadDeliveryPlatform))
        }
        return nil
    }
    res["report"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSimulationReportFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReport(val.(SimulationReportable))
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSimulationStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(*SimulationStatus))
        }
        return nil
    }
    res["trainingSetting"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTrainingSettingFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTrainingSetting(val.(TrainingSettingable))
        }
        return nil
    }
    return res
}
// GetIncludedAccountTarget gets the includedAccountTarget property value. Users targeted in the simulation.
// returns a AccountTargetContentable when successful
func (m *Simulation) GetIncludedAccountTarget()(AccountTargetContentable) {
    val, err := m.GetBackingStore().Get("includedAccountTarget")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AccountTargetContentable)
    }
    return nil
}
// GetIsAutomated gets the isAutomated property value. Flag that represents if the attack simulation and training campaign was created from a simulation automation flow. Supports $filter and $orderby.
// returns a *bool when successful
func (m *Simulation) GetIsAutomated()(*bool) {
    val, err := m.GetBackingStore().Get("isAutomated")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLandingPage gets the landingPage property value. The landing page associated with a simulation during its creation.
// returns a LandingPageable when successful
func (m *Simulation) GetLandingPage()(LandingPageable) {
    val, err := m.GetBackingStore().Get("landingPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(LandingPageable)
    }
    return nil
}
// GetLastModifiedBy gets the lastModifiedBy property value. Identity of the user who most recently modified the attack simulation and training campaign.
// returns a EmailIdentityable when successful
func (m *Simulation) GetLastModifiedBy()(EmailIdentityable) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EmailIdentityable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. Date and time of the most recent modification of the attack simulation and training campaign.
// returns a *Time when successful
func (m *Simulation) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLaunchDateTime gets the launchDateTime property value. Date and time of the launch/start of the attack simulation and training campaign. Supports $filter and $orderby.
// returns a *Time when successful
func (m *Simulation) GetLaunchDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("launchDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLoginPage gets the loginPage property value. The login page associated with a simulation during its creation.
// returns a LoginPageable when successful
func (m *Simulation) GetLoginPage()(LoginPageable) {
    val, err := m.GetBackingStore().Get("loginPage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(LoginPageable)
    }
    return nil
}
// GetOAuthConsentAppDetail gets the oAuthConsentAppDetail property value. OAuth app details for the OAuth technique.
// returns a OAuthConsentAppDetailable when successful
func (m *Simulation) GetOAuthConsentAppDetail()(OAuthConsentAppDetailable) {
    val, err := m.GetBackingStore().Get("oAuthConsentAppDetail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OAuthConsentAppDetailable)
    }
    return nil
}
// GetPayload gets the payload property value. The payload associated with a simulation during its creation.
// returns a Payloadable when successful
func (m *Simulation) GetPayload()(Payloadable) {
    val, err := m.GetBackingStore().Get("payload")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Payloadable)
    }
    return nil
}
// GetPayloadDeliveryPlatform gets the payloadDeliveryPlatform property value. Method of delivery of the phishing payload used in the attack simulation and training campaign. Possible values are: unknown, sms, email, teams, unknownFutureValue.
// returns a *PayloadDeliveryPlatform when successful
func (m *Simulation) GetPayloadDeliveryPlatform()(*PayloadDeliveryPlatform) {
    val, err := m.GetBackingStore().Get("payloadDeliveryPlatform")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PayloadDeliveryPlatform)
    }
    return nil
}
// GetReport gets the report property value. Report of the attack simulation and training campaign.
// returns a SimulationReportable when successful
func (m *Simulation) GetReport()(SimulationReportable) {
    val, err := m.GetBackingStore().Get("report")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SimulationReportable)
    }
    return nil
}
// GetStatus gets the status property value. Status of the attack simulation and training campaign. Supports $filter and $orderby. Possible values are: unknown, draft, running, scheduled, succeeded, failed, cancelled, excluded, unknownFutureValue.
// returns a *SimulationStatus when successful
func (m *Simulation) GetStatus()(*SimulationStatus) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SimulationStatus)
    }
    return nil
}
// GetTrainingSetting gets the trainingSetting property value. Details about the training settings for a simulation.
// returns a TrainingSettingable when successful
func (m *Simulation) GetTrainingSetting()(TrainingSettingable) {
    val, err := m.GetBackingStore().Get("trainingSetting")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TrainingSettingable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Simulation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAttackTechnique() != nil {
        cast := (*m.GetAttackTechnique()).String()
        err = writer.WriteStringValue("attackTechnique", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetAttackType() != nil {
        cast := (*m.GetAttackType()).String()
        err = writer.WriteStringValue("attackType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("automationId", m.GetAutomationId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("completionDateTime", m.GetCompletionDateTime())
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
    {
        err = writer.WriteStringValue("description", m.GetDescription())
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
    {
        err = writer.WriteInt32Value("durationInDays", m.GetDurationInDays())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("endUserNotificationSetting", m.GetEndUserNotificationSetting())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("excludedAccountTarget", m.GetExcludedAccountTarget())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("includedAccountTarget", m.GetIncludedAccountTarget())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isAutomated", m.GetIsAutomated())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("landingPage", m.GetLandingPage())
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
        err = writer.WriteTimeValue("launchDateTime", m.GetLaunchDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("loginPage", m.GetLoginPage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("oAuthConsentAppDetail", m.GetOAuthConsentAppDetail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("payload", m.GetPayload())
        if err != nil {
            return err
        }
    }
    if m.GetPayloadDeliveryPlatform() != nil {
        cast := (*m.GetPayloadDeliveryPlatform()).String()
        err = writer.WriteStringValue("payloadDeliveryPlatform", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("report", m.GetReport())
        if err != nil {
            return err
        }
    }
    if m.GetStatus() != nil {
        cast := (*m.GetStatus()).String()
        err = writer.WriteStringValue("status", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("trainingSetting", m.GetTrainingSetting())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAttackTechnique sets the attackTechnique property value. The social engineering technique used in the attack simulation and training campaign. Supports $filter and $orderby. Possible values are: unknown, credentialHarvesting, attachmentMalware, driveByUrl, linkInAttachment, linkToMalwareFile, unknownFutureValue, oAuthConsentGrant. Note that you must use the Prefer: include-unknown-enum-members request header to get the following values from this evolvable enum: oAuthConsentGrant. For more information on the types of social engineering attack techniques, see simulations.
func (m *Simulation) SetAttackTechnique(value *SimulationAttackTechnique)() {
    err := m.GetBackingStore().Set("attackTechnique", value)
    if err != nil {
        panic(err)
    }
}
// SetAttackType sets the attackType property value. Attack type of the attack simulation and training campaign. Supports $filter and $orderby. Possible values are: unknown, social, cloud, endpoint, unknownFutureValue.
func (m *Simulation) SetAttackType(value *SimulationAttackType)() {
    err := m.GetBackingStore().Set("attackType", value)
    if err != nil {
        panic(err)
    }
}
// SetAutomationId sets the automationId property value. Unique identifier for the attack simulation automation.
func (m *Simulation) SetAutomationId(value *string)() {
    err := m.GetBackingStore().Set("automationId", value)
    if err != nil {
        panic(err)
    }
}
// SetCompletionDateTime sets the completionDateTime property value. Date and time of completion of the attack simulation and training campaign. Supports $filter and $orderby.
func (m *Simulation) SetCompletionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("completionDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedBy sets the createdBy property value. Identity of the user who created the attack simulation and training campaign.
func (m *Simulation) SetCreatedBy(value EmailIdentityable)() {
    err := m.GetBackingStore().Set("createdBy", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Date and time of creation of the attack simulation and training campaign.
func (m *Simulation) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDescription sets the description property value. Description of the attack simulation and training campaign.
func (m *Simulation) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Display name of the attack simulation and training campaign. Supports $filter and $orderby.
func (m *Simulation) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetDurationInDays sets the durationInDays property value. Simulation duration in days.
func (m *Simulation) SetDurationInDays(value *int32)() {
    err := m.GetBackingStore().Set("durationInDays", value)
    if err != nil {
        panic(err)
    }
}
// SetEndUserNotificationSetting sets the endUserNotificationSetting property value. Details about the end user notification setting.
func (m *Simulation) SetEndUserNotificationSetting(value EndUserNotificationSettingable)() {
    err := m.GetBackingStore().Set("endUserNotificationSetting", value)
    if err != nil {
        panic(err)
    }
}
// SetExcludedAccountTarget sets the excludedAccountTarget property value. Users excluded from the simulation.
func (m *Simulation) SetExcludedAccountTarget(value AccountTargetContentable)() {
    err := m.GetBackingStore().Set("excludedAccountTarget", value)
    if err != nil {
        panic(err)
    }
}
// SetIncludedAccountTarget sets the includedAccountTarget property value. Users targeted in the simulation.
func (m *Simulation) SetIncludedAccountTarget(value AccountTargetContentable)() {
    err := m.GetBackingStore().Set("includedAccountTarget", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAutomated sets the isAutomated property value. Flag that represents if the attack simulation and training campaign was created from a simulation automation flow. Supports $filter and $orderby.
func (m *Simulation) SetIsAutomated(value *bool)() {
    err := m.GetBackingStore().Set("isAutomated", value)
    if err != nil {
        panic(err)
    }
}
// SetLandingPage sets the landingPage property value. The landing page associated with a simulation during its creation.
func (m *Simulation) SetLandingPage(value LandingPageable)() {
    err := m.GetBackingStore().Set("landingPage", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. Identity of the user who most recently modified the attack simulation and training campaign.
func (m *Simulation) SetLastModifiedBy(value EmailIdentityable)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. Date and time of the most recent modification of the attack simulation and training campaign.
func (m *Simulation) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLaunchDateTime sets the launchDateTime property value. Date and time of the launch/start of the attack simulation and training campaign. Supports $filter and $orderby.
func (m *Simulation) SetLaunchDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("launchDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLoginPage sets the loginPage property value. The login page associated with a simulation during its creation.
func (m *Simulation) SetLoginPage(value LoginPageable)() {
    err := m.GetBackingStore().Set("loginPage", value)
    if err != nil {
        panic(err)
    }
}
// SetOAuthConsentAppDetail sets the oAuthConsentAppDetail property value. OAuth app details for the OAuth technique.
func (m *Simulation) SetOAuthConsentAppDetail(value OAuthConsentAppDetailable)() {
    err := m.GetBackingStore().Set("oAuthConsentAppDetail", value)
    if err != nil {
        panic(err)
    }
}
// SetPayload sets the payload property value. The payload associated with a simulation during its creation.
func (m *Simulation) SetPayload(value Payloadable)() {
    err := m.GetBackingStore().Set("payload", value)
    if err != nil {
        panic(err)
    }
}
// SetPayloadDeliveryPlatform sets the payloadDeliveryPlatform property value. Method of delivery of the phishing payload used in the attack simulation and training campaign. Possible values are: unknown, sms, email, teams, unknownFutureValue.
func (m *Simulation) SetPayloadDeliveryPlatform(value *PayloadDeliveryPlatform)() {
    err := m.GetBackingStore().Set("payloadDeliveryPlatform", value)
    if err != nil {
        panic(err)
    }
}
// SetReport sets the report property value. Report of the attack simulation and training campaign.
func (m *Simulation) SetReport(value SimulationReportable)() {
    err := m.GetBackingStore().Set("report", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. Status of the attack simulation and training campaign. Supports $filter and $orderby. Possible values are: unknown, draft, running, scheduled, succeeded, failed, cancelled, excluded, unknownFutureValue.
func (m *Simulation) SetStatus(value *SimulationStatus)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
// SetTrainingSetting sets the trainingSetting property value. Details about the training settings for a simulation.
func (m *Simulation) SetTrainingSetting(value TrainingSettingable)() {
    err := m.GetBackingStore().Set("trainingSetting", value)
    if err != nil {
        panic(err)
    }
}
type Simulationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttackTechnique()(*SimulationAttackTechnique)
    GetAttackType()(*SimulationAttackType)
    GetAutomationId()(*string)
    GetCompletionDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetCreatedBy()(EmailIdentityable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetDurationInDays()(*int32)
    GetEndUserNotificationSetting()(EndUserNotificationSettingable)
    GetExcludedAccountTarget()(AccountTargetContentable)
    GetIncludedAccountTarget()(AccountTargetContentable)
    GetIsAutomated()(*bool)
    GetLandingPage()(LandingPageable)
    GetLastModifiedBy()(EmailIdentityable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLaunchDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLoginPage()(LoginPageable)
    GetOAuthConsentAppDetail()(OAuthConsentAppDetailable)
    GetPayload()(Payloadable)
    GetPayloadDeliveryPlatform()(*PayloadDeliveryPlatform)
    GetReport()(SimulationReportable)
    GetStatus()(*SimulationStatus)
    GetTrainingSetting()(TrainingSettingable)
    SetAttackTechnique(value *SimulationAttackTechnique)()
    SetAttackType(value *SimulationAttackType)()
    SetAutomationId(value *string)()
    SetCompletionDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetCreatedBy(value EmailIdentityable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetDurationInDays(value *int32)()
    SetEndUserNotificationSetting(value EndUserNotificationSettingable)()
    SetExcludedAccountTarget(value AccountTargetContentable)()
    SetIncludedAccountTarget(value AccountTargetContentable)()
    SetIsAutomated(value *bool)()
    SetLandingPage(value LandingPageable)()
    SetLastModifiedBy(value EmailIdentityable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLaunchDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLoginPage(value LoginPageable)()
    SetOAuthConsentAppDetail(value OAuthConsentAppDetailable)()
    SetPayload(value Payloadable)()
    SetPayloadDeliveryPlatform(value *PayloadDeliveryPlatform)()
    SetReport(value SimulationReportable)()
    SetStatus(value *SimulationStatus)()
    SetTrainingSetting(value TrainingSettingable)()
}
