import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

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
    possibleValues: ['static', 'dynamic'],
  }),
  ttl: attr({
    editType: 'ttl',
    defaultValue: '1h',
    label: 'Generated credentials’s Time-to-Live (TTL)',
    subText: 'Vault will use the engine default of 1 hour',
  }),
  max_ttl: attr({
    editType: 'ttl',
    defaultValue: '24h',
    label: 'Generated credentials’s maximum Time-to-Live (Max TTL)',
    subText: 'Vault will use the engine default of 24 hours',
  }),
  username: attr('string', {
    subText: 'The database username that this Vault role corresponds to.',
  }),
  rotation_period: attr({
    editType: 'ttl',
    defaultValue: '5s',
    subText:
      'Specifies the amount of time Vault should wait before rotating the password. The minimum is 5 seconds.',
  }),

  // path is used for capabilities only
  path: attr('string', { readOnly: true }),

  get fieldAttrs() {
    let fields = ['database', 'name', 'type'];
    return expandAttributeMeta(this, fields);
  },

  roleSettingAttrs: computed(function() {
    // actual display for which ones is set on DatabaseRoleSettingForm
    let allRoleSettingFields = ['ttl', 'max_ttl', 'username', 'rotation_period'];
    return expandAttributeMeta(this, allRoleSettingFields);
  }),

  statementsAttrs: computed('type', function() {
    if (!this.type) {
      return null;
    }
    let statementFields = ['creation_statements', 'revocation_statements'];
    if (this.type === 'static') {
      statementFields = ['rotation_statements'];
    }
    return expandAttributeMeta(this, statementFields);
  }),

  secretPath: lazyCapabilities(apiPath`${'backend'}/${'path'}/${'id'}`, 'backend', 'path', 'id'),
  canEditRole: alias('secretPath.canUpdate'),
  canDelete: alias('secretPath.canDelete'),
  dynamicPath: lazyCapabilities(apiPath`${'backend'}/roles/+`, 'backend'),
  canCreateDynamic: alias('dynamicPath.canCreate'),
  staticPath: lazyCapabilities(apiPath`${'backend'}/static-roles/+`, 'backend'),
  canCreateStatic: alias('staticPath.canCreate'),
  credentialPath: lazyCapabilities(apiPath`${'backend'}/creds/${'id'}`, 'backend', 'id'),
  canGenerateCredentials: alias('credentialPath.canRead'),
  databasePath: lazyCapabilities(apiPath`${'backend'}/config/${'database[0]'}`, 'backend', 'database'),
  canUpdateDb: alias('databasePath.canUpdate'),
});
