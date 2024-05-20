/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';

export default class KvDataFields extends Component {
  @tracked progressDecimal = null;
  @tracked size = 20;
  @tracked strokeWidth = 1;

  get viewBox() {
    const s = this.size;
    return `0 0 ${s} ${s}`;
  }

  get centerValue() {
    return this.size / 2;
  }

  get r() {
    return (this.size - this.strokeWidth) / 2;
  }

  get c() {
    return 2 * Math.PI * this.r;
  }

  get dashArrayOffset() {
    return this.c * (1 - this.progressDecimal);
  }
}
