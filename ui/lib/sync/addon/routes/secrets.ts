/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';

interface ConfigResponse {
  data: {
    disabled: boolean;
  };
}

export default class SyncSecretsRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;

  model() {
    return hash({
      featureEnabled: this.store
        .adapterFor('application')
        .ajax('/v1/sys/sync/config', 'GET')
        .then(({ data: { disabled } }: ConfigResponse) => !disabled),
    });
  }

  afterModel(model: { featureEnabled: boolean }) {
    if (!model.featureEnabled) {
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
    }
  }
}
