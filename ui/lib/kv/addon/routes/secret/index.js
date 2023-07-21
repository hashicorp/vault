/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class SecretIndex extends Route {
  @service router;

  redirect() {
    this.router.transitionTo('vault.cluster.secrets.backend.kv.secret.details');
  }
}
