/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class SecretsBackendCreateRoute extends Route {
  @service router;

  beforeModel() {
    const { secret, initialKey } = this.paramsFor(this.routeName);
    const qp = initialKey || secret;
    return this.router.transitionTo('vault.cluster.secrets.backend.create-root', {
      queryParams: { initialKey: qp },
    });
  }
}
