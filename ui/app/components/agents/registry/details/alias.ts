/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

import type { FlyoutData } from '../page';

interface Args {
  data: FlyoutData;
}

export default class AgentsRegistryDetailsAliasComponent extends Component<Args> {
  @tracked currentAliasIndex = 0;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.syncAliasIndex();
  }

  get aliases() {
    return this.args.data.entity?.aliases ?? [];
  }

  get aliasCount() {
    return this.aliases.length;
  }

  get selectedAlias() {
    return this.aliases[this.currentAliasIndex];
  }

  get selectedAliasPosition() {
    return this.currentAliasIndex + 1;
  }

  get isPreviousDisabled() {
    return this.currentAliasIndex === 0;
  }

  get isNextDisabled() {
    return this.currentAliasIndex >= this.aliasCount - 1;
  }

  @action
  syncAliasIndex() {
    if (!this.aliasCount) {
      this.currentAliasIndex = 0;
      return;
    }

    const selectedAliasId = this.args.data.aliasId;
    if (selectedAliasId) {
      const selectedAliasIndex = this.aliases.findIndex((alias) => alias.id === selectedAliasId);
      if (selectedAliasIndex >= 0) {
        this.currentAliasIndex = selectedAliasIndex;
        return;
      }
    }

    if (this.currentAliasIndex > this.aliasCount - 1) {
      this.currentAliasIndex = this.aliasCount - 1;
    }
  }

  @action
  selectPreviousAlias() {
    if (!this.isPreviousDisabled) {
      this.currentAliasIndex -= 1;
    }
  }

  @action
  selectNextAlias() {
    if (!this.isNextDisabled) {
      this.currentAliasIndex += 1;
    }
  }
}
