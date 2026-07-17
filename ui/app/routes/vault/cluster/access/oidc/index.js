/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { IdentityApiOidcListClientsListEnum } from '@hashicorp/vault-client-typescript';

export default class OidcConfigureRoute extends Route {
  @service api;
  @service router;

  async beforeModel() {
    try {
      const { keys } = await this.api.identity.oidcListClients(IdentityApiOidcListClientsListEnum.TRUE);
      if (keys?.length) {
        // transition to client list view if clients have been created
        this.router.transitionTo('vault.cluster.access.oidc.clients');
      }
    } catch (e) {
      // swallow error and remain on index route to show call to action
    }
  }
}
