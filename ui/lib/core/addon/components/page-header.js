/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

export default class PageHeader extends Component {
  get hasLevel() {
    return this.args.hasLevel === undefined ? true : this.args.hasLevel;
  }
}
