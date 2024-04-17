/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import { next } from '@ember/runloop';
import { service } from '@ember/service';
import { match, or } from '@ember/object/computed';
import { dasherize } from '@ember/string';
import Component from '@ember/component';
import { computed } from '@ember/object';
import { allSupportedAuthBackends, supportedAuthBackends } from 'vault/helpers/supported-auth-backends';
import { task, timeout } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { v4 as uuidv4 } from 'uuid';

/**
 * @module AuthForm
 * The `AuthForm` is used to sign users into Vault.
 *
 * @example ```js
 * // All properties are passed in via query params.
 * <AuthForm @wrappedToken={{wrappedToken}} @cluster={{model}} @namespace={{namespaceQueryParam}} @selectedAuth={{authMethod}} @onSuccess={{action this.onSuccess}}/>```
 *
 * @param {string} wrappedToken - The auth method that is currently selected in the dropdown.
 * @param {object} cluster - The auth method that is currently selected in the dropdown. This corresponds to an Ember Model.
 * @param {string} namespace- The currently active namespace.
 * @param {string} selectedAuth - The auth method that is currently selected in the dropdown.
 * @param {function} onSuccess - Fired on auth success.
 * @param {function} [setOktaNumberChallenge] - Sets whether we are waiting for okta number challenge to be used to sign in.
 * @param {boolean} [waitingForOktaNumberChallenge=false] - Determines if we are waiting for the Okta Number Challenge to sign in.
 * @param {function} [setCancellingAuth] - Sets whether we are cancelling or not the login authentication for Okta Number Challenge.
 * @param {boolean} [cancelAuthForOktaNumberChallenge=false] - Determines if we are cancelling the login authentication for the Okta Number Challenge.
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
  version: service(),

  //  passed in via a query param
  selectedAuth: null,
  methods: null,
  cluster: null,
  namespace: null,
  wrappedToken: null,
  // internal
  oldNamespace: null,

  // number answer for okta number challenge if applicable
  oktaNumberChallengeAnswer: null,

  authMethods: computed('version.isEnterprise', function () {
    return this.version.isEnterprise ? allSupportedAuthBackends() : supportedAuthBackends();
  }),

  didReceiveAttrs() {
    this._super(...arguments);
    const {
      wrappedToken: token,
      oldWrappedToken: oldToken,
      oldNamespace: oldNS,
      namespace: ns,
      selectedAuth: newMethod,
      oldSelectedAuth: oldMethod,
      cancelAuthForOktaNumberChallenge: cancelAuth,
    } = this;
    // if we are cancelling the login then we reset the number challenge answer and cancel the current authenticate and polling tasks
    if (cancelAuth) {
      this.set('oktaNumberChallengeAnswer', null);
      this.authenticate.cancelAll();
      this.pollForOktaNumberChallenge.cancelAll();
    }
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
    const activeEle = this.element.querySelector('li.is-active');
    if (activeEle) {
      activeEle.scrollIntoView();
    }

    next(() => {
      const firstMethod = this.firstMethod();
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
    const firstMethod = this.methodsToShow[0];
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
      return methods.find((m) => m.path === selected);
    }
    return this.authMethods.find((m) => m.type === selected);
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
    const templateName = dasherize(type);
    return templateName;
  }),

  cspError: computed('csp.connectionViolations.length', function () {
    if (this.csp.connectionViolations.length) {
      return `This is a standby Vault node but can't communicate with the active node via request forwarding. Sign in at the active node to use the Vault UI.`;
    }
    return '';
  }),

  allSupportedMethods: computed('methodsToShow', 'hasMethodsWithPath', 'authMethods', function () {
    const hasMethodsWithPath = this.hasMethodsWithPath;
    const methodsToShow = this.methodsToShow;
    return hasMethodsWithPath ? methodsToShow.concat(this.authMethods) : methodsToShow;
  }),

  hasMethodsWithPath: computed('methodsToShow', function () {
    return this.methodsToShow.isAny('path');
  }),
  methodsToShow: computed('methods', 'authMethods', function () {
    const methods = this.methods || [];
    const shownMethods = methods.filter((m) =>
      this.authMethods.find((b) => b.type.toLowerCase() === m.type.toLowerCase())
    );
    return shownMethods.length ? shownMethods : this.authMethods;
  }),

  unwrapToken: task(
    waitFor(function* (token) {
      // will be using the Token Auth Method, so set it here
      this.set('selectedAuth', 'token');
      const adapter = this.store.adapterFor('tools');
      try {
        const response = yield adapter.toolAction('unwrap', null, { clientToken: token });
        this.set('token', response.auth.client_token);
        this.send('doSubmit');
      } catch (e) {
        this.set('error', `Token unwrap failed: ${e.errors[0]}`);
      }
    })
  ),

  fetchMethods: task(
    waitFor(function* () {
      const store = this.store;
      try {
        const methods = yield store.findAll('auth-method', {
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
        // without unloading the records there will be an issue where all methods set to list when unauthenticated will appear for all namespaces
        // if possible, it would be more reliable to add a namespace attr to the model so we could filter against the current namespace rather than unloading all
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
        if (backendType === 'okta') {
          this.pollForOktaNumberChallenge.perform(data.nonce, data.path);
        } else {
          this.delayAuthMessageReminder.perform();
        }
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

  pollForOktaNumberChallenge: task(function* (nonce, mount) {
    // yield for 1s to wait to see if there is a login error before polling
    yield timeout(1000);
    if (this.error) {
      return;
    }
    let response = null;
    this.setOktaNumberChallenge(true);
    this.setCancellingAuth(false);
    // keep polling /auth/okta/verify/:nonce API every 1s until a response is given with the correct number for the Okta Number Challenge
    while (response === null) {
      // when testing, the polling loop causes promises to be rejected making acceptance tests fail
      // so disable the poll in tests
      if (Ember.testing) {
        return;
      }
      yield timeout(1000);
      response = yield this.auth.getOktaNumberChallengeAnswer(nonce, mount);
    }
    this.set('oktaNumberChallengeAnswer', response);
  }),

  delayAuthMessageReminder: task(function* () {
    if (Ember.testing) {
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
      // if callback from oidc, jwt, or saml we have a token at this point
      const backend = token ? this.getAuthBackend('token') : this.selectedAuthBackend || {};
      const backendMeta = this.authMethods.find(
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
      // add nonce field for okta backend
      if (backend.type === 'okta') {
        data.nonce = uuidv4();
        // add a default path of okta if it doesn't exist to be used for Okta Number Challenge
        if (!data.path) {
          data.path = 'okta';
        }
      }
      return this.authenticate.unlinked().perform(backend.type, data);
    },
    handleError(e) {
      this.setProperties({
        isLoading: false,
        error: e ? this.auth.handleError(e) : null,
      });
    },
    returnToLoginFromOktaNumberChallenge() {
      this.setOktaNumberChallenge(false);
      this.set('oktaNumberChallengeAnswer', null);
    },
  },
});
