/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import type AdapterError from '@ember-data/adapter/error';
import type FlagsService from 'vault/services/flags';
import type NamespaceService from 'vault/services/namespace';
import type Store from '@ember-data/store';
import type VersionService from 'vault/services/version';
import type { ModelFrom } from 'vault/vault/route';
import type ClientsRoute from '../clients';
import type ClientsCountsController from 'vault/controllers/vault/cluster/clients/counts';
import type ClientsActivityModel from 'vault/vault/models/clients/activity';

export interface ClientsCountsRouteParams {
  start_time?: string;
  end_time?: string;
  namespace_path?: string;
  mount_path?: string;
  mount_type?: string;
  month?: string;
}

interface ActivityAdapterQuery {
  start_time: string | undefined;
  end_time: string | undefined;
}

export type ClientsCountsRouteModel = ModelFrom<ClientsCountsRoute>;

export default class ClientsCountsRoute extends Route {
  @service declare readonly flags: FlagsService;
  @service declare readonly namespace: NamespaceService;
  @service declare readonly store: Store;
  @service declare readonly version: VersionService;

  queryParams = {
    // These query params make a new request to the API
    start_time: { refreshModel: true, replace: true },
    end_time: { refreshModel: true, replace: true },
    // These query params just filter client-side data
    namespace_path: { refreshModel: false, replace: true },
    mount_path: { refreshModel: false, replace: true },
    mount_type: { refreshModel: false, replace: true },
    month: { refreshModel: false, replace: true },
  };

  beforeModel() {
    return this.flags.fetchActivatedFlags();
  }

  async getActivity(params: ClientsCountsRouteParams): Promise<{
    activity?: ClientsActivityModel;
    activityError?: AdapterError;
  }> {
    let activity, activityError;
    // if CE without both start time and end time, we want to skip the activity call
    // so that the user is forced to choose a date range
    if (this.version.isEnterprise || (this.version.isCommunity && params.start_time && params.end_time)) {
      const query: ActivityAdapterQuery = {
        start_time: params?.start_time,
        end_time: params?.end_time,
      };
      try {
        activity = await this.store.queryRecord('clients/activity', query);
      } catch (error) {
        activityError = error as AdapterError;
      }
    }
    return {
      activity,
      activityError,
    };
  }

  async model(params: ClientsCountsRouteParams) {
    const { config, versionHistory } = this.modelFor('vault.cluster.clients') as ModelFrom<ClientsRoute>;
    const { activity, activityError } = await this.getActivity(params);
    return {
      activity,
      activityError,
      config,
      // We always want to return the start and end time from the activity response
      // so they serve as the source of truth for the time period of the displayed client count data
      startTimestamp: activity?.startTime,
      endTimestamp: activity?.endTime,
      versionHistory,
    };
  }

  resetController(controller: ClientsCountsController, isExiting: boolean) {
    if (isExiting) {
      controller.setProperties({
        start_time: '',
        end_time: '',
        namespace_path: '',
        mount_path: '',
        mount_type: '',
        month: '',
      });
    }
  }
}
