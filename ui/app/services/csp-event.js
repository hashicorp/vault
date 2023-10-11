/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

/*eslint-disable no-constant-condition*/
import { computed } from '@ember/object';

import Service from '@ember/service';
import { task, waitForEvent } from 'ember-concurrency';

export default Service.extend({
  events: computed(function () {
    return [];
  }),
  connectionViolations: computed('events.@each.violatedDirective', function () {
    return this.events.filter((e) => e.violatedDirective.startsWith('connect-src'));
  }),

  attach() {
    this.monitor.perform();
  },

  remove() {
    this.monitor.cancelAll();
  },

  monitor: task(function* () {
    this.events.clear();

    while (true) {
      const event = yield waitForEvent(window.document, 'securitypolicyviolation');
      this.events.addObject(event);
    }
  }),
});
