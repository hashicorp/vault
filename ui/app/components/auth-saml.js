/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Ember from 'ember';
import { inject as service } from '@ember/service';
import Component from './outer-html';
import { task, timeout, waitForEvent } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { computed } from '@ember/object';
import errorMessage from 'vault/utils/error-message';

const WAIT_TIME = 500;
const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete. Your web browser may have blocked or closed a pop-up window. Please check your settings and click Sign In to try again.';
const ERROR_MISSING_PARAMS =
  'The callback from the provider did not supply all of the required parameters. Please click Sign In to try again. If the problem persists, you may want to contact your administrator.';
export { ERROR_WINDOW_CLOSED, ERROR_MISSING_PARAMS };

export default Component.extend({
  store: service(),
  featureFlagService: service('featureFlag'),

  selectedAuthPath: null,
  selectedAuthType: null,
  roleName: null,
  role: null,
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

  fetchRole: task(
    waitFor(function* (roleName, options = { debounce: true }) {
      if (options.debounce) {
        this.onRoleName(roleName);
        // debounce
        yield timeout(Ember.testing ? 0 : WAIT_TIME);
      }

      const path = this.selectedAuthPath || this.selectedAuthType;
      const id = JSON.stringify([path, roleName]);
      let role = null;
      try {
        role = yield this.store.findRecord('role-saml', id, {
          adapterOptions: { namespace: this.namespace },
        });
      } catch (e) {
        this.set('errorMessage', errorMessage(e));
        return;
      }
      this.set('role', role);
    })
  ).restartable(),

  cancelLogin(samlWindow, errorMessage) {
    this.closeWindow(samlWindow);
    this.handleSAMLError(errorMessage);
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
        return this.handleSAMLError(ERROR_WINDOW_CLOSED);
      }
    }
  }),

  watchCurrent: task(function* (samlWindow) {
    // when user is about to change pages, close the popup window
    yield waitForEvent(this.getWindow(), 'beforeunload');
    samlWindow.close();
  }),

  exchangeSAMLTokenPollID: task(function* (samlWindow) {
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
      sleep(500);

      try {
        resp = yield adapter.pollSAMLToken(path, this.role.tokenPollID, this.role.clientVerifier);
        if (!resp?.auth) {
          continue;
        }

        // We've obtained the Vault token for the authentication flow, now log in.
        this.closeWindow(samlWindow);
        yield this.onSubmit(null, null, resp.auth.client_token);
        return;
      } catch (e) {
        if (e.httpStatus === 401) {
          // Continue to retry on 401 Unauthorized
          continue;
        }
        if (e.httpStatus === 403 || e.httpStatus === 400) {
          return this.cancelLogin(samlWindow, e.errors[0]);
        }
      }
    }
  }),

  actions: {
    async startSAMLAuth(data, e) {
      this.onError(null);
      if (e && e.preventDefault) {
        e.preventDefault();
      }

      try {
        await this.fetchRole.perform(this.roleName, { debounce: false });
      } catch (error) {
        // this task could be cancelled if the instances in didReceiveAttrs resolve after this was started
        if (error?.name !== 'TaskCancelation') {
          throw error;
        }
      }

      if (!this.role) {
        this.onError('Invalid role. Please try again.');
        return;
      }
      if (!this.role.ssoServiceURL) {
        this.onError(
          'Missing sso_service_url. Please check the acs_urls field of the auth method configuration.'
        );
        return;
      }

      const win = this.getWindow();
      const POPUP_WIDTH = 500;
      const POPUP_HEIGHT = 600;
      const left = win.screen.width / 2 - POPUP_WIDTH / 2;
      const top = win.screen.height / 2 - POPUP_HEIGHT / 2;
      const samlWindow = win.open(
        this.role.ssoServiceURL,
        'vaultSAMLWindow',
        `width=${POPUP_WIDTH},height=${POPUP_HEIGHT},resizable,scrollbars=yes,top=${top},left=${left}`
      );

      this.exchangeSAMLTokenPollID.perform(samlWindow);
    },
  },
});

function sleep(milliseconds) {
  const start = new Date().getTime();
  while (new Date().getTime() - start < milliseconds);
}
