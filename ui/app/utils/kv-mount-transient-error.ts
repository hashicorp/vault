/**
 * Copyright IBM Corp. 2016, 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

/**
 * A brand-new KV v2 mount can transiently 400 ("upgrading from non-versioned to versioned data") or
 * 404 ("no handler for route") for a short window right after creation.
 */

const KV_UPGRADE_MESSAGES = [
  'Waiting for the primary to upgrade from non-versioned to versioned data.',
  'Upgrading from non-versioned to versioned data.',
];

const ROUTE_NOT_YET_AVAILABLE_MESSAGE = 'no handler for route';

export function isTransientKvMountError(status: number | undefined, message: string | undefined): boolean {
  if (!message) return false;

  if (status === 400) {
    return KV_UPGRADE_MESSAGES.some((known) => message.includes(known));
  }

  if (status === 404) {
    return message.includes(ROUTE_NOT_YET_AVAILABLE_MESSAGE);
  }

  return false;
}
