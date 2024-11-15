package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/duration"
	"github.com/linode/linodego/internal/parseabletime"
)

// Event represents an action taken on the Account.
type Event struct {
	// The unique ID of this Event.
	ID int `json:"id"`

	// Current status of the Event, Enum: "failed" "finished" "notification" "scheduled" "started"
	Status EventStatus `json:"status"`

	// The action that caused this Event. New actions may be added in the future.
	Action EventAction `json:"action"`

	// A percentage estimating the amount of time remaining for an Event. Returns null for notification events.
	PercentComplete int `json:"percent_complete"`

	// The rate of completion of the Event. Only some Events will return rate; for example, migration and resize Events.
	Rate *string `json:"rate"`

	// If this Event has been read.
	Read bool `json:"read"`

	// If this Event has been seen.
	Seen bool `json:"seen"`

	// The estimated time remaining until the completion of this Event. This value is only returned for in-progress events.
	TimeRemaining *int `json:"-"`

	// The username of the User who caused the Event.
	Username string `json:"username"`

	// Detailed information about the Event's entity, including ID, type, label, and URL used to access it.
	Entity *EventEntity `json:"entity"`

	// Detailed information about the Event's secondary or related entity, including ID, type, label, and URL used to access it.
	SecondaryEntity *EventEntity `json:"secondary_entity"`

	// When this Event was created.
	Created *time.Time `json:"-"`
}

// EventAction constants start with Action and include all known Linode API Event Actions.
type EventAction string

// EventAction constants represent the actions that cause an Event. New actions may be added in the future.
const (
	ActionAccountUpdate                           EventAction = "account_update"
	ActionAccountSettingsUpdate                   EventAction = "account_settings_update"
	ActionBackupsEnable                           EventAction = "backups_enable"
	ActionBackupsCancel                           EventAction = "backups_cancel"
	ActionBackupsRestore                          EventAction = "backups_restore"
	ActionCommunityQuestionReply                  EventAction = "community_question_reply"
	ActionCommunityLike                           EventAction = "community_like"
	ActionCreditCardUpdated                       EventAction = "credit_card_updated"
	ActionDatabaseCreate                          EventAction = "database_create"
	ActionDatabaseDegraded                        EventAction = "database_degraded"
	ActionDatabaseDelete                          EventAction = "database_delete"
	ActionDatabaseFailed                          EventAction = "database_failed"
	ActionDatabaseUpdate                          EventAction = "database_update"
	ActionDatabaseCreateFailed                    EventAction = "database_create_failed"
	ActionDatabaseUpdateFailed                    EventAction = "database_update_failed"
	ActionDatabaseBackupCreate                    EventAction = "database_backup_create"
	ActionDatabaseBackupRestore                   EventAction = "database_backup_restore"
	ActionDatabaseCredentialsReset                EventAction = "database_credentials_reset"
	ActionDiskCreate                              EventAction = "disk_create"
	ActionDiskDelete                              EventAction = "disk_delete"
	ActionDiskUpdate                              EventAction = "disk_update"
	ActionDiskDuplicate                           EventAction = "disk_duplicate"
	ActionDiskImagize                             EventAction = "disk_imagize"
	ActionDiskResize                              EventAction = "disk_resize"
	ActionDNSRecordCreate                         EventAction = "dns_record_create"
	ActionDNSRecordDelete                         EventAction = "dns_record_delete"
	ActionDNSRecordUpdate                         EventAction = "dns_record_update"
	ActionDNSZoneCreate                           EventAction = "dns_zone_create"
	ActionDNSZoneDelete                           EventAction = "dns_zone_delete"
	ActionDNSZoneUpdate                           EventAction = "dns_zone_update"
	ActionDNSZoneImport                           EventAction = "dns_zone_import"
	ActionEntityTransferAccept                    EventAction = "entity_transfer_accept"
	ActionEntityTransferCancel                    EventAction = "entity_transfer_cancel"
	ActionEntityTransferCreate                    EventAction = "entity_transfer_create"
	ActionEntityTransferFail                      EventAction = "entity_transfer_fail"
	ActionEntityTransferStale                     EventAction = "entity_transfer_stale"
	ActionFirewallCreate                          EventAction = "firewall_create"
	ActionFirewallDelete                          EventAction = "firewall_delete"
	ActionFirewallDisable                         EventAction = "firewall_disable"
	ActionFirewallEnable                          EventAction = "firewall_enable"
	ActionFirewallUpdate                          EventAction = "firewall_update"
	ActionFirewallDeviceAdd                       EventAction = "firewall_device_add"
	ActionFirewallDeviceRemove                    EventAction = "firewall_device_remove"
	ActionHostReboot                              EventAction = "host_reboot"
	ActionImageDelete                             EventAction = "image_delete"
	ActionImageUpdate                             EventAction = "image_update"
	ActionImageUpload                             EventAction = "image_upload"
	ActionIPAddressUpdate                         EventAction = "ipaddress_update"
	ActionLassieReboot                            EventAction = "lassie_reboot"
	ActionLinodeAddIP                             EventAction = "linode_addip"
	ActionLinodeBoot                              EventAction = "linode_boot"
	ActionLinodeClone                             EventAction = "linode_clone"
	ActionLinodeCreate                            EventAction = "linode_create"
	ActionLinodeDelete                            EventAction = "linode_delete"
	ActionLinodeUpdate                            EventAction = "linode_update"
	ActionLinodeDeleteIP                          EventAction = "linode_deleteip"
	ActionLinodeMigrate                           EventAction = "linode_migrate"
	ActionLinodeMigrateDatacenter                 EventAction = "linode_migrate_datacenter"
	ActionLinodeMigrateDatacenterCreate           EventAction = "linode_migrate_datacenter_create"
	ActionLinodeMutate                            EventAction = "linode_mutate"
	ActionLinodeMutateCreate                      EventAction = "linode_mutate_create"
	ActionLinodeReboot                            EventAction = "linode_reboot"
	ActionLinodeRebuild                           EventAction = "linode_rebuild"
	ActionLinodeResize                            EventAction = "linode_resize"
	ActionLinodeResizeCreate                      EventAction = "linode_resize_create"
	ActionLinodeShutdown                          EventAction = "linode_shutdown"
	ActionLinodeSnapshot                          EventAction = "linode_snapshot"
	ActionLinodeConfigCreate                      EventAction = "linode_config_create"
	ActionLinodeConfigDelete                      EventAction = "linode_config_delete"
	ActionLinodeConfigUpdate                      EventAction = "linode_config_update"
	ActionLishBoot                                EventAction = "lish_boot"
	ActionLKENodeCreate                           EventAction = "lke_node_create"
	ActionLKEControlPlaneACLCreate                EventAction = "lke_control_plane_acl_create"
	ActionLKEControlPlaneACLUpdate                EventAction = "lke_control_plane_acl_update"
	ActionLKEControlPlaneACLDelete                EventAction = "lke_control_plane_acl_delete"
	ActionLongviewClientCreate                    EventAction = "longviewclient_create"
	ActionLongviewClientDelete                    EventAction = "longviewclient_delete"
	ActionLongviewClientUpdate                    EventAction = "longviewclient_update"
	ActionManagedDisabled                         EventAction = "managed_disabled"
	ActionManagedEnabled                          EventAction = "managed_enabled"
	ActionManagedServiceCreate                    EventAction = "managed_service_create"
	ActionManagedServiceDelete                    EventAction = "managed_service_delete"
	ActionNodebalancerCreate                      EventAction = "nodebalancer_create"
	ActionNodebalancerDelete                      EventAction = "nodebalancer_delete"
	ActionNodebalancerUpdate                      EventAction = "nodebalancer_update"
	ActionNodebalancerConfigCreate                EventAction = "nodebalancer_config_create"
	ActionNodebalancerConfigDelete                EventAction = "nodebalancer_config_delete"
	ActionNodebalancerConfigUpdate                EventAction = "nodebalancer_config_update"
	ActionNodebalancerFirewallModificationSuccess EventAction = "nodebalancer_firewall_modification_success"
	ActionNodebalancerFirewallModificationFailed  EventAction = "nodebalancer_firewall_modification_failed"
	ActionNodebalancerNodeCreate                  EventAction = "nodebalancer_node_create"
	ActionNodebalancerNodeDelete                  EventAction = "nodebalancer_node_delete"
	ActionNodebalancerNodeUpdate                  EventAction = "nodebalancer_node_update"
	ActionOAuthClientCreate                       EventAction = "oauth_client_create"
	ActionOAuthClientDelete                       EventAction = "oauth_client_delete"
	ActionOAuthClientSecretReset                  EventAction = "oauth_client_secret_reset" //#nosec G101
	ActionOAuthClientUpdate                       EventAction = "oauth_client_update"
	ActionOBJAccessKeyCreate                      EventAction = "obj_access_key_create"
	ActionOBJAccessKeyDelete                      EventAction = "obj_access_key_delete"
	ActionOBJAccessKeyUpdate                      EventAction = "obj_access_key_update"
	ActionPaymentMethodAdd                        EventAction = "payment_method_add"
	ActionPaymentSubmitted                        EventAction = "payment_submitted"
	ActionPasswordReset                           EventAction = "password_reset"
	ActionPlacementGroupCreate                    EventAction = "placement_group_create"
	ActionPlacementGroupUpdate                    EventAction = "placement_group_update"
	ActionPlacementGroupDelete                    EventAction = "placement_group_delete"
	ActionPlacementGroupAssign                    EventAction = "placement_group_assign"
	ActionPlacementGroupUnassign                  EventAction = "placement_group_unassign"
	ActionPlacementGroupBecameNonCompliant        EventAction = "placement_group_became_non_compliant"
	ActionPlacementGroupBecameCompliant           EventAction = "placement_group_became_compliant"
	ActionProfileUpdate                           EventAction = "profile_update"
	ActionStackScriptCreate                       EventAction = "stackscript_create"
	ActionStackScriptDelete                       EventAction = "stackscript_delete"
	ActionStackScriptUpdate                       EventAction = "stackscript_update"
	ActionStackScriptPublicize                    EventAction = "stackscript_publicize"
	ActionStackScriptRevise                       EventAction = "stackscript_revise"
	ActionTaxIDInvalid                            EventAction = "tax_id_invalid"
	ActionTagCreate                               EventAction = "tag_create"
	ActionTagDelete                               EventAction = "tag_delete"
	ActionTFADisabled                             EventAction = "tfa_disabled"
	ActionTFAEnabled                              EventAction = "tfa_enabled"
	ActionTicketAttachmentUpload                  EventAction = "ticket_attachment_upload"
	ActionTicketCreate                            EventAction = "ticket_create"
	ActionTicketUpdate                            EventAction = "ticket_update"
	ActionTokenCreate                             EventAction = "token_create"
	ActionTokenDelete                             EventAction = "token_delete"
	ActionTokenUpdate                             EventAction = "token_update"
	ActionUserCreate                              EventAction = "user_create"
	ActionUserDelete                              EventAction = "user_delete"
	ActionUserUpdate                              EventAction = "user_update"
	ActionUserSSHKeyAdd                           EventAction = "user_ssh_key_add"
	ActionUserSSHKeyDelete                        EventAction = "user_ssh_key_delete"
	ActionUserSSHKeyUpdate                        EventAction = "user_ssh_key_update"
	ActionVLANAttach                              EventAction = "vlan_attach"
	ActionVLANDetach                              EventAction = "vlan_detach"
	ActionVolumeAttach                            EventAction = "volume_attach"
	ActionVolumeClone                             EventAction = "volume_clone"
	ActionVolumeCreate                            EventAction = "volume_create"
	ActionVolumeDelete                            EventAction = "volume_delete"
	ActionVolumeUpdate                            EventAction = "volume_update"
	ActionVolumeDetach                            EventAction = "volume_detach"
	ActionVolumeResize                            EventAction = "volume_resize"
	ActionVPCCreate                               EventAction = "vpc_create"
	ActionVPCDelete                               EventAction = "vpc_delete"
	ActionVPCUpdate                               EventAction = "vpc_update"
	ActionVPCSubnetCreate                         EventAction = "subnet_create"
	ActionVPCSubnetDelete                         EventAction = "subnet_delete"
	ActionVPCSubnetUpdate                         EventAction = "subnet_update"

	// Deprecated: incorrect spelling,
	// to be removed in the next major version release.
	ActionVolumeDelte EventAction = "volume_delete"

	// Deprecated: incorrect spelling,
	// to be removed in the next major version
	ActionCreateCardUpdated = ActionCreditCardUpdated
)

// EntityType constants start with Entity and include Linode API Event Entity Types
type EntityType string

// EntityType contants are the entities an Event can be related to.
const (
	EntityAccount        EntityType = "account"
	EntityBackups        EntityType = "backups"
	EntityCommunity      EntityType = "community"
	EntityDatabase       EntityType = "database"
	EntityDisk           EntityType = "disk"
	EntityDomain         EntityType = "domain"
	EntityTransfer       EntityType = "entity_transfer"
	EntityFirewall       EntityType = "firewall"
	EntityImage          EntityType = "image"
	EntityIPAddress      EntityType = "ipaddress"
	EntityLinode         EntityType = "linode"
	EntityLongview       EntityType = "longview"
	EntityManagedService EntityType = "managed_service"
	EntityNodebalancer   EntityType = "nodebalancer"
	EntityOAuthClient    EntityType = "oauth_client"
	EntityPlacementGroup EntityType = "placement_group"
	EntityProfile        EntityType = "profile"
	EntityStackscript    EntityType = "stackscript"
	EntityTag            EntityType = "tag"
	EntityTicket         EntityType = "ticket"
	EntityToken          EntityType = "token"
	EntityUser           EntityType = "user"
	EntityUserSSHKey     EntityType = "user_ssh_key"
	EntityVolume         EntityType = "volume"
	EntityVPC            EntityType = "vpc"
	EntityVPCSubnet      EntityType = "subnet"
)

// EventStatus constants start with Event and include Linode API Event Status values
type EventStatus string

// EventStatus constants reflect the current status of an Event
const (
	EventFailed       EventStatus = "failed"
	EventFinished     EventStatus = "finished"
	EventNotification EventStatus = "notification"
	EventScheduled    EventStatus = "scheduled"
	EventStarted      EventStatus = "started"
)

// EventEntity provides detailed information about the Event's
// associated entity, including ID, Type, Label, and a URL that
// can be used to access it.
type EventEntity struct {
	// ID may be a string or int, it depends on the EntityType
	ID     any        `json:"id"`
	Label  string     `json:"label"`
	Type   EntityType `json:"type"`
	Status string     `json:"status"`
	URL    string     `json:"url"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *Event) UnmarshalJSON(b []byte) error {
	type Mask Event

	p := struct {
		*Mask
		Created       *parseabletime.ParseableTime `json:"created"`
		TimeRemaining json.RawMessage              `json:"time_remaining"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.TimeRemaining = duration.UnmarshalTimeRemaining(p.TimeRemaining)

	return nil
}

// ListEvents gets a collection of Event objects representing actions taken
// on the Account. The Events returned depend on the token grants and the grants
// of the associated user.
func (c *Client) ListEvents(ctx context.Context, opts *ListOptions) ([]Event, error) {
	response, err := getPaginatedResults[Event](ctx, c, "account/events", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetEvent gets the Event with the Event ID
func (c *Client) GetEvent(ctx context.Context, eventID int) (*Event, error) {
	e := formatAPIPath("account/events/%d", eventID)
	response, err := doGETRequest[Event](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// MarkEventRead marks a single Event as read.
func (c *Client) MarkEventRead(ctx context.Context, event *Event) error {
	e := formatAPIPath("account/events/%d/read", event.ID)
	_, err := doPOSTRequest[Event](ctx, c, e, []any{})
	return err
}

// MarkEventsSeen marks all Events up to and including this Event by ID as seen.
func (c *Client) MarkEventsSeen(ctx context.Context, event *Event) error {
	e := formatAPIPath("account/events/%d/seen", event.ID)
	_, err := doPOSTRequest[Event](ctx, c, e, []any{})
	return err
}
