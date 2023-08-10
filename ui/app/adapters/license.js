/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import ClusterAdapter from './cluster';

export default ClusterAdapter.extend({
  pathForType() {
    return 'license/status';
  },
});
