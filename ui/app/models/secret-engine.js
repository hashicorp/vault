import Model, { attr, belongsTo } from '@ember-data/model';
import { computed } from '@ember/object'; // eslint-disable-line
import { equal } from '@ember/object/computed'; // eslint-disable-line
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

// identity will be managed separately and the inclusion
// of the system backend is an implementation detail
const LIST_EXCLUDED_BACKENDS = ['system', 'identity'];

const validations = {
  path: [{ type: 'presence', message: "Path can't be blank." }],
  maxVersions: [
    { type: 'number', message: 'Maximum versions must be a number.' },
    { type: 'length', options: { min: 1, max: 16 }, message: 'You cannot go over 16 characters.' },
  ],
};

@withModelValidations(validations)
export default class SecretEngineModel extends Model {
  @belongsTo('mount-config', { async: false, inverse: null }) config;
  @attr('string') path;
  @attr('string', { label: 'Secret engine type' }) type;
  @attr('string', { editType: 'textarea' }) description;
  @attr('string') accessor;
  @attr('string') name;

  @attr('number', {
    label: 'Version',
    helpText:
      'The KV Secrets Engine can operate in different modes. Version 1 is the original generic Secrets Engine the allows for storing of static key/value pairs. Version 2 added more features including data versioning, TTLs, and check and set.',
    possibleValues: [2, 1],
    // This shouldn't be defaultValue because if no version comes back from API we should assume it's v1
    defaultFormValue: 2, // Set the form to 2 by default
  })
  version; // will only have value for type=kv

  @attr('boolean', {
    helpText:
      'When Replication is enabled, a local mount will not be replicated across clusters. This can only be specified at mount time.',
  })
  local;

  @attr('boolean', {
    helpText:
      'When enabled - if a seal supporting seal wrapping is specified in the configuration, all critical security parameters (CSPs) in this backend will be seal wrapped. (For K/V mounts, all values will be seal wrapped.) This can only be specified at mount time.',
  })
  sealWrap;

  // KV 2 additional config default options
  @attr('number', {
    defaultValue: 0,
    label: 'Maximum number of versions',
    subText:
      'The number of versions to keep per key. Once the number of keys exceeds the maximum number set here, the oldest version will be permanently deleted. This value applies to all keys, but a key’s metadata settings can overwrite this value. When 0 is used or the value is unset, Vault will keep 10 versions.',
  })
  maxVersions;

  @attr('boolean', {
    defaultValue: false,
    label: 'Require Check and Set',
    subText:
      'If checked, all keys will require the cas parameter to be set on all write requests. A key’s metadata settings can overwrite this value.',
  })
  casRequired;

  @attr({
    defaultValue: 0,
    editType: 'ttl',
    label: 'Automate secret deletion',
    helperTextDisabled: 'A secret’s version must be manually deleted.',
    helperTextEnabled: 'Delete all new versions of this secret after',
  })
  deleteVersionAfter;

  // end KV 2 additional config options

  get modelTypeForKV() {
    const type = this.engineType;
    let modelType = 'secret';
    if ((type === 'kv' || type === 'generic') && this.version === 2) {
      modelType = 'secret-v2';
    }
    return modelType;
  }

  get isV2KV() {
    return this.modelTypeForKV === 'secret-v2';
  }

  get engineFormFields() {
    if (!this.engineType) return null;
    const type = this.engineType;
    const fields = ['type', 'path', 'description', 'accessor', 'local', 'sealWrap'];
    // no ttl options for keymgmt
    const ttl = type !== 'keymgmt' ? 'defaultLeaseTtl,maxLeaseTtl,' : '';
    fields.push(
      `config.{${ttl}auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders}`
    );
    if (type === 'kv' || type === 'generic') {
      fields.push('version');
    }
    // version comes in as number not string
    if (type === 'kv' && this.version === 2) {
      fields.push('casRequired', 'deleteVersionAfter', 'maxVersions');
    }
    return fields;
  }

  get engineFormFieldGroups() {
    if (!this.engineType) return null;
    let defaultFields = ['path'];
    let optionFields;
    const CORE_OPTIONS = ['description', 'config.listingVisibility', 'local', 'sealWrap'];

    switch (this.engineType) {
      case 'kv':
        defaultFields = ['path', 'maxVersions', 'casRequired', 'deleteVersionAfter'];
        optionFields = [
          'version',
          ...CORE_OPTIONS,
          `config.{defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders}`,
        ];
        break;
      case 'generic':
        optionFields = [
          'version',
          ...CORE_OPTIONS,
          `config.{defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders}`,
        ];
        break;
      case 'database':
        // Highlight TTLs in default
        defaultFields = ['path', 'config.{defaultLeaseTtl}', 'config.{maxLeaseTtl}'];
        optionFields = [
          ...CORE_OPTIONS,
          'config.{auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders}',
        ];
        break;
      case 'keymgmt':
        // no ttl options for keymgmt
        optionFields = [
          ...CORE_OPTIONS,
          'config.{auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders}',
        ];
        break;
      default:
        defaultFields = ['path'];
        optionFields = [
          ...CORE_OPTIONS,
          `config.{defaultLeaseTtl,maxLeaseTtl,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders}`,
        ];
        break;
    }

    return [
      { default: defaultFields },
      {
        'Method Options': optionFields,
      },
    ];
  }

  get attrs() {
    if (!this.engineType) return null;
    return expandAttributeMeta(this, this.engineFormFields);
  }

  get fieldGroups() {
    if (!this.engineType) return null;
    return fieldToAttrs(this, this.engineFormFieldGroups);
  }

  get icon() {
    if (!this.engineType) return 'secrets';
    switch (this.engineType) {
      case 'keymgmt':
        return 'key';
      default:
        return this.engineType;
    }
  }

  // namespaces introduced types with a `ns_` prefix for built-in engines
  // so we need to strip that to normalize the type
  get engineType() {
    return (this.type || '').replace(/^ns_/, '');
  }

  get shouldIncludeInList() {
    return !LIST_EXCLUDED_BACKENDS.includes(this.engineType);
  }

  get localDisplay() {
    return this.local ? 'local' : 'replicated';
  }

  // EDITING SSH CA config
  @attr('string') privateKey;
  @attr('string') publicKey;
  @attr('boolean', { defaultValue: true }) generateSigningKey;

  saveCA(options) {
    if (this.type !== 'ssh') {
      return;
    }
    if (options.isDelete) {
      this.privateKey = null;
      this.publicKey = null;
      this.generateSigningKey = false;
    }
    return this.save({
      adapterOptions: {
        options: options,
        apiPath: 'config/ca',
        attrsToSend: ['privateKey', 'publicKey', 'generateSigningKey'],
      },
    });
  }

  saveZeroAddressConfig() {
    return this.save({
      adapterOptions: {
        adapterMethod: 'saveZeroAddressConfig',
      },
    });
  }

  // EDITING AWS config
  @attr('string') lease;
  @attr('string') leaseMax;
}
