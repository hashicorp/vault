/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class OidcProviderClientsRoute extends Route {
  @service store;

  async model() {
    const { allowedClientIds } = this.modelFor('vault.cluster.access.oidc.providers.provider');
    return await this.store.query('oidc/client', { paramKey: 'client_id', filterFor: allowedClientIds });
  }
}
