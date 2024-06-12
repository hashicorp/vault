/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import Ember from 'ember';
import { service } from '@ember/service';
import { task, timeout } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class AuthPageComponent extends Component {
  @service auth;

  @tracked authError = null;
  @tracked oktaNumberChallengeAnswer = '';
  @tracked waitingForOktaNumberChallenge = false;

  @task
  @waitFor
  *delayAuthMessageReminder() {
    if (Ember.testing) {
      yield timeout(0);
    } else {
      yield timeout(5000);
    }
  }

  @task
  @waitFor
  *authenticate(backendType, data) {
    const {
      selectedAuth,
      cluster: { id: clusterId },
    } = this.args;
    try {
      if (backendType === 'okta') {
        this.pollForOktaNumberChallenge.perform(data.nonce, data.path);
      } else {
        this.delayAuthMessageReminder.perform();
      }
      const authResponse = yield this.auth.authenticate({
        clusterId,
        backend: backendType,
        data,
        selectedAuth,
      });
      this.args.onSuccess(authResponse, backendType, data);
    } catch (e) {
      if (!this.auth.mfaError) {
        this.authError = `Authentication failed: ${this.auth.handleError(e)}`;
      }
    }
  }

  @task
  @waitFor
  *pollForOktaNumberChallenge(nonce, mount) {
    // yield for 1s to wait to see if there is a login error before polling
    yield timeout(1000);
    if (this.authError) return;

    let response = null;
    this.waitingForOktaNumberChallenge = true;
    // this.setCancellingAuth(false);
    // keep polling /auth/okta/verify/:nonce API every 1s until a response is given with the correct number for the Okta Number Challenge
    while (response === null) {
      // disable polling for tests otherwise promises reject and acceptance tests fail
      if (Ember.testing) return;

      yield timeout(1000);
      response = yield this.auth.getOktaNumberChallengeAnswer(nonce, mount);
    }
    this.oktaNumberChallengeAnswer = response;
  }

  @action
  onCancel() {
    // if we are cancelling the login then we reset the number challenge answer and cancel the current authenticate and polling tasks
    this.oktaNumberChallengeAnswer = null;
    this.waitingForOktaNumberChallenge = false;
    this.authenticate.cancelAll();
    this.pollForOktaNumberChallenge.cancelAll();
  }
}
