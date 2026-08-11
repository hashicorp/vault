/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type RouterService from '@ember/routing/router-service';
import type VersionService from 'vault/services/version';

export default class AgentsRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly version: VersionService;

  redirect() {
    // if not enterprise redirect to dashboard, otherwise registry is the current landing route
    const route = this.version.isEnterprise ? 'vault.cluster.agents.registry' : 'vault.cluster.dashboard';
    this.router.replaceWith(route);
  }
}
