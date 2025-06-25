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
    const fields = ['type', 'path', 'description', 'accessor', 'local', 'sealWrap'];
    // no ttl options for keymgmt
    if (engineType !== 'keymgmt') {
      fields.push('config.defaultLeaseTtl', 'config.maxLeaseTtl');
    }
    fields.push(
      'config.allowedManagedKeys',
      'config.auditNonHmacRequestKeys',
      'config.auditNonHmacResponseKeys',
      'config.passthroughRequestHeaders',
      'config.allowedResponseHeaders'
    );
    if (engineType === 'kv' || engineType === 'generic') {
      fields.push('version');
    }
    // For WIF Secret engines, allow users to set the identity token key when mounting the engine.
    if (engineDisplayData(engineType)?.isWIF) {
      fields.push('config.identityTokenKey');
    }
    return fields;
  }

  label = (field) => {
    const key = field.replace('config.', '');
    const label = toLabel([key]);
    // map specific fields to custom labels
    return (
      {
        defaultLeaseTtl: 'Default Lease TTL',
        maxLeaseTtl: 'Max Lease TTL',
        auditNonHmacRequestKeys: 'Request keys excluded from HMACing in audit',
        auditNonHmacResponseKeys: 'Response keys excluded from HMACing in audit',
        passthroughRequestHeaders: 'Allowed passthrough request headers',
      }[key] || label
    );
  };
}
