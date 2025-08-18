/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import { action } from '@ember/object';
import { task, timeout } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import uuid from 'core/utils/uuid';

import type { OktaVerifyApiResponse, UsernameLoginResponse } from 'vault/vault/auth/methods';

/**
 * @module Auth::Form::Okta
 * see Auth::Base
 * */

export default class AuthFormOkta extends AuthBase {
  @tracked challengeAnswer = '';
  @tracked oktaVerifyError = '';
  @tracked showNumberChallenge = false;

  loginFields = [{ name: 'username' }, { name: 'password' }];

  async loginRequest(formData: { path: string; username: string; password: string }) {
    const { path, username, password } = formData;
    // wait for 1s to wait to see if there is a login error before polling
    await timeout(1000);

    const nonce = uuid();
    this.pollForOktaNumberChallenge.perform(nonce, path);

    // If an Okta MFA challenge is configured for the end user this request resolves when it is completed.
    // If a user fails the MFA challenge (e.g. Okta number challenge) this POST login request fails.
    const { auth } = (await this.api.auth.oktaLogin(username, path, {
      nonce,
      password,
    })) as UsernameLoginResponse;

    return this.normalizeAuthResponse(auth, {
      authMountPath: path,
      displayName: auth?.metadata?.username,
      token: auth.client_token,
      ttl: auth.lease_duration,
    });
  }

  pollForOktaNumberChallenge = task(
    waitFor(async (nonce, mountPath) => {
      this.showNumberChallenge = true;

      // keep polling /auth/okta/verify/:nonce API every 1s until response returns with correct_number
      let verifyNumber = null;
      while (verifyNumber === null) {
        await timeout(1000);
        verifyNumber = await this.requestOktaVerify(nonce, mountPath);
      }

      // display correct number so user can select on personal MFA device
      this.challengeAnswer = verifyNumber?.toString() ?? '';
    })
  );

  @action
  async requestOktaVerify(nonce: string, mountPath: string) {
    try {
      const { data } = (await this.api.auth.oktaVerify(nonce, mountPath)) as OktaVerifyApiResponse;
      return data.correct_answer;
    } catch (e) {
      const { status, message } = await this.api.parseError(e);
      if (status === 404) {
        // if error status is 404 return null to keep polling for a response
        return null;
      } else {
        // this would be unusual, but handling just in case
        this.oktaVerifyError = message;
        return;
      }
    }
  }

  @action
  cancelLogin() {
    // reset tracked variables and stop polling tasks
    this.challengeAnswer = '';
    this.oktaVerifyError = '';
    this.showNumberChallenge = false;
    this.pollForOktaNumberChallenge.cancelAll();
    this.login.cancelAll();
  }
}
