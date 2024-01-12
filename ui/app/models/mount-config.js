/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';

export default class MountConfigModel extends Model {
  @attr({
    label: 'Default Lease TTL',
    editType: 'ttl',
  })
  defaultLeaseTtl;

  @attr({
    label: 'Max Lease TTL',
    editType: 'ttl',
  })
  maxLeaseTtl;

  @attr({
    label: 'Request keys excluded from HMACing in audit',
    editType: 'stringArray',
    helpText: "Keys that will not be HMAC'd by audit devices in the request data object.",
  })
  auditNonHmacRequestKeys;

  @attr({
    label: 'Response keys excluded from HMACing in audit',
    editType: 'stringArray',
    helpText: "Keys that will not be HMAC'd by audit devices in the response data object.",
  })
  auditNonHmacResponseKeys;

  @attr('mountVisibility', {
    editType: 'boolean',
    label: 'List method when unauthenticated',
    defaultValue: false,
  })
  listingVisibility;

  @attr({
    label: 'Allowed passthrough request headers',
    helpText: 'Headers to allow and pass from the request to the backend',
    editType: 'stringArray',
  })
  passthroughRequestHeaders;

  @attr({
    label: 'Allowed response headers',
    helpText: 'Headers to allow, allowing a plugin to include them in the response.',
    editType: 'stringArray',
  })
  allowedResponseHeaders;

  @attr('string', {
    label: 'Token Type',
    helpText:
      'The type of token that should be generated via this role. For `default-service` and `default-batch` service and batch tokens will be issued respectively, unless the auth method explicitly requests a different type.',
    possibleValues: ['default-service', 'default-batch', 'batch', 'service'],
    noDefault: true,
  })
  tokenType;

  @attr({
    editType: 'stringArray',
  })
  allowedManagedKeys;
}
