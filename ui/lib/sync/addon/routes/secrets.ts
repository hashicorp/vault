/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';
import type AdapterError from '@ember-data/adapter';

interface ActivationFlagsResponse {
  data: {
    activated: Array<string>;
    unactivated: Array<string>;
  };
}

export default class SyncSecretsRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;

  model() {
    return hash({
      featureEnabled: this.store
      .adapterFor('application')
      .ajax('/v1/sys/activation-flags', 'GET')
      .then(({ data: { activated } }: ActivationFlagsResponse) => {
        activated.includes('secrets-sync');
      }).catch((error: AdapterError) => {
        // break out the error in the args to the component
        return error;
      })
    });
  }

  afterModel(model: { featureEnabled: boolean }) {
    if (!model.featureEnabled) {
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
    }
  }
}
