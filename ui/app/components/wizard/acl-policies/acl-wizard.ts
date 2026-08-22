/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import Component from '@glimmer/component';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';
import { INTRO_ACL_POLICIES_CTA_CLICKED } from 'vault/utils/analytic-events';

import type WizardService from 'vault/services/wizard';
import type AnalyticsService from 'vault/services/analytics';

interface Args {
  isIntroModal: boolean;
  onRefresh: CallableFunction;
}

export default class WizardNamespacesWizardComponent extends Component<Args> {
  @service declare readonly wizard: WizardService;
  @service declare readonly analytics: AnalyticsService;

  wizardId = WIZARD_ID_MAP.aclPolicy;

  @action
  async onDismiss() {
    this.trackCtaEvent(this.args.isIntroModal ? 'Close' : 'Skip', 'dismissed', 'intro-dismiss-button');
    this.wizard.dismiss(this.wizardId);
    await this.args.onRefresh();
  }

  @action
  trackClickEvent(cta: string) {
    this.trackCtaEvent(cta, 'clicked', 'intro-cta-button');
  }

  // `variation` distinguishes the modal from the full-page intro, which is
  // otherwise only implied by the CTA label.
  private trackCtaEvent(CTA: string, action: string, uiElement: string) {
    this.analytics.trackEvent(INTRO_ACL_POLICIES_CTA_CLICKED, {
      CTA,
      channel: 'webpage',
      location: 'intro-page',
      objectType: 'acl-policy',
      variation: this.args.isIntroModal ? 'modal' : 'page',
      uiElement,
      type: 'Button',
      action,
    });
  }
}
