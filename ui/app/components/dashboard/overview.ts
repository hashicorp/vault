/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import { action } from '@ember/object';
import { service } from '@ember/service';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { CLUSTER_STARTUP_CHECKLIST } from 'vault/utils/constants/checklist';

import type ChecklistStateService from 'vault/services/checklist-state';
import NamespaceService from 'vault/services/namespace';
import type ThemeService from 'vault/services/theme';

export type Args = {
  isRootNamespace: boolean;
  replication: unknown;
  replicationUpdatedAt: unknown;
  secretsEngines: unknown;
  vaultConfiguration: unknown;
  version: { isEnterprise: boolean; hasPKIOnly: boolean; hasConsumptionBilling: boolean };
};

export default class OverviewComponent extends Component<Args> {
  @service declare readonly namespace: NamespaceService;
  @service('checklist-state') declare readonly checklistState: ChecklistStateService;
  @service declare readonly theme: ThemeService;
  @tracked showSetupGuideInCompleteState = false;

  /** Returns the dark-mode variant of the agent registry dashboard image when dark mode is active. */
  get featureSpotlightImageSrc(): string {
    const base = '/ui/images/agent-registry-dashboard.png';
    return this.theme.isDarkMode ? '/ui/images/agent-registry-dashboard-dark.png' : base;
  }

  get startupChecklist() {
    return CLUSTER_STARTUP_CHECKLIST;
  }

  /**
   * the client count card should show in the following conditions
   * Self Managed clusters that are running enterprise and showing the `root` namespace
   * Managed clusters that are running enterprise and show the `admin` namespace
   */
  // for self managed clusters, this is the `root` namespace
  // for HVD clusters, this is the `admin` namespace
  get shouldShowClientCount() {
    const { version, isRootNamespace } = this.args;
    const { namespace } = this;

    // don't show client count if this isn't an enterprise cluster
    if (!version.isEnterprise) return false;

    // don't show client count if this is a PKI-only Secrets cluster
    if (version.hasPKIOnly) return false;

    // don't show client count if this is a Consumption Billing cluster
    if (version.hasConsumptionBilling) return false;

    // HVD clusters
    if (namespace.inHvdAdminNamespace) return true;

    // SM clusters
    if (isRootNamespace) return true;

    return false;
  }

  /** Checklist currently supports enterprise root (self-managed) namespace only. */
  get isChecklistNamespaceScope(): boolean {
    return this.args.isRootNamespace;
  }

  /** Ordered step IDs visible to the current user based on permissions and cluster type. */
  get visibleSteps(): string[] {
    return this.checklistState.getVisibleSteps(this.startupChecklist.id, [...this.startupChecklist.order]);
  }

  /**
   * The checklist lifecycle state to render in the left column, or null when
   * the checklist should not be shown (CE, fetch failed, or no visible steps).
   */
  get checklistLifecycleState() {
    if (!this.checklistState.isEnterpriseFeature) return null;
    if (!this.isChecklistNamespaceScope) return null;
    if (!this.checklistState.hasChecklistEntryAccess) return null;
    if (!this.checklistState.isAvailable) return null;
    if (this.visibleSteps.length === 0) return null;
    return this.checklistState.getLifecycleState(this.startupChecklist.id, this.visibleSteps);
  }

  /**
   * Allows users to temporarily return to the checklist after completion
   * without mutating completion state.
   */
  get isShowingCompletedChecklist(): boolean {
    return this.checklistLifecycleState === 'complete' && this.showSetupGuideInCompleteState;
  }

  /** True when the checklist widget should be rendered in the left column. */
  get shouldShowChecklist(): boolean {
    return this.checklistLifecycleState === 'active' || this.isShowingCompletedChecklist;
  }

  @action hideChecklist() {
    // Remember whether the checklist or the completion panel was showing so
    // restoring later (from the Explore Vault card) returns to the same view.
    this.checklistState.setLastView(
      this.startupChecklist.id,
      this.shouldShowChecklist ? 'checklist' : 'complete-banner'
    );
    this.showSetupGuideInCompleteState = false;
    this.checklistState.hideChecklist(this.startupChecklist.id);
  }

  @action restoreChecklist() {
    const lastView = this.checklistState.getLastView(this.startupChecklist.id);
    this.showSetupGuideInCompleteState = lastView === 'checklist';
    this.checklistState.showChecklist(this.startupChecklist.id);
  }

  @action showSetupGuide() {
    this.showSetupGuideInCompleteState = true;
  }

  @action
  handleChecklistStepCompletionChange() {
    if (this.checklistLifecycleState !== 'complete') {
      this.showSetupGuideInCompleteState = false;
    }
  }
}
