/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { CredentialType, DestinationType, GcpEncryptionType } from 'sync/utils/constants';
import { gcpRegions } from 'vault/helpers/gcp-regions';

import type { SystemWriteSyncDestinationsGcpSmNameRequest } from '@hashicorp/vault-client-typescript';
import type { FormOptions } from '../form';
import type { Validations } from 'vault/app-types';
import CreateDestinationForm from './create-destination';

type GcpSmFormData = SystemWriteSyncDestinationsGcpSmNameRequest & {
  name: string;
  credential_type: CredentialType;
  encryption_type: GcpEncryptionType;
  kms_key_id?: string;
  replica_regions?: Record<string, string>;
};

export default class GcpSmForm extends CreateDestinationForm<GcpSmFormData> {
  // the "credentials" param is not checked for "isAccountPluginConfigured" because it's never return by the API
  // additionally credentials can be set via GOOGLE_APPLICATION_CREDENTIALS env var so we cannot call it a required field in the ui.
  // thus we can never say for sure if the account accessType has been configured so we always return false
  isAccountPluginConfigured = false;

  constructor(data: Partial<GcpSmFormData> = {}, options: FormOptions = {}, validations?: Validations) {
    super(data, options, validations);
    // the API doesn't return the encryption method, so we derive it from which field is populated.
    // replica_regions is used by both google-managed-encryption (regions only) and regional-kms (regions + keys),
    // so a populated KMS key value is what distinguishes the selected radio option on Edit page.
    if (!this.data.encryption_type) {
      if (this.data.kms_key_id) {
        this.data.encryption_type = GcpEncryptionType.GLOBAL_KMS;
      } else if (Object.values(this.data.replica_regions ?? {}).some((value) => !!value)) {
        this.data.encryption_type = GcpEncryptionType.REGIONAL_KMS;
      } else {
        this.data.encryption_type = GcpEncryptionType.GOOGLE_MANAGED;
      }
    }
  }

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

  get replicaRegionsAndEncryptionGroup() {
    const isRegionalKms = this.data.encryption_type === GcpEncryptionType.REGIONAL_KMS;
    return new FormFieldGroup('Replica regions and encryption', [
      new FormField('encryption_type', 'string', {
        label: 'Encryption method',
        editType: 'radio',
        editDisabled: true,
        possibleValues: [
          {
            label: 'Google-managed encryption',
            value: GcpEncryptionType.GOOGLE_MANAGED,
            subText:
              'Use Google-managed encryption. No encryption key is required. You can still add replica regions if needed.',
          },
          {
            label: 'Global KMS key',
            value: GcpEncryptionType.GLOBAL_KMS,
            subText: 'Use one customer-managed Cloud KMS key for automatically replicated secrets.',
          },
          {
            label: 'Regional KMS keys',
            value: GcpEncryptionType.REGIONAL_KMS,
            subText: 'Add region-specific Cloud KMS keys for user-managed replication.',
          },
        ],
      }),
      new FormField('kms_key_id', 'string', {
        label: 'KMS key ID',
        editDisabled: true,
        isRequired: true,
        subText:
          'Enter the full Cloud KMS key resource name to encrypt secrets with automatic replication. The key must be in the global location.',
      }),
      new FormField('replica_regions', 'object', {
        label: isRegionalKms ? 'Replica regions and KMS keys' : 'Replica regions',
        editType: 'keyValueInputs',
        editDisabled: true,
        addRowButtonText: 'Add region',
        subText: isRegionalKms
          ? 'Add each replica region and the Cloud KMS key used to encrypt secrets in that region. Each KMS key must be in the same location as its replica region.'
          : 'Add replica regions for user-managed replication. If no regions are added, Google Secret Manager uses automatic replication. No KMS key is required when using Google-managed encryption.',
        // google-managed-encryption only stores the selected regions (empty KMS key values), so the KMS key ID input isn't rendered at all
        keyValueFields: isRegionalKms
          ? [
              {
                name: 'key',
                label: 'Replica region',
                type: 'select',
                possibleValues: gcpRegions(),
                noDefault: true,
                isRequired: true,
              },
              {
                name: 'value',
                label: 'KMS key ID',
                type: 'text',
                placeholder: 'KMS key ID',
                isRequired: true,
              },
            ]
          : [
              {
                name: 'key',
                label: 'Region',
                type: 'select',
                possibleValues: gcpRegions(),
                noDefault: true,
                isRequired: false,
              },
            ],
      }),
    ]);
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
      this.replicaRegionsAndEncryptionGroup,
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
    // encryption_type is only used to determine which encryption field to render, it isn't part of the API payload
    delete (data as Partial<GcpSmFormData>).encryption_type;
    return { ...formState, data };
  }
}
