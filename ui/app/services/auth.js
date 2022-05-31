import Ember from 'ember';
import { resolve, reject } from 'rsvp';
import { assign } from '@ember/polyfills';
import { isArray } from '@ember/array';
import { computed, get } from '@ember/object';
import { capitalize } from '@ember/string';

import fetch from 'fetch';
import { getOwner } from '@ember/application';
import Service, { inject as service } from '@ember/service';
import getStorage from '../lib/token-storage';
import ENV from 'vault/config/environment';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import { task, timeout } from 'ember-concurrency';
const TOKEN_SEPARATOR = '☃';
const TOKEN_PREFIX = 'vault-';
const ROOT_PREFIX = '_root_';
const BACKENDS = supportedAuthBackends();

export { TOKEN_SEPARATOR, TOKEN_PREFIX, ROOT_PREFIX };

export default Service.extend({
  permissions: service(),
  namespaceService: service('namespace'),
  IDLE_TIMEOUT: 3 * 60e3,
  expirationCalcTS: null,
  isRenewing: false,
  mfaErrors: null,

  init() {
    this._super(...arguments);
    this.checkForRootToken();
  },

  clusterAdapter() {
    return getOwner(this).lookup('adapter:cluster');
  },
  // eslint-disable-next-line
  tokens: computed({
    get() {
      return this._tokens || this.getTokensFromStorage() || [];
    },
    set(key, value) {
      return (this._tokens = value);
    },
  }),

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
    this.set('activeCluster', clusterId);
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

    let namespace = typeof options.namespace === 'undefined' ? this.namespaceService.path : options.namespace;
    if (namespace) {
      defaults.headers['X-Vault-Namespace'] = namespace;
    }
    let opts = assign(defaults, options);

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
    let namespace = this.authData.userRootNamespace;
    const url = '/v1/auth/token/renew-self';
    return this.ajax(url, 'POST', { namespace });
  },

  revokeCurrentToken() {
    let namespace = this.authData.userRootNamespace;
    const url = '/v1/auth/token/revoke-self';
    return this.ajax(url, 'POST', { namespace });
  },

  calculateExpiration(resp) {
    let now = this.now();
    const ttl = resp.ttl || resp.lease_duration;
    const tokenExpirationEpoch = now + ttl * 1e3;
    this.set('expirationCalcTS', now);
    return {
      ttl,
      tokenExpirationEpoch,
    };
  },

  persistAuthData() {
    let [firstArg, resp] = arguments;
    let tokens = this.tokens;
    let currentNamespace = this.namespaceService.path || '';
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

    let currentBackend = BACKENDS.findBy('type', backend);
    let displayName;
    if (isArray(currentBackend.displayNamePath)) {
      displayName = currentBackend.displayNamePath.map((name) => get(resp, name)).join('/');
    } else {
      displayName = get(resp, currentBackend.displayNamePath);
    }

    let { entity_id, policies, renewable, namespace_path } = resp;
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
    let data = {
      userRootNamespace,
      displayName,
      backend: currentBackend,
      token: resp.client_token || get(resp, currentBackend.tokenPath),
      policies,
      renewable,
      entity_id,
    };

    tokenName = this.generateTokenName(
      {
        backend,
        clusterId: (options && options.clusterId) || this.activeCluster,
      },
      resp.policies
    );

    if (resp.renewable) {
      assign(data, this.calculateExpiration(resp));
    }

    if (!data.displayName) {
      data.displayName = (this.getTokenData(tokenName) || {}).displayName;
    }
    tokens.addObject(tokenName);
    this.set('tokens', tokens);
    this.set('allowExpiration', false);
    this.setTokenData(tokenName, data);
    return resolve({
      namespace: currentNamespace || data.userRootNamespace,
      token: tokenName,
      isRoot: policies.includes('root'),
    });
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

  tokenExpirationDate: computed('currentTokenName', 'expirationCalcTS', function () {
    const tokenName = this.currentTokenName;
    if (!tokenName) {
      return;
    }
    const { tokenExpirationEpoch } = this.getTokenData(tokenName);
    const expirationDate = new Date(0);
    return tokenExpirationEpoch ? expirationDate.setUTCMilliseconds(tokenExpirationEpoch) : null;
  }),

  get tokenExpired() {
    const expiration = this.tokenExpirationDate;
    return expiration ? this.now() >= expiration : null;
  },

  renewAfterEpoch: computed('currentTokenName', 'expirationCalcTS', function () {
    const tokenName = this.currentTokenName;
    let { expirationCalcTS } = this;
    const data = this.getTokenData(tokenName);
    if (!tokenName || !data || !expirationCalcTS) {
      return null;
    }
    const { ttl, renewable } = data;
    // renew after last expirationCalc time + half of the ttl (in ms)
    return renewable ? Math.floor((ttl * 1e3) / 2) + expirationCalcTS : null;
  }),

  renew() {
    const tokenName = this.currentTokenName;
    const currentlyRenewing = this.isRenewing;
    if (currentlyRenewing) {
      return;
    }
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
    if (this.allowExpiration && now >= this.renewAfterEpoch) {
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
      for (let key in mfa_constraints) {
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
    // persist selectedAuth to sessionStorage to rehydrate auth form on logout
    sessionStorage.setItem('selectedAuth', options.selectedAuth);
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
    // check sessionStorage first
    const selectedAuth = sessionStorage.getItem('selectedAuth');
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

  // returns the key for the token to use
  currentTokenName: computed('activeCluster', 'tokens', 'tokens.[]', function () {
    const regex = new RegExp(this.activeCluster);
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

    return assign(stored, {
      backend: BACKENDS.findBy('type', backend),
    });
  }),
});
