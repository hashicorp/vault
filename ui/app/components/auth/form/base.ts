/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';

import type AuthService from 'vault/vault/services/auth';
import type ClusterModel from 'vault/models/cluster';
import type { HTMLElementEvent } from 'vault/forms';

/**
 * @module Auth::Base
 *
 * @param {string} authType - chosen login method type
 * @param {object} cluster - The cluster model which contains information such as cluster id, name and boolean for if the cluster is in standby
 * @param {function} onError - callback if there is a login error
 * @param {function} onSuccess - calls onAuthResponse in auth/page redirects if successful
 */

interface Args {
  authType: string;
  cluster: ClusterModel;
  onError: CallableFunction;
  onSuccess: CallableFunction;
}

export default class AuthBase extends Component<Args> {
  @service declare readonly auth: AuthService;

  @action
  onSubmit(event: HTMLElementEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.target as HTMLFormElement);
    const data: Record<string, FormDataEntryValue | null> = {};

    for (const key of formData.keys()) {
      data[key] = formData.get(key);
    }

    // If path is not included in the submitted form data,
    // set it as the auth type which is the default path Vault expects.
    // The "token" auth method does not support custom login paths.
    if (this.args.authType !== 'token' && !Object.keys(data).includes('path')) {
      data['path'] = this.args.authType;
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
        this.onError(error as Error);
      }
    })
  );

  // if we move auth service authSuccess method here (or to each auth method component)
  // then call that before calling parent this.args.onSuccess
  onSuccess(authResponse: object) {
    //  responsible for redirect after auth data is persisted
    this.args.onSuccess(authResponse, this.args.authType);
  }

  onError(error: Error) {
    if (!this.auth.mfaErrors) {
      const errorMessage = `Authentication failed: ${this.auth.handleError(error)}`;
      this.args.onError(errorMessage);
    }
  }
}
