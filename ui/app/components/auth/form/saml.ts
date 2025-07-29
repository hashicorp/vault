/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import Ember from 'ember';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import { SamlWriteSsoServiceUrlRequestClientTypeEnum } from '@hashicorp/vault-client-typescript';
import { sanitizePath } from 'core/utils/sanitize-path';
import { task, timeout, waitForEvent } from 'ember-concurrency';
import uuid from 'core/utils/uuid';
import { ERROR_POPUP_FAILED, ERROR_TIMEOUT, ERROR_WINDOW_CLOSED } from 'vault/utils/auth-form-helpers';

import type { SamlLoginApiResponse, SamlSsoServiceUrlApiResponse } from 'vault/vault/auth/methods';

/**
 * @module Auth::Form::Saml
 * see Auth::Base
 */

interface SamlRole {
  sso_service_url: string;
  token_poll_id: string;
  client_verifier: string;
}
export default class AuthFormSaml extends AuthBase {
  loginFields = [
    {
      name: 'role',
      helperText: 'Vault will use the default role to sign in if this field is left blank.',
    },
  ];

  get canLoginSaml() {
    return window.isSecureContext;
  }

  /* Saml auth flow on login button click:
   * 1. find role-saml record which returns role info
   * 2. open popup at url defined returned from role
   * 3. watch popup window for close (and cancel polling if it closes)
   * 4. poll vault for 200 token response
   * 5. close popup, stop polling, and trigger onSubmit with token data
   */
  async loginRequest(formData: { namespace: string; path: string; role: string }) {
    // submit data is parsed by base.ts and a path will always have a value.
    // either the default of auth type, or the custom inputted path
    const { namespace, path, role } = formData;
    const fetchedRole = await this.fetchSamlRole({ namespace, path, role });
    const samlWindow = await this.startSAMLAuth(fetchedRole.sso_service_url);
    if (samlWindow) {
      try {
        // start watching the popup window and the current one
        this.watchPopup.perform(samlWindow);
        this.watchCurrent.perform(samlWindow);

        const { auth } = await this.exchangeSAMLTokenPollID(fetchedRole, { path });

        // displayName is not included in auth response - it is set in persistAuthData
        return this.normalizeAuthResponse(auth, {
          authMountPath: path,
          token: auth.client_token,
          ttl: auth.lease_duration,
        });
      } finally {
        this.closeWindow(samlWindow);
      }
    } else {
      throw `Failed to open SAML popup window. ${ERROR_POPUP_FAILED}`;
    }
  }

  cancelLogin() {
    this.login.cancelAll();
  }

  // Fetch role to get sso_service_url which is where popup is opened
  async fetchSamlRole({ namespace = '', path = '', role = '' }): Promise<SamlRole> {
    // Create the client verifier and challenge
    const verifier = uuid();
    const client_challenge = await this.generateClientChallenge(verifier);
    const acs_url = this.generateAcsUrl(path, namespace);
    const client_type = SamlWriteSsoServiceUrlRequestClientTypeEnum.BROWSER; // 'browser'
    // Kick off the authentication flow by generating the SSO service URL
    // It requires the client challenge generated from the verifier. We'll
    // later provide the verifier to match up with the challenge on the server
    // when we poll for the Vault token by its returned token poll ID.
    const { data } = (await this.api.auth.samlWriteSsoServiceUrl(path, {
      acs_url,
      client_challenge,
      client_type,
      role,
    })) as SamlSsoServiceUrlApiResponse;
    return {
      ...data,
      client_verifier: verifier,
    };
  }

  async startSAMLAuth(ssoServiceUrl: string): Promise<Window | null> {
    const win = window;
    const POPUP_WIDTH = 500;
    const POPUP_HEIGHT = 600;
    const left = win.screen.width / 2 - POPUP_WIDTH / 2;
    const top = win.screen.height / 2 - POPUP_HEIGHT / 2;
    const samlWindow = win.open(
      ssoServiceUrl,
      'vaultSAMLWindow',
      `width=${POPUP_WIDTH},height=${POPUP_HEIGHT},resizable,scrollbars=yes,top=${top},left=${left}`
    );

    return samlWindow;
  }

  async exchangeSAMLTokenPollID(fetchedRole: SamlRole, { path = '' }) {
    const WAIT_TIME = Ember.testing ? 50 : 1000;
    const MAX_TRIES = Ember.testing ? 3 : 180; // 180 is 3 minutes in seconds

    // Wait up to 3 minutes for a token to become available
    for (let attempt = 0; attempt < MAX_TRIES; attempt++) {
      // Poll every one second for the token to become available
      await timeout(WAIT_TIME);

      try {
        const { client_verifier, token_poll_id } = fetchedRole;
        // Exit loop if there's a response
        return (await this.api.auth.samlWriteToken(path, {
          client_verifier,
          token_poll_id,
        })) as SamlLoginApiResponse;
      } catch (e) {
        const { message, status } = await this.api.parseError(e);
        if (status === 401) {
          // Continue to retry on 401 Unauthorized
          continue;
        }
        // Just throw the message string because parent onError method will fail if it attempts to re-parse an error.
        throw message;
      }
    }

    throw ERROR_TIMEOUT;
  }

  // MANAGE POPUPS
  watchPopup = task(async (samlWindow) => {
    // eslint-disable-next-line no-constant-condition
    while (true) {
      const WAIT_TIME = Ember.testing ? 50 : 500;
      await timeout(WAIT_TIME);

      if (!samlWindow || samlWindow.closed) {
        // Since watchPopup isn't awaited, errors thrown here won't bubble up
        // and so we must call onError directly instead.
        this.onError(ERROR_WINDOW_CLOSED);
        return;
      }
    }
  });

  watchCurrent = task(async (samlWindow) => {
    // when user is about to change pages, close the popup window
    await waitForEvent(window, 'beforeunload');
    samlWindow?.close();
  });

  closeWindow(samlWindow: Window) {
    this.watchPopup.cancelAll();
    this.watchCurrent.cancelAll();
    samlWindow.close();
  }

  // generates a client challenge from a verifier for PKCE (Proof Key for Code Exchange).
  // The client challenge is the base64(sha256(verifier)). The verifier is
  // later presented to the server to obtain the resulting Vault token.
  async generateClientChallenge(verifier: string) {
    const encoder = new TextEncoder();
    const data = encoder.encode(verifier);
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
    const hashArray = new Uint8Array(hashBuffer);
    return btoa(String.fromCharCode(...hashArray));
  }

  generateAcsUrl(path: string, namespace: string) {
    const baseUrl = `${window.location.origin}/v1`;
    const ns = namespace ? `${encodePath(sanitizePath(namespace))}/` : '';
    const mountPath = encodePath(sanitizePath(path));
    // example with "admin" namespace: '${VAULT_ADDR}/v1/admin/auth/saml/callback';
    // example with "root" namespace: '${VAULT_ADDR}/v1/auth/saml/callback';
    return `${baseUrl}/${ns}auth/${mountPath}/callback`;
  }
}
