/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { computed } from '@ember/object';
import Controller from '@ember/controller';

const DEFAULTS = {
  keyData: null,
  secret_shares: null,
  secret_threshold: null,
  pgp_keys: null,
  use_pgp: false,
  loading: false,
};

export default Controller.extend(DEFAULTS, {
  reset() {
    this.setProperties(DEFAULTS);
  },

  initSuccess(resp) {
    this.set('loading', false);
    this.set('keyData', resp);
    this.model.reload();
  },

  initError(e) {
    this.set('loading', false);
    if (e.httpStatus === 400) {
      this.set('errors', e.errors);
    } else {
      throw e;
    }
  },

  keyFilename: computed('model.name', function () {
    return `vault-cluster-${this.model.name}`;
  }),

  actions: {
    initCluster(payload) {
      const data = { ...payload };
      const isCloudSeal = !!this.model.sealType && this.model.sealType !== 'shamir';
      if (data.secret_shares) {
        const shares = parseInt(data.secret_shares, 10);
        data.secret_shares = shares;
        if (isCloudSeal) {
          data.stored_shares = 1;
          data.recovery_shares = shares;
          delete data.secret_shares; // API will throw an error if secret_shares is passed for seal types other than shamir (transit, AWSKMS etc.)
        }
      }
      if (data.secret_threshold) {
        const threshold = parseInt(data.secret_threshold, 10);
        data.secret_threshold = threshold;
        if (isCloudSeal) {
          data.recovery_threshold = threshold;
          delete data.secret_threshold; // API will throw an error if secret_threshold is passed for seal types other than shamir (transit, AWSKMS etc.)
        }
      }
      if (!data.use_pgp) {
        delete data.pgp_keys;
      }
      if (data.use_pgp && isCloudSeal) {
        data.recovery_pgp_keys = data.pgp_keys;
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
        .then(
          (resp) => this.initSuccess(resp),
          (...errArgs) => this.initError(...errArgs)
        );
    },

    setKeys(data) {
      this.set('pgp_keys', data);
    },

    setRootKey([key]) {
      this.set('root_token_pgp_key', key);
    },
  },
});
