import Ember from 'ember';
import ClusterRoute from 'vault/mixins/cluster-route';
import ModelBoundaryRoute from 'vault/mixins/model-boundary-route';

const POLL_INTERVAL_MS = 10000;
const { inject, Route, getOwner } = Ember;

export default Route.extend(ModelBoundaryRoute, ClusterRoute, {
  namespaceService: inject.service('namespace'),
  version: inject.service(),
  store: inject.service(),
  auth: inject.service(),
  currentCluster: inject.service(),
  modelTypes: ['node', 'secret', 'secret-engine'],
  globalNamespaceModels: ['node', 'cluster'],

  queryParams: {
    namespaceQueryParam: {
      refreshModel: true,
    },
  },

  getClusterId(params) {
    const { cluster_name } = params;
    const cluster = this.modelFor('vault').findBy('name', cluster_name);
    return cluster ? cluster.get('id') : null;
  },

  clearNonGlobalModels() {
    // this method clears all of the ember data cached models except
    // the model types blacklisted in `globalNamespaceModels`
    let store = this.store;
    let modelsToKeep = this.get('globalNamespaceModels');
    for (let model of getOwner(this).lookup('data-adapter:main').getModelTypes()) {
      let { name } = model;
      if (modelsToKeep.includes(name)) {
        return;
      }
      store.unloadAll(name);
    }
  },

  beforeModel() {
    const params = this.paramsFor(this.routeName);
    this.clearNonGlobalModels();
    this.get('namespaceService').setNamespace(params.namespaceQueryParam);
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

    // Check that namespaces is enabled and if not,
    // clear the namespace by transition to this route w/o it
    if (this.get('namespaceService.path') && !this.get('version.hasNamespaces')) {
      return this.transitionTo(this.routeName, { queryParams: { namespace: '' } });
    }
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
