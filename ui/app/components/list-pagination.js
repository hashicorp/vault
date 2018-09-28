import { gt } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { range } from 'ember-composable-helpers/helpers/range';

export default Component.extend({
  classNames: ['box', 'is-shadowless', 'list-pagination'],
  page: null,
  lastPage: null,
  link: null,
  model: null,
  // number of links to show on each side of page
  spread: 2,
  hasNext: computed('page', 'lastPage', function() {
    return this.get('page') < this.get('lastPage');
  }),
  hasPrevious: computed('page', 'lastPage', function() {
    return this.get('page') > 1;
  }),

  segmentLinks: gt('lastPage', 10),

  pageRange: computed('page', 'lastPage', function() {
    const { spread, page, lastPage } = this.getProperties('spread', 'page', 'lastPage');

    let lower = Math.max(2, page - spread);
    let upper = Math.min(lastPage - 1, lower + spread * 2);
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
