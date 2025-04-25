/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { sanitizePath } from 'core/utils/sanitize-path';

import type AuthService from 'vault/vault/services/auth';
import type { AuthData } from 'vault/vault/services/auth';
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
      const value = formData.get(key);
      // strip leading or trailing slashes from path for consistency
      data[key] = key === 'path' ? sanitizePath(value) : value;
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
    waitFor(async (formData) => {
      try {
        const authResponse = await this.auth.authenticate({
          clusterId: this.args.cluster.id,
          backend: this.args.authType,
          data: formData,
          selectedAuth: this.args.authType,
        });

        this.handleAuthResponse(authResponse);
      } catch (error) {
        this.onError(error as Error);
      }
    })
  );

  handleAuthResponse(authResponse: AuthData) {
    // calls onAuthResponse in parent auth/page.js component
    this.args.onSuccess(authResponse);
  }

  onError(error: Error | string) {
    if (!this.auth.mfaErrors) {
      const errorMessage = `Authentication failed: ${this.auth.handleError(error)}`;
      this.args.onError(errorMessage);
    }
  }
}
