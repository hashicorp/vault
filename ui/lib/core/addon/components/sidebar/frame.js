/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { inject as controller } from '@ember/controller';
import { TOGGLE_WEB_REPL } from 'vault/utils/analytic-events';
import { FEEDBACK_SURVEY_URL } from 'vault/utils/constants/links';

export default class SidebarNavComponent extends Component {
  @service analytics;
  @service currentCluster;
  @service console;
  @controller('vault.cluster') clusterController;

  feedbackSurveyUrl = FEEDBACK_SURVEY_URL;

  trackReplToggle = () => {
    this.analytics.trackEvent(TOGGLE_WEB_REPL, {
      namespace: 'nav',
      action: 'clicked',
      elementId: 'web-repl-toggle',
      channel: 'webpage',
    });
  };

  closeConsole = (event) => {
    if (event?.key === 'Escape') {
      this.console.isOpen = false;
    }
  };
}
