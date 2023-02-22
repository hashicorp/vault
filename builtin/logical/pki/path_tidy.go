package pki

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

var tidyCancelledError = errors.New("tidy operation cancelled")

type tidyStatusState int

const (
	tidyStatusInactive   tidyStatusState = iota
	tidyStatusStarted                    = iota
	tidyStatusFinished                   = iota
	tidyStatusError                      = iota
	tidyStatusCancelling                 = iota
	tidyStatusCancelled                  = iota
)

type tidyStatus struct {
	// Parameters used to initiate the operation
	safetyBuffer          int
	issuerSafetyBuffer    int
	tidyCertStore         bool
	tidyRevokedCerts      bool
	tidyRevokedAssocs     bool
	tidyExpiredIssuers    bool
	tidyBackupBundle      bool
	tidyRevocationQueue   bool
	tidyCrossRevokedCerts bool
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
}

type tidyConfig struct {
	Enabled            bool          `json:"enabled"`
	Interval           time.Duration `json:"interval_duration"`
	CertStore          bool          `json:"tidy_cert_store"`
	RevokedCerts       bool          `json:"tidy_revoked_certs"`
	IssuerAssocs       bool          `json:"tidy_revoked_cert_issuer_associations"`
	ExpiredIssuers     bool          `json:"tidy_expired_issuers"`
	BackupBundle       bool          `json:"tidy_move_legacy_ca_bundle"`
	SafetyBuffer       time.Duration `json:"safety_buffer"`
	IssuerSafetyBuffer time.Duration `json:"issuer_safety_buffer"`
	PauseDuration      time.Duration `json:"pause_duration"`
	MaintainCount      bool          `json:"maintain_stored_certificate_counts"`
	PublishMetrics     bool          `json:"publish_stored_certificate_count_metrics"`
	RevocationQueue    bool          `json:"tidy_revocation_queue"`
	QueueSafetyBuffer  time.Duration `json:"revocation_queue_safety_buffer"`
	CrossRevokedCerts  bool          `json:"tidy_cross_cluster_revoked_certs"`
}

var defaultTidyConfig = tidyConfig{
	Enabled:            false,
	Interval:           12 * time.Hour,
	CertStore:          false,
	RevokedCerts:       false,
	IssuerAssocs:       false,
	ExpiredIssuers:     false,
	BackupBundle:       false,
	SafetyBuffer:       72 * time.Hour,
	IssuerSafetyBuffer: 365 * 24 * time.Hour,
	PauseDuration:      0 * time.Second,
	MaintainCount:      false,
	PublishMetrics:     false,
	RevocationQueue:    false,
	QueueSafetyBuffer:  48 * time.Hour,
	CrossRevokedCerts:  false,
}

func pathTidy(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "tidy$",
		Fields:  addTidyFields(map[string]*framework.FieldSchema{}),
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathTidyWrite,
				Responses: map[int][]framework.Response{
					http.StatusAccepted: {{
						Description: "Accepted",
						Fields:      map[string]*framework.FieldSchema{
							"http_content_type": {
								Type: framework.TypeString,
								Required: true,
							},
							"http_raw_body": {
								Type: framework.TypeString,
								Required: true,
							},
							"http_status_code": {
								Type: framework.TypeString,
								Required: true,
							},
						},
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
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathTidyCancelWrite,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"safety_buffer": {
								Type:        framework.TypeInt,
								Description: `Safety buffer time duration`,
								Required:    false,
							},
							"issuer_safety_buffer": {
								Type:        framework.TypeInt,
								Description: `Issuer safety buffer`,
								Required:    false,
							},
							"tidy_cert_store": {
								Type:        framework.TypeBool,
								Description: `Tidy certificate store`,
								Required:    false,
							},
							"tidy_revoked_certs": {
								Type:        framework.TypeBool,
								Description: `Tidy revoked certificates`,
								Required:    false,
							},
							"tidy_revoked_cert_issuer_associations": {
								Type:        framework.TypeBool,
								Description: `Tidy revoked certificate issuer associations`,
								Required:    false,
							},
							"tidy_expired_issuers": {
								Type:        framework.TypeBool,
								Description: `Tidy expired issuers`,
								Required:    false,
							},
							"pause_duration": {
								Type:        framework.TypeString,
								Description: `Duration to pause between tidying certificates`,
								Required:    false,
							},
							"state": {
								Type:        framework.TypeString,
								Description: `One of Inactive, Running, Finished, or Error`,
								Required:    false,
							},
							"error": {
								Type:        framework.TypeString,
								Description: `The error message`,
								Required:    false,
							},
							"time_started": {
								Type:        framework.TypeString,
								Description: `Time the operation started`,
								Required:    false,
							},
							"time_finished": {
								Type:        framework.TypeString,
								Description: `Time the operation finished`,
								Required:    false,
							},
							"message": {
								Type:        framework.TypeString,
								Description: `Message of the operation`,
								Required:    false,
							},
							"cert_store_deleted_count": {
								Type:        framework.TypeInt,
								Description: `The number of certificate storage entries deleted`,
								Required:    false,
							},
							"revoked_cert_deleted_count": {
								Type:        framework.TypeInt,
								Description: `The number of revoked certificate entries deleted`,
								Required:    false,
							},
							"current_cert_store_count": {
								Type:        framework.TypeInt,
								Description: `The number of revoked certificate entries deleted`,
								Required:    false,
							},
							"current_revoked_cert_count": {
								Type:        framework.TypeInt,
								Description: `The number of revoked certificate entries deleted`,
								Required:    false,
							},
							"missing_issuer_cert_count": {
								Type:     framework.TypeInt,
								Required: false,
							},
							"tidy_move_legacy_ca_bundle": {
								Type:     framework.TypeBool,
								Required: false,
							},
							"tidy_cross_cluster_revoked_certs": {
								Type:     framework.TypeBool,
								Required: false,
							},
							"tidy_revocation_queue": {
								Type:     framework.TypeBool,
								Required: false,
							},
							"revocation_queue_deleted_count": {
								Type:     framework.TypeInt,
								Required: false,
							},
							"cross_revoked_cert_deleted_count": {
								Type:     framework.TypeInt,
								Required: false,
							},
							"internal_backend_uuid": {
								Type:     framework.TypeString,
								Required: false,
							},
						},
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
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathTidyStatusRead,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
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
								Type:        framework.TypeString,
								Description: ``,
								Required:    false,
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
						},
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
	return &framework.Path{
		Pattern: "config/auto-tidy",
		Fields: addTidyFields(map[string]*framework.FieldSchema{
			"enabled": {
				Type:        framework.TypeBool,
				Description: `Set to true to enable automatic tidy operations.`,
			},
			"interval_duration": {
				Type:        framework.TypeDurationSecond,
				Description: `Interval at which to run an auto-tidy operation. This is the time between tidy invocations (after one finishes to the start of the next). Running a manual tidy will reset this duration.`,
				Default:     int(defaultTidyConfig.Interval / time.Second), // TypeDurationSecond currently requires the default to be an int.
			},
		}),
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigAutoTidyRead,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"enabled": {
								Type:        framework.TypeBool,
								Description: `Specifies whether automatic tidy is enabled or not`,
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
							"pause_duration": {
								Type:        framework.TypeString,
								Description: `Duration to pause between tidying certificates`,
								Required:    true,
							},
							"tidy_move_legacy_ca_bundle": {
								Type:     framework.TypeBool,
								Required: true,
							},
							"tidy_cross_cluster_revoked_certs": {
								Type:     framework.TypeBool,
								Required: true,
							},
							"tidy_revocation_queue": {
								Type:     framework.TypeBool,
								Required: true,
							},
							"revocation_queue_safety_buffer": {
								Type:     framework.TypeDurationSecond,
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
						},
					}},
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigAutoTidyWrite,
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
						Fields: map[string]*framework.FieldSchema{
							"enabled": {
								Type:        framework.TypeBool,
								Description: `Specifies whether automatic tidy is enabled or not`,
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
							"pause_duration": {
								Type:        framework.TypeString,
								Description: `Duration to pause between tidying certificates`,
								Required:    true,
							},
							"tidy_cross_cluster_revoked_certs": {
								Type:     framework.TypeBool,
								Required: true,
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
								Type:     framework.TypeDurationSecond,
								Required: true,
							},
						},
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

	if safetyBuffer < 1 {
		return logical.ErrorResponse("safety_buffer must be greater than zero"), nil
	}

	if issuerSafetyBuffer < 1 {
		return logical.ErrorResponse("issuer_safety_buffer must be greater than zero"), nil
	}

	if queueSafetyBuffer < 1 {
		return logical.ErrorResponse("revocation_queue_safety_buffer must be greater than zero"), nil
	}

	if pauseDurationStr != "" {
		var err error
		pauseDuration, err = time.ParseDuration(pauseDurationStr)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("Error parsing pause_duration: %v", err)), nil
		}

		if pauseDuration < (0 * time.Second) {
			return logical.ErrorResponse("received invalid, negative pause_duration"), nil
		}
	}

	bufferDuration := time.Duration(safetyBuffer) * time.Second
	issuerBufferDuration := time.Duration(issuerSafetyBuffer) * time.Second
	queueSafetyBufferDuration := time.Duration(queueSafetyBuffer) * time.Second

	// Manual run with constructed configuration.
	config := &tidyConfig{
		Enabled:            true,
		Interval:           0 * time.Second,
		CertStore:          tidyCertStore,
		RevokedCerts:       tidyRevokedCerts,
		IssuerAssocs:       tidyRevokedAssocs,
		ExpiredIssuers:     tidyExpiredIssuers,
		BackupBundle:       tidyBackupBundle,
		SafetyBuffer:       bufferDuration,
		IssuerSafetyBuffer: issuerBufferDuration,
		PauseDuration:      pauseDuration,
		RevocationQueue:    tidyRevocationQueue,
		QueueSafetyBuffer:  queueSafetyBufferDuration,
		CrossRevokedCerts:  tidyCrossRevokedCerts,
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

	// Mark the last tidy operation as relatively recent, to ensure we don't
	// try to trigger the periodic function.
	b.tidyStatusLock.Lock()
	b.lastTidy = time.Now()
	b.tidyStatusLock.Unlock()

	// Kick off the actual tidy.
	b.startTidyOperation(req, config)

	resp := &logical.Response{}
	if !tidyCertStore && !tidyRevokedCerts && !tidyRevokedAssocs && !tidyExpiredIssuers && !tidyBackupBundle && !tidyRevocationQueue && !tidyCrossRevokedCerts {
		resp.AddWarning("No targets to tidy; specify tidy_cert_store=true or tidy_revoked_certs=true or tidy_revoked_cert_issuer_associations=true or tidy_expired_issuers=true or tidy_move_legacy_ca_bundle=true or tidy_revocation_queue=true or tidy_cross_cluster_revoked_certs=true to start a tidy operation.")
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
			b.tidyStatusLock.Lock()
			b.lastTidy = time.Now()
			b.tidyStatusLock.Unlock()
		}
	}()
}

func (b *backend) doTidyCertStore(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	serials, err := req.Storage.List(ctx, "certs/")
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

		certEntry, err := req.Storage.Get(ctx, "certs/"+serial)
		if err != nil {
			return fmt.Errorf("error fetching certificate %q: %w", serial, err)
		}

		if certEntry == nil {
			logger.Warn("certificate entry is nil; tidying up since it is no longer useful for any server operations", "serial", serial)
			if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
				return fmt.Errorf("error deleting nil entry with serial %s: %w", serial, err)
			}
			b.tidyStatusIncCertStoreCount()
			continue
		}

		if certEntry.Value == nil || len(certEntry.Value) == 0 {
			logger.Warn("certificate entry has no value; tidying up since it is no longer useful for any server operations", "serial", serial)
			if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
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
			if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
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
	b.revokeStorageLock.Lock()
	defer b.revokeStorageLock.Unlock()

	// Fetch and parse our issuers so we can associate them if necessary.
	sc := b.makeStorageContext(ctx, req.Storage)
	issuerIDCertMap, err := fetchIssuerMapForRevocationChecking(sc)
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

	var revInfo revocationInfo
	for i, serial := range revokedSerials {
		b.tidyStatusMessage(fmt.Sprintf("Tidying revoked certificates: checking certificate %d of %d", i, len(revokedSerials)))
		metrics.SetGauge([]string{"secrets", "pki", "tidy", "revoked_cert_current_entry"}, float32(i))

		// Check for cancel before continuing.
		if atomic.CompareAndSwapUint32(b.tidyCancelCAS, 1, 0) {
			return tidyCancelledError
		}

		// Check for pause duration to reduce resource consumption.
		if config.PauseDuration > (0 * time.Second) {
			b.revokeStorageLock.Unlock()
			time.Sleep(config.PauseDuration)
			b.revokeStorageLock.Lock()
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
				revInfo.CertificateIssuer = issuerID("")
				storeCert = true
				if associateRevokedCertWithIsssuer(&revInfo, revokedCert, issuerIDCertMap) {
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
				if err := req.Storage.Delete(ctx, "certs/"+serial); err != nil {
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
			if err := b.crlBuilder.rebuild(sc, false); err != nil {
				return err
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
	if b.useLegacyBundleCaStorage() {
		return nil
	}

	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	// Fetch and parse our issuers so we have their expiration date.
	sc := b.makeStorageContext(ctx, req.Storage)
	issuerIDCertMap, err := fetchIssuerMapForRevocationChecking(sc)
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
		b.revokeStorageLock.Lock()
		defer b.revokeStorageLock.Unlock()

		if err := b.crlBuilder.rebuild(sc, false); err != nil {
			return err
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
	if b.useLegacyBundleCaStorage() {
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
	b.revokeStorageLock.Lock()
	defer b.revokeStorageLock.Unlock()

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
				b.revokeStorageLock.Unlock()
				time.Sleep(config.PauseDuration)
				b.revokeStorageLock.Lock()
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
	b.revokeStorageLock.Lock()
	defer b.revokeStorageLock.Unlock()

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
				b.revokeStorageLock.Unlock()
				time.Sleep(config.PauseDuration)
				b.revokeStorageLock.Lock()
			}

			ePath := cPath + serial
			entry, err := sc.Storage.Get(sc.Context, ePath)
			if err != nil {
				return fmt.Errorf("error reading cross-cluster revocation entry (%v) to tidy: %w", ePath, err)
			}
			if entry == nil || entry.Value == nil {
				continue
			}

			var details unifiedRevocationEntry
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
		},
	}

	resp.Data["internal_backend_uuid"] = b.backendUUID

	if b.certCountEnabled.Load() {
		resp.Data["current_cert_store_count"] = b.certCount.Load()
		resp.Data["current_revoked_cert_count"] = b.revokedCertCount.Load()
		if !b.certsCounted.Load() {
			resp.AddWarning("Certificates in storage are still being counted, current counts provided may be " +
				"inaccurate")
		}
		if b.certCountError != "" {
			resp.Data["certificate_counting_error"] = b.certCountError
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
	resp.Data["pause_duration"] = b.tidyStatus.pauseDuration
	resp.Data["time_started"] = b.tidyStatus.timeStarted
	resp.Data["message"] = b.tidyStatus.message
	resp.Data["cert_store_deleted_count"] = b.tidyStatus.certStoreDeletedCount
	resp.Data["revoked_cert_deleted_count"] = b.tidyStatus.revokedCertDeletedCount
	resp.Data["missing_issuer_cert_count"] = b.tidyStatus.missingIssuerCertCount
	resp.Data["revocation_queue_deleted_count"] = b.tidyStatus.revQueueDeletedCount
	resp.Data["cross_revoked_cert_deleted_count"] = b.tidyStatus.crossRevokedDeletedCount

	switch b.tidyStatus.state {
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
		Data: map[string]interface{}{
			"enabled":                                  config.Enabled,
			"interval_duration":                        int(config.Interval / time.Second),
			"tidy_cert_store":                          config.CertStore,
			"tidy_revoked_certs":                       config.RevokedCerts,
			"tidy_revoked_cert_issuer_associations":    config.IssuerAssocs,
			"tidy_expired_issuers":                     config.ExpiredIssuers,
			"safety_buffer":                            int(config.SafetyBuffer / time.Second),
			"issuer_safety_buffer":                     int(config.IssuerSafetyBuffer / time.Second),
			"pause_duration":                           config.PauseDuration.String(),
			"publish_stored_certificate_count_metrics": config.PublishMetrics,
			"maintain_stored_certificate_counts":       config.MaintainCount,
			"tidy_move_legacy_ca_bundle":               config.BackupBundle,
			"tidy_revocation_queue":                    config.RevocationQueue,
			"revocation_queue_safety_buffer":           int(config.QueueSafetyBuffer / time.Second),
			"tidy_cross_cluster_revoked_certs":         config.CrossRevokedCerts,
		},
	}, nil
}

func (b *backend) pathConfigAutoTidyWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	sc := b.makeStorageContext(ctx, req.Storage)
	config, err := sc.getAutoTidyConfig()
	if err != nil {
		return nil, err
	}

	if enabledRaw, ok := d.GetOk("enabled"); ok {
		config.Enabled = enabledRaw.(bool)
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
		config.PauseDuration, err = time.ParseDuration(pauseDurationRaw.(string))
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

	if config.Enabled && !(config.CertStore || config.RevokedCerts || config.IssuerAssocs || config.ExpiredIssuers || config.BackupBundle || config.RevocationQueue || config.CrossRevokedCerts) {
		return logical.ErrorResponse("Auto-tidy enabled but no tidy operations were requested. Enable at least one tidy operation to be run (tidy_cert_store / tidy_revoked_certs / tidy_revoked_cert_issuer_associations / tidy_expired_issuers / tidy_move_legacy_ca_bundle / tidy_revocation_queue / tidy_cross_cluster_revoked_certs)."), nil
	}

	if maintainCountEnabledRaw, ok := d.GetOk("maintain_stored_certificate_counts"); ok {
		config.MaintainCount = maintainCountEnabledRaw.(bool)
	}

	if runningStorageMetricsEnabledRaw, ok := d.GetOk("publish_stored_certificate_count_metrics"); ok {
		if config.MaintainCount == false {
			return logical.ErrorResponse("Can not publish a running storage metrics count to metrics without first maintaining that count.  Enable `maintain_stored_certificate_counts` to enable `publish_stored_certificate_count_metrics."), nil
		}
		config.PublishMetrics = runningStorageMetricsEnabledRaw.(bool)
	}

	if err := sc.writeAutoTidyConfig(config); err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"enabled":                               config.Enabled,
			"interval_duration":                     int(config.Interval / time.Second),
			"tidy_cert_store":                       config.CertStore,
			"tidy_revoked_certs":                    config.RevokedCerts,
			"tidy_revoked_cert_issuer_associations": config.IssuerAssocs,
			"tidy_expired_issuers":                  config.ExpiredIssuers,
			"tidy_move_legacy_ca_bundle":            config.BackupBundle,
			"safety_buffer":                         int(config.SafetyBuffer / time.Second),
			"issuer_safety_buffer":                  int(config.IssuerSafetyBuffer / time.Second),
			"pause_duration":                        config.PauseDuration.String(),
			"tidy_revocation_queue":                 config.RevocationQueue,
			"revocation_queue_safety_buffer":        int(config.QueueSafetyBuffer / time.Second),
			"tidy_cross_cluster_revoked_certs":      config.CrossRevokedCerts,
		},
	}, nil
}

func (b *backend) tidyStatusStart(config *tidyConfig) {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus = &tidyStatus{
		safetyBuffer:          int(config.SafetyBuffer / time.Second),
		issuerSafetyBuffer:    int(config.IssuerSafetyBuffer / time.Second),
		tidyCertStore:         config.CertStore,
		tidyRevokedCerts:      config.RevokedCerts,
		tidyRevokedAssocs:     config.IssuerAssocs,
		tidyExpiredIssuers:    config.ExpiredIssuers,
		tidyBackupBundle:      config.BackupBundle,
		tidyRevocationQueue:   config.RevocationQueue,
		tidyCrossRevokedCerts: config.CrossRevokedCerts,
		pauseDuration:         config.PauseDuration.String(),

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
	} else if err == tidyCancelledError {
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

	b.ifCountEnabledDecrementTotalCertificatesCountReport()
}

func (b *backend) tidyStatusIncRevokedCertCount() {
	b.tidyStatusLock.Lock()
	defer b.tidyStatusLock.Unlock()

	b.tidyStatus.revokedCertDeletedCount++

	b.ifCountEnabledDecrementTotalRevokedCertificatesCountReport()
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
