import Model, { attr } from '@ember-data/model';
import { tracked } from '@glimmer/tracking';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const CRED_PROPS = {
  azurekeyvault: ['client_id', 'client_secret', 'tenant_id'],
  awskms: ['access_key', 'secret_key', 'session_token', 'endpoint'],
  gcpckms: ['service_account_file'],
};
const OPTIONAL_CRED_PROPS = ['session_token', 'endpoint'];
// since we have dynamic credential attributes based on provider we need a dynamic presence validator
// add validators for all cred props and return true for value if not associated with selected provider
const credValidators = Object.keys(CRED_PROPS).reduce((obj, providerKey) => {
  CRED_PROPS[providerKey].forEach((prop) => {
    if (!OPTIONAL_CRED_PROPS.includes(prop)) {
      obj[`credentials.${prop}`] = [
        {
          message: `${prop} is required`,
          validator(model) {
            return model.credentialProps.includes(prop) ? model.credentials[prop] : true;
          },
        },
      ];
    }
  });
  return obj;
}, {});
const validations = {
  name: [{ type: 'presence', message: 'Provider name is required' }],
  keyCollection: [{ type: 'presence', message: 'Key Vault instance name' }],
  ...credValidators,
};
@withModelValidations(validations)
export default class KeymgmtProviderModel extends Model {
  @attr('string') backend;
  @attr('string', {
    label: 'Provider name',
    subText: 'This is the name of the provider that will be displayed in Vault. This cannot be edited later.',
  })
  name;

  @attr('string', {
    label: 'Type',
    subText: 'Choose the provider type.',
    possibleValues: ['azurekeyvault', 'awskms', 'gcpckms'],
    noDefault: true,
  })
  provider;

  @attr('string', {
    label: 'Key Vault instance name',
    subText: 'The name of a Key Vault instance must be supplied. This cannot be edited later.',
  })
  keyCollection;

  idPrefix = 'provider/';
  type = 'provider';

  @tracked keys = [];
  @tracked credentials = null; // never returned from API -- set only during create/edit

  get icon() {
    return {
      azurekeyvault: 'azure-color',
      awskms: 'aws-color',
      gcpckms: 'gcp-color',
    }[this.provider];
  }
  get typeName() {
    return {
      azurekeyvault: 'Azure Key Vault',
      awskms: 'AWS Key Management Service',
      gcpckms: 'Google Cloud Key Management Service',
    }[this.provider];
  }
  get showFields() {
    const attrs = expandAttributeMeta(this, ['name', 'keyCollection']);
    attrs.splice(1, 0, { hasBlock: true, label: 'Type', value: this.typeName, icon: this.icon });
    const l = this.keys.length;
    const value = l
      ? `${l} ${l > 1 ? 'keys' : 'key'}`
      : this.canListKeys
      ? 'None'
      : 'You do not have permission to list keys';
    attrs.push({ hasBlock: true, isLink: l, label: 'Keys', value });
    return attrs;
  }
  get credentialProps() {
    if (!this.provider) return [];
    return CRED_PROPS[this.provider];
  }
  get credentialFields() {
    const [creds, fields] = this.credentialProps.reduce(
      ([creds, fields], prop) => {
        creds[prop] = null;
        let field = { name: `credentials.${prop}`, type: 'string', options: { label: prop } };
        if (prop === 'service_account_file') {
          field.options.subText = 'The path to a Google service account key file, not the file itself.';
        }
        fields.push(field);
        return [creds, fields];
      },
      [{}, []]
    );
    this.credentials = creds;
    return fields;
  }
  get createFields() {
    return expandAttributeMeta(this, ['provider', 'name', 'keyCollection']);
  }

  async fetchKeys(page) {
    if (this.canListKeys === false) {
      this.keys = [];
    } else {
      // try unless capabilities returns false
      try {
        this.keys = await this.store.lazyPaginatedQuery('keymgmt/key', {
          backend: 'keymgmt',
          provider: this.name,
          responsePath: 'data.keys',
          page,
        });
      } catch (error) {
        this.keys = [];
        if (error.httpStatus !== 404) {
          throw error;
        }
      }
    }
  }

  @lazyCapabilities(apiPath`${'backend'}/kms/${'id'}`, 'backend', 'id') providerPath;
  @lazyCapabilities(apiPath`${'backend'}/kms`, 'backend') providersPath;
  @lazyCapabilities(apiPath`${'backend'}/kms/${'id'}/key`, 'backend', 'id') providerKeysPath;

  get canCreate() {
    return this.providerPath.get('canCreate');
  }
  get canDelete() {
    return this.providerPath.get('canDelete');
  }
  get canEdit() {
    return this.providerPath.get('canUpdate');
  }
  get canRead() {
    return this.providerPath.get('canRead');
  }
  get canList() {
    return this.providersPath.get('canList');
  }
  get canListKeys() {
    return this.providerKeysPath.get('canList');
  }
  get canCreateKeys() {
    return this.providerKeysPath.get('canCreate');
  }
}
