/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

import type { UnauthMountsByType } from 'vault/vault/auth/form';

interface Args {
  authTabData: UnauthMountsByType;
  handleTabClick: CallableFunction;
  selectedAuthMethod: string;
}

export default class AuthTabs extends Component<Args> {
  get tabTypes() {
    return this.args.authTabData ? Object.keys(this.args.authTabData) : [];
  }

  get selectedTabIndex() {
    const index = this.tabTypes.indexOf(this.args.selectedAuthMethod);
    // negative index means the selected method isn't a tab, default to first tab
    return index < 0 ? 0 : index;
  }

  @action
  onClickTab(_event: Event, idx: number) {
    this.args.handleTabClick(this.tabTypes[idx]);
  }
}
