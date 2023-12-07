/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

import type StoreService from 'vault/services/store';

@withConfirmLeave()
export default class SyncSecretsDestinationsEditDestinationRoute extends Route {
  @service declare readonly store: StoreService;
}
