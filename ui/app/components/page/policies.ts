/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import Component from '@glimmer/component';
import { WIZARD_ID } from 'vault/components/wizard/acl-policies/acl-wizard';
import errorMessage from 'vault/utils/error-message';
import { PolicyTypes } from 'core/utils/code-generators/policy';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type NamespaceService from 'vault/services/namespace';
import type RouterService from '@ember/routing/router-service';
import type WizardService from 'vault/services/wizard';
import type PolicyModel from 'vault/vault/models/policy';
import type { PaginatedMetadata } from 'core/utils/paginate-list';

interface Args {
  filter: string | null;
  model: PolicyModel[] & PaginatedMetadata;
  policyType: string;
  onRefresh: CallableFunction;
}

export default class PagePoliciesComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly wizard: WizardService;
  @service declare readonly namespace: NamespaceService;

  @tracked filter = '';
  @tracked filterFocused = false;
  // set when clicking 'Delete' from popup menu
  @tracked policyToDelete = null;
  @tracked shouldRenderIntroModal = false;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.filter = this.args.filter || '';
  }

  // Check if the filter exactly matches a policy ID
  get filterMatchesKey(): boolean {
    const filter = this.filter;
    const content = this.args.model;
    return !!(content && content.length && content.find((c: PolicyModel) => c['id'] === filter));
  }

  // Find the first policy that partially matches the filter (starts with filter)
  get firstPartialMatch(): PolicyModel | undefined {
    const filter = this.filter;
    const content = this.args.model;
    if (!content) {
      return undefined;
    }
    const filterMatchesKey = this.filterMatchesKey;
    const re = new RegExp('^' + filter);
    return filterMatchesKey
      ? undefined
      : content.find((key: PolicyModel) => {
          return re.test(key['id'] as string);
        });
  }

  // starting policies are 'default' and if in the root namespace, 'root' or 'hcp-root'
  get hasOnlyDefaultPolicies() {
    const expectedLength = this.namespace.inRootNamespace ? 2 : 1;
    return this.args.model.meta?.total <= expectedLength;
  }

  // callback from HDS pagination to set the queryParams page
  get paginationQueryParams() {
    return (page: number) => {
      return {
        page,
      };
    };
  }

  get showContent() {
    // Show when the 1) wizard is not shown OR 2) wizard intro modal is shown
    // This ensures the wizard intro modal is shown on top of the list view and the background content is not blank behind the modal
    return !this.showWizard || (this.shouldRenderIntroModal && this.wizard.isIntroVisible(WIZARD_ID));
  }

  get showIntroButton() {
    return this.args.policyType === PolicyTypes.ACL && this.showContent && this.hasOnlyDefaultPolicies;
  }

  // Show when it is not in a dismissed state and there are no non-default policies and
  get showWizard() {
    if (this.args.policyType !== 'acl') return false;
    // Use total instead of filtered total to avoid flashing wizard when filtering with no results
    return !this.wizard.isDismissed(WIZARD_ID) && this.hasOnlyDefaultPolicies;
  }

  @action
  async deletePolicy(policyToDelete: PolicyModel) {
    try {
      const policyName = policyToDelete.name;
      const policyType = this.args.policyType;

      // Use the appropriate sys endpoint based on policy type
      if (policyType === 'egp') {
        await this.api.sys.systemDeletePoliciesEgpName(policyName);
      } else if (policyType === 'rgp') {
        await this.api.sys.systemDeletePoliciesRgpName(policyName);
      } else {
        await this.api.sys.policiesDeleteAclPolicy(policyName);
      }

      // Log success and optionally update the UI
      this.flashMessages.success(`Successfully deleted policy: ${policyName}`);

      // Call the refresh method to update the list
      this.refreshPolicyList();
    } catch (error) {
      const message = errorMessage(error);
      this.flashMessages.danger(message);
    }
    this.policyToDelete = null;
  }

  @action
  setFilter(val: string) {
    this.filter = val;
  }

  @action
  setFilterFocus(bool: boolean) {
    this.filterFocused = bool;
  }

  @action
  showIntroPage() {
    // Reset the wizard dismissal state to allow re-entering the wizard
    this.wizard.reset(WIZARD_ID);
    this.shouldRenderIntroModal = true;
  }

  @action
  refreshPolicyList() {
    try {
      this.args.onRefresh();
    } catch (error) {
      this.flashMessages.danger('There was an error refreshing the policy list.');
    }
  }
}
