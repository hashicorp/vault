/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';
import { hash } from 'rsvp';
import { action } from '@ember/object';

export default class KvSecretRoute extends Route {
  @service secretMountPath;
  @service store;
  @service capabilities;

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

  fetchSubkeys(backend, path) {
    const adapter = this.store.adapterFor('kv/data');
    // metadata will throw if the secret does not exist
    // always return here so we get deletion state and relevant metadata
    return adapter.fetchSubkeys(backend, path);
  }

  model() {
    const backend = this.secretMountPath.currentPath;
    const { name: path } = this.paramsFor('secret');

    return hash({
      path,
      backend,
      subkeys: this.fetchSubkeys(backend, path),
      metadata: this.fetchSecretMetadata(backend, path),
      canPatchSecret: this.capabilities.canPatch(`${backend}/data/${path}`),
      // for creating a new secret version
      canUpdateSecret: this.capabilities.canUpdate(`${backend}/data/${path}`),
    });
  }

  @action
  willTransition(transition) {
    // TODO update this comment and refreshes
    // refresh the route if transitioning to secret.index (which happens after delete, undelete or destroy)
    // or transitioning from editing either metadata or secret data (creating a new version)
    const isToIndex = transition.to.name === 'vault.cluster.secrets.backend.kv.secret.index';
    const isFromEdit = transition.from.localName === 'edit';
    if (isToIndex || isFromEdit) {
      this.refresh();
    }
  }
}
