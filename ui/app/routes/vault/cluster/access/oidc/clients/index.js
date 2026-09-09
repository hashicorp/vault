/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { IdentityApiOidcListClientsListEnum } from '@hashicorp/vault-client-typescript';

export default class OidcClientsRoute extends Route {
  @service api;
  @service capabilities;
  @service router;

  async model() {
    try {
      const response = await this.api.identity.oidcListClients(IdentityApiOidcListClientsListEnum.TRUE);
      const paths = response.keys.map((name) => this.capabilities.pathFor('oidcClient', { name }));
      const capabilities = paths ? await this.capabilities.fetch(paths) : {};

      return {
        clients: this.api.keyInfoToArray(response, 'name'),
        capabilities,
      };
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        return {
          clients: [],
          capabilities: {},
        };
      } else {
        throw error;
      }
    }
  }

  afterModel(model) {
    if (model.clients.length === 0) {
      this.router.transitionTo('vault.cluster.access.oidc');
    }
  }
}
