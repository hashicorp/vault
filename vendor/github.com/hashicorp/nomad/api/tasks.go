package api

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/nomad/helper"
)

// MemoryStats holds memory usage related stats
type MemoryStats struct {
	RSS            uint64
	Cache          uint64
	Swap           uint64
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
	Interval *time.Duration
	Attempts *int
	Delay    *time.Duration
	Mode     *string
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

// CheckRestart describes if and when a task should be restarted based on
// failing health checks.
type CheckRestart struct {
	Limit          int            `mapstructure:"limit"`
	Grace          *time.Duration `mapstructure:"grace"`
	IgnoreWarnings bool           `mapstructure:"ignore_warnings"`
}

// Canonicalize CheckRestart fields if not nil.
func (c *CheckRestart) Canonicalize() {
	if c == nil {
		return
	}

	if c.Grace == nil {
		c.Grace = helper.TimeToPtr(1 * time.Second)
	}
}

// Copy returns a copy of CheckRestart or nil if unset.
func (c *CheckRestart) Copy() *CheckRestart {
	if c == nil {
		return nil
	}

	nc := new(CheckRestart)
	nc.Limit = c.Limit
	if c.Grace != nil {
		g := *c.Grace
		nc.Grace = &g
	}
	nc.IgnoreWarnings = c.IgnoreWarnings
	return nc
}

// Merge values from other CheckRestart over default values on this
// CheckRestart and return merged copy.
func (c *CheckRestart) Merge(o *CheckRestart) *CheckRestart {
	if c == nil {
		// Just return other
		return o
	}

	nc := c.Copy()

	if o == nil {
		// Nothing to merge
		return nc
	}

	if o.Limit > 0 {
		nc.Limit = o.Limit
	}

	if o.Grace != nil {
		nc.Grace = o.Grace
	}

	if o.IgnoreWarnings {
		nc.IgnoreWarnings = o.IgnoreWarnings
	}

	return nc
}

// The ServiceCheck data model represents the consul health check that
// Nomad registers for a Task
type ServiceCheck struct {
	Id            string
	Name          string
	Type          string
	Command       string
	Args          []string
	Path          string
	Protocol      string
	PortLabel     string `mapstructure:"port"`
	AddressMode   string `mapstructure:"address_mode"`
	Interval      time.Duration
	Timeout       time.Duration
	InitialStatus string `mapstructure:"initial_status"`
	TLSSkipVerify bool   `mapstructure:"tls_skip_verify"`
	Header        map[string][]string
	Method        string
	CheckRestart  *CheckRestart `mapstructure:"check_restart"`
}

// The Service model represents a Consul service definition
type Service struct {
	Id           string
	Name         string
	Tags         []string
	PortLabel    string `mapstructure:"port"`
	AddressMode  string `mapstructure:"address_mode"`
	Checks       []ServiceCheck
	CheckRestart *CheckRestart `mapstructure:"check_restart"`
}

func (s *Service) Canonicalize(t *Task, tg *TaskGroup, job *Job) {
	if s.Name == "" {
		s.Name = fmt.Sprintf("%s-%s-%s", *job.Name, *tg.Name, t.Name)
	}

	// Default to AddressModeAuto
	if s.AddressMode == "" {
		s.AddressMode = "auto"
	}

	// Canonicallize CheckRestart on Checks and merge Service.CheckRestart
	// into each check.
	for i, check := range s.Checks {
		s.Checks[i].CheckRestart = s.CheckRestart.Merge(check.CheckRestart)
		s.Checks[i].CheckRestart.Canonicalize()
	}
}

// EphemeralDisk is an ephemeral disk object
type EphemeralDisk struct {
	Sticky  *bool
	Migrate *bool
	SizeMB  *int `mapstructure:"size"`
}

func DefaultEphemeralDisk() *EphemeralDisk {
	return &EphemeralDisk{
		Sticky:  helper.BoolToPtr(false),
		Migrate: helper.BoolToPtr(false),
		SizeMB:  helper.IntToPtr(300),
	}
}

func (e *EphemeralDisk) Canonicalize() {
	if e.Sticky == nil {
		e.Sticky = helper.BoolToPtr(false)
	}
	if e.Migrate == nil {
		e.Migrate = helper.BoolToPtr(false)
	}
	if e.SizeMB == nil {
		e.SizeMB = helper.IntToPtr(300)
	}
}

// TaskGroup is the unit of scheduling.
type TaskGroup struct {
	Name          *string
	Count         *int
	Constraints   []*Constraint
	Tasks         []*Task
	RestartPolicy *RestartPolicy
	EphemeralDisk *EphemeralDisk
	Update        *UpdateStrategy
	Meta          map[string]string
}

// NewTaskGroup creates a new TaskGroup.
func NewTaskGroup(name string, count int) *TaskGroup {
	return &TaskGroup{
		Name:  helper.StringToPtr(name),
		Count: helper.IntToPtr(count),
	}
}

func (g *TaskGroup) Canonicalize(job *Job) {
	if g.Name == nil {
		g.Name = helper.StringToPtr("")
	}
	if g.Count == nil {
		g.Count = helper.IntToPtr(1)
	}
	for _, t := range g.Tasks {
		t.Canonicalize(g, job)
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

	var defaultRestartPolicy *RestartPolicy
	switch *job.Type {
	case "service", "system":
		defaultRestartPolicy = &RestartPolicy{
			Delay:    helper.TimeToPtr(15 * time.Second),
			Attempts: helper.IntToPtr(2),
			Interval: helper.TimeToPtr(1 * time.Minute),
			Mode:     helper.StringToPtr("delay"),
		}
	default:
		defaultRestartPolicy = &RestartPolicy{
			Delay:    helper.TimeToPtr(15 * time.Second),
			Attempts: helper.IntToPtr(15),
			Interval: helper.TimeToPtr(7 * 24 * time.Hour),
			Mode:     helper.StringToPtr("delay"),
		}
	}

	if g.RestartPolicy != nil {
		defaultRestartPolicy.Merge(g.RestartPolicy)
	}
	g.RestartPolicy = defaultRestartPolicy
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

// RequireDisk adds a ephemeral disk to the task group
func (g *TaskGroup) RequireDisk(disk *EphemeralDisk) *TaskGroup {
	g.EphemeralDisk = disk
	return g
}

// LogConfig provides configuration for log rotation
type LogConfig struct {
	MaxFiles      *int `mapstructure:"max_files"`
	MaxFileSizeMB *int `mapstructure:"max_file_size"`
}

func DefaultLogConfig() *LogConfig {
	return &LogConfig{
		MaxFiles:      helper.IntToPtr(10),
		MaxFileSizeMB: helper.IntToPtr(10),
	}
}

func (l *LogConfig) Canonicalize() {
	if l.MaxFiles == nil {
		l.MaxFiles = helper.IntToPtr(10)
	}
	if l.MaxFileSizeMB == nil {
		l.MaxFileSizeMB = helper.IntToPtr(10)
	}
}

// DispatchPayloadConfig configures how a task gets its input from a job dispatch
type DispatchPayloadConfig struct {
	File string
}

// Task is a single process in a task group.
type Task struct {
	Name            string
	Driver          string
	User            string
	Config          map[string]interface{}
	Constraints     []*Constraint
	Env             map[string]string
	Services        []*Service
	Resources       *Resources
	Meta            map[string]string
	KillTimeout     *time.Duration `mapstructure:"kill_timeout"`
	LogConfig       *LogConfig     `mapstructure:"logs"`
	Artifacts       []*TaskArtifact
	Vault           *Vault
	Templates       []*Template
	DispatchPayload *DispatchPayloadConfig
	Leader          bool
	ShutdownDelay   time.Duration `mapstructure:"shutdown_delay"`
	KillSignal      string        `mapstructure:"kill_signal"`
}

func (t *Task) Canonicalize(tg *TaskGroup, job *Job) {
	if t.Resources == nil {
		t.Resources = &Resources{}
	}
	t.Resources.Canonicalize()
	if t.KillTimeout == nil {
		t.KillTimeout = helper.TimeToPtr(5 * time.Second)
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
}

// TaskArtifact is used to download artifacts before running a task.
type TaskArtifact struct {
	GetterSource  *string           `mapstructure:"source"`
	GetterOptions map[string]string `mapstructure:"options"`
	GetterMode    *string           `mapstructure:"mode"`
	RelativeDest  *string           `mapstructure:"destination"`
}

func (a *TaskArtifact) Canonicalize() {
	if a.GetterMode == nil {
		a.GetterMode = helper.StringToPtr("any")
	}
	if a.GetterSource == nil {
		// Shouldn't be possible, but we don't want to panic
		a.GetterSource = helper.StringToPtr("")
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
			a.RelativeDest = helper.StringToPtr("local/")
		}
	}
}

type Template struct {
	SourcePath   *string        `mapstructure:"source"`
	DestPath     *string        `mapstructure:"destination"`
	EmbeddedTmpl *string        `mapstructure:"data"`
	ChangeMode   *string        `mapstructure:"change_mode"`
	ChangeSignal *string        `mapstructure:"change_signal"`
	Splay        *time.Duration `mapstructure:"splay"`
	Perms        *string        `mapstructure:"perms"`
	LeftDelim    *string        `mapstructure:"left_delimiter"`
	RightDelim   *string        `mapstructure:"right_delimiter"`
	Envvars      *bool          `mapstructure:"env"`
	VaultGrace   *time.Duration `mapstructure:"vault_grace"`
}

func (tmpl *Template) Canonicalize() {
	if tmpl.SourcePath == nil {
		tmpl.SourcePath = helper.StringToPtr("")
	}
	if tmpl.DestPath == nil {
		tmpl.DestPath = helper.StringToPtr("")
	}
	if tmpl.EmbeddedTmpl == nil {
		tmpl.EmbeddedTmpl = helper.StringToPtr("")
	}
	if tmpl.ChangeMode == nil {
		tmpl.ChangeMode = helper.StringToPtr("restart")
	}
	if tmpl.ChangeSignal == nil {
		if *tmpl.ChangeMode == "signal" {
			tmpl.ChangeSignal = helper.StringToPtr("SIGHUP")
		} else {
			tmpl.ChangeSignal = helper.StringToPtr("")
		}
	} else {
		sig := *tmpl.ChangeSignal
		tmpl.ChangeSignal = helper.StringToPtr(strings.ToUpper(sig))
	}
	if tmpl.Splay == nil {
		tmpl.Splay = helper.TimeToPtr(5 * time.Second)
	}
	if tmpl.Perms == nil {
		tmpl.Perms = helper.StringToPtr("0644")
	}
	if tmpl.LeftDelim == nil {
		tmpl.LeftDelim = helper.StringToPtr("{{")
	}
	if tmpl.RightDelim == nil {
		tmpl.RightDelim = helper.StringToPtr("}}")
	}
	if tmpl.Envvars == nil {
		tmpl.Envvars = helper.BoolToPtr(false)
	}
	if tmpl.VaultGrace == nil {
		tmpl.VaultGrace = helper.TimeToPtr(15 * time.Second)
	}
}

type Vault struct {
	Policies     []string
	Env          *bool
	ChangeMode   *string `mapstructure:"change_mode"`
	ChangeSignal *string `mapstructure:"change_signal"`
}

func (v *Vault) Canonicalize() {
	if v.Env == nil {
		v.Env = helper.BoolToPtr(true)
	}
	if v.ChangeMode == nil {
		v.ChangeMode = helper.StringToPtr("restart")
	}
	if v.ChangeSignal == nil {
		v.ChangeSignal = helper.StringToPtr("SIGHUP")
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
	// DEPRECATION NOTICE: The following fields are all deprecated. see TaskEvent struct in structs.go for details.
	FailsTask        bool
	RestartReason    string
	SetupError       string
	DriverError      string
	DriverMessage    string
	ExitCode         int
	Signal           int
	Message          string
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
