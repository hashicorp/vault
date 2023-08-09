/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  'data-test-radial-progress': true,
  tagName: 'svg',
  classNames: 'radial-progress',
  attributeBindings: ['size:width', 'size:height', 'viewBox', 'data-test-radial-progress'],
  progressDecimal: null,
  size: 20,
  strokeWidth: 1,

  viewBox: computed('size', function () {
    const s = this.size;
    return `0 0 ${s} ${s}`;
  }),
  centerValue: computed('size', function () {
    return this.size / 2;
  }),
  r: computed('size', 'strokeWidth', function () {
    return (this.size - this.strokeWidth) / 2;
  }),
  c: computed('r', function () {
    return 2 * Math.PI * this.r;
  }),
  dashArrayOffset: computed('c', 'progressDecimal', function () {
    return this.c * (1 - this.progressDecimal);
  }),
});
