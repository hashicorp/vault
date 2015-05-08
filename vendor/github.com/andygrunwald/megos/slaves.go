package megos

import (
	"fmt"
)

// GetSlaveByID will return a a slave by its unique ID (slaveID).
//
// The list of slaves are provided by a state of a single node.
func (c *Client) GetSlaveByID(slaves []Slave, slaveID string) (*Slave, error) {
	for _, slave := range slaves {
		if slaveID == slave.ID {
			return &slave, nil
		}
	}

	return nil, fmt.Errorf("No slave found with id \"%s\"", slaveID)
}
