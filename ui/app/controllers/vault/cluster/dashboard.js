/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import localStorage from 'vault/lib/local-storage';

export default class VaultClusterDashboardController extends Controller {
  @service auth;
  @tracked showNewFeatureModal = true;
  @tracked alwaysHideNewFeaturesModal = false;

  get shouldHideNewFeaturesModal() {
    return localStorage.getItem('alwaysHideNewFeaturesModal');
  }

  @action
  toggleAlwaysHideNewFeaturesModal() {
    this.alwaysHideNewFeaturesModal = !this.alwaysHideNewFeaturesModal;

    this.alwaysHideNewFeaturesModal
      ? localStorage.setItem('alwaysHideNewFeaturesModal', true)
      : localStorage.removeItem('alwaysHideNewFeaturesModal');
  }
}
