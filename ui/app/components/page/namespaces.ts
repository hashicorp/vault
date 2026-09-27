/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import Component from '@glimmer/component';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';
import { INTRO_REOPEN_CLICKED } from 'vault/utils/analytic-events';

import type AnalyticsService from 'vault/services/analytics';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type NamespaceService from 'vault/services/namespace';
import type RouterService from '@ember/routing/router-service';
import type WizardService from 'vault/services/wizard';
import type { NamespacesIndexModel } from 'vault/routes/vault/cluster/access/namespaces/index';

/**
 * @module PageNamespaces
 * PageNamespaces component handles the display and management of namespaces,
 * including the namespace wizard for first-time users. Internally it delegates
 * list rendering to Page::ListView.
 *
 * @param {NamespacesIndexModel} model - route model containing namespaces array, listViewConfig, page, pageSize
 * @param {function} onRefresh - callback to refresh the namespace list from the route
 */

interface Args {
  model: NamespacesIndexModel;
  onRefresh: () => void;
}

export default class PageNamespacesComponent extends Component<Args> {
  @service declare readonly analytics: AnalyticsService;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;
  @service declare readonly wizard: WizardService;
  @service declare namespace: NamespaceService;

  @tracked nsToDelete: string | null = null;
  @tracked showSetupAlert = false;
  @tracked shouldRenderIntroModal = false;

  wizardId = WIZARD_ID_MAP.namespace;

  get hasNamespaces() {
    return this.args.model.namespaces.length > 0;
  }

  // Show header and breadcrumbs when viewing the intro page or during the list view.
  // Do not show during Guided Start as that has its own header.
  get showPageHeader() {
    return !this.showWizard || this.wizard.isIntroVisible(this.wizardId);
  }

  get showContent() {
    // Show when 1) wizard is not shown OR 2) wizard intro modal is visible.
    // This ensures the intro modal appears over the list view rather than a blank background.
    return !this.showWizard || (this.shouldRenderIntroModal && this.wizard.isIntroVisible(this.wizardId));
  }

  get showIntroButton() {
    return this.showContent && !this.hasNamespaces;
  }

  get showWizard() {
    return !this.wizard.isDismissed(this.wizardId) && !this.hasNamespaces;
  }

  @action
  switchNamespace(targetNamespace: string) {
    this.router.transitionTo('vault.cluster.dashboard', {
      queryParams: { namespace: targetNamespace },
    });
  }

  @action
  async deleteNamespace(namespaceId: string) {
    try {
      await this.api.sys.systemDeleteNamespacesPath(namespaceId);
      this.flashMessages.success(`Successfully deleted namespace: ${namespaceId}`);
      this.refreshNamespaceList();
    } catch (err) {
      const { message } = await this.api.parseError(err);
      this.flashMessages.danger(message);
    }
    this.nsToDelete = null;
  }

  @action
  async refreshNamespaceList() {
    try {
      await this.namespace.findNamespacesForUser.perform();
      this.args.onRefresh();
    } catch {
      this.flashMessages.danger('There was an error refreshing the namespace list.');
    }
  }

  @action
  showIntroPage() {
    this.analytics.trackEvent(INTRO_REOPEN_CLICKED, {
      namespace: 'intro-page',
      action: 'clicked',
      elementId: 'intro-reopen-button',
      channel: 'webpage',
      objectType: 'namespace',
    });
    this.wizard.reset(this.wizardId);
    this.shouldRenderIntroModal = true;
  }
}
