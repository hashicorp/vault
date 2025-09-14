/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { Factory } from 'miragejs';
import timestamp from 'core/utils/timestamp';
import { addHours } from 'date-fns';

export default Factory.extend({
  status: 'ready',
  // Snapshots expire after 72 hours
  expires_at: addHours(timestamp.now(), 72),
  snapshot_id: '9465df92-8236-4af9-8cc8-b7460d882e41s',
});
