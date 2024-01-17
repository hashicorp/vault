/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class KvRoute extends Route {
  @service router;

  redirect() {
    this.router.transitionTo('vault.cluster.secrets.backend.kv.list');
  }
}
