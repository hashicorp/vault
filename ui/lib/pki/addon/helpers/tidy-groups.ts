/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Helper from '@ember/component/helper';
import { service } from '@ember/service';

import type VersionService from 'vault/services/version';

export default class TidyGroups extends Helper {
  @service declare readonly version: VersionService;

  compute([group]: [group?: 'auto' | 'manual']) {
    // groups that are shared between auto and manual tidy types
    const shared: Array<Record<string, Array<string>>> = [
      {
        'Universal operations': [
          'tidy_cert_store',
          'tidy_cert_metadata',
          'tidy_revoked_certs',
          'tidy_revoked_cert_issuer_associations',
          'safety_buffer',
          'pause_duration',
        ],
      },
      {
        'ACME operations': ['tidy_acme', 'acme_account_safety_buffer'],
      },
      {
        'Issuer operations': ['tidy_expired_issuers', 'tidy_move_legacy_ca_bundle', 'issuer_safety_buffer'],
      },
    ];
    // cross cluster operations are only available in enterprise
    if (this.version.isEnterprise) {
      shared.push({
        'Cross-cluster operations': [
          'tidy_revocation_queue',
          'tidy_cross_cluster_revoked_certs',
          'tidy_cmpv2_nonce_store',
          'revocation_queue_safety_buffer',
        ],
      });
    }
    // auto tidy specific fields
    const auto = {
      default: ['interval_duration', 'min_startup_backoff_duration', 'max_startup_backoff_duration'],
    };

    if (group === 'auto') {
      return [auto];
    }
    if (group === 'manual') {
      return shared;
    }
    // if group is not specified, return combined groups
    // add enabled field to the top of auto tidy fields for details view
    auto.default.unshift('enabled');
    shared.unshift(auto);
    return shared;
  }
}
