/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { WHITESPACE_WARNING } from 'vault/utils/forms/validators';
import { tracked } from '@glimmer/tracking';
import { ALL_ENGINES } from 'vault/utils/all-engines-metadata';

import type { SecretsEngineFormData } from 'vault/secrets/engine';
import type { Validations } from 'vault/app-types';

export default class SecretsEngineForm extends Form<SecretsEngineFormData> {
  @tracked declare type: string;

  validations: Validations = {
    path: [
      { type: 'presence', message: "Path can't be blank." },
      {
        type: 'containsWhiteSpace',
        message: WHITESPACE_WARNING('path'),
        level: 'warn',
      },
    ],
    'kvConfig.maxVersions': [
      { type: 'number', message: 'Maximum versions must be a number.' },
      { type: 'length', options: { min: 1, max: 16 }, message: 'You cannot go over 16 characters.' },
    ],
  };

  get coreOptionFields() {
    return [
      new FormField('description', 'string', { editType: 'textarea' }),
      new FormField('config.listingVisibility', 'boolean', {
        label: 'Use as preferred UI login method',
        editType: 'toggleButton',
        helperTextEnabled:
          'This mount will be included in the unauthenticated UI login endpoint and display as a preferred login method.',
        helperTextDisabled:
          'Turn on the toggle to use this auth mount as a preferred login method during UI login.',
      }),
      new FormField('local', 'boolean', {
        helpText:
          'When Replication is enabled, a local mount will not be replicated across clusters. This can only be specified at mount time.',
      }),
      new FormField('sealWrap', 'boolean', {
        helpText:
          'When enabled - if a seal supporting seal wrapping is specified in the configuration, all critical security parameters (CSPs) in this backend will be seal wrapped. (For KV mounts, all values will be seal wrapped.) This can only be specified at mount time.',
      }),
    ];
  }

  get leaseConfigFields() {
    return [
      new FormField('config.defaultLeaseTtl', 'string', { label: 'Default Lease TTL', editType: 'ttl' }),
      new FormField('config.maxLeaseTtl', 'string', { label: 'Max Lease TTL', editType: 'ttl' }),
      new FormField('config.allowedManagedKeys', 'string', {
        label: 'Allowed managed keys',
        editType: 'stringArray',
      }),
    ];
  }

  get standardConfigFields() {
    return [
      new FormField('config.auditNonHmacRequestKeys', 'string', {
        label: 'Request keys excluded from HMACing in audit',
        editType: 'stringArray',
        helpText: "Keys that will not be HMAC'd by audit devices in the request data object.",
      }),
      new FormField('config.auditNonHmacResponseKeys', 'string', {
        label: 'Response keys excluded from HMACing in audit',
        editType: 'stringArray',
        helpText: "Keys that will not be HMAC'd by audit devices in the response data object.",
      }),
      new FormField('config.passthroughRequestHeaders', 'string', {
        label: 'Allowed passthrough request headers',
        helpText: 'Headers to allow and pass from the request to the backend',
        editType: 'stringArray',
      }),
      new FormField('config.allowedResponseHeaders', 'string', {
        label: 'Allowed response headers',
        helpText: 'Headers to allow, allowing a plugin to include them in the response.',
        editType: 'stringArray',
      }),
    ];
  }

  get engineType() {
    return (this.type || '').replace(/^ns_/, '');
  }

  get defaultFields() {
    const fields = [new FormField('path', 'string')];
    if (this.engineType === 'kv') {
      fields.push(
        new FormField('kvConfig.maxVersions', 'number', {
          label: 'Maximum number of versions',
          subText:
            'The number of versions to keep per key. Once the number of keys exceeds the maximum number set here, the oldest version will be permanently deleted. This value applies to all keys, but a key’s metadata settings can overwrite this value. When 0 is used or the value is unset, Vault will keep 10 versions.',
        }),
        new FormField('kvConfig.casRequired', 'boolean', {
          label: 'Require Check and Set',
          subText:
            'If checked, all keys will require the cas parameter to be set on all write requests. A key’s metadata settings can overwrite this value.',
        }),
        new FormField('kvConfig.deleteVersionAfter', 'string', {
          editType: 'ttl',
          label: 'Automate secret deletion',
          helperTextDisabled: 'A secret’s version must be manually deleted.',
          helperTextEnabled: 'Delete all new versions of this secret after',
        })
      );
    } else if (['database', 'pki'].includes(this.engineType)) {
      const [defaultTtl, maxTtl, managedKeys] = this.leaseConfigFields as [FormField, FormField, FormField];
      fields.push(defaultTtl, maxTtl);
      if (this.engineType === 'pki') {
        fields.push(managedKeys);
      }
    }

    return fields;
  }

  get optionFields() {
    const [defaultTtl, maxTtl, managedKeys] = this.leaseConfigFields as [FormField, FormField, FormField];

    if (['database', 'keymgmt'].includes(this.engineType)) {
      return [...this.coreOptionFields, managedKeys, ...this.standardConfigFields];
    }
    if (this.engineType === 'pki') {
      return [...this.coreOptionFields, ...this.standardConfigFields];
    }
    if (ALL_ENGINES.find((engine) => engine.type === this.engineType && engine.isWIF)?.type) {
      return [
        ...this.coreOptionFields,
        defaultTtl,
        maxTtl,
        new FormField('config.identityTokenKey', undefined, {
          label: 'Identity token key',
          subText: `A named key to sign tokens. If not provided, this will default to Vault's OIDC default key.`,
          editType: 'yield',
        }),
        managedKeys,
        ...this.standardConfigFields,
      ];
    }

    const options = [...this.coreOptionFields, ...this.leaseConfigFields, ...this.standardConfigFields];
    if (['kv', 'generic'].includes(this.engineType)) {
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

  toJSON() {
    const { config } = this.data;
    const data = {
      type: this.type,
      ...this.data,
      config: {
        ...(config || {}),
        forceNoCache: config?.forceNoCache ?? false,
        listingVisibility: config?.listingVisibility ? 'unauth' : 'hidden',
      },
    };
    // options are only relevant for kv/generic engines
    if (!['kv', 'generic'].includes(this.type)) {
      delete data.options;
    }

    return super.toJSON(data);
  }
}
