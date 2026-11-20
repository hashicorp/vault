/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';
import { INTRO_REOPEN_CLICKED } from 'vault/utils/analytic-events';

import type RouterService from '@ember/routing/router-service';
import type ThemeService from 'vault/services/theme';
import type WizardService from 'vault/services/wizard';
import type AnalyticsService from 'vault/services/analytics';

export default class MfaIndexController extends Controller {
  @service declare readonly router: RouterService;
  @service declare readonly theme: ThemeService;
  @service declare readonly wizard: WizardService;
  @service declare readonly analytics: AnalyticsService;

  @tracked shouldRenderIntroModal = false;
  wizardId = WIZARD_ID_MAP.mfa;

  // This landing page only renders when no MFA methods exist.
  get showWizard() {
    return !this.wizard.isDismissed(this.wizardId);
  }

  get showContent() {
    return !this.showWizard || (this.shouldRenderIntroModal && this.wizard.isIntroVisible(this.wizardId));
  }

  @action
  showIntroPage() {
    this.analytics.trackEvent(INTRO_REOPEN_CLICKED, {
      namespace: 'intro-page',
      action: 'clicked',
      elementId: 'intro-reopen-button',
      channel: 'webpage',
      objectType: 'mfa',
    });
    this.wizard.reset(this.wizardId);
    this.shouldRenderIntroModal = true;
  }

  @action
  refreshRoute() {
    this.router.refresh('vault.cluster.access.mfa.index');
  }
}
