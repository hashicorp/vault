/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

export default class SectionTabs extends Component {
  get tabType() {
    return this.args.tabType || 'authSettings';
  }
}
