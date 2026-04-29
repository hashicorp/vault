/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { IdentityApiOidcListClientsListEnum } from '@hashicorp/vault-client-typescript';

export default class OidcKeyClientsRoute extends Route {
  @service api;
  @service capabilities;

  async model() {
    const { allowedClientIds } = this.modelFor('vault.cluster.access.oidc.keys.key');
    const response = await this.api.identity.oidcListClients(IdentityApiOidcListClientsListEnum.TRUE);
    const clients = this.api.keyInfoToArray(response, 'name');
    // filter clients based on allowed_client_ids of provider
    const filteredClients = allowedClientIds.includes('*')
      ? clients
      : clients.filter((client) => allowedClientIds.includes(client.client_id));
    // fetch capabilities for filtered clients
    const paths = filteredClients.map(({ name }) => this.capabilities.pathFor('oidcClient', { name }));
    const capabilities = paths ? await this.capabilities.fetch(paths) : {};

    return { clients: filteredClients, capabilities };
  }
}
