/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type ApiService from 'vault/services/api';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type { Month } from 'vault/vault/billing/overview';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  pollBillingOverview: ReturnType<typeof import('ember-concurrency').task>;
  fetchBillingMetrics: () => Promise<Month[]>;
  months: Month[];
}

export default class BillingOverviewRoute extends Route {
  @service declare readonly api: ApiService;

  setupController(controller: RouteController, resolvedModel: Month[]) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Billing metrics' },
    ];
  }
}
