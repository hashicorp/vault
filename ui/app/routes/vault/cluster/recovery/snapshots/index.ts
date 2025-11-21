/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type Controller from '@ember/controller';
import type RouterService from '@ember/routing/router-service';
import { Breadcrumb } from 'vault/vault/app-types';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
}

export default class RecoverySnapshotsIndexRoute extends Route {
  @service declare readonly router: RouterService;

  // There is not a recovery.snapshots.index view because currently only one snapshot can be loaded at a time.
  // Redirect to the parent route so we can reuse its logic and send users to "recovery.snapshots.snapshot.manage"
  // if a snapshot is loaded.
  redirect() {
    this.router.transitionTo('vault.cluster.recovery.snapshots');
  }

  setupController(controller: RouteController, resolvedModel: unknown) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Secrets Recovery' },
    ];
  }
}
