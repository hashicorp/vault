package cidrutil

import (
	"testing"

	sockaddr "github.com/hashicorp/go-sockaddr"
)

func TestCIDRUtil_IPBelongsToCIDR(t *testing.T) {
	ip := "192.168.25.30"
	cidr := "192.168.26.30/16"

	belongs, err := IPBelongsToCIDR(ip, cidr)
	if err != nil {
		t.Fatal(err)
	}
	if !belongs {
		t.Fatalf("expected IP %q to belong to CIDR %q", ip, cidr)
	}

	ip = "10.197.192.6"
	cidr = "10.197.192.0/18"
	belongs, err = IPBelongsToCIDR(ip, cidr)
	if err != nil {
		t.Fatal(err)
	}
	if !belongs {
		t.Fatalf("expected IP %q to belong to CIDR %q", ip, cidr)
	}

	ip = "192.168.25.30"
	cidr = "192.168.26.30/24"
	belongs, err = IPBelongsToCIDR(ip, cidr)
	if err != nil {
		t.Fatal(err)
	}
	if belongs {
		t.Fatalf("expected IP %q to not belong to CIDR %q", ip, cidr)
	}

	ip = "192.168.25.30.100"
	cidr = "192.168.26.30/24"
	belongs, err = IPBelongsToCIDR(ip, cidr)
	if err == nil {
		t.Fatalf("expected an error")
	}
}

func TestCIDRUtil_IPBelongsToCIDRBlocksSlice(t *testing.T) {
	ip := "192.168.27.29"
	cidrList := []string{"172.169.100.200/18", "192.168.0.0/16", "10.10.20.20/24"}

	belongs, err := IPBelongsToCIDRBlocksSlice(ip, cidrList)
	if err != nil {
		t.Fatal(err)
	}
	if !belongs {
		t.Fatalf("expected IP %q to belong to one of the CIDRs in %q", ip, cidrList)
	}

	ip = "192.168.27.29"
	cidrList = []string{"172.169.100.200/18", "192.168.0.0.0/16", "10.10.20.20/24"}

	belongs, err = IPBelongsToCIDRBlocksSlice(ip, cidrList)
	if err == nil {
		t.Fatalf("expected an error")
	}

	ip = "30.40.50.60"
	cidrList = []string{"172.169.100.200/18", "192.168.0.0/16", "10.10.20.20/24"}

	belongs, err = IPBelongsToCIDRBlocksSlice(ip, cidrList)
	if err != nil {
		t.Fatal(err)
	}
	if belongs {
		t.Fatalf("expected IP %q to not belong to one of the CIDRs in %q", ip, cidrList)
	}
}

func TestCIDRUtil_ValidateCIDRListString(t *testing.T) {
	cidrList := "172.169.100.200/18,192.168.0.0/16,10.10.20.20/24"

	valid, err := ValidateCIDRListString(cidrList, ",")
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatalf("expected CIDR list %q to be valid", cidrList)
	}

	cidrList = "172.169.100.200,192.168.0.0/16,10.10.20.20/24"
	valid, err = ValidateCIDRListString(cidrList, ",")
	if err == nil {
		t.Fatal("expected an error")
	}

	cidrList = "172.169.100.200/18,192.168.0.0.0/16,10.10.20.20/24"
	valid, err = ValidateCIDRListString(cidrList, ",")
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestCIDRUtil_ValidateCIDRListSlice(t *testing.T) {
	cidrList := []string{"172.169.100.200/18", "192.168.0.0/16", "10.10.20.20/24"}

	valid, err := ValidateCIDRListSlice(cidrList)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatalf("expected CIDR list %q to be valid", cidrList)
	}

	cidrList = []string{"172.169.100.200", "192.168.0.0/16", "10.10.20.20/24"}
	valid, err = ValidateCIDRListSlice(cidrList)
	if err == nil {
		t.Fatal("expected an error")
	}

	cidrList = []string{"172.169.100.200/18", "192.168.0.0.0/16", "10.10.20.20/24"}
	valid, err = ValidateCIDRListSlice(cidrList)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestCIDRUtil_Subset(t *testing.T) {
	cidr1 := "192.168.27.29/24"
	cidr2 := "192.168.27.29/24"
	subset, err := Subset(cidr1, cidr2)
	if err != nil {
		t.Fatal(err)
	}
	if !subset {
		t.Fatalf("expected CIDR %q to be a subset of CIDR %q", cidr2, cidr1)
	}

	cidr1 = "192.168.27.29/16"
	cidr2 = "192.168.27.29/24"
	subset, err = Subset(cidr1, cidr2)
	if err != nil {
		t.Fatal(err)
	}
	if !subset {
		t.Fatalf("expected CIDR %q to be a subset of CIDR %q", cidr2, cidr1)
	}

	cidr1 = "192.168.27.29/24"
	cidr2 = "192.168.27.29/16"
	subset, err = Subset(cidr1, cidr2)
	if err != nil {
		t.Fatal(err)
	}
	if subset {
		t.Fatalf("expected CIDR %q to not be a subset of CIDR %q", cidr2, cidr1)
	}

	cidr1 = "192.168.0.128/25"
	cidr2 = "192.168.0.0/24"
	subset, err = Subset(cidr1, cidr2)
	if err != nil {
		t.Fatal(err)
	}
	if subset {
		t.Fatalf("expected CIDR %q to not be a subset of CIDR %q", cidr2, cidr1)
	}
	subset, err = Subset(cidr2, cidr1)
	if err != nil {
		t.Fatal(err)
	}
	if !subset {
		t.Fatalf("expected CIDR %q to be a subset of CIDR %q", cidr1, cidr2)
	}
}

func TestCIDRUtil_SubsetBlocks(t *testing.T) {
	cidrBlocks1 := []string{"192.168.27.29/16", "172.245.30.40/24", "10.20.30.40/30"}
	cidrBlocks2 := []string{"192.168.27.29/20", "172.245.30.40/25", "10.20.30.40/32"}

	subset, err := SubsetBlocks(cidrBlocks1, cidrBlocks2)
	if err != nil {
		t.Fatal(err)
	}
	if !subset {
		t.Fatalf("expected CIDR blocks %q to be a subset of CIDR blocks %q", cidrBlocks2, cidrBlocks1)
	}

	cidrBlocks1 = []string{"192.168.27.29/16", "172.245.30.40/25", "10.20.30.40/30"}
	cidrBlocks2 = []string{"192.168.27.29/20", "172.245.30.40/24", "10.20.30.40/32"}

	subset, err = SubsetBlocks(cidrBlocks1, cidrBlocks2)
	if err != nil {
		t.Fatal(err)
	}
	if subset {
		t.Fatalf("expected CIDR blocks %q to not be a subset of CIDR blocks %q", cidrBlocks2, cidrBlocks1)
	}
}

func TestCIDRUtil_RemoteAddrIsOk_NegativeTest(t *testing.T) {
	addr, err := sockaddr.NewSockAddr("127.0.0.1/8")
	if err != nil {
		t.Fatal(err)
	}
	boundCIDRs := []*sockaddr.SockAddrMarshaler{
		{addr},
	}
	if RemoteAddrIsOk("123.0.0.1", boundCIDRs) {
		t.Fatal("remote address of 123.0.0.1/2 should not be allowed for 127.0.0.1/8")
	}
}

func TestCIDRUtil_RemoteAddrIsOk_PositiveTest(t *testing.T) {
	addr, err := sockaddr.NewSockAddr("127.0.0.1/8")
	if err != nil {
		t.Fatal(err)
	}
	boundCIDRs := []*sockaddr.SockAddrMarshaler{
		{addr},
	}
	if !RemoteAddrIsOk("127.0.0.1", boundCIDRs) {
		t.Fatal("remote address of 127.0.0.1 should be allowed for 127.0.0.1/8")
	}
}
