/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Form from 'vault/forms/form';
import FormField from 'vault/utils/forms/field';

import type { PkiConfigureAcmeRequest } from '@hashicorp/vault-client-typescript';

export default class PkiConfigAcmeForm extends Form<PkiConfigureAcmeRequest> {
  formFields = [
    new FormField('enabled', 'boolean', {
      label: 'ACME enabled',
      subText: 'When ACME is disabled, all requests to ACME directory URLs will return 404.',
    }),
    new FormField('default_directory_policy', 'string', {
      subText:
        "Specifies the behavior of the default ACME directory. Can be 'forbid', 'sign-verbatim' or a role given by 'role:<role_name>'. If a role is used, it must be present in 'allowed_roles'.",
    }),
    new FormField('allowed_roles', 'string', {
      editType: 'stringArray',
      subText:
        "The default value '*' allows every role within the mount to be used. If the default_directory_policy specifies a role, it must be allowed under this configuration.",
    }),
    new FormField('allowed_role_ext_key_usage', 'boolean', {
      label: 'Allow role ExtKeyUsage',
      subText:
        "When enabled, respect the role's ExtKeyUsage flags. Otherwise, ACME certificates are forced to ServerAuth.",
    }),
    new FormField('allowed_issuers', 'string', {
      editType: 'stringArray',
      subText:
        "Specifies a list of issuers allowed to issue certificates via explicit ACME paths. If an allowed role specifies an issuer outside this list, it will be allowed. The default value '*' allows every issuer within the mount.",
    }),
    new FormField('eab_policy', 'string', {
      label: 'EAB policy',
      possibleValues: ['not-required', 'new-account-required', 'always-required'],
    }),
    new FormField('dns_resolver', 'string', {
      label: 'DNS resolver',
      subText:
        'An optional overriding DNS resolver to use for challenge verification lookups. When not specified, the default system resolver will be used. This allows domains on peered networks with an accessible DNS resolver to be validated.',
    }),
    new FormField('max_ttl', 'string', {
      label: 'Max TTL',
      editType: 'ttl',
      hideToggle: true,
      helperTextEnabled:
        'Specify the maximum TTL for ACME certificates. Role TTL values will be limited to this value.',
    }),
  ];
}
