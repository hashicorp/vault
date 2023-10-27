/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';

import type SyncDestinationModel from 'vault/models/sync/destination';
import type FlashMessageService from 'vault/services/flash-messages';

interface Args {
  destination: SyncDestinationModel;
}

export default class DestinationsTabsToolbar extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;

  @tracked isDeleteModalOpen = false;

  @action
  deleteDestination() {
    // TODO wire up delete + purge when endpoint exists
    this.isDeleteModalOpen = false;
  }
}
