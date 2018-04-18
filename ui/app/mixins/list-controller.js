import Ember from 'ember';

export default Ember.Mixin.create({
  queryParams: {
    page: 'page',
    pageFilter: 'pageFilter',
  },

  page: 1,
  pageFilter: null,
  filter: null,

  isLoading: false,

  filterMatchesKey: Ember.computed('filter', 'model', 'model.[]', function() {
    var filter = this.get('filter');
    var content = this.get('model');
    return !!(content.length && content.findBy('id', filter));
  }),

  firstPartialMatch: Ember.computed('filter', 'model', 'model.[]', 'filterMatchesKey', function() {
    var filter = this.get('filter');
    var content = this.get('model');
    var filterMatchesKey = this.get('filterMatchesKey');
    var re = new RegExp('^' + filter);
    return filterMatchesKey
      ? null
      : content.find(function(key) {
          return re.test(key.get('id'));
        });
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
