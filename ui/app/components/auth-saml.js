/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { inject as service } from '@ember/service';
import Component from './outer-html';
import { task, timeout, waitForEvent } from 'ember-concurrency';
import { computed } from '@ember/object';
import errorMessage from 'vault/utils/error-message';

const WAIT_TIME = 500;
const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete. Your web browser may have blocked or closed a pop-up window. Please check your settings and click Sign In to try again.';
const ERROR_TIMEOUT = 'The authentication request has timed out. Please click Sign In to try again.';
const ERROR_MISSING_PARAMS =
  'The callback from the provider did not supply all of the required parameters. Please click Sign In to try again. If the problem persists, you may want to contact your administrator.';
export { ERROR_WINDOW_CLOSED, ERROR_MISSING_PARAMS };

export default Component.extend({
  store: service(),
  featureFlagService: service('featureFlag'),

  selectedAuthPath: null,
  selectedAuthType: null,
  roleName: null,
  errorMessage: null,
  onRoleName() {},
  onLoading() {},
  onError() {},
  onNamespace() {},

  didReceiveAttrs() {
    this._super();
    this.set('errorMessage', null);
  },

  getWindow() {
    return this.window || window;
  },

  canLoginSaml: computed('getWindow', function () {
    return this.getWindow().isSecureContext;
  }),

  async fetchRole(roleName) {
    const path = this.selectedAuthPath || this.selectedAuthType;
    const id = JSON.stringify([path, roleName]);
    return this.store.findRecord('role-saml', id, {
      adapterOptions: { namespace: this.namespace },
    });
  },

  cancelLogin(samlWindow, errorMessage) {
    this.closeWindow(samlWindow);
    this.handleSAMLError(errorMessage);
    this.exchangeSAMLTokenPollID.cancelAll();
  },

  closeWindow(samlWindow) {
    this.watchPopup.cancelAll();
    this.watchCurrent.cancelAll();
    samlWindow.close();
  },

  handleSAMLError(err) {
    this.onLoading(false);
    this.onError(err);
  },

  watchPopup: task(function* (samlWindow) {
    while (true) {
      yield timeout(WAIT_TIME);
      if (!samlWindow || samlWindow.closed) {
        this.exchangeSAMLTokenPollID.cancelAll();
        return this.handleSAMLError(ERROR_WINDOW_CLOSED);
      }
    }
  }),

  watchCurrent: task(function* (samlWindow) {
    // when user is about to change pages, close the popup window
    yield waitForEvent(this.getWindow(), 'beforeunload');
    samlWindow?.close();
  }),

  exchangeSAMLTokenPollID: task(function* (samlWindow, role) {
    this.onLoading(true);

    // start watching the popup window and the current one
    this.watchPopup.perform(samlWindow);
    this.watchCurrent.perform(samlWindow);

    const path = this.selectedAuthPath || this.selectedAuthType;
    const adapter = this.store.adapterFor('auth-method');
    this.onNamespace(this.namespace);

    // Wait up to 3 minutes for the token to become available
    let resp;
    for (let i = 0; i < 180; i++) {
      yield timeout(WAIT_TIME);
      try {
        resp = yield adapter.pollSAMLToken(path, role.tokenPollID, role.clientVerifier);
        if (!resp?.auth) {
          continue;
        }
        // We've obtained the Vault token for the authentication flow, now log in.
        yield this.onSubmit(null, null, resp.auth.client_token);
        this.closeWindow(samlWindow);
        return;
      } catch (e) {
        if (e.httpStatus === 401) {
          // Continue to retry on 401 Unauthorized
          continue;
        }
        return this.cancelLogin(samlWindow, errorMessage(e));
      }
    }
    this.cancelLogin(samlWindow, ERROR_TIMEOUT);
  }),

  actions: {
    setRole(roleName) {
      this.onRoleName(roleName);
    },
    /* Saml auth flow on login button click:
     * 1. find role-saml record which returns role info
     * 2. open popup at url defined returned from role
     * 3. watch popup window for close (and cancel polling if it closes)
     * 4. poll vault for 200 token response
     * 5. close popup, stop polling, and trigger onSubmit with token data
     */
    async startSAMLAuth(callback, data, e) {
      this.onError(null);
      this.onLoading(true);
      if (e && e.preventDefault) {
        e.preventDefault();
      }
      const roleName = data.role;
      let role;
      try {
        role = await this.fetchRole(roleName);
      } catch (error) {
        this.handleSAMLError(error);
        return;
      }

      const win = this.getWindow();
      const POPUP_WIDTH = 500;
      const POPUP_HEIGHT = 600;
      const left = win.screen.width / 2 - POPUP_WIDTH / 2;
      const top = win.screen.height / 2 - POPUP_HEIGHT / 2;
      const samlWindow = win.open(
        role.ssoServiceURL,
        'vaultSAMLWindow',
        `width=${POPUP_WIDTH},height=${POPUP_HEIGHT},resizable,scrollbars=yes,top=${top},left=${left}`
      );

      this.exchangeSAMLTokenPollID.perform(samlWindow, role);
    },
  },
});
