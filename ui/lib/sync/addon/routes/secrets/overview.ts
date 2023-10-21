/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type StoreService from 'vault/services/store';

interface SyncSecretsRouteModel {
  promptConfig: boolean;
}

export default class SyncSecretsOverviewRoute extends Route {
  @service declare readonly store: StoreService;

  model() {
    const model = this.modelFor('secrets') as SyncSecretsRouteModel;
    if (model.promptConfig) {
      return model;
    }
    return 'query all the information for overview route';
  }
}
