package pki

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

var MaxChallengeTimeout = 1 * time.Minute

const MaxRetryAttempts = 5

type ChallengeValidation struct {
	// Account KID that this validation attempt is recorded under.
	Account string `json:"account"`

	// The authorization ID that this validation attempt is for.
	Authorization string            `json:"authorization"`
	ChallengeType ACMEChallengeType `json:"challenge_type"`

	// The token of this challenge and the JWS thumbprint of the account
	// we're validating against.
	Token      string `json:"token"`
	Thumbprint string `json:"thumbprint"`

	Initiated       time.Time `json:"initiated"`
	FirstValidation time.Time `json:"first_validation,omitempty"`
	RetryCount      int       `json:"retry_count,omitempty"`
	LastRetry       time.Time `json:"last_retry,omitempty"`
	RetryAfter      time.Time `json:"retry_after,omitempty"`
}

type ACMEChallengeEngine struct {
	NumWorkers int

	ValidationLock sync.Mutex
	NewValidation  chan string
	Closing        chan struct{}
	Validations    *list.List
}

func NewACMEChallengeEngine() *ACMEChallengeEngine {
	ace := &ACMEChallengeEngine{}
	ace.NewValidation = make(chan string, 1)
	ace.Closing = make(chan struct{}, 1)
	ace.Validations = list.New()
	ace.NumWorkers = 5

	return ace
}

func (ace *ACMEChallengeEngine) Initialize(b *backend, sc *storageContext) error {
	if err := ace.LoadFromStorage(b, sc); err != nil {
		return fmt.Errorf("failed loading initial in-progress validations: %w", err)
	}

	return nil
}

func (ace *ACMEChallengeEngine) LoadFromStorage(b *backend, sc *storageContext) error {
	items, err := sc.Storage.List(sc.Context, acmeValidationPrefix)
	if err != nil {
		return fmt.Errorf("failed loading list of validations from disk: %w", err)
	}

	ace.ValidationLock.Lock()
	defer ace.ValidationLock.Unlock()

	// Add them to our queue of validations to work through later.
	for _, item := range items {
		ace.Validations.PushBack(item)
	}

	return nil
}

func (ace *ACMEChallengeEngine) Run(b *backend) {
	for true {
		// err == nil on shutdown.
		b.Logger().Debug("Starting ACME challenge validation engine")
		err := ace._run(b)
		if err != nil {
			b.Logger().Error("Got unexpected error from ACME challenge validation engine", "err", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
}

func (ace *ACMEChallengeEngine) _run(b *backend) error {
	// This runner uses a background context for storage operations: we don't
	// want to tie it to a inbound request and we don't want to set a time
	// limit, so create a fresh background context.
	runnerSC := b.makeStorageContext(context.Background(), b.storage)

	// We want at most a certain number of workers operating to verify
	// challenges.
	var finishedWorkersChannels []chan bool
	for true {
		// Wait until we've got more work to do.
		select {
		case <-ace.Closing:
			b.Logger().Debug("shutting down ACME challenge validation engine")
			return nil
		case <-ace.NewValidation:
		}

		// First try to reap any finished workers. Read from their channels
		// and if not finished yet, add to a fresh slice.
		var newFinishedWorkersChannels []chan bool
		for _, channel := range finishedWorkersChannels {
			select {
			case <-channel:
			default:
				// This channel had not been written to, indicating that the
				// worker had not yet finished.
				newFinishedWorkersChannels = append(newFinishedWorkersChannels, channel)
			}
		}
		finishedWorkersChannels = newFinishedWorkersChannels

		// If we have space to take on another work item, do so.
		if len(finishedWorkersChannels) < ace.NumWorkers {
			ace.ValidationLock.Lock()
			element := ace.Validations.Front()
			if element != nil {
				ace.Validations.Remove(element)
			}
			ace.ValidationLock.Unlock()

			task := element.Value.(string)
			channel := make(chan bool, 1)
			go ace.VerifyChallenge(runnerSC, task, channel)
		}

		// If we have no more work.
		if len(finishedWorkersChannels) == ace.NumWorkers {
			time.Sleep(50 * time.Millisecond)
		}

		// Lastly, if we have more work to do, re-trigger ourselves.
		ace.ValidationLock.Lock()
		if ace.Validations.Front() != nil {
			select {
			case ace.NewValidation <- "retry":
			default:
			}
		}
		ace.ValidationLock.Unlock()
	}

	return fmt.Errorf("unexpectedly exited from ACMEChallengeEngine._run()")
}

func (ace *ACMEChallengeEngine) AcceptChallenge(sc *storageContext, account string, authz *ACMEAuthorization, challenge *ACMEChallenge, thumbprint string) error {
	name := authz.Id + "-" + string(challenge.Type)
	path := acmeValidationPrefix + name

	entry, err := sc.Storage.Get(sc.Context, path)
	if err == nil && entry != nil {
		// Challenge already in the queue; exit without re-adding it.
		return nil
	}

	if authz.Status != ACMEAuthorizationPending {
		return fmt.Errorf("cannot accept already validated authorization %v (%v)", authz.Id, authz.Status)
	}

	if challenge.Status != ACMEChallengePending && challenge.Status != ACMEChallengeProcessing {
		return fmt.Errorf("challenge is in invalid state (%v) in authorization %v", challenge.Status, authz.Id)
	}

	token := challenge.ChallengeFields["token"].(string)

	cv := &ChallengeValidation{
		Account:       account,
		Authorization: authz.Id,
		ChallengeType: challenge.Type,
		Token:         token,
		Thumbprint:    thumbprint,
		Initiated:     time.Now(),
	}

	json, err := logical.StorageEntryJSON(path, &cv)
	if err != nil {
		return fmt.Errorf("error creating challenge validation queue entry: %w", err)
	}

	if err := sc.Storage.Put(sc.Context, json); err != nil {
		return fmt.Errorf("error writing challenge validation entry: %w", err)
	}

	if challenge.Status == ACMEChallengePending {
		challenge.Status = ACMEChallengeProcessing

		authzPath := getAuthorizationPath(account, authz.Id)
		if err := saveAuthorizationAtPath(sc, authzPath, authz); err != nil {
			return fmt.Errorf("error saving updated authorization %v: %w", authz.Id, err)
		}
	}

	ace.ValidationLock.Lock()
	defer ace.ValidationLock.Unlock()
	ace.Validations.PushBack(name)

	select {
	case ace.NewValidation <- name:
	default:
	}

	return nil
}

func (ace *ACMEChallengeEngine) VerifyChallenge(runnerSc *storageContext, id string, finished chan bool) {
	sc, _ /* cancel func */ := runnerSc.WithFreshTimeout(MaxChallengeTimeout)
	runnerSc.Backend.Logger().Debug("Starting verification of challenge: %v", id)

	if retry, err := ace._verifyChallenge(sc, id); err != nil {
		// Because verification of this challenge failed, we need to retry
		// it in the future. Log the error and re-add the item to the queue
		// to try again later.
		sc.Backend.Logger().Error(fmt.Sprintf("ACME validation failed for %v: %v", id, err))

		if retry {
			ace.ValidationLock.Lock()
			defer ace.ValidationLock.Unlock()
			ace.Validations.PushBack(id)

			// Let the validator know there's a pending challenge.
			select {
			case ace.NewValidation <- id:
			default:
			}
		}

		// We're the only producer on this channel and it has a buffer size
		// of one element, so it is safe to directly write here.
		finished <- true
		return
	}

	// We're the only producer on this channel and it has a buffer size of one
	// element, so it is safe to directly write here.
	finished <- false
}

func (ace *ACMEChallengeEngine) _verifyChallenge(sc *storageContext, id string) (bool, error) {
	now := time.Now()
	path := acmeValidationPrefix + id
	challengeEntry, err := sc.Storage.Get(sc.Context, path)
	if err != nil {
		return true, fmt.Errorf("error loading challenge %v: %w", id, err)
	}

	if challengeEntry == nil {
		// Something must've successfully cleaned up our storage entry from
		// under us. Assume we don't need to rerun, else the client will
		// trigger us to re-run.
		err = nil
		return ace._verifyChallengeCleanup(sc, err, id)
	}

	var cv *ChallengeValidation
	if err := challengeEntry.DecodeJSON(&cv); err != nil {
		return true, fmt.Errorf("error decoding challenge %v: %w", id, err)
	}

	if now.Before(cv.RetryAfter) {
		time.Sleep(50 * time.Millisecond)
		return true, fmt.Errorf("retrying challenge %v too soon", id)
	}

	authzPath := getAuthorizationPath(cv.Account, cv.Authorization)
	authz, err := loadAuthorizationAtPath(sc, authzPath)
	if err != nil {
		return true, fmt.Errorf("error loading authorization %v/%v for challenge %v: %w", cv.Account, cv.Authorization, id, err)
	}

	if authz.Status != ACMEAuthorizationPending {
		// Something must've finished up this challenge for us. Assume we
		// don't need to rerun and exit instead.
		err = nil
		return ace._verifyChallengeCleanup(sc, err, id)
	}

	var challenge *ACMEChallenge
	for _, authzChallenge := range authz.Challenges {
		if authzChallenge.Type == cv.ChallengeType {
			challenge = authzChallenge
			break
		}
	}

	if challenge == nil {
		err = fmt.Errorf("no challenge of type %v in authorization %v/%v for challenge %v", cv.ChallengeType, cv.Account, cv.Authorization, id)
		return ace._verifyChallengeCleanup(sc, err, id)
	}

	if challenge.Status != ACMEChallengePending && challenge.Status != ACMEChallengeProcessing {
		err = fmt.Errorf("challenge is in invalid state %v in authorization %v/%v for challenge %v", challenge.Status, cv.Account, cv.Authorization, id)
		return ace._verifyChallengeCleanup(sc, err, id)
	}

	var valid bool
	switch challenge.Type {
	case ACMEHTTPChallenge:
		if authz.Identifier.Type != ACMEDNSIdentifier && authz.Identifier.Type != ACMEIPIdentifier {
			err = fmt.Errorf("unsupported identifier type for authorization %v/%v in challenge %v: %v", cv.Account, cv.Authorization, id, authz.Identifier.Type)
			return ace._verifyChallengeCleanup(sc, err, id)
		}

		if authz.Wildcard {
			err = fmt.Errorf("unable to validate wildcard authorization %v/%v in challenge %v via http-01 challenge", cv.Account, cv.Authorization, id)
			return ace._verifyChallengeCleanup(sc, err, id)
		}

		valid, err = ValidateHTTP01Challenge(authz.Identifier.Value, cv.Token, cv.Thumbprint)
		if err != nil {
			err = fmt.Errorf("error validating http-01 challenge %v: %w", id, err)
			return ace._verifyChallengeRetry(sc, cv, authz, err, id)
		}
	case ACMEDNSChallenge:
		if authz.Identifier.Type != ACMEDNSIdentifier {
			err = fmt.Errorf("unsupported identifier type for authorization %v/%v in challenge %v: %v", cv.Account, cv.Authorization, id, authz.Identifier.Type)
            return ace._verifyChallengeCleanup(sc, err, id)
		}

		valid, err = ValidateDNS01Challenge(authz.Identifier.Value, cv.Token, cv.Thumbprint)
        if err != nil {
            err = fmt.Errorf("error validating dns-01 challenge %v: %w", id, err)
            return ace._verifyChallengeRetry(sc, cv, authz, err, id)
        }
	default:
		err = fmt.Errorf("unsupported ACME challenge type %v for challenge %v", cv.ChallengeType, id)
		return ace._verifyChallengeCleanup(sc, err, id)
	}

	if !valid {
		err = fmt.Errorf("challenge failed with no additional information")
		return ace._verifyChallengeRetry(sc, cv, authz, err, id)
	}

	// If we got here, the challenge verification was successful. Update
	// the authorization appropriately.
	expires := now.Add(15 * 24 * time.Hour)
	challenge.Status = ACMEChallengeValid
	challenge.Validated = now.Format(time.RFC3339)
	authz.Status = ACMEAuthorizationValid
	authz.Expires = expires.Format(time.RFC3339)

	if err := saveAuthorizationAtPath(sc, authzPath, authz); err != nil {
		err = fmt.Errorf("error saving updated (validated) authorization %v/%v for challenge %v: %w", cv.Account, cv.Authorization, id, err)
		return ace._verifyChallengeRetry(sc, cv, authz, err, id)
	}

	return ace._verifyChallengeCleanup(sc, nil, id)
}

func (ace *ACMEChallengeEngine) _verifyChallengeRetry(sc *storageContext, cv *ChallengeValidation, authz *ACMEAuthorization, err error, id string) (bool, error) {
	now := time.Now()
	path := acmeValidationPrefix + id

	if cv.RetryCount > MaxRetryAttempts {
		err = fmt.Errorf("reached max error attempts for challenge %v: %w", id, err)
		return ace._verifyChallengeCleanup(sc, err, id)
	}

	if cv.FirstValidation.IsZero() {
		cv.FirstValidation = now
	}
	cv.RetryCount += 1
	cv.LastRetry = now
	cv.RetryAfter = now.Add(time.Duration(cv.RetryCount*5) * time.Second)

	json, jsonErr := logical.StorageEntryJSON(path, cv)
	if jsonErr != nil {
		return true, fmt.Errorf("error persisting updated challenge validation queue entry (error prior to retry, if any: %v): %w", err, jsonErr)
	}

	if putErr := sc.Storage.Put(sc.Context, json); putErr != nil {
		return true, fmt.Errorf("error writing updated challenge validation entry (error prior to retry, if any: %v): %w", err, putErr)
	}

	if err != nil {
		err = fmt.Errorf("retrying validation: %w", err)
	}

	return true, err
}

func (ace *ACMEChallengeEngine) _verifyChallengeCleanup(sc *storageContext, err error, id string) (bool, error) {
	// Remove our ChallengeValidation entry only.
	if deleteErr := sc.Storage.Delete(sc.Context, acmeValidationPrefix+id); deleteErr != nil {
		return true, fmt.Errorf("error deleting challenge %v (error prior to cleanup, if any: %v): %w", id, err, deleteErr)
	}

	if err != nil {
		err = fmt.Errorf("removing challenge validation attempt and not retrying %v; previous error: %w", id, err)
	}

	return false, err
}
