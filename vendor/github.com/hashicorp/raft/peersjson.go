package raft

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

// ReadPeersJSON consumes a legacy peers.json file in the format of the old JSON
// peer store and creates a new-style configuration structure. This can be used
// to migrate this data or perform manual recovery when running protocol versions
// that can interoperate with older, unversioned Raft servers. This should not be
// used once server IDs are in use, because the old peers.json file didn't have
// support for these, nor non-voter suffrage types.
func ReadPeersJSON(path string) (Configuration, error) {
	// Read in the file.
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return Configuration{}, err
	}

	// Parse it as JSON.
	var peers []string
	dec := json.NewDecoder(bytes.NewReader(buf))
	if err := dec.Decode(&peers); err != nil {
		return Configuration{}, err
	}

	// Map it into the new-style configuration structure. We can only specify
	// voter roles here, and the ID has to be the same as the address.
	var configuration Configuration
	for _, peer := range peers {
		server := Server{
			Suffrage: Voter,
			ID:       ServerID(peer),
			Address:  ServerAddress(peer),
		}
		configuration.Servers = append(configuration.Servers, server)
	}

	// We should only ingest valid configurations.
	if err := checkConfiguration(configuration); err != nil {
		return Configuration{}, err
	}
	return configuration, nil
}

// configEntry is used when decoding a new-style peers.json.
type configEntry struct {
	// ID is the ID of the server (a UUID, usually).
	ID ServerID `json:"id"`

	// Address is the host:port of the server.
	Address ServerAddress `json:"address"`

	// NonVoter controls the suffrage. We choose this sense so people
	// can leave this out and get a Voter by default.
	NonVoter bool `json:"non_voter"`
}

// ReadConfigJSON reads a new-style peers.json and returns a configuration
// structure. This can be used to perform manual recovery when running protocol
// versions that use server IDs.
func ReadConfigJSON(path string) (Configuration, error) {
	// Read in the file.
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return Configuration{}, err
	}

	// Parse it as JSON.
	var peers []configEntry
	dec := json.NewDecoder(bytes.NewReader(buf))
	if err := dec.Decode(&peers); err != nil {
		return Configuration{}, err
	}

	// Map it into the new-style configuration structure.
	var configuration Configuration
	for _, peer := range peers {
		suffrage := Voter
		if peer.NonVoter {
			suffrage = Nonvoter
		}
		server := Server{
			Suffrage: suffrage,
			ID:       peer.ID,
			Address:  peer.Address,
		}
		configuration.Servers = append(configuration.Servers, server)
	}

	// We should only ingest valid configurations.
	if err := checkConfiguration(configuration); err != nil {
		return Configuration{}, err
	}
	return configuration, nil
}
