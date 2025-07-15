/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { paginate } from 'core/utils/paginate-list';

import type ApiService from 'vault/services/api';
import type { DestinationRouteModel } from '../destination';
import type Controller from '@ember/controller';

type Params = {
  page?: string;
};

export default class SyncDestinationSecretsRoute extends Route {
  @service declare readonly api: ApiService;

  queryParams = {
    page: {
      refreshModel: true,
    },
  };

  async model({ page }: Params) {
    const { destination, capabilities } = this.modelFor(
      'secrets.destinations.destination'
    ) as DestinationRouteModel;

    const {
      associatedSecrets = {},
      storeName,
      storeType,
    } = await this.api.sys.systemReadSyncDestinationsTypeNameAssociations(destination.name, destination.type);

    const associations = Object.values(associatedSecrets).map((association) => ({
      destinationName: storeName,
      destinationType: storeType,
      ...association,
    }));

    return {
      destination,
      capabilities,
      associations: paginate(associations, { page: Number(page) || 1 }),
    };
  }

  resetController(controller: Controller, isExiting: boolean) {
    if (isExiting) {
      controller.set('page', undefined);
    }
  }
}
