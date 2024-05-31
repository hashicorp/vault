/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { belongsTo, hasMany, attr } from '@ember-data/model';
import { service } from '@ember/service';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import apiPath from 'vault/utils/api-path';
import { withModelValidations } from 'vault/decorators/model-validations';
import { allMethods } from 'vault/helpers/mountable-auth-methods';
import lazyCapabilities from 'vault/macros/lazy-capabilities';
import { action } from '@ember/object';
import { camelize } from '@ember/string';

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
};

@withModelValidations(validations)
export default class AuthMethodModel extends Model {
  @service store;

  @belongsTo('mount-config', { async: false, inverse: null }) config; // one-to-none that replaces former fragment
  @hasMany('auth-config', { polymorphic: true, inverse: 'backend', async: false }) authConfigs;
  @attr('string') path;
  @attr('string') accessor;
  @attr('string') name;
  @attr('string') type;
  // namespaces introduced types with a `ns_` prefix for built-in engines
  // so we need to strip that to normalize the type
  get methodType() {
    return this.type.replace(/^ns_/, '');
  }
  get icon() {
    const authMethods = allMethods().find((backend) => backend.type === this.methodType);

    return authMethods?.glyph || 'users';
  }
  @attr('string', {
    editType: 'textarea',
  })
  description;
  @attr('boolean', {
    helpText:
      'When Replication is enabled, a local mount will not be replicated across clusters. This can only be specified at mount time.',
  })
  local;
  @attr('boolean', {
    helpText:
      'When enabled - if a seal supporting seal wrapping is specified in the configuration, all critical security parameters (CSPs) in this backend will be seal wrapped. (For KV mounts, all values will be seal wrapped.) This can only be specified at mount time.',
  })
  sealWrap;

  // used when the `auth` prefix is important,
  // currently only when setting perf mount filtering
  get apiPath() {
    return `auth/${this.path}`;
  }
  get localDisplay() {
    return this.local ? 'local' : 'replicated';
  }

  get supportsUserLockoutConfig() {
    return ['approle', 'ldap', 'userpass'].includes(this.methodType);
  }

  userLockoutConfig = {
    modelAttrs: [
      'config.lockoutThreshold',
      'config.lockoutDuration',
      'config.lockoutCounterReset',
      'config.lockoutDisable',
    ],
    apiParams: ['lockout_threshold', 'lockout_duration', 'lockout_counter_reset', 'lockout_disable'],
  };

  get tuneAttrs() {
    // order here determines order tune fields render
    const tuneAttrs = [
      'listingVisibility',
      'defaultLeaseTtl',
      'maxLeaseTtl',
      ...(this.methodType === 'token' ? [] : ['tokenType']),
      'auditNonHmacRequestKeys',
      'auditNonHmacResponseKeys',
      'passthroughRequestHeaders',
      'allowedResponseHeaders',
      'pluginVersion',
      ...(this.supportsUserLockoutConfig ? this.userLockoutConfig.apiParams.map((a) => camelize(a)) : []),
    ];

    return expandAttributeMeta(this, ['description', `config.{${tuneAttrs.join(',')}}`]);
  }

  get formFields() {
    return [
      'type',
      'path',
      'description',
      'accessor',
      'local',
      'sealWrap',
      'config.{listingVisibility,defaultLeaseTtl,maxLeaseTtl,tokenType,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders,pluginVersion}',
    ];
  }

  get formFieldGroups() {
    return [
      { default: ['path'] },
      {
        'Method Options': [
          'description',
          'config.listingVisibility',
          'local',
          'sealWrap',
          'config.{defaultLeaseTtl,maxLeaseTtl,tokenType,auditNonHmacRequestKeys,auditNonHmacResponseKeys,passthroughRequestHeaders,allowedResponseHeaders,pluginVersion}',
        ],
      },
    ];
  }

  get attrs() {
    return expandAttributeMeta(this, this.formFields);
  }

  get fieldGroups() {
    return fieldToAttrs(this, this.formFieldGroups);
  }
  @lazyCapabilities(apiPath`sys/auth/${'id'}`, 'id') deletePath;
  @lazyCapabilities(apiPath`auth/${'id'}/config`, 'id') configPath;
  @lazyCapabilities(apiPath`auth/${'id'}/config/client`, 'id') awsConfigPath;
  get canDisable() {
    return this.deletePath.get('canDelete') !== false;
  }
  get canEdit() {
    return this.configPath.get('canUpdate') !== false;
  }
  get canEditAws() {
    return this.awsConfigPath.get('canUpdate') !== false;
  }

  @action
  tune(data) {
    return this.store.adapterFor('auth-method').tune(this.path, data);
  }
}
