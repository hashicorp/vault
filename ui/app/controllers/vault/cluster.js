/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class VaultClusterController extends Controller {
  @service auth;
  @service permissions;
  @service customMessages;
  @service flashMessages;
  @service('version') vaultVersion;

  queryParams = [{ namespaceQueryParam: { as: 'namespace' } }];
  @tracked namespaceQueryParam = '';

  get activeCluster() {
    return this.auth.activeCluster;
  }
}
