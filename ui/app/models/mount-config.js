/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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

  @attr('string', {
    editType: 'boolean',
    label: 'List method when unauthenticated',
    trueValue: 'unauth',
    falseValue: 'hidden',
  })
  listingVisibility;

  @attr({
    label: 'Allowed passthrough request headers',
    helpText: 'Headers to allow and pass from the request to the backend',
    editType: 'stringArray',
  })
  passthroughRequestHeaders;

  @attr('string', {
    label: 'Token Type',
    helpText:
      "The type of token that should be generated via this role. Can be `service`, `batch`, or `default` to use the mount's default (which unless changed will be `service` tokens).",
    possibleValues: ['default', 'batch', 'service'],
    defaultFormValue: 'default',
  })
  tokenType;
}
