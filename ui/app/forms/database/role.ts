/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import { getRoleFields } from 'vault/utils/model-helpers/database-helpers';

import type { FormOptions } from 'vault/forms/form';
import type { Validations } from 'vault/app-types';

export type DatabaseRoleData = {
  name: string;
  backend: string;
  type?: string;
  database?: string[];
  [key: string]: unknown;
};

// Field definitions for every attribute a database role may have. Which of these are relevant
// for a given role is determined by the role type (static or dynamic) and the connection plugin.
export const ROLE_FIELD_DEFINITIONS: Record<string, FormField> = {
  backend: new FormField('backend', 'string', { readOnly: true }),
  name: new FormField('name', 'string', { label: 'Role name' }),
  database: new FormField('database', 'array', {
    label: 'Connection name',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['database/connection'],
    selectLimit: 1,
    onlyAllowExisting: true,
    subText: 'The database connection for which credentials will be generated.',
  }),
  type: new FormField('type', 'string', {
    label: 'Type of role',
    noDefault: true,
    possibleValues: ['static', 'dynamic'],
  }),
  default_ttl: new FormField('default_ttl', undefined, {
    editType: 'ttl',
    defaultValue: '1h',
    label: 'Generated credentials’s Time-to-Live (TTL)',
    helperTextDisabled: 'Vault will use a TTL of 1 hour.',
    defaultShown: 'Engine default',
  }),
  max_ttl: new FormField('max_ttl', undefined, {
    editType: 'ttl',
    defaultValue: '24h',
    label: 'Generated credentials’s maximum Time-to-Live (Max TTL)',
    helperTextDisabled: 'Vault will use a TTL of 24 hours.',
    defaultShown: 'Engine default',
  }),
  username: new FormField('username', 'string', {
    subText: 'The database username that this Vault role corresponds to.',
  }),
  rotation_period: new FormField('rotation_period', undefined, {
    editType: 'ttl',
    defaultValue: '24h',
    helperTextDisabled:
      'Specifies the amount of time Vault should wait before rotating the password. The minimum is 5 seconds. Default is 24 hours.',
    helperTextEnabled: 'Vault will rotate password after.',
  }),
  creation_statements: new FormField('creation_statements', 'array', { editType: 'stringArray' }),
  revocation_statements: new FormField('revocation_statements', 'array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  rotation_statements: new FormField('rotation_statements', 'array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  rollback_statements: new FormField('rollback_statements', 'array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  renew_statements: new FormField('renew_statements', 'array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  // MongoDB takes a json string rather than a list of statements
  creation_statement: new FormField('creation_statement', 'string', {
    editType: 'json',
    allowReset: true,
    theme: 'hashi short',
    defaultShown: 'Default',
  }),
  revocation_statement: new FormField('revocation_statement', 'string', {
    editType: 'json',
    allowReset: true,
    theme: 'hashi short',
    defaultShown: 'Default',
  }),
  last_vault_rotation: new FormField('last_vault_rotation', 'string', { readOnly: true }),
  // ENTERPRISE ONLY
  skip_import_rotation: new FormField('skip_import_rotation', undefined, {
    label: 'Rotate immediately',
    editType: 'toggleButton',
    helperTextEnabled: 'Vault will rotate the password for this static role on creation.',
    helperTextDisabled: "Vault will not rotate this role's password on creation.",
    isOppositeValue: true,
  }),
  password: new FormField('password', 'string', {
    sensitive: true,
    subText: 'The database password that this Vault role corresponds to.',
  }),
};

const fields = (names: string[]) => names.map((name) => ROLE_FIELD_DEFINITIONS[name] as FormField);

export const ROLE_DEFAULTS = Object.values(ROLE_FIELD_DEFINITIONS).reduce(
  (defaults: Record<string, unknown>, field) => {
    if (field.options.defaultValue !== undefined) {
      defaults[field.name] = field.options.defaultValue;
    }
    return defaults;
  },
  {}
);

// Main fields on the create / edit form
export const getRoleFieldAttrs = () => fields(['name', 'database', 'type']);

// Role settings and statement fields. Which of them render is decided by DatabaseRoleSettingForm.
// Enterprise-only fields are removed on community edition.
export const getRoleSettingAttrs = (isEnterprise: boolean) => {
  const names = [
    'default_ttl',
    'max_ttl',
    'username',
    'password',
    'rotation_period',
    'skip_import_rotation',
    'creation_statements',
    'creation_statement',
    'revocation_statements',
    'revocation_statement',
    'rotation_statements',
    'rollback_statements',
    'renew_statements',
  ];
  return fields(
    isEnterprise ? names : names.filter((n) => !['skip_import_rotation', 'password'].includes(n))
  );
};

// Fields displayed on the show view
export const getRoleShowFields = (type?: string, database?: string[]) => {
  const names = ['name', 'database', 'type', ...(getRoleFields(type) ?? []), 'creation_statements'];
  // elasticsearch does not support revocation statements: https://developer.hashicorp.com/vault/api-docs/secret/databases/elasticdb#parameters-1
  if (database?.[0] !== 'elasticsearch') {
    names.push('revocation_statements');
  }
  return fields(names);
};

// Flattens a read response (the `data` of the role) into the shape used by the form and the show view.
export function normalizeDatabaseRole(
  responseData: Record<string, any> | undefined, // eslint-disable-line @typescript-eslint/no-explicit-any
  name: string,
  backend: string,
  type: 'static' | 'dynamic'
): DatabaseRoleData {
  const { db_name, ...rest } = responseData ?? {};
  return {
    ...rest,
    name,
    backend,
    type,
    database: db_name ? [db_name] : [],
    // copy to singular for MongoDB
    creation_statement: rest['creation_statements']?.[0] ?? '',
    revocation_statement: rest['revocation_statements']?.[0] ?? '',
    // secret-edit-layout reads idPrefix from the model to pick the role tab, plain show data has no class to define it
    idPrefix: 'role/',
  };
}

const STATIC_STATEMENTS = ['rotation_statements'];
const DYNAMIC_STATEMENTS = [
  'creation_statements',
  'revocation_statements',
  'rollback_statements',
  'renew_statements',
];
const isEmpty = (value: unknown) =>
  value === undefined || value === null || value === '' || (Array.isArray(value) && !value.length);

export default class DatabaseRoleForm extends Form<DatabaseRoleData> {
  idPrefix = 'role/';

  // values the role had when it was loaded, used to send only the attributes that changed on update
  declare initial: Partial<DatabaseRoleData>;

  validations: Validations = {
    name: [{ type: 'presence', message: 'Role name is required.' }],
    database: [{ type: 'presence', message: 'Database is required.' }],
    type: [{ type: 'presence', message: 'Type is required.' }],
    username: [
      {
        validator(data) {
          const { type, username } = data as DatabaseRoleData;
          if (!type || type === 'dynamic') return true;
          return !!username;
        },
        message: 'Username is required.',
      },
    ],
  };

  constructor(data: Partial<DatabaseRoleData> = {}, options: FormOptions = {}) {
    super(data, options);
    this.initial = { ...data };
  }

  // returns every possible field so the proxy can discover each data key regardless of role type
  get formFields() {
    return Object.values(ROLE_FIELD_DEFINITIONS);
  }

  /*
   * The request body for creating or updating a role. Only attributes that belong to the role type are sent.
   * On update only the changed attributes are sent, apart from the username which is required for static roles.
   * The password is only sent when one has been entered. Validation runs against the form data, not the request body.
   */
  toJSON() {
    const validation = super.toJSON();
    const { data, initial, isNew } = this;
    const isStatic = data.type === 'static';
    const hasChanged = (key: string) => isNew || JSON.stringify(data[key]) !== JSON.stringify(initial[key]);
    const isCleared = (key: string) => !isNew && isEmpty(data[key]) && !isEmpty(initial[key]);

    const payload: Record<string, unknown> = {};
    const attrs = [...(getRoleFields(data.type) ?? [])];
    for (const key of attrs) {
      if (key === 'username' || isEmpty(data[key]) || !hasChanged(key)) continue;
      payload[key] = data[key];
    }
    if (isStatic) {
      // required to update a static role
      payload['username'] = data['username'];
    }
    if (isNew && data.database?.[0]) {
      payload['db_name'] = data.database[0];
    }

    const statements = isStatic ? STATIC_STATEMENTS : DYNAMIC_STATEMENTS;
    for (const key of statements) {
      // the input for MongoDB is a json string rather than an array so it is transposed into a list
      const singular = key.slice(0, -1);
      if (data[singular] && hasChanged(singular)) {
        payload[key] = [data[singular]];
      } else if (!isEmpty(data[key]) && hasChanged(key)) {
        payload[key] = data[key];
      } else if (isCleared(singular) || isCleared(key)) {
        // Vault keeps statements that are omitted from an update, so a cleared statement is sent as an empty list
        payload[key] = [];
      }
    }
    return { ...validation, data: payload as unknown as DatabaseRoleData };
  }
}
