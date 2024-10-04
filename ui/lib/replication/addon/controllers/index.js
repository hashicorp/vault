/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ReplicationModeBaseController from './replication-mode';
import { tracked } from '@glimmer/tracking';

export default class ReplicationIndexController extends ReplicationModeBaseController {
  @tracked modeSelection = 'dr';

  getPerm(type) {
    if (this.modeSelection === 'dr') {
      // returns canEnablePrimaryDr or canEnableSecondaryDr
      return `canEnable${type}Dr`;
    }
    if (this.modeSelection === 'performance') {
      // returns canEnablePrimaryPerformance or canEnableSecondaryPerformance
      return `canEnable${type}Performance`;
    }
  }

  // if there's a problem checking capabilities, default to true
  // since the backend will gate as a fallback
  get canEnablePrimary() {
    return this.model[this.getPerm('Primary')] ?? true;
  }
  get canEnableSecondary() {
    return this.model[this.getPerm('Secondary')] ?? true;
  }
}
