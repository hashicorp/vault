/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import MountForm from 'vault/forms/mount';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';

import type { SecretsEngineFormData } from 'vault/secrets/engine';
import type Form from 'vault/forms/form';

export default class SecretsEngineForm extends MountForm<SecretsEngineFormData> {
  constructor(...args: ConstructorParameters<typeof Form>) {
    super(...args);
    // path validation is already defined on the MountForm class
    // add validation for kv max versions
    this.validations['kv_config.max_versions'] = [
      { type: 'number', message: 'Maximum versions must be a number.' },
      { type: 'length', options: { min: 1, max: 16 }, message: 'You cannot go over 16 characters.' },
    ];
  }

  coreOptionFields = [this.fields.description, this.fields.local, this.fields.sealWrap];

  leaseConfigFields = [
    this.fields.defaultLeaseTtl,
    this.fields.maxLeaseTtl,
    new FormField('config.allowed_managed_keys', 'string', {
      label: 'Allowed managed keys',
      editType: 'stringArray',
    }),
  ];

  standardConfigFields = [
    this.fields.auditNonHmacRequestKeys,
    this.fields.auditNonHmacResponseKeys,
    this.fields.passthroughRequestHeaders,
    this.fields.allowedResponseHeaders,
  ];

  get defaultFields() {
    const fields = [new FormField('path', 'string')];
    if (this.normalizedType === 'kv') {
      fields.push(
        new FormField('kv_config.max_versions', 'number', {
          label: 'Maximum number of versions',
          subText:
            'The number of versions to keep per key. Once the number of keys exceeds the maximum number set here, the oldest version will be permanently deleted. This value applies to all keys, but a key’s metadata settings can overwrite this value. When 0 is used or the value is unset, Vault will keep 10 versions.',
        }),
        new FormField('kv_config.cas_required', 'boolean', {
          label: 'Require Check and Set',
          subText:
            'If checked, all keys will require the cas parameter to be set on all write requests. A key’s metadata settings can overwrite this value.',
        }),
        new FormField('kv_config.delete_version_after', 'string', {
          editType: 'ttl',
          label: 'Automate secret deletion',
          helperTextDisabled: 'A secret’s version must be manually deleted.',
          helperTextEnabled: 'Delete all new versions of this secret after',
        })
      );
    } else if (['database', 'pki'].includes(this.normalizedType)) {
      const [defaultTtl, maxTtl, managedKeys] = this.leaseConfigFields as [FormField, FormField, FormField];
      fields.push(defaultTtl, maxTtl);
      if (this.normalizedType === 'pki') {
        fields.push(managedKeys);
      }
    }

    return fields;
  }

  get optionFields() {
    const [defaultTtl, maxTtl, managedKeys] = this.leaseConfigFields as [FormField, FormField, FormField];

    if (['database', 'keymgmt'].includes(this.normalizedType)) {
      return [...this.coreOptionFields, managedKeys, ...this.standardConfigFields];
    }
    if (this.normalizedType === 'pki') {
      return [...this.coreOptionFields, ...this.standardConfigFields];
    }
    if (ALL_ENGINES.find((engine) => engine.type === this.normalizedType && engine.isWIF)?.type) {
      return [
        ...this.coreOptionFields,
        defaultTtl,
        maxTtl,
        new FormField('config.identity_token_key', undefined, {
          label: 'Identity token key',
          subText: `A named key to sign tokens. If not provided, this will default to Vault's OIDC default key.`,
          editType: 'yield',
        }),
        managedKeys,
        ...this.standardConfigFields,
      ];
    }

    const options = [...this.coreOptionFields, ...this.leaseConfigFields, ...this.standardConfigFields];
    if (['kv', 'generic'].includes(this.normalizedType)) {
      options.unshift(
        new FormField('options.version', 'number', {
          label: 'Version',
          helpText:
            'The KV Secrets Engine can operate in different modes. Version 1 is the original generic Secrets Engine the allows for storing of static key/value pairs. Version 2 added more features including data versioning, TTLs, and check and set.',
          possibleValues: [2, 1],
        })
      );
    }

    return options;
  }

  get formFieldGroups() {
    return [
      new FormFieldGroup('default', this.defaultFields),
      new FormFieldGroup('Method Options', this.optionFields),
    ];
  }
}
