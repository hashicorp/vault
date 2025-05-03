/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

interface Args {
  selectedAuthMethod: string;
  authTabData: AuthTabData;
  authTabTypes: string[];
  handleTabClick: CallableFunction;
}

interface AuthTabData {
  // key is the auth method type
  [key: string]: MountData[];
}

interface MountData {
  path: string;
  type: string;
  description?: string;
  config?: object | null;
}

export default class AuthTabs extends Component<Args> {
  get selectedTabIndex() {
    const index = this.args.authTabTypes.indexOf(this.args.selectedAuthMethod);
    // negative index means the selected method isn't a tab, default to first tab
    return index < 0 ? 0 : index;
  }

  @action
  onClickTab(idx: number) {
    this.args.handleTabClick(this.args.authTabTypes[idx]);
  }
}
