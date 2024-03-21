/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';
import { DEBUG } from '@glimmer/env';

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
    // The read request to the activation-flags endpoint is unauthenticated and root namespace
    // but the POST is not which is why it's not in the NAMESPACE_ROOT_URLS list
    return await this.store
      .adapterFor('application')
      .ajax('/v1/sys/activation-flags', 'GET', { unauthenticated: true, namespace: null })
      .then((resp: ActivationFlagsResponse) => {
        return resp.data?.activated;
      })
      .catch((error: unknown) => {
        if (DEBUG) console.error(error); // eslint-disable-line no-console
        return [];
      });
  }

  async model() {
    const activatedFeatures = await this.fetchActivatedFeatures();
    return {
      activatedFeatures,
    };
  }

  afterModel(model: { activatedFeatures: Array<string> }) {
    if (!model.activatedFeatures) {
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
    }
  }
}
