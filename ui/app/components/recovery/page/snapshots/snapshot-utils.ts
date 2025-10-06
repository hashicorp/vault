/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */
import Ember from 'ember';
import { timeout } from 'ember-concurrency';

import { ROOT_NAMESPACE } from 'vault/services/namespace';

import type ApiService from 'vault/services/api';

export interface BadgeInfo {
  status: string;
  color: 'critical' | 'highlight' | 'success' | 'warning';
}

export function getSnapshotStatusBadge(status: string | undefined): BadgeInfo {
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
        status: status || 'Unknown',
        color: 'warning',
      };
  }
}

export function createPollingTask(
  snapshotId: string | undefined,
  api: ApiService,
  onStatus: (status: string) => void,
  onError: (e: unknown) => void
) {
  let cancelled = false;

  const pollingFunction = async () => {
    if (!snapshotId || cancelled) {
      return;
    }

    let wait = Ember.testing ? 0 : 5000;

    while (!cancelled) {
      await timeout(wait);

      if (cancelled) break;

      try {
        const response = await api.sys.systemReadStorageRaftSnapshotLoadId(
          snapshotId,
          api.buildHeaders({ namespace: ROOT_NAMESPACE })
        );

        const status = response.status ?? '';
        onStatus(status);

        // Stop polling if status reaches error state as it won't change again once erroring
        if (status === 'error') {
          break;
        }

        // Just poll once for tests
        if (Ember.testing) return;

        // We still want to poll occasionally in case of an error
        wait = 30000;
      } catch (e) {
        onError(e);
        break;
      }
    }
  };

  // Return an object with the function and a cancel method
  return {
    start: pollingFunction,
    cancel: () => {
      cancelled = true;
    },
  };
}
