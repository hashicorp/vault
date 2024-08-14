/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';
import { service } from '@ember/service';
import { ByNamespaceClients, MountClients } from 'core/utils/client-count-utils';
import { sanitizePath } from 'core/utils/sanitize-path';
import type FlagsService from 'vault/services/flags';
import type NamespaceService from 'vault/services/namespace';

export default class ClientsOverviewPageComponent extends ActivityComponent {
  @service declare readonly flags: FlagsService;
  @service declare readonly namespace: NamespaceService;

  get showNewClientAttribution(): boolean {
    return this.args.activity.byMonth.length === 1;
  }

  get newClientNamespaceAttribution(): ByNamespaceClients[] | undefined {
    if (this.showNewClientAttribution) {
      return this.args.activity.byMonth[0]?.namespaces;
    }
    return;
  }

  get newClientMountAttribution(): MountClients[] | undefined {
    if (this.showNewClientAttribution) {
      const currentNs = this.namespace.currentNamespace;
      const nsLabel = sanitizePath(this.args.namespace || currentNs || 'root');
      const singleNamespace = this.newClientNamespaceAttribution?.find(
        (namespace) => sanitizePath(namespace.label) === nsLabel
      );
      return singleNamespace?.mounts;
    }
    return;
  }

  get hasAttributionData() {
    // TODO: update this logic. when do we hide attribution?
    if (this.args.mountPath) return false;
    return true;
    // const { mountPath, namespace } = this.args;
    // if (!mountPath) {
    //   if (namespace) {
    //     const mounts = this.filteredActivityByNamespace?.mounts?.map((mount) => ({
    //       id: mount.label,
    //       name: mount.label,
    //     }));
    //     return mounts && mounts.length > 0;
    //   }
    //   return !!this.totalClientAttribution && this.totalUsageCounts && this.totalUsageCounts.clients !== 0;
    // }

    // return false;
  }

  get namespaceMountAttribution() {
    const { activity, namespace } = this.args;
    const currentNs = this.namespace.currentNamespace;
    const nsLabel = sanitizePath(namespace || currentNs || 'root');
    return activity.byNamespace?.find((ns) => sanitizePath(ns.label) === nsLabel)?.mounts;
  }
}
