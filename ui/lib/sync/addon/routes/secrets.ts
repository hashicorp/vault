/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type RouterService from '@ember/routing/router-service';
import type StoreService from 'vault/services/store';
import AdapterError from 'ember-data/adapter'; // eslint-disable-line ember/use-ember-data-rfc-395-imports

interface ActivationFlagsResponse {
  data: {
    activated: Array<string>;
  };
}

export default class SyncSecretsRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly store: StoreService;

  async model() {
    const response = await this.store
    .adapterFor('application')
    .ajax('/v1/sys/activation-flags', 'GET')
    .then(({ data: { activated } }: ActivationFlagsResponse) => {
      activated.includes('secrets-sync');
    }).catch((error: AdapterError) => {
      return error;
    })
  
    return hash({
      featureEnabled: response.isAdapterError ? false : true,
      adapterError: response.isAdapterError ? response : false,
    });
  }

  afterModel(model: { featureEnabled: boolean }) {
    if (!model.featureEnabled) {
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
    }
  }
}
