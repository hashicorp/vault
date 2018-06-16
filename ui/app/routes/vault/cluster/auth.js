import ClusterRouteBase from './cluster-route-base';
import Ember from 'ember';

const { RSVP } = Ember;

export default ClusterRouteBase.extend({
  model() {
    let cluster = this._super(...arguments);
    return this.store
      .findAll('auth-method', {
        adapterOptions: {
          unauthenticated: true,
        },
      })
      .then(result => {
        return RSVP.hash({
          cluster,
          methods: result,
        });
      });
  },
});
