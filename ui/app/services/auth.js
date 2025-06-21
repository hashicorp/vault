/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import { task, timeout } from 'ember-concurrency';
import { getOwner } from '@ember/owner';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import Service, { inject as service } from '@ember/service';
import { capitalize } from '@ember/string';
import { resolve, reject } from 'rsvp';

import getStorage from 'vault/lib/token-storage';
import ENV from 'vault/config/environment';
import { addToArray } from 'vault/helpers/add-to-array';

const TOKEN_SEPARATOR = '☃';
const TOKEN_PREFIX = 'vault-';
const ROOT_PREFIX = '_root_';

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
    // const { tokenExpirationEpoch } = this.getTokenData(tokenName);
    const tokenData = this.getTokenData(tokenName);
    const tokenExpirationEpoch = tokenData ? tokenData?.tokenExpirationEpoch : undefined;
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
    const stored = this.getTokenData(token);
    return Object.assign(stored);
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

  // ttl is either the "ttl" or "lease_duration" returned by the auth data response
  calculateExpiration(expirationData, now) {
    // calculateExpiration(ttl, expire_time, now) {
    const ttl = expirationData.ttl || expirationData.lease_duration;
    const tokenExpirationEpoch = expirationData.expire_time
      ? new Date(expirationData.expire_time).getTime()
      : now + ttl * 1e3;

    return { ttl, tokenExpirationEpoch };
  },

  // setExpirationSettings(resp, now) {
  //   if (resp.renewable) {
  //     this.set('expirationCalcTS', now);
  //     this.set('allowExpiration', false);
  //   } else {
  //     this.set('allowExpiration', true);
  //   }
  // },

  setExpirationSettings(renewable, now) {
    if (renewable) {
      this.set('expirationCalcTS', now);
      this.set('allowExpiration', false);
    } else {
      this.set('allowExpiration', true);
    }
  },

  calculateRootNamespace(currentNamespace, namespacePath, backend) {
    // here we prefer namespace_path if its defined,
    // else we look and see if there's already a namespace saved
    // and then finally we'll use the current query param if the others
    // haven't set a value yet
    // all of the typeof checks are necessary because the root namespace is ''
    let userRootNamespace = namespacePath && namespacePath.replace(/\/$/, '');
    // renew-self does not return namespace_path, so we manually setting in renew().
    // so if we're logging in with token and there's no namespace_path, we can assume
    // that the token belongs to the root namespace
    if (backend === 'token' && !userRootNamespace) {
      userRootNamespace = '';
    }
    if (typeof userRootNamespace === 'undefined' && this.authData) {
      userRootNamespace = this.authData.userRootNamespace;
    }
    if (typeof userRootNamespace === 'undefined') {
      userRootNamespace = currentNamespace;
    }
    return userRootNamespace;
  },

  // TODO CMB changes below are stopgaps until this method is un-abstracted
  // end goal is for each auth method's component to handle setting relevant parameters.
  // this method should just accept an arg of data to persist from the response as well as:
  // 1. generate token name and set token data
  // 2. calculate and set expiration
  // 3. (maybe) calculate root namespace
  // async persistAuthData() {
  //   const [firstArg, resp] = arguments;
  //   const currentNamespace = this.namespaceService.path || '';
  //   // dropdown vs tab format
  //   //
  //   // TODO adding ANOTHER conditional until this method is un-abstracted :(
  //   const mountPath = firstArg?.path || firstArg?.data?.path || firstArg?.selectedAuth;
  //   let tokenName;
  //   let options;
  //   let backend;

  //   // TODO move setting current backend, options, etc to method's component
  //   if (typeof firstArg === 'string') {
  //     tokenName = firstArg;
  //     backend = this.backendFromTokenName(tokenName);
  //   } else {
  //     options = firstArg;
  //     // backend is old news since it's confusing whether it refers to the auth mount path or auth type,
  //     // new auth flow explicitly defines "selectedAuth" and "path"
  //     backend = options?.backend || options.selectedAuth;
  //   }

  //   const currentBackend = {
  //     mountPath,
  //     ...BACKENDS.find((b) => b.type === backend),
  //   };

  //   const { entity_id, policies, renewable, namespace_path } = resp;
  //   const userRootNamespace = this.calculateRootNamespace(currentNamespace, namespace_path, backend);
  //   const data = {
  //     userRootNamespace,
  //     displayName: null, // set below
  //     token: resp.client_token || get(resp, currentBackend.tokenPath),
  //     policies,
  //     renewable,
  //     entityId: entity_id,
  //     authMethodType: options?.selectedAuth,
  //     authMountPath: mountPath,
  //   };

  //   tokenName = this.generateTokenName(
  //     {
  //       backend,
  //       clusterId: (options && options.clusterId) || this.activeClusterId,
  //     },
  //     resp.policies
  //   );

  //   const now = this.now();

  //   Object.assign(data, this.calculateExpiration(resp, now));
  //   this.setExpirationSettings(resp, now);

  //   // ensure we don't call renew-self within tests
  //   // this is intentionally not included in setExpirationSettings so we can unit test that method
  //   if (Ember.testing) this.set('allowExpiration', false);

  //   data.displayName = await this.setDisplayName(resp, currentBackend.displayNamePath, tokenName);

  //   this.set('tokens', addToArray(this.tokens, tokenName));
  //   this.setTokenData(tokenName, data);

  //   return resolve({
  //     namespace: currentNamespace || data.userRootNamespace,
  //     token: tokenName,
  //     isRoot: policies.includes('root'),
  //   });
  // },

  async persistAuthData(clusterId, authResponseData) {
    const currentNamespace = this.namespaceService.path || '';
    // displayName will be set by auth method (not here by auth service)
    // namespace_path only returned for methods that use a token exchange to authenticate (i.e. token, oidc)
    const { authMethodType, expireTime, namespacePath, policies, renewable, ttl } = authResponseData;

    const data = {
      userRootNamespace: this.calculateRootNamespace(currentNamespace, namespacePath, authMethodType),
      ...authResponseData,
    };

    const tokenName = this.generateTokenName(
      { backend: authMethodType, clusterId: clusterId || this.activeClusterId },
      policies
    );

    const now = this.now();
    const expirationData = { ttl, expire_time: expireTime };
    Object.assign(data, this.calculateExpiration(expirationData, now));
    this.setExpirationSettings(renewable, now);

    if (!data.displayName) {
      const lookupDisplayName = await this.setDisplayName(data.token);
      data.displayName = lookupDisplayName || authMethodType;
    }

    // ensure we don't call renew-self within tests
    // this is intentionally not included in setExpirationSettings so we can unit test that method
    if (Ember.testing) this.set('allowExpiration', false);

    this.set('tokens', addToArray(this.tokens, tokenName));
    this.setTokenData(tokenName, data);

    return resolve({
      namespace: currentNamespace || data.userRootNamespace,
      token: tokenName,
      isRoot: (policies || []).includes('root'),
    });
  },

  async setDisplayName(clientToken) {
    // this is a fallback for any methods that don't return a display name from the initial auth request (i.e. JWT)
    // or for OIDC/SAML with mfa configured because the mfa/validate endpoint does not consistently
    // return display_name (or metadata that includes something to be used as such).
    // this if block can be removed if/when the API consistently returns a display_name.
    // if still nothing, request token data as a last resort
    try {
      const { data } = await this.lookupSelf(clientToken);
      return data.display_name;
    } catch {
      // silently fail since we're just trying to set a display name
    }
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
        const namespacePath = this.namespaceService.path;
        const response = resp.data || resp.auth;
        // renew-self does not return namespace_path, so manually add it if it exists
        if (!response?.namespace_path && namespacePath) {
          response.namespace_path = namespacePath;
        }
        return this.persistAuthData(tokenName, response);
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

  _parseMfaResponse(mfaRequirement) {
    // mfaRequirement response comes back in a shape that is not easy to work with
    // convert to array of objects and add necessary properties to satisfy the view
    if (mfaRequirement) {
      const { mfaRequestId, mfaConstraints } = mfaRequirement;
      const constraints = [];
      for (const key in mfaConstraints) {
        const methods = mfaConstraints[key].any;
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
      return { mfaRequestId, mfaConstraints: constraints };
    }
    return {};
  },

  async authenticate(/*{clusterId, backend, data, selectedAuth}*/) {
    const [options] = arguments;
    const adapter = this.clusterAdapter();
    const resp = await adapter.authenticate(options);

    if (resp.auth?.mfa_requirement) {
      const mfaRequirement = resp.auth?.mfa_requirement;
      return this._parseMfaResponse(mfaRequirement);
    }

    return this.authSuccess(options, resp.auth || resp.data);
  },

  // async totpValidate({ mfaRequirement, ...options }) {
  //   const resp = await this.clusterAdapter().mfaValidate(mfaRequirement);
  //   return this.authSuccess(options, resp.auth || resp.data);
  // },

  async totpValidate({ clusterId, mfaRequirement, authMethodType, authMountPath }) {
    const resp = await this.clusterAdapter().mfaValidate(mfaRequirement);
    const data = resp?.data || resp?.auth;
    const normalizedAuthData = {
      authMethodType,
      authMountPath,
      entityId: data.entity_id,
      token: data?.client_token || data?.id,
      renewable: data.renewable,
      ttl: data?.lease_duration || data?.ttl,
      policies: data.policies,
      displayName: data?.display_name || data?.metadata?.username || data?.metadata?.org,
    };
    return this.authSuccess(clusterId, normalizedAuthData);
  },

  async authSuccess(clusterId, authResponse) {
    // persist selectedAuth to localStorage to rehydrate auth form on logout
    localStorage.setItem('selectedAuth', authResponse.authMethodType);
    const authData = await this.persistAuthData(clusterId, authResponse);
    // TODO why does this line ruin everything???
    // await this.permissions.getPaths.perform();
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
    return this.authData ? this.authData.authMethodType : null;
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
});
