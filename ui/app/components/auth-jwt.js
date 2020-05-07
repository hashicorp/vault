import Ember from 'ember';
import { inject as service } from '@ember/service';
import Component from './outer-html';
import { later } from '@ember/runloop';
import { task, timeout, waitForEvent } from 'ember-concurrency';
import { computed } from '@ember/object';

/* eslint-disable ember/no-ember-testing-in-module-scope */
const WAIT_TIME = Ember.testing ? 0 : 500;
const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete.  Please click Sign In to try again.';
const ERROR_MISSING_PARAMS =
  'The callback from the provider did not supply all of the required parameters.  Please click Sign In to try again. If the problem persists, you may want to contact your administrator.';
const ERROR_JWT_LOGIN = 'OIDC login is not configured for this mount';
export { ERROR_WINDOW_CLOSED, ERROR_MISSING_PARAMS, ERROR_JWT_LOGIN };

export default Component.extend({
  store: service(),
  selectedAuthPath: null,
  selectedAuthType: null,
  roleName: null,
  role: null,
  errorMessage: null,
  onRoleName() {},
  onLoading() {},
  onError() {},
  onToken() {},
  onNamespace() {},

  didReceiveAttrs() {
    console.log(`didReceiveAttrs: ${this.errorMessage}`);
    let { oldSelectedAuthPath, selectedAuthPath } = this;
    let shouldDebounce = !oldSelectedAuthPath && !selectedAuthPath;
    if (oldSelectedAuthPath !== selectedAuthPath) {
      this.set('role', null);
      this.onRoleName(this.roleName);
      this.fetchRole.perform(null, { debounce: false });
    } else if (shouldDebounce) {
      this.fetchRole.perform(this.roleName);
    }
    this.set('errorMessage', null);
    this.set('oldSelectedAuthPath', selectedAuthPath);
  },

  // Assumes authentication using OIDC until it's known that the mount is
  // configured for JWT authentication via static keys, JWKS, or OIDC discovery.
  isOIDC: computed('errorMessage', function() {
    console.log(`isOIDC: ${this.errorMessage}`);
    return this.errorMessage !== ERROR_JWT_LOGIN;
  }),

  getWindow() {
    return this.window || window;
  },

  fetchRole: task(function*(roleName, options = { debounce: true }) {
    if (options.debounce) {
      this.onRoleName(roleName);
      // debounce
      yield timeout(WAIT_TIME);
    }
    let path = this.selectedAuthPath || this.selectedAuthType;
    let id = JSON.stringify([path, roleName]);
    let role = null;
    try {
      role = yield this.store.findRecord('role-jwt', id, { adapterOptions: { namespace: this.namespace } });
    } catch (e) {
      if (!e.httpStatus || e.httpStatus !== 400) {
        throw e;
      }
      if (e.errors && e.errors.length > 0) {
        this.set('errorMessage', e.errors[0]);
      }
    }
    this.set('role', role);
  })
    .restartable()
    .withTestWaiter(),

  handleOIDCError(err) {
    this.onLoading(false);
    this.prepareForOIDC.cancelAll();
    this.onError(err);
  },

  prepareForOIDC: task(function*(oidcWindow) {
    // show the loading animation in the parent
    this.onLoading(true);
    // start watching the popup window and the current one
    this.watchPopup.perform(oidcWindow);
    this.watchCurrent.perform(oidcWindow);
    // and then wait for storage event to be fired from the popup
    // window setting a value in localStorage when the callback route is loaded
    let storageEvent = yield waitForEvent(this.getWindow(), 'storage');
    this.exchangeOIDC.perform(storageEvent, oidcWindow);
  }),

  watchPopup: task(function*(oidcWindow) {
    while (true) {
      yield timeout(WAIT_TIME);
      if (!oidcWindow || oidcWindow.closed) {
        return this.handleOIDCError(ERROR_WINDOW_CLOSED);
      }
    }
  }),

  watchCurrent: task(function*(oidcWindow) {
    yield waitForEvent(this.getWindow(), 'beforeunload');
    oidcWindow.close();
  }),

  closeWindow(oidcWindow) {
    this.watchPopup.cancelAll();
    this.watchCurrent.cancelAll();
    oidcWindow.close();
  },

  exchangeOIDC: task(function*(event, oidcWindow) {
    // in non-incognito mode we need to use a timeout because it takes time before oidcState is written to local storage.
    let oidcState = Ember.testing
      ? event.storageArea.getItem('oidcState')
      : yield timeout(1000).then(() => event.storageArea.getItem('oidcState'));

    if (oidcState === null || oidcState === undefined) {
      return;
    }
    this.onLoading(true);
    // get the info from the event fired by the other window and
    // then remove it from localStorage
    let { namespace, path, state, code } = JSON.parse(oidcState);
    this.getWindow().localStorage.removeItem('oidcState');

    // defer closing of the window, but continue executing the task
    later(() => {
      this.closeWindow(oidcWindow);
    }, WAIT_TIME);
    if (!path || !state || !code) {
      return this.handleOIDCError(ERROR_MISSING_PARAMS);
    }
    let adapter = this.store.adapterFor('auth-method');
    this.onNamespace(namespace);
    let resp;
    // do the OIDC exchange, set the token on the parent component
    // and submit auth form
    try {
      resp = yield adapter.exchangeOIDC(path, state, code);
    } catch (e) {
      return this.handleOIDCError(e);
    }
    let token = resp.auth.client_token;
    this.onSelectedAuth('token');
    this.onToken(token);
    yield this.onSubmit();
  }),

  actions: {
    async startOIDCAuth(data, e) {
      this.onError(null);
      if (e && e.preventDefault) {
        e.preventDefault();
      }
      if (!this.isOIDC) {
        return;
      }

      console.log(`Role: ${this.role}`);

      await this.fetchRole.perform(this.roleName, { debounce: false });
      let win = this.getWindow();

      const POPUP_WIDTH = 500;
      const POPUP_HEIGHT = 600;
      let left = win.screen.width / 2 - POPUP_WIDTH / 2;
      let top = win.screen.height / 2 - POPUP_HEIGHT / 2;
      let oidcWindow = win.open(
        this.role.authUrl,
        'vaultOIDCWindow',
        `width=${POPUP_WIDTH},height=${POPUP_HEIGHT},resizable,scrollbars=yes,top=${top},left=${left}`
      );

      this.prepareForOIDC.perform(oidcWindow);
    },
  },
});
