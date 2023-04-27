/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { gt } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { range } from 'ember-composable-helpers/helpers/range';
import { A } from '@ember/array';
import layout from '../templates/components/list-pagination';

export default Component.extend({
  layout,
  classNames: ['box', 'is-shadowless', 'list-pagination'],
  page: null,
  lastPage: null,
  link: null,
  models: A(),
  // number of links to show on each side of page
  spread: 2,
  hasNext: computed('page', 'lastPage', function () {
    return this.page < this.lastPage;
  }),
  hasPrevious: gt('page', 1),

  segmentLinks: gt('lastPage', 10),

  pageRange: computed('lastPage', 'page', 'spread', function () {
    const { spread, page, lastPage } = this;

    let lower = Math.max(2, page - spread);
    const upper = Math.min(lastPage - 1, lower + spread * 2);
    // we're closer to lastPage than the spread
    if (upper - lower < 5) {
      lower = upper - 4;
    }
    if (lastPage <= 10) {
      return range([1, lastPage, true]);
    }
    return range([lower, upper, true]);
  }),
});
