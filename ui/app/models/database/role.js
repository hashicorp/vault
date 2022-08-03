import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { getRoleFields } from 'vault/utils/database-helpers';

export default Model.extend({
  idPrefix: 'role/',
  backend: attr('string', { readOnly: true }),
  name: attr('string', {
    label: 'Role name',
  }),
  database: attr('array', {
    label: '',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['database/connection'],
    selectLimit: 1,
    onlyAllowExisting: true,
    subLabel: 'Database name',
    subText: 'The database for which credentials will be generated.',
  }),
  type: attr('string', {
    label: 'Type of role',
    noDefault: true,
    possibleValues: ['static', 'dynamic'],
  }),
  default_ttl: attr({
    editType: 'ttl',
    defaultValue: '1h',
    label: 'Generated credentials’s Time-to-Live (TTL)',
    helperTextDisabled: 'Vault will use a TTL of 1 hour.',
    defaultShown: 'Engine default',
  }),
  max_ttl: attr({
    editType: 'ttl',
    defaultValue: '24h',
    label: 'Generated credentials’s maximum Time-to-Live (Max TTL)',
    helperTextDisabled: 'Vault will use a TTL of 24 hours.',
    defaultShown: 'Engine default',
  }),
  username: attr('string', {
    subText: 'The database username that this Vault role corresponds to.',
  }),
  rotation_period: attr({
    editType: 'ttl',
    defaultValue: '24h',
    helperTextDisabled:
      'Specifies the amount of time Vault should wait before rotating the password. The minimum is 5 seconds. Default is 24 hours.',
    helperTextEnabled: 'Vault will rotate password after',
  }),
  creation_statements: attr('array', {
    editType: 'stringArray',
  }),
  revocation_statements: attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  rotation_statements: attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  rollback_statements: attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  renew_statements: attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  creation_statement: attr('string', {
    editType: 'json',
    allowReset: true,
    theme: 'hashi short',
    defaultShown: 'Default',
  }),
  revocation_statement: attr('string', {
    editType: 'json',
    allowReset: true,
    theme: 'hashi short',
    defaultShown: 'Default',
  }),

  /* FIELD ATTRIBUTES */
  get fieldAttrs() {
    // Main fields on edit/create form
    let fields = ['name', 'database', 'type'];
    return expandAttributeMeta(this, fields);
  },

  get showFields() {
    let fields = ['name', 'database', 'type'];
    fields = fields.concat(getRoleFields(this.type)).concat(['creation_statements']);
    // elasticsearch does not support revocation statements: https://www.vaultproject.io/api-docs/secret/databases/elasticdb#parameters-1
    if (this.database[0] !== 'elasticsearch') {
      fields = fields.concat(['revocation_statements']);
    }
    return expandAttributeMeta(this, fields);
  },

  roleSettingAttrs: computed(function () {
    // logic for which get displayed is on DatabaseRoleSettingForm
    let allRoleSettingFields = [
      'default_ttl',
      'max_ttl',
      'username',
      'rotation_period',
      'creation_statements',
      'creation_statement', // for editType: JSON
      'revocation_statements',
      'revocation_statement', // only for MongoDB (editType: JSON)
      'rotation_statements',
      'rollback_statements',
      'renew_statements',
    ];
    return expandAttributeMeta(this, allRoleSettingFields);
  }),

  /* CAPABILITIES */
  // only used for secretPath
  path: attr('string', { readOnly: true }),

  secretPath: lazyCapabilities(apiPath`${'backend'}/${'path'}/${'id'}`, 'backend', 'path', 'id'),
  canEditRole: alias('secretPath.canUpdate'),
  canDelete: alias('secretPath.canDelete'),
  dynamicPath: lazyCapabilities(apiPath`${'backend'}/roles/+`, 'backend'),
  canCreateDynamic: alias('dynamicPath.canCreate'),
  staticPath: lazyCapabilities(apiPath`${'backend'}/static-roles/+`, 'backend'),
  canCreateStatic: alias('staticPath.canCreate'),
  credentialPath: lazyCapabilities(apiPath`${'backend'}/creds/${'id'}`, 'backend', 'id'),
  staticCredentialPath: lazyCapabilities(apiPath`${'backend'}/static-creds/${'id'}`, 'backend', 'id'),
  canGenerateCredentials: alias('credentialPath.canRead'),
  canGetCredentials: alias('staticCredentialPath.canRead'),
  databasePath: lazyCapabilities(apiPath`${'backend'}/config/${'database[0]'}`, 'backend', 'database'),
  canUpdateDb: alias('databasePath.canUpdate'),
  rotateRolePath: lazyCapabilities(apiPath`${'backend'}/rotate-role/${'id'}`, 'backend', 'id'),
  canRotateRoleCredentials: alias('rotateRolePath.canUpdate'),
});
