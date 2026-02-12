/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import {
  formatExportData,
  formatQueryParams,
  destructureClientCounts,
  formatByMonths,
  formatByNamespace,
} from 'core/utils/client-counts/serializers';
import { ModelFrom } from 'vault/route';
import timestamp from 'core/utils/timestamp';

import type ApiService from 'vault/services/api';
import type FlagsService from 'vault/services/flags';
import type NamespaceService from 'vault/services/namespace';
import type VersionService from 'vault/services/version';
import type { ClientsRouteModel } from '../clients';
import type ClientsCountsController from 'vault/controllers/vault/cluster/clients/counts';
import type {
  ByNamespaceClients,
  NamespaceObject,
  Counts,
  ActivityMonthBlock,
} from 'vault/client-counts/activity-api';

export interface ClientsCountsRouteParams {
  start_time?: string;
  end_time?: string;
  namespace_path?: string;
  mount_path?: string;
  mount_type?: string;
  month?: string;
}

export type ClientsCountsRouteModel = ModelFrom<ClientsCountsRoute>;

export default class ClientsCountsRoute extends Route {
  @service declare readonly api: ApiService;
  @service declare readonly flags: FlagsService;
  @service declare readonly namespace: NamespaceService;
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

  async getActivity(params: ClientsCountsRouteParams) {
    // if CE without both start time and end time, we want to skip the activity call
    // so that the user is forced to choose a date range
    if (this.version.isEnterprise || (this.version.isCommunity && params.start_time && params.end_time)) {
      const response = await this.api.sys.internalClientActivityReportCounts(
        undefined,
        params?.end_time || undefined,
        undefined,
        params?.start_time || undefined
      );
      if (response) {
        return {
          ...response,
          by_namespace: formatByNamespace(response.by_namespace as NamespaceObject[] | null),
          by_month: formatByMonths(response.months as ActivityMonthBlock[]),
          total: destructureClientCounts(response.total as ByNamespaceClients | Counts),
        };
      }
    }
    return undefined;
  }

  async fetchAndFormatExportData(startTimestamp: Date | undefined, endTimestamp: Date | undefined) {
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
          end_time: end_time?.toISOString(),
          format: 'json', // the API only accepts json or csv
          start_time: start_time?.toISOString(),
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
    const { config, versionHistory } = this.modelFor('vault.cluster.clients') as ClientsRouteModel;
    const activity = await this.getActivity(params);
    const { exportData, cannotRequestExport } = await this.fetchAndFormatExportData(
      activity?.start_time,
      activity?.end_time
    );
    return {
      activity,
      cannotRequestExport,
      config,
      exportData,
      // We always want to return the start and end time from the activity response
      // so they serve as the source of truth for the time period of the displayed client count data
      startTimestamp: activity?.start_time,
      endTimestamp: activity?.end_time,
      responseTimestamp: timestamp.now(),
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
