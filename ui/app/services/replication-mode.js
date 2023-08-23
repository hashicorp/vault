/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class ReplicationModeService extends Service {
  @tracked mode = null;

  getMode() {
    return this.mode;
  }

  setMode(mode) {
    this.mode = mode;
  }
}
