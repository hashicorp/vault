/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';

@withFormFields()
export default class PkiConfigAcmeModel extends Model {
  // This model uses the backend value as the model ID
  get useOpenAPI() {
    return true;
  }

  getHelpUrl(backendPath) {
    return `/v1/${backendPath}/config/acme?help=1`;
  }

  // attrs order in the form is determined by order here

  @attr('boolean', {
    label: 'ACME enabled',
    subText: 'When ACME is disabled, all requests to ACME directory URLs will return 404.',
  })
  enabled;

  @attr('string', {
    subText:
      "Specifies the behavior of the default ACME directory. Can be 'forbid', 'sign-verbatim' or a role given by 'role:<role_name>'. If a role is used, it must be present in 'allowed_roles'.",
  })
  defaultDirectoryPolicy;

  @attr('array', {
    editType: 'stringArray',
    subText:
      "The default value '*' allows every role within the mount to be used. If the default_directory_policy specifies a role, it must be allowed under this configuration.",
  })
  allowedRoles;

  @attr('boolean', {
    label: 'Allow role ExtKeyUsage',
    subText:
      "When enabled, respect the role's ExtKeyUsage flags. Otherwise, ACME certificates are forced to ServerAuth.",
  })
  allowRoleExtKeyUsage;

  @attr('array', {
    editType: 'stringArray',
    subText:
      "Specifies a list of issuers allowed to issue certificates via explicit ACME paths. If an allowed role specifies an issuer outside this list, it will be allowed. The default value '*' allows every issuer within the mount.",
  })
  allowedIssuers;

  @attr('string', {
    label: 'EAB policy',
    possibleValues: ['not-required', 'new-account-required', 'always-required'],
  })
  eabPolicy;

  @attr('string', {
    label: 'DNS resolver',
    subText:
      'An optional overriding DNS resolver to use for challenge verification lookups. When not specified, the default system resolver will be used. This allows domains on peered networks with an accessible DNS resolver to be validated.',
  })
  dnsResolver;

  @attr({
    label: 'Max TTL',
    editType: 'ttl',
    hideToggle: true,
    helperTextEnabled:
      'Specify the maximum TTL for ACME certificates. Role TTL values will be limited to this value.',
  })
  maxTtl;

  @lazyCapabilities(apiPath`${'id'}/config/acme`, 'id')
  acmePath;

  get canSet() {
    return this.acmePath.get('canUpdate') !== false;
  }
}
