/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class PkiKeyRoute extends Route {
  @service secretMountPath;
  @service api;
  @service capabilities;

  async model() {
    const { key_id: keyId } = this.paramsFor('keys/key');
    const backend = this.secretMountPath.currentPath;
    const { canUpdate, canDelete } = await this.capabilities.for('pkiKey', { backend, keyId });
    const key = await this.api.secrets.pkiReadKey(keyId, this.secretMountPath.currentPath);
    return {
      backend,
      key,
      canUpdate,
      canDelete,
    };
  }
}
