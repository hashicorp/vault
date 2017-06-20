package api

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrRenewerMissingInput  = errors.New("missing input to renewer")
	ErrRenewerMissingSecret = errors.New("missing secret to renew")
	ErrRenewerNotRenewable  = errors.New("secret is not renewable")
	ErrRenewerNoSecretData  = errors.New("returned empty secret data")

	// DefaultRenewerGrace is the default grace period
	DefaultRenewerGrace = 15 * time.Second
)

// Renewer is a process for renewing a secret.
//
// 	renewer, err := client.NewRenewer(&RenewerInput{
// 		Secret: mySecret,
// 	})
// 	go renewer.Renew()
// 	defer renewer.Stop()
//
// 	for {
// 		select {
// 		case err := <-DoneCh():
// 			if err != nil {
// 				log.Fatal(err)
// 			}
//
// 			// Renewal is now over
// 		case renewal := <-RenewCh():
// 			log.Printf("Successfully renewed: %#v", renewal)
// 		default:
// 		}
// 	}
//
//
// The `DoneCh` will return if renewal fails or if the remaining lease duration
// after a renewal is less than or equal to the grace (in number of seconds). In
// both cases, the caller should attempt a re-read of the secret. Clients should
// check the return value of the channel to see if renewal was successful.
type Renewer struct {
	sync.Mutex

	client  *Client
	secret  *Secret
	grace   time.Duration
	doneCh  chan error
	renewCh chan *RenewOutput

	stopped bool
	stopCh  chan struct{}
}

// RenewerInput is used as input to the renew function.
type RenewerInput struct {
	// Secret is the secret to renew
	Secret *Secret

	// Grace is a minimum renewal (in seconds) before returring so the upstream
	// client can do a re-read. This can be used to prevent clients from waiting
	// too long to read a new credential and incur downtime.
	Grace time.Duration
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

// NewRenewer creates a new renewer from the given input.
func (c *Client) NewRenewer(i *RenewerInput) (*Renewer, error) {
	if i == nil {
		return nil, ErrRenewerMissingInput
	}

	secret := i.Secret
	if secret == nil {
		return nil, ErrRenewerMissingSecret
	}

	grace := i.Grace
	if grace == 0 {
		grace = DefaultRenewerGrace
	}

	return &Renewer{
		client:  c,
		secret:  secret,
		grace:   grace,
		doneCh:  make(chan error),
		renewCh: make(chan *RenewOutput, 5),

		stopped: false,
		stopCh:  make(chan struct{}),
	}, nil
}

// DoneCh returns the channel where the renewer will publish when renewal stops.
// If there is an error, this will be an error.
func (r *Renewer) DoneCh() <-chan error {
	return r.doneCh
}

// RenewCh is a channel that receives a message when a successful renewal takes
// place and includes metadata about the renewal.
func (r *Renewer) RenewCh() <-chan *RenewOutput {
	return r.renewCh
}

// Stop stops the renewer.
func (r *Renewer) Stop() {
	r.Lock()
	if !r.stopped {
		close(r.stopCh)
		r.stopped = true
	}
	r.Unlock()
}

// Renew starts a background process for renewing this secret. When the secret
// is has auth data, this attempts to renew the auth (token). When the secret
// has a lease, this attempts to renew the lease.
//
// This function will not return if nothing is reading from doneCh (it blocks)
// on a write to the channel.
func (r *Renewer) Renew() {
	if r.secret.Auth != nil {
		r.doneCh <- r.renewAuth()
	} else {
		r.doneCh <- r.renewLease()
	}
}

// renewAuth is a helper for renewing authentication.
func (r *Renewer) renewAuth() error {
	if !r.secret.Auth.Renewable || r.secret.Auth.ClientToken == "" {
		return ErrRenewerNotRenewable
	}

	client, token := r.client, r.secret.Auth.ClientToken

	for {
		// Check if we are stopped.
		select {
		case <-r.stopCh:
			return nil
		default:
		}

		// Renew the auth.
		renewal, err := client.Auth().Token().RenewTokenAsSelf(token, 0)
		if err != nil {
			return err
		}

		// Push a message that a renewal took place.
		select {
		case r.renewCh <- &RenewOutput{time.Now().UTC(), renewal}:
		default:
		}

		// Somehow, sometimes, this happens.
		if renewal == nil || renewal.Auth == nil {
			return ErrRenewerNoSecretData
		}

		// Do nothing if we are not renewable
		if !renewal.Auth.Renewable {
			return ErrRenewerNotRenewable
		}

		// Grab the lease duration - note that we grab the auth lease duration, not
		// the secret lease duration.
		leaseDuration := time.Duration(renewal.Auth.LeaseDuration) * time.Second

		// If we are within grace, return now.
		if leaseDuration <= r.grace {
			return nil
		}

		select {
		case <-r.stopCh:
			return nil
		case <-time.After(time.Duration(leaseDuration/2.0) * time.Second):
			continue
		}
	}
}

// renewLease is a helper for renewing a lease.
func (r *Renewer) renewLease() error {
	if !r.secret.Renewable || r.secret.LeaseID == "" {
		return ErrRenewerNotRenewable
	}

	client, leaseID := r.client, r.secret.LeaseID

	for {
		// Check if we are stopped.
		select {
		case <-r.stopCh:
			return nil
		default:
		}

		// Renew the lease.
		renewal, err := client.Sys().Renew(leaseID, 0)
		if err != nil {
			return err
		}

		// Push a message that a renewal took place.
		select {
		case r.renewCh <- &RenewOutput{time.Now().UTC(), renewal}:
		default:
		}

		// Somehow, sometimes, this happens.
		if renewal == nil {
			return ErrRenewerNoSecretData
		}

		// Do nothing if we are not renewable
		if !renewal.Renewable {
			return ErrRenewerNotRenewable
		}

		// Grab the lease duration
		leaseDuration := time.Duration(renewal.LeaseDuration) * time.Second

		// If we are within grace, return now.
		if leaseDuration <= r.grace {
			return nil
		}

		select {
		case <-r.stopCh:
			return nil
		case <-time.After(time.Duration(leaseDuration/2.0) * time.Second):
			continue
		}
	}
}
