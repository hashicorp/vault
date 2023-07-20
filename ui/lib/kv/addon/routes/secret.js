/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';

export default class KvSecretRoute extends Route {
  @service store;
  @service secretMountPath;

  async fetchSecretData(backend, path) {
    return await this.store.queryRecord('kv/data', { backend, path }).catch((error) => error.httpStatus);
  }

  async fetchSecretMetadata(backend, path) {
    return await this.store.queryRecord('kv/metadata', { backend, path }).catch((error) => error.httpStatus);
  }

  model() {
    const backend = this.secretMountPath.currentPath;
    const { name: path } = this.paramsFor('secret');
    return hash({
      path,
      backend,
      secret: this.fetchSecretData(backend, path),
      metadata: this.fetchSecretMetadata(backend, path),
    });
  }
}
