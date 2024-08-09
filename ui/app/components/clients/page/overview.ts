/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';
import { service } from '@ember/service';
import { hasMountsKey, hasNamespacesKey } from 'core/utils/client-count-utils';
import { sanitizePath } from 'core/utils/sanitize-path';
import type FlagsService from 'vault/services/flags';
import type NamespaceService from 'vault/services/namespace';

export default class ClientsOverviewPageComponent extends ActivityComponent {
  @service declare readonly flags: FlagsService;
  @service declare readonly namespace: NamespaceService;

  // new client data for horizontal bar chart @isHistoricalMonth only
  get newClientAttribution() {
    // new client attribution only available in a single, historical month (not a date range or current month)
    if (this.isDateRange || this.isCurrentMonth || !this.newClientCounts) return null;

    const newCounts = this.newClientCounts;
    if (this.args.namespace && hasMountsKey(newCounts)) return newCounts?.mounts;

    if (hasNamespacesKey(newCounts)) return newCounts?.namespaces;

    return null;
  }

  get filteredActivityByNamespace() {
    const { activity, namespace } = this.args;
    const currentNs = this.namespace.currentNamespace;
    const nsLabel = sanitizePath(namespace || currentNs || 'root');
    return activity.byNamespace.find((namespace) => sanitizePath(namespace.label) === nsLabel);
  }

  // TODO: replace with constant namespace & mount attribution
  // total client data for horizontal bar chart in attribution component
  // array of labels + counts for namespace or mounts
  get totalClientAttribution() {
    const { namespace, activity } = this.args;
    if (namespace) {
      return this.filteredActivityByNamespace?.mounts || null;
    } else {
      return activity.byNamespace || null;
    }
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
