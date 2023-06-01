/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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

  @attr('array', {
    editType: 'stringArray',
    subText:
      'Specifies a list of roles to allow to issue certificates via explicit ACME paths. If no default_role is specified, sign-verbatim-like issuance on the default ACME directory will still occur. The default value * allows every role within the mount. ',
  })
  allowedRoles;

  @attr('array', {
    editType: 'stringArray',
    subText:
      'Specifies a list issuers allowed to issue certificates via explicit ACME paths. If an allowed role specifies an issuer outside this list, it will be allowed. The default value * allows every issuer within the mount. ',
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

  @lazyCapabilities(apiPath`${'id'}/config/acme`, 'id')
  acmePath;

  get canSet() {
    return this.acmePath.get('canCreate') !== false;
  }
}
