import Ember from 'ember';
import ClusterRoute from 'vault/mixins/cluster-route';
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';

const POLL_INTERVAL_MS = 10000;
const { inject } = Ember;

export default Ember.Route.extend(ModelBoundaryRoute, ClusterRoute, {
  version: inject.service(),
  store: inject.service(),
  auth: inject.service(),
  currentCluster: Ember.inject.service(),
  modelTypes: ['node', 'secret', 'secret-engine'],

  getClusterId(params) {
    const { cluster_name } = params;
    const cluster = this.modelFor('vault').findBy('name', cluster_name);
    return cluster ? cluster.get('id') : null;
  },

  beforeModel() {
    const params = this.paramsFor(this.routeName);
    const id = this.getClusterId(params);
    if (id) {
      this.get('auth').setCluster(id);
      return this.get('version').fetchFeatures();
    } else {
      return Ember.RSVP.reject({ httpStatus: 404, message: 'not found', path: params.cluster_name });
    }
  },

  model(params) {
    const id = this.getClusterId(params);

    return this.get('store').findRecord('cluster', id);
  },

  stopPoll: Ember.on('deactivate', function() {
    Ember.run.cancel(this.get('timer'));
  }),

  poll() {
    // when testing, the polling loop causes promises to never settle so acceptance tests hang
    // to get around that, we just disable the poll in tests
    return Ember.testing
      ? null
      : Ember.run.later(() => {
          this.controller.get('model').reload().then(
            () => {
              this.set('timer', this.poll());
              return this.transitionToTargetRoute();
            },
            () => {
              this.set('timer', this.poll());
            }
          );
        }, POLL_INTERVAL_MS);
  },

  afterModel(model) {
    this.get('currentCluster').setCluster(model);
    this._super(...arguments);
    this.poll();
    return this.transitionToTargetRoute();
  },

  actions: {
    error(e) {
      if (e.httpStatus === 503 && e.errors[0] === 'Vault is sealed') {
        this.refresh();
      }
      return true;
    },
  },
});
