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

  beforeModel() {
    return this.flags.fetchActivatedFeatures();
  }

  model() {
    return {
      // TODO this is a half way solution until persona service is implemented. Additionally, we should move away from calling the response of this endpoint features, and instead use flags which the noun used in the endpoint.
      activatedFeatures: this.flags.activatedFlags,
    };
  }

  afterModel(model: { activatedFeatures: Array<string> }) {
    if (!model.activatedFeatures) {
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
    }
  }
}
