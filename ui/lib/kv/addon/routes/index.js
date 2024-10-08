/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class KvRoute extends Route {
  @service('app-router') router;

  redirect() {
    this.router.transitionTo('vault.cluster.secrets.backend.kv.list');
  }
}
