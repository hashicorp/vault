/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import OidcClientForm from 'vault/forms/oidc/client';
import { IdentityApiOidcListKeysListEnum } from '@hashicorp/vault-client-typescript';

export default class OidcClientsCreateRoute extends Route {
  @service api;

  async model() {
    // fetch keys to populate dropdown in form
    let keys = [];
    try {
      const response = await this.api.identity.oidcListKeys(IdentityApiOidcListKeysListEnum.TRUE);
      // SearchSelect requires options to be objects
      keys = response.keys?.map((key) => ({ id: key }));
    } catch (error) {
      // swallow error and return empty array for keys
    }
    const defaults = {
      key: 'default',
      id_token_ttl: '24h',
      access_token_ttl: '24h',
      client_type: 'confidential',
    };
    return {
      form: new OidcClientForm(defaults, { isNew: true }),
      keys,
    };
  }
}
