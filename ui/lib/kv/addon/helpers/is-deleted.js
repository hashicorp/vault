/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import timestamp from 'core/utils/timestamp';

export default function isDeleted(date) {
  // deletion_time does not always mean the secret has been deleted.
  // if the delete_version_after is set then the deletion_time will be UTC of that time, even if it's a future time from now.
  // to determine if the secret is deleted we check if deletion_time <= time right now.
  const deletionTime = new Date(date);
  const now = timestamp.now();
  return deletionTime <= now;
}
