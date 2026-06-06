/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import OidcProviderForm from 'vault/forms/oidc/provider';
import {
  IdentityApiOidcListScopesListEnum,
  IdentityApiOidcListClientsListEnum,
} from '@hashicorp/vault-client-typescript';

export default class OidcProvidersCreateRoute extends Route {
  @service api;

  async model() {
    // fetch scopes and clients to populate dropdowns in form
    const [scopesResult, clientsResult] = await Promise.allSettled([
      this.api.identity.oidcListScopes(IdentityApiOidcListScopesListEnum.TRUE),
      this.api.identity.oidcListClients(IdentityApiOidcListClientsListEnum.TRUE),
    ]);

    // SearchSelect requires options to be objects.
    const scopes = scopesResult.value?.keys?.map((key) => ({ id: key })) ?? [];
    const clients = this.api.keyInfoToArray(clientsResult.value, 'name');

    return {
      form: new OidcProviderForm({}, { isNew: true }),
      scopes,
      clients,
    };
  }
}
