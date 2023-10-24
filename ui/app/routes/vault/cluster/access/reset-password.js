import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { encodePath } from 'vault/utils/path-encoding-helpers';

const ERROR_UNAVAILABLE = 'Password reset is not available for your current auth mount.';
const ERROR_NO_ACCESS =
  'You do not have permissions to update your password. If you think this is a mistake ask your administrator to update your policy.';
export default class VaultClusterAccessResetPasswordRoute extends Route {
  @service auth;
  @service store;

  async model() {
    // Password reset is only available on userpass type auth mounts
    if (this.auth.authData?.backend?.type !== 'userpass') {
      throw new Error(ERROR_UNAVAILABLE);
    }
    const { backend, displayName } = this.auth.authData;
    if (!backend.mountPath || !displayName) {
      throw new Error(ERROR_UNAVAILABLE);
    }
    try {
      const capabilities = await this.store.findRecord(
        'capabilities',
        `auth/${encodePath(backend.mountPath)}/users/${encodePath(displayName)}/password`
      );
      // Check that the user has ability to update password
      if (!capabilities.canUpdate) {
        throw new Error(ERROR_NO_ACCESS);
      }
    } catch (e) {
      // If capabilities can't be queried, default to letting the API decide
    }
    return {
      backend: backend.mountPath,
      username: displayName,
    };
  }
}
