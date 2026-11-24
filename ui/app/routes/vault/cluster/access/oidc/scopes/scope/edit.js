/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import OidcScopeForm from 'vault/forms/oidc/scope';

export default class OidcScopeEditRoute extends Route {
  model() {
    const { scope } = this.modelFor('vault.cluster.access.oidc.scopes.scope');
    return new OidcScopeForm(scope);
  }
}
