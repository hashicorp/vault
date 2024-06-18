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

/**
 * @module AuthPage
 * The Auth::Page wraps OktaNumberChallenge and AuthForm to manage the login flow and is responsible for calling the authenticate method
 *
 * @example
 * <Auth::Page @wrappedToken={{this.wrappedToken}} @cluster={{this.model}} @namespace={{this.namespaceQueryParam}} @selectedAuth={{this.authMethod}} @onSuccess={{action "onAuthResponse"}} />
 *
 * @param {string} wrappedToken - Query param value of a wrapped token that can be used to login when added directly to the URL via the "wrapped_token" query param
 * @param {object} cluster - The route model which is the ember data cluster model. contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {string} namespace- Namespace query param, passed to AuthForm and set by typing in namespace input or URL
 * @param {string} selectedAuth - The auth method selected in the dropdown, passed to auth service's authenticate method
 * @param {function} onSuccess - Callback that fires the "onAuthResponse" action in the auth controller and handles transitioning after success
 */

export default class AuthPageComponent extends Component {
  @service auth;

  @tracked authError = null;
  @tracked oktaNumberChallengeAnswer = '';
  @tracked waitingForOktaNumberChallenge = false;

  @action
  performAuth(backendType, data) {
    this.authenticate.unlinked().perform(backendType, data);
  }

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
