import { computed } from '@ember/object';
import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import utils from 'vault/lib/key-utils';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';
import WithNavToNearestAncestor from 'vault/mixins/with-nav-to-nearest-ancestor';
import ListController from 'vault/mixins/list-controller';

export default Controller.extend(ListController, BackendCrumbMixin, WithNavToNearestAncestor, {
  flashMessages: service(),
  queryParams: ['page', 'pageFilter', 'tab'],

  tab: '',

  filterIsFolder: computed('filter', function() {
    return !!utils.keyIsFolder(this.get('filter'));
  }),

  actions: {
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
