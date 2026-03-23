/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import clearModelCache from 'vault/utils/shared-model-boundary';

export default class AccessRoute extends Route {
  @service store;

  modelTypes = ['capabilities', 'control-group', 'identity/group', 'identity/group-alias', 'identity/alias'];

  model() {
    return {};
  }

  deactivate() {
    clearModelCache(this.store, this.modelTypes);
  }
}
