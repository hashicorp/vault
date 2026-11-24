/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import type { Validations } from 'vault/app-types';
import { get } from '@ember/object';
import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import {
  KeyManagementWriteKmsProviderRequest,
  KeyManagementWriteKmsProviderRequestProviderEnum,
} from '@hashicorp/vault-client-typescript';

type ProviderFormData = KeyManagementWriteKmsProviderRequest & {
  name: string;
  keys?: Array<Record<string, unknown>>;
};

type ProviderCredentialValidationModel = {
  provider?: keyof typeof CRED_PROPS;
  credentialProps?: string[];
  credentials?: Record<string, string | undefined | null>;
};

interface Validator {
  message: string;
  validator(model: ProviderCredentialValidationModel): boolean | string;
}

export const PROVIDER_TYPES = Object.values(KeyManagementWriteKmsProviderRequestProviderEnum);

export const CRED_PROPS = {
  azurekeyvault: ['client_id', 'client_secret', 'tenant_id'],
  awskms: ['access_key', 'secret_key', 'session_token', 'endpoint'],
  gcpckms: ['service_account_file'],
};

const CRED_PROPS_LABEL = {
  client_id: 'Client ID',
  client_secret: 'Client secret',
  tenant_id: 'Tenant ID',
  access_key: 'Access key',
  secret_key: 'Secret key',
  session_token: 'Session token',
  endpoint: 'Endpoint',
  service_account_file: 'Service account file',
};

export const OPTIONAL_CRED_PROPS = ['session_token', 'endpoint'];

export default class KeymgmtProviderForm extends Form<ProviderFormData> {
  icon = 'key';
  idPrefix = 'provider/';
  type = 'provider';

  constructor(...args: ConstructorParameters<typeof Form<ProviderFormData>>) {
    super(...args);

    // Provider read responses do not include credentials, but the form binds nested
    // values like "credentials.client_id". Ensure the parent object always exists.
    if (!this.data.credentials) {
      this.data.credentials = {};
    }
  }

  get credentialProps() {
    const provider = this.data.provider;
    if (!provider) return [];
    return CRED_PROPS[provider as KeyManagementWriteKmsProviderRequestProviderEnum] || [];
  }

  get createFields() {
    return [
      new FormField('provider', 'string', {
        label: 'Type',
        subText: 'Choose the provider type.',
        possibleValues: PROVIDER_TYPES,
        noDefault: true,
      }),
      new FormField('name', 'string', {
        label: 'Provider name',
        subText:
          'This is the name of the provider that will be displayed in Vault. This cannot be edited later.',
      }),
      new FormField('key_collection', 'string', {
        label: 'Key Vault instance name',
        subText: 'The name of a Key Vault instance must be supplied. This cannot be edited later.',
      }),
    ];
  }

  get credentialFields() {
    return this.credentialProps.map((prop) => {
      const options = {
        label: CRED_PROPS_LABEL[prop as keyof typeof CRED_PROPS_LABEL],
        ...(prop === 'service_account_file'
          ? { subText: 'The path to a Google service account key file, not the file itself.' }
          : {}),
      };
      return new FormField(`credentials.${prop}`, 'string', options);
    });
  }

  get formFieldGroups() {
    const groups: FormFieldGroup[] = [];

    if (this.isNew) {
      // Create mode: provider type, name, key_collection, credentials
      groups.push(new FormFieldGroup('default', [...this.createFields]));

      // Add credential fields based on provider type
      if (this.credentialProps.length > 0) {
        groups.push(new FormFieldGroup('credentials', this.credentialFields));
      }
    } else {
      // Edit mode: only credentials if any
      if (this.credentialProps.length > 0) {
        groups.push(new FormFieldGroup('credentials', this.credentialFields));
      }
    }

    return groups;
  }

  // since we have dynamic credential attributes based on provider we need a dynamic presence validator
  // add validators for all cred props and return true for value if not associated with selected provider
  get credValidators() {
    return (Object.keys(CRED_PROPS) as Array<keyof typeof CRED_PROPS>).reduce(
      (obj, providerKey) => {
        CRED_PROPS[providerKey].forEach((prop: string) => {
          if (!OPTIONAL_CRED_PROPS.includes(prop)) {
            obj[`credentials.${prop}`] = [
              {
                message: `${CRED_PROPS_LABEL[prop as keyof typeof CRED_PROPS_LABEL]} is required`,
                validator(model: ProviderCredentialValidationModel) {
                  const selectedProvider = model.provider;
                  if (!selectedProvider) return true;

                  if (!CRED_PROPS[selectedProvider]?.includes(prop)) return true;

                  const value = get(model, `credentials.${prop}`);
                  return typeof value === 'string' ? value.trim().length > 0 : Boolean(value);
                },
              },
            ];
          }
        });
        return obj;
      },
      {} as Record<string, Validator[]>
    );
  }

  validations: Validations = {
    name: [{ type: 'presence', message: 'Provider name is required' }],
    key_collection: [{ type: 'presence', message: 'Key Vault instance name is required' }],
    ...this.credValidators,
  };
}
