/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';

export default Service.extend({
  cluster: null,

  setCluster(cluster) {
    this.set('cluster', cluster);
  },
});
