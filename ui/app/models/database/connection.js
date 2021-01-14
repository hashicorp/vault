import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { apiPath } from 'vault/macros/lazy-capabilities';
import attachCapabilities from 'vault/lib/attach-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const AVAILABLE_PLUGIN_TYPES = ['MongoDB', 'MongoDBA'];

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
    defaultFormValue: '',
  }),
  verify_connection: attr('boolean', {
    defaultValue: true,
  }),
  allowed_roles: attr('array', {
    readOnly: true,
  }),

  password_policy: attr('string'),

  hosts: attr('string', {
    defaultValue: '',
  }),
  host: attr('string', {
    defaultValue: '',
  }),
  url: attr('string', {
    defaultValue: '',
  }),
  port: attr('string', {
    defaultValue: '',
  }),
  // connection_details
  username: attr('string', {
    defaultValue: '',
  }),
  password: attr('string', {
    defaultValue: '',
  }),
  connection_url: attr('string', {
    defaultValue: '',
  }),

  write_concern: attr('string', {
    defaultValue: '',
  }),
  max_open_connections: attr('string', {
    defaultValue: '',
  }),
  max_idle_connections: attr('string'),
  max_connection_lifetime: attr('string'),
  tls: attr('string', {
    label: 'TLS Certificate Key',
    helpText: 'x509 certificate for connecting to the database.',
    editType: 'file',
  }),
  tls_ca: attr('string', {
    label: 'TLS CA',
    helpText: 'x509 CA file for validating the certificate presented by the MongoDB server.',
    editType: 'file',
  }),
  root_rotation_statements: attr('array'),

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

  allGroups: computed('plugin_name', function() {
    let groups = [
      { default: ['plugin_name', 'name', 'connection_url', 'verify_connection', 'password_policy'] },
      // { 'Plugin config': [] },
      // {}
    ];
    // Get plugin options
    if (this.plugin_name === 'MongoDB') {
      groups.push({
        'Plugin config': ['username', 'password', 'write_concern', { 'TLS options': ['tls', 'tls_ca'] }],
      });
    }
    // get other options
    groups.push({ 'Root rotation statements': ['root_rotation_statements'] });
    return groups;
  }),

  fieldAttrs: computed('allFields', function() {
    const expanded = expandAttributeMeta(this, this.allFields);
    console.log(expanded);
    return expanded;
  }),
});

export default attachCapabilities(M, {
  updatePath: apiPath`${'backend'}/config/${'id'}`,
});
