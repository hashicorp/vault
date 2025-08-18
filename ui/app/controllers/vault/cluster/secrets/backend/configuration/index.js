/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { toLabel } from 'core/helpers/to-label';
import engineDisplayData from 'vault/helpers/engines-display-data';

export default class SecretsBackendConfigurationController extends Controller {
  get displayFields() {
    const { engineType } = this.model.secretsEngine;
    const fields = ['type', 'path', 'description', 'accessor', 'local', 'seal_wrap'];
    // no ttl options for keymgmt
    if (engineType !== 'keymgmt') {
      fields.push('config.default_lease_ttl', 'config.max_lease_ttl');
    }
    fields.push(
      'config.allowed_managed_keys',
      'config.audit_non_hmac_request_keys',
      'config.audit_non_hmac_response_keys',
      'config.passthrough_request_headers',
      'config.allowed_response_headers'
    );
    if (engineType === 'kv' || engineType === 'generic') {
      fields.push('version');
    }
    // For WIF Secret engines, allow users to set the identity token key when mounting the engine.
    if (engineDisplayData(engineType)?.isWIF) {
      fields.push('config.identity_token_key');
    }
    return fields;
  }

  label = (field) => {
    const key = field.replace('config.', '');
    const label = toLabel([key]);
    // map specific fields to custom labels
    return (
      {
        default_lease_ttl: 'Default Lease TTL',
        max_lease_ttl: 'Max Lease TTL',
        audit_non_hmac_request_keys: 'Request keys excluded from HMACing in audit',
        audit_non_hmac_response_keys: 'Response keys excluded from HMACing in audit',
        passthrough_request_headers: 'Allowed passthrough request headers',
      }[key] || label
    );
  };
}
