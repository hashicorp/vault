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
import { POSSIBLE_FIELDS } from 'vault/utils/supported-login-methods';
import { ResponseError } from '@hashicorp/vault-client-typescript';

import type { HTMLElementEvent } from 'vault/forms';
import type { LoginFields, NormalizedAuthData } from 'vault/vault/auth/form';
import type ApiService from 'vault/services/api';
import type ClusterModel from 'vault/models/cluster';
import type FlagsService from 'vault/services/flags';
import type VersionService from 'vault/services/version';

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

// This an "abstract" class because it is not meant to be instantiated directly and should be extended from by each auth method type.
// If at any point the Vault UI wants to support a dynamic list of login methods (for example, via custom auth plugins) this class can be
// refactored to handle general auth types, but the Vault UI does not currently support this.
export default abstract class AuthBase extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flags: FlagsService;
  @service declare readonly version: VersionService;

  @action
  onSubmit(event: HTMLElementEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.target as HTMLFormElement);
    const data = this.parseFormData(formData);
    this.login.unlinked().perform(data);
  }

  login = task(
    waitFor(async (formData) => {
      try {
        const authResponse = await this.loginRequest(formData);
        this.args.onSuccess(authResponse);
      } catch (error) {
        this.onError(error as ResponseError);
      }
    })
  );

  // This method must be defined by child components and invokes the relevant login method
  abstract loginRequest(formData: Record<string, any>): Promise<NormalizedAuthData>;

  // Optional method for canceling any additional login items that may be relevant the
  // authentication workflow for that method, such as canceling polling tasks or closing popup windows.
  cancelLogin?(): void;

  async onError(error: ResponseError | string) {
    this.login.cancelAll();
    // Cancel any additional login items that may be relevant to that method.
    // For example, polling tasks or popup windows.
    if (this.cancelLogin) {
      this.cancelLogin();
    }
    // If error has not been parsed then parse and render error message
    if (error instanceof ResponseError) {
      const { message } = await this.api.parseError(error);
      this.args.onError(`Authentication failed: ${message}`);
      return;
    }

    if (typeof error === 'string') {
      this.args.onError(`Authentication failed: ${error}`);
    }
  }

  parseFormData(formData: FormData) {
    const data: LoginFields = {};

    // iterate over method specific fields
    for (const field of POSSIBLE_FIELDS) {
      const value = formData.get(field);
      if (value) {
        data[field] = typeof value === 'string' ? value : undefined;
      }
    }

    // path is supported by all auth methods except token
    if (this.args.authType !== 'token') {
      // strip leading or trailing slashes for consistency.
      // fallback to auth type which is the default path Vault expects.
      data['path'] = sanitizePath(formData?.get('path')) || this.args.authType;
    }

    if (this.version.isEnterprise) {
      // strip leading or trailing slashes for consistency
      let namespace = sanitizePath(formData?.get('namespace')) || '';

      const hvdRootNs = this.flags.hvdManagedNamespaceRoot; // if HVD managed, this is "admin"
      if (hvdRootNs) {
        // HVD managed clusters can only input child namespaces, manually prepend with the hvd root
        namespace = namespace ? `${hvdRootNs}/${namespace}` : hvdRootNs;
      }
      data['namespace'] = namespace;
    }

    return data;
  }

  normalizeAuthResponse = (
    authResponse: any,
    { path = '', tokenKey = '', ttlKey = '', displayName = '' }
  ) => {
    return {
      ...authResponse,
      authMethodType: this.args.authType,
      authMountPath: path,
      displayName,
      token: authResponse[tokenKey],
      ttl: authResponse[ttlKey],
    };
  };
}
