/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import AuthBase from './base';
import Ember from 'ember';
import { tracked } from '@glimmer/tracking';
import { task, timeout } from 'ember-concurrency';
import { action } from '@ember/object';
import { waitFor } from '@ember/test-waiters';
import { v4 as uuidv4 } from 'uuid';

/**
 * @module Auth::Form::Okta
 * see Auth::Base
 * */

export default class AuthFormOkta extends AuthBase {
  loginFields = ['username', 'password'];

  @tracked challengeAnswer = '';

  @action
  async login(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = {};

    for (const key of formData.keys()) {
      data[key] = formData.get(key);
    }
    data.nonce = uuidv4();

    try {
      // perform auth
      this.pollForOktaNumberChallenge.perform(data.nonce, data.path);
      // do stuff
    } catch (error) {
      this.onError(error);
    }
  }

  pollForOktaNumberChallenge = task(
    waitFor(async (nonce, mount) => {
      // wait for 1s to wait to see if there is a login error before polling
      await timeout(1000);
      if (this.authError) return;

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
    return this.ajax(url, 'GET', {}).then(
      (resp) => {
        return resp.data.correct_answer;
      },
      (e) => {
        // if error status is 404, return and keep polling for a response
        if (e.status === 404) {
          return null;
        } else {
          throw e;
        }
      }
    );
  }
}
