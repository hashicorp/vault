package sarama

// ConfigResourceType is a type for resources that have configs.
type ConfigResourceType int8

// Taken from:
// https://github.com/apache/kafka/blob/ed7c071e07f1f90e4c2895582f61ca090ced3c42/clients/src/main/java/org/apache/kafka/common/config/ConfigResource.java#L32-L55

const (
	// UnknownResource constant type
	UnknownResource ConfigResourceType = 0
	// TopicResource constant type
	TopicResource ConfigResourceType = 2
	// BrokerResource constant type
	BrokerResource ConfigResourceType = 4
	// BrokerLoggerResource constant type
	BrokerLoggerResource ConfigResourceType = 8
)
