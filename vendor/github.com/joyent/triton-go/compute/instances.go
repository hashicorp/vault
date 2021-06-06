//
// Copyright 2019 Joyent, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package compute

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/joyent/triton-go/client"
	"github.com/joyent/triton-go/errors"
	pkgerrors "github.com/pkg/errors"
)

type InstancesClient struct {
	client *client.Client
}

const (
	CNSTagDisable    = "triton.cns.disable"
	CNSTagReversePTR = "triton.cns.reverse_ptr"
	CNSTagServices   = "triton.cns.services"
)

// InstanceCNS is a container for the CNS-specific attributes.  In the API these
// values are embedded within a Instance's Tags attribute, however they are
// exposed to the caller as their native types.
type InstanceCNS struct {
	Disable    bool
	ReversePTR string
	Services   []string
}

type InstanceVolume struct {
	Name       string `json:"name,omitempty"`
	Type       string `json:"type,omitempty"`
	Mode       string `json:"mode,omitempty"`
	Mountpoint string `json:"mountpoint,omitempty"`
}

type NetworkObject struct {
	IPv4UUID string   `json:"ipv4_uuid"`
	IPv4IPs  []string `json:"ipv4_ips,omitempty"`
}

type Instance struct {
	ID                 string                 `json:"id"`
	Name               string                 `json:"name"`
	Type               string                 `json:"type"`
	Brand              string                 `json:"brand"`
	State              string                 `json:"state"`
	Image              string                 `json:"image"`
	Memory             int                    `json:"memory"`
	Disk               int                    `json:"disk"`
	Metadata           map[string]string      `json:"metadata"`
	Tags               map[string]interface{} `json:"tags"`
	Created            time.Time              `json:"created"`
	Updated            time.Time              `json:"updated"`
	Docker             bool                   `json:"docker"`
	IPs                []string               `json:"ips"`
	Networks           []string               `json:"networks"`
	PrimaryIP          string                 `json:"primaryIp"`
	FirewallEnabled    bool                   `json:"firewall_enabled"`
	ComputeNode        string                 `json:"compute_node"`
	Package            string                 `json:"package"`
	DomainNames        []string               `json:"dns_names"`
	DeletionProtection bool                   `json:"deletion_protection"`
	CNS                InstanceCNS
}

// _Instance is a private facade over Instance that handles the necessary API
// overrides from VMAPI's machine endpoint(s).
type _Instance struct {
	Instance
	Tags map[string]interface{} `json:"tags"`
}

type NIC struct {
	IP      string `json:"ip"`
	MAC     string `json:"mac"`
	Primary bool   `json:"primary"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	State   string `json:"state"`
	Network string `json:"network"`
}

type GetInstanceInput struct {
	ID string
}

func (gmi *GetInstanceInput) Validate() error {
	if gmi.ID == "" {
		return fmt.Errorf("machine ID can not be empty")
	}

	return nil
}

func (c *InstancesClient) Count(ctx context.Context, input *ListInstancesInput) (int, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines")

	reqInputs := client.RequestInput{
		Method: http.MethodHead,
		Path:   fullPath,
		Query:  buildQueryFilter(input),
	}

	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return -1, pkgerrors.Wrap(err, "unable to get machines count")
	}

	if response == nil {
		return -1, pkgerrors.New("request to get machines count has empty response")
	}
	defer response.Body.Close()

	var result int

	if count := response.Header.Get("X-Resource-Count"); count != "" {
		value, err := strconv.Atoi(count)
		if err != nil {
			return -1, pkgerrors.Wrap(err, "unable to decode machines count response")
		}
		result = value
	}

	return result, nil
}

func (c *InstancesClient) Get(ctx context.Context, input *GetInstanceInput) (*Instance, error) {
	if err := input.Validate(); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to get machine")
	}

	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID)
	reqInputs := client.RequestInput{
		Method:       http.MethodGet,
		Path:         fullPath,
		PreserveGone: true,
	}
	response, reqErr := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if response == nil {
		return nil, pkgerrors.Wrap(reqErr, "unable to get machine")
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	if reqErr != nil {
		reqErr = pkgerrors.Wrap(reqErr, "unable to get machine")

		// If this is not a HTTP 410 Gone error, return it immediately to the caller.  Otherwise, we'll return it alongside the instance below.
		if response.StatusCode != http.StatusGone {
			return nil, reqErr
		}
	}

	var result *_Instance
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to parse JSON in get machine response")
	}

	native, err := result.toNative()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode get machine response")
	}

	// To remain compatible with the existing interface, we'll return both an error and an instance object in some cases; e.g., for HTTP 410 Gone responses for deleted instances.
	return native, reqErr
}

type ListInstancesInput struct {
	Brand       string
	Alias       string
	Name        string
	Image       string
	State       string
	Memory      uint16
	Limit       uint16
	Offset      uint16
	Tags        map[string]interface{} // query by arbitrary tags prefixed with "tag."
	Tombstone   bool
	Docker      bool
	Credentials bool
}

func buildQueryFilter(input *ListInstancesInput) *url.Values {
	query := &url.Values{}
	if input.Brand != "" {
		query.Set("brand", input.Brand)
	}
	if input.Name != "" {
		query.Set("name", input.Name)
	}
	if input.Image != "" {
		query.Set("image", input.Image)
	}
	if input.State != "" {
		query.Set("state", input.State)
	}
	if input.Memory >= 1 {
		query.Set("memory", fmt.Sprintf("%d", input.Memory))
	}
	if input.Limit >= 1 && input.Limit <= 1000 {
		query.Set("limit", fmt.Sprintf("%d", input.Limit))
	}
	if input.Offset >= 0 {
		query.Set("offset", fmt.Sprintf("%d", input.Offset))
	}
	if input.Tombstone {
		query.Set("tombstone", "true")
	}
	if input.Docker {
		query.Set("docker", "true")
	}
	if input.Credentials {
		query.Set("credentials", "true")
	}
	if input.Tags != nil {
		for k, v := range input.Tags {
			query.Set(fmt.Sprintf("tag.%s", k), v.(string))
		}
	}

	return query
}

func (c *InstancesClient) List(ctx context.Context, input *ListInstancesInput) ([]*Instance, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines")

	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
		Query:  buildQueryFilter(input),
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to list machines")
	}

	var results []*_Instance
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&results); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode list machines response")
	}

	machines := make([]*Instance, 0, len(results))
	for _, machineAPI := range results {
		native, err := machineAPI.toNative()
		if err != nil {
			return nil, pkgerrors.Wrap(err, "unable to decode list machines response")
		}
		machines = append(machines, native)
	}

	return machines, nil
}

type CreateInstanceInput struct {
	Name            string
	NamePrefix      string
	Package         string
	Image           string
	Networks        []string
	NetworkObjects  []NetworkObject
	Affinity        []string
	LocalityStrict  bool
	LocalityNear    []string
	LocalityFar     []string
	Metadata        map[string]string
	Tags            map[string]string //
	FirewallEnabled bool              //
	CNS             InstanceCNS
	Volumes         []InstanceVolume
}

func buildInstanceName(namePrefix string) string {
	h := sha1.New()
	io.WriteString(h, namePrefix+time.Now().UTC().String())
	return fmt.Sprintf("%s%s", namePrefix, hex.EncodeToString(h.Sum(nil))[:8])
}

func (input *CreateInstanceInput) toAPI() (map[string]interface{}, error) {
	const numExtraParams = 8
	result := make(map[string]interface{}, numExtraParams+len(input.Metadata)+len(input.Tags))

	result["firewall_enabled"] = input.FirewallEnabled

	if input.Name != "" {
		result["name"] = input.Name
	} else if input.NamePrefix != "" {
		result["name"] = buildInstanceName(input.NamePrefix)
	}

	if input.Package != "" {
		result["package"] = input.Package
	}

	if input.Image != "" {
		result["image"] = input.Image
	}

	// If we are passed []string from input.Networks that do not conflict with networks provided by NetworkObjects, add them to the request sent to CloudAPI
	var networks []NetworkObject

	if len(input.NetworkObjects) > 0 {
		networks = append(networks, input.NetworkObjects...)
	}

	for _, netuuid := range input.Networks {
		found := false

		for _, net := range networks {
			if net.IPv4UUID == netuuid {
				found = true
			}
		}

		if !found {
			networks = append(networks, NetworkObject{
				IPv4UUID: netuuid,
			})
		}
	}

	if len(networks) > 0 {
		result["networks"] = networks
	}

	if len(input.Volumes) > 0 {
		result["volumes"] = input.Volumes
	}

	// validate that affinity and locality are not included together
	hasAffinity := len(input.Affinity) > 0
	hasLocality := len(input.LocalityNear) > 0 || len(input.LocalityFar) > 0
	if hasAffinity && hasLocality {
		return nil, fmt.Errorf("Cannot include both Affinity and Locality")
	}

	// affinity takes precedence over locality regardless
	if len(input.Affinity) > 0 {
		result["affinity"] = input.Affinity
	} else {
		locality := struct {
			Strict bool     `json:"strict"`
			Near   []string `json:"near,omitempty"`
			Far    []string `json:"far,omitempty"`
		}{
			Strict: input.LocalityStrict,
			Near:   input.LocalityNear,
			Far:    input.LocalityFar,
		}
		result["locality"] = locality
	}

	for key, value := range input.Tags {
		result[fmt.Sprintf("tag.%s", key)] = value
	}

	// NOTE(justinwr): CNSTagServices needs to be a tag if available. No other
	// CNS tags will be handled at this time.
	input.CNS.toTags(result)
	if val, found := result[CNSTagServices]; found {
		result["tag."+CNSTagServices] = val
		delete(result, CNSTagServices)
	}

	for key, value := range input.Metadata {
		result[fmt.Sprintf("metadata.%s", key)] = value
	}

	return result, nil
}

func (c *InstancesClient) Create(ctx context.Context, input *CreateInstanceInput) (*Instance, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines")
	body, err := input.toAPI()
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to prepare for machine creation")
	}

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   body,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to create machine")
	}

	var result *Instance
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode create machine response")
	}

	return result, nil
}

type DeleteInstanceInput struct {
	ID string
}

func (c *InstancesClient) Delete(ctx context.Context, input *DeleteInstanceInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if response == nil {
		return pkgerrors.Wrap(err, "unable to delete machine")
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusGone {
		return nil
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to decode delete machine response")
	}

	return nil
}

type DeleteTagsInput struct {
	ID string
}

func (c *InstancesClient) DeleteTags(ctx context.Context, input *DeleteTagsInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID, "tags")
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return pkgerrors.Wrap(err, "unable to delete tags from machine")
	}
	if response == nil {
		return fmt.Errorf("DeleteTags request has empty response")
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	if response.StatusCode == http.StatusNotFound {
		return nil
	}

	return nil
}

type DeleteTagInput struct {
	ID  string
	Key string
}

func (c *InstancesClient) DeleteTag(ctx context.Context, input *DeleteTagInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID, "tags", input.Key)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return pkgerrors.Wrap(err, "unable to delete tag from machine")
	}
	if response == nil {
		return fmt.Errorf("DeleteTag request has empty response")
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	if response.StatusCode == http.StatusNotFound {
		return nil
	}

	return nil
}

type RenameInstanceInput struct {
	ID   string
	Name string
}

func (c *InstancesClient) Rename(ctx context.Context, input *RenameInstanceInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID)

	params := &url.Values{}
	params.Set("action", "rename")
	params.Set("name", input.Name)

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to rename machine")
	}

	return nil
}

type ReplaceTagsInput struct {
	ID   string
	Tags map[string]string
	CNS  InstanceCNS
}

// toAPI is used to join Tags and CNS tags into the same JSON object before
// sending an API request to the API gateway.
func (input ReplaceTagsInput) toAPI() map[string]interface{} {
	result := map[string]interface{}{}
	for key, value := range input.Tags {
		result[key] = value
	}
	input.CNS.toTags(result)
	return result
}

func (c *InstancesClient) ReplaceTags(ctx context.Context, input *ReplaceTagsInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID, "tags")
	reqInputs := client.RequestInput{
		Method: http.MethodPut,
		Path:   fullPath,
		Body:   input.toAPI(),
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to replace machine tags")
	}

	return nil
}

type AddTagsInput struct {
	ID   string
	Tags map[string]string
}

func (c *InstancesClient) AddTags(ctx context.Context, input *AddTagsInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID, "tags")
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   input.Tags,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to add tags to machine")
	}

	return nil
}

type GetTagInput struct {
	ID  string
	Key string
}

func (c *InstancesClient) GetTag(ctx context.Context, input *GetTagInput) (string, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID, "tags", input.Key)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if err != nil {
		return "", pkgerrors.Wrap(err, "unable to get tag")
	}
	if respReader != nil {
		defer respReader.Close()
	}

	var result string
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return "", pkgerrors.Wrap(err, "unable to decode get tag response")
	}

	return result, nil
}

type ListTagsInput struct {
	ID string
}

func (c *InstancesClient) ListTags(ctx context.Context, input *ListTagsInput) (map[string]interface{}, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID, "tags")
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to list machine tags")
	}

	var result map[string]interface{}
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, pkgerrors.Wrap(err, "unable decode list machine tags response")
	}

	_, tags := TagsExtractMeta(result)
	return tags, nil
}

type GetMetadataInput struct {
	ID  string
	Key string
}

// GetMetadata returns a single metadata entry associated with an instance.
func (c *InstancesClient) GetMetadata(ctx context.Context, input *GetMetadataInput) (string, error) {
	if input.Key == "" {
		return "", fmt.Errorf("Missing metadata Key from input: %s", input.Key)
	}

	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID, "metadata", input.Key)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return "", pkgerrors.Wrap(err, "unable to get machine metadata")
	}
	if response != nil {
		defer response.Body.Close()
	}
	if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusGone {
		return "", &errors.APIError{
			StatusCode: response.StatusCode,
			Code:       "ResourceNotFound",
		}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", pkgerrors.Wrap(err, "unable to decode get machine metadata response")
	}

	return fmt.Sprintf("%s", body), nil
}

type ListMetadataInput struct {
	ID          string
	Credentials bool
}

func (c *InstancesClient) ListMetadata(ctx context.Context, input *ListMetadataInput) (map[string]string, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID, "metadata")

	query := &url.Values{}
	if input.Credentials {
		query.Set("credentials", "true")
	}

	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
		Query:  query,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to list machine metadata")
	}

	var result map[string]string
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode list machine metadata response")
	}

	return result, nil
}

type UpdateMetadataInput struct {
	ID       string
	Metadata map[string]string
}

func (c *InstancesClient) UpdateMetadata(ctx context.Context, input *UpdateMetadataInput) (map[string]string, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID, "metadata")
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   input.Metadata,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to update machine metadata")
	}

	var result map[string]string
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode update machine metadata response")
	}

	return result, nil
}

type DeleteMetadataInput struct {
	ID  string
	Key string
}

// DeleteMetadata deletes a single metadata key from an instance
func (c *InstancesClient) DeleteMetadata(ctx context.Context, input *DeleteMetadataInput) error {
	if input.Key == "" {
		return fmt.Errorf("Missing metadata Key from input: %s", input.Key)
	}

	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID, "metadata", input.Key)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if err != nil {
		return pkgerrors.Wrap(err, "unable to delete machine metadata")
	}
	if respReader != nil {
		defer respReader.Close()
	}

	return nil
}

type DeleteAllMetadataInput struct {
	ID string
}

// DeleteAllMetadata deletes all metadata keys from this instance
func (c *InstancesClient) DeleteAllMetadata(ctx context.Context, input *DeleteAllMetadataInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID, "metadata")
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if err != nil {
		return pkgerrors.Wrap(err, "unable to delete all machine metadata")
	}
	if respReader != nil {
		defer respReader.Close()
	}

	return nil
}

type ResizeInstanceInput struct {
	ID      string
	Package string
}

func (c *InstancesClient) Resize(ctx context.Context, input *ResizeInstanceInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID)

	params := &url.Values{}
	params.Set("action", "resize")
	params.Set("package", input.Package)

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to resize machine")
	}

	return nil
}

type EnableFirewallInput struct {
	ID string
}

func (c *InstancesClient) EnableFirewall(ctx context.Context, input *EnableFirewallInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID)

	params := &url.Values{}
	params.Set("action", "enable_firewall")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to enable machine firewall")
	}

	return nil
}

type DisableFirewallInput struct {
	ID string
}

func (c *InstancesClient) DisableFirewall(ctx context.Context, input *DisableFirewallInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.ID)

	params := &url.Values{}
	params.Set("action", "disable_firewall")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to disable machine firewall")
	}

	return nil
}

type ListNICsInput struct {
	InstanceID string
}

func (c *InstancesClient) ListNICs(ctx context.Context, input *ListNICsInput) ([]*NIC, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.InstanceID, "nics")
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	respReader, err := c.client.ExecuteRequest(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to list machine NICs")
	}

	var result []*NIC
	decoder := json.NewDecoder(respReader)
	if err = decoder.Decode(&result); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode list machine NICs response")
	}

	return result, nil
}

type GetNICInput struct {
	InstanceID string
	MAC        string
}

func (c *InstancesClient) GetNIC(ctx context.Context, input *GetNICInput) (*NIC, error) {
	mac := strings.Replace(input.MAC, ":", "", -1)
	fullPath := path.Join("/", c.client.AccountName, "machines", input.InstanceID, "nics", mac)
	reqInputs := client.RequestInput{
		Method: http.MethodGet,
		Path:   fullPath,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to get machine NIC")
	}
	if response != nil {
		defer response.Body.Close()
	}
	switch response.StatusCode {
	case http.StatusNotFound:
		return nil, &errors.APIError{
			StatusCode: response.StatusCode,
			Code:       "ResourceNotFound",
		}
	}

	var result *NIC
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&result); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode get machine NIC response")
	}

	return result, nil
}

type AddNICInput struct {
	InstanceID    string
	Network       string
	NetworkObject NetworkObject
}

// toAPI is used to build up the JSON Object to send to the API gateway.  It
// also will resolve the scenario where a user provides both a NetworkObject
// and a Network. If both are provided, NetworkObject wins.
func (input AddNICInput) toAPI() map[string]interface{} {
	result := map[string]interface{}{}

	var network NetworkObject

	if input.NetworkObject.IPv4UUID != "" {
		network = input.NetworkObject
	} else {
		network = NetworkObject{
			IPv4UUID: input.Network,
		}
	}

	result["network"] = network

	return result
}

// AddNIC asynchronously adds a NIC to a given instance.  If a NIC for a given
// network already exists, a ResourceFound error will be returned.  The status
// of the addition of a NIC can be polled by calling GetNIC()'s and testing NIC
// until its state is set to "running".  Only one NIC per network may exist.
// Warning: this operation causes the instance to restart.
func (c *InstancesClient) AddNIC(ctx context.Context, input *AddNICInput) (*NIC, error) {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.InstanceID, "nics")
	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Body:   input.toAPI(),
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "unable to add NIC to machine")
	}
	if response != nil {
		defer response.Body.Close()
	}
	switch response.StatusCode {
	case http.StatusFound:
		return nil, &errors.APIError{
			StatusCode: response.StatusCode,
			Code:       "ResourceFound",
			Message:    response.Header.Get("Location"),
		}
	}

	var result *NIC
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&result); err != nil {
		return nil, pkgerrors.Wrap(err, "unable to decode add NIC to machine response")
	}

	return result, nil
}

type RemoveNICInput struct {
	InstanceID string
	MAC        string
}

// RemoveNIC removes a given NIC from a machine asynchronously.  The status of
// the removal can be polled via GetNIC().  When GetNIC() returns a 404, the NIC
// has been removed from the instance.  Warning: this operation causes the
// machine to restart.
func (c *InstancesClient) RemoveNIC(ctx context.Context, input *RemoveNICInput) error {
	mac := strings.Replace(input.MAC, ":", "", -1)
	fullPath := path.Join("/", c.client.AccountName, "machines", input.InstanceID, "nics", mac)
	reqInputs := client.RequestInput{
		Method: http.MethodDelete,
		Path:   fullPath,
	}
	response, err := c.client.ExecuteRequestRaw(ctx, reqInputs)
	if err != nil {
		return pkgerrors.Wrap(err, "unable to remove NIC from machine")
	}
	if response == nil {
		return pkgerrors.Wrap(err, "unable to remove NIC from machine")
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	switch response.StatusCode {
	case http.StatusNotFound:
		return &errors.APIError{
			StatusCode: response.StatusCode,
			Code:       "ResourceNotFound",
		}
	}

	return nil
}

type StopInstanceInput struct {
	InstanceID string
}

func (c *InstancesClient) Stop(ctx context.Context, input *StopInstanceInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.InstanceID)

	params := &url.Values{}
	params.Set("action", "stop")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to stop machine")
	}

	return nil
}

type StartInstanceInput struct {
	InstanceID string
}

func (c *InstancesClient) Start(ctx context.Context, input *StartInstanceInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.InstanceID)

	params := &url.Values{}
	params.Set("action", "start")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to start machine")
	}

	return nil
}

type RebootInstanceInput struct {
	InstanceID string
}

func (c *InstancesClient) Reboot(ctx context.Context, input *RebootInstanceInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.InstanceID)

	params := &url.Values{}
	params.Set("action", "reboot")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to reboot machine")
	}

	return nil
}

type EnableDeletionProtectionInput struct {
	InstanceID string
}

func (c *InstancesClient) EnableDeletionProtection(ctx context.Context, input *EnableDeletionProtectionInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.InstanceID)

	params := &url.Values{}
	params.Set("action", "enable_deletion_protection")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to enable deletion protection")
	}

	return nil
}

type DisableDeletionProtectionInput struct {
	InstanceID string
}

func (c *InstancesClient) DisableDeletionProtection(ctx context.Context, input *DisableDeletionProtectionInput) error {
	fullPath := path.Join("/", c.client.AccountName, "machines", input.InstanceID)

	params := &url.Values{}
	params.Set("action", "disable_deletion_protection")

	reqInputs := client.RequestInput{
		Method: http.MethodPost,
		Path:   fullPath,
		Query:  params,
	}
	respReader, err := c.client.ExecuteRequestURIParams(ctx, reqInputs)
	if respReader != nil {
		defer respReader.Close()
	}
	if err != nil {
		return pkgerrors.Wrap(err, "unable to disable deletion protection")
	}

	return nil
}

var reservedInstanceCNSTags = map[string]struct{}{
	CNSTagDisable:    {},
	CNSTagReversePTR: {},
	CNSTagServices:   {},
}

// TagsExtractMeta extracts all of the misc parameters from Tags and returns a
// clean CNS and Tags struct.
func TagsExtractMeta(tags map[string]interface{}) (InstanceCNS, map[string]interface{}) {
	nativeCNS := InstanceCNS{}
	nativeTags := make(map[string]interface{}, len(tags))
	for k, raw := range tags {
		if _, found := reservedInstanceCNSTags[k]; found {
			switch k {
			case CNSTagDisable:
				b := raw.(bool)
				nativeCNS.Disable = b
			case CNSTagReversePTR:
				s := raw.(string)
				nativeCNS.ReversePTR = s
			case CNSTagServices:
				nativeCNS.Services = strings.Split(raw.(string), ",")
			default:
				// TODO(seanc@): should assert, logic fail
			}
		} else {
			nativeTags[k] = raw
		}
	}

	return nativeCNS, nativeTags
}

// toNative() exports a given _Instance (API representation) to its native object
// format.
func (api *_Instance) toNative() (*Instance, error) {
	m := Instance(api.Instance)
	m.CNS, m.Tags = TagsExtractMeta(api.Tags)
	return &m, nil
}

// toTags() injects its state information into a Tags map suitable for use to
// submit an API call to the vmapi machine endpoint
func (cns *InstanceCNS) toTags(m map[string]interface{}) {
	if cns.Disable {
		// NOTE(justinwr): The JSON encoder and API require the CNSTagDisable
		// attribute to be an actual boolean, not a bool string.
		m[CNSTagDisable] = cns.Disable
	}
	if cns.ReversePTR != "" {
		m[CNSTagReversePTR] = cns.ReversePTR
	}
	if len(cns.Services) > 0 {
		m[CNSTagServices] = strings.Join(cns.Services, ",")
	}
}
