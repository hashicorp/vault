import { inject as service } from '@ember/service';
import { computed } from '@ember/object';
import Controller, { inject as controller } from '@ember/controller';
import utils from 'vault/lib/key-utils';

export default Controller.extend({
  flashMessages: service(),
  store: service(),
  clusterController: controller('vault.cluster'),
  queryParams: {
    page: 'page',
    pageFilter: 'pageFilter',
  },

  page: 1,
  pageFilter: null,
  filter: null,

  backendCrumb: computed(function() {
    return {
      label: 'leases',
      text: 'leases',
      path: 'vault.cluster.access.leases.list-root',
      model: this.get('clusterController.model.name'),
    };
  }),

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

  emptyTitle: computed('baseKey.id', 'filter', 'filterIsFolder', function() {
    let id = this.get('baseKey.id');
    let filter = this.filter;
    if (id === '') {
      return 'There are currently no leases.';
    }
    if (this.filterIsFolder) {
      if (filter === id) {
        return `There are no leases under &quot;${filter}&quot;.`;
      } else {
        return `We couldn't find a prefix matching &quot;${filter}&quot;.`;
      }
    }
  }),

  actions: {
    setFilter(val) {
      this.set('filter', val);
    },

    setFilterFocus(bool) {
      this.set('filterFocused', bool);
    },

    revokePrefix(prefix, isForce) {
      const adapter = this.get('store').adapterFor('lease');
      const method = isForce ? 'forceRevokePrefix' : 'revokePrefix';
      const fn = adapter[method];
      fn.call(adapter, prefix)
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
