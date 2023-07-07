/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import { withModelValidations } from 'vault/decorators/model-validations';

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
  @attr('string') type; // this must be set to either static or dynamic in order for the adapter to build the correct url and for the correct form fields to display

  @attr('string') backend; // dynamic path of secret -- set on response from value passed to queryRecord

  @attr('string', {
    label: 'Role Name',
    subText: 'The name of the role that will be used in Vault.',
  })
  name;

  // static role properties
  @attr('string', {
    label: 'Distinguished Name',
    subText: 'Distinguished name (DN) of entry Vault should manage.',
  })
  dn;

  @attr('string', {
    label: 'Username',
    subText:
      "The name of the user to be used when logging in. This is useful when DN isn't used for login purposes.",
  })
  username;

  @attr('string', {
    editType: 'ttl',
    label: 'Rotation Period',
    subText:
      'Specifies the amount of time Vault should wait before rotating the password. The minimum is 5 seconds.',
  })
  rotation_period;

  // dynamic role properties
  @attr('number', {
    editType: 'ttl',
    label: 'Generated credentials’s Time-to-Live (TTL)',
    detailsLabel: 'TTL',
    helperTextDisabled: 'Vault will use the default of 1 hour',
    defaultValue: '1h',
    defaultShown: 'Engine default',
  })
  default_ttl;

  @attr('number', {
    editType: 'ttl',
    label: 'Generated credentials’s maximum Time-to-Live (Max TTL)',
    detailsLabel: 'Max TTL',
    helperTextDisabled: 'Vault will use the engine default of 1 hour',
    defaultValue: '24h',
    defaultShown: 'Engine default',
  })
  max_ttl;

  @attr('string', {
    editType: 'optionalText',
    label: 'Username Template',
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
    defaultValue: creationLdifExample,
  })
  creation_ldif;

  @attr('string', {
    editType: 'json',
    label: 'Deletion LDIF',
    helpText:
      'Specifies the LDIF statements executed to delete a user once its TTL has expired. May optionally be base64 encoded.',
    defaultValue: deletionLdifExample,
  })
  deletion_ldif;

  @attr('string', {
    editType: 'json',
    label: 'Rollback LDIF',
    helpText:
      'Specifies the LDIF statement to attempt to rollback any changes if the creation results in an error. May optionally be base64 encoded.',
    defaultValue: rollbackLdifExample,
  })
  rollback_ldif;

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
}
