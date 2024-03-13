/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { findDestination } from 'core/helpers/sync-destinations';

import type StoreService from 'vault/services/store';
import type { SyncDestinationType } from 'vault/vault/helpers/sync-destinations';

interface Params {
  type: SyncDestinationType;
}

export default class SyncSecretsDestinationsCreateDestinationRoute extends Route {
  @service declare readonly store: StoreService;

  model(params: Params) {
    const { type } = params;
    const defaultValues = findDestination(type)?.defaultValues;
    return this.store.createRecord(`sync/destinations/${type}`, { type, ...defaultValues });
  }
}
