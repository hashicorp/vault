/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { NAV_AUTH_METHODS, NAV_ACL_POLICIES, NAV_NAMESPACES } from 'vault/utils/analytic-events';
import type AnalyticsService from 'vault/services/analytics';

export default class SidebarNavAccessComponent extends Component {
  @service declare readonly analytics: AnalyticsService;

  navEvents = {
    aclPolicies: NAV_ACL_POLICIES,
    authMethods: NAV_AUTH_METHODS,
    namespaces: NAV_NAMESPACES,
  };

  @action
  trackNavClick(eventName: string, elementId: string, cta: string) {
    this.analytics.trackEvent(eventName, {
      namespace: 'nav',
      action: 'clicked',
      elementId,
      CTA: cta,
      channel: 'webpage',
    });
  }
}
