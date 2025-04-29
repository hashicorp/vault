/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import Ember from 'ember';
import { service } from '@ember/service';
import { task, timeout, waitForEvent } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import errorMessage from 'vault/utils/error-message';

import type AdapterError from 'vault/@ember-data/adapter/error';
import type AuthMethodAdapter from 'vault/vault/adapters/auth-method';
import type AuthService from 'vault/vault/services/auth';
import type RoleSamlModel from 'vault/models/role-saml';
import type Store from '@ember-data/store';

/**
 * @module Auth::Form::Saml
 * see Auth::Base
 */

interface SamlLoginData {
  token: string;
  mfa_requirement: object | null;
}

const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete. Your web browser may have blocked or closed a pop-up window. Please check your settings and click "Sign in" to try again.';
const ERROR_TIMEOUT = 'The authentication request has timed out. Please click "Sign in" to try again.';

export { ERROR_WINDOW_CLOSED };

export default class AuthFormSaml extends AuthBase {
  @service declare readonly auth: AuthService;
  @service declare readonly store: Store;

  @tracked fetchedRole: RoleSamlModel | null = null;

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
  login = task(async (formData) => {
    // submit data is parsed by base.ts and a path will always have a value.
    // either the default of auth type, or the custom inputted path
    const { role, path, namespace } = formData;

    await this.startSAMLAuth({ role, path, namespace });
  });

  // Fetch role to get sso_service_url and open popup
  async startSAMLAuth({ role = '', path = '', namespace = '' }) {
    try {
      const id = JSON.stringify([path, role]);
      this.fetchedRole = await this.store.findRecord('role-saml', id, {
        adapterOptions: { namespace },
      });
    } catch (error) {
      this.onError(errorMessage(error));
      return;
    }

    if (this.fetchedRole) {
      const win = window;
      const POPUP_WIDTH = 500;
      const POPUP_HEIGHT = 600;
      const left = win.screen.width / 2 - POPUP_WIDTH / 2;
      const top = win.screen.height / 2 - POPUP_HEIGHT / 2;
      const samlWindow = win.open(
        this.fetchedRole.ssoServiceURL,
        'vaultSAMLWindow',
        `width=${POPUP_WIDTH},height=${POPUP_HEIGHT},resizable,scrollbars=yes,top=${top},left=${left}`
      );

      await this.exchangeSAMLTokenPollID.perform(samlWindow, { path });
    }
  }

  exchangeSAMLTokenPollID = task(async (samlWindow, { path }) => {
    // start watching the popup window and the current one
    this.watchPopup.perform(samlWindow);
    this.watchCurrent.perform(samlWindow);

    // TODO CMB - when wiring up components check if this is still necessary
    // pass namespace from state back to AuthForm
    // this.args.onNamespace(namespace);
    try {
      const authData = await this.pollForToken(samlWindow, { path });
      // We've obtained the Vault token for the authentication flow now log in or pass MFA data

      if (authData?.client_token) {
        const samlExchangeData = { token: authData.client_token, mfa_requirement: authData?.mfa_requirement };
        await this.continueLogin(samlExchangeData);
        this.closeWindow(samlWindow);
      }
    } catch (error) {
      this.cancelLogin(samlWindow, errorMessage(error));
    }
  });

  async pollForToken(samlWindow: Window, { path = '' }) {
    // Poll every one second for the token to become available
    const WAIT_TIME = Ember.testing ? 50 : 1000;
    const MAX_TIME = 180; // 3 minutes in seconds

    const adapter = this.store.adapterFor('auth-method') as AuthMethodAdapter;
    // Wait up to 3 minutes for a token to become available
    for (let attempt = 0; attempt < MAX_TIME; attempt++) {
      await timeout(WAIT_TIME);

      try {
        const resp = await adapter.pollSAMLToken(
          path,
          this.fetchedRole?.tokenPollID,
          this.fetchedRole?.clientVerifier
        );

        if (resp?.auth) {
          // Exit loop if response
          return resp.auth;
        }
      } catch (e) {
        const error = e as AdapterError;
        if (error.httpStatus === 401) {
          // Continue to retry on 401 Unauthorized
          continue;
        }
        throw error;
      }
    }

    this.cancelLogin(samlWindow, ERROR_TIMEOUT);
    return;
  }

  async continueLogin(data: SamlLoginData) {
    try {
      const authResponse = await this.auth.authenticate({
        clusterId: this.args.cluster.id,
        backend: 'token',
        data,
        selectedAuth: this.args.authType,
      });

      // responsible for redirect after auth data is persisted
      this.handleAuthResponse(authResponse, this.args.authType);
    } catch (e) {
      const error = e as AdapterError;
      this.onError(errorMessage(error));
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

  cancelLogin(samlWindow: Window, errorMessage: string) {
    this.closeWindow(samlWindow);
    this.handleSAMLError(errorMessage);
  }

  closeWindow(samlWindow: Window) {
    this.watchPopup.cancelAll();
    this.watchCurrent.cancelAll();
    samlWindow.close();
  }

  handleSAMLError(errorMessage: string) {
    this.exchangeSAMLTokenPollID.cancelAll();
    this.onError(errorMessage);
  }
}
