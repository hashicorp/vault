/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class CurrentClusterService extends Service {
  @tracked cluster = null;

  setCluster(cluster) {
    this.cluster = cluster;
  }
}
