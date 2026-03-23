/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import { toLabel } from 'core/helpers/to-label';

export default function tidyFieldLabel(field: string) {
  const label = toLabel([field]);
  return (
    {
      acme_account_safety_buffer: 'ACME account safety buffer',
      tidy_acme: 'Tidy ACME',
      enabled: 'Automatic tidy enabled',
      tidy_cert_store: 'Tidy the certificate store',
      tidy_cross_cluster_revoked_certs: 'Tidy cross-cluster revoked certificates',
      tidy_cmpv2_nonce_store: 'Tidy CMPv2 nonce store',
      tidy_move_legacy_ca_bundle: 'Tidy legacy CA bundle',
      tidy_revocation_queue: 'Tidy cross-cluster revocation requests',
      tidy_revoked_cert_issuer_associations: 'Tidy revoked certificate issuer associations',
      tidy_revoked_certs: 'Tidy revoked certificates',
    }[field] || label
  );
}
