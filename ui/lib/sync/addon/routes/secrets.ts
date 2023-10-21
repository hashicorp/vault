/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

import type StoreService from 'vault/services/store';

export default class SyncSecretsRoute extends Route {
  @service declare readonly store: StoreService;

  async model() {
    try {
      await this.store.query('sync/destination', {});
      return { promptConfig: false };
    } catch (error) {
      return { promptConfig: true };
    }
  }
}
