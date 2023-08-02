package vault

import (
	"context"
	"fmt"
	"path"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/vault/helper/fairshare"
	"github.com/hashicorp/vault/helper/locking"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// rotationSubPath is the sub-path used for the rotation manager view. This
	// is nested under the system view.
	rotationSubPath = "rotate/"

	// rotationViewPrefix is the prefix used for the ID based lookup of rotations.
	rotationViewPrefix = "rotation-id/"
)

// RotationManager is used by the Core to manage rotations. When a rotation
// period expires the RotationManager will handle doing automatic rotations.
type RotationManager struct {
	core   *Core
	router *Router
	idView *BarrierView
	logger log.Logger

	coreStateLock locking.RWMutex

	rotateFunc RotationStrategy
	jobManager *fairshare.JobManager
}

type RotationStrategy func(context.Context, *RotationManager, string, *namespace.Namespace)

// rotationJob should only be created through newRotationJob()
type rotationJob struct {
	rotationID string
	ns         *namespace.Namespace
	mgr        *RotationManager
	nsCtx      context.Context
	startTime  time.Time
}

func newRotationJob(nsCtx context.Context, rotationID string, ns *namespace.Namespace, r *RotationManager) (*rotationJob, error) {
	if rotationID == "" {
		return nil, fmt.Errorf("cannot have empty lease id")
	}
	if r == nil {
		return nil, fmt.Errorf("cannot have nil rotation manager")
	}
	if nsCtx == nil {
		return nil, fmt.Errorf("cannot have nil namespace context.Context")
	}

	return &rotationJob{
		rotationID: rotationID,
		ns:         ns,
		mgr:        r,
		nsCtx:      nsCtx,
		startTime:  time.Now(),
	}, nil
}

// NewRotationManager creates a new RotationManager that is backed
// using a given view, and uses the provided router for rotation.
func NewRotationManager(c *Core, view *BarrierView, r RotationStrategy, logger log.Logger) *RotationManager {
	managerLogger := logger.Named("job-manager")
	jobManager := fairshare.NewJobManager("rotate", getNumRotationWorkers(c, logger), managerLogger, c.metricSink)
	jobManager.Start()

	c.AddLogger(managerLogger)

	rot := &RotationManager{
		core:       c,
		router:     c.router,
		idView:     view.SubView(leaseViewPrefix),
		logger:     logger,
		rotateFunc: r,
		jobManager: jobManager,
	}

	if rot.logger == nil {
		opts := log.LoggerOptions{Name: "rotation_manager"}
		rot.logger = log.New(&opts)
	}

	return rot
}

// setupRotation is invoked after we've loaded the mount table to
// initialize the rotation manager
func (c *Core) setupRotation(r RotationStrategy) error {
	// Create a sub-view
	view := c.systemBarrierView.SubView(rotationSubPath)

	// Create the manager
	logger := c.baseLogger.Named("rotation")
	c.AddLogger(logger)
	mgr := NewRotationManager(c, view, r, logger)
	c.rotation = mgr

	// TODO: Restore the existing rotation state?

	c.logger.Debug("rotation setup completed")
	return nil
}

func (r *rotationJob) Execute() error {
	r.mgr.logger.Error("rotationJob.Execute called")

	// don't start the timer until the revocation is being executed
	revokeCtx, cancel := context.WithTimeout(r.nsCtx, DefaultMaxRequestDuration)
	defer cancel()

	r.mgr.coreStateLock.RLock()
	err := r.mgr.Rotate(revokeCtx, r.rotationID)
	r.mgr.coreStateLock.RUnlock()

	return err
}

func (r *rotationJob) OnFailure(err error) {
	r.mgr.logger.Error("rotationJob.OnFailure called")
	// TODO: handle failures?
}

// Rotate is used to rotate a credential named by the given rotationID
func (r *RotationManager) Rotate(ctx context.Context, rotationID string) error {
	r.logger.Debug("RotationManager.Rotate called")
	// TODO: rotate

	return nil
}

// Register takes a request and response with an associated StaticSecret. The
// secret gets assigned a RotationID and the management of the rotation is
// assumed by the rotation manager.
func (r *RotationManager) Register(ctx context.Context, req *logical.Request, resp *logical.Response) (id string, retErr error) {
	// Ignore if there is no static secret
	if resp == nil || resp.StaticSecret == nil {
		return "", nil
	}

	// TODO: Validate the secret?

	// Create a rotation entry. We use TokenLength because that is what is used
	// by ExpirationManager
	rotationRand, err := base62.Random(TokenLength)
	if err != nil {
		return "", err
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return "", err
	}

	rotationID := path.Join(req.Path, rotationRand)

	if ns.ID != namespace.RootNamespaceID {
		rotationID = fmt.Sprintf("%s.%s", rotationID, ns.ID)
	}

	re := &rotationEntry{
		RotationID: rotationID,
		Path:       req.Path,
		Data:       resp.Data,
		Secret:     resp.StaticSecret,
		IssueTime:  time.Now(),
		// TODO: handle cron schedule
		ExpireTime: time.Now().Add(*resp.StaticSecret.Period),
		namespace:  ns,
	}

	defer func() {
		if retErr != nil {
			// TODO: handle errors?
		}
	}()

	// TODO: lock and add the rotation entry?

	// We round here, see ExpirationManager's Register()
	period := re.ExpireTime.Sub(time.Now()).Round(time.Second)
	resp.StaticSecret.Period = &period

	return re.RotationID, nil
}

func rotateStrategyFairsharing(ctx context.Context, r *RotationManager, rotationID string, ns *namespace.Namespace) {
	nsCtx := namespace.ContextWithNamespace(ctx, ns)

	mountAccessor := r.getRotationAccessorLocked(ctx, rotationID)

	job, err := newRotationJob(nsCtx, rotationID, ns, r)
	if err != nil {
		r.logger.Warn("error creating rotation job", "error", err)
		return
	}

	r.jobManager.AddJob(job, mountAccessor)
	r.logger.Debug("jobManager added job", "mountAccessor", mountAccessor)
}

func (r *RotationManager) getRotationAccessorLocked(ctx context.Context, rotationID string) string {
	r.coreStateLock.RLock()
	defer r.coreStateLock.RUnlock()
	return r.getRotationAccessor(ctx, rotationID)
}

// note: this function must be called with r.coreStateLock held for read
func (r *RotationManager) getRotationAccessor(ctx context.Context, rotationID string) string {
	mount := r.core.router.MatchingMountEntry(ctx, rotationID)

	var mountAccessor string
	if mount == nil {
		mountAccessor = "mount-accessor-not-found"
	} else {
		mountAccessor = mount.Accessor
	}

	return mountAccessor
}

func getNumRotationWorkers(_ *Core, _ log.Logger) int {
	// TODO: make this configurable? See getNumExpirationWorkers
	numWorkers := 200
	return numWorkers
}

// rotationEntry is used to structure the values the expiration
// manager stores. This is used to handle renew and revocation.
type rotationEntry struct {
	RotationID string                 `json:"rotation_id"`
	Path       string                 `json:"path"`
	Data       map[string]interface{} `json:"data"`
	Secret     *logical.StaticSecret  `json:"static_secret"`
	IssueTime  time.Time              `json:"issue_time"`
	ExpireTime time.Time              `json:"expire_time"`

	namespace *namespace.Namespace
}
