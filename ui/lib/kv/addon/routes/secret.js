/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import { action } from '@ember/object';
import { isDeleted } from 'kv/utils/kv-deleted';

export default class KvSecretRoute extends Route {
  @service secretMountPath;
  @service store;
  @service capabilities;
  @service version;

  fetchSecretMetadata(backend, path) {
    // catch error and only return 404 which indicates the secret truly does not exist.
    // control group error is handled by the metadata route
    return this.store.queryRecord('kv/metadata', { backend, path }).catch((e) => {
      if (e.httpStatus === 404) {
        throw e;
      }
      return null;
    });
  }

  // this request always returns subkeys for the latest version
  fetchSubkeys(backend, path) {
    if (this.version.isEnterprise) {
      const adapter = this.store.adapterFor('kv/data');
      // metadata will throw if the secret does not exist
      // always return here so we get deletion state and relevant metadata
      return adapter.fetchSubkeys(backend, path);
    }
    return null;
  }

  isPatchAllowed({ capabilities, subkeysMeta }) {
    if (!this.version.isEnterprise) return false;
    const canReadSubkeys = capabilities.subkeys.canRead;
    const canPatchData = capabilities.data.canPatch;
    if (canReadSubkeys && canPatchData) {
      const { deletion_time, destroyed } = subkeysMeta;
      const isLatestActive = isDeleted(deletion_time) || destroyed ? false : true;
      // only the latest secret version can be patched and it must not be deleted or destroyed
      return isLatestActive;
    }
    return false;
  }

  async fetchCapabilities(backend, path) {
    const metadataPath = `${backend}/metadata/${path}`;
    const dataPath = `${backend}/data/${path}`;
    const subkeysPath = `${backend}/subkeys/${path}`;
    const perms = await this.capabilities.fetchMultiplePaths([metadataPath, dataPath, subkeysPath]);
    return {
      metadata: perms[metadataPath],
      data: perms[dataPath],
      subkeys: perms[subkeysPath],
    };
  }

  async model() {
    const backend = this.secretMountPath.currentPath;
    const { name: path } = this.paramsFor('secret');
    const capabilities = await this.fetchCapabilities(backend, path);
    const subkeys = await this.fetchSubkeys(backend, path);
    return hash({
      path,
      backend,
      subkeys,
      metadata: this.fetchSecretMetadata(backend, path),
      isPatchAllowed: this.isPatchAllowed({ capabilities, subkeysMeta: subkeys?.metadata }),
      canUpdateData: capabilities.data.canUpdate,
      canReadData: capabilities.data.canRead,
      canReadMetadata: capabilities.metadata.canRead,
      canDeleteMetadata: capabilities.metadata.canDelete,
      canUpdateMetadata: capabilities.metadata.canUpdate,
    });
  }

  @action
  willTransition(transition) {
    // refresh the route if transitioning to secret.index (which happens after delete, undelete or destroy)
    // or transitioning from editing either metadata or secret data (creating a new version)
    const isToIndex = transition.to.name === 'vault.cluster.secrets.backend.kv.secret.index';
    const isFromEdit = transition.from.localName === 'edit';
    if (isToIndex || isFromEdit) {
      this.refresh();
    }
  }
}
