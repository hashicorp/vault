/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Route from '@ember/routing/route';
import { service } from '@ember/service';

export default class KmipScopesCreate extends Route {
  @service store;
  @service secretMountPath;

  beforeModel() {
    this.store.unloadAll('kmip/scope');
  }

  model() {
    return this.store.createRecord('kmip/scope', {
      backend: this.secretMountPath.currentPath,
    });
  }
}
