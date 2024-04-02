/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Route from '@ember/routing/route';
import { toolsActions } from 'vault/helpers/tools-actions';

export default class ToolsIndexRouter extends Route {
  @service currentCluster;
  @service router;

  beforeModel(transition) {
    const currentCluster = this.currentCluster.cluster.name;
    const supportedActions = toolsActions();
    if (transition.targetName === this.routeName) {
      transition.abort();
      return this.router.replaceWith('vault.cluster.tools.tool', currentCluster, supportedActions[0]);
    }
  }
}
