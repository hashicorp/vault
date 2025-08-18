/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import SecretsEngineResource from 'vault/resources/secrets/engine';
import { sanitizePath } from 'core/utils/sanitize-path';
import errorMessage from 'vault/utils/error-message';

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
  @tracked showRecoveryBanner = false;
  @tracked mountError = '';
  @tracked resourcePathError = '';

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    this.selectedNamespace = this.namespace.path === '' ? 'root' : null;
    this.fetchMounts();
  }

  @action
  handleSelect(item: 'selectedNamespace' | 'selectedMount', selection: string) {
    if (item === 'selectedNamespace') {
      this.selectedMount = '';
      this[item] = selection;
      this.fetchMounts();
    }
    this[item] = selection;

    // Clear errors when user makes selections
    if (item === 'selectedMount') {
      this.mountError = '';
    }
  }

  @action
  clearSelections() {
    this.selectedNamespace = this.namespace.inRootNamespace ? '' : null;
    this.selectedMount = '';
    this.resourcePath = '';
    this.mountError = '';
    this.resourcePathError = '';
    this.secretData = undefined;
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
  async recover() {
    // Validate fields before attempting recovery
    const isValid = await this.validateFields(true);
    if (!isValid) {
      return;
    }

    try {
      const { snapshot_id } = this.args.model.snapshot as { snapshot_id: string };

      // TODO pass in snapshot_id once spec is updated
      if (this.selectedMount === 'kv/') {
        // TODO: Implement actual recovery logic for KV
      }
      if (this.selectedMount === 'cubbyhole/') {
        // TODO: Implement actual recovery logic for cubbyhole
      }

      this.showRecoveryBanner = true;
    } catch (e) {
      const error = await this.api.parseError(e);
      console.error('Failed to recover resource:', error);
    }
  }

  @action
  async readResource() {
    // Validate fields before attempting to read
    const isValid = await this.validateFields(false);
    if (!isValid) {
      return;
    }

    try {
      const { snapshot_id } = this.args.model.snapshot as { snapshot_id: string };

      // TODO pass in snapshot_id once spec is updated
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

  @action
  updateResourcePath({ target }: { target: HTMLInputElement }) {
    this.resourcePath = target.value;
    // Clear error when user starts typing
    this.resourcePathError = '';
  }

  @action
  async validateFields(isRecovery = false): Promise<boolean> {
    this.mountError = '';
    this.resourcePathError = '';
    let hasErrors = false;

    if (!this.selectedMount) {
      this.mountError = 'Please select a secret mount';
      hasErrors = true;
    }

    if (!this.resourcePath.trim()) {
      this.resourcePathError = 'Please enter a resource path';
      hasErrors = true;
    }

    if (hasErrors) {
      return false;
    }

    try {
      await this.checkResourceExists();

      return true;
    } catch (e) {
      const error = await this.api.parseError(e);

      if (error.status === 404) {
        this.resourcePathError = 'Resource does not exist at this path';
      } else if (error.status === 403) {
        if (isRecovery) {
          this.resourcePathError = 'You do not have permission to recover secrets to this path';
        } else {
          this.resourcePathError = 'You do not have permission to read from this path';
        }
      } else {
        this.resourcePathError = errorMessage(error) || 'Failed to validate resource path';
      }

      return false;
    }
  }

  async checkResourceExists(): Promise<void> {
    if (this.selectedMount === 'kv/') {
      await this.api.secrets.kvV1Read(this.resourcePath, 'kv');
    } else if (this.selectedMount === 'cubbyhole/') {
      await this.api.secrets.cubbyholeRead(this.resourcePath);
    }
  }

  // async checkReadPermission(): Promise<void> {}

  get hasValidationErrors() {
    return !!(this.mountError || this.resourcePathError);
  }

  get isFormValid() {
    return this.selectedMount && this.resourcePath.trim() && !this.hasValidationErrors;
  }

  get mountWithoutSlash() {
    return sanitizePath(this.selectedMount);
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
