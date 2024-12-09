/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Ember from 'ember';
import { service } from '@ember/service';
// ARG NOTE: Once you remove outer-html after glimmerizing you can remove the outer-html component
import Component from './outer-html';
import { task, timeout, waitForEvent } from 'ember-concurrency';
import { debounce } from '@ember/runloop';

const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete. Your web browser may have blocked or closed a pop-up window. Please check your settings and click Sign In to try again.';
const ERROR_MISSING_PARAMS =
  'The callback from the provider did not supply all of the required parameters.  Please click Sign In to try again. If the problem persists, you may want to contact your administrator.';
const ERROR_JWT_LOGIN = 'OIDC login is not configured for this mount';
export { ERROR_WINDOW_CLOSED, ERROR_MISSING_PARAMS, ERROR_JWT_LOGIN };

export default Component.extend({
  store: service(),
  flagsService: service('flags'),

  selectedAuthPath: null,
  selectedAuthType: null,
  roleName: null,
  role: null,
  errorMessage: null,
  isOIDC: true,

  onRoleName() {},
  onLoading() {},
  onError() {},
  onNamespace() {},

  didReceiveAttrs() {
    this._super(...arguments);
    // if mount path or type changes we need to check again for JWT configuration
    const didChangePath = this._authPath !== this.selectedAuthPath;
    const didChangeType = this._authType !== this.selectedAuthType;

    if (didChangePath || didChangeType) {
      // path updates as the user types so we need to debounce that event
      const wait = didChangePath ? 500 : 0;
      debounce(this, 'fetchRole', wait);
    }
    this._authPath = this.selectedAuthPath;
    this._authType = this.selectedAuthType;
  },

  getWindow() {
    return this.window || window;
  },

  async fetchRole() {
    const path = this.selectedAuthPath || this.selectedAuthType;
    const id = JSON.stringify([path, this.roleName]);
    this.setProperties({ role: null, errorMessage: null, isOIDC: true });

    try {
      const role = await this.store.findRecord('role-jwt', id, {
        adapterOptions: { namespace: this.namespace },
      });
      this.set('role', role);
    } catch (e) {
      const error = (e.errors || [])[0];
      const errorMessage =
        e.httpStatus === 400 ? 'Invalid role. Please try again.' : `Error fetching role: ${error}`;
      // assume OIDC until it's known that the mount is configured for JWT authentication via static keys, JWKS, or OIDC discovery.
      this.setProperties({ isOIDC: error !== ERROR_JWT_LOGIN, errorMessage });
    }
  },

  cancelLogin(oidcWindow, errorMessage) {
    this.closeWindow(oidcWindow);
    this.handleOIDCError(errorMessage);
  },

  closeWindow(oidcWindow) {
    this.watchPopup.cancelAll();
    this.watchCurrent.cancelAll();
    oidcWindow.close();
  },

  handleOIDCError(err) {
    this.onLoading(false);
    this.prepareForOIDC.cancelAll();
    this.onError(err);
  },

  // NOTE TO DEVS: Be careful when updating the OIDC flow and ensure the updates
  // work with implicit flow. See issue https://github.com/hashicorp/vault-plugin-auth-jwt/pull/192
  prepareForOIDC: task(function* (oidcWindow) {
    const thisWindow = this.getWindow();
    // show the loading animation in the parent
    this.onLoading(true);
    // start watching the popup window and the current one
    this.watchPopup.perform(oidcWindow);
    this.watchCurrent.perform(oidcWindow);
    // wait for message posted from oidc callback
    // see issue https://github.com/hashicorp/vault/issues/12436
    // ensure that postMessage event is from expected source
    while (true) {
      const event = yield waitForEvent(thisWindow, 'message');
      if (event.origin === thisWindow.origin && event.isTrusted && event.data.source === 'oidc-callback') {
        return this.exchangeOIDC.perform(event.data, oidcWindow);
      }
      // continue to wait for the correct message
    }
  }),

  watchPopup: task(function* (oidcWindow) {
    while (true) {
      const WAIT_TIME = Ember.testing ? 50 : 500;

      yield timeout(WAIT_TIME);
      if (!oidcWindow || oidcWindow.closed) {
        return this.handleOIDCError(ERROR_WINDOW_CLOSED);
      }
    }
  }),

  watchCurrent: task(function* (oidcWindow) {
    // when user is about to change pages, close the popup window
    yield waitForEvent(this.getWindow(), 'beforeunload');
    oidcWindow.close();
  }),

  exchangeOIDC: task(function* (oidcState, oidcWindow) {
    if (oidcState === null || oidcState === undefined) {
      return;
    }
    this.onLoading(true);

    let { namespace, path, state, code } = oidcState;

    // The namespace can be either be passed as a query parameter, or be embedded
    // in the state param in the format `<state_id>,ns=<namespace>`. So if
    // `namespace` is empty, check for namespace in state as well.
    if (namespace === '' || this.flagsService.hvdManagedNamespaceRoot) {
      const i = state.indexOf(',ns=');
      if (i >= 0) {
        // ",ns=" is 4 characters
        namespace = state.substring(i + 4);
        state = state.substring(0, i);
      }
    }

    if (!path || !state || !code) {
      return this.cancelLogin(oidcWindow, ERROR_MISSING_PARAMS);
    }
    const adapter = this.store.adapterFor('auth-method');
    this.onNamespace(namespace);
    let resp;
    // do the OIDC exchange, set the token on the parent component
    // and submit auth form
    try {
      resp = yield adapter.exchangeOIDC(path, state, code);
      this.closeWindow(oidcWindow);
    } catch (e) {
      // If there was an error on Vault's end, close the popup
      // and show the error on the login screen
      return this.cancelLogin(oidcWindow, e);
    }
    const { mfa_requirement, client_token } = resp.auth;
    // onSubmit calls doSubmit in auth-form.js
    yield this.onSubmit({ mfa_requirement }, null, client_token);
  }),

  async startOIDCAuth() {
    this.onError(null);

    await this.fetchRole();

    const error =
      this.role && !this.role.authUrl
        ? 'Missing auth_url. Please check that allowed_redirect_uris for the role include this mount path.'
        : this.errorMessage || null;

    if (error) {
      this.onError(error);
    } else {
      const win = this.getWindow();
      const POPUP_WIDTH = 500;
      const POPUP_HEIGHT = 600;
      const left = win.screen.width / 2 - POPUP_WIDTH / 2;
      const top = win.screen.height / 2 - POPUP_HEIGHT / 2;
      const oidcWindow = win.open(
        this.role.authUrl,
        'vaultOIDCWindow',
        `width=${POPUP_WIDTH},height=${POPUP_HEIGHT},resizable,scrollbars=yes,top=${top},left=${left}`
      );

      this.prepareForOIDC.perform(oidcWindow);
    }
  },

  actions: {
    onRoleChange(event) {
      this.onRoleName(event.target.value);
      debounce(this, 'fetchRole', 500);
    },
    signIn(event) {
      event.preventDefault();

      if (this.isOIDC) {
        this.startOIDCAuth();
      } else {
        const { jwt, roleName: role } = this;
        this.onSubmit({ role, jwt });
      }
    },
  },
});
