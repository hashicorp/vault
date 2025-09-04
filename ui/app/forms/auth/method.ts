/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import MountForm from 'vault/forms/mount';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { AuthMethodFormData } from 'vault/auth/methods';

export default class AuthMethodForm extends MountForm<AuthMethodFormData> {
  fieldProps = ['tuneFields', 'userLockoutConfigFields'];

  userLockoutConfigFields = [
    new FormField('user_lockout_config.lockout_threshold', 'string', {
      label: 'Lockout threshold',
      subText: 'Specifies the number of failed login attempts after which the user is locked out, e.g. 15.',
    }),
    new FormField('user_lockout_config.lockout_duration', undefined, {
      label: 'Lockout duration',
      helperTextEnabled: 'The duration for which a user will be locked out, e.g. "5s" or "30m".',
      editType: 'ttl',
      helperTextDisabled: 'No lockout duration configured.',
    }),

    new FormField('user_lockout_config.lockout_counter_reset', undefined, {
      label: 'Lockout counter reset',
      helperTextEnabled:
        'The duration after which the lockout counter is reset with no failed login attempts, e.g. "5s" or "30m".',
      editType: 'ttl',
      helperTextDisabled: 'No reset duration configured.',
    }),
    new FormField('user_lockout_config.lockout_disable', 'boolean', {
      label: 'Disable lockout for this mount',
      subText: 'If checked, disables the user lockout feature for this mount.',
    }),
  ];

  get tuneFields() {
    const readOnly = ['local', 'seal_wrap'];
    return this.formFieldGroups[1]?.['Method Options']?.filter((field) => {
      const isTuneable = !readOnly.includes(field.name);
      return isTuneable || (field.name === 'token_type' && this.normalizedType === 'token');
    });
  }

  formFieldGroups = [
    new FormFieldGroup('default', [this.fields.path]),
    new FormFieldGroup('Method Options', [
      this.fields.description,
      this.fields.listingVisibility,
      this.fields.local,
      this.fields.sealWrap,
      this.fields.defaultLeaseTtl,
      this.fields.maxLeaseTtl,
      new FormField('config.token_type', 'string', {
        label: 'Token type',
        helpText:
          'The type of token that should be generated via this role. For `default-service` and `default-batch` service and batch tokens will be issued respectively, unless the auth method explicitly requests a different type.',
        possibleValues: ['default-service', 'default-batch', 'batch', 'service'],
        noDefault: true,
      }),
      this.fields.auditNonHmacRequestKeys,
      this.fields.auditNonHmacResponseKeys,
      this.fields.passthroughRequestHeaders,
      this.fields.allowedResponseHeaders,
      this.fields.pluginVersion,
    ]),
  ];
}
