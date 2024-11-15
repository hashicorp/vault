package godo

import (
	"context"
	"net/http"
)

// SizesService is an interface for interfacing with the size
// endpoints of the DigitalOcean API
// See: https://docs.digitalocean.com/reference/api/api-reference/#tag/Sizes
type SizesService interface {
	List(context.Context, *ListOptions) ([]Size, *Response, error)
}

// SizesServiceOp handles communication with the size related methods of the
// DigitalOcean API.
type SizesServiceOp struct {
	client *Client
}

var _ SizesService = &SizesServiceOp{}

// Size represents a DigitalOcean Size
type Size struct {
	Slug         string     `json:"slug,omitempty"`
	Memory       int        `json:"memory,omitempty"`
	Vcpus        int        `json:"vcpus,omitempty"`
	Disk         int        `json:"disk,omitempty"`
	PriceMonthly float64    `json:"price_monthly,omitempty"`
	PriceHourly  float64    `json:"price_hourly,omitempty"`
	Regions      []string   `json:"regions,omitempty"`
	Available    bool       `json:"available,omitempty"`
	Transfer     float64    `json:"transfer,omitempty"`
	Description  string     `json:"description,omitempty"`
	GPUInfo      *GPUInfo   `json:"gpu_info,omitempty"`
	DiskInfo     []DiskInfo `json:"disk_info,omitempty"`
}

// DiskInfo containing information about the disks available to Droplets created
// with this size.
type DiskInfo struct {
	Type string    `json:"type,omitempty"`
	Size *DiskSize `json:"size,omitempty"`
}

// DiskSize provides information about the size of a disk.
type DiskSize struct {
	Amount int    `json:"amount,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

// GPUInfo provides information about the GPU available to Droplets created with this size.
type GPUInfo struct {
	Count int    `json:"count,omitempty"`
	VRAM  *VRAM  `json:"vram,omitempty"`
	Model string `json:"model,omitempty"`
}

// VRAM provides information about the amount of VRAM available to the GPU.
type VRAM struct {
	Amount int    `json:"amount,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

func (s Size) String() string {
	return Stringify(s)
}

type sizesRoot struct {
	Sizes []Size
	Links *Links `json:"links"`
	Meta  *Meta  `json:"meta"`
}

// List all images
func (s *SizesServiceOp) List(ctx context.Context, opt *ListOptions) ([]Size, *Response, error) {
	path := "v2/sizes"
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(sizesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.Sizes, resp, err
}
