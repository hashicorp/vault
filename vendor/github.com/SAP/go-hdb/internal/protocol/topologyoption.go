// SPDX-FileCopyrightText: 2014-2020 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package protocol

//go:generate stringer -type=topologyOption

type topologyOption int8

const (
	toHostName         topologyOption = 1
	toHostPortnumber   topologyOption = 2
	toTenantName       topologyOption = 3
	toLoadfactor       topologyOption = 4
	toVolumeID         topologyOption = 5
	toIsPrimary        topologyOption = 6
	toIsCurrentSession topologyOption = 7
	toServiceType      topologyOption = 8
	toNetworkDomain    topologyOption = 9
	toIsStandby        topologyOption = 10
	toAllIPAddresses   topologyOption = 11
	toAllHostNames     topologyOption = 12
)
