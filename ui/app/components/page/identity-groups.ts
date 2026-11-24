/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import type RouterService from '@ember/routing/router-service';
import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type {
  IdentityGroupsIndexModel,
  GroupListItem,
} from 'vault/routes/vault/cluster/access/identity/groups/index';
import WizardService from 'vault/services/wizard';
import AnalyticsService from 'vault/services/analytics';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';
import { INTRO_REOPEN_CLICKED } from 'vault/utils/analytic-events';

interface Args {
  model: IdentityGroupsIndexModel;
  page: number;
  pageSize: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  onSortChange?: (sortBy: string, sortOrder: 'asc' | 'desc') => void;
}

export default class PageIdentityGroupsComponent extends Component<Args> {
  @service declare readonly analytics: AnalyticsService;
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;
  @service declare readonly wizard: WizardService;

  @tracked groupToDelete: GroupListItem | null = null;
  // Optimistic local copy — updated on delete so the list updates without a route refresh.
  @tracked localGroups: GroupListItem[] | null = null;

  @tracked shouldRenderIntroModal = false;
  wizardId = WIZARD_ID_MAP.identityGroups;

  get canCreateGroup(): boolean {
    return this.args.model.canCreateGroup;
  }

  get groups(): GroupListItem[] {
    return this.localGroups ?? this.args.model.groups;
  }

  get isInitialState() {
    return !(this.groups.length > 0);
  }

  get showWizard() {
    return !this.wizard.isDismissed(this.wizardId) && this.isInitialState;
  }

  get showContent() {
    return !this.showWizard || (this.shouldRenderIntroModal && this.wizard.isIntroVisible(this.wizardId));
  }

  get showIntroButton() {
    return this.showContent && this.isInitialState;
  }

  @action
  setGroupToDelete(item: GroupListItem) {
    this.groupToDelete = item;
  }

  @action
  showIntroPage() {
    this.analytics.trackEvent(INTRO_REOPEN_CLICKED, {
      namespace: 'intro-page',
      action: 'clicked',
      elementId: 'intro-reopen-button',
      channel: 'webpage',
      objectType: 'identity-group',
    });
    // Reset the wizard dismissal state to allow re-entering the wizard
    this.wizard.reset(this.wizardId);
    this.shouldRenderIntroModal = true;
  }

  @action
  refreshList() {
    this.router.refresh('vault.cluster.access.identity.groups.index');
  }

  @action
  async deleteGroup(item: GroupListItem) {
    try {
      await this.api.identity.groupDeleteById(item.id);
      this.localGroups = this.groups.filter((g) => g.id !== item.id);
      this.flashMessages.success(`Successfully deleted group: ${item.name}`);
    } catch (err) {
      const { message } = await this.api.parseError(err);
      this.flashMessages.danger(message);
    }
    this.groupToDelete = null;
  }
}
