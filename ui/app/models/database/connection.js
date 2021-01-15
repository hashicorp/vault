import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { apiPath } from 'vault/macros/lazy-capabilities';
import attachCapabilities from 'vault/lib/attach-capabilities';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const AVAILABLE_PLUGIN_TYPES = [
  {
    value: 'mongodb-database-plugin',
    displayName: 'MongoDB',
  },
  {
    value: 'mongodbatlas-database-plugin',
    displayName: 'MongoDBA',
  },
];

const M = Model.extend({
  // ARG TODO API docs for connection https://www.vaultproject.io/api-docs/secret/databases#configure-connection
  // URL: http://127.0.0.1:8200/v1/database/config/my-db
  backend: attr('string', {
    readOnly: true,
  }),

  name: attr('string', {
    label: 'Connection Name',
    fieldValue: 'id',
    readOnly: true,
  }),
  plugin_name: attr('string', {
    label: 'Database plugin',
    possibleValues: AVAILABLE_PLUGIN_TYPES,
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
  password: attr('string', {}),
  connection_url: attr('string', {
    subText:
      'The connection string used to connect to the database. This allows for simple templating of username and password of the root user.',
  }),

  write_concern: attr('string', {}),
  max_open_connections: attr('string', {}),
  max_idle_connections: attr('string'),
  max_connection_lifetime: attr('string'),
  tls: attr('string', {
    label: 'TLS Certificate Key',
    subText: 'x509 certificate for connecting to the database.',
    editType: 'file',
  }),
  tls_ca: attr('string', {
    label: 'TLS CA',
    subText: 'x509 CA file for validating the certificate presented by the MongoDB server.',
    editType: 'file',
  }),
  root_rotation_statements: attr({
    label: 'Root rotation statements',
    editType: 'stringArray',
  }),

  allFields: computed(function() {
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

  mainFields: computed('plugin_name', function() {
    return [
      'plugin_name',
      'name',
      'connection_url',
      'verify_connection',
      'password_policy',
      'plugin_config',
      'root_rotation_statements',
    ];
  }),

  pluginGroups: computed('plugin_name', function() {
    let groups = [{ default: ['username', 'password', 'write_concern'] }];
    // Get plugin options based on plugin
    console.log(this.plugin_name);
    groups.push({
      'TLS options': ['tls', 'tls_ca'],
    });
    // get other options
    // groups.push({ 'Root rotation statements': ['root_rotation_statements'] });
    return groups;
  }),

  pluginFieldGroups: computed('pluginGroups', function() {
    return fieldToAttrs(this, this.pluginGroups);
  }),

  // TODO: Experimental
  formSections: computed('pluginGroups', function() {
    const plugin_config = fieldToAttrs(this, this.pluginGroups);
    const rotation_config = fieldToAttrs(this, [
      {
        default: ['root_rotation_statements'],
      },
    ]);
    return [{ 'Plugin config': plugin_config }, { default: rotation_config }];
  }),

  fieldAttrs: computed('mainFields', function() {
    // Main Field Attrs only
    const expanded = expandAttributeMeta(this, this.mainFields);
    console.log({ expanded });
    return expanded;
  }),
});

export default attachCapabilities(M, {
  updatePath: apiPath`${'backend'}/config/${'id'}`,
});
