package models

import (
	i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22 "github.com/google/uuid"
	i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
	ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
	i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
)

// ChangeNotification
type ChangeNotification struct {
	// Stores model information.
	backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}

// NewChangeNotification instantiates a new changeNotification and sets the default values.
func NewChangeNotification() *ChangeNotification {
	m := &ChangeNotification{}
	m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance()
	m.SetAdditionalData(make(map[string]any))
	return m
}

// CreateChangeNotificationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateChangeNotificationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) (i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
	return NewChangeNotification(), nil
}

// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ChangeNotification) GetAdditionalData() map[string]any {
	val, err := m.backingStore.Get("additionalData")
	if err != nil {
		panic(err)
	}
	if val == nil {
		var value = make(map[string]any)
		m.SetAdditionalData(value)
	}
	return val.(map[string]any)
}

// GetBackingStore gets the backingStore property value. Stores model information.
func (m *ChangeNotification) GetBackingStore() ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore {
	return m.backingStore
}

// GetChangeType gets the changeType property value. The changeType property
func (m *ChangeNotification) GetChangeType() *ChangeType {
	val, err := m.GetBackingStore().Get("changeType")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(*ChangeType)
	}
	return nil
}

// GetClientState gets the clientState property value. Value of the clientState property sent in the subscription request (if any). The maximum length is 255 characters. The client can check whether the change notification came from the service by comparing the values of the clientState property. The value of the clientState property sent with the subscription is compared with the value of the clientState property received with each change notification. Optional.
func (m *ChangeNotification) GetClientState() *string {
	val, err := m.GetBackingStore().Get("clientState")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(*string)
	}
	return nil
}

// GetEncryptedContent gets the encryptedContent property value. (Preview) Encrypted content attached with the change notification. Only provided if encryptionCertificate and includeResourceData were defined during the subscription request and if the resource supports it. Optional.
func (m *ChangeNotification) GetEncryptedContent() ChangeNotificationEncryptedContentable {
	val, err := m.GetBackingStore().Get("encryptedContent")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(ChangeNotificationEncryptedContentable)
	}
	return nil
}

// GetFieldDeserializers the deserialization information for the current model
func (m *ChangeNotification) GetFieldDeserializers() map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
	res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error)
	res["changeType"] = func(n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
		val, err := n.GetEnumValue(ParseChangeType)
		if err != nil {
			return err
		}
		if val != nil {
			m.SetChangeType(val.(*ChangeType))
		}
		return nil
	}
	res["clientState"] = func(n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetClientState(val)
		}
		return nil
	}
	res["encryptedContent"] = func(n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
		val, err := n.GetObjectValue(CreateChangeNotificationEncryptedContentFromDiscriminatorValue)
		if err != nil {
			return err
		}
		if val != nil {
			m.SetEncryptedContent(val.(ChangeNotificationEncryptedContentable))
		}
		return nil
	}
	res["id"] = func(n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetId(val)
		}
		return nil
	}
	res["lifecycleEvent"] = func(n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
		val, err := n.GetEnumValue(ParseLifecycleEventType)
		if err != nil {
			return err
		}
		if val != nil {
			m.SetLifecycleEvent(val.(*LifecycleEventType))
		}
		return nil
	}
	res["@odata.type"] = func(n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetOdataType(val)
		}
		return nil
	}
	res["resource"] = func(n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
		val, err := n.GetStringValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetResource(val)
		}
		return nil
	}
	res["resourceData"] = func(n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
		val, err := n.GetObjectValue(CreateResourceDataFromDiscriminatorValue)
		if err != nil {
			return err
		}
		if val != nil {
			m.SetResourceData(val.(ResourceDataable))
		}
		return nil
	}
	res["subscriptionExpirationDateTime"] = func(n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
		val, err := n.GetTimeValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetSubscriptionExpirationDateTime(val)
		}
		return nil
	}
	res["subscriptionId"] = func(n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
		val, err := n.GetUUIDValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetSubscriptionId(val)
		}
		return nil
	}
	res["tenantId"] = func(n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
		val, err := n.GetUUIDValue()
		if err != nil {
			return err
		}
		if val != nil {
			m.SetTenantId(val)
		}
		return nil
	}
	return res
}

// GetId gets the id property value. Unique ID for the notification. Optional.
func (m *ChangeNotification) GetId() *string {
	val, err := m.GetBackingStore().Get("id")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(*string)
	}
	return nil
}

// GetLifecycleEvent gets the lifecycleEvent property value. The type of lifecycle notification if the current notification is a lifecycle notification. Optional. Supported values are missed, subscriptionRemoved, reauthorizationRequired. Optional.
func (m *ChangeNotification) GetLifecycleEvent() *LifecycleEventType {
	val, err := m.GetBackingStore().Get("lifecycleEvent")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(*LifecycleEventType)
	}
	return nil
}

// GetOdataType gets the @odata.type property value. The OdataType property
func (m *ChangeNotification) GetOdataType() *string {
	val, err := m.GetBackingStore().Get("odataType")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(*string)
	}
	return nil
}

// GetResource gets the resource property value. The URI of the resource that emitted the change notification relative to https://graph.microsoft.com. Required.
func (m *ChangeNotification) GetResource() *string {
	val, err := m.GetBackingStore().Get("resource")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(*string)
	}
	return nil
}

// GetResourceData gets the resourceData property value. The content of this property depends on the type of resource being subscribed to. Optional.
func (m *ChangeNotification) GetResourceData() ResourceDataable {
	val, err := m.GetBackingStore().Get("resourceData")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(ResourceDataable)
	}
	return nil
}

// GetSubscriptionExpirationDateTime gets the subscriptionExpirationDateTime property value. The expiration time for the subscription. Required.
func (m *ChangeNotification) GetSubscriptionExpirationDateTime() *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time {
	val, err := m.GetBackingStore().Get("subscriptionExpirationDateTime")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
	}
	return nil
}

// GetSubscriptionId gets the subscriptionId property value. The unique identifier of the subscription that generated the notification.Required.
func (m *ChangeNotification) GetSubscriptionId() *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID {
	val, err := m.GetBackingStore().Get("subscriptionId")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
	}
	return nil
}

// GetTenantId gets the tenantId property value. The unique identifier of the tenant from which the change notification originated. Required.
func (m *ChangeNotification) GetTenantId() *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID {
	val, err := m.GetBackingStore().Get("tenantId")
	if err != nil {
		panic(err)
	}
	if val != nil {
		return val.(*i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
	}
	return nil
}

// Serialize serializes information the current object
func (m *ChangeNotification) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter) error {
	if m.GetChangeType() != nil {
		cast := (*m.GetChangeType()).String()
		err := writer.WriteStringValue("changeType", &cast)
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteStringValue("clientState", m.GetClientState())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteObjectValue("encryptedContent", m.GetEncryptedContent())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteStringValue("id", m.GetId())
		if err != nil {
			return err
		}
	}
	if m.GetLifecycleEvent() != nil {
		cast := (*m.GetLifecycleEvent()).String()
		err := writer.WriteStringValue("lifecycleEvent", &cast)
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
		err := writer.WriteStringValue("resource", m.GetResource())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteObjectValue("resourceData", m.GetResourceData())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteTimeValue("subscriptionExpirationDateTime", m.GetSubscriptionExpirationDateTime())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteUUIDValue("subscriptionId", m.GetSubscriptionId())
		if err != nil {
			return err
		}
	}
	{
		err := writer.WriteUUIDValue("tenantId", m.GetTenantId())
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

// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ChangeNotification) SetAdditionalData(value map[string]any) {
	err := m.GetBackingStore().Set("additionalData", value)
	if err != nil {
		panic(err)
	}
}

// SetBackingStore sets the backingStore property value. Stores model information.
func (m *ChangeNotification) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
	m.backingStore = value
}

// SetChangeType sets the changeType property value. The changeType property
func (m *ChangeNotification) SetChangeType(value *ChangeType) {
	err := m.GetBackingStore().Set("changeType", value)
	if err != nil {
		panic(err)
	}
}

// SetClientState sets the clientState property value. Value of the clientState property sent in the subscription request (if any). The maximum length is 255 characters. The client can check whether the change notification came from the service by comparing the values of the clientState property. The value of the clientState property sent with the subscription is compared with the value of the clientState property received with each change notification. Optional.
func (m *ChangeNotification) SetClientState(value *string) {
	err := m.GetBackingStore().Set("clientState", value)
	if err != nil {
		panic(err)
	}
}

// SetEncryptedContent sets the encryptedContent property value. (Preview) Encrypted content attached with the change notification. Only provided if encryptionCertificate and includeResourceData were defined during the subscription request and if the resource supports it. Optional.
func (m *ChangeNotification) SetEncryptedContent(value ChangeNotificationEncryptedContentable) {
	err := m.GetBackingStore().Set("encryptedContent", value)
	if err != nil {
		panic(err)
	}
}

// SetId sets the id property value. Unique ID for the notification. Optional.
func (m *ChangeNotification) SetId(value *string) {
	err := m.GetBackingStore().Set("id", value)
	if err != nil {
		panic(err)
	}
}

// SetLifecycleEvent sets the lifecycleEvent property value. The type of lifecycle notification if the current notification is a lifecycle notification. Optional. Supported values are missed, subscriptionRemoved, reauthorizationRequired. Optional.
func (m *ChangeNotification) SetLifecycleEvent(value *LifecycleEventType) {
	err := m.GetBackingStore().Set("lifecycleEvent", value)
	if err != nil {
		panic(err)
	}
}

// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ChangeNotification) SetOdataType(value *string) {
	err := m.GetBackingStore().Set("odataType", value)
	if err != nil {
		panic(err)
	}
}

// SetResource sets the resource property value. The URI of the resource that emitted the change notification relative to https://graph.microsoft.com. Required.
func (m *ChangeNotification) SetResource(value *string) {
	err := m.GetBackingStore().Set("resource", value)
	if err != nil {
		panic(err)
	}
}

// SetResourceData sets the resourceData property value. The content of this property depends on the type of resource being subscribed to. Optional.
func (m *ChangeNotification) SetResourceData(value ResourceDataable) {
	err := m.GetBackingStore().Set("resourceData", value)
	if err != nil {
		panic(err)
	}
}

// SetSubscriptionExpirationDateTime sets the subscriptionExpirationDateTime property value. The expiration time for the subscription. Required.
func (m *ChangeNotification) SetSubscriptionExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
	err := m.GetBackingStore().Set("subscriptionExpirationDateTime", value)
	if err != nil {
		panic(err)
	}
}

// SetSubscriptionId sets the subscriptionId property value. The unique identifier of the subscription that generated the notification.Required.
func (m *ChangeNotification) SetSubscriptionId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
	err := m.GetBackingStore().Set("subscriptionId", value)
	if err != nil {
		panic(err)
	}
}

// SetTenantId sets the tenantId property value. The unique identifier of the tenant from which the change notification originated. Required.
func (m *ChangeNotification) SetTenantId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID) {
	err := m.GetBackingStore().Set("tenantId", value)
	if err != nil {
		panic(err)
	}
}

// ChangeNotificationable
type ChangeNotificationable interface {
	i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
	ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
	i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
	GetBackingStore() ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
	GetChangeType() *ChangeType
	GetClientState() *string
	GetEncryptedContent() ChangeNotificationEncryptedContentable
	GetId() *string
	GetLifecycleEvent() *LifecycleEventType
	GetOdataType() *string
	GetResource() *string
	GetResourceData() ResourceDataable
	GetSubscriptionExpirationDateTime() *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time
	GetSubscriptionId() *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID
	GetTenantId() *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID
	SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
	SetChangeType(value *ChangeType)
	SetClientState(value *string)
	SetEncryptedContent(value ChangeNotificationEncryptedContentable)
	SetId(value *string)
	SetLifecycleEvent(value *LifecycleEventType)
	SetOdataType(value *string)
	SetResource(value *string)
	SetResourceData(value ResourceDataable)
	SetSubscriptionExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
	SetSubscriptionId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
	SetTenantId(value *i561e97a8befe7661a44c8f54600992b4207a3a0cf6770e5559949bc276de2e22.UUID)
}
