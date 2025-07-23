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
    'kv_config.max_versions': [
      { type: 'number', message: 'Maximum versions must be a number.' },
      { type: 'length', options: { min: 1, max: 16 }, message: 'You cannot go over 16 characters.' },
    ],
  };

  get coreOptionFields() {
    return [
      new FormField('description', 'string', { editType: 'textarea' }),
      new FormField('config.listing_visibility', 'boolean', {
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
      new FormField('seal_wrap', 'boolean', {
        helpText:
          'When enabled - if a seal supporting seal wrapping is specified in the configuration, all critical security parameters (CSPs) in this backend will be seal wrapped. (For KV mounts, all values will be seal wrapped.) This can only be specified at mount time.',
      }),
    ];
  }

  get leaseConfigFields() {
    return [
      new FormField('config.default_lease_ttl', 'string', { label: 'Default Lease TTL', editType: 'ttl' }),
      new FormField('config.max_lease_ttl', 'string', { label: 'Max Lease TTL', editType: 'ttl' }),
      new FormField('config.allowed_managed_keys', 'string', {
        label: 'Allowed managed keys',
        editType: 'stringArray',
      }),
    ];
  }

  get standardConfigFields() {
    return [
      new FormField('config.audit_non_hmac_request_keys', 'string', {
        label: 'Request keys excluded from HMACing in audit',
        editType: 'stringArray',
        helpText: "Keys that will not be HMAC'd by audit devices in the request data object.",
      }),
      new FormField('config.audit_non_hmac_response_keys', 'string', {
        label: 'Response keys excluded from HMACing in audit',
        editType: 'stringArray',
        helpText: "Keys that will not be HMAC'd by audit devices in the response data object.",
      }),
      new FormField('config.passthrough_request_headers', 'string', {
        label: 'Allowed passthrough request headers',
        helpText: 'Headers to allow and pass from the request to the backend',
        editType: 'stringArray',
      }),
      new FormField('config.allowed_response_headers', 'string', {
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
    } else if (['database', 'pki'].includes(this.engineType)) {
      const [defaultTtl, maxTtl, managedKeys] = this.leaseConfigFields as [FormField, FormField, FormField];
      fields.push(defaultTtl, maxTtl);
      if (this.engineType === 'pki') {
        fields.push(managedKeys);
      }
    } else if (this.engineType === 'custom-plugin') {
      const customFields = fields;
      customFields.push(
        new FormField('plugin_name', 'string', {
          subText:
            'Specifies the name for this plugin. Enterprise plugin names must match the name listed on the HashiCorp releases page ie. "vault-plugin-secrets-kv..."',
        })
      );
      customFields.push(
        new FormField('sha256', 'string', {
          subText:
            'SHA256 checksum of the plugin binary. - ex. run "shasum -a 256 {{your plugin binary}}" in CLI ',
        })
      );
      customFields.push(
        new FormField('plugin_zip', 'file', {
          label: 'Upload Plugin Zip',
          editType: 'file',
          subText:
            '- Upload a zip file executable from HashiCorp releases page or compiled from github repository',
        })
      );
      customFields.push(new FormField('description', 'string', { editType: 'textarea' }));
      return customFields;
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
    if (this.engineType === 'custom-plugin') {
      return [new FormFieldGroup('default', this.defaultFields)];
    }
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
        force_no_cache: config?.force_no_cache ?? false,
        listing_visibility: config?.listing_visibility ? 'unauth' : 'hidden',
      },
    };
    // options are only relevant for kv/generic engines
    if (!['kv', 'generic'].includes(this.type)) {
      delete data.options;
    }

    return super.toJSON(data);
  }
}
