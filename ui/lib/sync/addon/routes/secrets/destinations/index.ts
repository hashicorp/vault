/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';

import type StoreService from 'vault/services/store';
import type RouterService from '@ember/routing/router-service';
import type { ModelFrom } from 'vault/vault/route';
import type SyncDestinationModel from 'vault/vault/models/sync/destination';
import type Controller from '@ember/controller';

interface SyncSecretsDestinationsIndexRouteParams {
  name: string;
  type: string;
  page: string;
}

interface SyncSecretsDestinationsRouteModel {
  destinations: SyncDestinationModel[];
  nameFilter: string | undefined;
  typeFilter: string | undefined;
}

interface SyncSecretsDestinationsController extends Controller {
  model: SyncSecretsDestinationsRouteModel;
  page: number | undefined;
  name: number | undefined;
  type: number | undefined;
}

export default class SyncSecretsDestinationsIndexRoute extends Route {
  @service declare readonly store: StoreService;
  @service declare readonly router: RouterService;

  queryParams = {
    page: {
      refreshModel: true,
    },
    name: {
      refreshModel: true,
    },
    type: {
      refreshModel: true,
    },
  };

  redirect(model: ModelFrom<SyncSecretsDestinationsIndexRoute>) {
    if (!model.destinations.meta.total) {
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
    }
  }

  filterData(dataset: Array<SyncDestinationModel>, name: string, type: string): Array<SyncDestinationModel> {
    let filteredDataset = dataset;
    const filter = (key: keyof SyncDestinationModel, value: string) => {
      return filteredDataset.filter((model) => {
        return model[key].toLowerCase().includes(value.toLowerCase());
      });
    };
    if (type) {
      filteredDataset = filter('type', type);
    }
    if (name) {
      filteredDataset = filter('name', name);
    }
    return filteredDataset;
  }

  async model(params: SyncSecretsDestinationsIndexRouteParams) {
    const { name, type, page } = params;
    return hash({
      destinations: this.store.lazyPaginatedQuery('sync/destination', {
        page: Number(page) || 1,
        pageFilter: (dataset: Array<SyncDestinationModel>) => this.filterData(dataset, name, type),
        responsePath: 'data.keys',
      }),
      nameFilter: params.name,
      typeFilter: params.type,
    });
  }

  resetController(controller: SyncSecretsDestinationsController, isExiting: boolean) {
    if (isExiting) {
      controller.setProperties({
        page: undefined,
        name: undefined,
        type: undefined,
      });
    }
  }
}
