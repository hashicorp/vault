/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { IdentityApiOidcListProvidersListEnum } from '@hashicorp/vault-client-typescript';

export default class OidcClientProvidersRoute extends Route {
  @service api;
  @service capabilities;

  async model() {
    try {
      const { client } = this.modelFor('vault.cluster.access.oidc.clients.client');
      const response = await this.api.identity.oidcListProviders(
        IdentityApiOidcListProvidersListEnum.TRUE,
        client.client_id // use allowed_client_id query param to filter providers for this client
      );
      const paths = response.keys.map((name) => this.capabilities.pathFor('oidcProvider', { name }));
      const capabilities = paths ? await this.capabilities.fetch(paths) : {};

      return {
        providers: this.api.keyInfoToArray(response, 'name'),
        capabilities,
      };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return {
          providers: [],
          capabilities: {},
        };
      } else {
        throw error;
      }
    }
  }
}
