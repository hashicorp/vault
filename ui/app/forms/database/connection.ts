/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';
import { AVAILABLE_PLUGIN_TYPES } from 'vault/utils/model-helpers/database-helpers';

import type { FormOptions } from 'vault/forms/form';

export type DatabaseConnectionData = {
  name: string;
  backend: string;
  plugin_name?: string;
  allowed_roles?: string[];
  [key: string]: unknown;
};

type PluginField = {
  attr: string;
  group?: string;
  subgroup?: string;
  show?: boolean;
  isEnterprise?: boolean;
};

type DatabaseConnectionFormOptions = FormOptions & {
  // enterprise-only plugin fields are hidden on community edition
  isEnterprise?: boolean;
};

// Field definitions for every attribute a database connection may have. Which of these are
// relevant for a given connection is determined by the selected plugin (see AVAILABLE_PLUGIN_TYPES).
export const CONNECTION_FIELDS: Record<string, FormField> = {
  backend: new FormField('backend', 'string'),
  name: new FormField('name', 'string', { label: 'Connection name' }),
  plugin_name: new FormField('plugin_name', 'string', {
    label: 'Database plugin',
    possibleValues: AVAILABLE_PLUGIN_TYPES,
    noDefault: true,
  }),
  verify_connection: new FormField('verify_connection', 'boolean', {
    label: 'Connection will be verified',
    defaultValue: true,
  }),
  allowed_roles: new FormField('allowed_roles', 'array'),
  password_policy: new FormField('password_policy', 'string', {
    label: 'Use custom password policy',
    editType: 'optionalText',
    subText: 'Specify the name of an existing password policy.',
    defaultSubText:
      'Unless a custom policy is specified, Vault will use a default: 20 characters with at least 1 uppercase, 1 lowercase, 1 number, and 1 dash character.',
    defaultShown: 'Default',
    docLink: '/vault/docs/concepts/password-policies',
  }),
  connection_url: new FormField('connection_url', 'string', {
    label: 'Connection URL',
    subText:
      'The connection string used to connect to the database. This allows for simple templating of username and password of the root user in the {{field_name}} format.',
  }),
  url: new FormField('url', 'string', {
    label: 'URL',
    subText: `The URL for Elasticsearch's API ("https://localhost:9200").`,
  }),
  username: new FormField('username', 'string', {
    subText: `The name of the user to use as the "root" user when connecting to the database.`,
  }),
  password: new FormField('password', 'string', {
    subText: 'The password to use when connecting with the above username.',
    editType: 'password',
  }),
  disable_escaping: new FormField('disable_escaping', 'boolean', {
    defaultValue: false,
    subText: 'Turns off the escaping of special characters inside of the username and password fields.',
    docLink: '/vault/docs/secrets/databases#disable-character-escaping',
  }),
  ca_cert: new FormField('ca_cert', 'string', {
    label: 'CA certificate',
    subText: `The path to a PEM-encoded CA cert file to use to verify the Elasticsearch server's identity.`,
  }),
  ca_path: new FormField('ca_path', 'string', {
    label: 'CA path',
    subText: `The path to a directory of PEM-encoded CA cert files to use to verify the Elasticsearch server's identity.`,
  }),
  client_cert: new FormField('client_cert', 'string', {
    label: 'Client certificate',
    subText: 'The path to the certificate for the Elasticsearch client to present for communication.',
  }),
  client_key: new FormField('client_key', 'string', {
    subText: 'The path to the key for the Elasticsearch client to use for communication.',
  }),
  hosts: new FormField('hosts', 'string'),
  host: new FormField('host', 'string'),
  port: new FormField('port', 'string'),
  write_concern: new FormField('write_concern', 'string', {
    subText: 'Optional. Must be in JSON. See our documentation for help.',
    allowReset: true,
    editType: 'json',
    theme: 'hashi short',
    defaultShown: 'Default',
  }),
  username_template: new FormField('username_template', 'string', {
    editType: 'optionalText',
    subText: 'Enter the custom username template to use.',
    defaultSubText:
      'Template describing how dynamic usernames are generated. Vault will use the default for this plugin.',
    docLink: '/vault/docs/concepts/username-templating',
    defaultShown: 'Default',
  }),
  max_open_connections: new FormField('max_open_connections', 'number', { defaultValue: 4 }),
  max_idle_connections: new FormField('max_idle_connections', 'number', { defaultValue: 0 }),
  max_connection_lifetime: new FormField('max_connection_lifetime', 'string', { defaultValue: '0s' }),
  insecure: new FormField('insecure', 'boolean', {
    label: 'Disable SSL verification',
    defaultValue: false,
  }),
  password_authentication: new FormField('password_authentication', 'string', {
    defaultValue: 'password',
    editType: 'radio',
    subText: 'The default is "password."',
    possibleValues: [
      {
        value: 'password',
        helpText:
          'Passwords will be sent to PostgreSQL in plaintext format and may appear in PostgreSQL logs as-is.',
      },
      {
        value: 'scram-sha-256',
        helpText:
          'When set to "scram-sha-256", passwords will be hashed by Vault and stored as-is by PostgreSQL. Using "scram-sha-256" requires a minimum version of PostgreSQL 10.',
      },
    ],
    docLink: '/vault/api-docs/secret/databases/postgresql#password_authentication',
  }),
  auth_type: new FormField('auth_type', 'string', {
    subText: 'If set to "gcp_iam", will enable IAM authentication to a Google CloudSQL instance.',
    docLink: '/vault/api-docs/secret/databases/postgresql#auth_type',
  }),
  service_account_json: new FormField('service_account_json', 'string', {
    label: 'Service account JSON',
    subText:
      'JSON encoded credentials for a GCP Service Account to use for IAM authentication. Requires "auth_type" to be "gcp_iam".',
    editType: 'file',
  }),
  use_private_ip: new FormField('use_private_ip', 'boolean', {
    label: 'Use private IP',
    subText:
      'Enables the option to connect to CloudSQL Instances with Private IP. Requires auth_type to be "gcp_iam".',
    defaultValue: false,
  }),
  private_key: new FormField('private_key', 'string', {
    helpText: 'The secret key used for the x509 client certificate. Must be PEM encoded.',
    editType: 'file',
  }),
  tls: new FormField('tls', 'string', {
    label: 'TLS Certificate Key',
    helpText:
      'x509 certificate for connecting to the database. This must be a PEM encoded version of the private key and the certificate combined.',
    editType: 'file',
  }),
  tls_certificate: new FormField('tls_certificate', 'string', {
    label: 'TLS Certificate Key',
    helpText: 'The x509 certificate for connecting to the database. Must be PEM encoded.',
    editType: 'file',
  }),
  tls_ca: new FormField('tls_ca', 'string', {
    label: 'TLS CA',
    helpText:
      'x509 CA file for validating the certificate presented by the database server. Must be PEM encoded.',
    editType: 'file',
  }),
  tls_server_name: new FormField('tls_server_name', 'string', {
    label: 'TLS server name',
    subText: 'If set, this name is used to set the SNI host when connecting via 1TLS.',
  }),
  root_rotation_statements: new FormField('root_rotation_statements', undefined, {
    subText: `The database statements to be executed to rotate the root user's credentials. If nothing is entered, Vault will use a reasonable default.`,
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  // ENTERPRISE ONLY
  skip_static_role_import_rotation: new FormField('skip_static_role_import_rotation', undefined, {
    editType: 'toggleButton',
    label: 'Rotate static roles immediately',
    helperTextEnabled: 'Vault automatically rotates static roles upon their initial creation.',
    helperTextDisabled: 'Vault will not automatically rotate static role passwords upon creation.',
    defaultValue: false,
    isOppositeValue: true,
  }),
  self_managed: new FormField('self_managed', 'boolean', {
    subText:
      'Allows onboarding static roles with a rootless connection configuration. Mutually exclusive with username and password. If true, will force verify_connection to be false.',
    defaultValue: false,
  }),
};

// Values the form should start with, taken from each field's default.
export const CONNECTION_DEFAULTS = Object.values(CONNECTION_FIELDS).reduce(
  (defaults: Record<string, unknown>, field) => {
    if (field.options.defaultValue !== undefined) {
      defaults[field.name] = field.options.defaultValue;
    }
    return defaults;
  },
  {}
);

// Fields shown before a plugin has been selected.
const NO_PLUGIN_FIELDS = ['plugin_name', 'name', 'connection_url', 'verify_connection', 'password_policy'];

export function isAvailablePlugin(pluginName?: string) {
  return !!AVAILABLE_PLUGIN_TYPES.find((plugin) => plugin.value === pluginName);
}

// Returns the form fields of the selected plugin that satisfy the filter.
// Enterprise-only fields are excluded on community edition.
function getPluginFields(
  pluginName: string | undefined,
  isEnterprise: boolean,
  filter: (field: PluginField) => boolean
) {
  const plugin = AVAILABLE_PLUGIN_TYPES.find((p) => p.value === pluginName);
  const fields: PluginField[] = plugin?.fields ?? [];
  return fields
    .filter((field) => (field.isEnterprise && !isEnterprise ? false : filter(field)))
    .map((field) => CONNECTION_FIELDS[field.attr] as FormField)
    .filter(Boolean);
}

// Fields displayed on the show view. Shared with the show view so it renders the same labels as the form.
export function getShowAttrs(pluginName: string | undefined, isEnterprise: boolean) {
  return [
    ...getPluginFields(pluginName, isEnterprise, (field) => field.show !== false),
    CONNECTION_FIELDS['allowed_roles'] as FormField,
  ];
}

// Flattens a read response (the `data` of the connection configuration) into the shape used by the form
// and the show view.
export function normalizeDatabaseConnection(
  responseData: Record<string, any> | undefined, // eslint-disable-line @typescript-eslint/no-explicit-any
  name: string,
  backend: string
): DatabaseConnectionData {
  const { connection_details, root_credentials_rotate_statements, ...rest } = responseData ?? {};
  const connection = { ...rest, ...connection_details };
  if (connection.connection_url) {
    // this url can include interpolated data, such as: "{{username}}/{{password}}@localhost:1521/OraDoc.localhost"
    // these curly brackets are returned by the API encoded: "%7B%7Busername%7D%7D/%7B%7Bpassword%7D%7D@localhost:1521/OraDoc.localhost"
    // we decode here so the UI displays and submits the url in the correct format
    connection.connection_url = decodeURI(connection.connection_url);
  }
  if (root_credentials_rotate_statements) {
    connection.root_rotation_statements = root_credentials_rotate_statements;
  }
  return { ...CONNECTION_DEFAULTS, ...connection, name, backend };
}

export default class DatabaseConnectionForm extends Form<DatabaseConnectionData> {
  declare isEnterprise: boolean;

  constructor(
    data: Partial<DatabaseConnectionData> = {},
    { isEnterprise = false, ...options }: DatabaseConnectionFormOptions = {}
  ) {
    super(data, options);
    this.isEnterprise = isEnterprise;
  }

  /*
   * formFieldGroups always returns every possible field so the proxy can discover each data key
   * regardless of which plugin is selected. This is called on the raw target inside the proxy handler
   * so it must not read any proxied data properties (e.g. this.plugin_name).
   */
  get formFieldGroups() {
    return [new FormFieldGroup('default', Object.values(CONNECTION_FIELDS))];
  }

  // fields rendered for both create and edit, ahead of the plugin config section
  get fieldAttrs() {
    const { plugin_name } = this.data;
    if (!plugin_name) {
      return NO_PLUGIN_FIELDS.map((name) => CONNECTION_FIELDS[name] as FormField);
    }
    return getPluginFields(plugin_name, this.isEnterprise, (field) => !field.group);
  }

  // plugin config fields grouped by subgroup, fields without a subgroup are in "default"
  get pluginFieldGroups() {
    const { plugin_name } = this.data;
    if (!plugin_name) return null;
    const plugin = AVAILABLE_PLUGIN_TYPES.find((p) => p.value === plugin_name);
    const fields: PluginField[] = plugin?.fields ?? [];
    const bySubgroup = fields
      .filter((field) => field.group === 'pluginConfig' && (this.isEnterprise || !field.isEnterprise))
      .reduce((groups: Record<string, FormField[]>, field) => {
        const group = field.subgroup || 'default';
        (groups[group] = groups[group] || []).push(CONNECTION_FIELDS[field.attr] as FormField);
        return groups;
      }, {});
    return Object.entries(bySubgroup).map(([group, groupFields]) => new FormFieldGroup(group, groupFields));
  }

  get statementFields() {
    const { plugin_name } = this.data;
    if (!plugin_name) return [CONNECTION_FIELDS['root_rotation_statements'] as FormField];
    return getPluginFields(plugin_name, this.isEnterprise, (field) => field.group === 'statements');
  }

  /*
   * The request body for creating or updating a connection. Only attributes that belong to the selected
   * plugin (plus allowed_roles) are sent. name and backend are part of the request url so they are omitted.
   */
  toJSON(data = this.data) {
    const plugin = AVAILABLE_PLUGIN_TYPES.find((p) => p.value === data.plugin_name);
    const pluginAttrs = plugin ? plugin.fields.map((field: PluginField) => field.attr) : null;
    const payload = Object.fromEntries(
      Object.entries(data).filter(([key, value]) => {
        if (value === undefined || key === 'name' || key === 'backend') return false;
        return !pluginAttrs || key === 'allowed_roles' || pluginAttrs.includes(key);
      })
    );
    return super.toJSON(payload as DatabaseConnectionData);
  }
}
