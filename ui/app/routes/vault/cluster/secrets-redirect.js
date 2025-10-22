/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class SecretsRedirectRoute extends Route {
  @service router;

  beforeModel() {
    // Redirect to secrets page under /secrets-engines
    // if the user navigates to the legacy path /secrets
    this.router.replaceWith('vault.cluster.secrets');
  }
}
