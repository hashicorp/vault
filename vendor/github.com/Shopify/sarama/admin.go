package sarama

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// ClusterAdmin is the administrative client for Kafka, which supports managing and inspecting topics,
// brokers, configurations and ACLs. The minimum broker version required is 0.10.0.0.
// Methods with stricter requirements will specify the minimum broker version required.
// You MUST call Close() on a client to avoid leaks
type ClusterAdmin interface {
	// Creates a new topic. This operation is supported by brokers with version 0.10.1.0 or higher.
	// It may take several seconds after CreateTopic returns success for all the brokers
	// to become aware that the topic has been created. During this time, listTopics
	// may not return information about the new topic.The validateOnly option is supported from version 0.10.2.0.
	CreateTopic(topic string, detail *TopicDetail, validateOnly bool) error

	// List the topics available in the cluster with the default options.
	ListTopics() (map[string]TopicDetail, error)

	// Describe some topics in the cluster.
	DescribeTopics(topics []string) (metadata []*TopicMetadata, err error)

	// Delete a topic. It may take several seconds after the DeleteTopic to returns success
	// and for all the brokers to become aware that the topics are gone.
	// During this time, listTopics  may continue to return information about the deleted topic.
	// If delete.topic.enable is false on the brokers, deleteTopic will mark
	// the topic for deletion, but not actually delete them.
	// This operation is supported by brokers with version 0.10.1.0 or higher.
	DeleteTopic(topic string) error

	// Increase the number of partitions of the topics  according to the corresponding values.
	// If partitions are increased for a topic that has a key, the partition logic or ordering of
	// the messages will be affected. It may take several seconds after this method returns
	// success for all the brokers to become aware that the partitions have been created.
	// During this time, ClusterAdmin#describeTopics may not return information about the
	// new partitions. This operation is supported by brokers with version 1.0.0 or higher.
	CreatePartitions(topic string, count int32, assignment [][]int32, validateOnly bool) error

	// Alter the replica assignment for partitions.
	// This operation is supported by brokers with version 2.4.0.0 or higher.
	AlterPartitionReassignments(topic string, assignment [][]int32) error

	// Provides info on ongoing partitions replica reassignments.
	// This operation is supported by brokers with version 2.4.0.0 or higher.
	ListPartitionReassignments(topics string, partitions []int32) (topicStatus map[string]map[int32]*PartitionReplicaReassignmentsStatus, err error)

	// Delete records whose offset is smaller than the given offset of the corresponding partition.
	// This operation is supported by brokers with version 0.11.0.0 or higher.
	DeleteRecords(topic string, partitionOffsets map[int32]int64) error

	// Get the configuration for the specified resources.
	// The returned configuration includes default values and the Default is true
	// can be used to distinguish them from user supplied values.
	// Config entries where ReadOnly is true cannot be updated.
	// The value of config entries where Sensitive is true is always nil so
	// sensitive information is not disclosed.
	// This operation is supported by brokers with version 0.11.0.0 or higher.
	DescribeConfig(resource ConfigResource) ([]ConfigEntry, error)

	// Update the configuration for the specified resources with the default options.
	// This operation is supported by brokers with version 0.11.0.0 or higher.
	// The resources with their configs (topic is the only resource type with configs
	// that can be updated currently Updates are not transactional so they may succeed
	// for some resources while fail for others. The configs for a particular resource are updated automatically.
	AlterConfig(resourceType ConfigResourceType, name string, entries map[string]*string, validateOnly bool) error

	// Creates access control lists (ACLs) which are bound to specific resources.
	// This operation is not transactional so it may succeed for some ACLs while fail for others.
	// If you attempt to add an ACL that duplicates an existing ACL, no error will be raised, but
	// no changes will be made. This operation is supported by brokers with version 0.11.0.0 or higher.
	CreateACL(resource Resource, acl Acl) error

	// Lists access control lists (ACLs) according to the supplied filter.
	// it may take some time for changes made by createAcls or deleteAcls to be reflected in the output of ListAcls
	// This operation is supported by brokers with version 0.11.0.0 or higher.
	ListAcls(filter AclFilter) ([]ResourceAcls, error)

	// Deletes access control lists (ACLs) according to the supplied filters.
	// This operation is not transactional so it may succeed for some ACLs while fail for others.
	// This operation is supported by brokers with version 0.11.0.0 or higher.
	DeleteACL(filter AclFilter, validateOnly bool) ([]MatchingAcl, error)

	// List the consumer groups available in the cluster.
	ListConsumerGroups() (map[string]string, error)

	// Describe the given consumer groups.
	DescribeConsumerGroups(groups []string) ([]*GroupDescription, error)

	// List the consumer group offsets available in the cluster.
	ListConsumerGroupOffsets(group string, topicPartitions map[string][]int32) (*OffsetFetchResponse, error)

	// Delete a consumer group.
	DeleteConsumerGroup(group string) error

	// Get information about the nodes in the cluster
	DescribeCluster() (brokers []*Broker, controllerID int32, err error)

	// Get information about all log directories on the given set of brokers
	DescribeLogDirs(brokers []int32) (map[int32][]DescribeLogDirsResponseDirMetadata, error)

	// Close shuts down the admin and closes underlying client.
	Close() error
}

type clusterAdmin struct {
	client Client
	conf   *Config
}

// NewClusterAdmin creates a new ClusterAdmin using the given broker addresses and configuration.
func NewClusterAdmin(addrs []string, conf *Config) (ClusterAdmin, error) {
	client, err := NewClient(addrs, conf)
	if err != nil {
		return nil, err
	}
	return NewClusterAdminFromClient(client)
}

// NewClusterAdminFromClient creates a new ClusterAdmin using the given client.
// Note that underlying client will also be closed on admin's Close() call.
func NewClusterAdminFromClient(client Client) (ClusterAdmin, error) {
	//make sure we can retrieve the controller
	_, err := client.Controller()
	if err != nil {
		return nil, err
	}

	ca := &clusterAdmin{
		client: client,
		conf:   client.Config(),
	}
	return ca, nil
}

func (ca *clusterAdmin) Close() error {
	return ca.client.Close()
}

func (ca *clusterAdmin) Controller() (*Broker, error) {
	return ca.client.Controller()
}

func (ca *clusterAdmin) refreshController() (*Broker, error) {
	return ca.client.RefreshController()
}

// isErrNoController returns `true` if the given error type unwraps to an
// `ErrNotController` response from Kafka
func isErrNoController(err error) bool {
	switch e := err.(type) {
	case *TopicError:
		return e.Err == ErrNotController
	case *TopicPartitionError:
		return e.Err == ErrNotController
	case KError:
		return e == ErrNotController
	}
	return false
}

// retryOnError will repeatedly call the given (error-returning) func in the
// case that its response is non-nil and retriable (as determined by the
// provided retriable func) up to the maximum number of tries permitted by
// the admin client configuration
func (ca *clusterAdmin) retryOnError(retriable func(error) bool, fn func() error) error {
	var err error
	for attempt := 0; attempt < ca.conf.Admin.Retry.Max; attempt++ {
		err = fn()
		if err == nil || !retriable(err) {
			return err
		}
		Logger.Printf(
			"admin/request retrying after %dms... (%d attempts remaining)\n",
			ca.conf.Admin.Retry.Backoff/time.Millisecond, ca.conf.Admin.Retry.Max-attempt)
		time.Sleep(ca.conf.Admin.Retry.Backoff)
		continue
	}
	return err
}

func (ca *clusterAdmin) CreateTopic(topic string, detail *TopicDetail, validateOnly bool) error {
	if topic == "" {
		return ErrInvalidTopic
	}

	if detail == nil {
		return errors.New("you must specify topic details")
	}

	topicDetails := make(map[string]*TopicDetail)
	topicDetails[topic] = detail

	request := &CreateTopicsRequest{
		TopicDetails: topicDetails,
		ValidateOnly: validateOnly,
		Timeout:      ca.conf.Admin.Timeout,
	}

	if ca.conf.Version.IsAtLeast(V0_11_0_0) {
		request.Version = 1
	}
	if ca.conf.Version.IsAtLeast(V1_0_0_0) {
		request.Version = 2
	}

	return ca.retryOnError(isErrNoController, func() error {
		b, err := ca.Controller()
		if err != nil {
			return err
		}

		rsp, err := b.CreateTopics(request)
		if err != nil {
			return err
		}

		topicErr, ok := rsp.TopicErrors[topic]
		if !ok {
			return ErrIncompleteResponse
		}

		if topicErr.Err != ErrNoError {
			if topicErr.Err == ErrNotController {
				_, _ = ca.refreshController()
			}
			return topicErr
		}

		return nil
	})
}

func (ca *clusterAdmin) DescribeTopics(topics []string) (metadata []*TopicMetadata, err error) {
	controller, err := ca.Controller()
	if err != nil {
		return nil, err
	}

	request := &MetadataRequest{
		Topics:                 topics,
		AllowAutoTopicCreation: false,
	}

	if ca.conf.Version.IsAtLeast(V1_0_0_0) {
		request.Version = 5
	} else if ca.conf.Version.IsAtLeast(V0_11_0_0) {
		request.Version = 4
	}

	response, err := controller.GetMetadata(request)
	if err != nil {
		return nil, err
	}
	return response.Topics, nil
}

func (ca *clusterAdmin) DescribeCluster() (brokers []*Broker, controllerID int32, err error) {
	controller, err := ca.Controller()
	if err != nil {
		return nil, int32(0), err
	}

	request := &MetadataRequest{
		Topics: []string{},
	}

	if ca.conf.Version.IsAtLeast(V0_10_0_0) {
		request.Version = 1
	}

	response, err := controller.GetMetadata(request)
	if err != nil {
		return nil, int32(0), err
	}

	return response.Brokers, response.ControllerID, nil
}

func (ca *clusterAdmin) findBroker(id int32) (*Broker, error) {
	brokers := ca.client.Brokers()
	for _, b := range brokers {
		if b.ID() == id {
			return b, nil
		}
	}
	return nil, fmt.Errorf("could not find broker id %d", id)
}

func (ca *clusterAdmin) findAnyBroker() (*Broker, error) {
	brokers := ca.client.Brokers()
	if len(brokers) > 0 {
		index := rand.Intn(len(brokers))
		return brokers[index], nil
	}
	return nil, errors.New("no available broker")
}

func (ca *clusterAdmin) ListTopics() (map[string]TopicDetail, error) {
	// In order to build TopicDetails we need to first get the list of all
	// topics using a MetadataRequest and then get their configs using a
	// DescribeConfigsRequest request. To avoid sending many requests to the
	// broker, we use a single DescribeConfigsRequest.

	// Send the all-topic MetadataRequest
	b, err := ca.findAnyBroker()
	if err != nil {
		return nil, err
	}
	_ = b.Open(ca.client.Config())

	metadataReq := &MetadataRequest{}
	metadataResp, err := b.GetMetadata(metadataReq)
	if err != nil {
		return nil, err
	}

	topicsDetailsMap := make(map[string]TopicDetail)

	var describeConfigsResources []*ConfigResource

	for _, topic := range metadataResp.Topics {
		topicDetails := TopicDetail{
			NumPartitions: int32(len(topic.Partitions)),
		}
		if len(topic.Partitions) > 0 {
			topicDetails.ReplicaAssignment = map[int32][]int32{}
			for _, partition := range topic.Partitions {
				topicDetails.ReplicaAssignment[partition.ID] = partition.Replicas
			}
			topicDetails.ReplicationFactor = int16(len(topic.Partitions[0].Replicas))
		}
		topicsDetailsMap[topic.Name] = topicDetails

		// we populate the resources we want to describe from the MetadataResponse
		topicResource := ConfigResource{
			Type: TopicResource,
			Name: topic.Name,
		}
		describeConfigsResources = append(describeConfigsResources, &topicResource)
	}

	// Send the DescribeConfigsRequest
	describeConfigsReq := &DescribeConfigsRequest{
		Resources: describeConfigsResources,
	}

	if ca.conf.Version.IsAtLeast(V1_1_0_0) {
		describeConfigsReq.Version = 1
	}

	if ca.conf.Version.IsAtLeast(V2_0_0_0) {
		describeConfigsReq.Version = 2
	}

	describeConfigsResp, err := b.DescribeConfigs(describeConfigsReq)
	if err != nil {
		return nil, err
	}

	for _, resource := range describeConfigsResp.Resources {
		topicDetails := topicsDetailsMap[resource.Name]
		topicDetails.ConfigEntries = make(map[string]*string)

		for _, entry := range resource.Configs {
			// only include non-default non-sensitive config
			// (don't actually think topic config will ever be sensitive)
			if entry.Default || entry.Sensitive {
				continue
			}
			topicDetails.ConfigEntries[entry.Name] = &entry.Value
		}

		topicsDetailsMap[resource.Name] = topicDetails
	}

	return topicsDetailsMap, nil
}

func (ca *clusterAdmin) DeleteTopic(topic string) error {
	if topic == "" {
		return ErrInvalidTopic
	}

	request := &DeleteTopicsRequest{
		Topics:  []string{topic},
		Timeout: ca.conf.Admin.Timeout,
	}

	if ca.conf.Version.IsAtLeast(V0_11_0_0) {
		request.Version = 1
	}

	return ca.retryOnError(isErrNoController, func() error {
		b, err := ca.Controller()
		if err != nil {
			return err
		}

		rsp, err := b.DeleteTopics(request)
		if err != nil {
			return err
		}

		topicErr, ok := rsp.TopicErrorCodes[topic]
		if !ok {
			return ErrIncompleteResponse
		}

		if topicErr != ErrNoError {
			if topicErr == ErrNotController {
				_, _ = ca.refreshController()
			}
			return topicErr
		}

		return nil
	})
}

func (ca *clusterAdmin) CreatePartitions(topic string, count int32, assignment [][]int32, validateOnly bool) error {
	if topic == "" {
		return ErrInvalidTopic
	}

	topicPartitions := make(map[string]*TopicPartition)
	topicPartitions[topic] = &TopicPartition{Count: count, Assignment: assignment}

	request := &CreatePartitionsRequest{
		TopicPartitions: topicPartitions,
		Timeout:         ca.conf.Admin.Timeout,
	}

	return ca.retryOnError(isErrNoController, func() error {
		b, err := ca.Controller()
		if err != nil {
			return err
		}

		rsp, err := b.CreatePartitions(request)
		if err != nil {
			return err
		}

		topicErr, ok := rsp.TopicPartitionErrors[topic]
		if !ok {
			return ErrIncompleteResponse
		}

		if topicErr.Err != ErrNoError {
			if topicErr.Err == ErrNotController {
				_, _ = ca.refreshController()
			}
			return topicErr
		}

		return nil
	})
}

func (ca *clusterAdmin) AlterPartitionReassignments(topic string, assignment [][]int32) error {
	if topic == "" {
		return ErrInvalidTopic
	}

	request := &AlterPartitionReassignmentsRequest{
		TimeoutMs: int32(60000),
		Version:   int16(0),
	}

	for i := 0; i < len(assignment); i++ {
		request.AddBlock(topic, int32(i), assignment[i])
	}

	return ca.retryOnError(isErrNoController, func() error {
		b, err := ca.Controller()
		if err != nil {
			return err
		}

		errs := make([]error, 0)

		rsp, err := b.AlterPartitionReassignments(request)

		if err != nil {
			errs = append(errs, err)
		} else {
			if rsp.ErrorCode > 0 {
				errs = append(errs, errors.New(rsp.ErrorCode.Error()))
			}

			for topic, topicErrors := range rsp.Errors {
				for partition, partitionError := range topicErrors {
					if partitionError.errorCode != ErrNoError {
						errStr := fmt.Sprintf("[%s-%d]: %s", topic, partition, partitionError.errorCode.Error())
						errs = append(errs, errors.New(errStr))
					}
				}
			}
		}

		if len(errs) > 0 {
			return ErrReassignPartitions{MultiError{&errs}}
		}

		return nil
	})
}

func (ca *clusterAdmin) ListPartitionReassignments(topic string, partitions []int32) (topicStatus map[string]map[int32]*PartitionReplicaReassignmentsStatus, err error) {
	if topic == "" {
		return nil, ErrInvalidTopic
	}

	request := &ListPartitionReassignmentsRequest{
		TimeoutMs: int32(60000),
		Version:   int16(0),
	}

	request.AddBlock(topic, partitions)

	b, err := ca.Controller()
	if err != nil {
		return nil, err
	}
	_ = b.Open(ca.client.Config())

	rsp, err := b.ListPartitionReassignments(request)

	if err == nil && rsp != nil {
		return rsp.TopicStatus, nil
	} else {
		return nil, err
	}
}

func (ca *clusterAdmin) DeleteRecords(topic string, partitionOffsets map[int32]int64) error {
	if topic == "" {
		return ErrInvalidTopic
	}
	partitionPerBroker := make(map[*Broker][]int32)
	for partition := range partitionOffsets {
		broker, err := ca.client.Leader(topic, partition)
		if err != nil {
			return err
		}
		if _, ok := partitionPerBroker[broker]; ok {
			partitionPerBroker[broker] = append(partitionPerBroker[broker], partition)
		} else {
			partitionPerBroker[broker] = []int32{partition}
		}
	}
	errs := make([]error, 0)
	for broker, partitions := range partitionPerBroker {
		topics := make(map[string]*DeleteRecordsRequestTopic)
		recordsToDelete := make(map[int32]int64)
		for _, p := range partitions {
			recordsToDelete[p] = partitionOffsets[p]
		}
		topics[topic] = &DeleteRecordsRequestTopic{PartitionOffsets: recordsToDelete}
		request := &DeleteRecordsRequest{
			Topics:  topics,
			Timeout: ca.conf.Admin.Timeout,
		}

		rsp, err := broker.DeleteRecords(request)
		if err != nil {
			errs = append(errs, err)
		} else {
			deleteRecordsResponseTopic, ok := rsp.Topics[topic]
			if !ok {
				errs = append(errs, ErrIncompleteResponse)
			} else {
				for _, deleteRecordsResponsePartition := range deleteRecordsResponseTopic.Partitions {
					if deleteRecordsResponsePartition.Err != ErrNoError {
						errs = append(errs, errors.New(deleteRecordsResponsePartition.Err.Error()))
					}
				}
			}
		}
	}
	if len(errs) > 0 {
		return ErrDeleteRecords{MultiError{&errs}}
	}
	//todo since we are dealing with couple of partitions it would be good if we return slice of errors
	//for each partition instead of one error
	return nil
}

// Returns a bool indicating whether the resource request needs to go to a
// specific broker
func dependsOnSpecificNode(resource ConfigResource) bool {
	return (resource.Type == BrokerResource && resource.Name != "") ||
		resource.Type == BrokerLoggerResource
}

func (ca *clusterAdmin) DescribeConfig(resource ConfigResource) ([]ConfigEntry, error) {
	var entries []ConfigEntry
	var resources []*ConfigResource
	resources = append(resources, &resource)

	request := &DescribeConfigsRequest{
		Resources: resources,
	}

	if ca.conf.Version.IsAtLeast(V1_1_0_0) {
		request.Version = 1
	}

	if ca.conf.Version.IsAtLeast(V2_0_0_0) {
		request.Version = 2
	}

	var (
		b   *Broker
		err error
	)

	// DescribeConfig of broker/broker logger must be sent to the broker in question
	if dependsOnSpecificNode(resource) {
		id, _ := strconv.Atoi(resource.Name)
		b, err = ca.findBroker(int32(id))
	} else {
		b, err = ca.findAnyBroker()
	}
	if err != nil {
		return nil, err
	}

	_ = b.Open(ca.client.Config())
	rsp, err := b.DescribeConfigs(request)
	if err != nil {
		return nil, err
	}

	for _, rspResource := range rsp.Resources {
		if rspResource.Name == resource.Name {
			if rspResource.ErrorMsg != "" {
				return nil, errors.New(rspResource.ErrorMsg)
			}
			if rspResource.ErrorCode != 0 {
				return nil, KError(rspResource.ErrorCode)
			}
			for _, cfgEntry := range rspResource.Configs {
				entries = append(entries, *cfgEntry)
			}
		}
	}
	return entries, nil
}

func (ca *clusterAdmin) AlterConfig(resourceType ConfigResourceType, name string, entries map[string]*string, validateOnly bool) error {
	var resources []*AlterConfigsResource
	resources = append(resources, &AlterConfigsResource{
		Type:          resourceType,
		Name:          name,
		ConfigEntries: entries,
	})

	request := &AlterConfigsRequest{
		Resources:    resources,
		ValidateOnly: validateOnly,
	}

	var (
		b   *Broker
		err error
	)

	// AlterConfig of broker/broker logger must be sent to the broker in question
	if dependsOnSpecificNode(ConfigResource{Name: name, Type: resourceType}) {
		id, _ := strconv.Atoi(name)
		b, err = ca.findBroker(int32(id))
	} else {
		b, err = ca.findAnyBroker()
	}
	if err != nil {
		return err
	}

	_ = b.Open(ca.client.Config())
	rsp, err := b.AlterConfigs(request)
	if err != nil {
		return err
	}

	for _, rspResource := range rsp.Resources {
		if rspResource.Name == name {
			if rspResource.ErrorMsg != "" {
				return errors.New(rspResource.ErrorMsg)
			}
			if rspResource.ErrorCode != 0 {
				return KError(rspResource.ErrorCode)
			}
		}
	}
	return nil
}

func (ca *clusterAdmin) CreateACL(resource Resource, acl Acl) error {
	var acls []*AclCreation
	acls = append(acls, &AclCreation{resource, acl})
	request := &CreateAclsRequest{AclCreations: acls}

	if ca.conf.Version.IsAtLeast(V2_0_0_0) {
		request.Version = 1
	}

	b, err := ca.Controller()
	if err != nil {
		return err
	}

	_, err = b.CreateAcls(request)
	return err
}

func (ca *clusterAdmin) ListAcls(filter AclFilter) ([]ResourceAcls, error) {
	request := &DescribeAclsRequest{AclFilter: filter}

	if ca.conf.Version.IsAtLeast(V2_0_0_0) {
		request.Version = 1
	}

	b, err := ca.Controller()
	if err != nil {
		return nil, err
	}

	rsp, err := b.DescribeAcls(request)
	if err != nil {
		return nil, err
	}

	var lAcls []ResourceAcls
	for _, rAcl := range rsp.ResourceAcls {
		lAcls = append(lAcls, *rAcl)
	}
	return lAcls, nil
}

func (ca *clusterAdmin) DeleteACL(filter AclFilter, validateOnly bool) ([]MatchingAcl, error) {
	var filters []*AclFilter
	filters = append(filters, &filter)
	request := &DeleteAclsRequest{Filters: filters}

	if ca.conf.Version.IsAtLeast(V2_0_0_0) {
		request.Version = 1
	}

	b, err := ca.Controller()
	if err != nil {
		return nil, err
	}

	rsp, err := b.DeleteAcls(request)
	if err != nil {
		return nil, err
	}

	var mAcls []MatchingAcl
	for _, fr := range rsp.FilterResponses {
		for _, mACL := range fr.MatchingAcls {
			mAcls = append(mAcls, *mACL)
		}
	}
	return mAcls, nil
}

func (ca *clusterAdmin) DescribeConsumerGroups(groups []string) (result []*GroupDescription, err error) {
	groupsPerBroker := make(map[*Broker][]string)

	for _, group := range groups {
		controller, err := ca.client.Coordinator(group)
		if err != nil {
			return nil, err
		}
		groupsPerBroker[controller] = append(groupsPerBroker[controller], group)
	}

	for broker, brokerGroups := range groupsPerBroker {
		response, err := broker.DescribeGroups(&DescribeGroupsRequest{
			Groups: brokerGroups,
		})
		if err != nil {
			return nil, err
		}

		result = append(result, response.Groups...)
	}
	return result, nil
}

func (ca *clusterAdmin) ListConsumerGroups() (allGroups map[string]string, err error) {
	allGroups = make(map[string]string)

	// Query brokers in parallel, since we have to query *all* brokers
	brokers := ca.client.Brokers()
	groupMaps := make(chan map[string]string, len(brokers))
	errors := make(chan error, len(brokers))
	wg := sync.WaitGroup{}

	for _, b := range brokers {
		wg.Add(1)
		go func(b *Broker, conf *Config) {
			defer wg.Done()
			_ = b.Open(conf) // Ensure that broker is opened

			response, err := b.ListGroups(&ListGroupsRequest{})
			if err != nil {
				errors <- err
				return
			}

			groups := make(map[string]string)
			for group, typ := range response.Groups {
				groups[group] = typ
			}

			groupMaps <- groups
		}(b, ca.conf)
	}

	wg.Wait()
	close(groupMaps)
	close(errors)

	for groupMap := range groupMaps {
		for group, protocolType := range groupMap {
			allGroups[group] = protocolType
		}
	}

	// Intentionally return only the first error for simplicity
	err = <-errors
	return
}

func (ca *clusterAdmin) ListConsumerGroupOffsets(group string, topicPartitions map[string][]int32) (*OffsetFetchResponse, error) {
	coordinator, err := ca.client.Coordinator(group)
	if err != nil {
		return nil, err
	}

	request := &OffsetFetchRequest{
		ConsumerGroup: group,
		partitions:    topicPartitions,
	}

	if ca.conf.Version.IsAtLeast(V0_10_2_0) {
		request.Version = 2
	} else if ca.conf.Version.IsAtLeast(V0_8_2_2) {
		request.Version = 1
	}

	return coordinator.FetchOffset(request)
}

func (ca *clusterAdmin) DeleteConsumerGroup(group string) error {
	coordinator, err := ca.client.Coordinator(group)
	if err != nil {
		return err
	}

	request := &DeleteGroupsRequest{
		Groups: []string{group},
	}

	resp, err := coordinator.DeleteGroups(request)
	if err != nil {
		return err
	}

	groupErr, ok := resp.GroupErrorCodes[group]
	if !ok {
		return ErrIncompleteResponse
	}

	if groupErr != ErrNoError {
		return groupErr
	}

	return nil
}

func (ca *clusterAdmin) DescribeLogDirs(brokerIds []int32) (allLogDirs map[int32][]DescribeLogDirsResponseDirMetadata, err error) {
	allLogDirs = make(map[int32][]DescribeLogDirsResponseDirMetadata)

	// Query brokers in parallel, since we may have to query multiple brokers
	logDirsMaps := make(chan map[int32][]DescribeLogDirsResponseDirMetadata, len(brokerIds))
	errors := make(chan error, len(brokerIds))
	wg := sync.WaitGroup{}

	for _, b := range brokerIds {
		wg.Add(1)
		broker, err := ca.findBroker(b)
		if err != nil {
			Logger.Printf("Unable to find broker with ID = %v\n", b)
			continue
		}
		go func(b *Broker, conf *Config) {
			defer wg.Done()
			_ = b.Open(conf) // Ensure that broker is opened

			response, err := b.DescribeLogDirs(&DescribeLogDirsRequest{})
			if err != nil {
				errors <- err
				return
			}
			logDirs := make(map[int32][]DescribeLogDirsResponseDirMetadata)
			logDirs[b.ID()] = response.LogDirs
			logDirsMaps <- logDirs
		}(broker, ca.conf)
	}

	wg.Wait()
	close(logDirsMaps)
	close(errors)

	for logDirsMap := range logDirsMaps {
		for id, logDirs := range logDirsMap {
			allLogDirs[id] = logDirs
		}
	}

	// Intentionally return only the first error for simplicity
	err = <-errors
	return
}
