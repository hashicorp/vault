/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';

export default class MessagesEditController extends Controller {
  queryParams = ['authenticated'];

  authenticated = true;

  @action
  onUpdateEndTime(endTime) {
    this.model.message.endTime = endTime;
  }
}
