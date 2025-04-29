/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import Ember from 'ember';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task, timeout, waitForEvent } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

/**
 * @module Auth::Form::Saml
 * see Auth::Base
 */

const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete. Your web browser may have blocked or closed a pop-up window. Please check your settings and click "Sign in" to try again.';
const ERROR_TIMEOUT = 'The authentication request has timed out. Please click "Sign in" to try again.';

export { ERROR_WINDOW_CLOSED };

export default class AuthFormSaml extends AuthBase {
  @service store;

  @tracked errorMessage;
  @tracked fetchedRole;

  loginFields = [
    {
      name: 'role',
      helperText: 'Vault will use the default role to sign in if this field is left blank.',
    },
  ];

  get canLoginSaml() {
    return window.isSecureContext;
  }

  get tasksAreRunning() {
    return this.login.isRunning || this.exchangeSAMLTokenPollID.isRunning;
  }

  /* Saml auth flow on login button click:
   * 1. find role-saml record which returns role info
   * 2. open popup at url defined returned from role
   * 3. watch popup window for close (and cancel polling if it closes)
   * 4. poll vault for 200 token response
   * 5. close popup, stop polling, and trigger onSubmit with token data
   */
  login = task(async (submitData) => {
    const { role, path, namespace } = submitData;
    try {
      const id = JSON.stringify([path, role]);
      this.fetchedRole = await this.store.findRecord('role-saml', id, {
        adapterOptions: { namespace },
      });
    } catch (error) {
      this.onError(errorMessage(error));
      return;
    }

    // if role is successfully fetched start SAML auth
    const samlWindow = await this.startSAMLAuth();
    await this.exchangeSAMLTokenPollID.perform(samlWindow, submitData);
  });

  @action
  async startSAMLAuth() {
    const win = window;
    const POPUP_WIDTH = 500;
    const POPUP_HEIGHT = 600;
    const left = win.screen.width / 2 - POPUP_WIDTH / 2;
    const top = win.screen.height / 2 - POPUP_HEIGHT / 2;
    return win.open(
      this.fetchedRole.ssoServiceURL,
      'vaultSAMLWindow',
      `width=${POPUP_WIDTH},height=${POPUP_HEIGHT},resizable,scrollbars=yes,top=${top},left=${left}`
    );
  }

  exchangeSAMLTokenPollID = task(async (samlWindow, submitData) => {
    // start watching the popup window and the current one
    this.watchPopup.perform(samlWindow);
    this.watchCurrent.perform(samlWindow);

    const { path } = submitData;

    // TODO CMB - when wiring up components check if this is still necessary
    // pass namespace from state back to AuthForm
    // this.args.onNamespace(namespace);

    const adapter = this.store.adapterFor('auth-method');
    // Wait up to 3 minutes (180 seconds) for the token to become available
    for (let i = 0; i < 180; i++) {
      const WAIT_TIME = Ember.testing ? 50 : 1000;
      await timeout(WAIT_TIME);

      try {
        const resp = await adapter.pollSAMLToken(
          path,
          this.fetchedRole.tokenPollID,
          this.fetchedRole.clientVerifier
        );

        if (resp?.auth) {
          // We've obtained the Vault token for the authentication flow now log in or pass MFA data
          const { mfa_requirement, client_token } = resp.auth;
          // onSubmit calls doSubmit in auth-form.js
          const samlExchangeData = { token: client_token, mfa_requirement };
          await this.continueLogin(samlExchangeData);
          this.closeWindow(samlWindow);
          return;
        } else {
          continue;
        }
      } catch (e) {
        if (e.httpStatus === 401) {
          // Continue to retry on 401 Unauthorized
          continue;
        }
        return this.cancelLogin(samlWindow, errorMessage(e));
      }
    }
    this.cancelLogin(samlWindow, ERROR_TIMEOUT);
  });

  async continueLogin(data) {
    try {
      const authResponse = await this.auth.authenticate({
        clusterId: this.args.cluster.id,
        backend: 'token',
        data,
        selectedAuth: this.args.authType,
      });

      // responsible for redirect after auth data is persisted
      this.handleAuthResponse(authResponse, this.args.authType);
    } catch (error) {
      this.onError(error);
    }
  }

  // MANAGE POPUPS
  watchPopup = task(async (samlWindow) => {
    // eslint-disable-next-line no-constant-condition
    while (true) {
      const WAIT_TIME = Ember.testing ? 50 : 500;

      await timeout(WAIT_TIME);
      if (!samlWindow || samlWindow.closed) {
        return this.handleSAMLError(ERROR_WINDOW_CLOSED);
      }
    }
  });

  watchCurrent = task(async (samlWindow) => {
    // when user is about to change pages, close the popup window
    await waitForEvent(window, 'beforeunload');
    samlWindow?.close();
  });

  cancelLogin(samlWindow, errorMessage) {
    this.closeWindow(samlWindow);
    this.handleSAMLError(errorMessage);
  }

  closeWindow(samlWindow) {
    this.watchPopup.cancelAll();
    this.watchCurrent.cancelAll();
    samlWindow.close();
  }

  handleSAMLError(err) {
    this.exchangeSAMLTokenPollID.cancelAll();
    this.onError(err);
  }
}
