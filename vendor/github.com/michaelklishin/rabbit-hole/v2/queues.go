package rabbithole

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// BackingQueueStatus exposes backing queue (queue storage engine) metrics.
// They can change in a future version of RabbitMQ.
type BackingQueueStatus struct {
	Q1 int `json:"q1,omitempty"`
	Q2 int `json:"q2,omitempty"`
	Q3 int `json:"q3,omitempty"`
	Q4 int `json:"q4,omitempty"`
	// Total queue length
	Length int64 `json:"len,omitempty"`
	// Number of pending acks from consumers
	PendingAcks int64 `json:"pending_acks,omitempty"`
	// Number of messages held in RAM
	RAMMessageCount int64 `json:"ram_msg_count,omitempty"`
	// Number of outstanding acks held in RAM
	RAMAckCount int64 `json:"ram_ack_count,omitempty"`
	// Number of persistent messages in the store
	PersistentCount int64 `json:"persistent_count,omitempty"`
	// Average ingress (inbound) rate, not including messages
	// that straight through to auto-acking consumers.
	AverageIngressRate float64 `json:"avg_ingress_rate,omitempty"`
	// Average egress (outbound) rate, not including messages
	// that straight through to auto-acking consumers.
	AverageEgressRate float64 `json:"avg_egress_rate,omitempty"`
	// rate at which unacknowledged message records enter RAM,
	// e.g. because messages are delivered requiring acknowledgement
	AverageAckIngressRate float32 `json:"avg_ack_ingress_rate,omitempty"`
	// rate at which unacknowledged message records leave RAM,
	// e.g. because acks arrive or unacked messages are paged out
	AverageAckEgressRate float32 `json:"avg_ack_egress_rate,omitempty"`
}

// OwnerPidDetails describes an exclusive queue owner (connection).
type OwnerPidDetails struct {
	Name     string `json:"name,omitempty"`
	PeerPort Port   `json:"peer_port,omitempty"`
	PeerHost string `json:"peer_host,omitempty"`
}

// ConsumerDetail describe consumer information with a queue
type ConsumerDetail struct {
	Arguments      map[string]interface{} `json:"arguments"`
	ChannelDetails ChannelDetails         `json:"channel_details"`
	AckRequired    bool                   `json:"ack_required"`
	Active         bool                   `json:"active"`
	ActiveStatus   string                 `json:"active_status"`
	ConsumerTag    string                 `json:"consumer_tag"`
	Exclusive      bool                   `json:"exclusive,omitempty"`
	PrefetchCount  uint                   `json:"prefetch_count"`
	Queue          QueueDetail            `json:"queue"`
}

// ChannelDetails describe channel information with a consumer
type ChannelDetails struct {
	ConnectionName string `json:"connection_name"`
	Name           string `json:"name"`
	Node           string `json:"node"`
	Number         uint   `json:"number"`
	PeerHost       string `json:"peer_host"`
	PeerPort       Port   `json:"peer_port"`
	User           string `json:"user"`
}

// Handles special case where `ChannelDetails` is an empty array
// See https://github.com/rabbitmq/rabbitmq-server/issues/2684
func (c *ChannelDetails) UnmarshalJSON(data []byte) error {
	if string(data) == "[]" {
		*c = ChannelDetails{}
		return nil
	}
	type Alias ChannelDetails
	return json.Unmarshal(data, (*Alias)(c))
}

// QueueDetail describe queue information with a consumer
type QueueDetail struct {
	Name  string `json:"name"`
	Vhost string `json:"vhost,omitempty"`
}

// GarbageCollectionDetail describe queue garbage collection information
type GarbageCollectionDetails struct {
	FullSweepAfter  int `json:"fullsweep_after"`
	MaxHeapSize     int `json:"max_heap_size"`
	MinBinVheapSize int `json:"min_bin_vheap_size"`
	MinHeapSize     int `json:"min_heap_size"`
	MinorGCs        int `json:"minor_gcs"`
}

// QueueInfo represents a queue, its properties and key metrics.
type QueueInfo struct {
	// Queue name
	Name string `json:"name"`
	// Queue type
	Type string `json:"type,omitempty"`
	// Virtual host this queue belongs to
	Vhost string `json:"vhost,omitempty"`
	// Is this queue durable?
	Durable bool `json:"durable"`
	// Is this queue auto-deleted?
	AutoDelete AutoDelete `json:"auto_delete"`
	// Is this queue exclusive?
	Exclusive bool `json:"exclusive,omitempty"`
	// Extra queue arguments
	Arguments map[string]interface{} `json:"arguments"`

	// RabbitMQ node that hosts the leader replica for this queue
	Node string `json:"node,omitempty"`
	// Queue status
	Status string `json:"state,omitempty"`
	// Queue leader when it is quorum queue
	Leader string `json:"leader,omitempty"`
	// Queue members when it is quorum queue
	Members []string `json:"members,omitempty"`
	// Queue online members when it is quorum queue
	Online []string `json:"online,omitempty"`

	// Total amount of RAM used by this queue
	Memory int64 `json:"memory,omitempty"`
	// How many consumers this queue has
	Consumers int `json:"consumers,omitempty"`
	// Detail information of consumers
	ConsumerDetails *[]ConsumerDetail `json:"consumer_details,omitempty"`
	// Utilisation of all the consumers
	ConsumerUtilisation float64 `json:"consumer_utilisation,omitempty"`
	// If there is an exclusive consumer, its consumer tag
	ExclusiveConsumerTag string `json:"exclusive_consumer_tag,omitempty"`

	// GarbageCollection metrics
	GarbageCollection *GarbageCollectionDetails `json:"garbage_collection,omitempty"`

	// Policy applied to this queue, if any
	Policy string `json:"policy,omitempty"`

	// Total bytes of messages in this queues
	MessagesBytes               int64 `json:"message_bytes,omitempty"`
	MessagesBytesPersistent     int64 `json:"message_bytes_persistent,omitempty"`
	MessagesBytesRAM            int64 `json:"message_bytes_ram,omitempty"`
	MessagesBytesReady          int64 `json:"message_bytes_ready,omitempty"`
	MessagesBytesUnacknowledged int64 `json:"message_bytes_unacknowledged,omitempty"`

	// Total number of messages in this queue
	Messages           int          `json:"messages,omitempty"`
	MessagesDetails    *RateDetails `json:"messages_details,omitempty"`
	MessagesPersistent int          `json:"messages_persistent,omitempty"`
	MessagesRAM        int          `json:"messages_ram,omitempty"`

	// Number of messages ready to be delivered
	MessagesReady        int          `json:"messages_ready,omitempty"`
	MessagesReadyDetails *RateDetails `json:"messages_ready_details,omitempty"`

	// Number of messages delivered and pending acknowledgements from consumers
	MessagesUnacknowledged        int          `json:"messages_unacknowledged,omitempty"`
	MessagesUnacknowledgedDetails *RateDetails `json:"messages_unacknowledged_details,omitempty"`

	MessageStats *MessageStats `json:"message_stats,omitempty"`

	OwnerPidDetails *OwnerPidDetails `json:"owner_pid_details,omitempty"`

	BackingQueueStatus *BackingQueueStatus `json:"backing_queue_status,omitempty"`

	ActiveConsumers int64 `json:"active_consumers,omitempty"`
}

// PagedQueueInfo is additional context returned for paginated requests.
type PagedQueueInfo struct {
	Page          int         `json:"page"`
	PageCount     int         `json:"page_count"`
	PageSize      int         `json:"page_size"`
	FilteredCount int         `json:"filtered_count"`
	ItemCount     int         `json:"item_count"`
	TotalCount    int         `json:"total_count"`
	Items         []QueueInfo `json:"items"`
}

// DetailedQueueInfo is an alias for QueueInfo
type DetailedQueueInfo QueueInfo

//
// GET /api/queues
//

// [
//   {
//     "owner_pid_details": {
//       "name": "127.0.0.1:46928 -> 127.0.0.1:5672",
//       "peer_port": 46928,
//       "peer_host": "127.0.0.1"
//     },
//     "message_stats": {
//       "publish": 19830,
//       "publish_details": {
//         "rate": 5
//       }
//     },
//     "messages": 15,
//     "messages_details": {
//       "rate": 0
//     },
//     "messages_ready": 15,
//     "messages_ready_details": {
//       "rate": 0
//     },
//     "messages_unacknowledged": 0,
//     "messages_unacknowledged_details": {
//       "rate": 0
//     },
//     "policy": "",
//     "exclusive_consumer_tag": "",
//     "consumers": 0,
//     "memory": 143112,
//     "backing_queue_status": {
//       "q1": 0,
//       "q2": 0,
//       "delta": [
//         "delta",
//         "undefined",
//         0,
//         "undefined"
//       ],
//       "q3": 0,
//       "q4": 15,
//       "len": 15,
//       "pending_acks": 0,
//       "target_ram_count": "infinity",
//       "ram_msg_count": 15,
//       "ram_ack_count": 0,
//       "next_seq_id": 19830,
//       "persistent_count": 0,
//       "avg_ingress_rate": 4.9920127795527,
//       "avg_egress_rate": 4.9920127795527,
//       "avg_ack_ingress_rate": 0,
//       "avg_ack_egress_rate": 0
//     },
//     "status": "running",
//     "name": "amq.gen-QLEaT5Rn_ogbN3O8ZOQt3Q",
//     "vhost": "rabbit\/hole",
//     "durable": false,
//     "auto_delete": false,
//     "arguments": {
//       "x-message-ttl": 5000
//     },
//     "node": "rabbit@marzo"
//   }
// ]

// ListQueues lists all queues in the cluster. This only includes queues in the
// virtual hosts accessible to the user.
func (c *Client) ListQueues() (rec []QueueInfo, err error) {
	req, err := newGETRequest(c, "queues")
	if err != nil {
		return []QueueInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []QueueInfo{}, err
	}

	return rec, nil
}

// ListQueuesWithParameters lists queues with a list of query string values.
func (c *Client) ListQueuesWithParameters(params url.Values) (rec []QueueInfo, err error) {
	req, err := newGETRequestWithParameters(c, "queues", params)
	if err != nil {
		return []QueueInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []QueueInfo{}, err
	}

	return rec, nil
}

// ListQueuesWithParametersIs lists queues with a list of query string values in the vhost vhost.
func (c *Client) ListQueuesWithParametersIn(vhost string, params url.Values) (rec []QueueInfo, err error) {
	req, err := newGETRequestWithParameters(c, "queues/"+url.PathEscape(vhost), params)
	if err != nil {
		return []QueueInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []QueueInfo{}, err
	}

	return rec, nil
}

// PagedListQueuesWithParameters lists queues with pagination.
func (c *Client) PagedListQueuesWithParameters(params url.Values) (rec PagedQueueInfo, err error) {
	req, err := newGETRequestWithParameters(c, "queues", params)
	if err != nil {
		return PagedQueueInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return PagedQueueInfo{}, err
	}

	return rec, nil
}

// PagedListQueuesWithParameters lists queues with pagination in the vhost vhost.
func (c *Client) PagedListQueuesWithParametersIn(vhost string, params url.Values) (rec PagedQueueInfo, err error) {
	req, err := newGETRequestWithParameters(c, "queues/"+url.PathEscape(vhost), params)
	if err != nil {
		return PagedQueueInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return PagedQueueInfo{}, err
	}

	return rec, nil
}

//
// GET /api/queues/{vhost}
//

// ListQueuesIn lists all queues in a virtual host.
func (c *Client) ListQueuesIn(vhost string) (rec []QueueInfo, err error) {
	req, err := newGETRequest(c, "queues/"+url.PathEscape(vhost))
	if err != nil {
		return []QueueInfo{}, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return []QueueInfo{}, err
	}

	return rec, nil
}

//
// GET /api/queues/{vhost}/{name}
//

// GetQueue returns information about a queue.
func (c *Client) GetQueue(vhost, queue string) (rec *DetailedQueueInfo, err error) {
	req, err := newGETRequest(c, "queues/"+url.PathEscape(vhost)+"/"+url.PathEscape(queue))
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// GET /api/queues/{vhost}/{name}?{query}

// GetQueueWithParameters returns information about a queue. Compared to the regular GetQueue function,
// this one accepts additional query string values.
func (c *Client) GetQueueWithParameters(vhost, queue string, qs url.Values) (rec *DetailedQueueInfo, err error) {
	req, err := newGETRequestWithParameters(c, "queues/"+url.PathEscape(vhost)+"/"+url.PathEscape(queue), qs)
	if err != nil {
		return nil, err
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return nil, err
	}

	return rec, nil
}

//
// PUT /api/exchanges/{vhost}/{exchange}
//

// QueueSettings represents queue properties. Use it to declare a queue.
type QueueSettings struct {
	Type       string                 `json:"type"`
	Durable    bool                   `json:"durable"`
	AutoDelete bool                   `json:"auto_delete,omitempty"`
	Arguments  map[string]interface{} `json:"arguments,omitempty"`
}

// DeclareQueue declares a queue.
func (c *Client) DeclareQueue(vhost, queue string, info QueueSettings) (res *http.Response, err error) {
	if info.Arguments == nil {
		info.Arguments = make(map[string]interface{})
	}

	if info.Type != "" {
		info.Arguments["x-queue-type"] = info.Type
	}

	body, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "PUT", "queues/"+url.PathEscape(vhost)+"/"+url.PathEscape(queue), body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/queues/{vhost}/{name}
//

// Options for deleting a queue. Use it with DeleteQueue.
type QueueDeleteOptions struct {
	// Only delete the queue if there are no messages.
	IfEmpty bool
	// Only delete the queue if there are no consumers.
	IfUnused bool
}

// DeleteQueue deletes a queue.
func (c *Client) DeleteQueue(vhost, queue string, opts ...QueueDeleteOptions) (res *http.Response, err error) {
	query := url.Values{}
	for _, o := range opts {
		if o.IfEmpty {
			query["if-empty"] = []string{"true"}
		}
		if o.IfUnused {
			query["if-unused"] = []string{"true"}
		}
	}

	req, err := newRequestWithBody(c, "DELETE", "queues/"+url.PathEscape(vhost)+"/"+url.PathEscape(queue)+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

//
// DELETE /api/queues/{vhost}/{name}/contents
//

// PurgeQueue purges a queue (deletes all messages ready for delivery in it).
func (c *Client) PurgeQueue(vhost, queue string) (res *http.Response, err error) {
	req, err := newRequestWithBody(c, "DELETE", "queues/"+url.PathEscape(vhost)+"/"+url.PathEscape(queue)+"/contents", nil)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}

// queueAction represents an action that can be performed on a queue (sync/cancel_sync)
type queueAction struct {
	Action string `json:"action"`
}

// SyncQueue synchronises queue contents with the mirrors remaining in the cluster.
func (c *Client) SyncQueue(vhost, queue string) (res *http.Response, err error) {
	return c.sendQueueAction(vhost, queue, queueAction{"sync"})
}

// CancelSyncQueue cancels queue synchronisation process.
func (c *Client) CancelSyncQueue(vhost, queue string) (res *http.Response, err error) {
	return c.sendQueueAction(vhost, queue, queueAction{"cancel_sync"})
}

// POST /api/queues/{vhost}/{name}/actions
func (c *Client) sendQueueAction(vhost string, queue string, action queueAction) (res *http.Response, err error) {
	body, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}

	req, err := newRequestWithBody(c, "POST", "queues/"+url.PathEscape(vhost)+"/"+url.PathEscape(queue)+"/actions", body)
	if err != nil {
		return nil, err
	}

	if res, err = executeRequest(c, req); err != nil {
		return nil, err
	}

	return res, nil
}
