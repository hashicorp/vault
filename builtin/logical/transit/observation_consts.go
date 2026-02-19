// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package transit

const (

	// key type observations

	// ObservationTypeTransitKeyRotateSuccess is emitted when a key is successfully rotated.
	// Metadata: key_name, type, derived, deletion_allowed, min_available_version,
	// min_decryption_version, min_encryption_version, latest_version, exportable,
	// allow_plaintext_backup, auto_rotate_period, imported_key, kdf (if derived),
	// kdf_mode (if derived), convergent_encryption (if derived), managed_key_id (if managed key)
	ObservationTypeTransitKeyRotateSuccess = "transit/key/rotate/success"

	// ObservationTypeTransitKeyRotateFail is emitted when a key rotation fails.
	// Metadata: key_name, type, derived, deletion_allowed, min_available_version,
	// min_decryption_version, min_encryption_version, latest_version, exportable,
	// allow_plaintext_backup, auto_rotate_period, imported_key, kdf (if derived),
	// kdf_mode (if derived), convergent_encryption (if derived), managed_key_id (if managed key)
	ObservationTypeTransitKeyRotateFail = "transit/key/rotate/success"

	// ObservationTypeTransitKeyWrite is emitted when a new key is created.
	// Metadata: key_name, type, derived, deletion_allowed, min_available_version,
	// min_decryption_version, min_encryption_version, latest_version, exportable,
	// allow_plaintext_backup, auto_rotate_period, imported_key, kdf (if derived),
	// kdf_mode (if derived), convergent_encryption (if derived), managed_key_id (if managed key)
	ObservationTypeTransitKeyWrite = "transit/key/write"

	// ObservationTypeTransitKeyRead is emitted when a key is read.
	// Metadata: key_name, type, derived, deletion_allowed, min_available_version,
	// min_decryption_version, min_encryption_version, latest_version, exportable,
	// allow_plaintext_backup, auto_rotate_period, imported_key, kdf (if derived),
	// kdf_mode (if derived), convergent_encryption (if derived)
	ObservationTypeTransitKeyRead = "transit/key/read"

	// ObservationTypeTransitKeyDelete is emitted when a key is deleted.
	// Metadata: key_name
	ObservationTypeTransitKeyDelete = "transit/key/delete"

	// ObservationTypeTransitKeyImport is emitted when a key is imported.
	// For new key imports, metadata includes: key_name, type, derived, exportable,
	// allow_plaintext_backup, auto_rotate_period
	// For version imports, metadata includes: key_name, type, derived, deletion_allowed,
	// min_available_version, min_decryption_version, min_encryption_version, latest_version,
	// exportable, allow_plaintext_backup, auto_rotate_period, imported_key, kdf (if derived),
	// kdf_mode (if derived), convergent_encryption (if derived), import_version
	ObservationTypeTransitKeyImport = "transit/key/import"

	// ObservationTypeTransitKeyExport is emitted when a key is exported.
	// For full key exports, metadata includes: key_name, type, derived, deletion_allowed,
	// min_available_version, min_decryption_version, min_encryption_version, latest_version,
	// exportable, allow_plaintext_backup, auto_rotate_period, imported_key, kdf (if derived),
	// kdf_mode (if derived), convergent_encryption (if derived)
	// For single version exports, metadata also includes: export_version
	ObservationTypeTransitKeyExport = "transit/key/export"

	// ObservationTypeTransitKeyExportBYOK is emitted when a key is exported using BYOK (Bring Your Own Key).
	// Metadata: key_name, type, derived, deletion_allowed, min_available_version,
	// min_decryption_version, min_encryption_version, latest_version, exportable,
	// allow_plaintext_backup, auto_rotate_period, imported_key, kdf (if derived),
	// kdf_mode (if derived), convergent_encryption (if derived), export_version (if specified),
	// destination_key
	ObservationTypeTransitKeyExportBYOK = "transit/key/export/byok"

	// ObservationTypeTransitKeyBackup is emitted when a key is backed up.
	// Metadata: key_name
	ObservationTypeTransitKeyBackup = "transit/key/backup"

	// ObservationTypeTransitKeyRestore is emitted when a key is restored from backup.
	// Metadata: key_name, force
	ObservationTypeTransitKeyRestore = "transit/key/restore"

	// ObservationTypeTransitKeyTrim is emitted when old key versions are trimmed.
	// Metadata: key_name, type, derived, deletion_allowed, min_available_version,
	// min_decryption_version, min_encryption_version, latest_version, exportable,
	// allow_plaintext_backup, auto_rotate_period, imported_key, kdf (if derived),
	// kdf_mode (if derived), convergent_encryption (if derived)
	ObservationTypeTransitKeyTrim = "transit/key/trim"
)
