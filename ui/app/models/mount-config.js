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
    label: 'Token type',
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

  @attr('string', {
    label: 'Plugin version',
    subText:
      'Specifies the semantic version of the plugin to use, e.g. "v1.0.0". If unspecified, the server will select any matching un-versioned plugin that may have been registered, the latest versioned plugin registered, or a built-in plugin in that order of precedence.',
  })
  pluginVersion;

  // Auth mount userLockoutConfig params, added to user_lockout_config object in saveModel method
  @attr('string', {
    label: 'Lockout threshold',
    subText: 'Specifies the number of failed login attempts after which the user is locked out, e.g. 15.',
  })
  lockoutThreshold;

  @attr({
    label: 'Lockout duration',
    helperTextEnabled: 'The duration for which a user will be locked out, e.g. "5s" or "30m".',
    editType: 'ttl',
    helperTextDisabled: 'No lockout duration configured.',
  })
  lockoutDuration;

  @attr({
    label: 'Lockout counter reset',
    helperTextEnabled:
      'The duration after which the lockout counter is reset with no failed login attempts, e.g. "5s" or "30m".',
    editType: 'ttl',
    helperTextDisabled: 'No reset duration configured.',
  })
  lockoutCounterReset;

  @attr('boolean', {
    label: 'Disable lockout for this mount',
    subText: 'If checked, disables the user lockout feature for this mount.',
  })
  lockoutDisable;
  // end of user_lockout_config params
}
