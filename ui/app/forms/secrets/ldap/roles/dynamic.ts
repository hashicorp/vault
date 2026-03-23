/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { Validations } from 'vault/app-types';
import type { LdapWriteDynamicRoleRequest } from '@hashicorp/vault-client-typescript';

type LdapDynamicRoleFormData = LdapWriteDynamicRoleRequest & {
  name: string;
};

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

export default class LdapDynamicRoleForm extends Form<LdapDynamicRoleFormData> {
  formFields = [
    new FormField('name', 'string', {
      label: 'Role name',
      subText: 'The name of the role that will be used in Vault.',
      editDisabled: true,
    }),
    new FormField('default_ttl', 'number', {
      editType: 'ttl',
      label: 'Generated credential’s time-to-live (TTL)',
      helperTextDisabled: 'Vault will use the default of 1 hour.',
      defaultValue: '1h',
      defaultShown: 'Engine default',
    }),

    new FormField('max_ttl', 'number', {
      editType: 'ttl',
      label: 'Generated credential’s maximum time-to-live (Max TTL)',
      helperTextDisabled: 'Vault will use the engine default of 24 hours.',
      defaultValue: '24h',
      defaultShown: 'Engine default',
    }),
    new FormField('username_template', 'string', {
      editType: 'optionalText',
      label: 'Username template',
      subText: 'Enter the custom username template to use.',
      defaultSubText:
        'Template describing how dynamic usernames are generated. Vault will use the default for this plugin.',
      docLink: '/vault/docs/concepts/username-templating',
      defaultShown: 'Default',
    }),
    new FormField('creation_ldif', 'string', {
      editType: 'json',
      label: 'Creation LDIF',
      helpText: 'Specifies the LDIF statements executed to create a user. May optionally be base64 encoded.',
      example: creationLdifExample,
      mode: 'ruby',
    }),
    new FormField('deletion_ldif', 'string', {
      editType: 'json',
      label: 'Deletion LDIF',
      helpText:
        'Specifies the LDIF statements executed to delete a user once its TTL has expired. May optionally be base64 encoded.',
      example: deletionLdifExample,
      mode: 'ruby',
    }),
    new FormField('rollback_ldif', 'string', {
      editType: 'json',
      label: 'Rollback LDIF',
      helpText:
        'Specifies the LDIF statement to attempt to rollback any changes if the creation results in an error. May optionally be base64 encoded.',
      example: rollbackLdifExample,
      mode: 'ruby',
    }),
  ];

  validations: Validations = {
    name: [{ type: 'presence', message: 'Name is required' }],
    creation_ldif: [{ type: 'presence', message: 'Creation LDIF is required.' }],
    deletion_ldif: [{ type: 'presence', message: 'Deletion LDIF is required.' }],
  };
}
