/*eslint-disable no-constant-condition*/
import { computed } from '@ember/object';
import { filterBy } from '@ember/object/computed';

import Service from '@ember/service';
import { task, waitForEvent } from 'ember-concurrency';

export default Service.extend({
  events: computed(function() {
    return [];
  }),
  connectionViolations: filterBy('events', 'violatedDirective', 'connect-src'),

  attach() {
    this.monitor.perform();
  },

  remove() {
    this.monitor.cancelAll();
  },

  monitor: task(function*() {
    this.events.clear();

    while (true) {
      let event = yield waitForEvent(window.document, 'securitypolicyviolation');
      this.events.addObject(event);
    }
  }),
});
