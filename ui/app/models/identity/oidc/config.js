/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Model, { attr } from '@ember-data/model';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default class IdentityOidcConfig extends Model {
  @attr('string', {
    label: 'Issuer',
    subText:
      "The Issuer URL to be used in configuring Vault as an identity provider. If not set, Vault's default issuer will be used.",
    docLink: '/vault/api-docs/secret/identity/tokens#configure-the-identity-tokens-backend',
    placeholder: 'https://vault-test.com',
  })
  issuer;

  get attrs() {
    const keys = ['issuer'];
    return expandAttributeMeta(this, keys);
  }
}
