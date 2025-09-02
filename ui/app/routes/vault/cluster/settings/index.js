/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class SettingsIndexRouter extends Route {
  @service router;

  beforeModel(transition) {
    if (transition.targetName === this.routeName) {
      transition.abort();
      return this.router.replaceWith('vault.cluster.settings.mount-secret-backend');
    }
  }
}
