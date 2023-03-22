/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import PkiKeysIndexRoute from './index';

export default class PkiKeyRoute extends PkiKeysIndexRoute {
  model() {
    const { key_id } = this.paramsFor('keys/key');
    return this.store.queryRecord('pki/key', {
      backend: this.secretMountPath.currentPath,
      id: key_id,
    });
  }
}
