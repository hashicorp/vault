/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { formatExportData, formatQueryParams } from 'core/utils/client-counts/serializers';

import type AdapterError from '@ember-data/adapter/error';
import type ApiService from 'vault/services/api';
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
  @service declare readonly api: ApiService;
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

  async fetchAndFormatExportData(startTimestamp: string | undefined, endTimestamp: string | undefined) {
    // The "Client List" tab is only available on enterprise versions
    // For now, it is also hidden on HVD managed clusters
    if (this.version.isEnterprise && !this.flags.isHvdManaged) {
      const { start_time, end_time } = formatQueryParams({
        start_time: startTimestamp,
        end_time: endTimestamp,
      });
      let exportData, cannotRequestExport;
      try {
        const { raw } = await this.api.sys.internalClientActivityExportRaw({
          end_time,
          format: 'json', // the API only accepts json or csv
          start_time,
        });

        // If it's not a 200 but didn't throw an error then it's likely a 204 (empty response).
        exportData = raw.status === 200 ? await formatExportData(raw, { isDownload: false }) : null;
      } catch (e) {
        const { status, path, response } = await this.api.parseError(e);
        // Show a custom error message when the user does not have permission
        if (status === 403) {
          cannotRequestExport = true;
        } else {
          // re-throw if not a permissions error
          throw {
            httpStatus: status,
            path,
            message: response?.message,
            errors: response?.errors || [],
            error: response?.error,
          };
        }
      }
      return { exportData, cannotRequestExport };
    }
    return { exportData: null, exportError: null };
  }

  async model(params: ClientsCountsRouteParams) {
    const { config, versionHistory } = this.modelFor('vault.cluster.clients') as ModelFrom<ClientsRoute>;
    const { activity, activityError } = await this.getActivity(params);
    const { exportData, cannotRequestExport } = await this.fetchAndFormatExportData(
      activity?.startTime,
      activity?.endTime
    );
    return {
      activity,
      activityError,
      cannotRequestExport,
      config,
      exportData,
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
