/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';
import { service } from '@ember/service';
import { hasMountsKey, hasNamespacesKey } from 'core/utils/client-count-utils';
import type FlagsService from 'vault/services/flags';

export default class ClientsOverviewPageComponent extends ActivityComponent {
  @service declare readonly flags: FlagsService;

  // (object) single month new client data with total counts and array of
  // either namespaces or mounts
  get newClientCounts() {
    if (this.isDateRange || this.byMonthActivityData.length === 0) {
      return null;
    }

    return this.byMonthActivityData[0]?.new_clients;
  }

  // total client data for horizontal bar chart in attribution component
  get totalClientAttribution() {
    const { namespace, activity } = this.args;
    if (namespace) {
      return this.filteredActivityByNamespace?.mounts || null;
    } else {
      return activity.byNamespace || null;
    }
  }

  // new client data for horizontal bar chart
  get newClientAttribution() {
    // new client attribution only available in a single, historical month (not a date range or current month)
    if (this.isDateRange || this.isCurrentMonth || !this.newClientCounts) return null;

    const newCounts = this.newClientCounts;
    if (this.args.namespace && hasMountsKey(newCounts)) return newCounts?.mounts;

    if (hasNamespacesKey(newCounts)) return newCounts?.namespaces;

    return null;
  }

  get hasAttributionData() {
    const { mountPath, namespace } = this.args;
    if (!mountPath) {
      if (namespace) {
        const mounts = this.filteredActivityByNamespace?.mounts?.map((mount) => ({
          id: mount.label,
          name: mount.label,
        }));
        return mounts && mounts.length > 0;
      }
      return !!this.totalClientAttribution && this.totalUsageCounts && this.totalUsageCounts.clients !== 0;
    }

    return false;
  }
}
