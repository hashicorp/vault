/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */
import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';
import type { Breadcrumb } from 'vault/vault/app-types';

export default class VaultClusterSupportUpgradeInfoController extends Controller {
  queryParams = ['tab'];
  @tracked tab = '0';
  @tracked breadcrumbs: Array<Breadcrumb> = [];
  //TODO: upgradeInfo is tracked here so the tab change doesn't rebuild HDS tabs
  //When endpoint is set up, we will fetch in the parent model so re-analyzing still refreshes the data.
  @tracked upgradeInfo: unknown[] | null = null;
}
