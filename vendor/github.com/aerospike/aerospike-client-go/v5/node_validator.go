// Copyright 2014-2021 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/aerospike/aerospike-client-go/v5/logger"
	"github.com/aerospike/aerospike-client-go/v5/types"
)

type nodesToAddT map[string]*Node

func (nta nodesToAddT) addNodeIfNotExists(ndv *nodeValidator, cluster *Cluster) bool {
	_, exists := nta[ndv.name]
	if !exists {
		// found a new node
		node := cluster.createNode(ndv)
		nta[ndv.name] = node
	}
	return exists
}

// Validates a Database server node
type nodeValidator struct {
	name        string
	aliases     []*Host
	primaryHost *Host

	detectLoadBalancer bool

	sessionInfo *sessionInfo

	features int
}

func (ndv *nodeValidator) seedNodes(cluster *Cluster, host *Host, nodesToAdd nodesToAddT) Error {
	if err := ndv.setAliases(host); err != nil {
		return err
	}

	found := false
	var resultErr Error
	for _, alias := range ndv.aliases {
		if resultErr = ndv.validateAlias(cluster, alias); resultErr != nil {
			logger.Logger.Debug("Alias %s failed: %s", alias, resultErr)
			continue
		}

		found = true
		nodesToAdd.addNodeIfNotExists(ndv, cluster)
	}

	if !found {
		return resultErr
	}
	return nil
}

func (ndv *nodeValidator) validateNode(cluster *Cluster, host *Host) Error {
	if clusterNodes := cluster.GetNodes(); cluster.clientPolicy.IgnoreOtherSubnetAliases && len(clusterNodes) > 0 {
		masterHostname := clusterNodes[0].host.Name
		ip, ipnet, err := net.ParseCIDR(masterHostname + "/24")
		if err != nil {
			logger.Logger.Error(err.Error())
			return newError(types.NO_AVAILABLE_CONNECTIONS_TO_NODE, "Failed parsing hostname...")
		}

		stop := ip.Mask(ipnet.Mask)
		stop[3] += 255
		if bytes.Compare(net.ParseIP(host.Name).To4(), ip.Mask(ipnet.Mask).To4()) >= 0 && bytes.Compare(net.ParseIP(host.Name).To4(), stop.To4()) >= 0 {
			return newError(types.NO_AVAILABLE_CONNECTIONS_TO_NODE, "Ignored hostname from other subnet...")
		}
	}

	if err := ndv.setAliases(host); err != nil {
		return err
	}

	var resultErr Error
	for _, alias := range ndv.aliases {
		if err := ndv.validateAlias(cluster, alias); err != nil {
			resultErr = chainErrors(err, resultErr)
			logger.Logger.Debug("Aliases %s failed: %s", alias, err)
			continue
		}
		return nil
	}

	return resultErr
}

func (ndv *nodeValidator) setAliases(host *Host) Error {
	ndv.detectLoadBalancer = true

	// IP addresses do not need a lookup
	ip := net.ParseIP(host.Name)
	if ip != nil {
		// avoid detecting load balancer on localhost
		ndv.detectLoadBalancer = !ip.IsLoopback()

		aliases := make([]*Host, 1)
		aliases[0] = NewHost(host.Name, host.Port)
		aliases[0].TLSName = host.TLSName
		ndv.aliases = aliases
	} else {
		addresses, err := net.LookupHost(host.Name)
		if err != nil {
			logger.Logger.Error("Host lookup failed with error: %s", err.Error())
			return errToAerospikeErr(nil, err)
		}
		aliases := make([]*Host, len(addresses))
		for idx, addr := range addresses {
			aliases[idx] = NewHost(addr, host.Port)
			aliases[idx].TLSName = host.TLSName

			// avoid detecting load balancer on localhost
			if ip := net.ParseIP(host.Name); ip != nil && ip.IsLoopback() {
				ndv.detectLoadBalancer = false
			}
		}
		ndv.aliases = aliases
	}
	logger.Logger.Debug("Node Validator has %d nodes and they are: %v", len(ndv.aliases), ndv.aliases)
	return nil
}

func (ndv *nodeValidator) validateAlias(cluster *Cluster, alias *Host) Error {
	clientPolicy := cluster.clientPolicy
	clientPolicy.Timeout /= 2

	conn, err := NewConnection(&clientPolicy, alias)
	if err != nil {
		return err
	}
	defer conn.Close()

	if clientPolicy.RequiresAuthentication() {
		// need to authenticate
		acmd := newLoginCommand(conn.dataBuffer)
		err = acmd.login(&clientPolicy, conn, cluster.Password())
		if err != nil {
			return err
		}

		ndv.sessionInfo = acmd.sessionInfo()
	}

	// check to make sure we have actually connected
	info, err := conn.RequestInfo("build")
	if err != nil {
		return err
	}

	if _, exists := info["ERROR:80:not authenticated"]; exists {
		return ErrNotAuthenticated.err()
	}

	hasClusterName := len(clientPolicy.ClusterName) > 0

	infoKeys := []string{"node", "partition-generation", "features"}
	if hasClusterName {
		infoKeys = append(infoKeys, "cluster-name")
	}

	addressCommand := clientPolicy.serviceString()
	if ndv.detectLoadBalancer {
		infoKeys = append(infoKeys, addressCommand)
	}

	infoMap, err := conn.RequestInfo(infoKeys...)
	if err != nil {
		return err
	}

	nodeName, exists := infoMap["node"]
	if !exists {
		return newError(types.INVALID_NODE_ERROR, "Invalid node alias:"+alias.String())
	}

	genStr, exists := infoMap["partition-generation"]
	if !exists {
		return newError(types.INVALID_NODE_ERROR, "Invalid partition-generation for node:"+alias.String())
	}

	gen, nerr := strconv.Atoi(genStr)
	if nerr != nil {
		return newError(types.PARSE_ERROR, fmt.Sprintf("Invalid partition-generation for Node %s (%s), value: %s", nodeName, alias.String(), genStr))
	}

	if gen == -1 {
		return newError(types.INVALID_NODE_ERROR, fmt.Sprintf("Node %s (%s) is not yet fully initialized", nodeName, alias.String()))
	}

	if hasClusterName {
		id := infoMap["cluster-name"]

		if len(id) == 0 || id != clientPolicy.ClusterName {
			return newError(types.CLUSTER_NAME_MISMATCH_ERROR, fmt.Sprintf("Node %s (%s) expected cluster name `%s` but received `%s`", nodeName, alias.String(), clientPolicy.ClusterName, id))
		}
	}

	// set features
	if features, exists := infoMap["features"]; exists {
		ndv.setFeatures(features)
	}

	// This client requires partition scan support. Partition scans were first
	// supported in server version 4.9. Do not allow any server node into the
	// cluster that is running server version < 4.9.
	if (ndv.features & _SUPPORTS_PARTITION_SCAN) == 0 {
		return newError(types.INVALID_NODE_ERROR, fmt.Sprintf("Node %s (%s) is version < 4.9. This client supports server versions >= 4.9", nodeName, alias.String()))
	}

	// check if the host is a load-balancer
	if peersStr, exists := infoMap[addressCommand]; exists {
		var hostAddress []*Host
		peerParser := peerListParser{buf: []byte("[" + peersStr + "]")}
		if hostAddress, err = peerParser.readHosts(alias.TLSName); err != nil {
			logger.Logger.Error("Failed to parse `%s` results... err: %s", alias.String(), err.Error())
		}

		if len(hostAddress) > 0 {
			isLoadBalancer := true
		LOAD_BALANCER:
			for _, h := range hostAddress {
				for _, a := range ndv.aliases {
					if h.equals(a) {
						// one of the aliases were the same as an advertised service
						// no need to replace the seed host with the alias
						isLoadBalancer = false
						break LOAD_BALANCER
					}
				}
			}

			if isLoadBalancer && ndv.detectLoadBalancer {
				aliasFound := false

				// take the seed out of the aliases if it is load balancer
				logger.Logger.Info("Host `%s` seems to be a load balancer. It is going to be replace by `%v`", alias.String(), hostAddress[0])
				// try to connect to the aliases, and coose the first one that connects
				for _, h := range hostAddress {
					hconn, err := NewConnection(&clientPolicy, h)
					if err != nil {
						continue
					}
					defer hconn.Close()

					if clientPolicy.RequiresAuthentication() {
						// need to authenticate
						acmd := newLoginCommand(hconn.dataBuffer)
						err = acmd.login(&clientPolicy, hconn, cluster.Password())
						if err != nil {
							continue
						}

						ndv.sessionInfo = acmd.sessionInfo()
					}

					alias = h
					ndv.aliases = hostAddress
					aliasFound = true

					// found one, no need to try the rest
					break
				}

				// Failed to find a valid address to connect. IP Address is probably internal on the cloud
				// because the server access-address is not configured.  Log warning and continue
				// with original seed.
				if !aliasFound {
					logger.Logger.Info("Inaccessible address `%s` as cluster seed. access-address is probably not configured on server.", alias.String())
				}
			}
		}
	}

	ndv.name = nodeName
	ndv.primaryHost = alias

	return nil
}

func (ndv *nodeValidator) setFeatures(features string) {
	featureList := strings.Split(features, ";")
	for i := range featureList {
		switch featureList[i] {
		case "pscans":
			ndv.features |= _SUPPORTS_PARTITION_SCAN
		case "query-show":
			ndv.features |= _SUPPORTS_QUERY_SHOW
		}
	}
}
