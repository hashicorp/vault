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

  @tracked selectedNamespace: string | null = null;
  @tracked selectedMount = '';
  @tracked resourcePath = '';
  @tracked mountOptions: string[] = [];
  @tracked secretData: { [key: string]: string } | undefined;

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.selectedNamespace = this.namespace.path;
    this.fetchMounts();
  }

  @action
  handleSelect(item: 'selectedNamespace' | 'selectedMount', selection: string) {
    if (item === 'selectedNamespace') {
      this.selectedMount = '';
      this[item] = selection;
    }
    this[item] = selection;
  }

  @action
  updateResourcePath({ target }: { target: HTMLInputElement }) {
    this.resourcePath = target.value;
  }

  @action
  clearSelections() {
    this.selectedNamespace = this.namespace.inRootNamespace ? '' : null;
    this.selectedMount = '';
    this.resourcePath = '';
    // Refetch mounts after clearing
    if (this.namespace.inRootNamespace) {
      this.fetchMounts();
    }
  }

  // TODO make a task?
  @action
  async fetchMounts() {
    try {
      const { secret } = await this.api.sys.internalUiListEnabledVisibleMounts(
        // confirm that "path" is the full namespace path e.g. admin/dev not just "dev"
        this.api.buildHeaders({ namespace: this.selectedNamespace || this.namespace.path })
      );
      // TODO cleanup all this iteration business
      const secretEngines = this.api
        .responseObjectToArray(secret, 'path')
        // Filter for support engines
        .map((e) => new SecretsEngineResource(e));
      this.mountOptions = secretEngines.filter((e) => e.supportsRecovery).map((s) => s.path);
    } catch (error) {
      // Render basic input to manually input path
      this.mountOptions = [];
    }
  }

  @action
  async readResource() {
    // TODO validate required fields

    try {
      const { snapshot_id } = this.args.model.snapshot as { snapshot_id: string };
      console.log(snapshot_id);

      // TODO pass in snapshot_id, does the api spec need to be regenerated? why does swagger have a different version?
      if (this.selectedMount === 'kv/') {
        const { data } = await this.api.secrets.kvV1Read(this.resourcePath, 'kv');
        this.secretData = data as any;
      }

      if (this.selectedMount === 'cubbyhole/') {
        const { data } = await this.api.secrets.cubbyholeRead(this.resourcePath);
        this.secretData = data as any;
      }
    } catch (e) {
      const error = await this.api.parseError(e);
      console.error('Failed to read resource:', error);
    }
  }

  // format secret data for display
  get modelForData() {
    return {
      secretData: this.secretData,
      secretKeyAndValue: this.secretKeyAndValue,
    };
  }

  get secretKeyAndValue() {
    if (!this.secretData || typeof this.secretData !== 'object') {
      return [];
    }

    return Object.entries(this.secretData).map(([key, value]) => ({
      key,
      value: typeof value === 'string' ? value : JSON.stringify(value),
    }));
  }

  // TODO will need to poll for status updates
  get badge() {
    const { status } = this.args.model.snapshot as { status: string };
    switch (status) {
      case 'error':
        return {
          status: 'Error',
          color: 'critical',
        };
      case 'loading':
        return {
          status: 'Loading',
          color: 'highlight',
        };
      case 'ready':
        return {
          status: 'Ready',
          color: 'success',
        };
      default:
        return {
          status,
          color: 'warning',
        };
    }
  }
}

// TODO
// 1. show read secrets in kv view
// 2. set up recover flow (probably same signature issues though)
// 3. clean up
// 4. reach out around that discrepancy
// 5. tests
