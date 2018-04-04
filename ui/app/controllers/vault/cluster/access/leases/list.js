import Ember from 'ember';
import utils from 'vault/lib/key-utils';

export default Ember.Controller.extend({
  flashMessages: Ember.inject.service(),
  clusterController: Ember.inject.controller('vault.cluster'),
  queryParams: {
    page: 'page',
    pageFilter: 'pageFilter',
  },

  page: 1,
  pageFilter: null,
  filter: null,

  backendCrumb: Ember.computed(function() {
    return {
      label: 'leases',
      text: 'leases',
      path: 'vault.cluster.access.leases.list-root',
      model: this.get('clusterController.model.name'),
    };
  }),

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

    revokePrefix(prefix, isForce) {
      const adapter = this.model.store.adapterFor('lease');
      const method = isForce ? 'forceRevokePrefix' : 'revokePrefix';
      const fn = adapter[method];
      fn
        .call(adapter, prefix)
        .then(() => {
          return this.transitionToRoute('vault.cluster.access.leases.list-root').then(() => {
            this.get('flashMessages').success(`All of the leases under ${prefix} will be revoked.`);
          });
        })
        .catch(e => {
          const errString = e.errors.join('.');
          this.get('flashMessages').danger(
            `There was an error attempting to revoke the prefix: ${prefix}. ${errString}.`
          );
        });
    },
  },
});
