import Ember from 'ember';
import ClusterRoute from './cluster-route-base';

const { inject } = Ember;

export default ClusterRoute.extend({
  wizard: inject.service(),

  activate() {
    this.get('wizard').set('currentState', 'idle');
  },
});
