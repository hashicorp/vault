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

  fetchSecretData(backend, path) {
    // This will always return a record unless 404 not found (show error) or control group
    return this.store.queryRecord('kv/data', { backend, path });
  }

  fetchSecretMetadata(backend, path) {
    // catch error and do nothing because kv/data model handles metadata capabilities
    return this.store.queryRecord('kv/metadata', { backend, path }).catch(() => {});
  }

  fetchSubkeys(backend, path) {
    const adapter = this.store.adapterFor('kv/data');
    return adapter.fetchSubkeys(backend, path);
  }

  model() {
    const backend = this.secretMountPath.currentPath;
    const { name: path } = this.paramsFor('secret');

    return hash({
      path,
      backend,
      // TODO remove and move to details tab
      secret: this.fetchSecretData(backend, path),
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
