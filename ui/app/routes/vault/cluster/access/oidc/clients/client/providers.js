/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class OidcClientProvidersRoute extends Route {
  @service store;

  model() {
    const { client } = this.modelFor('vault.cluster.access.oidc.clients.client');
    return this.store
      .query('oidc/provider', {
        allowed_client_id: client.client_id,
      })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  }
}
