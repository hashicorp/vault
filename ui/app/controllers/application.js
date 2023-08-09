/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import config from '../config/environment';
import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import localStorage from 'vault/lib/local-storage';
export default class ApplicationController extends Controller {
  @service auth;
  @service store;
  @tracked showNewFeatureModal = true;
  @tracked alwaysHideNewFeaturesModal = false;

  env = config.environment;

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
