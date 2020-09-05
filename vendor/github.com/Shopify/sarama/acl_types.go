package sarama

type (
	AclOperation int

	AclPermissionType int

	AclResourceType int

	AclResourcePatternType int
)

// ref: https://github.com/apache/kafka/blob/trunk/clients/src/main/java/org/apache/kafka/common/acl/AclOperation.java
const (
	AclOperationUnknown AclOperation = iota
	AclOperationAny
	AclOperationAll
	AclOperationRead
	AclOperationWrite
	AclOperationCreate
	AclOperationDelete
	AclOperationAlter
	AclOperationDescribe
	AclOperationClusterAction
	AclOperationDescribeConfigs
	AclOperationAlterConfigs
	AclOperationIdempotentWrite
)

// ref: https://github.com/apache/kafka/blob/trunk/clients/src/main/java/org/apache/kafka/common/acl/AclPermissionType.java
const (
	AclPermissionUnknown AclPermissionType = iota
	AclPermissionAny
	AclPermissionDeny
	AclPermissionAllow
)

// ref: https://github.com/apache/kafka/blob/trunk/clients/src/main/java/org/apache/kafka/common/resource/ResourceType.java
const (
	AclResourceUnknown AclResourceType = iota
	AclResourceAny
	AclResourceTopic
	AclResourceGroup
	AclResourceCluster
	AclResourceTransactionalID
)

// ref: https://github.com/apache/kafka/blob/trunk/clients/src/main/java/org/apache/kafka/common/resource/PatternType.java
const (
	AclPatternUnknown AclResourcePatternType = iota
	AclPatternAny
	AclPatternMatch
	AclPatternLiteral
	AclPatternPrefixed
)
