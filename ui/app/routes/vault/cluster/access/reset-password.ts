/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

import type AuthService from 'vault/vault/services/auth';
import type CapabilitiesService from 'vault/services/capabilities';

const ERROR_UNAVAILABLE = 'Password reset is not available for the current user.';
const ERROR_NO_ACCESS =
  'You do not have permissions to update your password. If you think this is a mistake ask your administrator to update your policy.';

export default class VaultClusterAccessResetPasswordRoute extends Route {
  @service declare readonly auth: AuthService;
  @service declare readonly capabilities: CapabilitiesService;

  async model() {
    const { authMethodType, authMountPath, displayName } = this.auth.authData;
    // Password reset is only available on userpass type auth mounts
    if (authMethodType !== 'userpass') {
      throw new Error(ERROR_UNAVAILABLE);
    }

    // Both of these are necessary to build the reset password URL
    if (!authMountPath || !displayName) {
      throw new Error(ERROR_UNAVAILABLE);
    }

    const capabilities = await this.capabilities.fetchPathCapabilities(
      `auth/${authMountPath}/users/${displayName}/password`
    );

    // Throw an error if we know for certain the user doesn't have permission
    if (!capabilities.canUpdate) {
      throw new Error(ERROR_NO_ACCESS);
    }
    return {
      backend: authMountPath,
      username: displayName,
    };
  }
}
