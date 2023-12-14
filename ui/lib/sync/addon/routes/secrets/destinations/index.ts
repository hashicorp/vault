/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

import type StoreService from 'vault/services/store';
import type SyncDestinationModel from 'vault/vault/models/sync/destination';

interface SyncSecretsDestinationsIndexRouteParams {
  name: string;
  type: string;
  page: string;
}

export default class SyncSecretsDestinationsIndexRoute extends Route {
  @service declare readonly store: StoreService;

  filterData(dataset: Array<SyncDestinationModel>, name: string, type: string): Array<SyncDestinationModel> {
    let filteredDataset = dataset;
    const filter = (key: keyof SyncDestinationModel, value: string) => {
      return dataset.filter((model) => {
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
}
