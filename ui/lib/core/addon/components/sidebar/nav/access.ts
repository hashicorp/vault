/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { NAV_AUTH_METHODS, NAV_ACL_POLICIES, NAV_NAMESPACES } from 'vault/utils/analytic-events';

import type { AnalyticsEventName } from 'vault/utils/analytic-events';
import type AnalyticsService from 'vault/services/analytics';
import type RouterService from '@ember/routing/router-service';

export default class SidebarNavAccessComponent extends Component {
  @service declare readonly analytics: AnalyticsService;
  @service declare readonly router: RouterService;

  // Returns the policy type segment from the current URL (e.g. "acl", "rgp", "egp"),
  // or null when the current route is not a policy show/edit route.
  get currentPolicyType(): string | null {
    const match = this.router.currentURL?.match(/\/policy\/(acl|rgp|egp)\//);
    return match ? match[1] ?? null : null;
  }

  navEvents = {
    aclPolicies: NAV_ACL_POLICIES,
    authMethods: NAV_AUTH_METHODS,
    namespaces: NAV_NAMESPACES,
  };

  @action
  trackNavClick(eventName: AnalyticsEventName, elementId: string, cta: string) {
    this.analytics.trackEvent(eventName, {
      namespace: 'nav',
      action: 'clicked',
      elementId,
      CTA: cta,
      channel: 'webpage',
    });
  }
}
