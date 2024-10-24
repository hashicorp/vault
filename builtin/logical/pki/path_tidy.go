// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/revocation"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

var tidyCancelledError = errors.New("tidy operation cancelled")

//go:generate enumer -type=tidyStatusState -trimprefix=tidyStatus
type tidyStatusState int

const (
	tidyStatusInactive tidyStatusState = iota
	tidyStatusStarted
	tidyStatusFinished
	tidyStatusError
	tidyStatusCancelling
	tidyStatusCancelled
)

type tidyStatus struct {
	// Parameters used to initiate the operation
	safetyBuffer            int
	issuerSafetyBuffer      int
	revQueueSafetyBuffer    int
	acmeAccountSafetyBuffer int

	tidyCertStore         bool
	tidyRevokedCerts      bool
	tidyRevokedAssocs     bool
	tidyExpiredIssuers    bool
	tidyBackupBundle      bool
	tidyRevocationQueue   bool
	tidyCrossRevokedCerts bool
	tidyAcme              bool
	tidyCertMetadata      bool
	tidyCMPV2NonceStore   bool
	pauseDuration         string

	// Status
	state        tidyStatusState
	err          error
	timeStarted  time.Time
	timeFinished time.Time
	message      string

	// These counts use a custom incrementer that grab and release
	// a lock prior to reading.
	certStoreDeletedCount    uint
	revokedCertDeletedCount  uint
	missingIssuerCertCount   uint
	revQueueDeletedCount     uint
	crossRevokedDeletedCount uint
	certMetadataDeletedCount uint
	cmpv2NonceDeletedCount   uint

	acmeAccountsCount        uint
	acmeAccountsRevokedCount uint
	acmeAccountsDeletedCount uint
	acmeOrdersDeletedCount   uint
}

type tidyConfig struct {
	// AutoTidy config
	Enabled           bool          `json:"enabled"`
	Interval          time.Duration `json:"interval_duration"`
	MinStartupBackoff time.Duration `json:"min_startup_backoff_duration"`
	MaxStartupBackoff time.Duration `json:"max_startup_backoff_duration"`

	// Tidy Operations
	CertStore         bool `json:"tidy_cert_store"`
	RevokedCerts      bool `json:"tidy_revoked_certs"`
	IssuerAssocs      bool `json:"tidy_revoked_cert_issuer_associations"`
	ExpiredIssuers    bool `json:"tidy_expired_issuers"`
	BackupBundle      bool `json:"tidy_move_legacy_ca_bundle"`
	RevocationQueue   bool `json:"tidy_revocation_queue"`
	CrossRevokedCerts bool `json:"tidy_cross_cluster_revoked_certs"`
	TidyAcme          bool `json:"tidy_acme"`
	CertMetadata      bool `json:"tidy_cert_metadata"`
	CMPV2NonceStore   bool `json:"tidy_cmpv2_nonce_store"`

	// Safety Buffers
	SafetyBuffer            time.Duration `json:"safety_buffer"`
	IssuerSafetyBuffer      time.Duration `json:"issuer_safety_buffer"`
	QueueSafetyBuffer       time.Duration `json:"revocation_queue_safety_buffer"`
	AcmeAccountSafetyBuffer time.Duration `json:"acme_account_safety_buffer"`
	PauseDuration           time.Duration `json:"pause_duration"`

	// Metrics.
	MaintainCount  bool `json:"maintain_stored_certificate_counts"`
	PublishMetrics bool `json:"publish_stored_certificate_count_metrics"`
}

func (tc *tidyConfig) IsAnyTidyEnabled() bool {
	return tc.CertStore || tc.RevokedCerts || tc.IssuerAssocs || tc.ExpiredIssuers || tc.BackupBundle || tc.TidyAcme || tc.CrossRevokedCerts || tc.RevocationQueue || tc.CertMetadata || tc.CMPV2NonceStore
}

func (tc *tidyConfig) AnyTidyConfig() string {
	return "tidy_cert_store / tidy_revoked_certs / tidy_revoked_cert_issuer_associations / tidy_expired_issuers / tidy_move_legacy_ca_bundle / tidy_revocation_queue / tidy_cross_cluster_revoked_certs / tidy_acme"
}

func (tc *tidyConfig) CalculateStartupBackoff(mountStartup time.Time) time.Time {
	minBackoff := int64(tc.MinStartupBackoff.Seconds())
	maxBackoff := int64(tc.MaxStartupBackoff.Seconds())

	maxNumber := maxBackoff - minBackoff
	if maxNumber <= 0 {
		return mountStartup.Add(tc.MinStartupBackoff)
	}

	backoffSecs := rand.Int64N(maxNumber) + minBackoff
	return mountStartup.Add(time.Duration(backoffSecs) * time.Second)
}

var defaultTidyConfig = tidyConfig{
	Enabled:                 false,
	Interval:                12 * time.Hour,
	MinStartupBackoff:       5 * time.Minute,
	MaxStartupBackoff:       15 * time.Minute,
	CertStore:               false,
	RevokedCerts:            false,
	IssuerAssocs:            false,
	ExpiredIssuers:          false,
	BackupBundle:            false,
	TidyAcme:                false,
	SafetyBuffer:            72 * time.Hour,
	IssuerSafetyBuffer:      365 * 24 * time.Hour,
	AcmeAccountSafetyBuffer: 30 * 24 * time.Hour,
	PauseDuration:           0 * time.Second,
	MaintainCount:           false,
	PublishMetrics:          false,
	RevocationQueue:         false,
	QueueSafetyBuffer:       48 * time.Hour,
	CrossRevokedCerts:       false,
	CertMetadata:            false,
	CMPV2NonceStore:         false,
}

var tidyStatusResponseFields = map[string]*framework.FieldSchema{
	"safety_buffer": {
		Type:        framework.TypeInt,
		Description: `Safety buffer time duration`,
		Required:    true,
	},
	"issuer_safety_buffer": {
		Type:        framework.TypeInt,
		Description: `Issuer safety buffer`,
		Required:    true,
	},
	"revocation_queue_safety_buffer": {
		Type:        framework.TypeInt,
		Description: `Revocation queue safety buffer`,
		Required:    true,
	},
	"acme_account_safety_buffer": {
		Type:        framework.TypeInt,
		Description: `Safety buffer after creation after which accounts lacking orders are revoked`,
		Required:    false,
	},
	"tidy_cert_store": {
		Type:        framework.TypeBool,
		Description: `Tidy certificate store`,
		Required:    true,
	},
	"tidy_revoked_certs": {
		Type:        framework.TypeBool,
		Description: `Tidy revoked certificates`,
		Required:    true,
	},
	"tidy_revoked_cert_issuer_associations": {
		Type:        framework.TypeBool,
		Description: `Tidy revoked certificate issuer associations`,
		Required:    true,
	},
	"tidy_expired_issuers": {
		Type:        framework.TypeBool,
		Description: `Tidy expired issuers`,
		Required:    true,
	},
	"tidy_cross_cluster_revoked_certs": {
		Type:        framework.TypeBool,
		Description: `Tidy the cross-cluster revoked certificate store`,
		Required:    false,
	},
	"tidy_acme": {
		Type:        framework.TypeBool,
		Description: `Tidy Unused Acme Accounts, and Orders`,
		Required:    true,
	},
	"tidy_cert_metadata": {
		Type:        framework.TypeBool,
		Description: `Tidy cert metadata`,
		Required:    true,
	},
	"tidy_cmpv2_nonce_store": {
		Type:        framework.TypeBool,
		Description: `Tidy CMPv2 nonce store`,
		Required:    true,
	},
	"pause_duration": {
		Type:        framework.TypeString,
		Description: `Duration to pause between tidying certificates`,
		Required:    true,
	},
	"state": {
		Type:        framework.TypeString,
		Description: `One of Inactive, Running, Finished, or Error`,
		Required:    true,
	},
	"error": {
		Type:        framework.TypeString,
		Description: `The error message`,
		Required:    true,
	},
	"time_started": {
		Type:        framework.TypeString,
		Description: `Time the operation started`,
		Required:    true,
	},
	"time_finished": {
		Type:        framework.TypeString,
		Description: `Time the operation finished`,
		Required:    false,
	},
	"last_auto_tidy_finished": {
		Type:        framework.TypeString,
		Description: `Time the last auto-tidy operation finished`,
		Required:    true,
	},
	"message": {
		Type:        framework.TypeString,
		Description: `Message of the operation`,
		Required:    true,
	},
	"cert_store_deleted_count": {
		Type:        framework.TypeInt,
		Description: `The number of certificate storage entries deleted`,
		Required:    true,
	},
	"revoked_cert_deleted_count": {
		Type:        framework.TypeInt,
		Description: `The number of revoked certificate entries deleted`,
		Required:    true,
	},
	"current_cert_store_count": {
		Type:        framework.TypeInt,
		Description: `The number of revoked certificate entries deleted`,
		Required:    true,
	},
	"cross_revoked_cert_deleted_count": {
		Type:        framework.TypeInt,
		Description: ``,
		Required:    true,
	},
	"current_revoked_cert_count": {
		Type:        framework.TypeInt,
		Description: `The number of revoked certificate entries deleted`,
		Required:    true,
	},
	"revocation_queue_deleted_count": {
		Type:     framework.TypeInt,
		Required: true,
	},
	"tidy_move_legacy_ca_bundle": {
		Type:     framework.TypeBool,
		Required: true,
	},
	"tidy_revocation_queue": {
		Type:     framework.TypeBool,
		Required: true,
	},
	"missing_issuer_cert_count": {
		Type:     framework.TypeInt,
		Required: true,
	},
	"internal_backend_uuid": {
		Type:     framework.TypeString,
		Required: true,
	},
	"total_acme_account_count": {
		Type:        framework.TypeInt,
		Description: `Total number of acme accounts iterated over`,
		Required:    false,
	},
	"acme_account_deleted_count": {
		Type:        framework.TypeInt,
		Description: `The number of revoked acme accounts removed`,
		Required:    false,
	},
	"acme_account_revoked_count": {
		Type:        framework.TypeInt,
		Description: `The number of unused acme accounts revoked`,
		Required:    false,
	},
	"acme_orders_deleted_count": {
		Type:        framework.TypeInt,
		Description: `The number of expired, unused acme orders removed`,
		Required:    false,
	},
	"cert_metadata_deleted_count": {
		Type:        framework.TypeInt,
		Description: `The number of metadata entries removed`,
		Required:    false,
	},
	"cmpv2_nonce_deleted_count": {
		Type:        framework.TypeInt,
		Description: `The number of CMPv2 nonces removed`,
		Required:    false,
	},
}

func pathTidy(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationVerb:   "tidy",
		},

		Fields: addTidyFields(map[string]*framework.FieldSchema{}),
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathTidyWrite,
				Responses: map[int][]framework.Response{
					http.StatusAccepted: {{
						Description: "Accepted",
						Fields:      map[string]*framework.FieldSchema{},
					}},
				},
				ForwardPerformanceStandby: true,
			},
		},
		HelpSynopsis:    pathTidyHelpSyn,
		HelpDescription: pathTidyHelpDesc,
	}
}

func pathTidyCancel(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy-cancel$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationVerb:   "tidy",
			OperationSuffix: "cancel",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathTidyCancelWrite,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields:      tidyStatusResponseFields,
					}},
				},
				ForwardPerformanceStandby: true,
			},
		},
		HelpSynopsis:    pathTidyCancelHelpSyn,
		HelpDescription: pathTidyCancelHelpDesc,
	}
}

func pathTidyStatus(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy-status$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
			OperationVerb:   "tidy",
			OperationSuffix: "status",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathTidyStatusRead,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields:      tidyStatusResponseFields,
					}},
				},
				ForwardPerformanceStandby: true,
			},
		},
		HelpSynopsis:    pathTidyStatusHelpSyn,
		HelpDescription: pathTidyStatusHelpDesc,
	}
}

func pathConfigAutoTidy(b *backend) *framework.Path {
	autoTidyResponseFields := map[string]*framework.FieldSchema{
		"enabled": {
			Type:        framework.TypeBool,
			Description: `Specifies whether automatic tidy is enabled or not`,
			Required:    true,
		},
		"min_startup_backoff_duration": {
			Type:        framework.TypeInt,
			Description: `The minimum amount of time in seconds auto-tidy will be delayed after startup`,
			Required:    true,
		},
		"max_startup_backoff_duration": {
			Type:        framework.TypeInt,
			Description: `The maximum amount of time in seconds auto-tidy will be delayed after startup`,
			Required:    true,
		},
		"interval_duration": {
			Type:        framework.TypeInt,
			Description: `Specifies the duration between automatic tidy operation`,
			Required:    true,
		},
		"tidy_cert_store": {
			Type:        framework.TypeBool,
			Description: `Specifies whether to tidy up the certificate store`,
			Required:    true,
		},
		"tidy_revoked_certs": {
			Type:        framework.TypeBool,
			Description: `Specifies whether to remove all invalid and expired certificates from storage`,
			Required:    true,
		},
		"tidy_revoked_cert_issuer_associations": {
			Type:        framework.TypeBool,
			Description: `Specifies whether to associate revoked certificates with their corresponding issuers`,
			Required:    true,
		},
		"tidy_expired_issuers": {
			Type:        framework.TypeBool,
			Description: `Specifies whether tidy expired issuers`,
			Required:    true,
		},
		"tidy_acme": {
			Type:        framework.TypeBool,
			Description: `Tidy Unused Acme Accounts, and Orders`,
			Required:    true,
		},
		"tidy_cert_metadata": {
			Type:        framework.TypeBool,
			Description: `Tidy cert metadata`,
			Required:    true,
		},
		"tidy_cmpv2_nonce_store": {
			Type:        framework.TypeBool,
			Description: `Tidy CMPv2 nonce store`,
			Required:    true,
		},
		"safety_buffer": {
			Type:        framework.TypeInt,
			Description: `Safety buffer time duration`,
			Required:    true,
		},
		"issuer_safety_buffer": {
			Type:        framework.TypeInt,
			Description: `Issuer safety buffer`,
			Required:    true,
		},
		"acme_account_safety_buffer": {
			Type:        framework.TypeInt,
			Description: `Safety buffer after creation after which accounts lacking orders are revoked`,
			Required:    true,
		},
		"pause_duration": {
			Type:        framework.TypeString,
			Description: `Duration to pause between tidying certificates`,
			Required:    true,
		},
		"tidy_cross_cluster_revoked_certs": {
			Type:        framework.TypeBool,
			Description: `Tidy the cross-cluster revoked certificate store`,
			Required:    true,
		},
		"tidy_revocation_queue": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"tidy_move_legacy_ca_bundle": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"revocation_queue_safety_buffer": {
			Type:     framework.TypeInt,
			Required: true,
		},
		"publish_stored_certificate_count_metrics": {
			Type:     framework.TypeBool,
			Required: true,
		},
		"maintain_stored_certificate_counts": {
			Type:     framework.TypeBool,
			Required: true,
		},
	}
	return &framework.Path{
		Pattern: "config/auto-tidy",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixPKI,
		},
		Fields: addTidyFields(map[string]*framework.FieldSchema{
			"enabled": {
				Type:        framework.TypeBool,
				Description: `Set to true to enable automatic tidy operations.`,
			},
			"min_startup_backoff_duration": {
				Type:        framework.TypeDurationSecond,
				Description: `The minimum amount of time in seconds auto-tidy will be delayed after startup.`,
				Default:     int(defaultTidyConfig.MinStartupBackoff.Seconds()),
			},
			"max_startup_backoff_duration": {
				Type:        framework.TypeDurationSecond,
				Description: `The maximum amount of time in seconds auto-tidy will be delayed after startup.`,
				Default:     int(defaultTidyConfig.MaxStartupBackoff.Seconds()),
			},
			"interval_duration": {
				Type:        framework.TypeDurationSecond,
				Description: `Interval at which to run an auto-tidy operation. This is the time between tidy invocations (after one finishes to the start of the next). Running a manual tidy will reset this duration.`,
				Default:     int(defaultTidyConfig.Interval / time.Second), // TypeDurationSecond currently requires the default to be an int.
			},
			"maintain_stored_certificate_counts": {
				Type: framework.TypeBool,
				Description: `This configures whether stored certificates
are counted upon initialization of the backend, and whether during
normal operation, a running count of certificates stored is maintained.`,
				Default: false,
			},
			"publish_stored_certificate_count_metrics": {
				Type: framework.TypeBool,
				Description: `This configures whether the stored certificate
count is published to the metrics consumer.  It does not affect if the
stored certificate count is maintained, and if maintained, it will be
available on the tidy-status endpoint.`,
				Default: false,
			},
		}),
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigAutoTidyRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "auto-tidy-configuration",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields:      autoTidyResponseFields,
					}},
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigAutoTidyWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "configure",
					OperationSuffix: "auto-tidy",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields:      autoTidyResponseFields,
					}},
				},
				// Read more about why these flags are set in backend.go.
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},
		HelpSynopsis:    pathConfigAutoTidySyn,
		HelpDescription: pathConfigAutoTidyDesc,
	}
}

func (b *backend) pathTidyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	safetyBuffer := d.Get("safety_buffer").(int)
	tidyCertStore := d.Get("tidy_cert_store").(bool)
	tidyRevokedCerts := d.Get("tidy_revoked_certs").(bool) || d.Get("tidy_revocation_list").(bool)
	tidyRevokedAssocs := d.Get("tidy_revoked_cert_issuer_associations").(bool)
	tidyExpiredIssuers := d.Get("tidy_expired_issuers").(bool)
	tidyBackupBundle := d.Get("tidy_move_legacy_ca_bundle").(bool)
	issuerSafetyBuffer := d.Get("issuer_safety_buffer").(int)
	pauseDurationStr := d.Get("pause_duration").(string)
	pauseDuration := 0 * time.Second
	tidyRevocationQueue := d.Get("tidy_revocation_queue").(bool)
	queueSafetyBuffer := d.Get("revocation_queue_safety_buffer").(int)
	tidyCrossRevokedCerts := d.Get("tidy_cross_cluster_revoked_certs").(bool)
	tidyAcme := d.Get("tidy_acme").(bool)
	acmeAccountSafetyBuffer := d.Get("acme_account_safety_buffer").(int)
	tidyCertMetadata := d.Get("tidy_cert_metadata").(bool)
	tidyCMPV2NonceStore := d.Get("tidy_cmpv2_nonce_store").(bool)

	if safetyBuffer < 1 {
		return logical.ErrorResponse("safety_buffer must be greater than zero"), nil
	}

	if issuerSafetyBuffer < 1 {
		return logical.ErrorResponse("issuer_safety_buffer must be greater than zero"), nil
	}

	if queueSafetyBuffer < 1 {
		return logical.ErrorResponse("revocation_queue_safety_buffer must be greater than zero"), nil
	}

	if acmeAccountSafetyBuffer < 1 {
		return logical.ErrorResponse("acme_account_safety_buffer must be greater than zero"), nil
	}

	if pauseDurationStr != "" {
		var err error
		pauseDuration, err = parseutil.ParseDurationSecond(pauseDurationStr)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error parsing pause_duration: %v", err)), nil
		}

		if pauseDuration < (0 * time.Second) {
			return logical.ErrorResponse("received invalid, negative pause_duration"), nil
		}
	}

	if tidyCertMetadata && !constants.IsEnterprise {
		return logical.ErrorResponse("certificate metadata is only supported on Vault Enterprise"), nil
	}

	bufferDuration := time.Duration(safetyBuffer) * time.Second
	issuerBufferDuration := time.Duration(issuerSafetyBuffer) * time.Second
	queueSafetyBufferDuration := time.Duration(queueSafetyBuffer) * time.Second
	acmeAccountSafetyBufferDuration := time.Duration(acmeAccountSafetyBuffer) * time.Second

	// Manual run with constructed configuration.
	config := &tidyConfig{
		Enabled:                 true,
		Interval:                0 * time.Second,
		CertStore:               tidyCertStore,
		RevokedCerts:            tidyRevokedCerts,
		IssuerAssocs:            tidyRevokedAssocs,
		ExpiredIssuers:          tidyExpiredIssuers,
		BackupBundle:            tidyBackupBundle,
		SafetyBuffer:            bufferDuration,
		IssuerSafetyBuffer:      issuerBufferDuration,
		PauseDuration:           pauseDuration,
		RevocationQueue:         tidyRevocationQueue,
		QueueSafetyBuffer:       queueSafetyBufferDuration,
		CrossRevokedCerts:       tidyCrossRevokedCerts,
		TidyAcme:                tidyAcme,
		AcmeAccountSafetyBuffer: acmeAccountSafetyBufferDuration,
		CertMetadata:            tidyCertMetadata,
		CMPV2NonceStore:         tidyCMPV2NonceStore,
	}

	if !atomic.CompareAndSwapUint32(b.tidyCASGuard, 0, 1) {
		resp := &logical.Response{}
		resp.AddWarning("Tidy operation already in progress.")
		return resp, nil
	}

	// Tests using framework will screw up the storage so make a locally
	// scoped req to hold a reference
	req = &logical.Request{
		Storage: req.Storage,
	}

	resp := &logical.Response{}
	// Mark the last tidy operation as relatively recent, to ensure we don't
	// try to trigger the periodic function.
	// NOTE: not sure this is correct as we are updating the auto tidy time with this manual run. Ideally we
	//       could track when we ran each type of tidy was last run which would allow manual runs and auto
	//       runs to properly impact each other.
	sc := b.makeStorageContext(ctx, req.Storage)
	if err := b.updateLastAutoTidyTime(sc, time.Now()); err != nil {
		resp.AddWarning(fmt.Sprintf("failed persisting tidy last run time: %v", err))
	}

	// Kick off the actual tidy.
	b.startTidyOperation(req, config)

	if !config.IsAnyTidyEnabled() {
		resp.AddWarning("Manual tidy requested but no tidy operations were set. Enable at least one tidy operation to be run (" + config.AnyTidyConfig() + ").")
	} else {
		resp.AddWarning("Tidy operation successfully started. Any information from the operation will be printed to Vault's server logs.")
	}

	if tidyRevocationQueue || tidyCrossRevokedCerts {
		isNotPerfPrimary := b.System().ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
			(!b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary))
		if isNotPerfPrimary {
			resp.AddWarning("tidy_revocation_queue=true and tidy_cross_cluster_revoked_certs=true can only be set on the active node of the primary cluster unless a local mount is used; this option has been ignored.")
		}
	}

	return logical.RespondWithStatusCode(resp, req, http.StatusAccepted)
}

func (b *backend) startTidyOperation(req *logical.Request, config *tidyConfig) {
	go func() {
		atomic.StoreUint32(b.tidyCancelCAS, 0)
		defer atomic.StoreUint32(b.tidyCASGuard, 0)

		b.tidyStatusStart(config)

		// Don't cancel when the original client request goes away.
		ctx := context.Background()

		logger := b.Logger().Named("tidy")

		doTidy := func() error {
			if config.CertStore {
				if err := b.doTidyCertStore(ctx, req, logger, config); err != nil {
					return err
				}
			}

			// Check for cancel before continuing.
			if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
				return tidyCancelledError
			}

			if config.RevokedCerts || config.IssuerAssocs {
				if err := b.doTidyRevocationStore(ctx, req, logger, config); err != nil {
					return err
				}
			}

			// Check for cancel before continuing.
			if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
				return tidyCancelledError
			}

			if config.ExpiredIssuers {
				if err := b.doTidyExpiredIssuers(ctx, req, logger, config); err != nil {
					return err
				}
			}

			// Check for cancel before continuing.
			if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
				return tidyCancelledError
			}

			if config.BackupBundle {
				if err := b.doTidyMoveCABundle(ctx, req, logger, config); err != nil {
					return err
				}
			}

			// Check for cancel before continuing.
			if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
				return tidyCancelledError
			}

			if config.RevocationQueue {
				if err := b.doTidyRevocationQueue(ctx, req, logger, config); err != nil {
					return err
				}
			}

			// Check for cancel before continuing.
			if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
				return tidyCancelledError
			}

			if config.CrossRevokedCerts {
				if err := b.doTidyCrossRevocationStore(ctx, req, logger, config); err != nil {
					return err
				}
			}

			// Check for cancel before continuing.
			if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
				return tidyCancelledError
			}

			if config.TidyAcme {
				if err := b.doTidyAcme(ctx, req, logger, config); err != nil {
					return err
				}
			}

			// Check for cancel before continuing.
			if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
				return tidyCancelledError
			}

			if config.CertMetadata {
				if err := b.doTidyCertMetadata(ctx, req, logger, config); err != nil {
					return err
				}
			}

			// Check for cancel before continuing.
			if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
				return tidyCancelledError
			}

			if config.CMPV2NonceStore {
				if err := b.doTidyCMPV2NonceStore(ctx, req.Storage); err != nil {
					return err
				}
			}

			return nil
		}

		if err := doTidy(); err != nil {
			logger.Error("error running tidy", "error", err)
			b.tidyStatusStop(err)
		} else {
			b.tidyStatusStop(nil)

			// Since the tidy operation finished without an error, we don't
			// really want to start another tidy right away (if the interval
			// is too short). So mark the last tidy as now.
			sc := b.makeStorageContext(ctx, req.Storage)
			if err := b.updateLastAutoTidyTime(sc, time.Now()); err != nil {
				logger.Error("error persisting last tidy run time", "error", err)
			}
		}
	}()
}

func (b *backend) doTidyCertStore(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	serials, err := req.Storage.List(ctx, issuing.PathCerts)
	if err != nil {
		return fmt.Errorf("error fetching list of certs: %w", err)
	}

	serialCount := len(serials)
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "cert_store_total_entries"}, float32(serialCount))
	for i, serial := range serials {
		b.tidyStatusMessage(fmt.Sprintf("Tidying certificate store: checking entry %d of %d", i, serialCount))
		metrics.SetGauge([]string{"secrets", "pki", "tidy", "cert_store_current_entry"}, float32(i))

		// Check for cancel before continuing.
		if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
			return tidyCancelledError
		}

		// Check for pause duration to reduce resource consumption.
		if config.PauseDuration > (0 * time.Second) {
			time.Sleep(config.PauseDuration)
		}

		certEntry, err := req.Storage.Get(ctx, issuing.PathCerts+serial)
		if err != nil {
			return fmt.Errorf("error fetching certificate %q: %w", serial, err)
		}

		if certEntry == nil {
			logger.Warn("certificate entry is nil; tidying up since it is no longer useful for any server operations", "serial", serial)
			if err := req.Storage.Delete(ctx, issuing.PathCerts+serial); err != nil {
				return fmt.Errorf("error deleting nil entry with serial %s: %w", serial, err)
			}
			b.tidyStatusIncCertStoreCount()
			continue
		}

		if certEntry.Value == nil || len(certEntry.Value) == 0 {
			logger.Warn("certificate entry has no value; tidying up since it is no longer useful for any server operations", "serial", serial)
			if err := req.Storage.Delete(ctx, issuing.PathCerts+serial); err != nil {
				return fmt.Errorf("error deleting entry with nil value with serial %s: %w", serial, err)
			}
			b.tidyStatusIncCertStoreCount()
			continue
		}

		cert, err := x509.ParseCertificate(certEntry.Value)
		if err != nil {
			return fmt.Errorf("unable to parse stored certificate with serial %q: %w", serial, err)
		}

		if time.Since(cert.NotAfter) > config.SafetyBuffer {
			if err := req.Storage.Delete(ctx, issuing.PathCerts+serial); err != nil {
				return fmt.Errorf("error deleting serial %q from storage: %w", serial, err)
			}
			b.tidyStatusIncCertStoreCount()
		}
	}

	b.tidyStatusLock.RLock()
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "cert_store_total_entries_remaining"}, float32(uint(serialCount)-b.tidyStatus.certStoreDeletedCount))
	b.tidyStatusLock.RUnlock()

	return nil
}

func (b *backend) doTidyRevocationStore(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	b.GetRevokeStorageLock().Lock()
	defer b.GetRevokeStorageLock().Unlock()

	// Fetch and parse our issuers so we can associate them if necessary.
	sc := b.makeStorageContext(ctx, req.Storage)
	issuerIDCertMap, err := revocation.FetchIssuerMapForRevocationChecking(sc)
	if err != nil {
		return err
	}

	rebuildCRL := false

	revokedSerials, err := req.Storage.List(ctx, "revoked/")
	if err != nil {
		return fmt.Errorf("error fetching list of revoked certs: %w", err)
	}

	revokedSerialsCount := len(revokedSerials)
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "revoked_cert_total_entries"}, float32(revokedSerialsCount))

	fixedIssuers := 0

	var revInfo revocation.RevocationInfo
	for i, serial := range revokedSerials {
		b.tidyStatusMessage(fmt.Sprintf("Tidying revoked certificates: checking certificate %d of %d", i, len(revokedSerials)))
		metrics.SetGauge([]string{"secrets", "pki", "tidy", "revoked_cert_current_entry"}, float32(i))

		// Check for cancel before continuing.
		if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
			return tidyCancelledError
		}

		// Check for pause duration to reduce resource consumption.
		if config.PauseDuration > (0 * time.Second) {
			b.GetRevokeStorageLock().Unlock()
			time.Sleep(config.PauseDuration)
			b.GetRevokeStorageLock().Lock()
		}

		revokedEntry, err := req.Storage.Get(ctx, "revoked/"+serial)
		if err != nil {
			return fmt.Errorf("unable to fetch revoked cert with serial %q: %w", serial, err)
		}

		if revokedEntry == nil {
			logger.Warn("revoked entry is nil; tidying up since it is no longer useful for any server operations", "serial", serial)
			if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
				return fmt.Errorf("error deleting nil revoked entry with serial %s: %w", serial, err)
			}
			b.tidyStatusIncRevokedCertCount()
			continue
		}

		if revokedEntry.Value == nil || len(revokedEntry.Value) == 0 {
			logger.Warn("revoked entry has nil value; tidying up since it is no longer useful for any server operations", "serial", serial)
			if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
				return fmt.Errorf("error deleting revoked entry with nil value with serial %s: %w", serial, err)
			}
			b.tidyStatusIncRevokedCertCount()
			continue
		}

		err = revokedEntry.DecodeJSON(&revInfo)
		if err != nil {
			return fmt.Errorf("error decoding revocation entry for serial %q: %w", serial, err)
		}

		revokedCert, err := x509.ParseCertificate(revInfo.CertificateBytes)
		if err != nil {
			return fmt.Errorf("unable to parse stored revoked certificate with serial %q: %w", serial, err)
		}

		// Tidy operations over revoked certs should execute prior to
		// tidyRevokedCerts as that may remove the entry. If that happens,
		// we won't persist the revInfo changes (as it was deleted instead).
		var storeCert bool
		if config.IssuerAssocs {
			if !isRevInfoIssuerValid(&revInfo, issuerIDCertMap) {
				b.tidyStatusIncMissingIssuerCertCount()
				revInfo.CertificateIssuer = issuing.IssuerID("")
				storeCert = true
				if revInfo.AssociateRevokedCertWithIsssuer(revokedCert, issuerIDCertMap) {
					fixedIssuers += 1
				}
			}
		}

		if config.RevokedCerts {
			// Only remove the entries from revoked/ and certs/ if we're
			// past its NotAfter value. This is because we use the
			// information on revoked/ to build the CRL and the
			// information on certs/ for lookup.
			if time.Since(revokedCert.NotAfter) > config.SafetyBuffer {
				if err := req.Storage.Delete(ctx, "revoked/"+serial); err != nil {
					return fmt.Errorf("error deleting serial %q from revoked list: %w", serial, err)
				}
				if err := req.Storage.Delete(ctx, issuing.PathCerts+serial); err != nil {
					return fmt.Errorf("error deleting serial %q from store when tidying revoked: %w", serial, err)
				}
				rebuildCRL = true
				storeCert = false
				b.tidyStatusIncRevokedCertCount()
			}
		}

		// If the entry wasn't removed but was otherwise modified,
		// go ahead and write it back out.
		if storeCert {
			revokedEntry, err = logical.StorageEntryJSON("revoked/"+serial, revInfo)
			if err != nil {
				return fmt.Errorf("error building entry to persist changes to serial %v from revoked list: %w", serial, err)
			}

			err = req.Storage.Put(ctx, revokedEntry)
			if err != nil {
				return fmt.Errorf("error persisting changes to serial %v from revoked list: %w", serial, err)
			}
		}
	}

	b.tidyStatusLock.RLock()
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "revoked_cert_total_entries_remaining"}, float32(uint(revokedSerialsCount)-b.tidyStatus.revokedCertDeletedCount))
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "revoked_cert_entries_incorrect_issuers"}, float32(b.tidyStatus.missingIssuerCertCount))
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "revoked_cert_entries_fixed_issuers"}, float32(fixedIssuers))
	b.tidyStatusLock.RUnlock()

	if rebuildCRL {
		// Expired certificates isn't generally an important
		// reason to trigger a CRL rebuild for. Check if
		// automatic CRL rebuilds have been enabled and defer
		// the rebuild if so.
		config, err := sc.getRevocationConfig()
		if err != nil {
			return err
		}

		if !config.AutoRebuild {
			warnings, err := b.CrlBuilder().Rebuild(sc, false)
			if err != nil {
				return err
			}
			if len(warnings) > 0 {
				msg := "During rebuild of CRL for tidy, got the following warnings:"
				for index, warning := range warnings {
					msg = fmt.Sprintf("%v\n %d. %v", msg, index+1, warning)
				}
				b.Logger().Warn(msg)
			}
		}
	}

	return nil
}

func (b *backend) doTidyExpiredIssuers(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	// We do not support cancelling within the expired issuers operation.
	// Any cancellation will occur before or after this operation.

	if b.System().ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary)) {
		b.Logger().Debug("skipping expired issuer tidy as we're not on the primary or secondary with a local mount")
		return nil
	}

	// Short-circuit to avoid having to deal with the legacy mounts. While we
	// could handle this case and remove these issuers, its somewhat
	// unexpected behavior and we'd prefer to finish the migration first.
	if b.UseLegacyBundleCaStorage() {
		return nil
	}

	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	// Fetch and parse our issuers so we have their expiration date.
	sc := b.makeStorageContext(ctx, req.Storage)
	issuerIDCertMap, err := revocation.FetchIssuerMapForRevocationChecking(sc)
	if err != nil {
		return err
	}

	// Fetch the issuer config to find the default; we don't want to remove
	// the current active issuer automatically.
	iConfig, err := sc.getIssuersConfig()
	if err != nil {
		return err
	}

	// We want certificates which have expired before this date by a given
	// safety buffer.
	rebuildChainsAndCRL := false

	for issuer, cert := range issuerIDCertMap {
		if time.Since(cert.NotAfter) <= config.IssuerSafetyBuffer {
			continue
		}

		entry, err := sc.fetchIssuerById(issuer)
		if err != nil {
			return nil
		}

		// This issuer's certificate has expired. We explicitly persist the
		// key, but log both the certificate and the keyId to the
		// informational logs so an admin can recover the removed cert if
		// necessary or remove the key (and know which cert it belonged to),
		// if desired.
		msg := "[Tidy on mount: %v] Issuer %v has expired by %v and is being removed."
		idAndName := fmt.Sprintf("[id:%v/name:%v]", entry.ID, entry.Name)
		msg = fmt.Sprintf(msg, b.backendUUID, idAndName, config.IssuerSafetyBuffer)

		// Before we log, check if we're the default. While this is late, and
		// after we read it from storage, we have more info here to tell the
		// user that their default has expired AND has passed the safety
		// buffer.
		if iConfig.DefaultIssuerId == issuer {
			msg = "[Tidy on mount: %v] Issuer %v has expired and would be removed via tidy, but won't be, as it is currently the default issuer."
			msg = fmt.Sprintf(msg, b.backendUUID, idAndName)
			b.Logger().Warn(msg)
			continue
		}

		// Log the above message..
		b.Logger().Info(msg, "serial_number", entry.SerialNumber, "key_id", entry.KeyID, "certificate", entry.Certificate)

		wasDefault, err := sc.deleteIssuer(issuer)
		if err != nil {
			b.Logger().Error(fmt.Sprintf("failed to remove %v: %v", idAndName, err))
			return err
		}
		if wasDefault {
			b.Logger().Warn(fmt.Sprintf("expired issuer %v was default; it is strongly encouraged to choose a new default issuer for backwards compatibility", idAndName))
		}

		rebuildChainsAndCRL = true
	}

	if rebuildChainsAndCRL {
		// When issuers are removed, there's a chance chains change as a
		// result; remove them.
		if err := sc.rebuildIssuersChains(nil); err != nil {
			return err
		}

		// Removal of issuers is generally a good reason to rebuild the CRL,
		// even if auto-rebuild is enabled.
		b.GetRevokeStorageLock().Lock()
		defer b.GetRevokeStorageLock().Unlock()

		warnings, err := b.CrlBuilder().Rebuild(sc, false)
		if err != nil {
			return err
		}
		if len(warnings) > 0 {
			msg := "During rebuild of CRL for tidy, got the following warnings:"
			for index, warning := range warnings {
				msg = fmt.Sprintf("%v\n %d. %v", msg, index+1, warning)
			}
			b.Logger().Warn(msg)
		}
	}

	return nil
}

func (b *backend) doTidyMoveCABundle(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	// We do not support cancelling within this operation; any cancel will
	// occur before or after this operation.

	if b.System().ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary)) {
		b.Logger().Debug("skipping moving the legacy CA bundle as we're not on the primary or secondary with a local mount")
		return nil
	}

	// Short-circuit to avoid moving the legacy bundle from under a legacy
	// mount.
	if b.UseLegacyBundleCaStorage() {
		return nil
	}

	// If we've already run, exit.
	_, bundle, err := getLegacyCertBundle(ctx, req.Storage)
	if err != nil {
		return fmt.Errorf("failed to fetch the legacy CA bundle: %w", err)
	}

	if bundle == nil {
		b.Logger().Debug("No legacy CA bundle available; nothing to do.")
		return nil
	}

	log, err := getLegacyBundleMigrationLog(ctx, req.Storage)
	if err != nil {
		return fmt.Errorf("failed to fetch the legacy bundle migration log: %w", err)
	}

	if log == nil {
		return fmt.Errorf("refusing to tidy with an empty legacy migration log but present CA bundle: %w", err)
	}

	if time.Since(log.Created) <= config.IssuerSafetyBuffer {
		b.Logger().Debug("Migration was created too recently to remove the legacy bundle; refusing to move legacy CA bundle to backup location.")
		return nil
	}

	// Do the write before the delete.
	entry, err := logical.StorageEntryJSON(legacyCertBundleBackupPath, bundle)
	if err != nil {
		return fmt.Errorf("failed to create new backup storage entry: %w", err)
	}

	err = req.Storage.Put(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to write new backup legacy CA bundle: %w", err)
	}

	err = req.Storage.Delete(ctx, legacyCertBundlePath)
	if err != nil {
		return fmt.Errorf("failed to remove old legacy CA bundle path: %w", err)
	}

	b.Logger().Info("legacy CA bundle successfully moved to backup location")
	return nil
}

func (b *backend) doTidyRevocationQueue(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	if b.System().ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary)) {
		b.Logger().Debug("skipping cross-cluster revocation queue tidy as we're not on the primary or secondary with a local mount")
		return nil
	}

	sc := b.makeStorageContext(ctx, req.Storage)
	clusters, err := sc.Storage.List(sc.Context, crossRevocationPrefix)
	if err != nil {
		return fmt.Errorf("failed to list cross-cluster revocation queue participating clusters: %w", err)
	}

	// Grab locks as we're potentially modifying revocation-related storage.
	b.GetRevokeStorageLock().Lock()
	defer b.GetRevokeStorageLock().Unlock()

	for cIndex, cluster := range clusters {
		if cluster[len(cluster)-1] == '/' {
			cluster = cluster[0 : len(cluster)-1]
		}

		cPath := crossRevocationPrefix + cluster + "/"
		serials, err := sc.Storage.List(sc.Context, cPath)
		if err != nil {
			return fmt.Errorf("failed to list cross-cluster revocation queue entries for cluster %v (%v): %w", cluster, cIndex, err)
		}

		for _, serial := range serials {
			// Check for cancellation.
			if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
				return tidyCancelledError
			}

			// Check for pause duration to reduce resource consumption.
			if config.PauseDuration > (0 * time.Second) {
				b.GetRevokeStorageLock().Unlock()
				time.Sleep(config.PauseDuration)
				b.GetRevokeStorageLock().Lock()
			}

			// Confirmation entries _should_ be handled by this cluster's
			// processRevocationQueue(...) invocation; if not, when the plugin
			// reloads, maybeGatherQueueForFirstProcess(...) will remove all
			// stale confirmation requests. However, we don't want to force an
			// operator to reload their in-use plugin, so allow tidy to also
			// clean up confirmation values without reloading.
			if serial[len(serial)-1] == '/' {
				// Check if we have a confirmed entry.
				confirmedPath := cPath + serial + "confirmed"
				removalEntry, err := sc.Storage.Get(sc.Context, confirmedPath)
				if err != nil {
					return fmt.Errorf("error reading revocation confirmation (%v) during tidy: %w", confirmedPath, err)
				}
				if removalEntry == nil {
					continue
				}

				// Remove potential revocation requests from all clusters.
				for _, subCluster := range clusters {
					if subCluster[len(subCluster)-1] == '/' {
						subCluster = subCluster[0 : len(subCluster)-1]
					}

					reqPath := subCluster + "/" + serial[0:len(serial)-1]
					if err := sc.Storage.Delete(sc.Context, reqPath); err != nil {
						return fmt.Errorf("failed to remove confirmed revocation request on candidate cluster (%v): %w", reqPath, err)
					}
				}

				// Then delete the confirmation.
				if err := sc.Storage.Delete(sc.Context, confirmedPath); err != nil {
					return fmt.Errorf("failed to remove confirmed revocation confirmation (%v): %w", confirmedPath, err)
				}

				// No need to handle a revocation request at this path: it can't
				// still exist on this cluster after we deleted it above.
				continue
			}

			ePath := cPath + serial
			entry, err := sc.Storage.Get(sc.Context, ePath)
			if err != nil {
				return fmt.Errorf("error reading revocation request (%v) to tidy: %w", ePath, err)
			}
			if entry == nil || entry.Value == nil {
				continue
			}

			var revRequest revocationRequest
			if err := entry.DecodeJSON(&revRequest); err != nil {
				return fmt.Errorf("error reading revocation request (%v) to tidy: %w", ePath, err)
			}

			if time.Since(revRequest.RequestedAt) <= config.QueueSafetyBuffer {
				continue
			}

			// Safe to remove this entry.
			if err := sc.Storage.Delete(sc.Context, ePath); err != nil {
				return fmt.Errorf("error deleting revocation request (%v): %w", ePath, err)
			}

			// Assumption: there should never be a need to remove this from
			// the processing queue on this node. We're on the active primary,
			// so our writes don't cause invalidations. This means we'd have
			// to have slated it for deletion very quickly after it'd been
			// sent (i.e., inside of the 1-minute boundary that periodicFunc
			// executes at). While this is possible, because we grab the
			// revocationStorageLock above, we can't execute interleaved
			// with that periodicFunc, so the periodicFunc would've had to
			// finished before we actually did this deletion (or it wouldn't
			// have ignored this serial because our deletion would've
			// happened prior to it reading the storage entry). Thus we should
			// be safe to ignore the revocation queue removal here.
			b.tidyStatusIncRevQueueCount()
		}
	}

	return nil
}

func (b *backend) doTidyCrossRevocationStore(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	if b.System().ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!b.System().LocalMount() && b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary)) {
		b.Logger().Debug("skipping cross-cluster revoked certificate store tidy as we're not on the primary or secondary with a local mount")
		return nil
	}

	sc := b.makeStorageContext(ctx, req.Storage)
	clusters, err := sc.Storage.List(sc.Context, unifiedRevocationReadPathPrefix)
	if err != nil {
		return fmt.Errorf("failed to list cross-cluster revoked certificate store participating clusters: %w", err)
	}

	// Grab locks as we're potentially modifying revocation-related storage.
	b.GetRevokeStorageLock().Lock()
	defer b.GetRevokeStorageLock().Unlock()

	for cIndex, cluster := range clusters {
		if cluster[len(cluster)-1] == '/' {
			cluster = cluster[0 : len(cluster)-1]
		}

		cPath := unifiedRevocationReadPathPrefix + cluster + "/"
		serials, err := sc.Storage.List(sc.Context, cPath)
		if err != nil {
			return fmt.Errorf("failed to list cross-cluster revoked certificate store entries for cluster %v (%v): %w", cluster, cIndex, err)
		}

		for _, serial := range serials {
			// Check for cancellation.
			if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
				return tidyCancelledError
			}

			// Check for pause duration to reduce resource consumption.
			if config.PauseDuration > (0 * time.Second) {
				b.GetRevokeStorageLock().Unlock()
				time.Sleep(config.PauseDuration)
				b.GetRevokeStorageLock().Lock()
			}

			ePath := cPath + serial
			entry, err := sc.Storage.Get(sc.Context, ePath)
			if err != nil {
				return fmt.Errorf("error reading cross-cluster revocation entry (%v) to tidy: %w", ePath, err)
			}
			if entry == nil || entry.Value == nil {
				continue
			}

			var details revocation.UnifiedRevocationEntry
			if err := entry.DecodeJSON(&details); err != nil {
				return fmt.Errorf("error decoding cross-cluster revocation entry (%v) to tidy: %w", ePath, err)
			}

			if time.Since(details.CertExpiration) <= config.SafetyBuffer {
				continue
			}

			// Safe to remove this entry.
			if err := sc.Storage.Delete(sc.Context, ePath); err != nil {
				return fmt.Errorf("error deleting revocation request (%v): %w", ePath, err)
			}

			b.tidyStatusIncCrossRevCertCount()
		}
	}

	return nil
}

func (b *backend) doTidyAcme(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	b.acmeAccountLock.Lock()
	defer b.acmeAccountLock.Unlock()

	sc := b.makeStorageContext(ctx, req.Storage)
	thumbprints, err := sc.Storage.List(ctx, acmeThumbprintPrefix)
	if err != nil {
		return err
	}

	b.tidyStatusLock.Lock()
	b.tidyStatus.acmeAccountsCount = uint(len(thumbprints))
	b.tidyStatusLock.Unlock()

	for _, thumbprint := range thumbprints {
		err := b.tidyAcmeAccountByThumbprint(b.GetAcmeState(), sc, thumbprint, config.SafetyBuffer, config.AcmeAccountSafetyBuffer)
		if err != nil {
			logger.Warn("error tidying account %v: %v", thumbprint, err.Error())
		}

		// Check for cancel before continuing.
		if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
			return tidyCancelledError
		}

		// Check for pause duration to reduce resource consumption.
		if config.PauseDuration > (0 * time.Second) {
			b.acmeAccountLock.Unlock() // Correct the Lock
			time.Sleep(config.PauseDuration)
			b.acmeAccountLock.Lock()
		}

	}

	// Clean up any unused EAB
	eabIds, err := b.GetAcmeState().ListEabIds(sc)
	if err != nil {
		return fmt.Errorf("failed listing EAB ids: %w", err)
	}

	for _, eabId := range eabIds {
		eab, err := b.GetAcmeState().LoadEab(sc, eabId)
		if err != nil {
			if errors.Is(err, ErrStorageItemNotFound) {
				// We don't need to worry about a consumed EAB
				continue
			}
			return err
		}

		eabExpiration := eab.CreatedOn.Add(config.AcmeAccountSafetyBuffer)
		if time.Now().After(eabExpiration) {
			_, err := b.GetAcmeState().DeleteEab(sc, eabId)
			if err != nil {
				return fmt.Errorf("failed to tidy eab %s: %w", eabId, err)
			}
		}

		// Check for cancel before continuing.
		if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
			return tidyCancelledError
		}

		// Check for pause duration to reduce resource consumption.
		if config.PauseDuration > (0 * time.Second) {
			b.acmeAccountLock.Unlock() // Correct the Lock
			time.Sleep(config.PauseDuration)
			b.acmeAccountLock.Lock()
		}
	}

	return nil
}

func (b *backend) pathTidyCancelWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if atomic.LoadUint32(b.tidyCASGuard) == 0 {
		resp := &logical.Response{}
		resp.AddWarning("Tidy operation cannot be cancelled as none is currently running.")
		return resp, nil
	}

	// Grab the status lock before writing the cancel atomic. This lets us
	// update the status correctly as well, avoiding writing it if we're not
	// presently running.
	//
	// Unlock needs to occur prior to calling read.
	b.tidyStatusLock.Lock()
	if b.tidyStatus.state == tidyStatusStarted || atomic.LoadUint32(b.tidyCASGuard) == 1 {
		if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 0, 1) {
			b.tidyStatus.state = tidyStatusCancelling
		}
	}
	b.tidyStatusLock.Unlock()

	return b.pathTidyStatusRead(ctx, req, d)
}

func (b *backend) pathTidyStatusRead(_ context.Context, _ *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	b.tidyStatusLock.RLock()
	defer b.tidyStatusLock.RUnlock()

	resp := &logical.Response{
		Data: map[string]interface{}{
			"safety_buffer":                         nil,
			"issuer_safety_buffer":                  nil,
			"tidy_cert_store":                       nil,
			"tidy_revoked_certs":                    nil,
			"tidy_revoked_cert_issuer_associations": nil,
			"tidy_expired_issuers":                  nil,
			"tidy_move_legacy_ca_bundle":            nil,
			"tidy_revocation_queue":                 nil,
			"tidy_cross_cluster_revoked_certs":      nil,
			"tidy_acme":                             nil,
			"tidy_cert_metadata":                    nil,
			"tidy_cmpv2_nonce_store":                nil,
			"pause_duration":                        nil,
			"state":                                 "Inactive",
			"error":                                 nil,
			"time_started":                          nil,
			"time_finished":                         nil,
			"message":                               nil,
			"cert_store_deleted_count":              nil,
			"revoked_cert_deleted_count":            nil,
			"missing_issuer_cert_count":             nil,
			"current_cert_store_count":              nil,
			"current_revoked_cert_count":            nil,
			"internal_backend_uuid":                 nil,
			"revocation_queue_deleted_count":        nil,
			"cross_revoked_cert_deleted_count":      nil,
			"total_acme_account_count":              nil,
			"acme_account_deleted_count":            nil,
			"acme_account_revoked_count":            nil,
			"acme_orders_deleted_count":             nil,
			"acme_account_safety_buffer":            nil,
			"cert_metadata_deleted_count":           nil,
			"cmpv2_nonce_deleted_count":             nil,
			"last_auto_tidy_finished":               b.getLastAutoTidyTimeWithoutLock(), // we acquired the tidyStatusLock above.
		},
	}

	resp.Data["internal_backend_uuid"] = b.backendUUID

	certCounter := b.GetCertificateCounter()
	if certCounter.IsEnabled() {
		resp.Data["current_cert_store_count"] = certCounter.CertificateCount()
		resp.Data["current_revoked_cert_count"] = certCounter.RevokedCount()
		if !certCounter.IsInitialized() {
			resp.AddWarning("Certificates in storage are still being counted, current counts provided may be " +
				"inaccurate")
		}
		certError := certCounter.Error()
		if certError != nil {
			resp.Data["certificate_counting_error"] = certError.Error()
		}
	}

	if b.tidyStatus.state == tidyStatusInactive {
		return resp, nil
	}

	resp.Data["safety_buffer"] = b.tidyStatus.safetyBuffer
	resp.Data["issuer_safety_buffer"] = b.tidyStatus.issuerSafetyBuffer
	resp.Data["tidy_cert_store"] = b.tidyStatus.tidyCertStore
	resp.Data["tidy_revoked_certs"] = b.tidyStatus.tidyRevokedCerts
	resp.Data["tidy_revoked_cert_issuer_associations"] = b.tidyStatus.tidyRevokedAssocs
	resp.Data["tidy_expired_issuers"] = b.tidyStatus.tidyExpiredIssuers
	resp.Data["tidy_move_legacy_ca_bundle"] = b.tidyStatus.tidyBackupBundle
	resp.Data["tidy_revocation_queue"] = b.tidyStatus.tidyRevocationQueue
	resp.Data["tidy_cross_cluster_revoked_certs"] = b.tidyStatus.tidyCrossRevokedCerts
	resp.Data["tidy_acme"] = b.tidyStatus.tidyAcme
	resp.Data["tidy_cert_metadata"] = b.tidyStatus.tidyCertMetadata
	resp.Data["tidy_cmpv2_nonce_store"] = b.tidyStatus.tidyCMPV2NonceStore
	resp.Data["pause_duration"] = b.tidyStatus.pauseDuration
	resp.Data["time_started"] = b.tidyStatus.timeStarted
	resp.Data["message"] = b.tidyStatus.message
	resp.Data["cert_store_deleted_count"] = b.tidyStatus.certStoreDeletedCount
	resp.Data["revoked_cert_deleted_count"] = b.tidyStatus.revokedCertDeletedCount
	resp.Data["missing_issuer_cert_count"] = b.tidyStatus.missingIssuerCertCount
	resp.Data["revocation_queue_deleted_count"] = b.tidyStatus.revQueueDeletedCount
	resp.Data["cross_revoked_cert_deleted_count"] = b.tidyStatus.crossRevokedDeletedCount
	resp.Data["revocation_queue_safety_buffer"] = b.tidyStatus.revQueueSafetyBuffer
	resp.Data["total_acme_account_count"] = b.tidyStatus.acmeAccountsCount
	resp.Data["acme_account_deleted_count"] = b.tidyStatus.acmeAccountsDeletedCount
	resp.Data["acme_account_revoked_count"] = b.tidyStatus.acmeAccountsRevokedCount
	resp.Data["acme_orders_deleted_count"] = b.tidyStatus.acmeOrdersDeletedCount
	resp.Data["acme_account_safety_buffer"] = b.tidyStatus.acmeAccountSafetyBuffer
	resp.Data["cert_metadata_deleted_count"] = b.tidyStatus.certMetadataDeletedCount
	resp.Data["cmpv2_nonce_deleted_count"] = b.tidyStatus.cmpv2NonceDeletedCount

	switch b.tidyStatus.state {
	case tidyStatusInactive:
		resp.Data["state"] = "Inactive"
	case tidyStatusStarted:
		resp.Data["state"] = "Running"
	case tidyStatusFinished:
		resp.Data["state"] = "Finished"
		resp.Data["time_finished"] = b.tidyStatus.timeFinished
		resp.Data["message"] = nil
	case tidyStatusError:
		resp.Data["state"] = "Error"
		resp.Data["time_finished"] = b.tidyStatus.timeFinished
		resp.Data["error"] = b.tidyStatus.err.Error()
		// Don't clear the message so that it serves as a hint about when
		// the error occurred.
	case tidyStatusCancelling:
		resp.Data["state"] = "Cancelling"
	case tidyStatusCancelled:
		resp.Data["state"] = "Cancelled"
		resp.Data["time_finished"] = b.tidyStatus.timeFinished
	}

	return resp, nil
}

func (b *backend) pathConfigAutoTidyRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)
	config, err := sc.getAutoTidyConfig()
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: getTidyConfigData(*config),
	}, nil
}

func (b *backend) pathConfigAutoTidyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)
	config, err := sc.getAutoTidyConfig()
	if err != nil {
		return nil, err
	}

	isAutoTidyBeingEnabled := false

	if enabledRaw, ok := d.GetOk("enabled"); ok {
		enabled, err := parseutil.ParseBool(enabledRaw)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to parse enabled flag as a boolean: %s", err.Error())), nil
		}
		if !config.Enabled && enabled {
			// we are turning on auto-tidy reset our persisted time to now
			isAutoTidyBeingEnabled = true
		}
		config.Enabled = enabled
	}

	if minStartupBackoffRaw, ok := d.GetOk("min_startup_backoff_duration"); ok {
		minDuration, err := parseutil.ParseDurationSecond(minStartupBackoffRaw)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to parse min_startup_backoff_duration flag as a duration: %s", err.Error())), nil
		}
		if minDuration.Seconds() < 1 {
			return logical.ErrorResponse(fmt.Sprintf("min_startup_backoff_duration must be at least 1 second: parsed: %v", minDuration)), nil
		}
		config.MinStartupBackoff = minDuration
	}

	if maxStartupBackoffRaw, ok := d.GetOk("max_startup_backoff_duration"); ok {
		maxDuration, err := parseutil.ParseDurationSecond(maxStartupBackoffRaw)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to parse max_startup_backoff_duration flag as a duration: %s", err.Error())), nil
		}
		if maxDuration.Seconds() < 1 {
			return logical.ErrorResponse(fmt.Sprintf("max_startup_backoff_duration must be at least 1 second: parsed: %v", maxDuration)), nil
		}
		config.MaxStartupBackoff = maxDuration
	}

	if config.MinStartupBackoff > config.MaxStartupBackoff {
		return logical.ErrorResponse(fmt.Sprintf("max_startup_backoff_duration %v must be greater or equal to min_startup_backoff_duration %v", config.MaxStartupBackoff, config.MinStartupBackoff)), nil
	}

	if intervalRaw, ok := d.GetOk("interval_duration"); ok {
		config.Interval = time.Duration(intervalRaw.(int)) * time.Second
		if config.Interval < 0 {
			return logical.ErrorResponse(fmt.Sprintf("given interval_duration must be greater than or equal to zero seconds; got: %v", intervalRaw)), nil
		}
	}

	if certStoreRaw, ok := d.GetOk("tidy_cert_store"); ok {
		config.CertStore = certStoreRaw.(bool)
	}

	if revokedCertsRaw, ok := d.GetOk("tidy_revoked_certs"); ok {
		config.RevokedCerts = revokedCertsRaw.(bool)
	}

	if issuerAssocRaw, ok := d.GetOk("tidy_revoked_cert_issuer_associations"); ok {
		config.IssuerAssocs = issuerAssocRaw.(bool)
	}

	if safetyBufferRaw, ok := d.GetOk("safety_buffer"); ok {
		config.SafetyBuffer = time.Duration(safetyBufferRaw.(int)) * time.Second
		if config.SafetyBuffer < 1*time.Second {
			return logical.ErrorResponse(fmt.Sprintf("given safety_buffer must be at least one second; got: %v", safetyBufferRaw)), nil
		}
	}

	if pauseDurationRaw, ok := d.GetOk("pause_duration"); ok {
		config.PauseDuration, err = parseutil.ParseDurationSecond(pauseDurationRaw.(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("unable to parse given pause_duration: %v", err)), nil
		}

		if config.PauseDuration < (0 * time.Second) {
			return logical.ErrorResponse("received invalid, negative pause_duration"), nil
		}
	}

	if expiredIssuers, ok := d.GetOk("tidy_expired_issuers"); ok {
		config.ExpiredIssuers = expiredIssuers.(bool)
	}

	if issuerSafetyBufferRaw, ok := d.GetOk("issuer_safety_buffer"); ok {
		config.IssuerSafetyBuffer = time.Duration(issuerSafetyBufferRaw.(int)) * time.Second
		if config.IssuerSafetyBuffer < 1*time.Second {
			return logical.ErrorResponse(fmt.Sprintf("given safety_buffer must be at least one second; got: %v", issuerSafetyBufferRaw)), nil
		}
	}

	if backupBundle, ok := d.GetOk("tidy_move_legacy_ca_bundle"); ok {
		config.BackupBundle = backupBundle.(bool)
	}

	if revocationQueueRaw, ok := d.GetOk("tidy_revocation_queue"); ok {
		config.RevocationQueue = revocationQueueRaw.(bool)
	}

	if queueSafetyBufferRaw, ok := d.GetOk("revocation_queue_safety_buffer"); ok {
		config.QueueSafetyBuffer = time.Duration(queueSafetyBufferRaw.(int)) * time.Second
		if config.QueueSafetyBuffer < 1*time.Second {
			return logical.ErrorResponse(fmt.Sprintf("given revocation_queue_safety_buffer must be at least one second; got: %v", queueSafetyBufferRaw)), nil
		}
	}

	if crossRevokedRaw, ok := d.GetOk("tidy_cross_cluster_revoked_certs"); ok {
		config.CrossRevokedCerts = crossRevokedRaw.(bool)
	}

	if tidyAcmeRaw, ok := d.GetOk("tidy_acme"); ok {
		config.TidyAcme = tidyAcmeRaw.(bool)
	}

	if acmeAccountSafetyBufferRaw, ok := d.GetOk("acme_account_safety_buffer"); ok {
		config.AcmeAccountSafetyBuffer = time.Duration(acmeAccountSafetyBufferRaw.(int)) * time.Second
		if config.AcmeAccountSafetyBuffer < 1*time.Second {
			return logical.ErrorResponse(fmt.Sprintf("given acme_account_safety_buffer must be at least one second; got: %v", acmeAccountSafetyBufferRaw)), nil
		}
	}

	if tidyCertMetadataRaw, ok := d.GetOk("tidy_cert_metadata"); ok {
		config.CertMetadata = tidyCertMetadataRaw.(bool)

		if config.CertMetadata && !constants.IsEnterprise {
			return logical.ErrorResponse("certificate metadata is only supported on Vault Enterprise"), nil
		}
	}

	if config.Enabled && !config.IsAnyTidyEnabled() {
		return logical.ErrorResponse("Auto-tidy enabled but no tidy operations were requested. Enable at least one tidy operation to be run (" + config.AnyTidyConfig() + ")."), nil
	}

	if maintainCountEnabledRaw, ok := d.GetOk("maintain_stored_certificate_counts"); ok {
		config.MaintainCount = maintainCountEnabledRaw.(bool)
	}

	if runningStorageMetricsEnabledRaw, ok := d.GetOk("publish_stored_certificate_count_metrics"); ok {
		config.PublishMetrics = runningStorageMetricsEnabledRaw.(bool)
	}

	if config.PublishMetrics && !config.MaintainCount {
		return logical.ErrorResponse("Can not publish a running storage metrics count to metrics without first maintaining that count.  Enable `maintain_stored_certificate_counts` to enable `publish_stored_certificate_count_metrics`."), nil
	}

	if err := sc.writeAutoTidyConfig(config); err != nil {
		return nil, err
	}

	if isAutoTidyBeingEnabled {
		if err := b.updateLastAutoTidyTime(sc, time.Now()); err != nil {
			b.Logger().Warn("failed to update last auto tidy run time to now, the first auto-tidy "+
				"might run soon and not at the next delay provided", "error", err.Error())
		}
	}

	return &logical.Response{
		Data: getTidyConfigData(*config),
	}, nil
}

func (b *backend) tidyStatusStart(config *tidyConfig) {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus = &tidyStatus{
		safetyBuffer:            int(config.SafetyBuffer / time.Second),
		issuerSafetyBuffer:      int(config.IssuerSafetyBuffer / time.Second),
		revQueueSafetyBuffer:    int(config.QueueSafetyBuffer / time.Second),
		acmeAccountSafetyBuffer: int(config.AcmeAccountSafetyBuffer / time.Second),
		tidyCertStore:           config.CertStore,
		tidyRevokedCerts:        config.RevokedCerts,
		tidyRevokedAssocs:       config.IssuerAssocs,
		tidyExpiredIssuers:      config.ExpiredIssuers,
		tidyBackupBundle:        config.BackupBundle,
		tidyRevocationQueue:     config.RevocationQueue,
		tidyCrossRevokedCerts:   config.CrossRevokedCerts,
		tidyAcme:                config.TidyAcme,
		tidyCertMetadata:        config.CertMetadata,
		pauseDuration:           config.PauseDuration.String(),

		state:       tidyStatusStarted,
		timeStarted: time.Now(),
	}

	metrics.SetGauge([]string{"secrets", "pki", "tidy", "start_time_epoch"}, float32(b.tidyStatus.timeStarted.Unix()))
}

func (b *backend) tidyStatusStop(err error) {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.timeFinished = time.Now()
	b.tidyStatus.err = err
	if err == nil {
		b.tidyStatus.state = tidyStatusFinished
	} else if errors.Is(err, tidyCancelledError) {
		b.tidyStatus.state = tidyStatusCancelled
	} else {
		b.tidyStatus.state = tidyStatusError
	}

	metrics.MeasureSince([]string{"secrets", "pki", "tidy", "duration"}, b.tidyStatus.timeStarted)
	metrics.SetGauge([]string{"secrets", "pki", "tidy", "start_time_epoch"}, 0)
	metrics.IncrCounter([]string{"secrets", "pki", "tidy", "cert_store_deleted_count"}, float32(b.tidyStatus.certStoreDeletedCount))
	metrics.IncrCounter([]string{"secrets", "pki", "tidy", "revoked_cert_deleted_count"}, float32(b.tidyStatus.revokedCertDeletedCount))

	if err != nil {
		metrics.IncrCounter([]string{"secrets", "pki", "tidy", "failure"}, 1)
	} else {
		metrics.IncrCounter([]string{"secrets", "pki", "tidy", "success"}, 1)
	}
}

func (b *backend) tidyStatusMessage(msg string) {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.message = msg
}

func (b *backend) tidyStatusIncCertStoreCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.certStoreDeletedCount++

	b.GetCertificateCounter().DecrementTotalCertificatesCountReport()
}

func (b *backend) tidyStatusIncRevokedCertCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.revokedCertDeletedCount++

	b.GetCertificateCounter().DecrementTotalRevokedCertificatesCountReport()
}

func (b *backend) tidyStatusIncMissingIssuerCertCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.missingIssuerCertCount++
}

func (b *backend) tidyStatusIncRevQueueCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.revQueueDeletedCount++
}

func (b *backend) tidyStatusIncCrossRevCertCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.crossRevokedDeletedCount++
}

func (b *backend) tidyStatusIncRevAcmeAccountCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.acmeAccountsRevokedCount++
}

func (b *backend) tidyStatusIncDeletedAcmeAccountCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.acmeAccountsDeletedCount++
}

func (b *backend) tidyStatusIncDelAcmeOrderCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.acmeOrdersDeletedCount++
}

func (b *backend) tidyStatusIncCertMetadataCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.certMetadataDeletedCount++
}

func (b *backend) tidyStatusIncCMPV2NonceDeletedCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.cmpv2NonceDeletedCount++
}

// updateLastAutoTidyTime should be used to update b.lastAutoTidy as the required locks
// are acquired and the auto tidy time is persisted to storage to work across restarts
func (b *backend) updateLastAutoTidyTime(sc *storageContext, lastRunTime time.Time) error {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.lastAutoTidy = lastRunTime
	return sc.writeAutoTidyLastRun(lastRunTime)
}

// getLastAutoTidyTime should be used to read from b.lastAutoTidy as the required locks
// are acquired prior to reading
func (b *backend) getLastAutoTidyTime() time.Time {
	b.tidyStatusLock.RLock()
	defer b.tidyStatusLock.RUnlock()
	return b.getLastAutoTidyTimeWithoutLock()
}

// getLastAutoTidyTimeWithoutLock should be used to read from b.lastAutoTidy with the
// b.tidyStatusLock being acquired, normally use getLastAutoTidyTime
func (b *backend) getLastAutoTidyTimeWithoutLock() time.Time {
	return b.lastAutoTidy
}

const pathTidyHelpSyn = `
Tidy up the backend by removing expired certificates, revocation information,
or both.
`

const pathTidyHelpDesc = `
This endpoint allows expired certificates and/or revocation information to be
removed from the backend, freeing up storage and shortening CRLs.

For safety, this function is a noop if called without parameters; cleanup from
normal certificate storage must be enabled with 'tidy_cert_store' and cleanup
from revocation information must be enabled with 'tidy_revocation_list'.

The 'safety_buffer' parameter is useful to ensure that clock skew amongst your
hosts cannot lead to a certificate being removed from the CRL while it is still
considered valid by other hosts (for instance, if their clocks are a few
minutes behind). The 'safety_buffer' parameter can be an integer number of
seconds or a string duration like "72h".

All certificates and/or revocation information currently stored in the backend
will be checked when this endpoint is hit. The expiration of the
certificate/revocation information of each certificate being held in
certificate storage or in revocation information will then be checked. If the
current time, minus the value of 'safety_buffer', is greater than the
expiration, it will be removed.
`

const pathTidyCancelHelpSyn = `
Cancels a currently running tidy operation.
`

const pathTidyCancelHelpDesc = `
This endpoint allows cancelling a currently running tidy operation.

Periodically throughout the invocation of tidy, we'll check if the operation
has been requested to be cancelled. If so, we'll stop the currently running
tidy operation.
`

const pathTidyStatusHelpSyn = `
Returns the status of the tidy operation.
`

const pathTidyStatusHelpDesc = `
This is a read only endpoint that returns information about the current tidy
operation, or the most recent if none is currently running.

The result includes the following fields:
* 'safety_buffer': the value of this parameter when initiating the tidy operation
* 'tidy_cert_store': the value of this parameter when initiating the tidy operation
* 'tidy_revoked_certs': the value of this parameter when initiating the tidy operation
* 'tidy_revoked_cert_issuer_associations': the value of this parameter when initiating the tidy operation
* 'state': one of "Inactive", "Running", "Finished", "Error"
* 'error': the error message, if the operation ran into an error
* 'time_started': the time the operation started
* 'time_finished': the time the operation finished
* 'message': One of "Tidying certificate store: checking entry N of TOTAL" or
  "Tidying revoked certificates: checking certificate N of TOTAL"
* 'cert_store_deleted_count': The number of certificate storage entries deleted
* 'revoked_cert_deleted_count': The number of revoked certificate entries deleted
* 'missing_issuer_cert_count': The number of revoked certificates which were missing a valid issuer reference
* 'tidy_expired_issuers': the value of this parameter when initiating the tidy operation
* 'issuer_safety_buffer': the value of this parameter when initiating the tidy operation
* 'tidy_move_legacy_ca_bundle': the value of this parameter when initiating the tidy operation
* 'tidy_revocation_queue': the value of this parameter when initiating the tidy operation
* 'revocation_queue_deleted_count': the number of revocation queue entries deleted
* 'tidy_cross_cluster_revoked_certs': the value of this parameter when initiating the tidy operation
* 'cross_revoked_cert_deleted_count': the number of cross-cluster revoked certificate entries deleted
* 'revocation_queue_safety_buffer': the value of this parameter when initiating the tidy operation
* 'tidy_acme': the value of this parameter when initiating the tidy operation
* 'acme_account_safety_buffer': the value of this parameter when initiating the tidy operation
* 'total_acme_account_count': the total number of acme accounts in the list to be iterated over
* 'acme_account_deleted_count': the number of revoked acme accounts deleted during the operation
* 'acme_account_revoked_count': the number of acme accounts revoked during the operation
* 'acme_orders_deleted_count': the number of acme orders deleted during the operation
`

const pathConfigAutoTidySyn = `
Modifies the current configuration for automatic tidy execution.
`

const pathConfigAutoTidyDesc = `
This endpoint accepts parameters to a tidy operation (see /tidy) that
will be used for automatic tidy execution. This takes two extra parameters,
enabled (to enable or disable auto-tidy) and interval_duration (which
controls the frequency of auto-tidy execution).

Once enabled, a tidy operation will be kicked off automatically, as if it
were executed with the posted configuration.
`

func getTidyConfigData(config tidyConfig) map[string]interface{} {
	return map[string]interface{}{
		// This map is in the same order as tidyConfig to ensure that all fields are accounted for
		"enabled":                                  config.Enabled,
		"interval_duration":                        int(config.Interval / time.Second),
		"min_startup_backoff_duration":             int(config.MinStartupBackoff.Seconds()),
		"max_startup_backoff_duration":             int(config.MaxStartupBackoff.Seconds()),
		"tidy_cert_store":                          config.CertStore,
		"tidy_revoked_certs":                       config.RevokedCerts,
		"tidy_revoked_cert_issuer_associations":    config.IssuerAssocs,
		"tidy_expired_issuers":                     config.ExpiredIssuers,
		"tidy_move_legacy_ca_bundle":               config.BackupBundle,
		"tidy_acme":                                config.TidyAcme,
		"safety_buffer":                            int(config.SafetyBuffer / time.Second),
		"issuer_safety_buffer":                     int(config.IssuerSafetyBuffer / time.Second),
		"acme_account_safety_buffer":               int(config.AcmeAccountSafetyBuffer / time.Second),
		"pause_duration":                           config.PauseDuration.String(),
		"publish_stored_certificate_count_metrics": config.PublishMetrics,
		"maintain_stored_certificate_counts":       config.MaintainCount,
		"tidy_revocation_queue":                    config.RevocationQueue,
		"revocation_queue_safety_buffer":           int(config.QueueSafetyBuffer / time.Second),
		"tidy_cross_cluster_revoked_certs":         config.CrossRevokedCerts,
		"tidy_cert_metadata":                       config.CertMetadata,
		"tidy_cmpv2_nonce_store":                   config.CMPV2NonceStore,
	}
}
