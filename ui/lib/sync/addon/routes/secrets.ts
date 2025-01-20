/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

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

  async fetchActivatedFeatures() {
    return await this.store
      .adapterFor('application')
      .ajax('/v1/sys/activation-flags', 'GET')
      .then((resp: ActivationFlagsResponse) => {
        return resp.data?.activated;
      })
      .catch((error: AdapterError) => {
        return error;
      });
  }

  async model() {
    const activatedFeatures = await this.fetchActivatedFeatures();
    const { isAdapterError } = activatedFeatures;
    return {
      activatedFeatures: isAdapterError ? [] : activatedFeatures,
      adapterError: isAdapterError ? activatedFeatures : null,
    };
  }

  afterModel(model: { activatedFeatures: Array<string> }) {
    if (!model.activatedFeatures) {
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
    }
  }
}
