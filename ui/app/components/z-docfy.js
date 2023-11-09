/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class ZDocfy extends Component {
  @tracked filterValue;

  @action
  filterComponents({ target }) {
    this.filterValue = target.value;
  }

  get componentList() {
    return this.filterValue
      ? this.args.components.filter((c) => c.title.toLowerCase().includes(this.filterValue.toLowerCase()))
      : this.args.components;
  }
}
