/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ActivityComponent from '../activity';
import { service } from '@ember/service';
import { sanitizePath } from 'core/utils/sanitize-path';
import type FlagsService from 'vault/services/flags';

export default class ClientsOverviewPageComponent extends ActivityComponent {
  @service declare readonly flags: FlagsService;

  get hasAttributionData() {
    // we hide attribution data when mountPath filter present
    // or if there's no data
    if (this.args.mountPath || !this.totalUsageCounts.clients) return false;
    return true;
  }

  // mounts attribution
  get namespaceMountAttribution() {
    const { activity } = this.args;
    const nsLabel = this.namespacePathForFilter;
    return activity?.byNamespace?.find((ns) => sanitizePath(ns.label) === nsLabel)?.mounts || [];
  }
}
