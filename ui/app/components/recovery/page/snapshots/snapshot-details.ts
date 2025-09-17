/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { service } from '@ember/service';
import {
  createPollingTask,
  getSnapshotStatusBadge,
} from 'vault/components/recovery/page/snapshots/snapshot-utils';
import { dateFormat } from 'core/helpers/date-format';
import { tracked } from '@glimmer/tracking';

import type ApiService from 'vault/services/api';
import type RouterService from '@ember/routing/router-service';
import type { SnapshotManageModel } from 'vault/routes/vault/cluster/recovery/snapshots/snapshot/manage';
import type FlashMessageService from 'vault/services/flash-messages';

interface Args {
  model: SnapshotManageModel;
}

export default class SnapshotDetails extends Component<Args> {
  @service declare readonly api: ApiService;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly router: RouterService;

  @tracked snapshotStatus?: string = '';

  private pollingController: { start: () => Promise<void>; cancel: () => void } | null = null;

  constructor(owner: unknown, args: Args) {
    super(owner, args);

    // Create and start polling task
    this.pollingController = createPollingTask(
      this.args.model.snapshot.snapshot_id,
      this.api,
      this.onPollSuccess,
      this.onPollError
    );
    this.pollingController.start();
  }

  willDestroy() {
    super.willDestroy();
    if (this.pollingController) {
      this.pollingController.cancel();
    }
  }

  get badge() {
    // Use polled status if available, otherwise fall back to initial model status
    const status = this.snapshotStatus || this.args.model.snapshot?.status;

    return getSnapshotStatusBadge(status);
  }

  get tableColumns() {
    const snapshot = this.args.model.snapshot;

    const columns = [
      {
        label: 'Snapshot ID',
        key: 'snapshot_id',
        value: snapshot.snapshot_id,
      },
      {
        label: 'Expiring at',
        key: 'expires_at',
        value: dateFormat([snapshot.expires_at, 'MMM d, yyyy hh:mm aaa'], {
          withTimeZone: true,
        }),
      },
    ];

    if (this.args.model.snapshot.auto_snapshot_config) {
      const automatedSnapshotCols = [
        {
          label: 'Auto config name',
          key: 'auto_config_name',
          value: snapshot.auto_snapshot_config,
        },
        {
          label: 'URL',
          key: 'url',
          value: snapshot.url,
        },
      ];
      columns.splice(1, 0, ...automatedSnapshotCols);
    }

    return columns;
  }

  onPollError = async () => {
    this.snapshotStatus = 'error';
    // Update the model directly for reactive updates across routes
    this.args.model.snapshot.status = 'error';
  };

  onPollSuccess = async (status: string) => {
    this.snapshotStatus = status;
    // Update the model directly for reactive updates across routes
    this.args.model.snapshot.status = status as 'ready' | 'error' | 'loading';
  };

  @action
  async unloadSnapshot() {
    try {
      const { snapshot_id } = this.args.model.snapshot as { snapshot_id: string };
      await this.api.sys.systemDeleteStorageRaftSnapshotLoadId(snapshot_id);

      this.router.transitionTo('vault.cluster.recovery.snapshots');
    } catch (e) {
      const { message } = await this.api.parseError(e);
      this.flashMessages.danger(`Snapshot was not unloaded: ${message}`);
    }
  }
}
