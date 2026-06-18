/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { IdentityApiOidcListClientsListEnum } from '@hashicorp/vault-client-typescript';
import OidcKeyForm from 'vault/forms/oidc/key';

export default class OidcKeyEditRoute extends Route {
  @service api;

  async model() {
    const { key } = this.modelFor('vault.cluster.access.oidc.keys.key');
    // fetch clients to populate dropdown in form
    const response = await this.api.identity.oidcListClients(IdentityApiOidcListClientsListEnum.TRUE);
    const clients = this.api.keyInfoToArray(response, 'name');
    // filter clients that are associated with this key
    const filteredClients = clients.filter((client) => client.key === key.name);

    return {
      clients: filteredClients,
      form: new OidcKeyForm(key),
    };
  }
}
