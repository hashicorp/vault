import { computed } from '@ember/object';
import Mixin from '@ember/object/mixin';
import escapeStringRegexp from 'escape-string-regexp';
import commonPrefix from 'vault/utils/common-prefix';

export default Mixin.create({
  queryParams: {
    page: 'page',
    pageFilter: 'pageFilter',
  },

  page: 1,
  pageFilter: null,
  filter: null,
  filterFocused: false,

  isLoading: false,

  filterMatchesKey: computed('filter', 'model', 'model.[]', function() {
    var filter = this.get('filter');
    var content = this.get('model');
    return !!(content.length && content.findBy('id', filter));
  }),

  firstPartialMatch: computed('filter', 'model', 'model.[]', 'filterMatchesKey', function() {
    var filter = this.get('filter');
    var content = this.get('model');
    var filterMatchesKey = this.get('filterMatchesKey');
    var re = new RegExp('^' + escapeStringRegexp(filter));
    let matchSet = content.filter(key => re.test(key.id));
    let match = matchSet.firstObject;

    if (filterMatchesKey || !match) {
      return null;
    }

    let sharedPrefix = commonPrefix(content);
    // if we already are filtering the prefix, then next we want
    // the exact match
    if (filter === sharedPrefix || matchSet.length === 1) {
      return match;
    }
    return { id: sharedPrefix };
  }),

  actions: {
    setFilter(val) {
      this.set('filter', val);
    },

    setFilterFocus(bool) {
      this.set('filterFocused', bool);
    },
  },
});
