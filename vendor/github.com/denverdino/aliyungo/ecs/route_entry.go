package ecs

import "github.com/denverdino/aliyungo/common"

type DescribeRouteEntryListArgs struct {
	RegionId             string
	RouteTableId         string
	DestinationCidrBlock string
	IpVersion            string
	MaxResult            int
	NextHopId            string
	NextHopType          string
	NextToken            string
	RouteEntryId         string
	RouteEntryName       string
	RouteEntryType       string
}

type DescribeRouteEntryListResponse struct {
	common.Response
	NextToken   string
	RouteEntrys struct {
		RouteEntry []RouteEntry
	}
}

type RouteEntry struct {
	DestinationCidrBlock string
	IpVersion            string
	RouteEntryId         string
	RouteEntryName       string
	RouteTableId         string
	Status               string
	Type                 string
	NextHops             struct {
		NextHop []NextHop
	}
}

type NextHop struct {
	Enabled            int
	Weight             int
	NextHopId          string
	NextHopRegionId    string
	NextHopType        string
	NextHopRelatedInfo NextHopRelatedInfo
}

type NextHopRelatedInfo struct {
	RegionId     string
	InstanceId   string
	InstanceType string
}

// DescribeRouteEntryList describes route entries
//
func (client *Client) DescribeRouteEntryList(args *DescribeRouteEntryListArgs) (*DescribeRouteEntryListResponse, error) {
	response := &DescribeRouteEntryListResponse{}
	err := client.Invoke("DescribeRouteEntryList", args, &response)
	return response, err
}
