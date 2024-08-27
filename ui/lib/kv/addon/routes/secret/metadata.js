/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class KvSecretMetadataRoute extends Route {
  @service store;
  @service secretMountPath;
  @service capabilities;

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
      if (metadata) {
        return {
          ...parentModel,
          metadata,
        };
      }
      // only request secret data if they can read it AND cannot read metadata
      const canReadSecretData = await this.capabilities.canRead(`${backend}/data/${path}`);
      if (canReadSecretData) {
        return {
          ...parentModel,
          secret: this.store.queryRecord('kv/data', { backend, path }),
        };
      }
    }
    return parentModel;
  }
}
