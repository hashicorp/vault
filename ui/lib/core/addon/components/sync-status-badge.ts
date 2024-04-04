/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

interface Args {
  status: string; // https://developer.hashicorp.com/vault/docs/sync#sync-statuses
}

export default class SyncStatusBadge extends Component<Args> {
  get state() {
    switch (this.args.status) {
      case 'SYNCING':
        return {
          icon: 'sync',
          color: 'neutral',
        };
      case 'SYNCED':
        return {
          icon: 'check-circle',
          color: 'success',
        };
      case 'UNSYNCING':
        return {
          icon: 'sync-reverse',
          color: 'neutral',
        };
      case 'UNSYNCED':
        return {
          icon: 'sync-alert',
          color: 'warning',
        };
      case 'INTERNAL_VAULT_ERROR':
        return {
          icon: 'x-circle',
          color: 'critical',
        };
      case 'CLIENT_SIDE_ERROR':
        return {
          icon: 'x-circle',
          color: 'critical',
        };
      case 'EXTERNAL_SERVICE_ERROR':
        return {
          icon: 'x-circle',
          color: 'critical',
        };
      case 'UNKNOWN':
        return {
          icon: 'help',
          color: 'neutral',
        };
      default:
        return {
          icon: 'help',
          color: 'neutral',
        };
    }
  }
}
