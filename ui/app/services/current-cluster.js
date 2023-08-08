/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class CurrentClusterService extends Service {
  @tracked cluster = null;

  setCluster(cluster) {
    this.cluster = cluster;
  }
}
