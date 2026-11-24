/**
 * Copyright IBM Corp. 2026
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';
import { service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class ControlGroupsController extends Controller {
  @service router;

  @tracked accessor = '';

  @action
  lookup() {
    this.router.transitionTo('vault.cluster.access.control-group-accessor', this.accessor);
  }
}
