/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import SecretsEngineResource from 'vault/resources/secrets/engine';

import type ApiService from 'vault/services/api';
import type NamespaceService from 'vault/services/namespace';
import type { SnapshotManageModel } from 'vault/routes/vault/cluster/recovery/snapshots/snapshot/manage';

interface Args {
  model: SnapshotManageModel;
}

export default class SnapshotManage extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly namespace: NamespaceService;

  @tracked selectedNamespace = '';
  @tracked selectedMount = '';
  @tracked mountOptions: string[] = [];

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.fetchMounts();
  }

  @action
  handleSelect(item: 'selectedNamespace' | 'selectedMount', selection: string) {
    this[item] = selection;
  }

  // TODO make a task?
  @action
  async fetchMounts() {
    try {
      const { secret } = await this.api.sys.internalUiListEnabledVisibleMounts(
        // confirm that "path" is the full namespace path e.g. admin/dev not just "dev"
        this.api.buildHeaders({ namespace: this.namespace.path })
      );
      // TODO cleanup all this iteration business
      const secretEngines = this.api
        .responseObjectToArray(secret, 'path')
        // Filter for support engines
        .map((e) => new SecretsEngineResource(e));
      this.mountOptions = secretEngines.filter((e) => e.supportsRecovery).map((e) => e.path);
    } catch (error) {
      // Render basic input to manually input path
      this.mountOptions = [];
    }
  }
}
