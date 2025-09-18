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
    const tokenData = this.getTokenData(tokenName);
    const tokenExpirationEpoch = tokenData ? tokenData?.tokenExpirationEpoch : undefined;
    const expirationDate = new Date(0); // Creates a "zeroed" date object

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

  // ttl is originally either the "ttl" or "lease_duration" returned by the auth data response
  calculateExpiration({ now, ttl, expireTime }) {
    // First check if the ttl is falsy, including 0, before converting to milliseconds.
    // Obviously a ttl of zero seconds is not recommended, but root tokens have a `0` ttl because they never expire.
    // Note - this is different from mount configurations where a `ttl: 0` actually means the value is "unset" and to use system defaults.
    const convertToMilliseconds = () => (ttl ? now + ttl * 1e3 : null);
    const tokenExpirationEpoch = expireTime ? new Date(expireTime).getTime() : convertToMilliseconds();
    // To avoid confusion, if a TTL is `0` return null
    return { ttl: ttl || null, tokenExpirationEpoch };
  },

  setExpirationSettings(renewable, now) {
    if (renewable) {
      this.set('expirationCalcTS', now);
      this.set('allowExpiration', false);
    } else {
      this.set('allowExpiration', true);
    }
  },

  calculateRootNamespace(currentNamespace, namespacePath, backend) {
    // namespace_path is only returned for methods that use a token exchange to authenticate (i.e. token, oidc)
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

  async persistAuthData(clusterId, authResponseData) {
    // An empty string denotes the "root" namespace
    const currentNamespace = this.namespaceService.path || '';
    // Only pull out the necessary data
    const { authMethodType, authMountPath, entityId, policies, renewable, token, ttl } = authResponseData;

    // Lookup token for additional data that may be missing from the method's login response
    const { displayName, expireTime, namespacePath } = await this.lookupTokenData(token, !!currentNamespace, {
      displayName: authResponseData?.displayName,
      expireTime: authResponseData?.expireTime,
      namespacePath: authResponseData?.namespacePath,
    });

    const userRootNamespace = this.calculateRootNamespace(currentNamespace, namespacePath, authMethodType);

    const persistedTokenData = {
      authMethodType,
      authMountPath,
      displayName: displayName || authMethodType,
      entityId,
      policies,
      renewable,
      token,
      userRootNamespace,
      // Only include namespacePath if it exists
      ...(namespacePath && { namespacePath }),
    };

    // Set stored ttl and tokenExpirationEpoch
    const now = this.now();
    const { ttl: calculatedTtl, tokenExpirationEpoch } = this.calculateExpiration({ now, ttl, expireTime });
    persistedTokenData.ttl = calculatedTtl;
    persistedTokenData.tokenExpirationEpoch = tokenExpirationEpoch;
    this.setExpirationSettings(renewable, now);
    // ensure we don't call renew-self within tests
    // this is intentionally not included in setExpirationSettings so we can unit test that method
    if (Ember.testing) this.set('allowExpiration', false);

    // Set token name and store data
    const tokenName = this.generateTokenName({ backend: authMethodType, clusterId }, policies);
    this.set('tokens', addToArray(this.tokens, tokenName));
    this.setTokenData(tokenName, persistedTokenData);
    return resolve({
      namespace: currentNamespace || persistedTokenData.userRootNamespace,
      token: tokenName,
      isRoot: (policies || []).includes('root'),
    });
  },

  async lookupTokenData(token, hasNamespace, { displayName, expireTime, namespacePath }) {
    // Only lookup if we're missing displayName or namespacePath in a non-root namespace
    if (!displayName || (!namespacePath && hasNamespace)) {
      try {
        const { data } = await this.lookupSelf(token);
        return {
          displayName: displayName || data?.display_name,
          namespacePath: namespacePath || data?.namespace_path,
          expireTime: expireTime || data?.expire_time,
        };
      } catch {
        // It would be unusual for this request to fail, but swallowing it because we're
        // essentially setting "nice to have" data here.
      }
    }
    // Return original values as fallback
    return { displayName, namespacePath, expireTime };
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
    const currentlyRenewing = this.isRenewing;
    if (currentlyRenewing) return;

    this.isRenewing = true;
    return this.renewCurrentToken().then(
      async (resp) => {
        this.isRenewing = false;
        // If we renewing, authData already exists so all we really need to update are the token and expiration details
        const { authMethodType, authMountPath, displayName } = this.authData;
        const normalizedAuthData = this.normalizeAuthData(resp.auth, {
          authMethodType,
          authMountPath,
          displayName,
        });
        return await this.persistAuthData(this.activeClusterId, normalizedAuthData);
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

  parseMfaResponse(mfaRequirement) {
    // mfa_requirement response comes back in a shape that is not easy to work with
    // convert to array of objects and add necessary properties to satisfy the view
    if (mfaRequirement) {
      const { mfa_request_id, mfa_constraints } = mfaRequirement;
      const constraints = [];
      for (const key in mfa_constraints) {
        const methods = mfa_constraints[key].any;
        // friendly label for display in MfaForm
        methods.forEach((m) => {
          const typeFormatted = m.type === 'totp' ? m.type.toUpperCase() : capitalize(m.type);
          m.label = `${typeFormatted} ${m.uses_passcode ? 'passcode' : 'push notification'}`;
        });
        constraints.push({ name: key, methods });
      }
      return { mfa_request_id, mfa_constraints: constraints };
    }
    return {};
  },

  async totpValidate({ clusterId, mfaRequirement, authMethodType, authMountPath }) {
    // mfa/validate consistently returns data inside the "auth" key
    const { auth } = await this.clusterAdapter().mfaValidate(mfaRequirement);
    const normalizedAuthData = this.normalizeAuthData(auth, { authMethodType, authMountPath });
    return this.authSuccess(clusterId, normalizedAuthData);
  },

  async authSuccess(clusterId, authResponse) {
    // persist selectedAuth to localStorage to rehydrate auth form on logout
    localStorage.setItem('selectedAuth', authResponse.authMethodType);
    const authData = await this.persistAuthData(clusterId, authResponse);
    this.permissions.getPaths.perform();
    return authData;
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

  // Depending on where auth happens (mfa/validate, renew-self or the method's login) the auth data
  // varies slightly (i.e. "ttl" vs "lease_duration"). Normalize it so stored authData contains consistent keys.
  // (Also, the API service returns camel cased keys and raw ajax requests return snake cased params.)
  normalizeAuthData(authData, { authMethodType, authMountPath, displayName, token, ttl }) {
    const displayNameFromMetadata = (metadata) =>
      metadata
        ? ['org', 'username']
            .map((key) => (key in metadata ? metadata[key] : null))
            .filter(Boolean)
            .join('/')
        : '';

    return {
      authMethodType,
      authMountPath,
      entityId: authData?.entity_id,
      expireTime: authData?.expire_time,
      token: token || authData?.client_token,
      renewable: authData?.renewable,
      ttl: ttl || authData?.lease_duration,
      policies: authData?.policies,
      mfaRequirement: authData?.mfa_requirement,
      // not all methods return a display name or metadata, if this is still empty it will be gleaned from lookup-self
      displayName: displayName || displayNameFromMetadata(authData?.metadata),
    };
  },
});
