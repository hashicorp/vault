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
import type {
  IdentityGroupsIndexModel,
  GroupListItem,
} from 'vault/routes/vault/cluster/access/identity/groups/index';

interface Args {
  model: IdentityGroupsIndexModel;
  page: number;
  pageSize: number;
}

export default class PageIdentityGroupsComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked groupToDelete: GroupListItem | null = null;
  // Optimistic local copy — updated on delete so the list updates without a route refresh.
  @tracked localGroups: GroupListItem[] | null = null;

  get groups(): GroupListItem[] {
    return this.localGroups ?? this.args.model.groups;
  }

  @action
  setGroupToDelete(item: GroupListItem) {
    this.groupToDelete = item;
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
