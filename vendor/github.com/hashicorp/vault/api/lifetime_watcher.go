package api

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrLifetimeWatcherMissingInput  = errors.New("missing input")
	ErrLifetimeWatcherMissingSecret = errors.New("missing secret")
	ErrLifetimeWatcherNotRenewable  = errors.New("secret is not renewable")
	ErrLifetimeWatcherNoSecretData  = errors.New("returned empty secret data")

	// Deprecated; kept for compatibility
	ErrRenewerMissingInput  = errors.New("missing input to renewer")
	ErrRenewerMissingSecret = errors.New("missing secret to renew")
	ErrRenewerNotRenewable  = errors.New("secret is not renewable")
	ErrRenewerNoSecretData  = errors.New("returned empty secret data")

	// DefaultLifetimeWatcherRenewBuffer is the default size of the buffer for renew
	// messages on the channel.
	DefaultLifetimeWatcherRenewBuffer = 5
	// Deprecated: kept for backwards compatibility
	DefaultRenewerRenewBuffer = 5
)

type RenewBehavior uint

const (
	// RenewBehaviorIgnoreErrors means we will attempt to keep renewing until
	// we hit the lifetime threshold. It also ignores errors stemming from
	// passing a non-renewable lease in. In practice, this means you simply
	// reauthenticate/refetch credentials when the watcher exits. This is the
	// default.
	RenewBehaviorIgnoreErrors RenewBehavior = iota

	// RenewBehaviorRenewDisabled turns off renewal attempts entirely. This
	// allows you to simply watch lifetime and have the watcher return at a
	// reasonable threshold without actually making Vault calls.
	RenewBehaviorRenewDisabled

	// RenewBehaviorErrorOnErrors is the "legacy" behavior which always exits
	// on some kind of error
	RenewBehaviorErrorOnErrors
)

// LifetimeWatcher is a process for watching lifetime of a secret.
//
// 	watcher, err := client.NewLifetimeWatcher(&LifetimeWatcherInput{
// 		Secret: mySecret,
// 	})
// 	go watcher.Start()
// 	defer watcher.Stop()
//
// 	for {
// 		select {
// 		case err := <-watcher.DoneCh():
// 			if err != nil {
// 				log.Fatal(err)
// 			}
//
// 			// Renewal is now over
// 		case renewal := <-watcher.RenewCh():
// 			log.Printf("Successfully renewed: %#v", renewal)
// 		}
// 	}
//
//
// `DoneCh` will return if renewal fails, or if the remaining lease duration is
// under a built-in threshold and either renewing is not extending it or
// renewing is disabled.  In both cases, the caller should attempt a re-read of
// the secret. Clients should check the return value of the channel to see if
// renewal was successful.
type LifetimeWatcher struct {
	l sync.Mutex

	client        *Client
	secret        *Secret
	grace         time.Duration
	random        *rand.Rand
	increment     int
	doneCh        chan error
	renewCh       chan *RenewOutput
	renewBehavior RenewBehavior

	stopped bool
	stopCh  chan struct{}

	errLifetimeWatcherNotRenewable error
	errLifetimeWatcherNoSecretData error
}

// LifetimeWatcherInput is used as input to the renew function.
type LifetimeWatcherInput struct {
	// Secret is the secret to renew
	Secret *Secret

	// DEPRECATED: this does not do anything.
	Grace time.Duration

	// Rand is the randomizer to use for underlying randomization. If not
	// provided, one will be generated and seeded automatically. If provided, it
	// is assumed to have already been seeded.
	Rand *rand.Rand

	// RenewBuffer is the size of the buffered channel where renew messages are
	// dispatched.
	RenewBuffer int

	// The new TTL, in seconds, that should be set on the lease. The TTL set
	// here may or may not be honored by the vault server, based on Vault
	// configuration or any associated max TTL values.
	Increment int

	// RenewBehavior controls what happens when a renewal errors or the
	// passed-in secret is not renewable.
	RenewBehavior RenewBehavior
}

// RenewOutput is the metadata returned to the client (if it's listening) to
// renew messages.
type RenewOutput struct {
	// RenewedAt is the timestamp when the renewal took place (UTC).
	RenewedAt time.Time

	// Secret is the underlying renewal data. It's the same struct as all data
	// that is returned from Vault, but since this is renewal data, it will not
	// usually include the secret itself.
	Secret *Secret
}

// NewLifetimeWatcher creates a new renewer from the given input.
func (c *Client) NewLifetimeWatcher(i *LifetimeWatcherInput) (*LifetimeWatcher, error) {
	if i == nil {
		return nil, ErrLifetimeWatcherMissingInput
	}

	secret := i.Secret
	if secret == nil {
		return nil, ErrLifetimeWatcherMissingSecret
	}

	random := i.Rand
	if random == nil {
		random = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	}

	renewBuffer := i.RenewBuffer
	if renewBuffer == 0 {
		renewBuffer = DefaultLifetimeWatcherRenewBuffer
	}

	return &LifetimeWatcher{
		client:        c,
		secret:        secret,
		increment:     i.Increment,
		random:        random,
		doneCh:        make(chan error, 1),
		renewCh:       make(chan *RenewOutput, renewBuffer),
		renewBehavior: i.RenewBehavior,

		stopped: false,
		stopCh:  make(chan struct{}),

		errLifetimeWatcherNotRenewable: ErrLifetimeWatcherNotRenewable,
		errLifetimeWatcherNoSecretData: ErrLifetimeWatcherNoSecretData,
	}, nil
}

// Deprecated: exists only for backwards compatibility. Calls
// NewLifetimeWatcher, and sets compatibility flags.
func (c *Client) NewRenewer(i *LifetimeWatcherInput) (*LifetimeWatcher, error) {
	if i == nil {
		return nil, ErrRenewerMissingInput
	}

	secret := i.Secret
	if secret == nil {
		return nil, ErrRenewerMissingSecret
	}

	renewer, err := c.NewLifetimeWatcher(i)
	if err != nil {
		return nil, err
	}

	renewer.renewBehavior = RenewBehaviorErrorOnErrors
	renewer.errLifetimeWatcherNotRenewable = ErrRenewerNotRenewable
	renewer.errLifetimeWatcherNoSecretData = ErrRenewerNoSecretData
	return renewer, err
}

// DoneCh returns the channel where the renewer will publish when renewal stops.
// If there is an error, this will be an error.
func (r *LifetimeWatcher) DoneCh() <-chan error {
	return r.doneCh
}

// RenewCh is a channel that receives a message when a successful renewal takes
// place and includes metadata about the renewal.
func (r *LifetimeWatcher) RenewCh() <-chan *RenewOutput {
	return r.renewCh
}

// Stop stops the renewer.
func (r *LifetimeWatcher) Stop() {
	r.l.Lock()
	defer r.l.Unlock()

	if !r.stopped {
		close(r.stopCh)
		r.stopped = true
	}
}

// Start starts a background process for watching the lifetime of this secret.
// If renewal is enabled, when the secret has auth data, this attempts to renew
// the auth (token); When the secret has a lease, this attempts to renew the
// lease.
func (r *LifetimeWatcher) Start() {
	r.doneCh <- r.doRenew()
}

// Renew is for comnpatibility with the legacy api.Renewer. Calling Renew
// simply chains to Start.
func (r *LifetimeWatcher) Renew() {
	r.Start()
}

// renewAuth is a helper for renewing authentication.
func (r *LifetimeWatcher) doRenew() error {
	var nonRenewable bool
	var tokenMode bool
	var initLeaseDuration int
	var credString string
	var renewFunc func(string, int) (*Secret, error)

	switch {
	case r.secret.Auth != nil:
		tokenMode = true
		nonRenewable = !r.secret.Auth.Renewable
		initLeaseDuration = r.secret.Auth.LeaseDuration
		credString = r.secret.Auth.ClientToken
		renewFunc = r.client.Auth().Token().RenewTokenAsSelf
	default:
		nonRenewable = !r.secret.Renewable
		initLeaseDuration = r.secret.LeaseDuration
		credString = r.secret.LeaseID
		renewFunc = r.client.Sys().Renew
	}

	if credString == "" ||
		(nonRenewable && r.renewBehavior == RenewBehaviorErrorOnErrors) {
		return r.errLifetimeWatcherNotRenewable
	}

	initialTime := time.Now()
	priorDuration := time.Duration(initLeaseDuration) * time.Second
	r.calculateGrace(priorDuration)

	for {
		// Check if we are stopped.
		select {
		case <-r.stopCh:
			return nil
		default:
		}

		var leaseDuration time.Duration
		fallbackLeaseDuration := initialTime.Add(priorDuration).Sub(time.Now())

		switch {
		case nonRenewable || r.renewBehavior == RenewBehaviorRenewDisabled:
			// Can't or won't renew, just keep the same expiration so we exit
			// when it's reauthentication time
			leaseDuration = fallbackLeaseDuration

		default:
			// Renew the token
			renewal, err := renewFunc(credString, r.increment)
			if err != nil || renewal == nil || (tokenMode && renewal.Auth == nil) {
				if r.renewBehavior == RenewBehaviorErrorOnErrors {
					if err != nil {
						return err
					}
					if renewal == nil || (tokenMode && renewal.Auth == nil) {
						return r.errLifetimeWatcherNoSecretData
					}
				}

				leaseDuration = fallbackLeaseDuration
				break
			}

			// Push a message that a renewal took place.
			select {
			case r.renewCh <- &RenewOutput{time.Now().UTC(), renewal}:
			default:
			}

			// Possibly error if we are not renewable
			if ((tokenMode && !renewal.Auth.Renewable) || (!tokenMode && !renewal.Renewable)) &&
				r.renewBehavior == RenewBehaviorErrorOnErrors {
				return r.errLifetimeWatcherNotRenewable
			}

			// Grab the lease duration
			newDuration := renewal.LeaseDuration
			if tokenMode {
				newDuration = renewal.Auth.LeaseDuration
			}

			leaseDuration = time.Duration(newDuration) * time.Second
		}

		// We keep evaluating a new grace period so long as the lease is
		// extending. Once it stops extending, we've hit the max and need to
		// rely on the grace duration.
		if leaseDuration > priorDuration {
			r.calculateGrace(leaseDuration)
		}
		priorDuration = leaseDuration

		// The sleep duration is set to 2/3 of the current lease duration plus
		// 1/3 of the current grace period, which adds jitter.
		sleepDuration := time.Duration(float64(leaseDuration.Nanoseconds())*2/3 + float64(r.grace.Nanoseconds())/3)

		// If we are within grace, return now; or, if the amount of time we
		// would sleep would land us in the grace period. This helps with short
		// tokens; for example, you don't want a current lease duration of 4
		// seconds, a grace period of 3 seconds, and end up sleeping for more
		// than three of those seconds and having a very small budget of time
		// to renew.
		if leaseDuration <= r.grace || leaseDuration-sleepDuration <= r.grace {
			return nil
		}

		select {
		case <-r.stopCh:
			return nil
		case <-time.After(sleepDuration):
			continue
		}
	}
}

// sleepDuration calculates the time to sleep given the base lease duration. The
// base is the resulting lease duration. It will be reduced to 1/3 and
// multiplied by a random float between 0.0 and 1.0. This extra randomness
// prevents multiple clients from all trying to renew simultaneously.
func (r *LifetimeWatcher) sleepDuration(base time.Duration) time.Duration {
	sleep := float64(base)

	// Renew at 1/3 the remaining lease. This will give us an opportunity to retry
	// at least one more time should the first renewal fail.
	sleep = sleep / 3.0

	// Use a randomness so many clients do not hit Vault simultaneously.
	sleep = sleep * (r.random.Float64() + 1) / 2.0

	return time.Duration(sleep)
}

// calculateGrace calculates the grace period based on a reasonable set of
// assumptions given the total lease time; it also adds some jitter to not have
// clients be in sync.
func (r *LifetimeWatcher) calculateGrace(leaseDuration time.Duration) {
	if leaseDuration == 0 {
		r.grace = 0
		return
	}

	leaseNanos := float64(leaseDuration.Nanoseconds())
	jitterMax := 0.1 * leaseNanos

	// For a given lease duration, we want to allow 80-90% of that to elapse,
	// so the remaining amount is the grace period
	r.grace = time.Duration(jitterMax) + time.Duration(uint64(r.random.Int63())%uint64(jitterMax))
}

type Renewer = LifetimeWatcher
type RenewerInput = LifetimeWatcherInput
