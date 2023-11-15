/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class KvSecretMetadataRoute extends Route {
  @service store;
  @service secretMountPath;

  fetchMetadata(backend, path) {
    return this.store.queryRecord('kv/metadata', { backend, path }).catch((error) => {
      if (error.message === 'Control Group encountered') {
        throw error;
      }
      return {};
    });
  }

  async model() {
    const backend = this.secretMountPath.currentPath;
    const { name: path } = this.paramsFor('secret');
    const parentModel = this.modelFor('secret');
    if (!parentModel.metadata) {
      // metadata read on the secret root fails silently
      // if there's no metadata, try again in case it's a control group
      const metadata = await this.fetchMetadata(backend, path);
      return {
        ...parentModel,
        metadata,
      };
    }
    return parentModel;
  }
}
