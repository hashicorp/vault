import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default class VaultClusterAccessResetPasswordRoute extends Route {
  @service auth;
  @service store;

  async model() {
    // Password reset is only available on userpass type auth mounts
    if (this.auth.authData?.backend?.type !== 'userpass') {
      throw new Error('Password reset is not available for your current auth mount.');
    }
    const { backend, displayName } = this.auth.authData;
    const capabilities = await this.store.findRecord(
      'capabilities',
      `auth/${encodePath(backend.mountPath)}/users/${encodePath(displayName)}/password`
    );
    // Check that the user has ability to update password
    if (!capabilities.canUpdate) {
      throw new Error('NO_UPDATE_ACCESS');
    }
    return {
      backend: backend.mountPath,
      username: displayName,
    };
  }
}
