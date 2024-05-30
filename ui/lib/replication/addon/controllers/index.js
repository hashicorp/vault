/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import ReplicationModeBaseController from './replication-mode';
import { tracked } from '@glimmer/tracking';

export default class ReplicationIndexController extends ReplicationModeBaseController {
  @tracked modeSelection = 'dr';
}
