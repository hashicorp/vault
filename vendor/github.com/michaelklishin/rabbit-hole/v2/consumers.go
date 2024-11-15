package rabbithole

import (
	"net/url"
)

// AcknowledgementMode specifies an acknowledgement mode used by
// a consumer. Learn more at https://www.rabbitmq.com/confirms.html.
type AcknowledgementMode bool

const (
	// ManualAcknowledgement requires the consumer to explicitly
	// acknowledge processed deliveries.
	ManualAcknowledgement AcknowledgementMode = true
	// AutomaticAcknowledgment means that deliveries sent
	// to the consumer will be considered processed immediately.
	// Explicit acks from the client are not needed or expected
	// by the server.
	AutomaticAcknowledgment AcknowledgementMode = false
)

// BriefQueueInfo represents a fully qualified queue name.
type BriefQueueInfo struct {
	Name  string `json:"name"`
	Vhost string `json:"vhost"`
}

// BriefChannelDetail represents a channel with a limited
// number of metrics.
type BriefChannelDetail struct {
	ConnectionName string `json:"connection_name"`
	Name           string `json:"name"`
	Node           string `json:"node"`
	Number         int    `json:"number"`
	PeerHost       string `json:"peer_host"`
	PeerPort       Port   `json:"peer_port"`
	User           string `json:"user"`
}

// ConsumerInfo represents a consumer.
type ConsumerInfo struct {
	Arguments           map[string]interface{} `json:"arguments"`
	AcknowledgementMode AcknowledgementMode    `json:"ack_required"`
	ChannelDetails      BriefChannelDetail     `json:"channel_details"`
	ConsumerTag         string                 `json:"consumer_tag"`
	Exclusive           bool                   `json:"exclusive"`
	PrefetchCount       int                    `json:"prefetch_count"`
	Queue               BriefQueueInfo         `json:"queue"`
}

// ListConsumers lists all consumers in the cluster.
func (c *Client) ListConsumers() (rec []ConsumerInfo, err error) {
	req, err := newGETRequest(c, "consumers")
	if err != nil {
		return
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return
	}

	return
}

// ListConsumersIn lists all consumers in a virtual host.
func (c *Client) ListConsumersIn(vhost string) (rec []ConsumerInfo, err error) {
	req, err := newGETRequest(c, "consumers/"+url.PathEscape(vhost))
	if err != nil {
		return
	}

	if err = executeAndParseRequest(c, req, &rec); err != nil {
		return
	}

	return
}
