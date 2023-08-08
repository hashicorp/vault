/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
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
