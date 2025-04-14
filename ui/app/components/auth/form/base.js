/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

/**
 * @module Auth::Base
 *
 * @param {string} authType - chosen login method type
 * @param {object} cluster - The cluster model which contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {function} onError - callback if there is a login error
 * @param {function} onSuccess - calls onAuthResponse in auth/page redirects if successful
 */

export default class AuthBase extends Component {
  @service auth;

  @action
  onSubmit(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const data = {};

    for (const key of formData.keys()) {
      data[key] = formData.get(key);
    }

    this.login.unlinked().perform(data);
  }

  login = task(
    waitFor(async (data) => {
      try {
        const authResponse = await this.auth.authenticate({
          clusterId: this.args.cluster.id,
          backend: this.args.authType,
          data,
          selectedAuth: this.args.authType,
        });

        // responsible for redirect after auth data is persisted
        this.onSuccess(authResponse);
      } catch (error) {
        this.onError(error);
      }
    })
  );

  // if we move auth service authSuccess method here (or to each auth method component)
  // then call that before calling parent this.args.onSuccess
  onSuccess(authResponse) {
    //  responsible for redirect after auth data is persisted
    this.args.onSuccess(authResponse, this.args.authType);
  }

  onError(error) {
    if (!this.auth.mfaError) {
      const errorMessage = `Authentication failed: ${this.auth.handleError(error)}`;
      this.args.onError(errorMessage);
    }
  }
}
