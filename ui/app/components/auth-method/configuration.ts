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
    'sealWrap',
    'config.listingVisibility',
    'config.defaultLeaseTtl',
    'config.maxLeaseTtl',
    'config.tokenType',
    'config.auditNonHmacRequestKeys',
    'config.auditNonHmacResponseKeys',
    'config.passthroughRequestHeaders',
    'config.allowedResponseHeaders',
    'config.pluginVersion',
  ];

  label = (field: string) => {
    const key = field.replace('config.', '');
    const label = toLabel([key]);
    // map specific fields to custom labels
    return (
      {
        listingVisibility: 'Use as preferred UI login method',
        defaultLeaseTtl: 'Default Lease TTL',
        maxLeaseTtl: 'Max Lease TTL',
        auditNonHmacRequestKeys: 'Request keys excluded from HMACing in audit',
        auditNonHmacResponseKeys: 'Response keys excluded from HMACing in audit',
        passthroughRequestHeaders: 'Allowed passthrough request headers',
      }[key] || label
    );
  };
  value = (field: string) => {
    const { method } = this.args;
    if (field === 'config.listingVisibility') {
      return method.config.listingVisibility === 'unauth';
    }
    return get(method, field);
  };

  isTtl = (field: string) => {
    return ['config.defaultLeaseTtl', 'config.maxLeaseTtl'].includes(field);
  };
}
