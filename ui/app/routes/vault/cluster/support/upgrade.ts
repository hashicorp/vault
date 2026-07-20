/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */
import Route from '@ember/routing/route';
import { action } from '@ember/object';
import { service } from '@ember/service';
import type Controller from '@ember/controller';
import type { Breadcrumb } from 'vault/vault/app-types';
import type RouterService from '@ember/routing/router-service';
import type VersionService from 'vault/services/version';

interface RouteController extends Controller {
  breadcrumbs: Array<Breadcrumb>;
  setUpgradeInfo: (info: unknown[]) => void;
}

interface UpgradeModel {
  upgradeInfo: unknown[] | null;
}

export default class VaultClusterSupportUpgradeRoute extends Route {
  @service declare readonly router: RouterService;
  @service declare readonly version: VersionService;

  beforeModel() {
    if (this.version.isCommunity) {
      this.router.transitionTo('vault.cluster.dashboard');
    }
  }

  model(): UpgradeModel {
    return { upgradeInfo: null };
  }

  setupController(controller: RouteController, resolvedModel: UpgradeModel) {
    super.setupController(controller, resolvedModel);

    controller.breadcrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      { label: 'Support', route: 'vault.cluster.support.upgrade' },
      { label: 'Upgrade path analyzer' },
    ];

    controller.setUpgradeInfo = this.setUpgradeInfo;
  }

  @action
  setUpgradeInfo(info: unknown[]) {
    const model = this.modelFor(this.routeName) as UpgradeModel;
    model.upgradeInfo = info;
  }
}
