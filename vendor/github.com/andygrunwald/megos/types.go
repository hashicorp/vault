package megos

// State represents the JSON from the state.json of a mesos node
type State struct {
	Version                string      `json:"version"`
	GitSHA                 string      `json:"git_sha"`
	GitTag                 string      `json:"git_tag"`
	BuildDate              string      `json:"build_date"`
	BuildTime              float64     `json:"build_time"`
	BuildUser              string      `json:"build_user"`
	StartTime              float64     `json:"start_time"`
	ElectedTime            float64     `json:"elected_time"`
	ID                     string      `json:"id"`
	PID                    string      `json:"pid"`
	Hostname               string      `json:"hostname"`
	ActivatedSlaves        float64     `json:"activated_slaves"`
	DeactivatedSlaves      float64     `json:"deactivated_slaves"`
	Cluster                string      `json:"cluster"`
	Leader                 string      `json:"leader"`
	CompletedFrameworks    []Framework `json:"completed_frameworks"`
	OrphanTasks            []Task      `json:"orphan_tasks"`
	UnregisteredFrameworks []string    `json:"unregistered_frameworks"`
	Flags                  Flags       `json:"flags"`
	Slaves                 []Slave     `json:"slaves"`
	Frameworks             []Framework `json:"frameworks"`
	GitBranch              string      `json:"git_branch"`
	LogDir                 string      `json:"log_dir"`
	ExternalLogFile        string      `json:"external_log_file"`
}

// Flags represents the flags of a mesos state
type Flags struct {
	AppcStoreDir                     string `json:"appc_store_dir"`
	AllocationInterval               string `json:"allocation_interval"`
	Allocator                        string `json:"allocator"`
	Authenticate                     string `json:"authenticate"`
	AuthenticateHTTP                 string `json:"authenticate_http"`
	Authenticatee                    string `json:"authenticatee"`
	AuthenticateSlaves               string `json:"authenticate_slaves"`
	Authenticators                   string `json:"authenticators"`
	Authorizers                      string `json:"authorizers"`
	CgroupsCPUEnablePIDsAndTIDsCount string `json:"cgroups_cpu_enable_pids_and_tids_count"`
	CgroupsEnableCfs                 string `json:"cgroups_enable_cfs"`
	CgroupsHierarchy                 string `json:"cgroups_hierarchy"`
	CgroupsLimitSwap                 string `json:"cgroups_limit_swap"`
	CgroupsRoot                      string `json:"cgroups_root"`
	Cluster                          string `json:"cluster"`
	ContainerDiskWatchInterval       string `json:"container_disk_watch_interval"`
	Containerizers                   string `json:"containerizers"`
	DefaultRole                      string `json:"default_role"`
	DiskWatchInterval                string `json:"disk_watch_interval"`
	Docker                           string `json:"docker"`
	DockerKillOrphans                string `json:"docker_kill_orphans"`
	DockerRegistry                   string `json:"docker_registry"`
	DockerRemoveDelay                string `json:"docker_remove_delay"`
	DockerSandboxDirectory           string `json:"docker_sandbox_directory"`
	DockerSocket                     string `json:"docker_socket"`
	DockerStoreDir                   string `json:"docker_store_dir"`
	DockerStopTimeout                string `json:"docker_stop_timeout"`
	EnforceContainerDiskQuota        string `json:"enforce_container_disk_quota"`
	ExecutorRegistrationTimeout      string `json:"executor_registration_timeout"`
	ExecutorShutdownGracePeriod      string `json:"executor_shutdown_grace_period"`
	FetcherCacheDir                  string `json:"fetcher_cache_dir"`
	FetcherCacheSize                 string `json:"fetcher_cache_size"`
	FrameworksHome                   string `json:"frameworks_home"`
	FrameworkSorter                  string `json:"framework_sorter"`
	GCDelay                          string `json:"gc_delay"`
	GCDiskHeadroom                   string `json:"gc_disk_headroom"`
	HadoopHome                       string `json:"hadoop_home"`
	Help                             string `json:"help"`
	Hostname                         string `json:"hostname"`
	HostnameLookup                   string `json:"hostname_lookup"`
	HTTPAuthenticators               string `json:"http_authenticators"`
	ImageProvisionerBackend          string `json:"image_provisioner_backend"`
	InitializeDriverLogging          string `json:"initialize_driver_logging"`
	IP                               string `json:"ip"`
	Isolation                        string `json:"isolation"`
	LauncherDir                      string `json:"launcher_dir"`
	LogAutoInitialize                string `json:"log_auto_initialize"`
	LogDir                           string `json:"log_dir"`
	Logbufsecs                       string `json:"logbufsecs"`
	LoggingLevel                     string `json:"logging_level"`
	MaxCompletedFrameworks           string `json:"max_completed_frameworks"`
	MaxCompletedTasksPerFramework    string `json:"max_completed_tasks_per_framework"`
	MaxSlavePingTimeouts             string `json:"max_slave_ping_timeouts"`
	Master                           string `json:"master"`
	PerfDuration                     string `json:"perf_duration"`
	PerfInterval                     string `json:"perf_interval"`
	Port                             string `json:"port"`
	Quiet                            string `json:"quiet"`
	Quorum                           string `json:"quorum"`
	QOSCorrectionIntervalMin         string `json:"qos_correction_interval_min"`
	Recover                          string `json:"recover"`
	RevocableCPULowPriority          string `json:"revocable_cpu_low_priority"`
	RecoverySlaveRemovalLimit        string `json:"recovery_slave_removal_limit"`
	RecoveryTimeout                  string `json:"recovery_timeout"`
	RegistrationBackoffFactor        string `json:"registration_backoff_factor"`
	Registry                         string `json:"registry"`
	RegistryFetchTimeout             string `json:"registry_fetch_timeout"`
	RegistryStoreTimeout             string `json:"registry_store_timeout"`
	RegistryStrict                   string `json:"registry_strict"`
	ResourceMonitoringInterval       string `json:"resource_monitoring_interval"`
	RootSubmissions                  string `json:"root_submissions"`
	SandboxDirectory                 string `json:"sandbox_directory"`
	SlavePingTimeout                 string `json:"slave_ping_timeout"`
	SlaveReregisterTimeout           string `json:"slave_reregister_timeout"`
	Strict                           string `json:"strict"`
	SystemdRuntimeDirectory          string `json:"systemd_runtime_directory"`
	SwitchUser                       string `json:"switch_user"`
	OversubscribedResourcesInterval  string `json:"oversubscribed_resources_interval"`
	UserSorter                       string `json:"user_sorter"`
	Version                          string `json:"version"`
	WebuiDir                         string `json:"webui_dir"`
	WorkDir                          string `json:"work_dir"`
	ZK                               string `json:"zk"`
	ZKSessionTimeout                 string `json:"zk_session_timeout"`
}

// Framework represent a single framework of a mesos node
type Framework struct {
	Active             bool       `json:"active"`
	Checkpoint         bool       `json:"checkpoint"`
	CompletedTasks     []Task     `json:"completed_tasks"`
	Executors          []Executor `json:"executors"`
	CompletedExecutors []Executor `json:"completed_executors"`
	FailoverTimeout    float64    `json:"failover_timeout"`
	Hostname           string     `json:"hostname"`
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	PID                string     `json:"pid"`
	OfferedResources   Resources  `json:"offered_resources"`
	Offers             []Offer    `json:"offers"`
	RegisteredTime     float64    `json:"registered_time"`
	ReregisteredTime   float64    `json:"reregistered_time"`
	Resources          Resources  `json:"resources"`
	Role               string     `json:"role"`
	Tasks              []Task     `json:"tasks"`
	UnregisteredTime   float64    `json:"unregistered_time"`
	UsedResources      Resources  `json:"used_resources"`
	User               string     `json:"user"`
	WebuiURL           string     `json:"webui_url"`
	Labels             []Label    `json:"label"`
	// Missing fields
	// TODO: "capabilities": [],
}

// Offer represents a single offer from a Mesos Slave to a Mesos master
type Offer struct {
	ID          string            `json:"id"`
	FrameworkID string            `json:"framework_id"`
	SlaveID     string            `json:"slave_id"`
	Hostname    string            `json:"hostname"`
	URL         URL               `json:"url"`
	Resources   Resources         `json:"resources"`
	Attributes  map[string]string `json:"attributes"`
}

// URL represents a single URL
type URL struct {
	Scheme     string      `json:"scheme"`
	Address    Address     `json:"address"`
	Path       string      `json:"path"`
	Parameters []Parameter `json:"parameters"`
}

// Address represents a single address.
// e.g. from a Slave or from a Master
type Address struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
}

// Parameter represents a single key / value pair for parameters
type Parameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Label represents a single key / value pair for labeling
type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Task represent a single Mesos task
type Task struct {
	// Missing fields
	// TODO: "labels": [],
	ExecutorID  string        `json:"executor_id"`
	FrameworkID string        `json:"framework_id"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Resources   Resources     `json:"resources"`
	SlaveID     string        `json:"slave_id"`
	State       string        `json:"state"`
	Statuses    []TaskStatus  `json:"statuses"`
	Discovery   TaskDiscovery `json:"discovery"`
}

// TaskDiscovery represents the dicovery information of a task
type TaskDiscovery struct {
	Visibility string `json:"visibility"`
	Name       string `json:"name"`
	Ports      Ports  `json:"ports"`
}

// Ports represents a number of PortDetails
type Ports struct {
	Ports []PortDetails `json:"ports"`
}

// PortDetails represents details about a single port
type PortDetails struct {
	Number   int    `json:"number"`
	Protocol string `json:"protocol"`
}

// Resources represents a resource type for a task
type Resources struct {
	CPUs  float64 `json:"cpus"`
	Disk  float64 `json:"disk"`
	Mem   float64 `json:"mem"`
	Ports string  `json:"ports"`
}

// TaskStatus represents the status of a single task
type TaskStatus struct {
	State           string          `json:"state"`
	Timestamp       float64         `json:"timestamp"`
	ContainerStatus ContainerStatus `json:"container_status"`
}

// ContainerStatus represents the status of a single container inside a task
type ContainerStatus struct {
	NetworkInfos []NetworkInfo `json:"network_infos"`
}

// NetworkInfo represents information about the network of a container
type NetworkInfo struct {
	IpAddress   string      `json:"ip_address"`
	IpAddresses []IpAddress `json:"ip_addresses"`
}

// IpAddress represents a single IpAddress
type IpAddress struct {
	IpAddress string `json:"ip_address"`
}

// Slave represents a single mesos slave node
type Slave struct {
	Active              bool                   `json:"active"`
	Hostname            string                 `json:"hostname"`
	ID                  string                 `json:"id"`
	PID                 string                 `json:"pid"`
	RegisteredTime      float64                `json:"registered_time"`
	Resources           Resources              `json:"resources"`
	UsedResources       Resources              `json:"used_resources"`
	OfferedResources    Resources              `json:"offered_resources"`
	ReservedResources   Resources              `json:"reserved_resources"`
	UnreservedResources Resources              `json:"unreserved_resources"`
	Attributes          map[string]interface{} `json:"attributes"`
	Version             string                 `json:"version"`
}

// Executor represents a single executor of a framework
type Executor struct {
	CompletedTasks []Task    `json:"completed_tasks"`
	Container      string    `json:"container"`
	Directory      string    `json:"directory"`
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Resources      Resources `json:"resources"`
	Source         string    `json:"source"`
	QueuedTasks    []Task    `json:"queued_tasks"`
	Tasks          []Task    `json:"tasks"`
}
