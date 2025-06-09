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

import type AuthService from 'vault/vault/services/auth';
import type ClusterModel from 'vault/models/cluster';
import type FlagsService from 'vault/services/flags';
import type VersionService from 'vault/services/version';
import type { AuthResponse } from 'vault/vault/services/auth';
import type { HTMLElementEvent } from 'vault/forms';
import type { LoginFields } from 'vault/vault/auth/form';
import type { MfaRequirementApiResponse, ParsedMfaRequirement } from 'vault/vault/auth/mfa';

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
        const authResponse = await this.auth.authenticate({
          clusterId: this.args.cluster.id,
          backend: this.args.authType,
          data: formData,
          selectedAuth: this.args.authType,
        });

        const path = formData?.path;
        this.handleAuthResponse(authResponse, path);
      } catch (error) {
        this.onError(error as Error);
      }
    })
  );

  // Standard methods get mfa_requirements from the authenticate method in the auth service
  // methodData is necessary if there's an MfaRequirement because persisting auth data happens after that
  handleAuthResponse(authResponse: AuthResponse | ParsedMfaRequirement, path?: string) {
    const methodData: { selectedAuth: string; path?: string } = { selectedAuth: this.args.authType, path };
    // calls onAuthResponse in parent auth/page.js component
    this.args.onSuccess(authResponse, methodData);
  }

  // SSO methods with a different token exchange workflow skip the auth service authenticate method
  // and need mfa handle separately
  handleMfa(mfaRequirement: MfaRequirementApiResponse, path: string) {
    const parsedMfaAuthResponse = this.auth._parseMfaResponse(mfaRequirement);
    this.handleAuthResponse(parsedMfaAuthResponse, path);
  }

  onError(error: Error | string) {
    if (!this.auth.mfaErrors) {
      const errorMessage = `Authentication failed: ${this.auth.handleError(error)}`;
      this.args.onError(errorMessage);
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
}
