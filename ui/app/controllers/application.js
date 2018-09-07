import { inject as service } from '@ember/service';
import Controller from '@ember/controller';
import { computed } from '@ember/object';
import config from '../config/environment';

export default Controller.extend({
  env: config.environment,
  auth: service(),
  store: service(),
  activeCluster: computed('auth.activeCluster', function() {
    return this.get('store').peekRecord('cluster', this.get('auth.activeCluster'));
  }),
  activeClusterName: computed('activeCluster', function() {
    const activeCluster = this.get('activeCluster');
    return activeCluster ? activeCluster.get('name') : null;
  }),
});
