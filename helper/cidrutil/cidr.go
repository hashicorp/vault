package cidrutil

import (
	"fmt"
	"net"
	"strings"

	"github.com/hashicorp/vault/helper/strutil"
)

// IPBelongsToCIDR checks if the given IP is encompassed by the given CIDR block
func IPBelongsToCIDR(ipAddr string, cidr string) (bool, error) {
	if ipAddr == "" {
		return false, fmt.Errorf("missing IP address")
	}

	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address")
	}

	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}

	if !ipnet.Contains(ip) {
		return false, nil
	}

	return true, nil
}

// IPBelongsToCIDRBlocksString checks if the given IP is encompassed by any of
// the given CIDR blocks, when the input is a string composed by joining all
// the CIDR blocks using a separator. The input is separated based on the given
// separator and the IP is checked to be belonged by any CIDR block.
func IPBelongsToCIDRBlocksString(ipAddr string, cidrList, separator string) (bool, error) {
	if ipAddr == "" {
		return false, fmt.Errorf("missing IP address")
	}

	if cidrList == "" {
		return false, fmt.Errorf("missing CIDR list")
	}

	if separator == "" {
		return false, fmt.Errorf("missing separator")
	}

	if ip := net.ParseIP(ipAddr); ip == nil {
		return false, fmt.Errorf("invalid IP address")
	}

	return IPBelongsToCIDRBlocksSlice(ipAddr, strutil.ParseDedupLowercaseAndSortStrings(cidrList, separator))
}

// IPBelongsToCIDRBlocksSlice checks if the given IP is encompassed by any of the given
// CIDR blocks
func IPBelongsToCIDRBlocksSlice(ipAddr string, cidrs []string) (bool, error) {
	if ipAddr == "" {
		return false, fmt.Errorf("missing IP address")
	}

	if len(cidrs) == 0 {
		return false, fmt.Errorf("missing CIDR blocks to be checked against")
	}

	if ip := net.ParseIP(ipAddr); ip == nil {
		return false, fmt.Errorf("invalid IP address")
	}

	for _, cidr := range cidrs {
		belongs, err := IPBelongsToCIDR(ipAddr, cidr)
		if err != nil {
			return false, err
		}
		if belongs {
			return true, nil
		}
	}

	return false, nil
}

// ValidateCIDRListString checks if the list of CIDR blocks are valid, given
// that the input is a string composed by joining all the CIDR blocks using a
// separator. The input is separated based on the given separator and validity
// of each is checked.
func ValidateCIDRListString(cidrList string, separator string) (bool, error) {
	if cidrList == "" {
		return false, fmt.Errorf("missing CIDR list that needs validation")
	}
	if separator == "" {
		return false, fmt.Errorf("missing separator")
	}

	return ValidateCIDRListSlice(strutil.ParseDedupLowercaseAndSortStrings(cidrList, separator))
}

// ValidateCIDRListSlice checks if the given list of CIDR blocks are valid
func ValidateCIDRListSlice(cidrBlocks []string) (bool, error) {
	if len(cidrBlocks) == 0 {
		return false, fmt.Errorf("missing CIDR blocks that needs validation")
	}

	for _, block := range cidrBlocks {
		if _, _, err := net.ParseCIDR(strings.TrimSpace(block)); err != nil {
			return false, err
		}
	}

	return true, nil
}

// Subset checks if the IPs belonging to a given CIDR block is a subset of IPs
// belonging to another CIDR block.
func Subset(cidr1, cidr2 string) (bool, error) {
	if cidr1 == "" {
		return false, fmt.Errorf("missing CIDR to be checked against")
	}

	if cidr2 == "" {
		return false, fmt.Errorf("missing CIDR that needs to be checked")
	}

	ip1, net1, err := net.ParseCIDR(cidr1)
	if err != nil {
		return false, fmt.Errorf("failed to parse the CIDR to be checked against: %q", err)
	}

	zeroAddr := false
	if ip := ip1.To4(); ip != nil && ip.Equal(net.IPv4zero) {
		zeroAddr = true
	}
	if ip := ip1.To16(); ip != nil && ip.Equal(net.IPv6zero) {
		zeroAddr = true
	}

	maskLen1, _ := net1.Mask.Size()
	if !zeroAddr && maskLen1 == 0 {
		return false, fmt.Errorf("CIDR to be checked against is not in its canonical form")
	}

	ip2, net2, err := net.ParseCIDR(cidr2)
	if err != nil {
		return false, fmt.Errorf("failed to parse the CIDR that needs to be checked: %q", err)
	}

	zeroAddr = false
	if ip := ip2.To4(); ip != nil && ip.Equal(net.IPv4zero) {
		zeroAddr = true
	}
	if ip := ip2.To16(); ip != nil && ip.Equal(net.IPv6zero) {
		zeroAddr = true
	}

	maskLen2, _ := net2.Mask.Size()
	if !zeroAddr && maskLen2 == 0 {
		return false, fmt.Errorf("CIDR that needs to be checked is not in its canonical form")
	}

	// If the mask length of the CIDR that needs to be checked is smaller
	// then the mask length of the CIDR to be checked against, then the
	// former will encompass more IPs than the latter, and hence can't be a
	// subset of the latter.
	if maskLen2 < maskLen1 {
		return false, nil
	}

	belongs, err := IPBelongsToCIDR(net2.IP.String(), cidr1)
	if err != nil {
		return false, err
	}

	return belongs, nil
}

// SubsetBlocks checks if each CIDR block of a given set of CIDR blocks, is a
// subset of at least one CIDR block belonging to another set of CIDR blocks.
// First parameter is the set of CIDR blocks to check against and the second
// parameter is the set of CIDR blocks that needs to be checked.
func SubsetBlocks(cidrBlocks1, cidrBlocks2 []string) (bool, error) {
	if len(cidrBlocks1) == 0 {
		return false, fmt.Errorf("missing CIDR blocks to be checked against")
	}

	if len(cidrBlocks2) == 0 {
		return false, fmt.Errorf("missing CIDR blocks that needs to be checked")
	}

	// Check if all the elements of cidrBlocks2 is a subset of at least one
	// element of cidrBlocks1
	for _, cidrBlock2 := range cidrBlocks2 {
		isSubset := false
		for _, cidrBlock1 := range cidrBlocks1 {
			subset, err := Subset(cidrBlock1, cidrBlock2)
			if err != nil {
				return false, err
			}
			// If CIDR is a subset of any of the CIDR block, its
			// good enough. Break out.
			if subset {
				isSubset = true
				break
			}
		}
		// CIDR block was not a subset of any of the CIDR blocks in the
		// set of blocks to check against
		if !isSubset {
			return false, nil
		}
	}

	return true, nil
}
