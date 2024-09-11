/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ReplicationModeBaseController from './replication-mode';
import { tracked } from '@glimmer/tracking';

export default class ReplicationIndexController extends ReplicationModeBaseController {
  @tracked modeSelection = 'dr';

  get canEnablePrimary() {
    switch (this.modeSelection) {
      case 'dr':
        return this.model.canEnablePrimaryDr;
      case 'performance':
        return this.model.canEnablePrimaryPerformance;
      default:
        // if there's a problem checking capabilities, default to true
        // since the backend will gate as a fallback
        return true;
    }
  }
  get canEnableSecondary() {
    switch (this.modeSelection) {
      case 'dr':
        return this.model.canEnableSecondaryDr;
      case 'performance':
        return this.model.canEnableSecondaryPerformance;
      default:
        // if there's a problem checking capabilities, default to true
        // since the backend will gate as a fallback
        return true;
    }
  }
}
