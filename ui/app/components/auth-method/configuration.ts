/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { toLabel } from 'core/helpers/to-label';
import { get } from '@ember/object';

import type AuthMethodResource from 'vault/resources/auth/method';

interface Args {
  method: AuthMethodResource;
}
export default class AuthMethodConfigurationComponent extends Component<Args> {
  displayFields = [
    'type',
    'path',
    'description',
    'accessor',
    'local',
    'seal_wrap',
    'config.listing_visibility',
    'config.default_lease_ttl',
    'config.max_lease_ttl',
    'config.token_type',
    'config.audit_non_hmac_request_keys',
    'config.audit_non_hmac_response_keys',
    'config.passthrough_request_headers',
    'config.allowed_response_headers',
    'config.plugin_version',
  ];

  label = (field: string) => {
    const key = field.replace('config.', '');
    const label = toLabel([key]);
    // map specific fields to custom labels
    return (
      {
        listing_visibility: 'Use as preferred UI login method',
        default_lease_ttl: 'Default Lease TTL',
        max_lease_ttl: 'Max Lease TTL',
        audit_non_hmac_request_keys: 'Request keys excluded from HMACing in audit',
        audit_non_hmac_response_keys: 'Response keys excluded from HMACing in audit',
        passthrough_request_headers: 'Allowed passthrough request headers',
      }[key] || label
    );
  };
  value = (field: string) => {
    const { method } = this.args;
    if (field === 'config.listing_visibility') {
      return method.config.listing_visibility === 'unauth';
    }
    return get(method, field);
  };

  isTtl = (field: string) => {
    return ['config.default_lease_ttl', 'config.max_lease_ttl'].includes(field);
  };
}
