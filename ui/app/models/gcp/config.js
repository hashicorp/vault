/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class GcpConfig extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord

  // GCP only field
  @attr('string', {
    label: 'JSON credentials',
    subText:
      'If empty, Vault will use the GOOGLE_APPLICATION_CREDENTIALS environment variable if configured.',
    editType: 'file',
    docLink: '/vault/docs/secrets/gcp#authentication',
  })
  credentials; // obfuscated, never returned by API.

  // WIF only fields
  @attr('string', {
    subText:
      'The audience claim value for plugin identity tokens. Must match an allowed audience configured for the target IAM OIDC identity provider.',
  })
  identityTokenAudience;

  @attr({
    label: 'Identity token TTL',
    helperTextDisabled:
      'The TTL of generated tokens. Defaults to 1 hour, turn on the toggle to specify a different value.',
    helperTextEnabled: 'The TTL of generated tokens.',
    editType: 'ttl',
  })
  identityTokenTtl;

  @attr('string', {
    subText: 'Email ID for the Service Account to impersonate for Workload Identity Federation.',
  })
  serviceAccountEmail;

  // Fields that show regardless of access type
  @attr({
    label: 'Config TTL',
    editType: 'ttl',
    helperTextDisabled: 'Vault will use the default config TTL (time-to-live) for long-lived credentials.',
    helperTextEnabled:
      'The default config TTL (time-to-live) for long-lived credentials (i.e. service account keys).',
  })
  ttl;

  @attr({
    label: 'Max TTL',
    editType: 'ttl',
    helperTextDisabled:
      'Vault will use the default maximum config TTL (time-to-live) for long-lived credentials.',
    helperTextEnabled:
      'The maximum config TTL (time-to-live) for long-lived credentials (i.e. service account keys).',
  })
  maxTtl;

  configurableParams = [
    'credentials',
    'serviceAccountEmail',
    'ttl',
    'maxTtl',
    'identityTokenAudience',
    'identityTokenTtl',
  ];

  get isWifPluginConfigured() {
    return !!this.identityTokenAudience || !!this.identityTokenTtl || !!this.serviceAccountEmail;
  }

  // the "credentials" param is not checked for "isAccountPluginConfigured" because it's never return by the API
  // additionally credentials can be set via GOOGLE_APPLICATION_CREDENTIALS env var so we cannot call it a required field in the ui.
  // thus we can never say for sure if the account accessType has been configured so we always return false
  isAccountPluginConfigured = false;

  get displayAttrs() {
    const formFields = expandAttributeMeta(this, this.configurableParams);
    return formFields.filter((attr) => attr.name !== 'credentials');
  }

  get fieldGroupsWif() {
    return fieldToAttrs(this, this.formFieldGroups('wif'));
  }

  get fieldGroupsAccount() {
    return fieldToAttrs(this, this.formFieldGroups('account'));
  }

  formFieldGroups(accessType = 'account') {
    const formFieldGroups = [];
    if (accessType === 'wif') {
      formFieldGroups.push({
        default: ['identityTokenAudience', 'serviceAccountEmail', 'identityTokenTtl'],
      });
    }
    if (accessType === 'account') {
      formFieldGroups.push({
        default: ['credentials'],
      });
    }
    formFieldGroups.push({
      'More options': ['ttl', 'maxTtl'],
    });
    return formFieldGroups;
  }
}
