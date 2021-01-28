// ARG TODO need to replace with Lazy capabilities, don't want to step on toes for now.
import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { apiPath } from 'vault/macros/lazy-capabilities';
import attachCapabilities from 'vault/lib/attach-capabilities';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const STANDARD_FIELDS = [
  { attr: 'name', required: true },
  { attr: 'plugin_name', required: true },
  { attr: 'verify_connection' },
  { attr: 'allowed_roles' },
];

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
  {
    value: 'mongodbatlas-database-plugin',
    displayName: 'MongoDBA',
  },
];

const M = Model.extend({
  // URL: http://127.0.0.1:8200/v1/database/config/my-db
  backend: attr('string', {
    readOnly: true,
  }),

  name: attr('string', {
    label: 'Connection Name',
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
    // defaultValue: '# For example: { "wmode": "majority", "wtimeout": 5000 }',
  }),
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
  root_rotation_statements: attr('string', {
    // label: 'Root rotation statements',
    subText: `The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.`,
    editType: 'json',
    theme: 'hashi short',
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
    return [
      'plugin_name',
      'name',
      'connection_url',
      'verify_connection',
      'password_policy',
      'pluginConfig',
      'root_rotation_statements',
    ];
  }),

  // showFields: computed('plugin_name', function() {
  //   const f = [
  //     'name',
  //     'plugin_name',
  //     'connection_url',
  //     'write_concern',
  //     'verify_connection',
  //     'root_rotation_statements',
  //     'allowed_roles',
  //   ];
  //   return fieldToAttrs(this, f);
  // }),
  showAttrs: computed('plugin_name', function() {
    const f = [
      'name',
      'plugin_name',
      'connection_url',
      'write_concern',
      'verify_connection',
      'root_rotation_statements',
      'allowed_roles',
    ];
    return expandAttributeMeta(this, f);
  }),

  // pluginGroups: computed('plugin_name', function() {
  //   let groups = [{ default: ['username', 'password', 'write_concern'] }];
  //   // TODO: Get plugin options based on plugin
  //   groups.push({
  //     'TLS options': ['tls', 'tls_ca'],
  //   });
  //   return groups;
  // }),

  pluginFieldGroups: computed('plugin_name', function() {
    let groups = [{ default: ['username', 'password', 'write_concern'] }];
    // TODO: Get plugin options based on plugin
    groups.push({
      'TLS options': ['tls', 'tls_ca'],
    });
    return fieldToAttrs(this, groups);
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
    return expandAttributeMeta(this, this.mainFields);
  }),
});

export default attachCapabilities(M, {
  updatePath: apiPath`${'backend'}/config/${'id'}`,
});
