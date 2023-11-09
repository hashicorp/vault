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
import { capitalize } from '@ember/string';
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

  isRootToken: alias('session.data.authenticated.isRootToken'),
  tokenExpirationDate: computed(
    'session.{isAuthenticated,data.authenticated.tokenExpirationEpoch}',
    'expirationCalcTS',
    function () {
      if (!this.session.isAuthenticated) {
        return;
      }
      const { tokenExpirationEpoch } = this.session.data.authenticated;
      const expirationDate = new Date(0);
      return tokenExpirationEpoch ? expirationDate.setUTCMilliseconds(tokenExpirationEpoch) : null;
    }
  ),

  renewAfterEpoch: computed(
    'session.data.authenticated.data.{renewable,ttl}',
    'expirationCalcTS',
    function () {
      const { expirationCalcTS } = this;
      const data = this.session.data.authenticated;
      if (!data?.renewable || !expirationCalcTS) {
        return null;
      }
      const { ttl } = data;
      // renew after last expirationCalc time + half of the ttl (in ms)
      return Math.floor((ttl * 1e3) / 2) + expirationCalcTS;
    }
  ),

  currentToken: alias('session.data.authenticated.token'),

  authData: alias('session.data.authenticated'),

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

  renewCurrentToken() {
    const namespace = this.authData.userRootNamespace;
    const url = '/v1/auth/token/renew-self';
    return this.ajax(url, 'POST', { namespace });
  },

  revokeCurrentToken() {
    const namespace = this.authData.userRootNamespace;
    const url = '/v1/auth/token/revoke-self';
    return this.ajax(url, 'POST', { namespace });
  },

  expirationCalcTS: alias('session.data.authenticated.expirationCalcTS'),

  calculateRootNamespace(currentNamespace, namespace_path, backend) {
    // here we prefer namespace_path if its defined,
    // else we look and see if there's already a namespace saved
    // and then finally we'll use the current query param if the others
    // haven't set a value yet
    // all of the typeof checks are necessary because the root namespace is ''
    let userRootNamespace = namespace_path && namespace_path.replace(/\/$/, '');
    // if we're logging in with token and there's no namespace_path, we can assume
    // that the token belongs to the root namespace
    if (backend === 'token' && !userRootNamespace) {
      userRootNamespace = '';
    }
    if (typeof userRootNamespace === 'undefined') {
      if (this.authData) {
        userRootNamespace = this.authData.userRootNamespace;
      }
    }
    if (typeof userRootNamespace === 'undefined') {
      userRootNamespace = currentNamespace;
    }
    return userRootNamespace;
  },

  renew() {
    const currentlyRenewing = this.isRenewing;
    if (currentlyRenewing) {
      return;
    }
    this.isRenewing = true;
    const { authenticator, backend, token, userRootNamespace } = this.session.data.authenticated;
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
    if (!this.session.isAuthenticated || this.tokenExpired || this.allowExpiration || !renewTime) {
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

  _parseMfaResponse(mfa_requirement) {
    // mfa_requirement response comes back in a shape that is not easy to work with
    // convert to array of objects and add necessary properties to satisfy the view
    if (mfa_requirement) {
      const { mfa_request_id, mfa_constraints } = mfa_requirement;
      const constraints = [];
      for (const key in mfa_constraints) {
        const methods = mfa_constraints[key].any;
        const isMulti = methods.length > 1;

        // friendly label for display in MfaForm
        methods.forEach((m) => {
          const typeFormatted = m.type === 'totp' ? m.type.toUpperCase() : capitalize(m.type);
          m.label = `${typeFormatted} ${m.uses_passcode ? 'passcode' : 'push notification'}`;
        });
        constraints.push({
          name: key,
          methods,
          selectedMethod: isMulti ? null : methods[0],
        });
      }

      return {
        mfa_requirement: { mfa_request_id, mfa_constraints: constraints },
      };
    }
    return {};
  },

  async authenticate(/*{clusterId, backend, data, selectedAuth}*/) {
    const [options] = arguments;
    const adapter = this.clusterAdapter();
    const resp = await adapter.authenticate(options);

    if (resp.auth?.mfa_requirement) {
      return this._parseMfaResponse(resp.auth?.mfa_requirement);
    }

    return this.authSuccess(options, resp.auth || resp.data);
  },

  async totpValidate({ mfa_requirement }) {
    return this.clusterAdapter().mfaValidate(mfa_requirement);
  },

  async authSuccess(options, response) {
    // persist selectedAuth to localStorage to rehydrate auth form on logout
    localStorage.setItem('selectedAuth', options.selectedAuth);
    const authData = await this.persistAuthData(options, response, this.namespaceService.path);
    await this.permissions.getPaths.perform();
    return authData;
  },

  handleError(e) {
    if (e.errors) {
      return e.errors.map((error) => {
        if (error.detail) {
          return error.detail;
        }
        return error;
      });
    }
    return [e];
  },

  getAuthType() {
    // check localStorage first
    const selectedAuth = localStorage.getItem('selectedAuth');
    if (selectedAuth) return selectedAuth;
    // fallback to authData which discerns backend type from token
    return this.authData ? this.authData.backend?.type : null;
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
