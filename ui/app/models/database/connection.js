import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import fieldToAttrs, { expandAttributeMeta } from 'vault/utils/field-to-attrs';

const AVAILABLE_PLUGIN_TYPES = [
  {
    value: 'mongodb-database-plugin',
    displayName: 'MongoDB',
    fields: [
      { attr: 'plugin_name' },
      { attr: 'name' },
      { attr: 'connection_url' },
      { attr: 'verify_connection' },
      { attr: 'password_policy' },
      { attr: 'username', group: 'pluginConfig' },
      { attr: 'password', group: 'pluginConfig', show: false },
      { attr: 'connection_url', group: 'pluginConfig' },
      { attr: 'write_concern', group: 'pluginConfig' },
      { attr: 'tls', group: 'pluginConfig', subgroup: 'TLS options' },
      { attr: 'tls_ca', group: 'pluginConfig', subgroup: 'TLS options' },
    ],
  },
  {
    value: 'mssql-database-plugin',
    displayName: 'MSSQL',
    fields: [
      { attr: 'plugin_name' },
      { attr: 'name' },
      { attr: 'connection_url' },
      { attr: 'verify_connection' },
      { attr: 'password_policy' },
      { attr: 'username', group: 'pluginConfig' },
      { attr: 'password', group: 'pluginConfig', show: false },
      { attr: 'max_open_connections', group: 'pluginConfig' },
      { attr: 'max_idle_connections', group: 'pluginConfig' },
      { attr: 'max_connection_lifetime', group: 'pluginConfig' },
    ],
  },
];

/**
 * fieldsToGroups helper fn
 * @param {array} arr any subset of "fields" from AVAILABLE_PLUGIN_TYPES
 * @param {*} key item by which to group the fields. If item has no group it will be under "default"
 * @returns array of objects where the key is default or the name of the option group, and the value is an array of attr names
 */
const fieldsToGroups = function(arr, key = 'subgroup') {
  const fieldGroups = [];
  const byGroup = arr.reduce(function(rv, x) {
    (rv[x[key]] = rv[x[key]] || []).push(x);
    return rv;
  }, {});
  Object.keys(byGroup).forEach(key => {
    const attrsArray = byGroup[key].map(obj => obj.attr);
    const group = key === 'undefined' ? 'default' : key;
    fieldGroups.push({ [group]: attrsArray });
  });
  return fieldGroups;
};

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
  }),
  verify_connection: attr('boolean', {
    label: 'Connection will be verified',
    defaultValue: true,
  }),
  allowed_roles: attr('array', {
    readOnly: true,
  }),

  password_policy: attr('string', {
    label: 'Use custom password policy',
    editType: 'optionalText',
    subText:
      'Unless a custom policy is specified, Vault will use a default: 20 characters with at least 1 uppercase, 1 lowercase, 1 number, and 1 dash character.',
  }),

  hosts: attr('string', {}),
  host: attr('string', {}),
  url: attr('string', {}),
  port: attr('string', {}),

  username: attr('string', {
    subText: 'Optional. The name of the user to use as the "root" user when connecting to the database.',
  }),
  password: attr('string', {
    subText:
      'Optional. The password to use when connecting to the database. Typically used in the connection_url field via the templating directive {{password}}.',
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
  max_open_connections: attr('number', {
    defaultValue: 4,
  }),
  max_idle_connections: attr('number', {
    defaultValue: 0,
  }),
  max_connection_lifetime: attr('string', {
    defaultValue: '0s',
  }),
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
    subText: `The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.`,
    editType: 'stringArray',
    defaultShown: 'Default',
  }),

  showAttrs: computed('plugin_name', function() {
    const fields = AVAILABLE_PLUGIN_TYPES.find(a => a.value === this.plugin_name)
      .fields.filter(f => f.show !== false)
      .map(f => f.attr);
    fields.push('allowed_roles');
    return expandAttributeMeta(this, fields);
  }),

  fieldAttrs: computed('plugin_name', function() {
    // for both create and edit fields
    let fields = ['plugin_name', 'name', 'connection_url', 'verify_connection', 'password_policy'];
    if (this.plugin_name) {
      fields = AVAILABLE_PLUGIN_TYPES.find(a => a.value === this.plugin_name)
        .fields.filter(f => !f.group)
        .map(field => field.attr);
    }
    return expandAttributeMeta(this, fields);
  }),

  pluginFieldGroups: computed('plugin_name', function() {
    if (!this.plugin_name) {
      return null;
    }
    let pluginFields = AVAILABLE_PLUGIN_TYPES.find(a => a.value === this.plugin_name).fields.filter(
      f => f.group === 'pluginConfig'
    );
    let groups = fieldsToGroups(pluginFields, 'subgroup');
    return fieldToAttrs(this, groups);
  }),

  statementFields: computed('plugin_name', function() {
    if (!this.plugin_name) {
      return null;
    }
    let fields = AVAILABLE_PLUGIN_TYPES.find(a => a.value === this.plugin_name)
      .fields.filter(f => f.group === 'statements')
      .map(field => field.attr);
    return expandAttributeMeta(this, fields);
  }),

  /* CAPABILITIES */
  editConnectionPath: lazyCapabilities(apiPath`${'backend'}/config/${'id'}`, 'backend', 'id'),
  canEdit: alias('editConnectionPath.canUpdate'),
  canDelete: alias('editConnectionPath.canDelete'),
  resetConnectionPath: lazyCapabilities(apiPath`${'backend'}/reset/${'id'}`, 'backend', 'id'),
  canReset: computed.or('resetConnectionPath.canUpdate', 'resetConnectionPath.canCreate'),
  rotateRootPath: lazyCapabilities(apiPath`${'backend'}/rotate-root/${'id'}`, 'backend', 'id'),
  canRotateRoot: computed.or('rotateRootPath.canUpdate', 'rotateRootPath.canCreate'),
  rolePath: lazyCapabilities(apiPath`${'backend'}/role/*`, 'backend'),
  staticRolePath: lazyCapabilities(apiPath`${'backend'}/static-role/*`, 'backend'),
  canAddRole: computed.or('rolePath.canCreate', 'staticRolePath.canCreate'),
});
