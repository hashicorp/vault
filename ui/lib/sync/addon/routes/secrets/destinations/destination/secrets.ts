/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type StoreService from 'vault/services/store';
import type SyncDestinationModel from 'vault/vault/models/sync/destination';
import type SyncAssociationModel from 'vault/vault/models/sync/association';
import type Controller from '@ember/controller';

interface SyncDestinationSecretsRouteParams {
  page: string;
}

interface SyncDestinationSecretsRouteModel {
  destination: SyncDestinationModel;
  associations: SyncAssociationModel[];
}

interface SyncDestinationSecretsController extends Controller {
  model: SyncDestinationSecretsRouteModel;
  page: number | undefined;
}

export default class SyncDestinationSecretsRoute extends Route {
  @service declare readonly store: StoreService;

  queryParams = {
    page: {
      refreshModel: true,
    },
  };

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

  resetController(controller: SyncDestinationSecretsController, isExiting: boolean) {
    if (isExiting) {
      controller.set('page', undefined);
    }
  }
}
