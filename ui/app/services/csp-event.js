/*eslint-disable no-constant-condition*/
import { filterBy } from '@ember/object/computed';

import Service from '@ember/service';
import { task, waitForEvent } from 'ember-concurrency';

export default Service.extend({
  events: [],
  connectionViolations: filterBy('events', 'violatedDirective', 'connect-src'),

  attach() {
    this.get('monitor').perform();
  },

  remove() {
    this.get('monitor').cancelAll();
  },

  monitor: task(function*() {
    this.get('events').clear();

    while (true) {
      let event = yield waitForEvent(window.document, 'securitypolicyviolation');
      this.get('events').addObject(event);
    }
  }),
});
