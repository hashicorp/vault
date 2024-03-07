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
      activatedFeatures: this.store
        .adapterFor('application')
        .ajax('/v1/sys/activation-flags', 'GET')
        .then((resp: ActivationFlagsResponse) => {
          return resp.data.activated;
        })
        .catch((error: AdapterError) => {
          // we break out this error while passing args to the component and handle the error in the overview template
          return error;
        }),
    });
  }

  afterModel(model: { activatedFeatures: Array<string> | AdapterError }) {
    if (!model.activatedFeatures) {
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
    }
  }
}
