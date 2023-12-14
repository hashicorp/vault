/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

import type StoreService from 'vault/services/store';
import SyncDestinationModel from 'vault/vault/models/sync/destination';

interface SyncDestinationSecretsRouteParams {
  page: string;
}

export default class SyncDestinationSecretsRoute extends Route {
  @service declare readonly store: StoreService;

  model(params: SyncDestinationSecretsRouteParams) {
    const destination = this.modelFor('secrets.destinations.destination') as SyncDestinationModel;
    return hash({
      destination,
      associations: this.store.lazyPaginatedQuery('sync/association', {
        responsePath: 'data.keys',
        page: Number(params.page) || 1,
        destinationType: destination.type,
        destinationName: destination.name,
      }),
    });
  }
}
