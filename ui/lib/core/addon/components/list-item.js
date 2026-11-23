/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

export default class ListItemComponent extends Component {
  get link() {
    if (!Array.isArray(this.args.linkParams) || !this.args.linkParams.length) return {};
    const [route, ...models] = this.args.linkParams;
    return { route, models };
  }
}
