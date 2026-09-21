/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';
import { WIZARD_ID_MAP } from 'vault/utils/constants/wizard';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type NamespaceService from 'vault/services/namespace';
import type RouterService from '@ember/routing/router-service';
import type WizardService from 'vault/services/wizard';
import type { Capabilities } from 'vault/app-types';
import type { PoliciesAclIndexModel } from 'vault/routes/vault/cluster/policies/acl/index';

/**
 * @module Page::AclPolicies
 *
 * Renders the ACL policies list view. Wraps Page::ListView with the
 * policies-specific named blocks (customTableItem for the two-model name
 * link, popupMenu with capability gating and root/default suppression) and
 * owns the wizard visibility state, download-at-click-time, and delete
 * confirmation modal.
 *
 * @param {PoliciesAclIndexModel} model - route model: policies array, listViewConfig, page, pageSize
 * @param {number} page - current page number from the controller query param
 * @param {number} pageSize - number of items per page from the controller query param
 */

interface PolicyRow {
  name: string;
  capabilities: Capabilities | null;
}

interface Args {
  model: PoliciesAclIndexModel;
  page: number;
  pageSize: number;
}

export default class PageAclPoliciesComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly namespace: NamespaceService;
  @service declare readonly router: RouterService;
  @service declare readonly wizard: WizardService;

  @tracked policyToDelete: PolicyRow | null = null;
  @tracked shouldRenderIntroModal = false;

  wizardId = WIZARD_ID_MAP.aclPolicy;

  // Policies are considered "only defaults" when the total is at or below the
  // expected number of built-in policies (root + default + default-ceiling in root
  // namespace; default + default-ceiling elsewhere).
  get hasOnlyDefaultPolicies() {
    const expectedLength = this.namespace.inRootNamespace ? 3 : 2;
    return (this.args.model?.policies?.length ?? 0) <= expectedLength;
  }

  get showWizard() {
    return !this.wizard.isDismissed(this.wizardId) && this.hasOnlyDefaultPolicies;
  }

  get showContent() {
    return !this.showWizard || (this.shouldRenderIntroModal && this.wizard.isIntroVisible(this.wizardId));
  }

  get showIntroButton() {
    return this.showContent && this.hasOnlyDefaultPolicies;
  }

  @action
  async downloadPolicy(item: PolicyRow) {
    try {
      const { policy } = await this.api.sys.policiesReadAclPolicy(item.name);
      const blob = new Blob([policy ?? ''], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${item.name}.hcl`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      this.flashMessages.danger(errorMessage(err));
    }
  }

  @action
  setPolicyToDelete(item: PolicyRow) {
    this.policyToDelete = item;
  }

  @action
  async deletePolicy(item: PolicyRow) {
    try {
      await this.api.sys.policiesDeleteAclPolicy(item.name);
      this.flashMessages.success(`Successfully deleted policy: ${item.name}`);
      this.router.refresh();
    } catch (err) {
      this.flashMessages.danger(errorMessage(err));
    }
    this.policyToDelete = null;
  }

  @action
  showIntroPage() {
    this.wizard.reset(this.wizardId);
    this.shouldRenderIntroModal = true;
  }

  @action
  refreshRoute() {
    this.router.refresh();
  }
}
