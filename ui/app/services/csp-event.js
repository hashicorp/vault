/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { task, waitForEvent } from 'ember-concurrency';
import { addToArray } from 'vault/helpers/add-to-array';

export default class CspEventService extends Service {
  @tracked connectionViolations = [];

  attach() {
    this.monitor.perform();
  }

  remove() {
    this.monitor.cancelAll();
  }

  handleEvent(event) {
    if (event.violatedDirective.startsWith('connect-src')) {
      this.connectionViolations = addToArray(this.connectionViolations, event);
    }
  }

  @task
  *monitor() {
    this.connectionViolations = [];

    while (true) {
      const event = yield waitForEvent(window.document, 'securitypolicyviolation');
      this.handleEvent(event);
    }
  }
}
