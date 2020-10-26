package rabbithole

import (
	"net/url"
)

type AcknowledgementMode bool

const (
	ManualAcknowledgement   AcknowledgementMode = true
	AutomaticAcknowledgment AcknowledgementMode = false
)

type BriefQueueInfo struct {
	Name  string `json:"name"`
	Vhost string `json:"vhost"`
}

type BriefChannelDetail struct {
	ConnectionName string `json:"connection_name"`
	Name           string `json:"name"`
	Node           string `json:"node"`
	Number         int    `json:"number"`
	PeerHost       string `json:"peer_host"`
	PeerPort       int    `json:"peer_port"`
	User           string `json:"user"`
}

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
