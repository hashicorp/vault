import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import { computed } from '@ember/object';
import config from '../config/environment';

export default Controller.extend({
  queryParams: [
    {
      wrappedToken: 'wrapped_token',
      redirectTo: 'redirect_to',
    },
  ],
  wrappedToken: '',
  redirectTo: '',
  env: config.environment,
  auth: service(),
  store: service(),
  activeCluster: computed('auth.activeCluster', function () {
    const id = this.auth.activeCluster;
    return id ? this.store.peekRecord('cluster', id) : null;
  }),
  activeClusterName: computed('activeCluster', function () {
    const activeCluster = this.activeCluster;
    return activeCluster ? activeCluster.get('name') : null;
  }),
});
