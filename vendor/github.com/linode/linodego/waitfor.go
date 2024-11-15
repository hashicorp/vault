package linodego

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var englishTitle = cases.Title(language.English)

type EventPoller struct {
	EntityID   any
	EntityType EntityType

	// Type is excluded here because it is implicitly determined
	// by the event action.
	SecondaryEntityID any

	Action EventAction

	client         Client
	previousEvents map[int]bool
}

// WaitForInstanceStatus waits for the Linode instance to reach the desired state
// before returning. It will timeout with an error after timeoutSeconds.
func (client Client) WaitForInstanceStatus(ctx context.Context, instanceID int, status InstanceStatus, timeoutSeconds int) (*Instance, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			instance, err := client.GetInstance(ctx, instanceID)
			if err != nil {
				return instance, err
			}
			complete := (instance.Status == status)

			if complete {
				return instance, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("Error waiting for Instance %d status %s: %w", instanceID, status, ctx.Err())
		}
	}
}

// WaitForInstanceDiskStatus waits for the Linode instance disk to reach the desired state
// before returning. It will timeout with an error after timeoutSeconds.
func (client Client) WaitForInstanceDiskStatus(ctx context.Context, instanceID int, diskID int, status DiskStatus, timeoutSeconds int) (*InstanceDisk, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// GetInstanceDisk will 404 on newly created disks. use List instead.
			// disk, err := client.GetInstanceDisk(ctx, instanceID, diskID)
			disks, err := client.ListInstanceDisks(ctx, instanceID, nil)
			if err != nil {
				return nil, err
			}

			for _, disk := range disks {
				if disk.ID == diskID {
					complete := (disk.Status == status)
					if complete {
						return &disk, nil
					}

					break
				}
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("Error waiting for Instance %d Disk %d status %s: %w", instanceID, diskID, status, ctx.Err())
		}
	}
}

// WaitForVolumeStatus waits for the Volume to reach the desired state
// before returning. It will timeout with an error after timeoutSeconds.
func (client Client) WaitForVolumeStatus(ctx context.Context, volumeID int, status VolumeStatus, timeoutSeconds int) (*Volume, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			volume, err := client.GetVolume(ctx, volumeID)
			if err != nil {
				return volume, err
			}
			complete := (volume.Status == status)

			if complete {
				return volume, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("Error waiting for Volume %d status %s: %w", volumeID, status, ctx.Err())
		}
	}
}

// WaitForSnapshotStatus waits for the Snapshot to reach the desired state
// before returning. It will timeout with an error after timeoutSeconds.
func (client Client) WaitForSnapshotStatus(ctx context.Context, instanceID int, snapshotID int, status InstanceSnapshotStatus, timeoutSeconds int) (*InstanceSnapshot, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			snapshot, err := client.GetInstanceSnapshot(ctx, instanceID, snapshotID)
			if err != nil {
				return snapshot, err
			}
			complete := (snapshot.Status == status)

			if complete {
				return snapshot, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("Error waiting for Instance %d Snapshot %d status %s: %w", instanceID, snapshotID, status, ctx.Err())
		}
	}
}

// WaitForVolumeLinodeID waits for the Volume to match the desired LinodeID
// before returning. An active Instance will not immediately attach or detach a volume, so
// the LinodeID must be polled to determine volume readiness from the API.
// WaitForVolumeLinodeID will timeout with an error after timeoutSeconds.
func (client Client) WaitForVolumeLinodeID(ctx context.Context, volumeID int, linodeID *int, timeoutSeconds int) (*Volume, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			volume, err := client.GetVolume(ctx, volumeID)
			if err != nil {
				return volume, err
			}

			switch {
			case linodeID == nil && volume.LinodeID == nil:
				return volume, nil
			case linodeID == nil || volume.LinodeID == nil:
				// continue waiting
			case *volume.LinodeID == *linodeID:
				return volume, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("Error waiting for Volume %d to have Instance %v: %w", volumeID, linodeID, ctx.Err())
		}
	}
}

// WaitForLKEClusterStatus waits for the LKECluster to reach the desired state
// before returning. It will timeout with an error after timeoutSeconds.
func (client Client) WaitForLKEClusterStatus(ctx context.Context, clusterID int, status LKEClusterStatus, timeoutSeconds int) (*LKECluster, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cluster, err := client.GetLKECluster(ctx, clusterID)
			if err != nil {
				return cluster, err
			}
			complete := (cluster.Status == status)

			if complete {
				return cluster, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("Error waiting for Cluster %d status %s: %w", clusterID, status, ctx.Err())
		}
	}
}

// LKEClusterPollOptions configures polls against LKE Clusters.
type LKEClusterPollOptions struct {
	// Retry will cause the Poll to ignore interimittent errors
	Retry bool

	// TimeoutSeconds is the number of Seconds to wait for the poll to succeed
	// before exiting.
	TimeoutSeconds int

	// TansportWrapper allows adding a transport middleware function that will
	// wrap the LKE Cluster client's underlying http.RoundTripper.
	TransportWrapper func(http.RoundTripper) http.RoundTripper
}

type ClusterConditionOptions struct {
	LKEClusterKubeconfig *LKEClusterKubeconfig
	TransportWrapper     func(http.RoundTripper) http.RoundTripper
}

// ClusterConditionFunc represents a function that tests a condition against an LKE cluster,
// returns true if the condition has been reached, false if it has not yet been reached.
type ClusterConditionFunc func(context.Context, ClusterConditionOptions) (bool, error)

// WaitForLKEClusterConditions waits for the given LKE conditions to be true
func (client Client) WaitForLKEClusterConditions(
	ctx context.Context,
	clusterID int,
	options LKEClusterPollOptions,
	conditions ...ClusterConditionFunc,
) error {
	ctx, cancel := context.WithCancel(ctx)
	if options.TimeoutSeconds != 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(options.TimeoutSeconds)*time.Second)
	}
	defer cancel()

	lkeKubeConfig, err := client.GetLKEClusterKubeconfig(ctx, clusterID)
	if err != nil {
		return fmt.Errorf("failed to get Kubeconfig for LKE cluster %d: %w", clusterID, err)
	}

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	conditionOptions := ClusterConditionOptions{LKEClusterKubeconfig: lkeKubeConfig, TransportWrapper: options.TransportWrapper}

	for _, condition := range conditions {
	ConditionSucceeded:
		for {
			select {
			case <-ticker.C:
				result, err := condition(ctx, conditionOptions)
				if err != nil {
					log.Printf("[WARN] Ignoring WaitForLKEClusterConditions conditional error: %s", err)
					if !options.Retry {
						return err
					}
				}

				if result {
					break ConditionSucceeded
				}

			case <-ctx.Done():
				return fmt.Errorf("Error waiting for cluster %d conditions: %w", clusterID, ctx.Err())
			}
		}
	}
	return nil
}

// WaitForEventFinished waits for an entity action to reach the 'finished' state
// before returning. It will timeout with an error after timeoutSeconds.
// If the event indicates a failure both the failed event and the error will be returned.
// nolint
func (client Client) WaitForEventFinished(
	ctx context.Context,
	id any,
	entityType EntityType,
	action EventAction,
	minStart time.Time,
	timeoutSeconds int,
) (*Event, error) {
	titledEntityType := englishTitle.String(string(entityType))
	filter := Filter{
		Order:   Descending,
		OrderBy: "created",
	}
	filter.AddField(Eq, "action", action)
	filter.AddField(Gte, "created", minStart.UTC().Format("2006-01-02T15:04:05"))

	// Optimistically restrict results to page 1.  We should remove this when more
	// precise filtering options exist.
	pages := 1

	// The API has limitted filtering support for Event ID and Event Type
	// Optimize the list, if possible
	switch entityType {
	case EntityDisk, EntityDatabase, EntityLinode, EntityDomain, EntityNodebalancer:
		// All of the filter supported types have int ids
		filterableEntityID, err := strconv.Atoi(fmt.Sprintf("%v", id))
		if err != nil {
			return nil, fmt.Errorf("error parsing Entity ID %q for optimized "+
				"WaitForEventFinished EventType %q: %w", id, entityType, err)
		}
		filter.AddField(Eq, "entity.id", filterableEntityID)
		filter.AddField(Eq, "entity.type", entityType)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	if deadline, ok := ctx.Deadline(); ok {
		duration := time.Until(deadline)
		log.Printf("[INFO] Waiting %d seconds for %s events since %v for %s %v", int(duration.Seconds()), action, minStart, titledEntityType, id)
	}

	ticker := time.NewTicker(client.pollInterval)

	// avoid repeating log messages
	nextLog := ""
	lastLog := ""
	lastEventID := 0

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if lastEventID > 0 {
				filter.AddField(Gte, "id", lastEventID)
			}

			filterStr, err := filter.MarshalJSON()
			if err != nil {
				return nil, err
			}

			listOptions := NewListOptions(pages, string(filterStr))

			events, err := client.ListEvents(ctx, listOptions)
			if err != nil {
				return nil, err
			}

			// If there are events for this instance + action, inspect them
			for _, event := range events {
				event := event

				if event.Entity == nil || event.Entity.Type != entityType {
					// log.Println("type mismatch", event.Entity.Type, entityType)
					continue
				}

				var entID string

				switch id := event.Entity.ID.(type) {
				case float64, float32:
					entID = fmt.Sprintf("%.f", id)
				case int:
					entID = strconv.Itoa(id)
				default:
					entID = fmt.Sprintf("%v", id)
				}

				var findID string
				switch id := id.(type) {
				case float64, float32:
					findID = fmt.Sprintf("%.f", id)
				case int:
					findID = strconv.Itoa(id)
				default:
					findID = fmt.Sprintf("%v", id)
				}

				if entID != findID {
					// log.Println("id mismatch", entID, findID)
					continue
				}

				if event.Created == nil {
					log.Printf("[WARN] event.Created is nil when API returned: %#+v", event.Created)
				}

				// This is the event we are looking for. Save our place.
				if lastEventID == 0 {
					lastEventID = event.ID
				}

				switch event.Status {
				case EventFailed:
					return &event, fmt.Errorf("%s %v action %s failed", titledEntityType, id, action)
				case EventFinished:
					log.Printf("[INFO] %s %v action %s is finished", titledEntityType, id, action)
					return &event, nil
				}

				nextLog = fmt.Sprintf("[INFO] %s %v action %s is %s", titledEntityType, id, action, event.Status)
			}

			// de-dupe logging statements
			if nextLog != lastLog {
				log.Print(nextLog)
				lastLog = nextLog
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("Error waiting for Event Status '%s' of %s %v action '%s': %w", EventFinished, titledEntityType, id, action, ctx.Err())
		}
	}
}

// WaitForImageStatus waits for the Image to reach the desired state
// before returning. It will timeout with an error after timeoutSeconds.
func (client Client) WaitForImageStatus(ctx context.Context, imageID string, status ImageStatus, timeoutSeconds int) (*Image, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			image, err := client.GetImage(ctx, imageID)
			if err != nil {
				return image, err
			}
			complete := image.Status == status

			if complete {
				return image, nil
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to wait for Image %s status %s: %w", imageID, status, ctx.Err())
		}
	}
}

// WaitForImageRegionStatus waits for an Image's replica to reach the desired state
// before returning.
func (client Client) WaitForImageRegionStatus(ctx context.Context, imageID, region string, status ImageRegionStatus) (*Image, error) {
	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			image, err := client.GetImage(ctx, imageID)
			if err != nil {
				return image, err
			}

			replicaIdx := slices.IndexFunc(
				image.Regions,
				func(r ImageRegion) bool {
					return r.Region == region
				},
			)

			// If no replica was found or the status doesn't match, try again
			if replicaIdx < 0 || image.Regions[replicaIdx].Status != status {
				continue
			}

			return image, nil

		case <-ctx.Done():
			return nil, fmt.Errorf("failed to wait for Image %s status %s: %w", imageID, status, ctx.Err())
		}
	}
}

// WaitForMySQLDatabaseBackup waits for the backup with the given label to be available.
func (client Client) WaitForMySQLDatabaseBackup(ctx context.Context, dbID int, label string, timeoutSeconds int) (*MySQLDatabaseBackup, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			backups, err := client.ListMySQLDatabaseBackups(ctx, dbID, nil)
			if err != nil {
				return nil, err
			}

			for _, backup := range backups {
				if backup.Label == label {
					return &backup, nil
				}
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to wait for backup %s: %w", label, ctx.Err())
		}
	}
}

// WaitForPostgresDatabaseBackup waits for the backup with the given label to be available.
func (client Client) WaitForPostgresDatabaseBackup(ctx context.Context, dbID int, label string, timeoutSeconds int) (*PostgresDatabaseBackup, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			backups, err := client.ListPostgresDatabaseBackups(ctx, dbID, nil)
			if err != nil {
				return nil, err
			}

			for _, backup := range backups {
				if backup.Label == label {
					return &backup, nil
				}
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to wait for backup %s: %w", label, ctx.Err())
		}
	}
}

type databaseStatusFunc func(ctx context.Context, client Client, dbID int) (DatabaseStatus, error)

var databaseStatusHandlers = map[DatabaseEngineType]databaseStatusFunc{
	DatabaseEngineTypeMySQL: func(ctx context.Context, client Client, dbID int) (DatabaseStatus, error) {
		db, err := client.GetMySQLDatabase(ctx, dbID)
		if err != nil {
			return "", err
		}

		return db.Status, nil
	},
	DatabaseEngineTypePostgres: func(ctx context.Context, client Client, dbID int) (DatabaseStatus, error) {
		db, err := client.GetPostgresDatabase(ctx, dbID)
		if err != nil {
			return "", err
		}

		return db.Status, nil
	},
}

// WaitForDatabaseStatus waits for the provided database to have the given status.
func (client Client) WaitForDatabaseStatus(
	ctx context.Context, dbID int, dbEngine DatabaseEngineType, status DatabaseStatus, timeoutSeconds int,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			statusHandler, ok := databaseStatusHandlers[dbEngine]
			if !ok {
				return fmt.Errorf("invalid db engine: %s", dbEngine)
			}

			currentStatus, err := statusHandler(ctx, client, dbID)
			if err != nil {
				return fmt.Errorf("failed to get db status: %w", err)
			}

			if currentStatus == status {
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("failed to wait for database %d status: %w", dbID, ctx.Err())
		}
	}
}

// NewEventPoller initializes a new Linode event poller. This should be run before the event is triggered as it stores
// the previous state of the entity's events.
func (client Client) NewEventPoller(
	ctx context.Context, id any, entityType EntityType, action EventAction,
) (*EventPoller, error) {
	result := EventPoller{
		EntityID:   id,
		EntityType: entityType,
		Action:     action,

		client: client,
	}

	if err := result.PreTask(ctx); err != nil {
		return nil, fmt.Errorf("failed to run pretask: %w", err)
	}

	return &result, nil
}

// NewEventPollerWithSecondary initializes a new Linode event poller with for events with a
// specific secondary entity.
func (client Client) NewEventPollerWithSecondary(
	ctx context.Context, id any, primaryEntityType EntityType, secondaryID int, action EventAction,
) (*EventPoller, error) {
	poller, err := client.NewEventPoller(ctx, id, primaryEntityType, action)
	if err != nil {
		return nil, err
	}

	poller.SecondaryEntityID = secondaryID

	return poller, nil
}

// NewEventPollerWithoutEntity initializes a new Linode event poller without a target entity ID.
// This is useful for create events where the ID of the entity is not yet known.
// For example:
// p, _ := client.NewEventPollerWithoutEntity(...)
// inst, _ := client.CreateInstance(...)
// p.EntityID = inst.ID
// ...
func (client Client) NewEventPollerWithoutEntity(entityType EntityType, action EventAction) (*EventPoller, error) {
	result := EventPoller{
		EntityType:     entityType,
		Action:         action,
		EntityID:       0,
		previousEvents: make(map[int]bool, 0),

		client: client,
	}

	return &result, nil
}

// PreTask stores all current events for the given entity to prevent them from being
// processed on subsequent runs.
func (p *EventPoller) PreTask(ctx context.Context) error {
	f := Filter{
		OrderBy: "created",
		Order:   Descending,
	}
	f.AddField(Eq, "entity.type", p.EntityType)
	f.AddField(Eq, "entity.id", p.EntityID)
	f.AddField(Eq, "action", p.Action)

	fBytes, err := f.MarshalJSON()
	if err != nil {
		return err
	}

	events, err := p.client.ListEvents(ctx, &ListOptions{
		Filter:      string(fBytes),
		PageOptions: &PageOptions{Page: 1},
	})
	if err != nil {
		return fmt.Errorf("failed to list events: %w", err)
	}

	eventIDs := make(map[int]bool, len(events))
	for _, event := range events {
		eventIDs[event.ID] = true
	}

	p.previousEvents = eventIDs

	return nil
}

func (p *EventPoller) WaitForLatestUnknownEvent(ctx context.Context) (*Event, error) {
	ticker := time.NewTicker(p.client.pollInterval)
	defer ticker.Stop()

	f := Filter{
		OrderBy: "created",
		Order:   Descending,
	}
	f.AddField(Eq, "entity.type", p.EntityType)
	f.AddField(Eq, "entity.id", p.EntityID)
	f.AddField(Eq, "action", p.Action)

	fBytes, err := f.MarshalJSON()
	if err != nil {
		return nil, err
	}

	listOpts := ListOptions{
		Filter:      string(fBytes),
		PageOptions: &PageOptions{Page: 1},
	}

	for {
		select {
		case <-ticker.C:
			events, err := p.client.ListEvents(ctx, &listOpts)
			if err != nil {
				return nil, fmt.Errorf("failed to list events: %w", err)
			}

			for _, event := range events {
				if p.SecondaryEntityID != nil && !eventMatchesSecondary(p.SecondaryEntityID, event) {
					continue
				}

				if _, ok := p.previousEvents[event.ID]; !ok {
					// Store this event so it is no longer picked up
					// on subsequent jobs
					p.previousEvents[event.ID] = true

					return &event, nil
				}
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to wait for event: %w", ctx.Err())
		}
	}
}

// WaitForFinished waits for a new event to be finished.
func (p *EventPoller) WaitForFinished(
	ctx context.Context, timeoutSeconds int,
) (*Event, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(p.client.pollInterval)
	defer ticker.Stop()

	event, err := p.WaitForLatestUnknownEvent(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for event: %w", err)
	}

	for {
		select {
		case <-ticker.C:
			event, err := p.client.GetEvent(ctx, event.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get event: %w", err)
			}

			switch event.Status {
			case EventFinished:
				return event, nil
			case EventFailed:
				return nil, fmt.Errorf("event %d has failed", event.ID)
			case EventScheduled, EventStarted, EventNotification:
				continue
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to wait for event finished: %w", ctx.Err())
		}
	}
}

// WaitForResourceFree waits for a resource to have no running events.
func (client Client) WaitForResourceFree(
	ctx context.Context, entityType EntityType, entityID any, timeoutSeconds int,
) error {
	apiFilter := Filter{
		Order:   Descending,
		OrderBy: "created",
	}
	apiFilter.AddField(Eq, "entity.id", entityID)
	apiFilter.AddField(Eq, "entity.type", entityType)

	filterStr, err := apiFilter.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to create filter: %s", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	ticker := time.NewTicker(client.pollInterval)
	defer ticker.Stop()

	// A helper function to determine whether a resource is busy
	checkIsBusy := func(events []Event) bool {
		for _, event := range events {
			if event.Status == EventStarted || event.Status == EventScheduled {
				return true
			}
		}

		return false
	}

	for {
		select {
		case <-ticker.C:
			events, err := client.ListEvents(ctx, &ListOptions{
				Filter: string(filterStr),
			})
			if err != nil {
				return fmt.Errorf("failed to list events: %s", err)
			}

			if !checkIsBusy(events) {
				return nil
			}

		case <-ctx.Done():
			return fmt.Errorf("failed to wait for resource free: %s", ctx.Err())
		}
	}
}

// eventMatchesSecondary returns whether the given event's secondary entity
// matches the configured secondary ID.
// This logic has been broken out to improve readability.
func eventMatchesSecondary(configuredID any, e Event) bool {
	// We should return false if the event has no secondary entity.
	// e.g. A previous disk deletion has completed.
	if e.SecondaryEntity == nil {
		return false
	}

	secondaryID := e.SecondaryEntity.ID

	// Evil hack to correct IDs parsed as floats
	if value, ok := secondaryID.(float64); ok {
		secondaryID = int(value)
	}

	return secondaryID == configuredID
}
