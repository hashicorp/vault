/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

import type ApiService from 'vault/services/api';
import type FlashMessageService from 'vault/services/flash-messages';
import type {
  IdentityGroupAliasesIndexModel,
  GroupAliasListItem,
} from 'vault/routes/vault/cluster/access/identity/groups/aliases/index';

interface Args {
  model: IdentityGroupAliasesIndexModel;
  page: number;
  pageSize: number;
}

export default class PageIdentityGroupAliasesComponent extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;

  @tracked aliasToDelete: GroupAliasListItem | null = null;
  // Optimistic local copy — updated on delete so the list updates without a route refresh.
  @tracked localAliases: GroupAliasListItem[] | null = null;

  get aliases(): GroupAliasListItem[] {
    return this.localAliases ?? this.args.model.aliases;
  }

  @action
  setAliasToDelete(item: GroupAliasListItem) {
    this.aliasToDelete = item;
  }

  @action
  async deleteAlias(item: GroupAliasListItem) {
    try {
      await this.api.identity.groupDeleteAliasById(item.id);
      this.localAliases = this.aliases.filter((a) => a.id !== item.id);
      this.flashMessages.success(`Successfully deleted alias: ${item.name}`);
    } catch (err) {
      const { message } = await this.api.parseError(err);
      this.flashMessages.danger(message);
    }
    this.aliasToDelete = null;
  }
}
