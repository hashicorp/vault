/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import Ember from 'ember';
import { action } from '@ember/object';
import { dasherize } from '@ember/string';
import {
  DOMAIN_PROVIDER_MAP,
  ERROR_JWT_LOGIN,
  ERROR_MISSING_PARAMS,
  ERROR_POPUP_FAILED,
  ERROR_WINDOW_CLOSED,
} from 'vault/utils/auth-form-helpers';
import { restartableTask, task, timeout, waitForEvent } from 'ember-concurrency';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import parseURL from 'core/utils/parse-url';

import type { HTMLElementEvent } from 'vault/forms';
import type { JwtOidcAuthUrlResponse, JwtOidcLoginApiResponse } from 'vault/vault/auth/methods';
import type RouterService from '@ember/routing/router-service';

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
  path: string;
  role?: string;
  jwt?: string;
}

interface UrlParseData {
  hostname: string;
}

export default class AuthFormOidcJwt extends AuthBase {
  @service declare readonly router: RouterService;

  loginFields = [
    {
      name: 'role',
      helperText: 'Vault will use the default role to sign in if this field is left blank.',
    },
  ];

  // set by form inputs
  _formData: FormData = new FormData();

  // set during auth prep and login workflow
  @tracked authUrl: string | null = null;
  @tracked errorMessage = '';
  @tracked isOIDC = true;

  get icon() {
    // Right now there is a bug in HDS where the name includes a space, this line can be removed when we
    // upgrade to an HDS version with the corrected icon name
    if (this.provider === 'Ping Identity') return 'ping-identity ';
    return this.provider ? dasherize(this.provider.toLowerCase()) : '';
  }

  get providerName() {
    return `with ${this.provider || 'OIDC Provider'}`;
  }

  get provider() {
    const { hostname } = parseURL(this.authUrl) as UrlParseData;
    if (hostname) {
      const firstMatch = Object.keys(DOMAIN_PROVIDER_MAP).find((name) => hostname.includes(name));
      return firstMatch ? DOMAIN_PROVIDER_MAP[firstMatch as keyof typeof DOMAIN_PROVIDER_MAP] : null;
    }
    return null;
  }

  @action
  initializeFormData(element: HTMLFormElement) {
    this._formData = new FormData(element);
    this.fetchAuthUrl.perform();
  }

  @action
  updateFormData(event: HTMLElementEvent<HTMLInputElement>) {
    const { name, value } = event.target;
    // the selectedAuthMethod dropdown is unrelated to login data so no need to track it in the form state.
    if (name === 'selectedAuthMethod') return;
    this._formData?.set(name, value);

    // re-request auth_url if the following inputs have changed. namespace is not included because
    // when it changes the route model refreshes and a new component instantiates.
    if (['path', 'role'].includes(name)) {
      this.fetchAuthUrl.perform(500);
    }
  }

  fetchAuthUrl = restartableTask(async (wait = 0) => {
    // task is restartable so if the user starts typing again,
    // it will cancel and restart from the beginning.
    if (wait) await timeout(wait);

    const { namespace = '', path = '', role = '' } = this.parseFormData(this._formData);
    const redirect_uri = this.generateRedirectUri(namespace, path);

    // reset state
    this.authUrl = null;
    this.errorMessage = '';

    try {
      const { data } = (await this.api.auth.jwtOidcRequestAuthorizationUrl(path, {
        role,
        redirect_uri,
      })) as JwtOidcAuthUrlResponse;
      this.authUrl = data.auth_url;
      this.isOIDC = true;
    } catch (e) {
      const { status, message } = await this.api.parseError(e);
      // errors are tracked here but they only display on submit
      this.errorMessage =
        // A 400 is returned if OIDC is configured but does not have a default role set.
        status === 400 ? 'Invalid role. Please try again.' : `Error fetching role: ${message}`;
      // If the mount is configured for JWT authentication via static keys, JWKS, or OIDC discovery
      // this specific error is returned. Flip the isOIDC boolean accordingly, otherwise assume OIDC.
      this.isOIDC = !message.includes(ERROR_JWT_LOGIN);
    }
  });

  generateRedirectUri(namespace = '', path = '') {
    const origin = window.location.origin;
    const qp = namespace ? { namespace } : {};
    const routeUrl = this.router.urlFor(
      'vault.cluster.oidc-callback',
      { auth_path: path },
      { queryParams: qp }
    );
    return `${origin}${routeUrl}`;
  }

  // * LOGIN WORKFLOW BEGINS
  async loginRequest(formData: JwtLoginData) {
    if (this.isOIDC) {
      return await this.loginOidc();
    } else {
      return await this.loginJwt(formData);
    }
  }

  async loginJwt(formData: JwtLoginData) {
    const { path, jwt, role } = formData;
    const { auth } = (await this.api.auth.jwtLogin(path, { jwt, role })) as JwtOidcLoginApiResponse;
    // displayName is not returned by auth response and is set in persistAuthData
    return this.normalizeAuthResponse(auth, {
      authMountPath: path,
      token: auth.client_token,
      ttl: auth.lease_duration,
    });
  }

  async loginOidc() {
    const oidcWindow = await this.startOIDCAuth();
    if (oidcWindow) {
      try {
        // Initiate watching for the popup and current window
        this.watchPopup.perform(oidcWindow);
        this.watchCurrent.perform(oidcWindow);
        const eventData = await this.prepareForOIDC();
        const { auth, path } = await this.exchangeOIDC(eventData);
        // displayName is not returned by auth response and is set in persistAuthData
        return this.normalizeAuthResponse(auth, {
          authMountPath: path,
          token: auth.client_token,
          ttl: auth.lease_duration,
        });
      } finally {
        this.closeWindow(oidcWindow);
      }
    } else {
      throw `Failed to open OIDC popup window. ${ERROR_POPUP_FAILED}`;
    }
  }

  // * OIDC AUTH PART 1
  // 1. request oidc/auth_url to check for config errors, if none continue
  // 2. open popup window at auth_url
  async startOIDCAuth() {
    await this.fetchAuthUrl.perform();

    if (!this.authUrl) {
      const error =
        // authUrl is an empty string if the request succeeds but a role is not properly configured.
        this.authUrl === ''
          ? 'Missing auth_url. Please check that allowed_redirect_uris for the role include this mount path.'
          : this.errorMessage || 'Unknown OIDC error. Check the Vault logs and try again.';
      throw error;
    }

    const win = window;
    const POPUP_WIDTH = 500;
    const POPUP_HEIGHT = 600;
    const left = win.screen.width / 2 - POPUP_WIDTH / 2;
    const top = win.screen.height / 2 - POPUP_HEIGHT / 2;
    const oidcWindow = win.open(
      this.authUrl,
      'vaultOIDCWindow',
      `width=${POPUP_WIDTH},height=${POPUP_HEIGHT},resizable,scrollbars=yes,top=${top},left=${left}`
    );

    return oidcWindow;
  }

  // * OIDC AUTH PART 2
  // 3. watch popups for premature closure
  // 4. wait for message event from window.postMessage() in oidc-callback route
  async prepareForOIDC() {
    // NOTE TO DEVS: Be careful when updating the OIDC flow and ensure the updates
    // work with implicit flow. See issue https://github.com/hashicorp/vault-plugin-auth-jwt/pull/192
    const thisWindow = window;

    // eslint-disable-next-line no-constant-condition
    while (true) {
      // wait for message posted from oidc callback, see issue https://github.com/hashicorp/vault/issues/12436
      // ensure that postMessage event is from expected source
      const event = (await waitForEvent(thisWindow, 'message')) as unknown as MessageEvent;
      if (event.origin === thisWindow.origin && event.isTrusted && event.data.source === 'oidc-callback') {
        // event.data are params from the oidc callback url parsed by getParamsForCallback in the oidc-callback route
        return event.data;
      }
    }
  }

  // * OIDC AUTH PART 3
  // 5. check parsed url for expected state params
  // 6. if successful, request client_token from oidc/callback
  // 7. close popups and continue login with client_token

  async exchangeOIDC(oidcState: { path: string; state: string; code: string }) {
    const { path, state, code } = oidcState;
    if (!path || !state || !code) {
      throw ERROR_MISSING_PARAMS;
    }

    // do the OIDC exchange, set the token and continue login flow
    const { auth } = (await this.api.auth.jwtOidcCallback(
      path,
      undefined,
      code,
      state
    )) as JwtOidcLoginApiResponse;
    return { auth, path };
  }
  //* END LOGIN METHODS

  // MANAGE POPUPS
  watchPopup = task(async (oidcWindow) => {
    // eslint-disable-next-line no-constant-condition
    while (true) {
      const WAIT_TIME = Ember.testing ? 50 : 500;

      await timeout(WAIT_TIME);
      if (!oidcWindow || oidcWindow.closed) {
        // Since watchPopup isn't awaited, errors thrown here won't bubble up
        // and so we must call onError directly instead.
        this.onError(ERROR_WINDOW_CLOSED);
        return;
      }
    }
  });

  watchCurrent = task(async (oidcWindow) => {
    // when user is about to change pages, close the popup window
    await waitForEvent(window, 'beforeunload');
    oidcWindow.close();
  });

  cancelLogin() {
    this.login.cancelAll();
  }

  closeWindow(oidcWindow: Window) {
    this.watchPopup.cancelAll();
    this.watchCurrent.cancelAll();
    oidcWindow.close();
  }
}
