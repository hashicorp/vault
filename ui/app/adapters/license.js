/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import ClusterAdapter from './cluster';

export default ClusterAdapter.extend({
  pathForType() {
    return 'license/status';
  },
});
