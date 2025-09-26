/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { HTMLElementEvent } from 'vault/forms';
import { SnapshotsLoadModel } from 'vault/routes/vault/cluster/recovery/snapshots/load';

import type ApiService from 'vault/services/api';
import type RouterService from '@ember/routing/router-service';

interface Args {
  model: SnapshotsLoadModel;
}

enum Methods {
  AUTOMATED = 'automated',
  MANUAL = 'manual',
}

export default class SnapshotsLoad extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly router: RouterService;

  @tracked selectedLoadMethod: Methods;
  @tracked selectedConfig = '';
  @tracked url = '';
  @tracked file?: ArrayBuffer;

  @tracked urlError = '';
  @tracked configError = '';
  @tracked fileError = '';
  @tracked bannerError = '';

  loadMethods = Methods;
  automatedConfigs: string[];

  constructor(owner: unknown, args: Args) {
    super(owner, args);
    const { configError } = this.args.model;

    if (configError && !this.isExpectedError) {
      this.bannerError = configError.message;
    }
    // We want to default to automated unless it is disabled (which is the case during 404)
    this.selectedLoadMethod = configError?.status === 404 ? Methods.MANUAL : Methods.AUTOMATED;
    this.automatedConfigs = this.args.model.configs;
  }

  // For 403, we'll allow manual entry
  // For 404, we'll show a badge + disable the input
  get isExpectedError() {
    const status = this.args.model.configError?.status;
    return status === 403 || status === 404;
  }

  validateFields(): boolean {
    let hasErrors = false;
    this.fileError = '';
    this.urlError = '';
    this.configError = '';

    switch (this.selectedLoadMethod) {
      case Methods.MANUAL: {
        if (!this.file) {
          this.fileError = 'Please upload a snapshot file';
          hasErrors = true;
        }
        break;
      }
      case Methods.AUTOMATED: {
        if (!this.url) {
          this.urlError = 'Please enter a url';
          hasErrors = true;
        }

        if (!this.selectedConfig) {
          this.configError = 'Please select a config';
          hasErrors = true;
        }
        break;
      }
    }

    return !hasErrors;
  }

  @action
  selectLoadMethod(event: HTMLElementEvent<HTMLInputElement>) {
    this.selectedLoadMethod = event.target.value as Methods;

    const { configError } = this.args.model;

    if (this.selectedLoadMethod === Methods.AUTOMATED && configError && !this.isExpectedError) {
      this.bannerError = configError.message;
    } else {
      this.bannerError = '';
    }
  }

  @action
  updateUrl(event: HTMLElementEvent<HTMLInputElement>) {
    this.url = event.target.value.trim();
  }

  @action
  async loadSnapshot(event: Event) {
    event.preventDefault();

    const isValid = this.validateFields();

    if (!isValid) return;
    try {
      switch (this.selectedLoadMethod) {
        case Methods.AUTOMATED: {
          await this.api.sys.systemWriteStorageRaftSnapshotAutoSnapshotLoadName(this.selectedConfig, {
            url: this.url,
          });

          break;
        }
        case Methods.MANUAL: {
          await this.api.sys.systemWriteStorageRaftSnapshotLoad({ body: this.file });
          break;
        }
        default: {
          // This should never be reached, but just in case
          throw new Error('Unsupported load method');
        }
      }

      this.router.transitionTo('vault.cluster.recovery.snapshots');
    } catch (e) {
      const error = await this.api.parseError(e);

      this.bannerError = `Snapshot load error: ${error.message}`;
    }
  }
}
