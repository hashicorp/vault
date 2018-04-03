import Ember from 'ember';

const DEFAULTS = {
  keyData: null,
  secret_shares: null,
  secret_threshold: null,
  pgp_keys: null,
  use_pgp: false,
  loading: false,
};

export default Ember.Controller.extend(DEFAULTS, {
  reset() {
    this.setProperties(DEFAULTS);
  },

  initSuccess(resp) {
    this.set('loading', false);
    this.set('keyData', resp);
  },

  initError(e) {
    this.set('loading', false);
    if (e.httpStatus === 400) {
      this.set('errors', e.errors);
    } else {
      throw e;
    }
  },

  keyFilename: Ember.computed('model.name', function() {
    return `vault-cluster-${this.get('model.name')}`;
  }),

  actions: {
    initCluster(data) {
      if (data.secret_shares) {
        data.secret_shares = parseInt(data.secret_shares);
      }
      if (data.secret_threshold) {
        data.secret_threshold = parseInt(data.secret_threshold);
      }
      if (!data.use_pgp) {
        delete data.pgp_keys;
      }
      if (!data.use_pgp_for_root) {
        delete data.root_token_pgp_key;
      }

      delete data.use_pgp;
      delete data.use_pgp_for_root;
      const store = this.model.store;
      this.setProperties({
        loading: true,
        errors: null,
      });
      store
        .adapterFor('cluster')
        .initCluster(data)
        .then(resp => this.initSuccess(resp), (...errArgs) => this.initError(...errArgs));
    },

    setKeys(data) {
      this.set('pgp_keys', data);
    },

    setRootKey([key]) {
      this.set('root_token_pgp_key', key);
    },
  },
});
