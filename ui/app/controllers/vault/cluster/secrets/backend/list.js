import Ember from 'ember';
import utils from 'vault/lib/key-utils';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';

export default Ember.Controller.extend(BackendCrumbMixin, {
  flashMessages: Ember.inject.service(),
  queryParams: ['page', 'pageFilter', 'tab'],

  tab: '',
  page: 1,
  pageFilter: null,
  filterFocused: false,

  // set via the route `loading` action
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

  filterIsFolder: Ember.computed('filter', function() {
    return !!utils.keyIsFolder(this.get('filter'));
  }),

  actions: {
    setFilter(val) {
      this.set('filter', val);
    },

    setFilterFocus(bool) {
      this.set('filterFocused', bool);
    },

    chooseAction(action) {
      this.set('selectedAction', action);
    },

    toggleZeroAddress(item, backend) {
      item.toggleProperty('zeroAddress');
      this.set('loading-' + item.id, true);
      backend
        .saveZeroAddressConfig()
        .catch(e => {
          item.set('zeroAddress', false);
          this.get('flashMessages').danger(e.message);
        })
        .finally(() => {
          this.set('loading-' + item.id, false);
        });
    },

    delete(item) {
      const name = item.id;
      item.destroyRecord().then(() => {
        this.send('reload');
        this.get('flashMessages').success(`${name} was successfully deleted.`);
      });
    },
  },
});
