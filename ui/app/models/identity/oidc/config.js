/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class IdentityOidcConfig extends Model {
  @attr('string', {
    label: 'Issuer',
    subText:
      "The Issuer URL to be used in configuring Vault as an identity provider in AWS. If not set, Vault's default issuer will be used",
    docLink: '/vault/api-docs/secret/identity/tokens#configure-the-identity-tokens-backend',
    placeholder: 'https://vault.prod/v1/identity/oidc',
  })
  issuer;

  get attrs() {
    const keys = ['issuer'];
    return expandAttributeMeta(this, keys);
  }
  get formFields() {
    const keys = ['issuer'];
    return expandAttributeMeta(this, keys);
  }

  // CAPABILITIES
  @lazyCapabilities(apiPath`identity/oidc/config`) issuerPath;
  get canRead() {
    return this.issuerPath.get('canRead') !== false;
  }
}
