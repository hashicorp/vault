/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  binddn: [{ type: 'presence', message: 'Administrator distinguished name is required.' }],
  bindpass: [{ type: 'presence', message: 'Administrator password is required.' }],
};
const formGroups = [
  { default: ['binddn', 'bindpass', 'url', 'password_policy'] },
  { 'TLS options': ['starttls', 'insecure_tls', 'certificate', 'client_tls_cert', 'client_tls_key'] },
  { 'More options': ['userdn', 'userattr', 'upndomain', 'connection_timeout', 'request_timeout'] },
];

@withModelValidations(validations)
@withFormFields(null, formGroups)
export default class LdapConfigModel extends Model {
  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord

  @attr('string', {
    label: 'Administrator Distinguished Name',
    subText:
      'Distinguished name of the administrator to bind (Bind DN) when performing user and group search. Example: cn=vault,ou=Users,dc=example,dc=com.',
  })
  binddn;

  @attr('string', {
    label: 'Administrator Password',
    subText: 'Password to use along with Bind DN when performing user search.',
  })
  bindpass;

  @attr('string', {
    label: 'URL',
    subText: 'The directory server to connect to.',
  })
  url;

  @attr('string', {
    editType: 'optionalText',
    label: 'Use custom password policy',
    subText: 'Specify the name of an existing password policy.',
    defaultSubText: 'Unless a custom policy is specified, Vault will use a default.',
    defaultShown: 'Default',
    docLink: '/vault/docs/concepts/password-policies',
  })
  password_policy;

  @attr('string') schema;

  @attr('boolean', {
    label: 'Start TLS',
    subText: 'If checked, or address contains “ldaps://”, creates an encrypted connection with LDAP.',
  })
  starttls;

  @attr('boolean', {
    label: 'Insecure TLS',
    subText: 'If checked, skips LDAP server SSL certificate verification - insecure, use with caution!',
  })
  insecure_tls;

  @attr('string', {
    editType: 'file',
    label: 'CA Certificate',
    helpText: 'CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded.',
  })
  certificate;

  @attr('string', {
    editType: 'file',
    label: 'Client TLS Certificate',
    helpText: 'Client certificate to provide to the LDAP server, must be x509 PEM encoded.',
  })
  client_tls_cert;

  @attr('string', {
    editType: 'file',
    label: 'Client TLS Key',
    helpText: 'Client key to provide to the LDAP server, must be x509 PEM encoded.',
  })
  client_tls_key;

  @attr('string', {
    label: 'Userdn',
    helpText: 'The base DN under which to perform user search in library management and static roles.',
  })
  userdn;

  @attr('string', {
    label: 'Userattr',
    subText: 'The attribute field name used to perform user search in library management and static roles.',
  })
  userattr;

  @attr('string', {
    label: 'Upndomain',
    subText: 'The domain (userPrincipalDomain) used to construct a UPN string for authentication.',
  })
  upndomain;

  @attr('number', {
    editType: 'optionalText',
    label: 'Connection Timeout',
    subText: 'Specify the connection timeout length in seconds.',
    defaultSubText: 'Vault will use the default of 30 seconds.',
    defaultShown: 'Default 30 seconds.',
  })
  connection_timeout;

  @attr('number', {
    editType: 'optionalText',
    label: 'Request Timeout',
    subText: 'Specify the connection timeout length in seconds.',
    defaultSubText: 'Vault will use the default of 90 seconds.',
    defaultShown: 'Default 90 seconds.',
  })
  request_timeout;

  async rotateRoot() {
    const adapter = this.store.adapterFor('ldap/config');
    return adapter.rotateRoot(this.backend);
  }
}
