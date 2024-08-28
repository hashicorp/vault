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

  async fetchCapabilities(backend, path) {
    const metadataPath = `${backend}/metadata/${path}`;
    const dataPath = `${backend}/data/${path}`;
    const resp = await this.capabilities.fetchMultiplePaths([metadataPath, dataPath]);

    const findPerms = (path) => {
      const model = resp.find((m) => m.id === path);
      const { canRead, canUpdate, canDelete } = model;
      return { canRead, canUpdate, canDelete };
    };

    return {
      metadata: findPerms(metadataPath),
      data: findPerms(dataPath),
    };
  }

  async model() {
    const parentModel = this.modelFor('secret');
    const { backend, path } = parentModel;
    const permissions = await this.fetchCapabilities(backend, path);
    const model = {
      ...parentModel,
      permissions,
    };
    if (!parentModel.metadata) {
      // metadata read on the secret root fails silently
      // if there's no metadata, try again in case it's a control group
      const metadata = await this.fetchMetadata(backend, path);
      if (metadata) {
        return {
          ...model,
          metadata,
        };
      }
      // only fetch secret data if metadata is unavailable and user can read endpoint
      if (permissions.data.canRead) {
        const secret = await this.store.queryRecord('kv/data', { backend, path });
        return {
          ...model,
          secret,
        };
      }
    }
    return model;
  }
}
