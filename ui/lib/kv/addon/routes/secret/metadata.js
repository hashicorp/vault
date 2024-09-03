/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class KvSecretMetadataRoute extends Route {
  @service capabilities;
  @service secretMountPath;
  @service store;

  fetchMetadata(backend, path) {
    return this.store.queryRecord('kv/metadata', { backend, path }).catch((error) => {
      if (error.message === 'Control Group encountered') {
        throw error;
      }
      return null;
    });
  }

  async model() {
    const parentModel = this.modelFor('secret');
    const { backend, path } = parentModel;
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
    }
    // if users can read secret data they can make an explicit request
    // to retrieve secret data in the component
    return parentModel;
  }
}
