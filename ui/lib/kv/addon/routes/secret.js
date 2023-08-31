/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { action } from '@ember/object';

export default class KvSecretRoute extends Route {
  @service secretMountPath;
  @service store;

  fetchSecretData(backend, path) {
    // This will always return a record unless 404 not found (show error) or control group
    return this.store.queryRecord('kv/data', { backend, path });
  }

  fetchSecretMetadata(backend, path) {
    // catch error and do nothing because kv/data model handles metadata capabilities
    return this.store.queryRecord('kv/metadata', { backend, path }).catch(() => {});
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

  @action
  willTransition(transition) {
    const { to, from } = transition;
    if (to.name === 'vault.cluster.secrets.backend.kv.secret.index' || from.localName === 'edit') {
      this.refresh();
    }
  }
}
