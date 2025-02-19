/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import { task, timeout, waitForEvent } from 'ember-concurrency';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

import errorMessage from 'vault/utils/error-message';

const WAIT_TIME = 500;
const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete. Your web browser may have blocked or closed a pop-up window. Please check your settings and click Sign In to try again.';
const ERROR_TIMEOUT = 'The authentication request has timed out. Please click Sign In to try again.';
const ERROR_MISSING_PARAMS =
  'The callback from the provider did not supply all of the required parameters. Please click Sign In to try again. If the problem persists, you may want to contact your administrator.';
export { ERROR_WINDOW_CLOSED, ERROR_MISSING_PARAMS };

export default class AuthSaml extends Component {
  @service store;
  @service flags;

  @tracked errorMessage;

  getWindow() {
    return this.window || window;
  }

  get canLoginSaml() {
    return this.getWindow().isSecureContext;
  }

  async fetchRole(roleName) {
    const path = this.args.selectedAuthPath || this.args.selectedAuthType;
    const id = JSON.stringify([path, roleName]);
    return this.store.findRecord('role-saml', id, {
      adapterOptions: { namespace: this.args.namespace },
    });
  }

  cancelLogin(samlWindow, errorMessage) {
    this.closeWindow(samlWindow);
    this.handleSAMLError(errorMessage);
    this.exchangeSAMLTokenPollID.cancelAll();
  }

  closeWindow(samlWindow) {
    this.watchPopup.cancelAll();
    this.watchCurrent.cancelAll();
    samlWindow.close();
  }

  handleSAMLError(err) {
    this.args.onLoading(false);
    this.args.onError(err);
  }

  @task
  *watchPopup(samlWindow) {
    while (true) {
      yield timeout(WAIT_TIME);
      if (!samlWindow || samlWindow.closed) {
        this.exchangeSAMLTokenPollID.cancelAll();
        return this.handleSAMLError(ERROR_WINDOW_CLOSED);
      }
    }
  }

  @task
  *watchCurrent(samlWindow) {
    // when user is about to change pages, close the popup window
    yield waitForEvent(this.getWindow(), 'beforeunload');
    samlWindow?.close();
  }

  @task
  *exchangeSAMLTokenPollID(samlWindow, role) {
    this.args.onLoading(true);

    // start watching the popup window and the current one
    this.watchPopup.perform(samlWindow);
    this.watchCurrent.perform(samlWindow);

    const path = this.args.selectedAuthPath || this.args.selectedAuthType;
    const adapter = this.store.adapterFor('auth-method');
    this.args.onNamespace(this.args.namespace);

    // Wait up to 3 minutes for the token to become available
    let resp;
    for (let i = 0; i < 180; i++) {
      yield timeout(WAIT_TIME);
      try {
        resp = yield adapter.pollSAMLToken(path, role.tokenPollID, role.clientVerifier);
        if (!resp?.auth) {
          continue;
        }
        // We've obtained the Vault token for the authentication flow now log in or pass MFA data
        const { mfa_requirement, client_token } = resp.auth;
        // onSubmit calls doSubmit in auth-form.js
        yield this.args.onSubmit({ mfa_requirement }, null, client_token);
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
  }

  @action
  setRole(event) {
    this.args.onRoleName(event.target.value);
  }
  /* Saml auth flow on login button click:
   * 1. find role-saml record which returns role info
   * 2. open popup at url defined returned from role
   * 3. watch popup window for close (and cancel polling if it closes)
   * 4. poll vault for 200 token response
   * 5. close popup, stop polling, and trigger onSubmit with token data
   */
  @action async startSAMLAuth(callback, data, e) {
    this.args.onError(null);
    this.args.onLoading(true);
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
  }
}
