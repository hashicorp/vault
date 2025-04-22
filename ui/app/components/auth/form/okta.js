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

/**
 * @module Auth::Form::Okta
 * see Auth::Base
 * */

export default class AuthFormOkta extends AuthBase {
  loginFields = [{ name: 'username' }, { name: 'password' }];

  @service auth;

  @tracked challengeAnswer = null;
  @tracked oktaVerifyError = '';
  @tracked showNumberChallenge = false;

  login = task(
    waitFor(async (data) => {
      // wait for 1s to wait to see if there is a login error before polling
      await timeout(1000);

      data.nonce = uuid();
      this.pollForOktaNumberChallenge.perform(data.nonce, data.path);

      try {
        const authResponse = await this.auth.authenticate({
          clusterId: this.args.cluster.id,
          backend: this.args.authType,
          data,
          selectedAuth: this.args.authType,
        });

        this.onSuccess(authResponse);
      } catch (error) {
        // if a user fails the okta verify challenge, the POST login request fails (made by this.auth.authenticate above)
        // bubble those up for consistency instead of managing error state in this component
        this.onError(error);
        // cancel polling tasks and reset state
        this.reset();
      }
    })
  );

  pollForOktaNumberChallenge = task(
    waitFor(async (nonce, mountPath) => {
      this.showNumberChallenge = true;

      // keep polling /auth/okta/verify/:nonce API every 1s until response returns with correct_number
      let response = null;
      while (response === null) {
        await timeout(1000);
        response = await this.requestOktaVerify(nonce, mountPath);
      }

      // display correct number so user can select on personal MFA device
      this.challengeAnswer = response;
    })
  );

  @action
  requestOktaVerify(nonce, mountPath) {
    const url = `/v1/auth/${mountPath}/verify/${nonce}`;
    return this.auth
      .ajax(url, 'GET', {})
      .then((resp) => resp.data.correct_answer)
      .catch((e) => {
        // if error status is 404 keep polling for a response
        if (e.status === 404) {
          return null;
        } else {
          // this would be unusual, but handling just in case
          this.oktaVerifyError = errorMessage(e);
        }
      });
  }

  @action
  reset() {
    // reset tracked variables and stop polling tasks
    this.challengeAnswer = null;
    this.oktaVerifyError = '';
    this.showNumberChallenge = false;
    this.login.cancelAll();
    this.pollForOktaNumberChallenge.cancelAll();
  }
}
