/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class GcpConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord

  /* GCP config fields */
  @attr({
    label: 'Config TTL',
    editType: 'ttl',
    helperTextDisabled: 'The TTL (time-to-live) of generated tokens.',
  })
  ttl;

  @attr({
    label: 'Max TTL',
    editType: 'ttl',
    helperTextDisabled:
      'Specifies the maximum config TTL (time-to-live) for long-lived credentials (i.e. service account keys).',
  })
  maxTtl;

  /* GCP credential config field */
  @attr('string', {
    label: 'JSON credentials',
    subText:
      'If empty, Vault will use the GOOGLE_APPLICATION_CREDENTIALS environment variable if configured.',
    editType: 'file',
    docLink: '/vault/docs/secrets/gcp#authentication',
  })
  credentials; // obfuscated, never returned by API.

  /* WIF config fields */
  @attr('string', {
    subText:
      'The audience claim value for plugin identity tokens. Must match an allowed audience configured for the target IAM OIDC identity provider.',
  })
  identityTokenAudience;

  @attr({
    label: 'Identity token TTL',
    helperTextDisabled:
      'The TTL of generated tokens. Defaults to 1 hour, toggle on to specify a different value.',
    helperTextEnabled: 'The TTL of generated tokens.',
    editType: 'ttl',
  })
  identityTokenTtl;

  @attr('string', {
    subText: 'Email ID for the Service Account to impersonate for Workload Identity Federation.',
  })
  serviceAccountEmail;

  configurableParams = [
    'credentials',
    'serviceAccountEmail',
    'ttl',
    'maxTtl',
    'identityTokenAudience',
    'identityTokenTtl',
  ];

  get displayAttrs() {
    const formFields = expandAttributeMeta(this, this.configurableParams);
    return formFields.filter((attr) => attr.name !== 'credentials');
  }
}
