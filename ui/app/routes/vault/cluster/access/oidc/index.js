/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class OidcConfigureRoute extends Route {
  @service store;
  @service router;

  beforeModel() {
    return this.store
      .query('oidc/client', {})
      .then(() => {
        // transition to client list view if clients have been created
        this.router.transitionTo('vault.cluster.access.oidc.clients');
      })
      .catch(() => {
        // adapter throws error for 404 - swallow and remain on index route to show call to action
      });
  }
}
