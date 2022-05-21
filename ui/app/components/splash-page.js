import { inject as service } from '@ember/service';
import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  version: service(),
  auth: service(),
  store: service(),
  tagName: '',
  showTruncatedNavBar: true,

  activeCluster: computed('auth.activeCluster', function () {
    return this.store.peekRecord('cluster', this.auth.activeCluster);
  }),
});
