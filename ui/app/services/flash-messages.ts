/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import FlashMessages from 'ember-cli-flash/services/flash-messages';

export default class FlashMessageService extends FlashMessages {
  stickyInfo(message: string) {
    return this.info(message, {
      sticky: true,
      priority: 300,
    });
  }
}
