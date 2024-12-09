/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import { task, timeout } from 'ember-concurrency';
import { getOwner } from '@ember/owner';
import { isArray } from '@ember/array';
import { computed, get } from '@ember/object';
import { alias } from '@ember/object/computed';
import Service, { inject as service } from '@ember/service';
import { capitalize } from '@ember/string';
import { resolve, reject } from 'rsvp';

import getStorage from 'vault/lib/token-storage';
import ENV from 'vault/config/environment';
import { allSupportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import { addToArray } from 'vault/helpers/add-to-array';

const TOKEN_SEPARATOR = '☃';
const TOKEN_PREFIX = 'vault-';
const ROOT_PREFIX = '_root_';
const BACKENDS = allSupportedAuthBackends();

export { TOKEN_SEPARATOR, TOKEN_PREFIX, ROOT_PREFIX };

export default Service.extend({
  permissions: service(),
  currentCluster: service(),
  router: service(),
  store: service(),
  namespaceService: service('namespace'),

  IDLE_TIMEOUT: 3 * 60e3,
  expirationCalcTS: null,
  isRenewing: false,
  mfaErrors: null,
  isRootToken: false,

  get tokenExpired() {
    const expiration = this.tokenExpirationDate;
    return expiration ? this.now() >= expiration : null;
  },

  activeCluster: alias('currentCluster.cluster'),

  // eslint-disable-next-line
  tokens: computed({
    get() {
      return this._tokens || this.getTokensFromStorage() || [];
    },
    set(key, value) {
      return (this._tokens = value);
    },
  }),

  isActiveSession: computed(
    'router.currentRouteName',
    'currentToken',
    'activeCluster.{dr.isSecondary,needsInit,sealed,name}',
    function () {
      if (this.activeCluster) {
        if (this.activeCluster.dr?.isSecondary || this.activeCluster.needsInit || this.activeCluster.sealed) {
          return false;
        }
        if (
          this.activeCluster.name &&
          this.currentToken &&
          this.router.currentRouteName !== 'vault.cluster.auth'
        ) {
          return true;
        }
      }
      return false;
    }
  ),

  tokenExpirationDate: computed('currentTokenName', 'expirationCalcTS', function () {
    const tokenName = this.currentTokenName;
    if (!tokenName) {
      return;
    }

    const { tokenExpirationEpoch } = this.getTokenData(tokenName);
    const expirationDate = new Date(0);

    return tokenExpirationEpoch ? expirationDate.setUTCMilliseconds(tokenExpirationEpoch) : null;
  }),

  renewAfterEpoch: computed('currentTokenName', 'expirationCalcTS', function () {
    const tokenName = this.currentTokenName;
    const { expirationCalcTS } = this;
    const data = this.getTokenData(tokenName);
    if (!tokenName || !data || !expirationCalcTS) {
      return null;
    }
    const { ttl, renewable } = data;
    // renew after last expirationCalc time + half of the ttl (in ms)
    return renewable ? Math.floor((ttl * 1e3) / 2) + expirationCalcTS : null;
  }),

  // returns the key for the token to use
  currentTokenName: computed('activeClusterId', 'tokens', 'tokens.[]', function () {
    const regex = new RegExp(this.activeClusterId);
    return this.tokens.find((key) => regex.test(key));
  }),

  currentToken: computed('currentTokenName', function () {
    const name = this.currentTokenName;
    const data = name && this.getTokenData(name);
    // data.token is undefined so that's why it returns current token undefined
    return name && data ? data.token : null;
  }),

  authData: computed('currentTokenName', function () {
    const token = this.currentTokenName;
    if (!token) {
      return;
    }
    const backend = this.backendFromTokenName(token);
    const stored = this.getTokenData(token);
    return Object.assign(stored, {
      backend: {
        // add mount path for password reset
        mountPath: stored.backend.mountPath,
        ...BACKENDS.find((b) => b.type === backend),
      },
    });
  }),

  init() {
    this._super(...arguments);
    this.checkForRootToken();
  },

  clusterAdapter() {
    return getOwner(this).lookup('adapter:cluster');
  },

  generateTokenName({ backend, clusterId }, policies) {
    return (policies || []).includes('root')
      ? `${TOKEN_PREFIX}${ROOT_PREFIX}${TOKEN_SEPARATOR}${clusterId}`
      : `${TOKEN_PREFIX}${backend}${TOKEN_SEPARATOR}${clusterId}`;
  },

  backendFromTokenName(tokenName) {
    return tokenName.includes(`${TOKEN_PREFIX}${ROOT_PREFIX}`)
      ? 'token'
      : tokenName.slice(TOKEN_PREFIX.length).split(TOKEN_SEPARATOR)[0];
  },

  storage(tokenName) {
    if (
      tokenName &&
      tokenName.indexOf(`${TOKEN_PREFIX}${ROOT_PREFIX}`) === 0 &&
      this.environment() !== 'development'
    ) {
      return getStorage('memory');
    } else {
      return getStorage();
    }
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

  async lookupSelf(token) {
    return this.store
      .adapterFor('application')
      .ajax('/v1/auth/token/lookup-self', 'GET', { headers: { 'X-Vault-Token': token } });
  },

  revokeCurrentToken() {
    const namespace = this.authData.userRootNamespace;
    const url = '/v1/auth/token/revoke-self';
    return this.ajax(url, 'POST', { namespace });
  },

  calculateExpiration(resp, now) {
    const ttl = resp.ttl || resp.lease_duration;
    const tokenExpirationEpoch = resp.expire_time ? new Date(resp.expire_time).getTime() : now + ttl * 1e3;

    return { ttl, tokenExpirationEpoch };
  },

  setExpirationSettings(resp, now) {
    if (resp.renewable) {
      this.set('expirationCalcTS', now);
      this.set('allowExpiration', false);
    } else {
      this.set('allowExpiration', true);
    }
  },

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

  async persistAuthData() {
    const [firstArg, resp] = arguments;
    const currentNamespace = this.namespaceService.path || '';
    // dropdown vs tab format
    const mountPath = firstArg?.data?.path || firstArg?.selectedAuth;
    let tokenName;
    let options;
    let backend;
    if (typeof firstArg === 'string') {
      tokenName = firstArg;
      backend = this.backendFromTokenName(tokenName);
    } else {
      options = firstArg;
      backend = options.backend;
    }

    const currentBackend = {
      mountPath,
      ...BACKENDS.find((b) => b.type === backend),
    };

    const { entity_id, policies, renewable, namespace_path } = resp;
    const userRootNamespace = this.calculateRootNamespace(currentNamespace, namespace_path, backend);
    const data = {
      userRootNamespace,
      displayName: null, // set below
      backend: currentBackend,
      token: resp.client_token || get(resp, currentBackend.tokenPath),
      policies,
      renewable,
      entity_id,
    };

    tokenName = this.generateTokenName(
      {
        backend,
        clusterId: (options && options.clusterId) || this.activeClusterId,
      },
      resp.policies
    );

    const now = this.now();

    Object.assign(data, this.calculateExpiration(resp, now));
    this.setExpirationSettings(resp, now);

    // ensure we don't call renew-self within tests
    // this is intentionally not included in setExpirationSettings so we can unit test that method
    if (Ember.testing) this.set('allowExpiration', false);

    data.displayName = await this.setDisplayName(resp, currentBackend.displayNamePath, tokenName);

    this.set('tokens', addToArray(this.tokens, tokenName));
    this.setTokenData(tokenName, data);

    return resolve({
      namespace: currentNamespace || data.userRootNamespace,
      token: tokenName,
      isRoot: policies.includes('root'),
    });
  },

  async setDisplayName(resp, displayNamePath, tokenName) {
    let displayName;

    // first check if auth response includes a display name
    displayName = isArray(displayNamePath)
      ? displayNamePath.map((name) => get(resp, name)).join('/')
      : get(resp, displayNamePath);

    // if not, check stored token data
    if (!displayName) {
      displayName = (this.getTokenData(tokenName) || {}).displayName;
    }

    // this is a workaround for OIDC/SAML methods WITH mfa configured. at this time mfa/validate endpoint does not
    // return display_name (or metadata that includes it) for this auth combination.
    // this if block can be removed if/when the API returns display_name on the mfa/validate response.
    if (!displayName) {
      // if still nothing, request token data as a last resort
      try {
        const { data } = await this.lookupSelf(resp.client_token);
        displayName = data.display_name;
      } catch {
        // silently fail since we're just trying to set a display name
      }
    }
    return displayName;
  },

  setTokenData(token, data) {
    this.storage(token).setItem(token, data);
  },

  getTokenData(token) {
    return this.storage(token).getItem(token);
  },

  removeTokenData(token) {
    return this.storage(token).removeItem(token);
  },

  renew() {
    const tokenName = this.currentTokenName;
    const currentlyRenewing = this.isRenewing;

    if (currentlyRenewing) return;

    this.isRenewing = true;
    return this.renewCurrentToken().then(
      (resp) => {
        this.isRenewing = false;
        return this.persistAuthData(tokenName, resp.data || resp.auth);
      },
      (e) => {
        this.isRenewing = false;
        throw e;
      }
    );
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
    // renewAfterEpoch is a unix timestamp of login time + half of ttl
    const renewTime = this.renewAfterEpoch;
    if (!this.currentTokenName || this.tokenExpired || this.allowExpiration || !renewTime) {
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
    if (this.allowExpiration && this.renewAfterEpoch && now >= this.renewAfterEpoch) {
      this.renew();
    }
    this.set('allowExpiration', false);
  },

  getTokensFromStorage(filterFn) {
    return this.storage()
      .keys()
      .reject((key) => {
        return key.indexOf(TOKEN_PREFIX) !== 0 || (filterFn && filterFn(key));
      });
  },

  checkForRootToken() {
    if (this.environment() === 'development') {
      return;
    }

    this.getTokensFromStorage().forEach((key) => {
      const data = this.getTokenData(key);
      if (data && data.policies && data.policies.includes('root')) {
        this.removeTokenData(key);
      }
    });
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

  async totpValidate({ mfa_requirement, ...options }) {
    const resp = await this.clusterAdapter().mfaValidate(mfa_requirement);
    return this.authSuccess(options, resp.auth || resp.data);
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
    return this.authData ? this.authData.backend.type : null;
  },

  deleteCurrentToken() {
    const tokenName = this.currentTokenName;
    this.deleteToken(tokenName);
    this.removeTokenData(tokenName);
  },

  deleteToken(tokenName) {
    const tokenNames = this.tokens.without(tokenName);
    this.removeTokenData(tokenName);
    this.set('tokens', tokenNames);
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
