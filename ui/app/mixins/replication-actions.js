import Ember from 'ember';
const { inject, computed } = Ember;

export default Ember.Mixin.create({
  store: inject.service(),
  routing: inject.service('-routing'),
  router: computed.alias('routing.router'),
  submitHandler(action, clusterMode, data, event) {
    let replicationMode = (data && data.replicationMode) || this.get('replicationMode');
    if (event && event.preventDefault) {
      event.preventDefault();
    }
    this.setProperties({
      loading: true,
      errors: [],
    });
    if (data) {
      data = Object.keys(data).reduce((newData, key) => {
        var val = data[key];
        if (Ember.isPresent(val)) {
          newData[key] = val;
        }
        return newData;
      }, {});
      delete data.replicationMode;
    }

    return this.get('store')
      .adapterFor('cluster')
      .replicationAction(action, replicationMode, clusterMode, data)
      .then(
        resp => {
          return this.submitSuccess(resp, action, clusterMode);
        },
        (...args) => this.submitError(...args)
      );
  },

  submitSuccess(resp, action, mode) {
    const cluster = this.get('cluster');
    const replicationMode = this.get('selectedReplicationMode') || this.get('replicationMode');
    const store = this.get('store');
    if (!cluster) {
      return;
    }

    if (resp && resp.wrap_info) {
      this.set('token', resp.wrap_info.token);
    }
    if (action === 'secondary-token') {
      this.setProperties({
        loading: false,
        primary_api_addr: null,
        primary_cluster_addr: null,
      });
      return cluster;
    }
    this.reset();
    if (action === 'enable') {
      // do something to show model is pending
      cluster.set(
        replicationMode,
        store.createFragment('replication-attributes', {
          mode: 'bootstrapping',
        })
      );
      if (mode === 'secondary' && replicationMode === 'performance') {
        // if we're enabing a secondary, there could be mount filtering,
        // so we should unload all of the backends
        store.unloadAll('secret-engine');
      }
    }
    const router = this.get('router');
    if (action === 'disable') {
      return router.transitionTo.call(router, 'vault.cluster.replication.mode', replicationMode);
    }
    return cluster
      .reload()
      .then(() => {
        cluster.rollbackAttributes();
        if (action === 'enable') {
          return router.transitionTo.call(router, 'vault.cluster.replication.mode', replicationMode);
        }

        if (mode === 'secondary' && replicationMode === 'dr') {
          return router.transitionTo.call(router, 'vault.cluster');
        }
      })
      .finally(() => {
        this.set('loading', false);
      });
  },

  submitError(e) {
    if (e.errors) {
      this.set('errors', e.errors);
    } else {
      throw e;
    }
  },
});
