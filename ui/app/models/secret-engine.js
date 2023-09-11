/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr, belongsTo } from '@ember-data/model';
import { computed } from '@ember/object'; // eslint-disable-line
import { equal } from '@ember/object/computed'; // eslint-disable-line
import { withModelValidations } from 'vault/decorators/model-validations';
import { withExpandedAttributes } from 'vault/decorators/model-expanded-attributes';
import { supportedSecretBackends } from 'vault/helpers/supported-secret-backends';
import { isAddonEngine, allEngines } from 'vault/helpers/mountable-secret-engines';

const LINKED_BACKENDS = supportedSecretBackends();

// identity will be managed separately and the inclusion
// of the system backend is an implementation detail
const LIST_EXCLUDED_BACKENDS = ['system', 'identity'];

const validations = {
  path: [
    { type: 'presence', message: "Path can't be blank." },
    {
      type: 'containsWhiteSpace',
      message:
        "Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests.",
      level: 'warn',
    },
  ],
  maxVersions: [
    { type: 'number', message: 'Maximum versions must be a number.' },
    { type: 'length', options: { min: 1, max: 16 }, message: 'You cannot go over 16 characters.' },
  ],
};

@withModelValidations(validations)
@withExpandedAttributes()
export default class SecretEngineModel extends Model {
  @attr('string') path;
  @attr('string') type;
  @attr('string', {
    editType: 'textarea',
  })
  description;
  @belongsTo('mount-config', { async: false, inverse: null }) config;

  // Enterprise options (still available on OSS)
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
  @attr('boolean') externalEntropyAccess;

  // options.version
  @attr('number', {
    label: 'Version',
    helpText:
      'The KV Secrets Engine can operate in different modes. Version 1 is the original generic Secrets Engine the allows for storing of static key/value pairs. Version 2 added more features including data versioning, TTLs, and check and set.',
    possibleValues: [2, 1],
    // This shouldn't be defaultValue because if no version comes back from API we should assume it's v1
    defaultFormValue: 2, // Set the form to 2 by default
  })
  version;

  // SSH specific attributes
  @attr('string') privateKey;
  @attr('string') publicKey;
  @attr('boolean', {
    defaultValue: true,
  })
  generateSigningKey;

  // AWS specific attributes
  @attr('string') lease;
  @attr('string') leaseMax;

  // Returned from API response
  @attr('string') accessor;

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

  /* GETTERS */
  get modelTypeForKV() {
    const engineType = this.engineType;
    if ((engineType === 'kv' || engineType === 'generic') && this.version === 2) {
      return 'secret-v2';
    }
    return 'secret';
  }
  get isV2KV() {
    return this.modelTypeForKV === 'secret-v2';
  }

  get attrs() {
    return this.formFields.map((fieldName) => {
      return this.allByKey[fieldName];
    });
  }

  get fieldGroups() {
    return this._expandGroups(this.formFieldGroups);
  }

  get icon() {
    const defaultIcon = this.engineType || 'secrets';
    return (
      {
        keymgmt: 'key',
        kmip: 'secrets',
        ldap: 'folder-users',
      }[this.engineType] || defaultIcon
    );
  }

  get engineType() {
    return (this.type || '').replace(/^ns_/, '');
  }

  get shouldIncludeInList() {
    return !LIST_EXCLUDED_BACKENDS.includes(this.engineType);
  }

  get isSupportedBackend() {
    return LINKED_BACKENDS.includes(this.engineType);
  }

  get backendLink() {
    if (this.engineType === 'database') {
      return 'vault.cluster.secrets.backend.overview';
    }
    if (isAddonEngine(this.engineType, this.version)) {
      const { engineRoute } = allEngines().findBy('type', this.engineType);
      return `vault.cluster.secrets.backend.${engineRoute}`;
    }
    return `vault.cluster.secrets.backend.list-root`;
  }

  get localDisplay() {
    return this.local ? 'local' : 'replicated';
  }

  get formFields() {
    const type = this.engineType;
    const fields = ['type', 'path', 'description', 'accessor', 'local', 'sealWrap'];
    // no ttl options for keymgmt
    if (type !== 'keymgmt') {
      fields.push('config.defaultLeaseTtl', 'config.maxLeaseTtl');
    }
    fields.push(
      'config.allowedManagedKeys',
      'config.auditNonHmacRequestKeys',
      'config.auditNonHmacResponseKeys',
      'config.passthroughRequestHeaders',
      'config.allowedResponseHeaders'
    );
    if (type === 'kv' || type === 'generic') {
      fields.push('version');
    }
    // version comes in as number not string
    if (type === 'kv' && parseInt(this.version, 10) === 2) {
      fields.push('casRequired', 'deleteVersionAfter', 'maxVersions');
    }
    return fields;
  }

  get formFieldGroups() {
    let defaultFields = ['path'];
    let optionFields;
    const CORE_OPTIONS = ['description', 'config.listingVisibility', 'local', 'sealWrap'];
    const STANDARD_CONFIG = [
      'config.auditNonHmacRequestKeys',
      'config.auditNonHmacResponseKeys',
      'config.passthroughRequestHeaders',
      'config.allowedResponseHeaders',
    ];

    switch (this.engineType) {
      case 'kv':
        defaultFields = ['path', 'maxVersions', 'casRequired', 'deleteVersionAfter'];
        optionFields = [
          'version',
          ...CORE_OPTIONS,
          'config.defaultLeaseTtl',
          'config.maxLeaseTtl',
          'config.allowedManagedKeys',
          ...STANDARD_CONFIG,
        ];
        break;
      case 'generic':
        optionFields = [
          'version',
          ...CORE_OPTIONS,
          'config.defaultLeaseTtl',
          'config.maxLeaseTtl',
          'config.allowedManagedKeys',
          ...STANDARD_CONFIG,
        ];
        break;
      case 'database':
        // Highlight TTLs in default
        defaultFields = ['path', 'config.defaultLeaseTtl', 'config.maxLeaseTtl'];
        optionFields = [...CORE_OPTIONS, 'config.allowedManagedKeys', ...STANDARD_CONFIG];
        break;
      case 'pki':
        defaultFields = ['path', 'config.defaultLeaseTtl', 'config.maxLeaseTtl', 'config.allowedManagedKeys'];
        optionFields = [...CORE_OPTIONS, ...STANDARD_CONFIG];
        break;
      case 'keymgmt':
        // no ttl options for keymgmt
        optionFields = [...CORE_OPTIONS, 'config.allowedManagedKeys', ...STANDARD_CONFIG];
        break;
      default:
        defaultFields = ['path'];
        optionFields = [
          ...CORE_OPTIONS,
          'config.defaultLeaseTtl',
          'config.maxLeaseTtl',
          'config.allowedManagedKeys',
          ...STANDARD_CONFIG,
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

  /* ACTIONS */
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
}
