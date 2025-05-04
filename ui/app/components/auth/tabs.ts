/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

import type { AuthTabData } from 'vault/vault/auth/form';

interface Args {
  authTabData: AuthTabData;
  authTabTypes: string[];
  handleTabClick: CallableFunction;
  selectedAuthMethod: string;
}

export default class AuthTabs extends Component<Args> {
  get selectedTabIndex() {
    const index = this.args.authTabTypes.indexOf(this.args.selectedAuthMethod);
    // negative index means the selected method isn't a tab, default to first tab
    return index < 0 ? 0 : index;
  }

  @action
  onClickTab(_event: Event, idx: number) {
    this.args.handleTabClick(this.args.authTabTypes[idx]);
  }
}
