/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

import type { MountsTuneConfigurationParametersRequest } from '@hashicorp/vault-client-typescript';

export default class GeneralSettingsForm extends Form<MountsTuneConfigurationParametersRequest> {
  formFieldGroups = [
    new FormFieldGroup('Version', [
      new FormField('engine_type', 'string'),
      new FormField('current_version', 'string'),
      new FormField('latest_version', 'string'),
      new FormField('plugin_version', 'string', { sensitive: true }),
    ]),
    new FormFieldGroup('Lease duration', [
      new FormField('default_lease_ttl', 'number', {
        label: 'Time-to-live (TTL) for secrets issued by this engine',
        editType: 'ttl',
      }),
      new FormField('max_lease_ttl', 'number', { label: 'Maximum time-to-live (TTL)', editType: 'ttl' }),
    ]),
    new FormFieldGroup('Metadata', [
      new FormField('path', 'string', { sensitive: true }),
      new FormField('accessor', 'string'),
      new FormField('description', 'string'),
    ]),
    new FormFieldGroup('Security', [
      new FormField('local', 'boolean'),
      new FormField('seal_wrap', 'boolean'),
    ]),
  ];
}
