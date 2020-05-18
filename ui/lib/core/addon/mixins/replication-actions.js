import { inject as service } from '@ember/service';
import { or } from '@ember/object/computed';
import { isPresent } from '@ember/utils';
import Mixin from '@ember/object/mixin';
import { task } from 'ember-concurrency';

export default Mixin.create({
  store: service(),
  router: service(),
  loading: or('save.isRunning', 'submitSuccess.isRunning'),
  onEnable() {},
  onDisable() {},
  submitHandler: task(function*(action, clusterMode, data, event) {
    let replicationMode = (data && data.replicationMode) || this.get('replicationMode');
    if (event && event.preventDefault) {
      event.preventDefault();
    }
    this.setProperties({
      errors: [],
    });
    if (data) {
      data = Object.keys(data).reduce((newData, key) => {
        var val = data[key];
        if (isPresent(val)) {
          newData[key] = val;
        }
        return newData;
      }, {});
      delete data.replicationMode;
    }
    return yield this.save.perform(action, replicationMode, clusterMode, data);
  }),

  save: task(function*(action, replicationMode, clusterMode, data) {
    let resp;
    try {
      resp = yield this.get('store')
        .adapterFor('cluster')
        .replicationAction(action, replicationMode, clusterMode, data);
    } catch (e) {
      return this.submitError(e);
    }
    return yield this.submitSuccess.perform(resp, action, clusterMode);
  }).drop(),

  submitSuccess: task(function*(resp, action, mode) {
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
    if (this.reset) {
      // resets mode to primary
      this.reset();
    }
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
    if (mode !== 'secondary') {
      // data is not ready at this point for reload during secondary and causes error
      try {
        yield cluster.reload();
      } catch (e) {
        // no error handling here
      }
    }
    cluster.rollbackAttributes();
    if (action === 'disable') {
      yield this.onDisable();
    }
    if (mode === 'secondary') {
      let modeObject = {};
      // return mode and replicationMode so you can properly handle the transition
      return (modeObject = {
        mode,
        replicationMode,
      });
    }
    // ARG TODO: I would argue onEnable should never be called here, because it uses the transitionTo method which demounts the cluster and makes it so this concurrency function never returns anything.
    // onEnable is a method off the controller/replication-mode.js
    if (action === 'enable') {
      try {
        yield this.onEnable(replicationMode);
      } catch (e) {
        console.log(e);
      }
    }
  }).drop(),

  submitError(e) {
    if (e.errors) {
      this.set('errors', e.errors);
    } else {
      throw e;
    }
  },
});
