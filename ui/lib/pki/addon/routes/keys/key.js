/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class PkiKeyRoute extends Route {
  @service secretMountPath;
  @service store;

  model() {
    const { key_id } = this.paramsFor('keys/key');
    return this.store.queryRecord('pki/key', {
      backend: this.secretMountPath.currentPath,
      id: key_id,
    });
  }
}
