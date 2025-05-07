/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import Ember from 'ember';
import { service } from '@ember/service';
import { restartableTask, task, timeout, waitForEvent } from 'ember-concurrency';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

const ERROR_WINDOW_CLOSED =
  'The provider window was closed before authentication was complete. Your web browser may have blocked or closed a pop-up window. Please check your settings and click Sign In to try again.';
const ERROR_MISSING_PARAMS =
  'The callback from the provider did not supply all of the required parameters.  Please click Sign In to try again. If the problem persists, you may want to contact your administrator.';
const ERROR_JWT_LOGIN = 'OIDC login is not configured for this mount';
export { ERROR_WINDOW_CLOSED, ERROR_MISSING_PARAMS, ERROR_JWT_LOGIN };

export default class AuthOidcJwt extends Component {
  @service store;
  @service flags;

  // cache values to determine whether or not to refire fetchRole task
  _authType;
  _authPath;
  // set by form inputs
  @tracked roleName = null;
  @tracked jwt;
  // set by auth workflow
  @tracked fetchedRole = null;
  @tracked errorMessage = null;
  @tracked isOIDC = true;

  constructor() {
    super(...arguments);
    this._authPath = this.args.selectedAuthPath;
    this._authType = this.args.selectedAuthType;
    this.fetchRole.perform();
  }

  get tasksAreRunning() {
    return this.prepareForOIDC.isRunning || this.exchangeOIDC.isRunning;
  }

  @action
  checkArgUpdate() {
    // if mount path or type changes we need to check again for JWT configuration
    const didChangePath = this._authPath !== this.args.selectedAuthPath;
    const didChangeType = this._authType !== this.args.selectedAuthType;

    if (didChangePath || didChangeType) {
      // path updates as the user types so we need to debounce that event
      const wait = didChangePath ? 500 : 0;
      this.fetchRole.perform(wait);
    }

    // update cached props
    this._authPath = this.args.selectedAuthPath;
    this._authType = this.args.selectedAuthType;
  }

  fetchRole = restartableTask(async (wait) => {
    // task is `restartable` so if the user starts typing again,
    // it will cancel and restart from the beginning.
    if (wait) await timeout(wait);

    // if we have a custom path is inputted use that,
    // otherwise fallback to type (which is the default path)
    const path = this.args.selectedAuthPath || this.args.selectedAuthType;

    const id = JSON.stringify([path, this.roleName]);

    this.fetchedRole = null;
    this.errorMessage = null;
    this.isOIDC = true;

    try {
      this.fetchedRole = await this.store.findRecord('role-jwt', id, {
        adapterOptions: { namespace: this.args.namespace },
      });
    } catch (e) {
      const error = (e.errors || [])[0];
      const errorMessage =
        e.httpStatus === 400 ? 'Invalid role. Please try again.' : `Error fetching role: ${error}`;
      // assume OIDC until it's known that the mount is configured for JWT authentication via static keys, JWKS, or OIDC discovery.
      // if the mount is configured for JWT this specific error is returned.
      this.isOIDC = error !== ERROR_JWT_LOGIN;
      this.errorMessage = errorMessage;
    }
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
    this.args.onError(err);
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
      const event = await waitForEvent(thisWindow, 'message');
      if (event.origin === thisWindow.origin && event.isTrusted && event.data.source === 'oidc-callback') {
        return this.exchangeOIDC.perform(event.data, oidcWindow);
      }
    }
  });

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

  exchangeOIDC = task(async (oidcState, oidcWindow) => {
    if (oidcState === null || oidcState === undefined) {
      return;
    }

    let { namespace, path, state, code } = oidcState;

    // The namespace can be either be passed as a query parameter, or be embedded
    // in the state param in the format `<state_id>,ns=<namespace>`. So if
    // `namespace` is empty, check for namespace in state as well.
    // TODO smoke test HVD flag here and add test
    if (namespace === '' || this.flags.hvdManagedNamespaceRoot) {
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
    // pass namespace from state back to AuthForm
    this.args.onNamespace(namespace);
    let resp;
    // do the OIDC exchange, set the token on the parent component
    // and submit auth form
    try {
      resp = await adapter.exchangeOIDC(path, state, code);
      this.closeWindow(oidcWindow);
    } catch (e) {
      // If there was an error on Vault's end, close the popup
      // and show the error on the login screen
      return this.cancelLogin(oidcWindow, e);
    }
    const { mfa_requirement, client_token } = resp.auth;
    // onSubmit calls doSubmit in auth-form.js
    await this.args.onSubmit({ mfa_requirement }, null, client_token);
  });

  async startOIDCAuth() {
    this.args.onError(null);

    await this.fetchRole.perform();

    const error =
      this.fetchedRole && !this.fetchedRole.authUrl
        ? 'Missing auth_url. Please check that allowed_redirect_uris for the role include this mount path.'
        : this.errorMessage || null;

    if (error) {
      this.args.onError(error);
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

  @action
  onRoleInput(event) {
    this.roleName = event.target.value;
    this.fetchRole.perform(500);
  }

  @action
  signIn(event) {
    event.preventDefault();

    if (this.isOIDC) {
      this.startOIDCAuth();
    } else {
      this.args.onSubmit({ role: this.roleName, jwt: this.jwt });
    }
  }
}
