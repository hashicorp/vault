package api

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	// RestartPolicyModeDelay causes an artificial delay till the next interval is
	// reached when the specified attempts have been reached in the interval.
	RestartPolicyModeDelay = "delay"

	// RestartPolicyModeFail causes a job to fail if the specified number of
	// attempts are reached within an interval.
	RestartPolicyModeFail = "fail"
)

// MemoryStats holds memory usage related stats
type MemoryStats struct {
	RSS            uint64
	Cache          uint64
	Swap           uint64
	Usage          uint64
	MaxUsage       uint64
	KernelUsage    uint64
	KernelMaxUsage uint64
	Measured       []string
}

// CpuStats holds cpu usage related stats
type CpuStats struct {
	SystemMode       float64
	UserMode         float64
	TotalTicks       float64
	ThrottledPeriods uint64
	ThrottledTime    uint64
	Percent          float64
	Measured         []string
}

// ResourceUsage holds information related to cpu and memory stats
type ResourceUsage struct {
	MemoryStats *MemoryStats
	CpuStats    *CpuStats
	DeviceStats []*DeviceGroupStats
}

// TaskResourceUsage holds aggregated resource usage of all processes in a Task
// and the resource usage of the individual pids
type TaskResourceUsage struct {
	ResourceUsage *ResourceUsage
	Timestamp     int64
	Pids          map[string]*ResourceUsage
}

// AllocResourceUsage holds the aggregated task resource usage of the
// allocation.
type AllocResourceUsage struct {
	ResourceUsage *ResourceUsage
	Tasks         map[string]*TaskResourceUsage
	Timestamp     int64
}

// RestartPolicy defines how the Nomad client restarts
// tasks in a taskgroup when they fail
type RestartPolicy struct {
	Interval *time.Duration `hcl:"interval,optional"`
	Attempts *int           `hcl:"attempts,optional"`
	Delay    *time.Duration `hcl:"delay,optional"`
	Mode     *string        `hcl:"mode,optional"`
}

func (r *RestartPolicy) Merge(rp *RestartPolicy) {
	if rp.Interval != nil {
		r.Interval = rp.Interval
	}
	if rp.Attempts != nil {
		r.Attempts = rp.Attempts
	}
	if rp.Delay != nil {
		r.Delay = rp.Delay
	}
	if rp.Mode != nil {
		r.Mode = rp.Mode
	}
}

// Reschedule configures how Tasks are rescheduled  when they crash or fail.
type ReschedulePolicy struct {
	// Attempts limits the number of rescheduling attempts that can occur in an interval.
	Attempts *int `mapstructure:"attempts" hcl:"attempts,optional"`

	// Interval is a duration in which we can limit the number of reschedule attempts.
	Interval *time.Duration `mapstructure:"interval" hcl:"interval,optional"`

	// Delay is a minimum duration to wait between reschedule attempts.
	// The delay function determines how much subsequent reschedule attempts are delayed by.
	Delay *time.Duration `mapstructure:"delay" hcl:"delay,optional"`

	// DelayFunction determines how the delay progressively changes on subsequent reschedule
	// attempts. Valid values are "exponential", "constant", and "fibonacci".
	DelayFunction *string `mapstructure:"delay_function" hcl:"delay_function,optional"`

	// MaxDelay is an upper bound on the delay.
	MaxDelay *time.Duration `mapstructure:"max_delay" hcl:"max_delay,optional"`

	// Unlimited allows rescheduling attempts until they succeed
	Unlimited *bool `mapstructure:"unlimited" hcl:"unlimited,optional"`
}

func (r *ReschedulePolicy) Merge(rp *ReschedulePolicy) {
	if rp == nil {
		return
	}
	if rp.Interval != nil {
		r.Interval = rp.Interval
	}
	if rp.Attempts != nil {
		r.Attempts = rp.Attempts
	}
	if rp.Delay != nil {
		r.Delay = rp.Delay
	}
	if rp.DelayFunction != nil {
		r.DelayFunction = rp.DelayFunction
	}
	if rp.MaxDelay != nil {
		r.MaxDelay = rp.MaxDelay
	}
	if rp.Unlimited != nil {
		r.Unlimited = rp.Unlimited
	}
}

func (r *ReschedulePolicy) Canonicalize(jobType string) {
	dp := NewDefaultReschedulePolicy(jobType)
	if r.Interval == nil {
		r.Interval = dp.Interval
	}
	if r.Attempts == nil {
		r.Attempts = dp.Attempts
	}
	if r.Delay == nil {
		r.Delay = dp.Delay
	}
	if r.DelayFunction == nil {
		r.DelayFunction = dp.DelayFunction
	}
	if r.MaxDelay == nil {
		r.MaxDelay = dp.MaxDelay
	}
	if r.Unlimited == nil {
		r.Unlimited = dp.Unlimited
	}
}

// Affinity is used to serialize task group affinities
type Affinity struct {
	LTarget string `hcl:"attribute,optional"` // Left-hand target
	RTarget string `hcl:"value,optional"`     // Right-hand target
	Operand string `hcl:"operator,optional"`  // Constraint operand (<=, <, =, !=, >, >=), set_contains_all, set_contains_any
	Weight  *int8  `hcl:"weight,optional"`    // Weight applied to nodes that match the affinity. Can be negative
}

func NewAffinity(LTarget string, Operand string, RTarget string, Weight int8) *Affinity {
	return &Affinity{
		LTarget: LTarget,
		RTarget: RTarget,
		Operand: Operand,
		Weight:  int8ToPtr(Weight),
	}
}

func (a *Affinity) Canonicalize() {
	if a.Weight == nil {
		a.Weight = int8ToPtr(50)
	}
}

func NewDefaultReschedulePolicy(jobType string) *ReschedulePolicy {
	var dp *ReschedulePolicy
	switch jobType {
	case "service":
		// This needs to be in sync with DefaultServiceJobReschedulePolicy
		// in nomad/structs/structs.go
		dp = &ReschedulePolicy{
			Delay:         timeToPtr(30 * time.Second),
			DelayFunction: stringToPtr("exponential"),
			MaxDelay:      timeToPtr(1 * time.Hour),
			Unlimited:     boolToPtr(true),

			Attempts: intToPtr(0),
			Interval: timeToPtr(0),
		}
	case "batch":
		// This needs to be in sync with DefaultBatchJobReschedulePolicy
		// in nomad/structs/structs.go
		dp = &ReschedulePolicy{
			Attempts:      intToPtr(1),
			Interval:      timeToPtr(24 * time.Hour),
			Delay:         timeToPtr(5 * time.Second),
			DelayFunction: stringToPtr("constant"),

			MaxDelay:  timeToPtr(0),
			Unlimited: boolToPtr(false),
		}

	case "system":
		dp = &ReschedulePolicy{
			Attempts:      intToPtr(0),
			Interval:      timeToPtr(0),
			Delay:         timeToPtr(0),
			DelayFunction: stringToPtr(""),
			MaxDelay:      timeToPtr(0),
			Unlimited:     boolToPtr(false),
		}

	default:
		// GH-7203: it is possible an unknown job type is passed to this
		// function and we need to ensure a non-nil object is returned so that
		// the canonicalization runs without panicking.
		dp = &ReschedulePolicy{
			Attempts:      intToPtr(0),
			Interval:      timeToPtr(0),
			Delay:         timeToPtr(0),
			DelayFunction: stringToPtr(""),
			MaxDelay:      timeToPtr(0),
			Unlimited:     boolToPtr(false),
		}
	}
	return dp
}

func (r *ReschedulePolicy) Copy() *ReschedulePolicy {
	if r == nil {
		return nil
	}
	nrp := new(ReschedulePolicy)
	*nrp = *r
	return nrp
}

func (p *ReschedulePolicy) String() string {
	if p == nil {
		return ""
	}
	if *p.Unlimited {
		return fmt.Sprintf("unlimited with %v delay, max_delay = %v", *p.DelayFunction, *p.MaxDelay)
	}
	return fmt.Sprintf("%v in %v with %v delay, max_delay = %v", *p.Attempts, *p.Interval, *p.DelayFunction, *p.MaxDelay)
}

// Spread is used to serialize task group allocation spread preferences
type Spread struct {
	Attribute    string          `hcl:"attribute,optional"`
	Weight       *int8           `hcl:"weight,optional"`
	SpreadTarget []*SpreadTarget `hcl:"target,block"`
}

// SpreadTarget is used to serialize target allocation spread percentages
type SpreadTarget struct {
	Value   string `hcl:",label"`
	Percent uint8  `hcl:"percent,optional"`
}

func NewSpreadTarget(value string, percent uint8) *SpreadTarget {
	return &SpreadTarget{
		Value:   value,
		Percent: percent,
	}
}

func NewSpread(attribute string, weight int8, spreadTargets []*SpreadTarget) *Spread {
	return &Spread{
		Attribute:    attribute,
		Weight:       int8ToPtr(weight),
		SpreadTarget: spreadTargets,
	}
}

func (s *Spread) Canonicalize() {
	if s.Weight == nil {
		s.Weight = int8ToPtr(50)
	}
}

// EphemeralDisk is an ephemeral disk object
type EphemeralDisk struct {
	Sticky  *bool `hcl:"sticky,optional"`
	Migrate *bool `hcl:"migrate,optional"`
	SizeMB  *int  `mapstructure:"size" hcl:"size,optional"`
}

func DefaultEphemeralDisk() *EphemeralDisk {
	return &EphemeralDisk{
		Sticky:  boolToPtr(false),
		Migrate: boolToPtr(false),
		SizeMB:  intToPtr(300),
	}
}

func (e *EphemeralDisk) Canonicalize() {
	if e.Sticky == nil {
		e.Sticky = boolToPtr(false)
	}
	if e.Migrate == nil {
		e.Migrate = boolToPtr(false)
	}
	if e.SizeMB == nil {
		e.SizeMB = intToPtr(300)
	}
}

// MigrateStrategy describes how allocations for a task group should be
// migrated between nodes (eg when draining).
type MigrateStrategy struct {
	MaxParallel     *int           `mapstructure:"max_parallel" hcl:"max_parallel,optional"`
	HealthCheck     *string        `mapstructure:"health_check" hcl:"health_check,optional"`
	MinHealthyTime  *time.Duration `mapstructure:"min_healthy_time" hcl:"min_healthy_time,optional"`
	HealthyDeadline *time.Duration `mapstructure:"healthy_deadline" hcl:"healthy_deadline,optional"`
}

func DefaultMigrateStrategy() *MigrateStrategy {
	return &MigrateStrategy{
		MaxParallel:     intToPtr(1),
		HealthCheck:     stringToPtr("checks"),
		MinHealthyTime:  timeToPtr(10 * time.Second),
		HealthyDeadline: timeToPtr(5 * time.Minute),
	}
}

func (m *MigrateStrategy) Canonicalize() {
	if m == nil {
		return
	}
	defaults := DefaultMigrateStrategy()
	if m.MaxParallel == nil {
		m.MaxParallel = defaults.MaxParallel
	}
	if m.HealthCheck == nil {
		m.HealthCheck = defaults.HealthCheck
	}
	if m.MinHealthyTime == nil {
		m.MinHealthyTime = defaults.MinHealthyTime
	}
	if m.HealthyDeadline == nil {
		m.HealthyDeadline = defaults.HealthyDeadline
	}
}

func (m *MigrateStrategy) Merge(o *MigrateStrategy) {
	if o.MaxParallel != nil {
		m.MaxParallel = o.MaxParallel
	}
	if o.HealthCheck != nil {
		m.HealthCheck = o.HealthCheck
	}
	if o.MinHealthyTime != nil {
		m.MinHealthyTime = o.MinHealthyTime
	}
	if o.HealthyDeadline != nil {
		m.HealthyDeadline = o.HealthyDeadline
	}
}

func (m *MigrateStrategy) Copy() *MigrateStrategy {
	if m == nil {
		return nil
	}
	nm := new(MigrateStrategy)
	*nm = *m
	return nm
}

// VolumeRequest is a representation of a storage volume that a TaskGroup wishes to use.
type VolumeRequest struct {
	Name         string           `hcl:"name,label"`
	Type         string           `hcl:"type,optional"`
	Source       string           `hcl:"source,optional"`
	ReadOnly     bool             `hcl:"read_only,optional"`
	MountOptions *CSIMountOptions `hcl:"mount_options,block"`
	ExtraKeysHCL []string         `hcl1:",unusedKeys,optional" json:"-"`
}

const (
	VolumeMountPropagationPrivate       = "private"
	VolumeMountPropagationHostToTask    = "host-to-task"
	VolumeMountPropagationBidirectional = "bidirectional"
)

// VolumeMount represents the relationship between a destination path in a task
// and the task group volume that should be mounted there.
type VolumeMount struct {
	Volume          *string `hcl:"volume,optional"`
	Destination     *string `hcl:"destination,optional"`
	ReadOnly        *bool   `mapstructure:"read_only" hcl:"read_only,optional"`
	PropagationMode *string `mapstructure:"propagation_mode" hcl:"propagation_mode,optional"`
}

func (vm *VolumeMount) Canonicalize() {
	if vm.PropagationMode == nil {
		vm.PropagationMode = stringToPtr(VolumeMountPropagationPrivate)
	}
	if vm.ReadOnly == nil {
		vm.ReadOnly = boolToPtr(false)
	}
}

// TaskGroup is the unit of scheduling.
type TaskGroup struct {
	Name                      *string                   `hcl:"name,label"`
	Count                     *int                      `hcl:"count,optional"`
	Constraints               []*Constraint             `hcl:"constraint,block"`
	Affinities                []*Affinity               `hcl:"affinity,block"`
	Tasks                     []*Task                   `hcl:"task,block"`
	Spreads                   []*Spread                 `hcl:"spread,block"`
	Volumes                   map[string]*VolumeRequest `hcl:"volume,block"`
	RestartPolicy             *RestartPolicy            `hcl:"restart,block"`
	ReschedulePolicy          *ReschedulePolicy         `hcl:"reschedule,block"`
	EphemeralDisk             *EphemeralDisk            `hcl:"ephemeral_disk,block"`
	Update                    *UpdateStrategy           `hcl:"update,block"`
	Migrate                   *MigrateStrategy          `hcl:"migrate,block"`
	Networks                  []*NetworkResource        `hcl:"network,block"`
	Meta                      map[string]string         `hcl:"meta,block"`
	Services                  []*Service                `hcl:"service,block"`
	ShutdownDelay             *time.Duration            `mapstructure:"shutdown_delay" hcl:"shutdown_delay,optional"`
	StopAfterClientDisconnect *time.Duration            `mapstructure:"stop_after_client_disconnect" hcl:"stop_after_client_disconnect,optional"`
	Scaling                   *ScalingPolicy            `hcl:"scaling,block"`
}

// NewTaskGroup creates a new TaskGroup.
func NewTaskGroup(name string, count int) *TaskGroup {
	return &TaskGroup{
		Name:  stringToPtr(name),
		Count: intToPtr(count),
	}
}

// Canonicalize sets defaults and merges settings that should be inherited from the job
func (g *TaskGroup) Canonicalize(job *Job) {
	if g.Name == nil {
		g.Name = stringToPtr("")
	}

	if g.Count == nil {
		if g.Scaling != nil && g.Scaling.Min != nil {
			g.Count = intToPtr(int(*g.Scaling.Min))
		} else {
			g.Count = intToPtr(1)
		}
	}
	if g.Scaling != nil {
		g.Scaling.Canonicalize(*g.Count)
	}
	if g.EphemeralDisk == nil {
		g.EphemeralDisk = DefaultEphemeralDisk()
	} else {
		g.EphemeralDisk.Canonicalize()
	}

	// Merge the update policy from the job
	if ju, tu := job.Update != nil, g.Update != nil; ju && tu {
		// Merge the jobs and task groups definition of the update strategy
		jc := job.Update.Copy()
		jc.Merge(g.Update)
		g.Update = jc
	} else if ju && !job.Update.Empty() {
		// Inherit the jobs as long as it is non-empty.
		jc := job.Update.Copy()
		g.Update = jc
	}

	if g.Update != nil {
		g.Update.Canonicalize()
	}

	// Merge the reschedule policy from the job
	if jr, tr := job.Reschedule != nil, g.ReschedulePolicy != nil; jr && tr {
		jobReschedule := job.Reschedule.Copy()
		jobReschedule.Merge(g.ReschedulePolicy)
		g.ReschedulePolicy = jobReschedule
	} else if jr {
		jobReschedule := job.Reschedule.Copy()
		g.ReschedulePolicy = jobReschedule
	}
	// Only use default reschedule policy for non system jobs
	if g.ReschedulePolicy == nil && *job.Type != "system" {
		g.ReschedulePolicy = NewDefaultReschedulePolicy(*job.Type)
	}
	if g.ReschedulePolicy != nil {
		g.ReschedulePolicy.Canonicalize(*job.Type)
	}
	// Merge the migrate strategy from the job
	if jm, tm := job.Migrate != nil, g.Migrate != nil; jm && tm {
		jobMigrate := job.Migrate.Copy()
		jobMigrate.Merge(g.Migrate)
		g.Migrate = jobMigrate
	} else if jm {
		jobMigrate := job.Migrate.Copy()
		g.Migrate = jobMigrate
	}

	// Merge with default reschedule policy
	if g.Migrate == nil && *job.Type == "service" {
		g.Migrate = &MigrateStrategy{}
	}
	if g.Migrate != nil {
		g.Migrate.Canonicalize()
	}

	var defaultRestartPolicy *RestartPolicy
	switch *job.Type {
	case "service", "system":
		defaultRestartPolicy = defaultServiceJobRestartPolicy()
	default:
		defaultRestartPolicy = defaultBatchJobRestartPolicy()
	}

	if g.RestartPolicy != nil {
		defaultRestartPolicy.Merge(g.RestartPolicy)
	}
	g.RestartPolicy = defaultRestartPolicy

	for _, t := range g.Tasks {
		t.Canonicalize(g, job)
	}

	for _, spread := range g.Spreads {
		spread.Canonicalize()
	}
	for _, a := range g.Affinities {
		a.Canonicalize()
	}
	for _, n := range g.Networks {
		n.Canonicalize()
	}
	for _, s := range g.Services {
		s.Canonicalize(nil, g, job)
	}
}

// These needs to be in sync with DefaultServiceJobRestartPolicy in
// in nomad/structs/structs.go
func defaultServiceJobRestartPolicy() *RestartPolicy {
	return &RestartPolicy{
		Delay:    timeToPtr(15 * time.Second),
		Attempts: intToPtr(2),
		Interval: timeToPtr(30 * time.Minute),
		Mode:     stringToPtr(RestartPolicyModeFail),
	}
}

// These needs to be in sync with DefaultBatchJobRestartPolicy in
// in nomad/structs/structs.go
func defaultBatchJobRestartPolicy() *RestartPolicy {
	return &RestartPolicy{
		Delay:    timeToPtr(15 * time.Second),
		Attempts: intToPtr(3),
		Interval: timeToPtr(24 * time.Hour),
		Mode:     stringToPtr(RestartPolicyModeFail),
	}
}

// Constrain is used to add a constraint to a task group.
func (g *TaskGroup) Constrain(c *Constraint) *TaskGroup {
	g.Constraints = append(g.Constraints, c)
	return g
}

// AddMeta is used to add a meta k/v pair to a task group
func (g *TaskGroup) SetMeta(key, val string) *TaskGroup {
	if g.Meta == nil {
		g.Meta = make(map[string]string)
	}
	g.Meta[key] = val
	return g
}

// AddTask is used to add a new task to a task group.
func (g *TaskGroup) AddTask(t *Task) *TaskGroup {
	g.Tasks = append(g.Tasks, t)
	return g
}

// AddAffinity is used to add a new affinity to a task group.
func (g *TaskGroup) AddAffinity(a *Affinity) *TaskGroup {
	g.Affinities = append(g.Affinities, a)
	return g
}

// RequireDisk adds a ephemeral disk to the task group
func (g *TaskGroup) RequireDisk(disk *EphemeralDisk) *TaskGroup {
	g.EphemeralDisk = disk
	return g
}

// AddSpread is used to add a new spread preference to a task group.
func (g *TaskGroup) AddSpread(s *Spread) *TaskGroup {
	g.Spreads = append(g.Spreads, s)
	return g
}

// LogConfig provides configuration for log rotation
type LogConfig struct {
	MaxFiles      *int `mapstructure:"max_files" hcl:"max_files,optional"`
	MaxFileSizeMB *int `mapstructure:"max_file_size" hcl:"max_file_size,optional"`
}

func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		MaxFiles:      intToPtr(10),
		MaxFileSizeMB: intToPtr(10),
	}
}

func (l *LogConfig) Canonicalize() {
	if l.MaxFiles == nil {
		l.MaxFiles = intToPtr(10)
	}
	if l.MaxFileSizeMB == nil {
		l.MaxFileSizeMB = intToPtr(10)
	}
}

// DispatchPayloadConfig configures how a task gets its input from a job dispatch
type DispatchPayloadConfig struct {
	File string `hcl:"file,optional"`
}

const (
	TaskLifecycleHookPrestart  = "prestart"
	TaskLifecycleHookPoststart = "poststart"
	TaskLifecycleHookPoststop  = "poststop"
)

type TaskLifecycle struct {
	Hook    string `mapstructure:"hook" hcl:"hook,optional"`
	Sidecar bool   `mapstructure:"sidecar" hcl:"sidecar,optional"`
}

// Determine if lifecycle has user-input values
func (l *TaskLifecycle) Empty() bool {
	return l == nil || (l.Hook == "")
}

// Task is a single process in a task group.
type Task struct {
	Name            string                 `hcl:"name,label"`
	Driver          string                 `hcl:"driver,optional"`
	User            string                 `hcl:"user,optional"`
	Lifecycle       *TaskLifecycle         `hcl:"lifecycle,block"`
	Config          map[string]interface{} `hcl:"config,block"`
	Constraints     []*Constraint          `hcl:"constraint,block"`
	Affinities      []*Affinity            `hcl:"affinity,block"`
	Env             map[string]string      `hcl:"env,block"`
	Services        []*Service             `hcl:"service,block"`
	Resources       *Resources             `hcl:"resources,block"`
	RestartPolicy   *RestartPolicy         `hcl:"restart,block"`
	Meta            map[string]string      `hcl:"meta,block"`
	KillTimeout     *time.Duration         `mapstructure:"kill_timeout" hcl:"kill_timeout,optional"`
	LogConfig       *LogConfig             `mapstructure:"logs" hcl:"logs,block"`
	Artifacts       []*TaskArtifact        `hcl:"artifact,block"`
	Vault           *Vault                 `hcl:"vault,block"`
	Templates       []*Template            `hcl:"template,block"`
	DispatchPayload *DispatchPayloadConfig `hcl:"dispatch_payload,block"`
	VolumeMounts    []*VolumeMount         `hcl:"volume_mount,block"`
	CSIPluginConfig *TaskCSIPluginConfig   `mapstructure:"csi_plugin" json:",omitempty" hcl:"csi_plugin,block"`
	Leader          bool                   `hcl:"leader,optional"`
	ShutdownDelay   time.Duration          `mapstructure:"shutdown_delay" hcl:"shutdown_delay,optional"`
	KillSignal      string                 `mapstructure:"kill_signal" hcl:"kill_signal,optional"`
	Kind            string                 `hcl:"kind,optional"`
	ScalingPolicies []*ScalingPolicy       `hcl:"scaling,block"`
}

func (t *Task) Canonicalize(tg *TaskGroup, job *Job) {
	if t.Resources == nil {
		t.Resources = &Resources{}
	}
	t.Resources.Canonicalize()

	if t.KillTimeout == nil {
		t.KillTimeout = timeToPtr(5 * time.Second)
	}
	if t.LogConfig == nil {
		t.LogConfig = DefaultLogConfig()
	} else {
		t.LogConfig.Canonicalize()
	}
	for _, artifact := range t.Artifacts {
		artifact.Canonicalize()
	}
	if t.Vault != nil {
		t.Vault.Canonicalize()
	}
	for _, tmpl := range t.Templates {
		tmpl.Canonicalize()
	}
	for _, s := range t.Services {
		s.Canonicalize(t, tg, job)
	}
	for _, a := range t.Affinities {
		a.Canonicalize()
	}
	for _, vm := range t.VolumeMounts {
		vm.Canonicalize()
	}
	if t.Lifecycle.Empty() {
		t.Lifecycle = nil
	}
	if t.CSIPluginConfig != nil {
		t.CSIPluginConfig.Canonicalize()
	}
	if t.RestartPolicy == nil {
		t.RestartPolicy = tg.RestartPolicy
	} else {
		tgrp := &RestartPolicy{}
		*tgrp = *tg.RestartPolicy
		tgrp.Merge(t.RestartPolicy)
		t.RestartPolicy = tgrp
	}
}

// TaskArtifact is used to download artifacts before running a task.
type TaskArtifact struct {
	GetterSource  *string           `mapstructure:"source" hcl:"source,optional"`
	GetterOptions map[string]string `mapstructure:"options" hcl:"options,block"`
	GetterHeaders map[string]string `mapstructure:"headers" hcl:"headers,block"`
	GetterMode    *string           `mapstructure:"mode" hcl:"mode,optional"`
	RelativeDest  *string           `mapstructure:"destination" hcl:"destination,optional"`
}

func (a *TaskArtifact) Canonicalize() {
	if a.GetterMode == nil {
		a.GetterMode = stringToPtr("any")
	}
	if a.GetterSource == nil {
		// Shouldn't be possible, but we don't want to panic
		a.GetterSource = stringToPtr("")
	}
	if len(a.GetterOptions) == 0 {
		a.GetterOptions = nil
	}
	if len(a.GetterHeaders) == 0 {
		a.GetterHeaders = nil
	}
	if a.RelativeDest == nil {
		switch *a.GetterMode {
		case "file":
			// File mode should default to local/filename
			dest := *a.GetterSource
			dest = path.Base(dest)
			dest = filepath.Join("local", dest)
			a.RelativeDest = &dest
		default:
			// Default to a directory
			a.RelativeDest = stringToPtr("local/")
		}
	}
}

type Template struct {
	SourcePath   *string        `mapstructure:"source" hcl:"source,optional"`
	DestPath     *string        `mapstructure:"destination" hcl:"destination,optional"`
	EmbeddedTmpl *string        `mapstructure:"data" hcl:"data,optional"`
	ChangeMode   *string        `mapstructure:"change_mode" hcl:"change_mode,optional"`
	ChangeSignal *string        `mapstructure:"change_signal" hcl:"change_signal,optional"`
	Splay        *time.Duration `mapstructure:"splay" hcl:"splay,optional"`
	Perms        *string        `mapstructure:"perms" hcl:"perms,optional"`
	LeftDelim    *string        `mapstructure:"left_delimiter" hcl:"left_delimiter,optional"`
	RightDelim   *string        `mapstructure:"right_delimiter" hcl:"right_delimiter,optional"`
	Envvars      *bool          `mapstructure:"env" hcl:"env,optional"`
	VaultGrace   *time.Duration `mapstructure:"vault_grace" hcl:"vault_grace,optional"`
}

func (tmpl *Template) Canonicalize() {
	if tmpl.SourcePath == nil {
		tmpl.SourcePath = stringToPtr("")
	}
	if tmpl.DestPath == nil {
		tmpl.DestPath = stringToPtr("")
	}
	if tmpl.EmbeddedTmpl == nil {
		tmpl.EmbeddedTmpl = stringToPtr("")
	}
	if tmpl.ChangeMode == nil {
		tmpl.ChangeMode = stringToPtr("restart")
	}
	if tmpl.ChangeSignal == nil {
		if *tmpl.ChangeMode == "signal" {
			tmpl.ChangeSignal = stringToPtr("SIGHUP")
		} else {
			tmpl.ChangeSignal = stringToPtr("")
		}
	} else {
		sig := *tmpl.ChangeSignal
		tmpl.ChangeSignal = stringToPtr(strings.ToUpper(sig))
	}
	if tmpl.Splay == nil {
		tmpl.Splay = timeToPtr(5 * time.Second)
	}
	if tmpl.Perms == nil {
		tmpl.Perms = stringToPtr("0644")
	}
	if tmpl.LeftDelim == nil {
		tmpl.LeftDelim = stringToPtr("{{")
	}
	if tmpl.RightDelim == nil {
		tmpl.RightDelim = stringToPtr("}}")
	}
	if tmpl.Envvars == nil {
		tmpl.Envvars = boolToPtr(false)
	}

	//COMPAT(0.12) VaultGrace is deprecated and unused as of Vault 0.5
	if tmpl.VaultGrace == nil {
		tmpl.VaultGrace = timeToPtr(0)
	}
}

type Vault struct {
	Policies     []string `hcl:"policies,optional"`
	Namespace    *string  `mapstructure:"namespace" hcl:"namespace,optional"`
	Env          *bool    `hcl:"env,optional"`
	ChangeMode   *string  `mapstructure:"change_mode" hcl:"change_mode,optional"`
	ChangeSignal *string  `mapstructure:"change_signal" hcl:"change_signal,optional"`
}

func (v *Vault) Canonicalize() {
	if v.Env == nil {
		v.Env = boolToPtr(true)
	}
	if v.Namespace == nil {
		v.Namespace = stringToPtr("")
	}
	if v.ChangeMode == nil {
		v.ChangeMode = stringToPtr("restart")
	}
	if v.ChangeSignal == nil {
		v.ChangeSignal = stringToPtr("SIGHUP")
	}
}

// NewTask creates and initializes a new Task.
func NewTask(name, driver string) *Task {
	return &Task{
		Name:   name,
		Driver: driver,
	}
}

// Configure is used to configure a single k/v pair on
// the task.
func (t *Task) SetConfig(key string, val interface{}) *Task {
	if t.Config == nil {
		t.Config = make(map[string]interface{})
	}
	t.Config[key] = val
	return t
}

// SetMeta is used to add metadata k/v pairs to the task.
func (t *Task) SetMeta(key, val string) *Task {
	if t.Meta == nil {
		t.Meta = make(map[string]string)
	}
	t.Meta[key] = val
	return t
}

// Require is used to add resource requirements to a task.
func (t *Task) Require(r *Resources) *Task {
	t.Resources = r
	return t
}

// Constraint adds a new constraints to a single task.
func (t *Task) Constrain(c *Constraint) *Task {
	t.Constraints = append(t.Constraints, c)
	return t
}

// AddAffinity adds a new affinity to a single task.
func (t *Task) AddAffinity(a *Affinity) *Task {
	t.Affinities = append(t.Affinities, a)
	return t
}

// SetLogConfig sets a log config to a task
func (t *Task) SetLogConfig(l *LogConfig) *Task {
	t.LogConfig = l
	return t
}

// TaskState tracks the current state of a task and events that caused state
// transitions.
type TaskState struct {
	State       string
	Failed      bool
	Restarts    uint64
	LastRestart time.Time
	StartedAt   time.Time
	FinishedAt  time.Time
	Events      []*TaskEvent
}

const (
	TaskSetup                  = "Task Setup"
	TaskSetupFailure           = "Setup Failure"
	TaskDriverFailure          = "Driver Failure"
	TaskDriverMessage          = "Driver"
	TaskReceived               = "Received"
	TaskFailedValidation       = "Failed Validation"
	TaskStarted                = "Started"
	TaskTerminated             = "Terminated"
	TaskKilling                = "Killing"
	TaskKilled                 = "Killed"
	TaskRestarting             = "Restarting"
	TaskNotRestarting          = "Not Restarting"
	TaskDownloadingArtifacts   = "Downloading Artifacts"
	TaskArtifactDownloadFailed = "Failed Artifact Download"
	TaskSiblingFailed          = "Sibling Task Failed"
	TaskSignaling              = "Signaling"
	TaskRestartSignal          = "Restart Signaled"
	TaskLeaderDead             = "Leader Task Dead"
	TaskBuildingTaskDir        = "Building Task Directory"
)

// TaskEvent is an event that effects the state of a task and contains meta-data
// appropriate to the events type.
type TaskEvent struct {
	Type           string
	Time           int64
	DisplayMessage string
	Details        map[string]string
	Message        string
	// DEPRECATION NOTICE: The following fields are all deprecated. see TaskEvent struct in structs.go for details.
	FailsTask        bool
	RestartReason    string
	SetupError       string
	DriverError      string
	DriverMessage    string
	ExitCode         int
	Signal           int
	KillReason       string
	KillTimeout      time.Duration
	KillError        string
	StartDelay       int64
	DownloadError    string
	ValidationError  string
	DiskLimit        int64
	DiskSize         int64
	FailedSibling    string
	VaultError       string
	TaskSignalReason string
	TaskSignal       string
	GenericSource    string
}

// CSIPluginType is an enum string that encapsulates the valid options for a
// CSIPlugin stanza's Type. These modes will allow the plugin to be used in
// different ways by the client.
type CSIPluginType string

const (
	// CSIPluginTypeNode indicates that Nomad should only use the plugin for
	// performing Node RPCs against the provided plugin.
	CSIPluginTypeNode CSIPluginType = "node"

	// CSIPluginTypeController indicates that Nomad should only use the plugin for
	// performing Controller RPCs against the provided plugin.
	CSIPluginTypeController CSIPluginType = "controller"

	// CSIPluginTypeMonolith indicates that Nomad can use the provided plugin for
	// both controller and node rpcs.
	CSIPluginTypeMonolith CSIPluginType = "monolith"
)

// TaskCSIPluginConfig contains the data that is required to setup a task as a
// CSI plugin. This will be used by the csi_plugin_supervisor_hook to configure
// mounts for the plugin and initiate the connection to the plugin catalog.
type TaskCSIPluginConfig struct {
	// ID is the identifier of the plugin.
	// Ideally this should be the FQDN of the plugin.
	ID string `mapstructure:"id" hcl:"id,optional"`

	// CSIPluginType instructs Nomad on how to handle processing a plugin
	Type CSIPluginType `mapstructure:"type" hcl:"type,optional"`

	// MountDir is the destination that nomad should mount in its CSI
	// directory for the plugin. It will then expect a file called CSISocketName
	// to be created by the plugin, and will provide references into
	// "MountDir/CSIIntermediaryDirname/VolumeName/AllocID for mounts.
	//
	// Default is /csi.
	MountDir string `mapstructure:"mount_dir" hcl:"mount_dir,optional"`
}

func (t *TaskCSIPluginConfig) Canonicalize() {
	if t.MountDir == "" {
		t.MountDir = "/csi"
	}
}
