import { computed } from '@ember/object';
import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import utils from 'vault/lib/key-utils';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';
import WithNavToNearestAncestor from 'vault/mixins/with-nav-to-nearest-ancestor';

export default Controller.extend(BackendCrumbMixin, WithNavToNearestAncestor, {
  flashMessages: service(),
  queryParams: ['page', 'pageFilter', 'tab'],

  tab: '',
  page: 1,
  pageFilter: null,
  filterFocused: false,

  // set via the route `loading` action
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
    var re = new RegExp('^' + filter);
    return filterMatchesKey
      ? null
      : content.find(function(key) {
          return re.test(key.get('id'));
        });
  }),

  filterIsFolder: computed('filter', function() {
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

    delete(item, type) {
      const name = item.id;
      item.destroyRecord().then(() => {
        this.get('flashMessages').success(`${name} was successfully deleted.`);
        this.send('reload');
        if (type === 'secret') {
          this.navToNearestAncestor.perform(name);
        }
      });
    },
  },
});
