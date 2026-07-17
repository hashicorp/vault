/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */
import Route from '@ember/routing/route';
import { service } from '@ember/service';
import type RouterService from '@ember/routing/router-service';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type { ModelFrom } from 'vault/vault/route';
import type VaultClusterSupportUpgradeRoute from 'vault/routes/vault/cluster/support/upgrade';
import type VersionService from 'vault/services/version';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  upgradeInfo: unknown[] | null;
}

type ParentModel = ModelFrom<VaultClusterSupportUpgradeRoute>;

export default class VaultClusterSupportUpgradeInfoRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly version: VersionService;

  beforeModel() {
    if (this.version.isCommunity) {
      this.router.transitionTo('vault.cluster.dashboard');
      return;
    }

    const parentModel = this.modelFor('vault.cluster.support.upgrade') as ParentModel;

    if (!parentModel?.upgradeInfo) {
      this.router.transitionTo('vault.cluster.support.upgrade');
    }
  }

  setupController(controller: RouteController) {
    super.setupController(controller, {});
    const parentModel = this.modelFor('vault.cluster.support.upgrade') as ParentModel;
    controller.upgradeInfo = parentModel?.upgradeInfo ?? null;
    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Support', route: 'vault.cluster.support.upgrade' },
      { label: 'Upgrade path analyzer', route: 'vault.cluster.support.upgrade' },
      { label: 'Issues' },
    ];
  }
}
