/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { action } from '@ember/object';
import isDeleted from 'kv/helpers/is-deleted';
import { kvErrorHandler } from 'kv/utils/kv-error-handler';

export default class KvSecretRoute extends Route {
  @service api;
  @service secretMountPath;
  @service capabilities;
  @service version;

  async fetchSecretMetadata(backend, path) {
    // catch error and only return 404 which indicates the secret truly does not exist.
    // control group error is handled by the metadata route
    try {
      return await this.api.secrets.kvV2ReadMetadata(path, backend);
    } catch (error) {
      const { status } = await this.api.parseError(error);
      if (status === 404) {
        throw error;
      }
      return null;
    }
  }

  // this request always returns subkeys for the latest version
  async fetchSubkeys(backend, path) {
    if (this.version.isEnterprise) {
      try {
        return await this.api.secrets.kvV2ReadSubkeys(path, backend);
      } catch (error) {
        // metadata will throw if the secret does not exist
        // kvErrorHandler will extract deletion state and relevant metadata from error
        const { status, response } = await this.api.parseError(error);
        return kvErrorHandler(status, response);
      }
    }
    return null;
  }

  isPatchAllowed({ capabilities, subkeysMeta = {} }) {
    if (this.version.isEnterprise) {
      const { canReadSubkeys, canPatchData } = capabilities;
      if (canReadSubkeys && canPatchData && subkeysMeta) {
        const { deletion_time, destroyed } = subkeysMeta;
        const isLatestActive = isDeleted(deletion_time) || destroyed ? false : true;
        // only the latest secret version can be patched and it must not be deleted or destroyed
        return isLatestActive;
      }
    }
    return false;
  }

  async fetchCapabilities(backend, path) {
    const metadataPath = `${backend}/metadata/${path}`;
    const dataPath = `${backend}/data/${path}`;
    const subkeysPath = `${backend}/subkeys/${path}`;
    const deletePath = `${backend}/delete/${path}`;
    const undeletePath = `${backend}/undelete/${path}`;
    const destroyPath = `${backend}/destroy/${path}`;

    const apiPaths = [metadataPath, dataPath, subkeysPath, deletePath, undeletePath, destroyPath];
    const perms = await this.capabilities.fetch(apiPaths, {
      routeForCache: 'vault.cluster.secrets.backend.kv.secret',
    });

    return {
      canReadData: perms[dataPath].canRead,
      canUpdateData: perms[dataPath].canUpdate,
      canPatchData: perms[dataPath].canPatch,
      canCreateVersionData: perms[dataPath].canUpdate,
      canDeleteVersion: perms[deletePath].canUpdate,
      canDeleteLatestVersion: perms[dataPath].canDelete,
      canDestroyVersion: perms[destroyPath].canUpdate,
      canReadMetadata: perms[metadataPath].canRead,
      canDeleteMetadata: perms[metadataPath].canDelete,
      canUpdateMetadata: perms[metadataPath].canUpdate,
      canUndelete: perms[undeletePath].canUpdate,
      canReadSubkeys: perms[subkeysPath].canRead,
    };
  }

  async model() {
    const backend = this.secretMountPath.currentPath;
    const { name: path } = this.paramsFor('secret');
    const capabilities = await this.fetchCapabilities(backend, path);
    const subkeys = await this.fetchSubkeys(backend, path);
    const metadata = await this.fetchSecretMetadata(backend, path);

    return {
      path,
      backend,
      subkeys,
      metadata,
      isPatchAllowed: this.isPatchAllowed({ capabilities, subkeysMeta: subkeys?.metadata }),
      capabilities,
    };
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
