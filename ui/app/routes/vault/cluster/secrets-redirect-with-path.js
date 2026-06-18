/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class SecretsRedirectWithPathRoute extends Route {
  @service router;

  beforeModel(transition) {
    // Redirect to secrets page under /secrets-engines
    // if the user navigates to the legacy path /secrets
    // Preserve the full path after /secrets (e.g., /secrets/kv/kv/list -> /secrets-engines/kv/kv/list)
    const params = transition.to.params;
    const path = params?.path;

    if (path) {
      // Construct the new URL with full path including /vault/secrets-engines/*path
      const newUrl = `/vault/secrets-engines/${path}`;
      this.router.replaceWith(newUrl);
    } else {
      // If no path, just redirect to the base secrets page
      this.router.replaceWith('vault.cluster.secrets');
    }
  }
}
