import { isPresent } from '@ember/utils';
import { alias } from '@ember/object/computed';
import { inject as service } from '@ember/service';
import Controller from '@ember/controller';

const DEFAULTS = {
  token: null,
  id: null,
  loading: false,
  errors: [],
  showFilterConfig: false,
  primary_api_addr: null,
  primary_cluster_addr: null,
  filterConfig: {
    mode: 'whitelist',
    paths: [],
  },
};

export default Controller.extend(DEFAULTS, {
  store: service(),
  rm: service('replication-mode'),
  replicationMode: alias('rm.mode'),

  submitError(e) {
    if (e.errors) {
      this.set('errors', e.errors);
    } else {
      throw e;
    }
  },

  saveFilterConfig() {
    const config = this.get('filterConfig');
    const id = this.get('id');
    config.id = id;
    const configRecord = this.get('store').createRecord('mount-filter-config', config);
    return configRecord.save().catch(e => this.submitError(e));
  },

  reset() {
    this.setProperties(DEFAULTS);
  },

  submitSuccess(resp, action) {
    const cluster = this.get('model');
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
    return store
      .adapterFor('cluster')
      .replicationStatus()
      .then(status => {
        return store.pushPayload('cluster', status);
      })
      .finally(() => {
        this.set('loading', false);
      });
  },

  submitHandler(action, clusterMode, data, event) {
    const replicationMode = this.get('replicationMode');
    let saveFilterConfig;
    if (event && event.preventDefault) {
      event.preventDefault();
    }
    if (data && isPresent(data.saveFilterConfig)) {
      saveFilterConfig = data.saveFilterConfig;
      delete data.saveFilterConfig;
    }
    this.setProperties({
      loading: true,
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
    }

    return this.get('store')
      .adapterFor('cluster')
      .replicationAction(action, replicationMode, clusterMode, data)
      .then(
        resp => {
          if (saveFilterConfig) {
            return this.saveFilterConfig().then(() => {
              return this.submitSuccess(resp, action, clusterMode);
            });
          } else {
            return this.submitSuccess(resp, action, clusterMode);
          }
        },
        (...args) => this.submitError(...args)
      );
  },

  actions: {
    onSubmit(/*action, mode, data, event*/) {
      return this.submitHandler(...arguments);
    },

    clear() {
      this.reset();
      this.setProperties({
        token: null,
        id: null,
      });
    },
  },
});
