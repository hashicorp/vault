/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { inject as controller } from '@ember/controller';

import { TOGGLE_WEB_REPL } from 'vault/utils/analytic-events';

export default class SidebarNavComponent extends Component {
  @service analytics;
  @service currentCluster;
  @service console;
  @controller('vault.cluster') clusterController;

  trackReplToggle = () => {
    this.analytics.trackEvent(TOGGLE_WEB_REPL);
  };

  closeConsole = (event) => {
    if (event?.key === 'Escape') {
      this.console.isOpen = false;
    }
  };
}
