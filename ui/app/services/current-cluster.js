/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Service from '@ember/service';

export default Service.extend({
  cluster: null,

  setCluster(cluster) {
    this.set('cluster', cluster);
  },
});
