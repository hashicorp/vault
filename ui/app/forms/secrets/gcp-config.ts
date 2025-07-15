/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import WifConfigForm from './wif-config';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { GcpConfigFormData } from 'vault/secrets/engine';

export default class AzureConfigForm extends WifConfigForm<GcpConfigFormData> {
  // the "credentials" param is not checked for "isAccountPluginConfigured" because it's never return by the API
  // additionally credentials can be set via GOOGLE_APPLICATION_CREDENTIALS env var so we cannot call it a required field in the ui.
  // thus we can never say for sure if the account accessType has been configured so we always return false
  isAccountPluginConfigured = false;

  get isWifPluginConfigured() {
    const { identityTokenAudience, identityTokenTtl, serviceAccountEmail } = this.data;
    return !!identityTokenAudience || !!identityTokenTtl || !!serviceAccountEmail;
  }

  accountFields = [
    new FormField('credentials', 'string', {
      label: 'JSON credentials',
      subText:
        'If empty, Vault will use the GOOGLE_APPLICATION_CREDENTIALS environment variable if configured.',
      editType: 'file',
      docLink: '/vault/docs/secrets/gcp#authentication',
    }),
  ];

  optionFields = [
    new FormField('ttl', 'string', {
      label: 'Config TTL',
      editType: 'ttl',
      helperTextDisabled: 'Vault will use the default config TTL (time-to-live) for long-lived credentials.',
      helperTextEnabled:
        'The default config TTL (time-to-live) for long-lived credentials (i.e. service account keys).',
    }),
    new FormField('maxTtl', 'string', {
      label: 'Max TTL',
      editType: 'ttl',
      helperTextDisabled:
        'Vault will use the default maximum config TTL (time-to-live) for long-lived credentials.',
      helperTextEnabled:
        'The maximum config TTL (time-to-live) for long-lived credentials (i.e. service account keys).',
    }),
  ];

  wifFields = [
    this.commonWifFields.issuer,
    this.commonWifFields.identityTokenAudience,
    new FormField('serviceAccountEmail', 'string', {
      subText: 'Email ID for the Service Account to impersonate for Workload Identity Federation.',
    }),
    this.commonWifFields.identityTokenTtl,
  ];

  get formFieldGroups() {
    const defaultFields = this.accessType === 'account' ? this.accountFields : this.wifFields;
    return [
      new FormFieldGroup('default', defaultFields),
      new FormFieldGroup('More options', this.optionFields),
    ];
  }
}
