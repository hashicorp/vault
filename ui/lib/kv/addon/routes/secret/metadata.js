/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import KvForm from 'vault/forms/secrets/kv';

export default class KvSecretMetadataRoute extends Route {
  @service secretMountPath;
  @service api;

  async fetchMetadata(backend, path) {
    try {
      return await this.api.secrets.kvV2ReadMetadata(path, backend);
    } catch (error) {
      const { response } = await this.api.parseError(error);
      if (response.isControlGroupError) {
        throw response;
      }
      // if users can read secret data they can make an explicit request to retrieve secret data in the component
      return null;
    }
  }

  async model() {
    const parentModel = this.modelFor('secret');
    const { backend, path } = parentModel;
    if (!parentModel.metadata) {
      // metadata read on the secret root fails silently
      // if there's no metadata, try again in case it's a control group
      parentModel.metadata = await this.fetchMetadata(backend, path);
    }

    const { custom_metadata, max_versions, cas_required, delete_version_after } = parentModel.metadata || {};
    return {
      ...parentModel,
      form: new KvForm({ path, custom_metadata, max_versions, cas_required, delete_version_after }),
    };
  }
}
