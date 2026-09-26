/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type RouterService from '@ember/routing/router-service';
import type { Capabilities } from 'vault/app-types';
import type { PoliciesRgpIndexModel } from 'vault/routes/vault/cluster/policies/rgp/index';

/**
 * @module Page::RgpPolicies
 *
 * Renders the RGP policies list view. Wraps Page::ListView with the
 * policies-specific named blocks (customTableItem for the policy name link,
 * popupMenu with capability gating) and owns the download-at-click-time and
 * delete confirmation modal.
 *
 * @param {PoliciesRgpIndexModel} model - route model: policies array, listViewConfig, page, pageSize
 * @param {number} page - current page number from the controller query param
 * @param {number} pageSize - number of items per page from the controller query param
 */

interface PolicyRow {
  name: string;
  capabilities: Capabilities | null;
}

interface Args {
  model: PoliciesRgpIndexModel;
  page: number;
  pageSize: number;
}

export default class PageRgpPoliciesComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;

  @tracked policyToDelete: PolicyRow | null = null;

  @action
  async downloadPolicy(item: PolicyRow) {
    try {
      const res = await this.api.sys.systemReadPoliciesRgpName(item.name);
      const policy = (res as unknown as { data: { policy: string } }).data?.policy ?? '';
      const blob = new Blob([policy], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${item.name}.sentinel`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (err) {
      const { message } = await this.api.parseError(err);
      this.flashMessages.danger(message);
    }
  }

  @action
  setPolicyToDelete(item: PolicyRow) {
    this.policyToDelete = item;
  }

  @action
  async deletePolicy(item: PolicyRow) {
    try {
      await this.api.sys.systemDeletePoliciesRgpName(item.name);
      this.flashMessages.success(`Successfully deleted policy: ${item.name}`);
      this.router.refresh();
    } catch (err) {
      const { message } = await this.api.parseError(err);
      this.flashMessages.danger(message);
    }
    this.policyToDelete = null;
  }
}
