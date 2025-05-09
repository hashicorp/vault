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
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

import type AdapterError from 'vault/@ember-data/adapter/error';
import type AuthService from 'vault/vault/services/auth';
import type FlagsService from 'vault/services/flags';
import type RoleJwtModel from 'vault/models/role-jwt';
import type Store from '@ember-data/store';
import type { HTMLElementEvent } from 'vault/forms';

/**
 * @module Auth::Form::OidcJwt
 * see Auth::Base
 *
 * OIDC can be configured at 'jwt' or 'oidc', see https://developer.hashicorp.com/vault/docs/auth/jwt
 * we use the same template because displaying the JWT token input depends on the error message
 * returned when fetching :path/oidc/auth_url
 */

interface JwtLoginData {
  namespace?: string;
  path?: string;
  role?: string;
  jwt?: string;
}

interface OidcLoginData {
  token: string;
}

const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete. Your web browser may have blocked or closed a pop-up window. Please check your settings and click Sign In to try again.';
const ERROR_MISSING_PARAMS =
  'The callback from the provider did not supply all of the required parameters.  Please click Sign In to try again. If the problem persists, you may want to contact your administrator.';
const ERROR_JWT_LOGIN = 'OIDC login is not configured for this mount';
export { ERROR_WINDOW_CLOSED, ERROR_MISSING_PARAMS, ERROR_JWT_LOGIN };

export default class AuthFormOidcJwt extends AuthBase {
  @service declare readonly auth: AuthService;
  @service declare readonly flags: FlagsService;
  @service declare readonly store: Store;

  loginFields = [
    {
      name: 'role',
      helperText: 'Vault will use the default role to sign in if this field is left blank.',
    },
  ];

  // set by form inputs
  _formData: FormData = new FormData();

  // set during auth prep and login workflow
  @tracked fetchedRole: RoleJwtModel | null = null;
  @tracked errorMessage = '';
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
  initializeFormData(element: HTMLFormElement) {
    this._formData = new FormData(element);
    this.fetchRole.perform();
  }

  @action
  updateFormData(event: HTMLElementEvent<HTMLInputElement>) {
    const { name, value } = event.target;
    this._formData?.set(name, value);

    // re-fetch role if the following inputs have changed. namespace is not included because
    // when it changes the route model refreshes and a new component instantiates.
    if (['path', 'role'].includes(name)) {
      this.fetchRole.perform(500);
    }
  }

  fetchRole = restartableTask(async (wait = 0) => {
    // task is restartable so if the user starts typing again,
    // it will cancel and restart from the beginning.
    if (wait) await timeout(wait);

    const { namespace = '', path = '', role = '' } = this.parseFormData(this._formData);
    const id = JSON.stringify([path, role]);

    // reset state
    this.fetchedRole = null;
    this.errorMessage = '';

    try {
      this.fetchedRole = await this.store.findRecord('role-jwt', id, {
        adapterOptions: { namespace },
      });
      this.isOIDC = true;
    } catch (e) {
      const { httpStatus } = e as AdapterError;
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
    waitFor(async (submitData) => {
      if (this.isOIDC) {
        this.startOIDCAuth();
      } else {
        this.continueLogin(submitData);
      }
    })
  );

  async continueLogin(data: JwtLoginData | OidcLoginData) {
    try {
      // TODO CMB backend should probably be path, but holding off refactor since api service may remove need all together
      // OIDC callback returns a token so authenticate with that
      const backend = this.isOIDC && 'token' in data ? 'token' : this.args.authType;

      const authResponse = await this.auth.authenticate({
        clusterId: this.args.cluster.id,
        backend,
        data,
        selectedAuth: this.args.authType,
      });

      // responsible for redirect after auth data is persisted
      this.handleAuthResponse(authResponse);
    } catch (error) {
      this.onError(error as Error);
    }
  }

  // * OIDC AUTH PART 1
  // 1. request oidc/auth_url to check for config errors, if none continue
  // 2. open popup window at auth_url
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
        this.fetchedRole?.authUrl,
        'vaultOIDCWindow',
        `width=${POPUP_WIDTH},height=${POPUP_HEIGHT},resizable,scrollbars=yes,top=${top},left=${left}`
      );

      this.prepareForOIDC.perform(oidcWindow);
    }
  }

  // * OIDC AUTH PART 2
  // 3. watch popups for premature closure
  // 4. wait message event from window.postMessage() in oidc-callback route
  prepareForOIDC = task(async (oidcWindow) => {
    // NOTE TO DEVS: Be careful when updating the OIDC flow and ensure the updates
    // work with implicit flow. See issue https://github.com/hashicorp/vault-plugin-auth-jwt/pull/192
    const thisWindow = window;

    // start watching the popup window and the current one
    this.watchPopup.perform(oidcWindow);
    this.watchCurrent.perform(oidcWindow);
    // eslint-disable-next-line no-constant-condition
    while (true) {
      // wait for message posted from oidc callback, see issue https://github.com/hashicorp/vault/issues/12436
      // ensure that postMessage event is from expected source
      const event = (await waitForEvent(thisWindow, 'message')) as unknown as MessageEvent;
      if (event.origin === thisWindow.origin && event.isTrusted && event.data.source === 'oidc-callback') {
        // event.data are params from the oidc callback url parsed by getParamsForCallback in the oidc-callback route
        return this.exchangeOIDC.perform(event.data, oidcWindow);
      }
    }
  });

  // * OIDC AUTH PART 3
  // 5. check parsed url for expected state params
  // 6. if successful, request client_token from oidc/callback
  // 7. close popups and continue login with client_token
  exchangeOIDC = task(async (oidcState, oidcWindow) => {
    if (oidcState === null || oidcState === undefined) {
      return;
    }

    const { path, state, code } = oidcState;
    if (!path || !state || !code) {
      return this.cancelLogin(oidcWindow, ERROR_MISSING_PARAMS);
    }

    let resp;
    // do the OIDC exchange, set the token and continue login flow
    try {
      const adapter = this.store.adapterFor('auth-method');
      resp = await adapter.exchangeOIDC(path, state, code);
      this.closeWindow(oidcWindow);
    } catch (e) {
      // If there was an error on Vault's end, close the popup
      // and show the error on the login screen
      return this.cancelLogin(oidcWindow, errorMessage(e));
    }

    const { client_token, mfa_requirement } = resp.auth;
    if (mfa_requirement) {
      return this.handleMfa(mfa_requirement, path);
    } else if (client_token) {
      return this.continueLogin({ token: client_token });
    } else {
      // If there's a problem with the OIDC exchange the auth workflow should fail earlier.
      // Including this catch just in case, though it's unlikely this will be hit.
      this.handleOIDCError('Missing token. Please try again.');
    }
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

  cancelLogin(oidcWindow: Window, errorMessage: string) {
    this.closeWindow(oidcWindow);
    this.handleOIDCError(errorMessage);
  }

  closeWindow(oidcWindow: Window) {
    this.watchPopup.cancelAll();
    this.watchCurrent.cancelAll();
    oidcWindow.close();
  }

  handleOIDCError(err: string) {
    this.prepareForOIDC.cancelAll();
    this.exchangeOIDC.cancelAll();
    this.onError(err);
  }
}
