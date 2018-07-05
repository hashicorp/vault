import ClusterRouteBase from './cluster-route-base';
import Ember from 'ember';

const { RSVP } = Ember;

export default ClusterRouteBase.extend({
  beforeModel() {
    return this.store.unloadAll('auth-method');
  },
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
  resetController(controller) {
    controller.set('wrappedToken', '');
    controller.set('authMethod', '');
  },
});
