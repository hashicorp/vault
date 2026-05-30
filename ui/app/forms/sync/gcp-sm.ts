/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { CredentialType, DestinationType } from 'sync/utils/constants';

import type { SystemWriteSyncDestinationsGcpSmNameRequest } from '@hashicorp/vault-client-typescript';
import CreateDestinationForm from './create-destination';

type GcpSmFormData = SystemWriteSyncDestinationsGcpSmNameRequest & {
  name: string;
  credential_type: CredentialType;
};

export default class GcpSmForm extends CreateDestinationForm<GcpSmFormData> {
  // the "credentials" param is not checked for "isAccountPluginConfigured" because it's never return by the API
  // additionally credentials can be set via GOOGLE_APPLICATION_CREDENTIALS env var so we cannot call it a required field in the ui.
  // thus we can never say for sure if the account accessType has been configured so we always return false
  isAccountPluginConfigured = false;

  get isWifPluginConfigured() {
    const { identity_token_audience, identity_token_ttl, service_account_email } = this.data;
    return !!identity_token_audience || !!identity_token_ttl || !!service_account_email;
  }

  accountCredentialGroup = new FormFieldGroup('JSON credentials', [
    new FormField('credentials', 'string', {
      label: 'JSON credentials',
      subText:
        'If empty, Vault will use the GOOGLE_APPLICATION_CREDENTIALS environment variable if configured.',
      editType: 'file',
      sensitive: true,
      docLink: '/vault/docs/secrets/gcp#authentication',
    }),
  ]);

  get wifCredentialGroup() {
    const serviceAccountField = new FormField('service_account_email', 'string', {
      label: 'Service account email',
    });
    return this.createWifCredentialGroup([serviceAccountField]);
  }

  get formFieldGroups() {
    const credentialGroup =
      this.credentialType === CredentialType.ACCOUNT ? this.accountCredentialGroup : this.wifCredentialGroup;
    return [
      new FormFieldGroup('Destination details', [
        this.commonFields.name,
        new FormField('project_id', 'string', {
          label: 'Project ID',
          subText:
            'The target project to manage secrets in. If set, overrides the project derived from the service account JSON credentials or application default credentials.',
        }),
      ]),
      credentialGroup,
      new FormFieldGroup('Advanced configuration', [
        this.commonFields.granularity,
        this.commonFields.secretNameTemplate,
        this.commonFields.customTags,
      ]),
    ];
  }

  toJSON() {
    const formState = super.toJSON();
    const data = this.getPayload<GcpSmFormData>(DestinationType.GcpSm, this.data, this.isNew);
    return { ...formState, data };
  }
}
