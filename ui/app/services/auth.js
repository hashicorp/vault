import Ember from 'ember';
import { resolve } from 'rsvp';
import { assign } from '@ember/polyfills';
import $ from 'jquery';
import { isArray } from '@ember/array';
import { computed, get } from '@ember/object';
import { getOwner } from '@ember/application';
import Service, { inject as service } from '@ember/service';
import getStorage from '../lib/token-storage';
import ENV from 'vault/config/environment';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import { task, timeout } from 'ember-concurrency';
const TOKEN_SEPARATOR = '☃';
const TOKEN_PREFIX = 'vault-';
const ROOT_PREFIX = '🗝';
const BACKENDS = supportedAuthBackends();

export { TOKEN_SEPARATOR, TOKEN_PREFIX, ROOT_PREFIX };

export default Service.extend({
  permissions: service(),
  namespace: service(),
  IDLE_TIMEOUT: 3 * 60e3,
  expirationCalcTS: null,
  init() {
    this._super(...arguments);
    this.checkForRootToken();
  },

  clusterAdapter() {
    return getOwner(this).lookup('adapter:cluster');
  },

  tokens: computed(function() {
    return this.getTokensFromStorage() || [];
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
        'X-Vault-Token': this.get('currentToken'),
      },
    };

    let namespace =
      typeof options.namespace === 'undefined' ? this.get('namespaceService.path') : options.namespace;
    if (namespace) {
      defaults.headers['X-Vault-Namespace'] = namespace;
    }
    return $.ajax(assign(defaults, options));
  },

  renewCurrentToken() {
    let namespace = this.get('authData.userRootNamespace');
    const url = '/v1/auth/token/renew-self';
    return this.ajax(url, 'POST', { namespace });
  },

  revokeCurrentToken() {
    let namespace = this.get('authData.userRootNamespace');
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
    let tokens = this.get('tokens');
    let currentNamespace = this.get('namespace.path') || '';
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
      displayName = currentBackend.displayNamePath.map(name => get(resp, name)).join('/');
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
      userRootNamespace = this.get('authData.userRootNamespace');
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
        clusterId: (options && options.clusterId) || this.get('activeCluster'),
      },
      resp.policies
    );

    if (resp.renewable) {
      assign(data, this.calculateExpiration(resp));
    }

    if (!data.displayName) {
      data.displayName = get(this.getTokenData(tokenName) || {}, 'displayName');
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

  tokenExpirationDate: computed('currentTokenName', 'expirationCalcTS', function() {
    const tokenName = this.get('currentTokenName');
    if (!tokenName) {
      return;
    }
    const { tokenExpirationEpoch } = this.getTokenData(tokenName);
    const expirationDate = new Date(0);
    return tokenExpirationEpoch ? expirationDate.setUTCMilliseconds(tokenExpirationEpoch) : null;
  }),

  tokenExpired: computed(function() {
    const expiration = this.get('tokenExpirationDate');
    return expiration ? this.now() >= expiration : null;
  }).volatile(),

  renewAfterEpoch: computed('currentTokenName', 'expirationCalcTS', function() {
    const tokenName = this.get('currentTokenName');
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
    const tokenName = this.get('currentTokenName');
    const currentlyRenewing = this.get('isRenewing');
    if (currentlyRenewing) {
      return;
    }
    this.set('isRenewing', true);
    return this.renewCurrentToken().then(
      resp => {
        this.set('isRenewing', false);
        return this.persistAuthData(tokenName, resp.data || resp.auth);
      },
      e => {
        this.set('isRenewing', false);
        throw e;
      }
    );
  },

  checkShouldRenew: task(function*() {
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
    const lastFetch = this.get('lastFetch');
    const renewTime = this.get('renewAfterEpoch');
    if (!this.currentTokenName || this.get('tokenExpired') || this.get('allowExpiration') || !renewTime) {
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
    this.set('lastFetch', timestamp);
    // if expiration was allowed we want to go ahead and renew here
    if (this.allowExpiration) {
      this.renew();
    }
    this.set('allowExpiration', false);
  },

  getTokensFromStorage(filterFn) {
    return this.storage()
      .keys()
      .reject(key => {
        return key.indexOf(TOKEN_PREFIX) !== 0 || (filterFn && filterFn(key));
      });
  },

  checkForRootToken() {
    if (this.environment() === 'development') {
      return;
    }
    this.getTokensFromStorage().forEach(key => {
      const data = this.getTokenData(key);
      if (data.policies.includes('root')) {
        this.removeTokenData(key);
      }
    });
  },

  authenticate(/*{clusterId, backend, data}*/) {
    const [options] = arguments;
    const adapter = this.clusterAdapter();

    return adapter.authenticate(options).then(resp => {
      return this.persistAuthData(options, resp.auth || resp.data, this.get('namespace.path')).then(
        authData => {
          return this.get('permissions')
            .getPaths.perform()
            .then(() => {
              return authData;
            });
        }
      );
    });
  },

  deleteCurrentToken() {
    const tokenName = this.get('currentTokenName');
    this.deleteToken(tokenName);
    this.removeTokenData(tokenName);
  },

  deleteToken(tokenName) {
    const tokenNames = this.get('tokens').without(tokenName);
    this.removeTokenData(tokenName);
    this.set('tokens', tokenNames);
  },

  // returns the key for the token to use
  currentTokenName: computed('activeCluster', 'tokens', 'tokens.[]', function() {
    const regex = new RegExp(this.get('activeCluster'));
    return this.get('tokens').find(key => regex.test(key));
  }),

  currentToken: computed('currentTokenName', function() {
    const name = this.get('currentTokenName');
    const data = name && this.getTokenData(name);
    return name && data ? data.token : null;
  }),

  authData: computed('currentTokenName', function() {
    const token = this.get('currentTokenName');
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
