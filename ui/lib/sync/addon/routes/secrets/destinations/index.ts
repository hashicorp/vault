/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { SystemListSyncDestinationsListEnum } from '@hashicorp/vault-client-typescript';
import { listDestinationsTransform } from 'sync/utils/api-transforms';
import { paginate } from 'core/utils/paginate-list';

import type PaginationService from 'vault/services/pagination';
import type RouterService from '@ember/routing/router-service';
import type { ModelFrom } from 'vault/vault/route';
import type SyncDestinationModel from 'vault/vault/models/sync/destination';
import type Controller from '@ember/controller';
import type ApiService from 'vault/services/api';
import type CapabilitiesService from 'vault/services/capabilities';
import { ListDestination } from 'vault/vault/sync';

interface SyncSecretsDestinationsIndexRouteParams {
  name?: string;
  type?: string;
  page?: string;
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
  @service declare readonly pagination: PaginationService;
  @service declare readonly api: ApiService;
  @service('app-router') declare readonly router: RouterService;
  @service declare readonly capabilities: CapabilitiesService;

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

  async model(params: SyncSecretsDestinationsIndexRouteParams) {
    const { name = '', type = '', page } = params;
    const response = await this.api.sys.systemListSyncDestinations(SystemListSyncDestinationsListEnum.TRUE);
    // transform and paginate destinations response
    const destinations = listDestinationsTransform(response, name, type);
    const paginatedDestinations = paginate(destinations, { page: Number(page) || 1 });
    // fetch capabilities for destinations
    const paths = paginatedDestinations.map((destination: ListDestination) =>
      this.capabilities.pathFor('syncDestination', destination)
    );
    const capabilities = await this.capabilities.fetch(paths);

    return {
      capabilities,
      destinations: paginatedDestinations,
      nameFilter: name,
      typeFilter: type,
    };
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
