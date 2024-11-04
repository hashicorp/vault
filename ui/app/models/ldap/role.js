/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const validations = {
  name: [{ type: 'presence', message: 'Name is required' }],
};

export const staticRoleFields = ['username', 'dn', 'rotation_period'];
export const dynamicRoleFields = [
  'default_ttl',
  'max_ttl',
  'username_template',
  'creation_ldif',
  'deletion_ldif',
  'rollback_ldif',
];

@withModelValidations(validations)
@withFormFields()
export default class LdapRoleModel extends Model {
  @attr backend; // dynamic path of secret -- set on response from value passed to queryRecord

  @attr('string', {
    label: 'Role name',
    subText: 'The name of the role that will be used in Vault.',
    editDisabled: true,
  })
  name;

  get isStatic() {
    return this.type === 'static';
  }
  get isDynamic() {
    return this.type === 'dynamic';
  }
  // this is used to build the form fields as well as serialize the correct payload based on type
  // if a new attr is added be sure to add it to the appropriate array
  get fieldsForType() {
    return this.isStatic
      ? ['username', 'dn', 'rotation_period']
      : ['default_ttl', 'max_ttl', 'username_template', 'creation_ldif', 'deletion_ldif', 'rollback_ldif'];
  }

  get displayFields() {
    // insert type after role name
    const [name, ...rest] = this.formFields;
    const typeField = { name: 'type', options: { label: 'Role type' } };
    return [name, typeField, ...rest];
  }

  @lazyCapabilities(apiPath`${'backend'}/${'roleUri'}/${'name'}`, 'backend', 'roleUri', 'name') rolePath;
  @lazyCapabilities(apiPath`${'backend'}/${'credsUri'}/${'name'}`, 'backend', 'credsUri', 'name') credsPath;

  get canCreate() {
    return this.rolePath.get('canCreate') !== false;
  }
  get canDelete() {
    return this.rolePath.get('canDelete') !== false;
  }
  get canEdit() {
    return this.rolePath.get('canUpdate') !== false;
  }
  get canRead() {
    return this.rolePath.get('canRead') !== false;
  }
  get canList() {
    return this.rolePath.get('canList') !== false;
  }
  get canReadCreds() {
    return this.credsPath.get('canRead') !== false;
  }

  fetchCredentials() {
    return this.store.adapterFor('ldap/role').fetchCredentials(this.backend, this.type, this.name);
  }
  rotateStaticPassword() {
    return this.store.adapterFor('ldap/role').rotateStaticPassword(this.backend, this.name);
  }
}
