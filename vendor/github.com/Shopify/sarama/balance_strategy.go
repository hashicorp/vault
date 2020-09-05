package sarama

import (
	"container/heap"
	"math"
	"sort"
	"strings"
)

const (
	// RangeBalanceStrategyName identifies strategies that use the range partition assignment strategy
	RangeBalanceStrategyName = "range"

	// RoundRobinBalanceStrategyName identifies strategies that use the round-robin partition assignment strategy
	RoundRobinBalanceStrategyName = "roundrobin"

	// StickyBalanceStrategyName identifies strategies that use the sticky-partition assignment strategy
	StickyBalanceStrategyName = "sticky"

	defaultGeneration = -1
)

// BalanceStrategyPlan is the results of any BalanceStrategy.Plan attempt.
// It contains an allocation of topic/partitions by memberID in the form of
// a `memberID -> topic -> partitions` map.
type BalanceStrategyPlan map[string]map[string][]int32

// Add assigns a topic with a number partitions to a member.
func (p BalanceStrategyPlan) Add(memberID, topic string, partitions ...int32) {
	if len(partitions) == 0 {
		return
	}
	if _, ok := p[memberID]; !ok {
		p[memberID] = make(map[string][]int32, 1)
	}
	p[memberID][topic] = append(p[memberID][topic], partitions...)
}

// --------------------------------------------------------------------

// BalanceStrategy is used to balance topics and partitions
// across members of a consumer group
type BalanceStrategy interface {
	// Name uniquely identifies the strategy.
	Name() string

	// Plan accepts a map of `memberID -> metadata` and a map of `topic -> partitions`
	// and returns a distribution plan.
	Plan(members map[string]ConsumerGroupMemberMetadata, topics map[string][]int32) (BalanceStrategyPlan, error)

	// AssignmentData returns the serialized assignment data for the specified
	// memberID
	AssignmentData(memberID string, topics map[string][]int32, generationID int32) ([]byte, error)
}

// --------------------------------------------------------------------

// BalanceStrategyRange is the default and assigns partitions as ranges to consumer group members.
// Example with one topic T with six partitions (0..5) and two members (M1, M2):
//   M1: {T: [0, 1, 2]}
//   M2: {T: [3, 4, 5]}
var BalanceStrategyRange = &balanceStrategy{
	name: RangeBalanceStrategyName,
	coreFn: func(plan BalanceStrategyPlan, memberIDs []string, topic string, partitions []int32) {
		step := float64(len(partitions)) / float64(len(memberIDs))

		for i, memberID := range memberIDs {
			pos := float64(i)
			min := int(math.Floor(pos*step + 0.5))
			max := int(math.Floor((pos+1)*step + 0.5))
			plan.Add(memberID, topic, partitions[min:max]...)
		}
	},
}

// BalanceStrategyRoundRobin assigns partitions to members in alternating order.
// Example with topic T with six partitions (0..5) and two members (M1, M2):
//   M1: {T: [0, 2, 4]}
//   M2: {T: [1, 3, 5]}
var BalanceStrategyRoundRobin = &balanceStrategy{
	name: RoundRobinBalanceStrategyName,
	coreFn: func(plan BalanceStrategyPlan, memberIDs []string, topic string, partitions []int32) {
		for i, part := range partitions {
			memberID := memberIDs[i%len(memberIDs)]
			plan.Add(memberID, topic, part)
		}
	},
}

// BalanceStrategySticky assigns partitions to members with an attempt to preserve earlier assignments
// while maintain a balanced partition distribution.
// Example with topic T with six partitions (0..5) and two members (M1, M2):
//   M1: {T: [0, 2, 4]}
//   M2: {T: [1, 3, 5]}
//
// On reassignment with an additional consumer, you might get an assignment plan like:
//   M1: {T: [0, 2]}
//   M2: {T: [1, 3]}
//   M3: {T: [4, 5]}
//
var BalanceStrategySticky = &stickyBalanceStrategy{}

// --------------------------------------------------------------------

type balanceStrategy struct {
	name   string
	coreFn func(plan BalanceStrategyPlan, memberIDs []string, topic string, partitions []int32)
}

// Name implements BalanceStrategy.
func (s *balanceStrategy) Name() string { return s.name }

// Plan implements BalanceStrategy.
func (s *balanceStrategy) Plan(members map[string]ConsumerGroupMemberMetadata, topics map[string][]int32) (BalanceStrategyPlan, error) {
	// Build members by topic map
	mbt := make(map[string][]string)
	for memberID, meta := range members {
		for _, topic := range meta.Topics {
			mbt[topic] = append(mbt[topic], memberID)
		}
	}

	// Sort members for each topic
	for topic, memberIDs := range mbt {
		sort.Sort(&balanceStrategySortable{
			topic:     topic,
			memberIDs: memberIDs,
		})
	}

	// Assemble plan
	plan := make(BalanceStrategyPlan, len(members))
	for topic, memberIDs := range mbt {
		s.coreFn(plan, memberIDs, topic, topics[topic])
	}
	return plan, nil
}

// AssignmentData simple strategies do not require any shared assignment data
func (s *balanceStrategy) AssignmentData(memberID string, topics map[string][]int32, generationID int32) ([]byte, error) {
	return nil, nil
}

type balanceStrategySortable struct {
	topic     string
	memberIDs []string
}

func (p balanceStrategySortable) Len() int { return len(p.memberIDs) }
func (p balanceStrategySortable) Swap(i, j int) {
	p.memberIDs[i], p.memberIDs[j] = p.memberIDs[j], p.memberIDs[i]
}
func (p balanceStrategySortable) Less(i, j int) bool {
	return balanceStrategyHashValue(p.topic, p.memberIDs[i]) < balanceStrategyHashValue(p.topic, p.memberIDs[j])
}

func balanceStrategyHashValue(vv ...string) uint32 {
	h := uint32(2166136261)
	for _, s := range vv {
		for _, c := range s {
			h ^= uint32(c)
			h *= 16777619
		}
	}
	return h
}

type stickyBalanceStrategy struct {
	movements partitionMovements
}

// Name implements BalanceStrategy.
func (s *stickyBalanceStrategy) Name() string { return StickyBalanceStrategyName }

// Plan implements BalanceStrategy.
func (s *stickyBalanceStrategy) Plan(members map[string]ConsumerGroupMemberMetadata, topics map[string][]int32) (BalanceStrategyPlan, error) {
	// track partition movements during generation of the partition assignment plan
	s.movements = partitionMovements{
		Movements:                 make(map[topicPartitionAssignment]consumerPair),
		PartitionMovementsByTopic: make(map[string]map[consumerPair]map[topicPartitionAssignment]bool),
	}

	// prepopulate the current assignment state from userdata on the consumer group members
	currentAssignment, prevAssignment, err := prepopulateCurrentAssignments(members)
	if err != nil {
		return nil, err
	}

	// determine if we're dealing with a completely fresh assignment, or if there's existing assignment state
	isFreshAssignment := false
	if len(currentAssignment) == 0 {
		isFreshAssignment = true
	}

	// create a mapping of all current topic partitions and the consumers that can be assigned to them
	partition2AllPotentialConsumers := make(map[topicPartitionAssignment][]string)
	for topic, partitions := range topics {
		for _, partition := range partitions {
			partition2AllPotentialConsumers[topicPartitionAssignment{Topic: topic, Partition: partition}] = []string{}
		}
	}

	// create a mapping of all consumers to all potential topic partitions that can be assigned to them
	// also, populate the mapping of partitions to potential consumers
	consumer2AllPotentialPartitions := make(map[string][]topicPartitionAssignment, len(members))
	for memberID, meta := range members {
		consumer2AllPotentialPartitions[memberID] = make([]topicPartitionAssignment, 0)
		for _, topicSubscription := range meta.Topics {
			// only evaluate topic subscriptions that are present in the supplied topics map
			if _, found := topics[topicSubscription]; found {
				for _, partition := range topics[topicSubscription] {
					topicPartition := topicPartitionAssignment{Topic: topicSubscription, Partition: partition}
					consumer2AllPotentialPartitions[memberID] = append(consumer2AllPotentialPartitions[memberID], topicPartition)
					partition2AllPotentialConsumers[topicPartition] = append(partition2AllPotentialConsumers[topicPartition], memberID)
				}
			}
		}

		// add this consumer to currentAssignment (with an empty topic partition assignment) if it does not already exist
		if _, exists := currentAssignment[memberID]; !exists {
			currentAssignment[memberID] = make([]topicPartitionAssignment, 0)
		}
	}

	// create a mapping of each partition to its current consumer, where possible
	currentPartitionConsumers := make(map[topicPartitionAssignment]string, len(currentAssignment))
	unvisitedPartitions := make(map[topicPartitionAssignment]bool, len(partition2AllPotentialConsumers))
	for partition := range partition2AllPotentialConsumers {
		unvisitedPartitions[partition] = true
	}
	var unassignedPartitions []topicPartitionAssignment
	for memberID, partitions := range currentAssignment {
		var keepPartitions []topicPartitionAssignment
		for _, partition := range partitions {
			// If this partition no longer exists at all, likely due to the
			// topic being deleted, we remove the partition from the member.
			if _, exists := partition2AllPotentialConsumers[partition]; !exists {
				continue
			}
			delete(unvisitedPartitions, partition)
			currentPartitionConsumers[partition] = memberID

			if !strsContains(members[memberID].Topics, partition.Topic) {
				unassignedPartitions = append(unassignedPartitions, partition)
				continue
			}
			keepPartitions = append(keepPartitions, partition)
		}
		currentAssignment[memberID] = keepPartitions
	}
	for unvisited := range unvisitedPartitions {
		unassignedPartitions = append(unassignedPartitions, unvisited)
	}

	// sort the topic partitions in order of priority for reassignment
	sortedPartitions := sortPartitions(currentAssignment, prevAssignment, isFreshAssignment, partition2AllPotentialConsumers, consumer2AllPotentialPartitions)

	// at this point we have preserved all valid topic partition to consumer assignments and removed
	// all invalid topic partitions and invalid consumers. Now we need to assign unassignedPartitions
	// to consumers so that the topic partition assignments are as balanced as possible.

	// an ascending sorted set of consumers based on how many topic partitions are already assigned to them
	sortedCurrentSubscriptions := sortMemberIDsByPartitionAssignments(currentAssignment)
	s.balance(currentAssignment, prevAssignment, sortedPartitions, unassignedPartitions, sortedCurrentSubscriptions, consumer2AllPotentialPartitions, partition2AllPotentialConsumers, currentPartitionConsumers)

	// Assemble plan
	plan := make(BalanceStrategyPlan, len(currentAssignment))
	for memberID, assignments := range currentAssignment {
		if len(assignments) == 0 {
			plan[memberID] = make(map[string][]int32)
		} else {
			for _, assignment := range assignments {
				plan.Add(memberID, assignment.Topic, assignment.Partition)
			}
		}
	}
	return plan, nil
}

// AssignmentData serializes the set of topics currently assigned to the
// specified member as part of the supplied balance plan
func (s *stickyBalanceStrategy) AssignmentData(memberID string, topics map[string][]int32, generationID int32) ([]byte, error) {
	return encode(&StickyAssignorUserDataV1{
		Topics:     topics,
		Generation: generationID,
	}, nil)
}

func strsContains(s []string, value string) bool {
	for _, entry := range s {
		if entry == value {
			return true
		}
	}
	return false
}

// Balance assignments across consumers for maximum fairness and stickiness.
func (s *stickyBalanceStrategy) balance(currentAssignment map[string][]topicPartitionAssignment, prevAssignment map[topicPartitionAssignment]consumerGenerationPair, sortedPartitions []topicPartitionAssignment, unassignedPartitions []topicPartitionAssignment, sortedCurrentSubscriptions []string, consumer2AllPotentialPartitions map[string][]topicPartitionAssignment, partition2AllPotentialConsumers map[topicPartitionAssignment][]string, currentPartitionConsumer map[topicPartitionAssignment]string) {
	initializing := false
	if len(sortedCurrentSubscriptions) == 0 || len(currentAssignment[sortedCurrentSubscriptions[0]]) == 0 {
		initializing = true
	}

	// assign all unassigned partitions
	for _, partition := range unassignedPartitions {
		// skip if there is no potential consumer for the partition
		if len(partition2AllPotentialConsumers[partition]) == 0 {
			continue
		}
		sortedCurrentSubscriptions = assignPartition(partition, sortedCurrentSubscriptions, currentAssignment, consumer2AllPotentialPartitions, currentPartitionConsumer)
	}

	// narrow down the reassignment scope to only those partitions that can actually be reassigned
	for partition := range partition2AllPotentialConsumers {
		if !canTopicPartitionParticipateInReassignment(partition, partition2AllPotentialConsumers) {
			sortedPartitions = removeTopicPartitionFromMemberAssignments(sortedPartitions, partition)
		}
	}

	// narrow down the reassignment scope to only those consumers that are subject to reassignment
	fixedAssignments := make(map[string][]topicPartitionAssignment)
	for memberID := range consumer2AllPotentialPartitions {
		if !canConsumerParticipateInReassignment(memberID, currentAssignment, consumer2AllPotentialPartitions, partition2AllPotentialConsumers) {
			fixedAssignments[memberID] = currentAssignment[memberID]
			delete(currentAssignment, memberID)
			sortedCurrentSubscriptions = sortMemberIDsByPartitionAssignments(currentAssignment)
		}
	}

	// create a deep copy of the current assignment so we can revert to it if we do not get a more balanced assignment later
	preBalanceAssignment := deepCopyAssignment(currentAssignment)
	preBalancePartitionConsumers := make(map[topicPartitionAssignment]string, len(currentPartitionConsumer))
	for k, v := range currentPartitionConsumer {
		preBalancePartitionConsumers[k] = v
	}

	reassignmentPerformed := s.performReassignments(sortedPartitions, currentAssignment, prevAssignment, sortedCurrentSubscriptions, consumer2AllPotentialPartitions, partition2AllPotentialConsumers, currentPartitionConsumer)

	// if we are not preserving existing assignments and we have made changes to the current assignment
	// make sure we are getting a more balanced assignment; otherwise, revert to previous assignment
	if !initializing && reassignmentPerformed && getBalanceScore(currentAssignment) >= getBalanceScore(preBalanceAssignment) {
		currentAssignment = deepCopyAssignment(preBalanceAssignment)
		currentPartitionConsumer = make(map[topicPartitionAssignment]string, len(preBalancePartitionConsumers))
		for k, v := range preBalancePartitionConsumers {
			currentPartitionConsumer[k] = v
		}
	}

	// add the fixed assignments (those that could not change) back
	for consumer, assignments := range fixedAssignments {
		currentAssignment[consumer] = assignments
	}
}

// Calculate the balance score of the given assignment, as the sum of assigned partitions size difference of all consumer pairs.
// A perfectly balanced assignment (with all consumers getting the same number of partitions) has a balance score of 0.
// Lower balance score indicates a more balanced assignment.
func getBalanceScore(assignment map[string][]topicPartitionAssignment) int {
	consumer2AssignmentSize := make(map[string]int, len(assignment))
	for memberID, partitions := range assignment {
		consumer2AssignmentSize[memberID] = len(partitions)
	}

	var score float64
	for memberID, consumerAssignmentSize := range consumer2AssignmentSize {
		delete(consumer2AssignmentSize, memberID)
		for _, otherConsumerAssignmentSize := range consumer2AssignmentSize {
			score += math.Abs(float64(consumerAssignmentSize - otherConsumerAssignmentSize))
		}
	}
	return int(score)
}

// Determine whether the current assignment plan is balanced.
func isBalanced(currentAssignment map[string][]topicPartitionAssignment, sortedCurrentSubscriptions []string, allSubscriptions map[string][]topicPartitionAssignment) bool {
	sortedCurrentSubscriptions = sortMemberIDsByPartitionAssignments(currentAssignment)
	min := len(currentAssignment[sortedCurrentSubscriptions[0]])
	max := len(currentAssignment[sortedCurrentSubscriptions[len(sortedCurrentSubscriptions)-1]])
	if min >= max-1 {
		// if minimum and maximum numbers of partitions assigned to consumers differ by at most one return true
		return true
	}

	// create a mapping from partitions to the consumer assigned to them
	allPartitions := make(map[topicPartitionAssignment]string)
	for memberID, partitions := range currentAssignment {
		for _, partition := range partitions {
			if _, exists := allPartitions[partition]; exists {
				Logger.Printf("Topic %s Partition %d is assigned more than one consumer", partition.Topic, partition.Partition)
			}
			allPartitions[partition] = memberID
		}
	}

	// for each consumer that does not have all the topic partitions it can get make sure none of the topic partitions it
	// could but did not get cannot be moved to it (because that would break the balance)
	for _, memberID := range sortedCurrentSubscriptions {
		consumerPartitions := currentAssignment[memberID]
		consumerPartitionCount := len(consumerPartitions)

		// skip if this consumer already has all the topic partitions it can get
		if consumerPartitionCount == len(allSubscriptions[memberID]) {
			continue
		}

		// otherwise make sure it cannot get any more
		potentialTopicPartitions := allSubscriptions[memberID]
		for _, partition := range potentialTopicPartitions {
			if !memberAssignmentsIncludeTopicPartition(currentAssignment[memberID], partition) {
				otherConsumer := allPartitions[partition]
				otherConsumerPartitionCount := len(currentAssignment[otherConsumer])
				if consumerPartitionCount < otherConsumerPartitionCount {
					return false
				}
			}
		}
	}
	return true
}

// Reassign all topic partitions that need reassignment until balanced.
func (s *stickyBalanceStrategy) performReassignments(reassignablePartitions []topicPartitionAssignment, currentAssignment map[string][]topicPartitionAssignment, prevAssignment map[topicPartitionAssignment]consumerGenerationPair, sortedCurrentSubscriptions []string, consumer2AllPotentialPartitions map[string][]topicPartitionAssignment, partition2AllPotentialConsumers map[topicPartitionAssignment][]string, currentPartitionConsumer map[topicPartitionAssignment]string) bool {
	reassignmentPerformed := false
	modified := false

	// repeat reassignment until no partition can be moved to improve the balance
	for {
		modified = false
		// reassign all reassignable partitions (starting from the partition with least potential consumers and if needed)
		// until the full list is processed or a balance is achieved
		for _, partition := range reassignablePartitions {
			if isBalanced(currentAssignment, sortedCurrentSubscriptions, consumer2AllPotentialPartitions) {
				break
			}

			// the partition must have at least two consumers
			if len(partition2AllPotentialConsumers[partition]) <= 1 {
				Logger.Printf("Expected more than one potential consumer for partition %s topic %d", partition.Topic, partition.Partition)
			}

			// the partition must have a consumer
			consumer := currentPartitionConsumer[partition]
			if consumer == "" {
				Logger.Printf("Expected topic %s partition %d to be assigned to a consumer", partition.Topic, partition.Partition)
			}

			if _, exists := prevAssignment[partition]; exists {
				if len(currentAssignment[consumer]) > (len(currentAssignment[prevAssignment[partition].MemberID]) + 1) {
					sortedCurrentSubscriptions = s.reassignPartition(partition, currentAssignment, sortedCurrentSubscriptions, currentPartitionConsumer, prevAssignment[partition].MemberID)
					reassignmentPerformed = true
					modified = true
					continue
				}
			}

			// check if a better-suited consumer exists for the partition; if so, reassign it
			for _, otherConsumer := range partition2AllPotentialConsumers[partition] {
				if len(currentAssignment[consumer]) > (len(currentAssignment[otherConsumer]) + 1) {
					sortedCurrentSubscriptions = s.reassignPartitionToNewConsumer(partition, currentAssignment, sortedCurrentSubscriptions, currentPartitionConsumer, consumer2AllPotentialPartitions)
					reassignmentPerformed = true
					modified = true
					break
				}
			}
		}
		if !modified {
			return reassignmentPerformed
		}
	}
}

// Identify a new consumer for a topic partition and reassign it.
func (s *stickyBalanceStrategy) reassignPartitionToNewConsumer(partition topicPartitionAssignment, currentAssignment map[string][]topicPartitionAssignment, sortedCurrentSubscriptions []string, currentPartitionConsumer map[topicPartitionAssignment]string, consumer2AllPotentialPartitions map[string][]topicPartitionAssignment) []string {
	for _, anotherConsumer := range sortedCurrentSubscriptions {
		if memberAssignmentsIncludeTopicPartition(consumer2AllPotentialPartitions[anotherConsumer], partition) {
			return s.reassignPartition(partition, currentAssignment, sortedCurrentSubscriptions, currentPartitionConsumer, anotherConsumer)
		}
	}
	return sortedCurrentSubscriptions
}

// Reassign a specific partition to a new consumer
func (s *stickyBalanceStrategy) reassignPartition(partition topicPartitionAssignment, currentAssignment map[string][]topicPartitionAssignment, sortedCurrentSubscriptions []string, currentPartitionConsumer map[topicPartitionAssignment]string, newConsumer string) []string {
	consumer := currentPartitionConsumer[partition]
	// find the correct partition movement considering the stickiness requirement
	partitionToBeMoved := s.movements.getTheActualPartitionToBeMoved(partition, consumer, newConsumer)
	return s.processPartitionMovement(partitionToBeMoved, newConsumer, currentAssignment, sortedCurrentSubscriptions, currentPartitionConsumer)
}

// Track the movement of a topic partition after assignment
func (s *stickyBalanceStrategy) processPartitionMovement(partition topicPartitionAssignment, newConsumer string, currentAssignment map[string][]topicPartitionAssignment, sortedCurrentSubscriptions []string, currentPartitionConsumer map[topicPartitionAssignment]string) []string {
	oldConsumer := currentPartitionConsumer[partition]
	s.movements.movePartition(partition, oldConsumer, newConsumer)

	currentAssignment[oldConsumer] = removeTopicPartitionFromMemberAssignments(currentAssignment[oldConsumer], partition)
	currentAssignment[newConsumer] = append(currentAssignment[newConsumer], partition)
	currentPartitionConsumer[partition] = newConsumer
	return sortMemberIDsByPartitionAssignments(currentAssignment)
}

// Determine whether a specific consumer should be considered for topic partition assignment.
func canConsumerParticipateInReassignment(memberID string, currentAssignment map[string][]topicPartitionAssignment, consumer2AllPotentialPartitions map[string][]topicPartitionAssignment, partition2AllPotentialConsumers map[topicPartitionAssignment][]string) bool {
	currentPartitions := currentAssignment[memberID]
	currentAssignmentSize := len(currentPartitions)
	maxAssignmentSize := len(consumer2AllPotentialPartitions[memberID])
	if currentAssignmentSize > maxAssignmentSize {
		Logger.Printf("The consumer %s is assigned more partitions than the maximum possible", memberID)
	}
	if currentAssignmentSize < maxAssignmentSize {
		// if a consumer is not assigned all its potential partitions it is subject to reassignment
		return true
	}
	for _, partition := range currentPartitions {
		if canTopicPartitionParticipateInReassignment(partition, partition2AllPotentialConsumers) {
			return true
		}
	}
	return false
}

// Only consider reassigning those topic partitions that have two or more potential consumers.
func canTopicPartitionParticipateInReassignment(partition topicPartitionAssignment, partition2AllPotentialConsumers map[topicPartitionAssignment][]string) bool {
	return len(partition2AllPotentialConsumers[partition]) >= 2
}

// The assignment should improve the overall balance of the partition assignments to consumers.
func assignPartition(partition topicPartitionAssignment, sortedCurrentSubscriptions []string, currentAssignment map[string][]topicPartitionAssignment, consumer2AllPotentialPartitions map[string][]topicPartitionAssignment, currentPartitionConsumer map[topicPartitionAssignment]string) []string {
	for _, memberID := range sortedCurrentSubscriptions {
		if memberAssignmentsIncludeTopicPartition(consumer2AllPotentialPartitions[memberID], partition) {
			currentAssignment[memberID] = append(currentAssignment[memberID], partition)
			currentPartitionConsumer[partition] = memberID
			break
		}
	}
	return sortMemberIDsByPartitionAssignments(currentAssignment)
}

// Deserialize topic partition assignment data to aid with creation of a sticky assignment.
func deserializeTopicPartitionAssignment(userDataBytes []byte) (StickyAssignorUserData, error) {
	userDataV1 := &StickyAssignorUserDataV1{}
	if err := decode(userDataBytes, userDataV1); err != nil {
		userDataV0 := &StickyAssignorUserDataV0{}
		if err := decode(userDataBytes, userDataV0); err != nil {
			return nil, err
		}
		return userDataV0, nil
	}
	return userDataV1, nil
}

// filterAssignedPartitions returns a map of consumer group members to their list of previously-assigned topic partitions, limited
// to those topic partitions currently reported by the Kafka cluster.
func filterAssignedPartitions(currentAssignment map[string][]topicPartitionAssignment, partition2AllPotentialConsumers map[topicPartitionAssignment][]string) map[string][]topicPartitionAssignment {
	assignments := deepCopyAssignment(currentAssignment)
	for memberID, partitions := range assignments {
		// perform in-place filtering
		i := 0
		for _, partition := range partitions {
			if _, exists := partition2AllPotentialConsumers[partition]; exists {
				partitions[i] = partition
				i++
			}
		}
		assignments[memberID] = partitions[:i]
	}
	return assignments
}

func removeTopicPartitionFromMemberAssignments(assignments []topicPartitionAssignment, topic topicPartitionAssignment) []topicPartitionAssignment {
	for i, assignment := range assignments {
		if assignment == topic {
			return append(assignments[:i], assignments[i+1:]...)
		}
	}
	return assignments
}

func memberAssignmentsIncludeTopicPartition(assignments []topicPartitionAssignment, topic topicPartitionAssignment) bool {
	for _, assignment := range assignments {
		if assignment == topic {
			return true
		}
	}
	return false
}

func sortPartitions(currentAssignment map[string][]topicPartitionAssignment, partitionsWithADifferentPreviousAssignment map[topicPartitionAssignment]consumerGenerationPair, isFreshAssignment bool, partition2AllPotentialConsumers map[topicPartitionAssignment][]string, consumer2AllPotentialPartitions map[string][]topicPartitionAssignment) []topicPartitionAssignment {
	unassignedPartitions := make(map[topicPartitionAssignment]bool, len(partition2AllPotentialConsumers))
	for partition := range partition2AllPotentialConsumers {
		unassignedPartitions[partition] = true
	}

	sortedPartitions := make([]topicPartitionAssignment, 0)
	if !isFreshAssignment && areSubscriptionsIdentical(partition2AllPotentialConsumers, consumer2AllPotentialPartitions) {
		// if this is a reassignment and the subscriptions are identical (all consumers can consumer from all topics)
		// then we just need to simply list partitions in a round robin fashion (from consumers with
		// most assigned partitions to those with least)
		assignments := filterAssignedPartitions(currentAssignment, partition2AllPotentialConsumers)

		// use priority-queue to evaluate consumer group members in descending-order based on
		// the number of topic partition assignments (i.e. consumers with most assignments first)
		pq := make(assignmentPriorityQueue, len(assignments))
		i := 0
		for consumerID, consumerAssignments := range assignments {
			pq[i] = &consumerGroupMember{
				id:          consumerID,
				assignments: consumerAssignments,
			}
			i++
		}
		heap.Init(&pq)

		for {
			// loop until no consumer-group members remain
			if pq.Len() == 0 {
				break
			}
			member := pq[0]

			// partitions that were assigned to a different consumer last time
			var prevPartitionIndex int
			for i, partition := range member.assignments {
				if _, exists := partitionsWithADifferentPreviousAssignment[partition]; exists {
					prevPartitionIndex = i
					break
				}
			}

			if len(member.assignments) > 0 {
				partition := member.assignments[prevPartitionIndex]
				sortedPartitions = append(sortedPartitions, partition)
				delete(unassignedPartitions, partition)
				if prevPartitionIndex == 0 {
					member.assignments = member.assignments[1:]
				} else {
					member.assignments = append(member.assignments[:prevPartitionIndex], member.assignments[prevPartitionIndex+1:]...)
				}
				heap.Fix(&pq, 0)
			} else {
				heap.Pop(&pq)
			}
		}

		for partition := range unassignedPartitions {
			sortedPartitions = append(sortedPartitions, partition)
		}
	} else {
		// an ascending sorted set of topic partitions based on how many consumers can potentially use them
		sortedPartitions = sortPartitionsByPotentialConsumerAssignments(partition2AllPotentialConsumers)
	}
	return sortedPartitions
}

func sortMemberIDsByPartitionAssignments(assignments map[string][]topicPartitionAssignment) []string {
	// sort the members by the number of partition assignments in ascending order
	sortedMemberIDs := make([]string, 0, len(assignments))
	for memberID := range assignments {
		sortedMemberIDs = append(sortedMemberIDs, memberID)
	}
	sort.SliceStable(sortedMemberIDs, func(i, j int) bool {
		ret := len(assignments[sortedMemberIDs[i]]) - len(assignments[sortedMemberIDs[j]])
		if ret == 0 {
			return sortedMemberIDs[i] < sortedMemberIDs[j]
		}
		return len(assignments[sortedMemberIDs[i]]) < len(assignments[sortedMemberIDs[j]])
	})
	return sortedMemberIDs
}

func sortPartitionsByPotentialConsumerAssignments(partition2AllPotentialConsumers map[topicPartitionAssignment][]string) []topicPartitionAssignment {
	// sort the members by the number of partition assignments in descending order
	sortedPartionIDs := make([]topicPartitionAssignment, len(partition2AllPotentialConsumers))
	i := 0
	for partition := range partition2AllPotentialConsumers {
		sortedPartionIDs[i] = partition
		i++
	}
	sort.Slice(sortedPartionIDs, func(i, j int) bool {
		if len(partition2AllPotentialConsumers[sortedPartionIDs[i]]) == len(partition2AllPotentialConsumers[sortedPartionIDs[j]]) {
			ret := strings.Compare(sortedPartionIDs[i].Topic, sortedPartionIDs[j].Topic)
			if ret == 0 {
				return sortedPartionIDs[i].Partition < sortedPartionIDs[j].Partition
			}
			return ret < 0
		}
		return len(partition2AllPotentialConsumers[sortedPartionIDs[i]]) < len(partition2AllPotentialConsumers[sortedPartionIDs[j]])
	})
	return sortedPartionIDs
}

func deepCopyAssignment(assignment map[string][]topicPartitionAssignment) map[string][]topicPartitionAssignment {
	copy := make(map[string][]topicPartitionAssignment, len(assignment))
	for memberID, subscriptions := range assignment {
		copy[memberID] = append(subscriptions[:0:0], subscriptions...)
	}
	return copy
}

func areSubscriptionsIdentical(partition2AllPotentialConsumers map[topicPartitionAssignment][]string, consumer2AllPotentialPartitions map[string][]topicPartitionAssignment) bool {
	curMembers := make(map[string]int)
	for _, cur := range partition2AllPotentialConsumers {
		if len(curMembers) == 0 {
			for _, curMembersElem := range cur {
				curMembers[curMembersElem]++
			}
			continue
		}

		if len(curMembers) != len(cur) {
			return false
		}

		yMap := make(map[string]int)
		for _, yElem := range cur {
			yMap[yElem]++
		}

		for curMembersMapKey, curMembersMapVal := range curMembers {
			if yMap[curMembersMapKey] != curMembersMapVal {
				return false
			}
		}
	}

	curPartitions := make(map[topicPartitionAssignment]int)
	for _, cur := range consumer2AllPotentialPartitions {
		if len(curPartitions) == 0 {
			for _, curPartitionElem := range cur {
				curPartitions[curPartitionElem]++
			}
			continue
		}

		if len(curPartitions) != len(cur) {
			return false
		}

		yMap := make(map[topicPartitionAssignment]int)
		for _, yElem := range cur {
			yMap[yElem]++
		}

		for curMembersMapKey, curMembersMapVal := range curPartitions {
			if yMap[curMembersMapKey] != curMembersMapVal {
				return false
			}
		}
	}
	return true
}

// We need to process subscriptions' user data with each consumer's reported generation in mind
// higher generations overwrite lower generations in case of a conflict
// note that a conflict could exist only if user data is for different generations
func prepopulateCurrentAssignments(members map[string]ConsumerGroupMemberMetadata) (map[string][]topicPartitionAssignment, map[topicPartitionAssignment]consumerGenerationPair, error) {
	currentAssignment := make(map[string][]topicPartitionAssignment)
	prevAssignment := make(map[topicPartitionAssignment]consumerGenerationPair)

	// for each partition we create a sorted map of its consumers by generation
	sortedPartitionConsumersByGeneration := make(map[topicPartitionAssignment]map[int]string)
	for memberID, meta := range members {
		consumerUserData, err := deserializeTopicPartitionAssignment(meta.UserData)
		if err != nil {
			return nil, nil, err
		}
		for _, partition := range consumerUserData.partitions() {
			if consumers, exists := sortedPartitionConsumersByGeneration[partition]; exists {
				if consumerUserData.hasGeneration() {
					if _, generationExists := consumers[consumerUserData.generation()]; generationExists {
						// same partition is assigned to two consumers during the same rebalance.
						// log a warning and skip this record
						Logger.Printf("Topic %s Partition %d is assigned to multiple consumers following sticky assignment generation %d", partition.Topic, partition.Partition, consumerUserData.generation())
						continue
					} else {
						consumers[consumerUserData.generation()] = memberID
					}
				} else {
					consumers[defaultGeneration] = memberID
				}
			} else {
				generation := defaultGeneration
				if consumerUserData.hasGeneration() {
					generation = consumerUserData.generation()
				}
				sortedPartitionConsumersByGeneration[partition] = map[int]string{generation: memberID}
			}
		}
	}

	// prevAssignment holds the prior ConsumerGenerationPair (before current) of each partition
	// current and previous consumers are the last two consumers of each partition in the above sorted map
	for partition, consumers := range sortedPartitionConsumersByGeneration {
		// sort consumers by generation in decreasing order
		var generations []int
		for generation := range consumers {
			generations = append(generations, generation)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(generations)))

		consumer := consumers[generations[0]]
		if _, exists := currentAssignment[consumer]; !exists {
			currentAssignment[consumer] = []topicPartitionAssignment{partition}
		} else {
			currentAssignment[consumer] = append(currentAssignment[consumer], partition)
		}

		// check for previous assignment, if any
		if len(generations) > 1 {
			prevAssignment[partition] = consumerGenerationPair{
				MemberID:   consumers[generations[1]],
				Generation: generations[1],
			}
		}
	}
	return currentAssignment, prevAssignment, nil
}

type consumerGenerationPair struct {
	MemberID   string
	Generation int
}

// consumerPair represents a pair of Kafka consumer ids involved in a partition reassignment.
type consumerPair struct {
	SrcMemberID string
	DstMemberID string
}

// partitionMovements maintains some data structures to simplify lookup of partition movements among consumers.
type partitionMovements struct {
	PartitionMovementsByTopic map[string]map[consumerPair]map[topicPartitionAssignment]bool
	Movements                 map[topicPartitionAssignment]consumerPair
}

func (p *partitionMovements) removeMovementRecordOfPartition(partition topicPartitionAssignment) consumerPair {
	pair := p.Movements[partition]
	delete(p.Movements, partition)

	partitionMovementsForThisTopic := p.PartitionMovementsByTopic[partition.Topic]
	delete(partitionMovementsForThisTopic[pair], partition)
	if len(partitionMovementsForThisTopic[pair]) == 0 {
		delete(partitionMovementsForThisTopic, pair)
	}
	if len(p.PartitionMovementsByTopic[partition.Topic]) == 0 {
		delete(p.PartitionMovementsByTopic, partition.Topic)
	}
	return pair
}

func (p *partitionMovements) addPartitionMovementRecord(partition topicPartitionAssignment, pair consumerPair) {
	p.Movements[partition] = pair
	if _, exists := p.PartitionMovementsByTopic[partition.Topic]; !exists {
		p.PartitionMovementsByTopic[partition.Topic] = make(map[consumerPair]map[topicPartitionAssignment]bool)
	}
	partitionMovementsForThisTopic := p.PartitionMovementsByTopic[partition.Topic]
	if _, exists := partitionMovementsForThisTopic[pair]; !exists {
		partitionMovementsForThisTopic[pair] = make(map[topicPartitionAssignment]bool)
	}
	partitionMovementsForThisTopic[pair][partition] = true
}

func (p *partitionMovements) movePartition(partition topicPartitionAssignment, oldConsumer, newConsumer string) {
	pair := consumerPair{
		SrcMemberID: oldConsumer,
		DstMemberID: newConsumer,
	}
	if _, exists := p.Movements[partition]; exists {
		// this partition has previously moved
		existingPair := p.removeMovementRecordOfPartition(partition)
		if existingPair.DstMemberID != oldConsumer {
			Logger.Printf("Existing pair DstMemberID %s was not equal to the oldConsumer ID %s", existingPair.DstMemberID, oldConsumer)
		}
		if existingPair.SrcMemberID != newConsumer {
			// the partition is not moving back to its previous consumer
			p.addPartitionMovementRecord(partition, consumerPair{
				SrcMemberID: existingPair.SrcMemberID,
				DstMemberID: newConsumer,
			})
		}
	} else {
		p.addPartitionMovementRecord(partition, pair)
	}
}

func (p *partitionMovements) getTheActualPartitionToBeMoved(partition topicPartitionAssignment, oldConsumer, newConsumer string) topicPartitionAssignment {
	if _, exists := p.PartitionMovementsByTopic[partition.Topic]; !exists {
		return partition
	}
	if _, exists := p.Movements[partition]; exists {
		// this partition has previously moved
		if oldConsumer != p.Movements[partition].DstMemberID {
			Logger.Printf("Partition movement DstMemberID %s was not equal to the oldConsumer ID %s", p.Movements[partition].DstMemberID, oldConsumer)
		}
		oldConsumer = p.Movements[partition].SrcMemberID
	}

	partitionMovementsForThisTopic := p.PartitionMovementsByTopic[partition.Topic]
	reversePair := consumerPair{
		SrcMemberID: newConsumer,
		DstMemberID: oldConsumer,
	}
	if _, exists := partitionMovementsForThisTopic[reversePair]; !exists {
		return partition
	}
	var reversePairPartition topicPartitionAssignment
	for otherPartition := range partitionMovementsForThisTopic[reversePair] {
		reversePairPartition = otherPartition
	}
	return reversePairPartition
}

func (p *partitionMovements) isLinked(src, dst string, pairs []consumerPair, currentPath []string) ([]string, bool) {
	if src == dst {
		return currentPath, false
	}
	if len(pairs) == 0 {
		return currentPath, false
	}
	for _, pair := range pairs {
		if src == pair.SrcMemberID && dst == pair.DstMemberID {
			currentPath = append(currentPath, src, dst)
			return currentPath, true
		}
	}

	for _, pair := range pairs {
		if pair.SrcMemberID == src {
			// create a deep copy of the pairs, excluding the current pair
			reducedSet := make([]consumerPair, len(pairs)-1)
			i := 0
			for _, p := range pairs {
				if p != pair {
					reducedSet[i] = pair
					i++
				}
			}

			currentPath = append(currentPath, pair.SrcMemberID)
			return p.isLinked(pair.DstMemberID, dst, reducedSet, currentPath)
		}
	}
	return currentPath, false
}

func (p *partitionMovements) in(cycle []string, cycles [][]string) bool {
	superCycle := make([]string, len(cycle)-1)
	for i := 0; i < len(cycle)-1; i++ {
		superCycle[i] = cycle[i]
	}
	superCycle = append(superCycle, cycle...)
	for _, foundCycle := range cycles {
		if len(foundCycle) == len(cycle) && indexOfSubList(superCycle, foundCycle) != -1 {
			return true
		}
	}
	return false
}

func (p *partitionMovements) hasCycles(pairs []consumerPair) bool {
	cycles := make([][]string, 0)
	for _, pair := range pairs {
		// create a deep copy of the pairs, excluding the current pair
		reducedPairs := make([]consumerPair, len(pairs)-1)
		i := 0
		for _, p := range pairs {
			if p != pair {
				reducedPairs[i] = pair
				i++
			}
		}
		if path, linked := p.isLinked(pair.DstMemberID, pair.SrcMemberID, reducedPairs, []string{pair.SrcMemberID}); linked {
			if !p.in(path, cycles) {
				cycles = append(cycles, path)
				Logger.Printf("A cycle of length %d was found: %v", len(path)-1, path)
			}
		}
	}

	// for now we want to make sure there is no partition movements of the same topic between a pair of consumers.
	// the odds of finding a cycle among more than two consumers seem to be very low (according to various randomized
	// tests with the given sticky algorithm) that it should not worth the added complexity of handling those cases.
	for _, cycle := range cycles {
		if len(cycle) == 3 {
			return true
		}
	}
	return false
}

func (p *partitionMovements) isSticky() bool {
	for topic, movements := range p.PartitionMovementsByTopic {
		movementPairs := make([]consumerPair, len(movements))
		i := 0
		for pair := range movements {
			movementPairs[i] = pair
			i++
		}
		if p.hasCycles(movementPairs) {
			Logger.Printf("Stickiness is violated for topic %s", topic)
			Logger.Printf("Partition movements for this topic occurred among the following consumer pairs: %v", movements)
			return false
		}
	}
	return true
}

func indexOfSubList(source []string, target []string) int {
	targetSize := len(target)
	maxCandidate := len(source) - targetSize
nextCand:
	for candidate := 0; candidate <= maxCandidate; candidate++ {
		j := candidate
		for i := 0; i < targetSize; i++ {
			if target[i] != source[j] {
				// Element mismatch, try next cand
				continue nextCand
			}
			j++
		}
		// All elements of candidate matched target
		return candidate
	}
	return -1
}

type consumerGroupMember struct {
	id          string
	assignments []topicPartitionAssignment
}

// assignmentPriorityQueue is a priority-queue of consumer group members that is sorted
// in descending order (most assignments to least assignments).
type assignmentPriorityQueue []*consumerGroupMember

func (pq assignmentPriorityQueue) Len() int { return len(pq) }

func (pq assignmentPriorityQueue) Less(i, j int) bool {
	// order asssignment priority queue in descending order using assignment-count/member-id
	if len(pq[i].assignments) == len(pq[j].assignments) {
		return strings.Compare(pq[i].id, pq[j].id) > 0
	}
	return len(pq[i].assignments) > len(pq[j].assignments)
}

func (pq assignmentPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *assignmentPriorityQueue) Push(x interface{}) {
	member := x.(*consumerGroupMember)
	*pq = append(*pq, member)
}

func (pq *assignmentPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	member := old[n-1]
	*pq = old[0 : n-1]
	return member
}
