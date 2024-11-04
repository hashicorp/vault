/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import LdapRoleModel from '../role';
import { attr } from '@ember-data/model';
import { withModelValidations } from 'vault/decorators/model-validations';
import { withFormFields } from 'vault/decorators/model-form-fields';

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
  creation_ldif: [
    {
      validator: (model) => (!model.creation_ldif ? false : true),
      message: 'Creation LDIF is required.',
    },
  ],
  deletion_ldif: [
    {
      validator: (model) => (!model.creation_ldif ? false : true),
      message: 'Deletion LDIF is required.',
    },
  ],
};

// determines form input rendering order
@withFormFields([
  'name',
  'default_ttl',
  'max_ttl',
  'username_template',
  'creation_ldif',
  'deletion_ldif',
  'rollback_ldif',
])
@withModelValidations(validations)
export default class LdapRoleDynamicModel extends LdapRoleModel {
  type = 'dynamic';
  roleUri = 'role';
  credsUri = 'creds';

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
}
