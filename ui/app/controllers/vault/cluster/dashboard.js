/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import timestamp from 'core/utils/timestamp';

export default class DashboardController extends Controller {
  @tracked replicationUpdatedAt = timestamp.now().toISOString();

  @action
  refreshModel(e) {
    if (e) e.preventDefault();
    this.replicationUpdatedAt = timestamp.now().toISOString();
    this.send('refreshRoute');
  }
}
