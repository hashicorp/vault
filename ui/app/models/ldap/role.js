/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

const creationLdifExample = `# The example below is treated as a comment and will not be submitted
# dn: cn={{.Username}},ou=users,dc=learn,dc=example
# objectClass: person
# objectClass: top
`;
const deletionLdifExample = `# The example below is treated as a comment and will not be submitted
# dn: cn={{.Username}},ou=users,dc=learn,dc=example
# changetype: delete
`;
const rollbackLdifExample = `# The example below is treated as a comment and will not be submitted
# dn: cn={{.Username}},ou=users,dc=learn,dc=example
# changetype: delete
`;

const validations = {
  name: [{ type: 'presence', message: 'Name is required' }],
  username: [
    {
      validator: (model) => (model.isStatic && !model.username ? false : true),
      message: 'Username is required.',
    },
  ],
  rotation_period: [
    {
      validator: (model) => (model.isStatic && !model.rotation_period ? false : true),
      message: 'Rotation Period is required.',
    },
  ],
  creation_ldif: [
    {
      validator: (model) => (model.isDynamic && !model.creation_ldif ? false : true),
      message: 'Creation LDIF is required.',
    },
  ],
  deletion_ldif: [
    {
      validator: (model) => (model.isDynamic && !model.creation_ldif ? false : true),
      message: 'Deletion LDIF is required.',
    },
  ],
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
  @attr('string') backend; // mount path of ldap engine -- set on response from value passed to queryRecord
  @attr('string') path_to_role; // ancestral path to the role added in the adapter (only exists for nested roles)

  @attr('string', {
    defaultValue: 'static',
  })
  type; // this must be set to either static or dynamic in order for the adapter to build the correct url and for the correct form fields to display

  @attr('string', {
    label: 'Role name',
    subText: 'The name of the role that will be used in Vault.',
    editDisabled: true,
  })
  name;

  // static role properties
  @attr('string', {
    label: 'Distinguished name',
    subText: 'Distinguished name (DN) of entry Vault should manage.',
  })
  dn;

  @attr('string', {
    label: 'Username',
    subText:
      "The name of the user to be used when logging in. This is useful when DN isn't used for login purposes.",
  })
  username;

  @attr({
    editType: 'ttl',
    label: 'Rotation period',
    helperTextEnabled:
      'Specifies the amount of time Vault should wait before rotating the password. The minimum is 5 seconds.',
    hideToggle: true,
  })
  rotation_period;

  // dynamic role properties
  @attr({
    editType: 'ttl',
    label: 'Generated credential’s time-to-live (TTL)',
    detailsLabel: 'TTL',
    helperTextDisabled: 'Vault will use the default of 1 hour.',
    defaultValue: '1h',
    defaultShown: 'Engine default',
  })
  default_ttl;

  @attr({
    editType: 'ttl',
    label: 'Generated credential’s maximum time-to-live (Max TTL)',
    detailsLabel: 'Max TTL',
    helperTextDisabled: 'Vault will use the engine default of 24 hours.',
    defaultValue: '24h',
    defaultShown: 'Engine default',
  })
  max_ttl;

  @attr('string', {
    editType: 'optionalText',
    label: 'Username template',
    subText: 'Enter the custom username template to use.',
    defaultSubText:
      'Template describing how dynamic usernames are generated. Vault will use the default for this plugin.',
    docLink: '/vault/docs/concepts/username-templating',
    defaultShown: 'Default',
  })
  username_template;

  @attr('string', {
    editType: 'json',
    label: 'Creation LDIF',
    helpText: 'Specifies the LDIF statements executed to create a user. May optionally be base64 encoded.',
    example: creationLdifExample,
    mode: 'ruby',
    sectionHeading: 'LDIF Statements', // render section heading before form field
  })
  creation_ldif;

  @attr('string', {
    editType: 'json',
    label: 'Deletion LDIF',
    helpText:
      'Specifies the LDIF statements executed to delete a user once its TTL has expired. May optionally be base64 encoded.',
    example: deletionLdifExample,
    mode: 'ruby',
  })
  deletion_ldif;

  @attr('string', {
    editType: 'json',
    label: 'Rollback LDIF',
    helpText:
      'Specifies the LDIF statement to attempt to rollback any changes if the creation results in an error. May optionally be base64 encoded.',
    example: rollbackLdifExample,
    mode: 'ruby',
  })
  rollback_ldif;

  get completeRoleName() {
    // if there is a path_to_role then the name is hierarchical
    // and we must concat the ancestors with the leaf name to get the full role path
    return this.path_to_role ? this.path_to_role + this.name : this.name;
  }

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
  get formFields() {
    // filter all fields and return only those relevant to type
    return this.allFields.filter((field) => {
      // name is the only common field
      return field.name === 'name' || this.fieldsForType.includes(field.name);
    });
  }

  get displayFields() {
    // insert type after role name
    const [name, ...rest] = this.formFields;
    const typeField = { name: 'type', options: { label: 'Role type' } };
    return [name, typeField, ...rest];
  }

  get roleUri() {
    return this.isStatic ? 'static-role' : 'role';
  }
  get credsUri() {
    return this.isStatic ? 'static-cred' : 'creds';
  }
  @lazyCapabilities(apiPath`${'backend'}/${'roleUri'}/${'name'}`, 'backend', 'roleUri', 'name') rolePath;
  @lazyCapabilities(apiPath`${'backend'}/${'credsUri'}/${'name'}`, 'backend', 'credsUri', 'name') credsPath;
  @lazyCapabilities(apiPath`${'backend'}/rotate-role/${'name'}`, 'backend', 'name') staticRotateCredsPath;

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
  get canRotateStaticCreds() {
    return this.isStatic && this.staticRotateCredsPath.get('canCreate') !== false;
  }

  fetchCredentials() {
    return this.store
      .adapterFor('ldap/role')
      .fetchCredentials(this.backend, this.type, this.completeRoleName);
  }
  rotateStaticPassword() {
    return this.store.adapterFor('ldap/role').rotateStaticPassword(this.backend, this.completeRoleName);
  }
}
