package ecs

import (
	"time"

	"github.com/denverdino/aliyungo/common"
)

type SnatEntryStatus string

const (
	SnatEntryStatusPending   = SnatEntryStatus("Pending")
	SnatEntryStatusAvailable = SnatEntryStatus("Available")
)

type CreateSnatEntryArgs struct {
	RegionId        common.Region
	SnatTableId     string
	SourceVSwitchId string
	SnatIp          string
	SourceCIDR      string
}

type CreateSnatEntryResponse struct {
	common.Response
	SnatEntryId string
}

type SnatEntrySetType struct {
	RegionId        common.Region
	SnatEntryId     string
	SnatIp          string
	SnatTableId     string
	SourceCIDR      string
	SourceVSwitchId string
	Status          SnatEntryStatus
}

type DescribeSnatTableEntriesArgs struct {
	RegionId        common.Region
	SnatTableId     string
	SnatEntryId     string
	SnatEntryName   string
	SnatIp          string
	SourceCIDR      string
	SourceVSwitchId string
	common.Pagination
}

type DescribeSnatTableEntriesResponse struct {
	common.Response
	common.PaginationResult
	SnatTableEntries struct {
		SnatTableEntry []SnatEntrySetType
	}
}

type ModifySnatEntryArgs struct {
	RegionId      common.Region
	SnatTableId   string
	SnatEntryId   string
	SnatIp        string
	SnatEntryName string
}

type ModifySnatEntryResponse struct {
	common.Response
}

type DeleteSnatEntryArgs struct {
	RegionId    common.Region
	SnatTableId string
	SnatEntryId string
}

type DeleteSnatEntryResponse struct {
	common.Response
}

func (client *Client) CreateSnatEntry(args *CreateSnatEntryArgs) (resp *CreateSnatEntryResponse, err error) {
	response := CreateSnatEntryResponse{}
	err = client.Invoke("CreateSnatEntry", args, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (client *Client) DescribeSnatTableEntries(args *DescribeSnatTableEntriesArgs) (snatTableEntries []SnatEntrySetType,
	pagination *common.PaginationResult, err error) {
	response, err := client.DescribeSnatTableEntriesWithRaw(args)
	if err != nil {
		return nil, nil, err
	}

	return response.SnatTableEntries.SnatTableEntry, &response.PaginationResult, nil
}

func (client *Client) DescribeSnatTableEntriesWithRaw(args *DescribeSnatTableEntriesArgs) (response *DescribeSnatTableEntriesResponse, err error) {
	args.Validate()
	response = &DescribeSnatTableEntriesResponse{}

	err = client.Invoke("DescribeSnatTableEntries", args, response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *Client) ModifySnatEntry(args *ModifySnatEntryArgs) error {
	response := ModifySnatEntryResponse{}
	return client.Invoke("ModifySnatEntry", args, &response)
}

func (client *Client) DeleteSnatEntry(args *DeleteSnatEntryArgs) error {
	response := DeleteSnatEntryResponse{}
	err := client.Invoke("DeleteSnatEntry", args, &response)
	return err
}

// WaitForSnatEntryAvailable waits for SnatEntry to available status
func (client *Client) WaitForSnatEntryAvailable(regionId common.Region, snatTableId, snatEntryId string, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	args := &DescribeSnatTableEntriesArgs{
		RegionId:    regionId,
		SnatTableId: snatTableId,
		SnatEntryId: snatEntryId,
	}

	for {
		snatEntries, _, err := client.DescribeSnatTableEntries(args)
		if err != nil {
			return err
		}

		if len(snatEntries) == 0 {
			return common.GetClientErrorFromString("Not found")
		}
		if snatEntries[0].Status == SnatEntryStatusAvailable {
			break
		}

		timeout = timeout - DefaultWaitForInterval
		if timeout <= 0 {
			return common.GetClientErrorFromString("Timeout")
		}
		time.Sleep(DefaultWaitForInterval * time.Second)
	}
	return nil
}
