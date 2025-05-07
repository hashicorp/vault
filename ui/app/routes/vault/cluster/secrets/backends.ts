/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type Store from '@ember-data/store';

export default class SecretsBackends extends Route {
  @service declare readonly store: Store;

  model() {
    return this.store.query('secret-engine', {});
  }
}
