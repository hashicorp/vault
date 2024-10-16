// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"container/list"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

var MaxChallengeTimeout = 1 * time.Minute

const MaxRetryAttempts = 5

const ChallengeAttemptFailedMsg = "this may occur if the validation target was misconfigured: check that challenge responses are available at the required locations and retry."

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

type ChallengeQueueEntry struct {
	Identifier string
	RetryAfter time.Time
	NumRetries int // Track if we are spinning on a corrupted challenge
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

func (ace *ACMEChallengeEngine) LoadFromStorage(b *backend, sc *storageContext) error {
	items, err := sc.Storage.List(sc.Context, acmeValidationPrefix)
	if err != nil {
		return fmt.Errorf("failed loading list of validations from disk: %w", err)
	}

	ace.ValidationLock.Lock()
	defer ace.ValidationLock.Unlock()

	// Add them to our queue of validations to work through later.
	foundExistingValidations := false
	for _, item := range items {
		ace.Validations.PushBack(&ChallengeQueueEntry{
			Identifier: item,
		})
		foundExistingValidations = true
	}

	if foundExistingValidations {
		ace.NewValidation <- "existing"
	}

	return nil
}

func (ace *ACMEChallengeEngine) Run(b *backend, state *acmeState, sc *storageContext) {
	// We load the existing ACME challenges within the Run thread to avoid
	// delaying the PKI mount initialization
	b.Logger().Debug("Loading existing challenge validations on disk")
	err := ace.LoadFromStorage(b, sc)
	if err != nil {
		b.Logger().Error("failed loading existing ACME challenge validations:", "err", err)
	}

	for {
		// err == nil on shutdown.
		b.Logger().Debug("Starting ACME challenge validation engine")
		err := ace._run(b, state)
		if err != nil {
			b.Logger().Error("Got unexpected error from ACME challenge validation engine", "err", err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
}

func (ace *ACMEChallengeEngine) _run(b *backend, state *acmeState) error {
	// This runner uses a background context for storage operations: we don't
	// want to tie it to a inbound request and we don't want to set a time
	// limit, so create a fresh background context.
	runnerSC := b.makeStorageContext(context.Background(), b.storage)

	// We want at most a certain number of workers operating to verify
	// challenges.
	var finishedWorkersChannels []chan bool
	for {
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
		firstIdentifier := ""
		startedWork := false
		now := time.Now()
		for len(finishedWorkersChannels) < ace.NumWorkers {
			var task *ChallengeQueueEntry

			// Find our next work item. We do all of these operations
			// while holding the queue lock, hence some repeated checks
			// afterwards. Out of this, we get a candidate task, using
			// element == nil as a sentinel for breaking our parent
			// loop.
			ace.ValidationLock.Lock()
			element := ace.Validations.Front()
			if element != nil {
				ace.Validations.Remove(element)
				task = element.Value.(*ChallengeQueueEntry)
				if !task.RetryAfter.IsZero() && now.Before(task.RetryAfter) {
					// We cannot work on this element yet; remove it to
					// the back of the queue. This allows us to potentially
					// select the next item in the next iteration.
					ace.Validations.PushBack(task)
				}

				if firstIdentifier != "" && task.Identifier == firstIdentifier {
					// We found and rejected this element before; exit the
					// loop by "claiming" we didn't find any work.
					element = nil
				} else if firstIdentifier == "" {
					firstIdentifier = task.Identifier
				}
			}
			ace.ValidationLock.Unlock()
			if element == nil {
				// There was no more work to do to fill up the queue; exit
				// this loop.
				break
			}
			if now.Before(task.RetryAfter) {
				// Here, while we found an element, we didn't want to
				// completely exit the loop (perhaps it was our first time
				// finding a work order), so retry without modifying
				// firstIdentifier.
				continue
			}

			config, err := state.getConfigWithUpdate(runnerSC)
			if err != nil {
				return fmt.Errorf("failed fetching ACME configuration: %w", err)
			}

			// Since this work item was valid, we won't expect to see it in
			// the validation queue again until it is executed. Here, we
			// want to avoid infinite looping above (if we removed the one
			// valid item and the remainder are all not immediately
			// actionable). At the worst, we'll spend a little more time
			// looping through the queue until we hit a repeat.
			firstIdentifier = ""

			// If we are no longer the active node, break out
			if b.System().ReplicationState().HasState(consts.ReplicationDRSecondary | consts.ReplicationPerformanceStandby) {
				break
			}

			// Here, we got a piece of work that is ready to check; create a
			// channel and a new go routine and run it. Note that this still
			// could have a RetryAfter date we're not aware of (e.g., if the
			// cluster restarted as we do not read the entries there).
			channel := make(chan bool, 1)
			go ace.VerifyChallenge(runnerSC, task.Identifier, task.NumRetries, channel, config)
			finishedWorkersChannels = append(finishedWorkersChannels, channel)
			startedWork = true
		}

		// If we have no more capacity for work, we should pause a little to
		// let the system catch up. Additionally, if we only had
		// non-actionable work items, we should pause until some time has
		// elapsed: not too much that we potentially starve any new incoming
		// items from validation, but not too short that we cause a busy loop.
		if len(finishedWorkersChannels) == ace.NumWorkers || !startedWork {
			time.Sleep(100 * time.Millisecond)
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
		return fmt.Errorf("%w: cannot accept already validated authorization %v (%v)", ErrMalformed, authz.Id, authz.Status)
	}

	for _, otherChallenge := range authz.Challenges {
		// We assume within an authorization we won't have multiple challenges of the same challenge type
		// and we want to limit a single challenge being in a processing state to avoid race conditions
		// failing one challenge and passing another.
		if otherChallenge.Type != challenge.Type && otherChallenge.Status != ACMEChallengePending {
			return fmt.Errorf("%w: only a single challenge within an authorization can be accepted (%v) in status %v", ErrMalformed, otherChallenge.Type, otherChallenge.Status)
		}

		// The requested challenge can ping us to wake us up, so allow pending and currently processing statuses
		if otherChallenge.Status != ACMEChallengePending && otherChallenge.Status != ACMEChallengeProcessing {
			return fmt.Errorf("%w: challenge is in invalid state (%v) in authorization %v", ErrMalformed, challenge.Status, authz.Id)
		}
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
	ace.Validations.PushBack(&ChallengeQueueEntry{
		Identifier: name,
	})

	select {
	case ace.NewValidation <- name:
	default:
	}

	return nil
}

func (ace *ACMEChallengeEngine) VerifyChallenge(runnerSc *storageContext, id string, validationQueueRetries int, finished chan bool, config *acmeConfigEntry) {
	sc, cancel := runnerSc.WithFreshTimeout(MaxChallengeTimeout)
	defer cancel()
	runnerSc.Logger().Debug("Starting verification of challenge", "id", id)

	if retry, retryAfter, err := ace._verifyChallenge(sc, id, config); err != nil {
		// Because verification of this challenge failed, we need to retry
		// it in the future. Log the error and re-add the item to the queue
		// to try again later.
		sc.Logger().Error(fmt.Sprintf("ACME validation failed for %v: %v", id, err))

		if retry {
			validationQueueRetries++

			// The retry logic within _verifyChallenge is dependent on being able to read and decode
			// the ACME challenge entries. If we encounter such failures we would retry forever, so
			// we have a secondary check here to see if we are consistently looping within the validation
			// queue that is larger than the normal retry attempts we would allow.
			if validationQueueRetries > MaxRetryAttempts*2 {
				sc.Logger().Warn("reached max error attempts within challenge queue: %v, giving up", id)
				_, _, err = ace._verifyChallengeCleanup(sc, nil, id)
				if err != nil {
					sc.Logger().Warn("Failed cleaning up challenge entry: %v", err)
				}
				finished <- true
				return
			}

			ace.ValidationLock.Lock()
			defer ace.ValidationLock.Unlock()
			ace.Validations.PushBack(&ChallengeQueueEntry{
				Identifier: id,
				RetryAfter: retryAfter,
				NumRetries: validationQueueRetries,
			})

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

func (ace *ACMEChallengeEngine) _verifyChallenge(sc *storageContext, id string, config *acmeConfigEntry) (bool, time.Time, error) {
	now := time.Now()
	backoffTime := now.Add(1 * time.Second)
	path := acmeValidationPrefix + id
	challengeEntry, err := sc.Storage.Get(sc.Context, path)
	if err != nil {
		return true, backoffTime, fmt.Errorf("error loading challenge %v: %w", id, err)
	}

	if challengeEntry == nil {
		// Something must've successfully cleaned up our storage entry from
		// under us. Assume we don't need to rerun, else the client will
		// trigger us to re-run.
		return ace._verifyChallengeCleanup(sc, nil, id)
	}

	var cv *ChallengeValidation
	if err := challengeEntry.DecodeJSON(&cv); err != nil {
		return true, backoffTime, fmt.Errorf("error decoding challenge %v: %w", id, err)
	}

	if now.Before(cv.RetryAfter) {
		return true, cv.RetryAfter, fmt.Errorf("retrying challenge %v too soon", id)
	}

	authzPath := getAuthorizationPath(cv.Account, cv.Authorization)
	authz, err := loadAuthorizationAtPath(sc, authzPath)
	if err != nil {
		return true, backoffTime, fmt.Errorf("error loading authorization %v/%v for challenge %v: %w", cv.Account, cv.Authorization, id, err)
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

		domain := encodeIdentifierForHTTP01Challenge(authz.Identifier)

		valid, err = ValidateHTTP01Challenge(domain, cv.Token, cv.Thumbprint, config)
		if err != nil {
			err = fmt.Errorf("%w: error validating http-01 challenge %v: %v; %v", ErrIncorrectResponse, id, err, ChallengeAttemptFailedMsg)
			return ace._verifyChallengeRetry(sc, cv, authzPath, authz, challenge, err, id)
		}
	case ACMEDNSChallenge:
		if authz.Identifier.Type != ACMEDNSIdentifier {
			err = fmt.Errorf("unsupported identifier type for authorization %v/%v in challenge %v: %v", cv.Account, cv.Authorization, id, authz.Identifier.Type)
			return ace._verifyChallengeCleanup(sc, err, id)
		}

		valid, err = ValidateDNS01Challenge(authz.Identifier.Value, cv.Token, cv.Thumbprint, config)
		if err != nil {
			err = fmt.Errorf("%w: error validating dns-01 challenge %v: %v; %v", ErrIncorrectResponse, id, err, ChallengeAttemptFailedMsg)
			return ace._verifyChallengeRetry(sc, cv, authzPath, authz, challenge, err, id)
		}
	case ACMEALPNChallenge:
		if authz.Identifier.Type != ACMEDNSIdentifier {
			err = fmt.Errorf("unsupported identifier type for authorization %v/%v in challenge %v: %v", cv.Account, cv.Authorization, id, authz.Identifier.Type)
			return ace._verifyChallengeCleanup(sc, err, id)
		}

		if authz.Wildcard {
			err = fmt.Errorf("unable to validate wildcard authorization %v/%v in challenge %v via tls-alpn-01 challenge", cv.Account, cv.Authorization, id)
			return ace._verifyChallengeCleanup(sc, err, id)
		}

		valid, err = ValidateTLSALPN01Challenge(authz.Identifier.Value, cv.Token, cv.Thumbprint, config)
		if err != nil {
			err = fmt.Errorf("%w: error validating tls-alpn-01 challenge %v: %s", ErrIncorrectResponse, id, err.Error())
			return ace._verifyChallengeRetry(sc, cv, authzPath, authz, challenge, err, id)
		}
	default:
		err = fmt.Errorf("unsupported ACME challenge type %v for challenge %v", cv.ChallengeType, id)
		return ace._verifyChallengeCleanup(sc, err, id)
	}

	if !valid {
		err = fmt.Errorf("%w: challenge failed with no additional information", ErrIncorrectResponse)
		return ace._verifyChallengeRetry(sc, cv, authzPath, authz, challenge, err, id)
	}

	// If we got here, the challenge verification was successful. Update
	// the authorization appropriately.
	expires := now.Add(15 * 24 * time.Hour)
	challenge.Status = ACMEChallengeValid
	challenge.Validated = now.Format(time.RFC3339)
	challenge.Error = nil
	authz.Status = ACMEAuthorizationValid
	authz.Expires = expires.Format(time.RFC3339)

	if err := saveAuthorizationAtPath(sc, authzPath, authz); err != nil {
		err = fmt.Errorf("error saving updated (validated) authorization %v/%v for challenge %v: %w", cv.Account, cv.Authorization, id, err)
		return ace._verifyChallengeRetry(sc, cv, authzPath, authz, challenge, fmt.Errorf("%w: %s", ErrServerInternal, err.Error()), id)
	}

	return ace._verifyChallengeCleanup(sc, nil, id)
}

func encodeIdentifierForHTTP01Challenge(identifier *ACMEIdentifier) string {
	if !(identifier.Type == ACMEIPIdentifier && identifier.IsV6IP) {
		return identifier.Value
	}

	// If our IPv6 identifier has a zone we need to encode the % to %25 otherwise
	// the http.Client won't properly use it. RFC6874 specifies how the zone
	// identifier in the URI should be handled. In section 2:
	//
	//    According to URI syntax [RFC3986], "%" is always treated as
	//    an escape character in a URI, so, according to the established URI
	//    syntax [RFC3986] any occurrences of literal "%" symbols in a URI MUST
	//    be percent-encoded and represented in the form "%25". Thus, the
	//    scoped address fe80::a%en1 would appear in a URI as
	//    http://[fe80::a%25en1].
	escapedIPv6 := strings.Replace(identifier.Value, "%", "%25", 1)

	// IPv6 addresses need to be surrounded by [] within URLs
	return fmt.Sprintf("[%s]", escapedIPv6)
}

func (ace *ACMEChallengeEngine) _verifyChallengeRetry(sc *storageContext, cv *ChallengeValidation, authzPath string, auth *ACMEAuthorization, challenge *ACMEChallenge, verificationErr error, id string) (bool, time.Time, error) {
	now := time.Now()
	path := acmeValidationPrefix + id

	if err := updateChallengeStatus(sc, cv, authzPath, auth, challenge, verificationErr); err != nil {
		return true, now, err
	}

	if cv.RetryCount > MaxRetryAttempts {
		err := fmt.Errorf("reached max error attempts for challenge %v: %w", id, verificationErr)
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
		return true, now, fmt.Errorf("error persisting updated challenge validation queue entry (error prior to retry, if any: %v): %w", verificationErr, jsonErr)
	}

	if putErr := sc.Storage.Put(sc.Context, json); putErr != nil {
		return true, now, fmt.Errorf("error writing updated challenge validation entry (error prior to retry, if any: %v): %w", verificationErr, putErr)
	}

	if verificationErr != nil {
		verificationErr = fmt.Errorf("retrying validation: %w", verificationErr)
	}

	return true, cv.RetryAfter, verificationErr
}

func updateChallengeStatus(sc *storageContext, cv *ChallengeValidation, authzPath string, auth *ACMEAuthorization, challenge *ACMEChallenge, verificationErr error) error {
	if verificationErr != nil {
		challengeError := TranslateErrorToErrorResponse(verificationErr)
		challenge.Error = challengeError.MarshalForStorage()
	}

	if cv.RetryCount > MaxRetryAttempts {
		challenge.Status = ACMEChallengeInvalid
		auth.Status = ACMEAuthorizationInvalid
	}

	if err := saveAuthorizationAtPath(sc, authzPath, auth); err != nil {
		return fmt.Errorf("error persisting authorization/challenge update: %w", err)
	}
	return nil
}

func (ace *ACMEChallengeEngine) _verifyChallengeCleanup(sc *storageContext, err error, id string) (bool, time.Time, error) {
	now := time.Now()

	// Remove our ChallengeValidation entry only.
	if deleteErr := sc.Storage.Delete(sc.Context, acmeValidationPrefix+id); deleteErr != nil {
		return true, now.Add(1 * time.Second), fmt.Errorf("error deleting challenge %v (error prior to cleanup, if any: %v): %w", id, err, deleteErr)
	}

	if err != nil {
		err = fmt.Errorf("removing challenge validation attempt and not retrying %v; previous error: %w", id, err)
	}

	return false, now, err
}
