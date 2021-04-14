import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { alias, or } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const AVAILABLE_PLUGIN_TYPES = [
  {
    value: 'mongodb-database-plugin',
    displayName: 'MongoDB',
    fields: [
      { attr: 'name' },
      { attr: 'plugin_name' },
      { attr: 'password_policy' },
      { attr: 'username', group: 'pluginConfig' },
      { attr: 'password', group: 'pluginConfig' },
      { attr: 'connection_url', group: 'pluginConfig' },
      { attr: 'write_concern' },
      { attr: 'creation_statements' },
    ],
  },
];

export default Model.extend({
  backend: attr('string', {
    readOnly: true,
  }),
  name: attr('string', {
    label: 'Connection Name',
  }),
  plugin_name: attr('string', {
    label: 'Database plugin',
    possibleValues: AVAILABLE_PLUGIN_TYPES,
    noDefault: true,
  }),
  verify_connection: attr('boolean', {
    defaultValue: true,
  }),
  allowed_roles: attr('array', {
    readOnly: true,
  }),

  password_policy: attr('string', {
    editType: 'optionalText',
    subText:
      'Unless a custom policy is specified, Vault will use a default: 20 characters with at least 1 uppercase, 1 lowercase, 1 number, and 1 dash character.',
  }),

  hosts: attr('string', {}),
  host: attr('string', {}),
  url: attr('string', {}),
  port: attr('string', {}),
  // connection_details
  username: attr('string', {}),
  password: attr('string', {
    editType: 'password',
  }),
  connection_url: attr('string', {
    subText:
      'The connection string used to connect to the database. This allows for simple templating of username and password of the root user.',
  }),

  write_concern: attr('string', {
    subText: 'Optional. Must be in JSON. See our documentation for help.',
    editType: 'json',
    theme: 'hashi short',
    defaultShown: 'Default',
    // defaultValue: '# For example: { "wmode": "majority", "wtimeout": 5000 }',
  }),
  max_open_connections: attr('string', {}),
  max_idle_connections: attr('string'),
  max_connection_lifetime: attr('string'),
  tls: attr('string', {
    label: 'TLS Certificate Key',
    helpText:
      'x509 certificate for connecting to the database. This must be a PEM encoded version of the private key and the certificate combined.',
    editType: 'file',
  }),
  tls_ca: attr('string', {
    label: 'TLS CA',
    helpText:
      'x509 CA file for validating the certificate presented by the MongoDB server. Must be PEM encoded.',
    editType: 'file',
  }),
  root_rotation_statements: attr({
    subText: `The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.`,
    editType: 'stringArray',
    defaultShown: 'Default',
  }),

  allowedFields: computed(function() {
    return [
      // required
      'plugin_name',
      'name',
      // fields
      'connection_url', // * MongoDB, HanaDB, MSSQL, MySQL/MariaDB, Oracle, PostgresQL, Redshift
      'verify_connection', // default true
      'password_policy', // default ""

      // plugin config
      'username',
      'password',

      'hosts',
      'host',
      'url',
      'port',
      'write_concern',
      'max_open_connections',
      'max_idle_connections',
      'max_connection_lifetime',
      'tls',
      'tls_ca',
    ];
  }),

  // for both create and edit fields
  mainFields: computed('plugin_name', function() {
    return ['plugin_name', 'name', 'connection_url', 'verify_connection', 'password_policy', 'pluginConfig'];
  }),

  showAttrs: computed('plugin_name', function() {
    const f = [
      'name',
      'plugin_name',
      'connection_url',
      'write_concern',
      'verify_connection',
      'allowed_roles',
    ];
    return expandAttributeMeta(this, f);
  }),

  pluginFieldGroups: computed('plugin_name', function() {
    if (!this.plugin_name) {
      return null;
    }
    let groups = [{ default: ['username', 'password', 'write_concern'] }];
    // TODO: Get plugin options based on plugin
    groups.push({
      'TLS options': ['tls', 'tls_ca'],
    });
    return fieldToAttrs(this, groups);
  }),

  fieldAttrs: computed('mainFields', function() {
    // Main Field Attrs only
    return expandAttributeMeta(this, this.mainFields);
  }),

  /* CAPABILITIES */
  editConnectionPath: lazyCapabilities(apiPath`${'backend'}/config/${'id'}`, 'backend', 'id'),
  canEdit: alias('editConnectionPath.canUpdate'),
  canDelete: alias('editConnectionPath.canDelete'),
  resetConnectionPath: lazyCapabilities(apiPath`${'backend'}/reset/${'id'}`, 'backend', 'id'),
  canReset: or('resetConnectionPath.canUpdate', 'resetConnectionPath.canCreate'),
  rotateRootPath: lazyCapabilities(apiPath`${'backend'}/rotate-root/${'id'}`, 'backend', 'id'),
  canRotateRoot: or('rotateRootPath.canUpdate', 'rotateRootPath.canCreate'),
  rolePath: lazyCapabilities(apiPath`${'backend'}/role/*`, 'backend'),
  staticRolePath: lazyCapabilities(apiPath`${'backend'}/static-role/*`, 'backend'),
  canAddRole: or('rolePath.canCreate', 'staticRolePath.canCreate'),
});
