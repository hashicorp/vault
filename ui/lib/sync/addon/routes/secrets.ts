/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';
import type FlagService from 'vault/services/flags';

export default class SyncSecretsRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly flags: FlagService;

  model() {
    return {
      activatedFeatures: this.flags.activatedFlags,
    };
  }

  afterModel(model: { activatedFeatures: Array<string> }) {
    if (!model.activatedFeatures) {
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
    }
  }
}
