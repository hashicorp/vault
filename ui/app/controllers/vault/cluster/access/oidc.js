/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { capitalize } from '@ember/string';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';
import { INTRO_REOPEN_CLICKED } from 'vault/utils/analytic-events';

export default class OidcConfigureController extends Controller {
  @service router;
  @service wizard;
  @service analytics;

  @tracked header = null;
  @tracked shouldRenderIntroModal = false;

  wizardId = WIZARD_ID_MAP.oidcProvider;

  constructor() {
    super(...arguments);
    this.router.on('routeDidChange', (transition) => this.setHeader(transition));
  }

  setHeader(transition) {
    // set correct header state based on child route
    // when no clients have been created, display create button as call to action
    // list views share the same header with tabs as resource links
    // the remaining routes are responsible for their own header
    const routeName = transition.to.name;
    if (routeName.includes('oidc.index')) {
      this.header = 'cta';
    } else {
      const isList = ['clients', 'assignments', 'keys', 'scopes', 'providers'].find((resource) => {
        return routeName.includes(`${resource}.index`);
      });
      this.header = isList ? 'list' : null;
    }
  }

  get isCta() {
    return this.header === 'cta';
  }

  // True when the user is on the initial "no clients yet" landing page.
  get isInitialState() {
    return this.isCta;
  }

  get showWizard() {
    return !this.wizard.isDismissed(this.wizardId) && this.isInitialState;
  }

  // Show page content when the full-page wizard is not displayed, or when the
  // intro modal is open on top of the existing content.
  get showContent() {
    return !this.showWizard || (this.shouldRenderIntroModal && this.wizard.isIntroVisible(this.wizardId));
  }

  get breadcrumbs() {
    // we check parent for the name as the currentRoute is always "index"
    const route = this.router.currentRoute.parent.localName ?? '';
    const isDefaultRoute = route === 'clients' || route == 'oidc';
    const showApplications = route !== 'oidc' ? ': Applications' : '';

    const baseCrumbs = [
      { label: 'Vault', route: 'vault.cluster.dashboard', icon: 'vault' },
      {
        label: `OIDC provider${showApplications}`,
        route: 'vault.cluster.access.oidc',
        current: isDefaultRoute,
      },
    ];

    // clients is the default view, so in order to match current patterns, we do not include it as a breadcrumb
    if (!isDefaultRoute && route) {
      return [
        ...baseCrumbs,
        {
          label: capitalize(route),
        },
      ];
    }

    return baseCrumbs;
  }

  @action
  showIntroPage() {
    this.analytics.trackEvent(INTRO_REOPEN_CLICKED, {
      namespace: 'intro-page',
      action: 'clicked',
      elementId: 'intro-reopen-button',
      channel: 'webpage',
      objectType: 'oidc-provider',
    });
    // Reset dismissal so the wizard is visible again as a modal
    this.wizard.reset(this.wizardId);
    this.shouldRenderIntroModal = true;
  }

  @action
  refreshRoute() {
    this.router.refresh('vault.cluster.access.oidc');
  }
}
