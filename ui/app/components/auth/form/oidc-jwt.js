/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import Ember from 'ember';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { restartableTask, task, timeout, waitForEvent } from 'ember-concurrency';
import { action } from '@ember/object';
import { sanitizePath } from 'core/utils/sanitize-path';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

/**
 * @module Auth::Form::OidcJwt
 * see Auth::Base
 *
 * OIDC can be configured at 'jwt' or 'oidc', see https://developer.hashicorp.com/vault/docs/auth/jwt
 * we use the same template because displaying the JWT token input depends on the error message returned when fetching
 * the role
 */

const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete. Your web browser may have blocked or closed a pop-up window. Please check your settings and click Sign In to try again.';
const ERROR_MISSING_PARAMS =
  'The callback from the provider did not supply all of the required parameters.  Please click Sign In to try again. If the problem persists, you may want to contact your administrator.';
const ERROR_JWT_LOGIN = 'OIDC login is not configured for this mount';
export { ERROR_WINDOW_CLOSED, ERROR_MISSING_PARAMS, ERROR_JWT_LOGIN };

export default class AuthFormOidcJwt extends AuthBase {
  @service auth;
  @service flags;
  @service store;

  // set by form inputs
  @tracked formData = null;

  // set by auth workflow
  @tracked fetchedRole = null;
  @tracked errorMessage = null;
  @tracked isOIDC = true;

  get tasksAreRunning() {
    return this.prepareForOIDC.isRunning || this.exchangeOIDC.isRunning;
  }

  get icon() {
    return this?.fetchedRole?.providerIcon || '';
  }

  get providerName() {
    return `with ${this?.fetchedRole?.providerName || 'OIDC Provider'}`;
  }

  @action
  initializeFormData(element) {
    this.formData = new FormData(element);
    this.fetchRole.perform();
  }

  @action
  updateFormData(event) {
    const { name, value } = event.target;
    this.formData.set(name, value);

    // only fetch role if the following inputs have changed
    if (['path', 'role', 'namespace'].includes(name)) {
      this.fetchRole.perform(500);
    }
  }

  fetchRole = restartableTask(async (wait) => {
    // task is restartable so if the user starts typing again,
    // it will cancel and restart from the beginning.
    if (wait) await timeout(wait);

    const namespace = this.formData.get('namespace') || '';
    const path = sanitizePath(this.formData.get('path')) || this.args.authType;
    const role = this.formData.get('role') || '';
    const id = JSON.stringify([path, role]);

    // reset state
    this.fetchedRole = null;
    this.errorMessage = null;

    try {
      this.fetchedRole = await this.store.findRecord('role-jwt', id, {
        adapterOptions: { namespace },
      });
      this.isOIDC = true;
    } catch (e) {
      const { httpStatus } = e;
      const message = errorMessage(e);
      // track errors but they only display on submit
      this.errorMessage =
        httpStatus === 400 ? 'Invalid role. Please try again.' : `Error fetching role: ${message}`;
      // if the mount is configured for JWT authentication via static keys, JWKS, or OIDC discovery
      // this specific error is returned. Flip the isOIDC boolean accordingly, otherwise assume OIDC.
      this.isOIDC = message !== ERROR_JWT_LOGIN;
    }
  });

  login = task(
    waitFor(async (formData) => {
      if (this.isOIDC) {
        this.startOIDCAuth();
      } else {
        this.continueLogin(formData);
      }
    })
  );

  async continueLogin(data) {
    // if (data?.mfa_requirement) {
    //   // calls onAuthResponse in parent auth/page.js component
    //   this.handleAuthResponse(data, this.args.authType);
    //   // return here because mfa-form.js will finish login/authentication flow after mfa validation
    //   return;
    // }

    try {
      // OIDC callback returns a token so authenticate with that
      const backend = this.isOIDC && data?.token ? 'token' : this.args.authType;

      const authResponse = await this.auth.authenticate({
        clusterId: this.args.cluster.id,
        backend,
        data,
        selectedAuth: this.args.authType,
      });
      // responsible for redirect after auth data is persisted
      this.handleAuthResponse(authResponse);
    } catch (error) {
      this.onError(error);
    }
  }

  //* OIDC AUTH BEGINS
  async startOIDCAuth() {
    await this.fetchRole.perform();

    const error =
      this.fetchedRole && !this.fetchedRole.authUrl
        ? 'Missing auth_url. Please check that allowed_redirect_uris for the role include this mount path.'
        : this.errorMessage || null;

    if (error) {
      this.onError(error);
    } else {
      const win = window;
      const POPUP_WIDTH = 500;
      const POPUP_HEIGHT = 600;
      const left = win.screen.width / 2 - POPUP_WIDTH / 2;
      const top = win.screen.height / 2 - POPUP_HEIGHT / 2;
      const oidcWindow = win.open(
        this.fetchedRole.authUrl,
        'vaultOIDCWindow',
        `width=${POPUP_WIDTH},height=${POPUP_HEIGHT},resizable,scrollbars=yes,top=${top},left=${left}`
      );

      this.prepareForOIDC.perform(oidcWindow);
    }
  }

  // NOTE TO DEVS: Be careful when updating the OIDC flow and ensure the updates
  // work with implicit flow. See issue https://github.com/hashicorp/vault-plugin-auth-jwt/pull/192
  prepareForOIDC = task(async (oidcWindow) => {
    const thisWindow = window;

    // start watching the popup window and the current one
    this.watchPopup.perform(oidcWindow);
    this.watchCurrent.perform(oidcWindow);
    // wait for message posted from oidc callback
    // see issue https://github.com/hashicorp/vault/issues/12436
    // ensure that postMessage event is from expected source
    // eslint-disable-next-line no-constant-condition
    while (true) {
      // the oidc-callback url is parsed by getParamsForCallback in the oidc-callback route
      // and the params are returned as event.data here
      const event = await waitForEvent(thisWindow, 'message');
      if (event.origin === thisWindow.origin && event.isTrusted && event.data.source === 'oidc-callback') {
        return this.exchangeOIDC.perform(event.data, oidcWindow);
      }
    }
  });

  exchangeOIDC = task(async (oidcState, oidcWindow) => {
    if (oidcState === null || oidcState === undefined) {
      return;
    }

    const { path, state, code } = oidcState;

    if (!path || !state || !code) {
      return this.cancelLogin(oidcWindow, ERROR_MISSING_PARAMS);
    }
    const adapter = this.store.adapterFor('auth-method');

    let resp;
    // do the OIDC exchange, set the token and continue login flow
    try {
      resp = await adapter.exchangeOIDC(path, state, code);
      this.closeWindow(oidcWindow);
    } catch (e) {
      // If there was an error on Vault's end, close the popup
      // and show the error on the login screen
      return this.cancelLogin(oidcWindow, e);
    }

    const { client_token, mfa_requirement } = resp.auth;
    const callbackData = { token: client_token, mfa_requirement };
    await this.continueLogin(callbackData);
  });

  // MANAGE POPUPS
  watchPopup = task(async (oidcWindow) => {
    // eslint-disable-next-line no-constant-condition
    while (true) {
      const WAIT_TIME = Ember.testing ? 50 : 500;

      await timeout(WAIT_TIME);
      if (!oidcWindow || oidcWindow.closed) {
        return this.handleOIDCError(ERROR_WINDOW_CLOSED);
      }
    }
  });

  watchCurrent = task(async (oidcWindow) => {
    // when user is about to change pages, close the popup window
    await waitForEvent(window, 'beforeunload');
    oidcWindow.close();
  });

  cancelLogin(oidcWindow, errorMessage) {
    this.closeWindow(oidcWindow);
    this.handleOIDCError(errorMessage);
  }

  closeWindow(oidcWindow) {
    this.watchPopup.cancelAll();
    this.watchCurrent.cancelAll();
    oidcWindow.close();
  }

  handleOIDCError(err) {
    this.prepareForOIDC.cancelAll();
    this.onError(err);
  }
}
