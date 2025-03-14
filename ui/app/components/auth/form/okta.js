/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task, timeout } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { v4 as uuidv4 } from 'uuid';
import { waitFor } from '@ember/test-waiters';
import Ember from 'ember';

/**
 * @module Auth::Form::Okta
 * see Auth::Base
 * */

export default class AuthFormOkta extends AuthBase {
  loginFields = ['username', 'password'];

  @service auth;

  @tracked challengeAnswer = null;
  @tracked waiting = false;

  login = task(
    waitFor(async (data) => {
      data.nonce = uuidv4();
      try {
        // wait for 1s to wait to see if there is a login error before polling
        await timeout(1000);
        this.pollForOktaNumberChallenge.perform(data.nonce, data.path);

        const authResponse = await this.auth.authenticate({
          clusterId: this.args.cluster.id,
          backend: this.args.authType,
          data,
          selectedAuth: this.args.authType,
        });

        this.onSuccess(authResponse);
      } catch (error) {
        this.onError(error);
        this.onCancel();
      }
    })
  );

  pollForOktaNumberChallenge = task(
    waitFor(async (nonce, mount) => {
      this.waiting = true;

      // keep polling /auth/okta/verify/:nonce API every 1s until response returns with correct_number
      let response = null;
      while (response === null) {
        // disable polling for tests otherwise promises reject and acceptance tests fail
        if (Ember.testing) return;

        await timeout(1000);
        response = await this.getOktaNumberChallengeAnswer(nonce, mount);
      }

      // display correct number so user can select on personal MFA device
      this.challengeAnswer = response;
    })
  );

  @action
  getOktaNumberChallengeAnswer(nonce, mount) {
    const url = `/v1/auth/${mount}/verify/${nonce}`;
    return this.auth
      .ajax(url, 'GET', {})
      .then((resp) => resp.data.correct_answer)
      .catch((e) => {
        // if error status is 404, return and keep polling for a response
        if (e.status === 404) {
          return null;
        } else {
          throw e;
        }
      });
  }

  @action
  onCancel() {
    // reset variables and stop polling tasks if canceling login
    this.challengeAnswer = null;
    this.waiting = false;
    this.login.cancelAll();
    this.pollForOktaNumberChallenge.cancelAll();
  }
}
