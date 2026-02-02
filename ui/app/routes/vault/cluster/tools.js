/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';

export default class ToolsRoute extends Route {
  model() {
    return this.modelFor('vault.cluster');
  }
}
