/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import EmberCliFlash from 'ember-cli-flash/services/flash-messages';

export default class FlashMessages extends EmberCliFlash {
  stickyInfo(message: string) {
    return this.info(message, {
      sticky: true,
      priority: 300,
    });
  }
}
