package dockertest

import "github.com/ory-am/common/env"

// Dockertest configuration
var (
	// Debug if set, prevents any container from being removed.
	Debug bool

	// DockerMachineAvailable if true, uses docker-machine to run docker commands (for running tests on Windows and Mac OS)
	DockerMachineAvailable bool

	// DockerMachineName is the machine's name. You might want to use a dedicated machine for running your tests.
	// You can set this variable either directly or by defining a DOCKERTEST_IMAGE_NAME env variable.
	DockerMachineName = env.Getenv("DOCKERTEST_IMAGE_NAME", "default")

	// BindDockerToLocalhost if set, forces docker to bind the image to localhost. This for example is required when running tests on travis-ci.
	// You can set this variable either directly or by defining a DOCKERTEST_BIND_LOCALHOST env variable.
	// FIXME DOCKER_BIND_LOCALHOST remove legacy support
	BindDockerToLocalhost = env.Getenv("DOCKERTEST_BIND_LOCALHOST", env.Getenv("DOCKER_BIND_LOCALHOST", ""))

	// ContainerPrefix will be prepended to all containers started by dockertest to make identification of these "test images" hassle-free.
	ContainerPrefix = env.Getenv("DOCKERTEST_CONTAINER_PREFIX", "dockertest-")
)

// Image configuration
var (
	// MongoDBImageName is the MongoDB image name on dockerhub.
	MongoDBImageName = env.Getenv("DOCKERTEST_MONGODB_IMAGE_NAME", "mongo")

	// MySQLImageName is the MySQL image name on dockerhub.
	MySQLImageName = env.Getenv("DOCKERTEST_MYSQL_IMAGE_NAME", "mysql")

	// PostgresImageName is the PostgreSQL image name on dockerhub.
	PostgresImageName = env.Getenv("DOCKERTEST_POSTGRES_IMAGE_NAME", "postgres")

	// ElasticSearchImageName is the ElasticSearch image name on dockerhub.
	ElasticSearchImageName = env.Getenv("DOCKERTEST_ELASTICSEARCH_IMAGE_NAME", "elasticsearch")

	// RedisImageName is the Redis image name on dockerhub.
	RedisImageName = env.Getenv("DOCKERTEST_REDIS_IMAGE_NAME", "redis")

	// NSQImageName is the NSQ image name on dockerhub.
	NSQImageName = env.Getenv("DOCKERTEST_NSQ_IMAGE_NAME", "nsqio/nsq")

	// RethinkDBImageName is the RethinkDB image name on dockerhub.
	RethinkDBImageName = env.Getenv("DOCKERTEST_RETHINKDB_IMAGE_NAME", "rethinkdb")

	// RabbitMQImage name is the RabbitMQ image name on dockerhub.
	RabbitMQImageName = env.Getenv("DOCKERTEST_RABBITMQ_IMAGE_NAME", "rabbitmq")

	// ActiveMQImage name is the ActiveMQ image name on dockerhub.
	ActiveMQImageName = env.Getenv("DOCKERTEST_ACTIVEMQ_IMAGE_NAME", "webcenter/activemq")

	// MockserverImageName name is the Mockserver image name on dockerhub.
	MockserverImageName = env.Getenv("DOCKERTEST_MOCKSERVER_IMAGE_NAME", "jamesdbloom/mockserver")

	// ConsulImageName is the Consul image name on dockerhub.
	ConsulImageName = env.Getenv("DOCKERTEST_CONSUL_IMAGE_NAME", "consul")
)

// Username and password configuration
var (
	// MySQLUsername must be passed as username when connecting to mysql
	MySQLUsername = "root"

	// MySQLPassword must be passed as password when connecting to mysql
	MySQLPassword = "root"

	// PostgresUsername must be passed as username when connecting to postgres
	PostgresUsername = "postgres"

	// PostgresPassword must be passed as password when connecting to postgres
	PostgresPassword = "docker"
)
