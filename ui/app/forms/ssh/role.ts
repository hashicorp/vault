/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';
import FormFieldGroup from 'vault/utils/forms/field-group';

const OTP_FIELD_ORDER = [
  'name',
  'key_type',
  'default_user',
  'admin_user',
  'port',
  'allowed_users',
  'cidr_list',
  'exclude_cidr_list',
];

const CA_FIELD_ORDER = [
  'name',
  'key_type',
  'allow_user_certificates',
  'allow_host_certificates',
  'default_user',
  'allowed_users',
  'allowed_users_template',
  'allowed_domains',
  'allowed_domains_template',
  'ttl',
  'max_ttl',
  'allowed_critical_options',
  'default_critical_options',
  'allowed_extensions',
  'default_extensions',
  'allow_bare_domains',
  'allow_subdomains',
  'allow_empty_principals',
  'allow_user_key_ids',
  'key_id_format',
  'not_before_duration',
  'algorithm_signer',
];

// All field names across both key types, deduplicated, for proxy discovery
const ALL_FIELD_NAMES = [...new Set([...OTP_FIELD_ORDER, ...CA_FIELD_ORDER])];

const FIELDS: Record<string, FormField> = {
  name: new FormField('name', 'string', {
    label: 'Role name',
    fieldValue: 'name',
  }),
  key_type: new FormField('key_type', 'string', {
    possibleValues: ['ca', 'otp'],
  }),
  admin_user: new FormField('admin_user', 'string', {
    helpText: 'Username of the admin user at the remote host',
  }),
  default_user: new FormField('default_user', 'string', {
    helpText: "Username to use when one isn't specified",
  }),
  allowed_users: new FormField('allowed_users', 'string', {
    helpText:
      'Create a list of users who are allowed to use this key (e.g. `admin, dev`, or use `*` to allow all.)',
  }),
  allowed_users_template: new FormField('allowed_users_template', 'boolean', {
    helpText:
      'Specifies that Allowed Users can be templated e.g. {{identity.entity.aliases.mount_accessor_xyz.name}}',
  }),
  allowed_domains: new FormField('allowed_domains', 'string', {
    helpText:
      'List of domains for which a client can request a certificate (e.g. `example.com`, or `*` to allow all)',
  }),
  allowed_domains_template: new FormField('allowed_domains_template', 'boolean', {
    helpText:
      'Specifies that Allowed Domains can be set using identity template policies. Non-templated domains are also permitted.',
  }),
  cidr_list: new FormField('cidr_list', 'string', {
    helpText: 'List of CIDR blocks for which this role is applicable',
  }),
  exclude_cidr_list: new FormField('exclude_cidr_list', 'string', {
    helpText: 'List of CIDR blocks that are not accepted by this role',
  }),
  port: new FormField('port', 'number', {
    helpText: 'Port number for the SSH connection (default is `22`)',
  }),
  allowed_critical_options: new FormField('allowed_critical_options', 'string', {
    helpText: 'List of critical options that certificates have when signed',
  }),
  default_critical_options: new FormField('default_critical_options', 'object', {
    helpText: 'Map of critical options certificates should have if none are provided when signing',
  }),
  allowed_extensions: new FormField('allowed_extensions', 'string', {
    helpText: 'List of extensions that certificates can have when signed',
  }),
  default_extensions: new FormField('default_extensions', 'object', {
    helpText: 'Map of extensions certificates should have if none are provided when signing',
  }),
  allow_user_certificates: new FormField('allow_user_certificates', 'boolean', {
    helpText: 'Specifies if certificates are allowed to be signed for us as a user',
  }),
  allow_host_certificates: new FormField('allow_host_certificates', 'boolean', {
    helpText: 'Specifies if certificates are allowed to be signed for us as a host',
  }),
  allow_bare_domains: new FormField('allow_bare_domains', 'boolean', {
    helpText:
      'Specifies if host certificates that are requested are allowed to use the base domains listed in Allowed Domains',
  }),
  allow_subdomains: new FormField('allow_subdomains', 'boolean', {
    helpText:
      'Specifies if host certificates that are requested are allowed to be subdomains of those listed in Allowed Domains',
  }),
  allow_empty_principals: new FormField('allow_empty_principals', 'boolean', {
    helpText:
      'Allow signing certificates with no valid principals (e.g. any valid principal). For backwards compatibility only. The default of false is highly recommended.',
  }),
  allow_user_key_ids: new FormField('allow_user_key_ids', 'boolean', {
    helpText: 'Specifies if users can override the key ID for a signed certificate with the "key_id" field',
  }),
  key_id_format: new FormField('key_id_format', 'string', {
    helpText: 'When supplied, this value specifies a custom format for the key id of a signed certificate',
  }),
  not_before_duration: new FormField('not_before_duration', 'string', {
    helpText: 'Specifies the duration by which to backdate the ValidAfter property',
    editType: 'ttl',
  }),
  ttl: new FormField('ttl', 'string', {
    editType: 'ttl',
  }),
  max_ttl: new FormField('max_ttl', 'string', {
    editType: 'ttl',
  }),
  algorithm_signer: new FormField('algorithm_signer', 'string', {
    helpText: 'When supplied, this value specifies a signing algorithm for the key',
    possibleValues: ['default', 'ssh-rsa', 'rsa-sha2-256', 'rsa-sha2-512'],
  }),
};

type SshRoleData = {
  name: string;
  id: string;
  backend: string;
  key_type: string;
  [key: string]: unknown;
};

export default class SshRoleForm extends Form<SshRoleData> {
  /*
   * formFieldGroups always returns all possible fields so the proxy can
   * discover every data key regardless of which key_type is active.
   * This is called directly on the raw target inside the proxy handler,
   * so it must not read any proxied data properties (e.g. this.key_type).
   */
  get formFieldGroups() {
    return [
      new FormFieldGroup(
        'default',
        ALL_FIELD_NAMES.map((name) => FIELDS[name]).filter((f): f is FormField => !!f)
      ),
    ];
  }

  /*
   * fieldGroups is what FormFieldGroupsLoop reads. It is accessed through
   * the proxy, so this.key_type correctly resolves to data.key_type.
   */
  get fieldGroups() {
    const isOtp = this.data.key_type === 'otp';
    const defaultFieldsNum = isOtp ? 3 : 4;
    const fieldOrder = isOtp ? [...OTP_FIELD_ORDER] : [...CA_FIELD_ORDER];
    const defaultFieldNames = fieldOrder.splice(0, defaultFieldsNum);
    return [
      new FormFieldGroup(
        'default',
        defaultFieldNames.map((name) => FIELDS[name]).filter((f): f is FormField => !!f)
      ),
      new FormFieldGroup(
        'Options',
        fieldOrder.map((name) => FIELDS[name]).filter((f): f is FormField => !!f)
      ),
    ];
  }

  // Flat ordered list of fields for the current key_type, used by the show view.
  get displayFields() {
    const fieldOrder = this.data.key_type === 'otp' ? OTP_FIELD_ORDER : CA_FIELD_ORDER;
    return fieldOrder.map((name) => FIELDS[name]).filter((f): f is FormField => !!f);
  }
}
