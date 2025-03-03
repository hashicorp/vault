/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { getRoleFields } from 'vault/utils/model-helpers/database-helpers';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';
const validations = {
  database: [{ type: 'presence', message: 'Database is required.' }],
  type: [{ type: 'presence', message: 'Type is required.' }],
  username: [
    {
      validator(model) {
        const { type, username } = model;
        if (!type || type === 'dynamic') return true;
        if (username) return true;
      },
      message: 'Username is required.',
    },
  ],
};
@withModelValidations(validations)
export default class RoleModel extends Model {
  idPrefix = 'role/';
  @attr('string', { readOnly: true }) backend;
  @attr('string', { label: 'Role name' }) name;
  @attr('array', {
    label: 'Connection name',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['database/connection'],
    selectLimit: 1,
    onlyAllowExisting: true,
    subText: 'The database connection for which credentials will be generated.',
  })
  database;
  @attr('string', {
    label: 'Type of role',
    noDefault: true,
    possibleValues: ['static', 'dynamic'],
  })
  type;
  @attr({
    editType: 'ttl',
    defaultValue: '1h',
    label: 'Generated credentials’s Time-to-Live (TTL)',
    helperTextDisabled: 'Vault will use a TTL of 1 hour.',
    defaultShown: 'Engine default',
  })
  default_ttl;
  @attr({
    editType: 'ttl',
    defaultValue: '24h',
    label: 'Generated credentials’s maximum Time-to-Live (Max TTL)',
    helperTextDisabled: 'Vault will use a TTL of 24 hours.',
    defaultShown: 'Engine default',
  })
  max_ttl;
  @attr('string', { subText: 'The database username that this Vault role corresponds to.' }) username;
  @attr({
    editType: 'ttl',
    defaultValue: '24h',
    helperTextDisabled:
      'Specifies the amount of time Vault should wait before rotating the password. The minimum is 5 seconds. Default is 24 hours.',
    helperTextEnabled: 'Vault will rotate password after.',
  })
  rotation_period;
  @attr({
    label: 'Skip initial rotation',
    editType: 'boolean',
    defaultValue: false,
    subText: 'When unchecked, Vault automatically rotates the password upon creation.',
  })
  skip_import_rotation;
  @attr('array', {
    editType: 'stringArray',
  })
  creation_statements;
  @attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  })
  revocation_statements;
  @attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  })
  rotation_statements;
  @attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  })
  rollback_statements;
  @attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  })
  renew_statements;
  @attr('string', {
    editType: 'json',
    allowReset: true,
    theme: 'hashi short',
    defaultShown: 'Default',
  })
  creation_statement;
  @attr('string', {
    editType: 'json',
    allowReset: true,
    theme: 'hashi short',
    defaultShown: 'Default',
  })
  revocation_statement;
  /* FIELD ATTRIBUTES */
  get fieldAttrs() {
    // Main fields on edit/create form
    const fields = ['name', 'database', 'type'];
    return expandAttributeMeta(this, fields);
  }
  get showFields() {
    let fields = ['name', 'database', 'type'];
    fields = fields.concat(getRoleFields(this.type)).concat(['creation_statements']);
    // elasticsearch does not support revocation statements: https://developer.hashicorp.com/vault/api-docs/secret/databases/elasticdb#parameters-1
    if (this.database[0] !== 'elasticsearch') {
      fields = fields.concat(['revocation_statements']);
    }
    return expandAttributeMeta(this, fields);
  }
  get roleSettingAttrs() {
    // logic for which get displayed is on DatabaseRoleSettingForm
    const allRoleSettingFields = [
      'default_ttl',
      'max_ttl',
      'username',
      'rotation_period',
      'skip_import_rotation',
      'creation_statements',
      'creation_statement', // for editType: JSON
      'revocation_statements',
      'revocation_statement', // only for MongoDB (editType: JSON)
      'rotation_statements',
      'rollback_statements',
      'renew_statements',
    ];
    return expandAttributeMeta(this, allRoleSettingFields);
  }
  /* CAPABILITIES */
  // only used for secretPath
  @attr('string', { readOnly: true }) path;
  @lazyCapabilities(apiPath`${'backend'}/${'path'}/${'id'}`, 'backend', 'path', 'id') secretPath;
  @lazyCapabilities(apiPath`${'backend'}/roles/+`, 'backend') dynamicPath;
  @lazyCapabilities(apiPath`${'backend'}/static-roles/+`, 'backend') staticPath;
  @lazyCapabilities(apiPath`${'backend'}/creds/${'id'}`, 'backend', 'id') credentialPath;
  @lazyCapabilities(apiPath`${'backend'}/static-creds/${'id'}`, 'backend', 'id') staticCredentialPath;
  @lazyCapabilities(apiPath`${'backend'}/config/${'database[0]'}`, 'backend', 'database') databasePath;
  @lazyCapabilities(apiPath`${'backend'}/rotate-role/${'id'}`, 'backend', 'id') rotateRolePath;

  get canEditRole() {
    return this.secretPath.get('canUpdate');
  }
  get canDelete() {
    return this.secretPath.get('canDelete');
  }
  get canCreateDynamic() {
    return this.dynamicPath.get('canCreate');
  }
  get canCreateStatic() {
    return this.staticPath.get('canCreate');
  }
  get canGenerateCredentials() {
    return this.credentialPath.get('canRead');
  }
  get canGetCredentials() {
    return this.staticCredentialPath.get('canRead');
  }
  get canUpdateDb() {
    return this.databasePath.get('canUpdate');
  }
  get canRotateRoleCredentials() {
    return this.rotateRolePath.get('canUpdate');
  }
}
