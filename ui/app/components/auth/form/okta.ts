/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task, timeout } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';
import uuid from 'core/utils/uuid';

import type AuthService from 'vault/vault/services/auth';

/**
 * @module Auth::Form::Okta
 * see Auth::Base
 * */

export default class AuthFormOkta extends AuthBase {
  @service declare readonly auth: AuthService;

  @tracked challengeAnswer = '';
  @tracked oktaVerifyError = '';
  @tracked showNumberChallenge = false;

  loginFields = [{ name: 'username' }, { name: 'password' }];

  login = task(
    waitFor(async (data) => {
      // wait for 1s to wait to see if there is a login error before polling
      await timeout(1000);

      data.nonce = uuid();
      this.pollForOktaNumberChallenge.perform(data.nonce, data.path);

      try {
        // selecting the correct okta verify answer on the personal device resolves this request
        const authResponse = await this.auth.authenticate({
          clusterId: this.args.cluster.id,
          backend: this.args.authType,
          data,
          selectedAuth: this.args.authType,
        });

        this.handleAuthResponse(authResponse);
      } catch (error) {
        // if a user fails the okta verify challenge, the POST login request fails (made by this.auth.authenticate above)
        // bubble those up for consistency instead of managing error state in this component
        this.onError(error as Error);
        // cancel polling tasks and reset state
        this.reset();
      }
    })
  );

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
      this.challengeAnswer = verifyNumber ?? '';
    })
  );

  @action
  async requestOktaVerify(nonce: string, mountPath: string) {
    const url = `/v1/auth/${mountPath}/verify/${nonce}`;
    try {
      const response = await this.auth.ajax(url, 'GET', {});
      return response.data.correct_answer;
    } catch (e) {
      const error = e as Response;
      if (error?.status === 404) {
        // if error status is 404 return null to keep polling for a response
        return null;
      } else {
        // this would be unusual, but handling just in case
        this.oktaVerifyError = errorMessage(e);
        return;
      }
    }
  }

  @action
  reset() {
    // reset tracked variables and stop polling tasks
    this.challengeAnswer = '';
    this.oktaVerifyError = '';
    this.showNumberChallenge = false;
    this.login.cancelAll();
    this.pollForOktaNumberChallenge.cancelAll();
  }
}
