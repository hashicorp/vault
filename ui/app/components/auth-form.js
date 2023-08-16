import Ember from 'ember';
import { next } from '@ember/runloop';
import { inject as service } from '@ember/service';
import { match, alias, or } from '@ember/object/computed';
import { dasherize } from '@ember/string';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import { task, timeout } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

const BACKENDS = supportedAuthBackends();

/**
 * @module AuthForm
 * The `AuthForm` is used to sign users into Vault.
 *
 * @example ```js
 * // All properties are passed in via query params.
 * <AuthForm @wrappedToken={{wrappedToken}} @cluster={{model}} @namespace={{namespaceQueryParam}} @selectedAuth={{authMethod}} @onSuccess={{action this.onSuccess}} />```
 *
 * @param {string} wrappedToken - The auth method that is currently selected in the dropdown.
 * @param {object} cluster - The auth method that is currently selected in the dropdown. This corresponds to an Ember Model.
 * @param {string} namespace- The currently active namespace.
 * @param {string} selectedAuth - The auth method that is currently selected in the dropdown.
 * @param {function} onSuccess - Fired on auth success
 */

const DEFAULTS = {
  token: null,
  username: null,
  password: null,
  customPath: null,
};

export default Component.extend(DEFAULTS, {
  router: service(),
  auth: service(),
  flashMessages: service(),
  store: service(),
  csp: service('csp-event'),

  //  passed in via a query param
  selectedAuth: null,
  methods: null,
  cluster: null,
  namespace: null,
  wrappedToken: null,
  // internal
  oldNamespace: null,
  authMethods: BACKENDS,

  didReceiveAttrs() {
    this._super(...arguments);
    let {
      wrappedToken: token,
      oldWrappedToken: oldToken,
      oldNamespace: oldNS,
      namespace: ns,
      selectedAuth: newMethod,
      oldSelectedAuth: oldMethod,
    } = this;

    next(() => {
      if (!token && (oldNS === null || oldNS !== ns)) {
        this.fetchMethods.perform();
      }
      this.set('oldNamespace', ns);
      // we only want to trigger this once
      if (token && !oldToken) {
        this.unwrapToken.perform(token);
        this.set('oldWrappedToken', token);
      }
      if (oldMethod && oldMethod !== newMethod) {
        this.resetDefaults();
      }
      this.set('oldSelectedAuth', newMethod);
    });
  },

  didRender() {
    this._super(...arguments);
    // on very narrow viewports the active tab may be overflowed, so we scroll it into view here
    let activeEle = this.element.querySelector('li.is-active');
    if (activeEle) {
      activeEle.scrollIntoView();
    }

    next(() => {
      let firstMethod = this.firstMethod();
      // set `with` to the first method
      if (
        !this.wrappedToken &&
        ((this.fetchMethods.isIdle && firstMethod && !this.selectedAuth) ||
          (this.selectedAuth && !this.selectedAuthBackend))
      ) {
        this.set('selectedAuth', firstMethod);
      }
    });
  },

  firstMethod() {
    let firstMethod = this.methodsToShow.firstObject;
    if (!firstMethod) return;
    // prefer backends with a path over those with a type
    return firstMethod.path || firstMethod.type;
  },

  resetDefaults() {
    this.setProperties(DEFAULTS);
  },

  getAuthBackend(type) {
    const { wrappedToken, methods, selectedAuth, selectedAuthIsPath: keyIsPath } = this;
    const selected = type || selectedAuth;
    if (!methods && !wrappedToken) {
      return {};
    }
    // if type is provided we can ignore path since we are attempting to lookup a specific backend by type
    if (keyIsPath && !type) {
      return methods.findBy('path', selected);
    }
    return BACKENDS.findBy('type', selected);
  },

  selectedAuthIsPath: match('selectedAuth', /\/$/),
  selectedAuthBackend: computed(
    'wrappedToken',
    'methods',
    'methods.[]',
    'selectedAuth',
    'selectedAuthIsPath',
    function () {
      return this.getAuthBackend();
    }
  ),

  providerName: computed('selectedAuthBackend.type', function () {
    if (!this.selectedAuthBackend) {
      return;
    }
    let type = this.selectedAuthBackend.type || 'token';
    type = type.toLowerCase();
    let templateName = dasherize(type);
    return templateName;
  }),

  hasCSPError: alias('csp.connectionViolations.firstObject'),

  cspErrorText: `This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.`,

  allSupportedMethods: computed('methodsToShow', 'hasMethodsWithPath', function () {
    let hasMethodsWithPath = this.hasMethodsWithPath;
    let methodsToShow = this.methodsToShow;
    return hasMethodsWithPath ? methodsToShow.concat(BACKENDS) : methodsToShow;
  }),

  hasMethodsWithPath: computed('methodsToShow', function () {
    return this.methodsToShow.isAny('path');
  }),
  methodsToShow: computed('methods', function () {
    let methods = this.methods || [];
    let shownMethods = methods.filter((m) =>
      BACKENDS.find((b) => b.type.toLowerCase() === m.type.toLowerCase())
    );
    return shownMethods.length ? shownMethods : BACKENDS;
  }),

  unwrapToken: task(
    waitFor(function* (token) {
      // will be using the Token Auth Method, so set it here
      this.set('selectedAuth', 'token');
      let adapter = this.store.adapterFor('tools');
      try {
        let response = yield adapter.toolAction('unwrap', null, { clientToken: token });
        this.set('token', response.auth.client_token);
        this.send('doSubmit');
      } catch (e) {
        this.set('error', `Token unwrap failed: ${e.errors[0]}`);
      }
    })
  ),

  fetchMethods: task(
    waitFor(function* () {
      let store = this.store;
      try {
        let methods = yield store.findAll('auth-method', {
          adapterOptions: {
            unauthenticated: true,
          },
        });
        this.set(
          'methods',
          methods.map((m) => {
            const method = m.serialize({ includeId: true });
            return {
              ...method,
              mountDescription: method.description,
            };
          })
        );
        next(() => {
          store.unloadAll('auth-method');
        });
      } catch (e) {
        this.set('error', `There was an error fetching Auth Methods: ${e.errors[0]}`);
      }
    })
  ),

  showLoading: or('isLoading', 'authenticate.isRunning', 'fetchMethods.isRunning', 'unwrapToken.isRunning'),

  authenticate: task(
    waitFor(function* (backendType, data) {
      const {
        selectedAuth,
        cluster: { id: clusterId },
      } = this;
      try {
        this.delayAuthMessageReminder.perform();
        const authResponse = yield this.auth.authenticate({
          clusterId,
          backend: backendType,
          data,
          selectedAuth,
        });
        this.onSuccess(authResponse, backendType, data);
      } catch (e) {
        this.set('isLoading', false);
        if (!this.auth.mfaError) {
          this.set('error', `Authentication failed: ${this.auth.handleError(e)}`);
        }
      }
    })
  ),

  delayAuthMessageReminder: task(function* () {
    if (Ember.testing) {
      this.showLoading = true;
      yield timeout(0);
    } else {
      yield timeout(5000);
    }
  }),

  actions: {
    doSubmit(passedData, event, token) {
      if (event) {
        event.preventDefault();
      }
      if (token) {
        this.set('token', token);
      }
      this.set('error', null);
      // if callback from oidc or jwt we have a token at this point
      const backend = token ? this.getAuthBackend('token') : this.selectedAuthBackend || {};
      const backendMeta = BACKENDS.find(
        (b) => (b.type || '').toLowerCase() === (backend.type || '').toLowerCase()
      );
      const attributes = (backendMeta || {}).formAttributes || [];
      const data = this.getProperties(...attributes);

      if (passedData) {
        Object.assign(data, passedData);
      }
      if (this.customPath || backend.id) {
        data.path = this.customPath || backend.id;
      }
      return this.authenticate.unlinked().perform(backend.type, data);
    },
    handleError(e) {
      this.setProperties({
        isLoading: false,
        error: e ? this.auth.handleError(e) : null,
      });
    },
  },
});
