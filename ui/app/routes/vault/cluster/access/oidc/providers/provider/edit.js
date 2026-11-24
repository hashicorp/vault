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
import parseURL from 'core/utils/parse-url';

export default class OidcProviderEditRoute extends Route {
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
    const { provider } = this.modelFor('vault.cluster.access.oidc.providers.provider');
    // parse issuer to only include scheme, host, and port for form field
    const formData = {
      ...provider,
      issuer: provider.issuer ? parseURL(provider.issuer).origin : '',
    };

    return {
      form: new OidcProviderForm(formData),
      scopes,
      clients,
    };
  }
}
