/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { helper as buildHelper } from '@ember/component/helper';

const REPLICATION_MODE_DESCRIPTIONS = {
  dr: 'Disaster Recovery Replication is designed to protect against catastrophic failure of entire clusters. Secondaries do not forward service requests until they are elected and become a new primary.',
  performance:
    'Performance Replication scales workloads horizontally across clusters to make requests faster. Local secondaries handle read requests but forward writes to the primary to be handled.',
};

function replicationModeDescription([mode]) {
  return REPLICATION_MODE_DESCRIPTIONS[mode];
}

export default buildHelper(replicationModeDescription);
