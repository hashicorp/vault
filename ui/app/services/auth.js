/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import { task, timeout } from 'ember-concurrency';
import { getOwner } from '@ember/application';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import Service, { inject as service } from '@ember/service';
import fetch from 'fetch';
import { resolve, reject } from 'rsvp';

import ENV from 'vault/config/environment';

const TOKEN_SEPARATOR = '☃';
const TOKEN_PREFIX = 'vault-';
const ROOT_PREFIX = '_root_';

export { TOKEN_SEPARATOR, TOKEN_PREFIX, ROOT_PREFIX };

export default Service.extend({
  permissions: service(),
  currentCluster: service(),
  router: service(),
  session: service(),
  namespaceService: service('namespace'),

  IDLE_TIMEOUT: 3 * 60e3,
  isRenewing: false,
  mfaErrors: null,

  get tokenExpired() {
    const expiration = this.tokenExpirationDate;
    return expiration ? this.now() >= expiration : null;
  },

  activeCluster: alias('currentCluster.cluster'),
  isAuthenticated: alias('session.isAuthenticated'),
  authData: alias('session.data.authenticated'),
  isRootToken: alias('session.data.authenticated.isRootToken'),
  currentToken: alias('session.data.authenticated.token'),
  expirationCalcTS: alias('session.data.authenticated.expirationCalcTS'),

  tokenExpirationDate: computed('isAuthenticated', 'authData', 'expirationCalcTS', function () {
    if (!this.isAuthenticated) {
      return;
    }
    const { tokenExpirationEpoch } = this.authData;
    const expirationDate = new Date(0);
    return tokenExpirationEpoch ? expirationDate.setUTCMilliseconds(tokenExpirationEpoch) : null;
  }),

  renewAfterEpoch: computed('authData.{renewable,ttl}', 'expirationCalcTS', function () {
    const { expirationCalcTS } = this;
    if (!this.authData?.renewable || !expirationCalcTS) {
      return null;
    }
    const { ttl } = this.authData;
    // renew after last expirationCalc time + half of the ttl (in ms)
    return Math.floor((ttl * 1e3) / 2) + expirationCalcTS;
  }),

  clusterAdapter() {
    return getOwner(this).lookup('adapter:cluster');
  },

  environment() {
    return ENV.environment;
  },

  now() {
    return Date.now();
  },

  setCluster(clusterId) {
    this.set('activeClusterId', clusterId);
  },

  ajax(url, method, options) {
    const defaults = {
      url,
      method,
      dataType: 'json',
      headers: {
        'X-Vault-Token': this.currentToken,
      },
    };

    const namespace =
      typeof options.namespace === 'undefined' ? this.namespaceService.path : options.namespace;
    if (namespace) {
      defaults.headers['X-Vault-Namespace'] = namespace;
    }
    const opts = Object.assign(defaults, options);

    return fetch(url, {
      method: opts.method || 'GET',
      headers: opts.headers || {},
    }).then((response) => {
      if (response.status === 204) {
        return resolve();
      } else if (response.status >= 200 && response.status < 300) {
        return resolve(response.json());
      } else {
        return reject(response);
      }
    });
  },

  renew() {
    const currentlyRenewing = this.isRenewing;
    if (currentlyRenewing) {
      return;
    }
    this.isRenewing = true;
    const { authenticator, backend, token, userRootNamespace } = this.authData;
    return this.session
      .authenticate(
        authenticator,
        { token },
        { renew: true, backend: backend.mountPath, namespace: userRootNamespace }
      )
      .finally(() => {
        this.isRenewing = false;
      });
  },

  checkShouldRenew: task(function* () {
    while (true) {
      if (Ember.testing) {
        return;
      }
      yield timeout(5000);
      if (this.shouldRenew()) {
        yield this.renew();
      }
    }
  }).on('init'),

  shouldRenew() {
    const now = this.now();
    const lastFetch = this.lastFetch;
    const renewTime = this.renewAfterEpoch;
    if (!this.isAuthenticated || this.tokenExpired || this.allowExpiration || !renewTime) {
      return false;
    }
    if (lastFetch && now - lastFetch >= this.IDLE_TIMEOUT) {
      this.set('allowExpiration', true);
      return false;
    }
    if (now >= renewTime) {
      return true;
    }
    return false;
  },

  setLastFetch(timestamp) {
    const now = this.now();
    this.set('lastFetch', timestamp);
    // if expiration was allowed and we're over half the ttl we want to go ahead and renew here
    if (this.allowExpiration && now >= this.renewAfterEpoch) {
      this.renew();
    }
    this.set('allowExpiration', false);
  },

  async totpValidate({ mfa_requirement }) {
    return this.clusterAdapter().mfaValidate(mfa_requirement);
  },

  getOktaNumberChallengeAnswer(nonce, mount) {
    const url = `/v1/auth/${mount}/verify/${nonce}`;
    return this.ajax(url, 'GET', {}).then(
      (resp) => {
        return resp.data.correct_answer;
      },
      (e) => {
        // if error status is 404, return and keep polling for a response
        if (e.status === 404) {
          return null;
        } else {
          throw e;
        }
      }
    );
  },
});
