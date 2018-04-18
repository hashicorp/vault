/*eslint-disable no-constant-condition*/
import Ember from 'ember';
import { task, waitForEvent } from 'ember-concurrency';

const { Service, computed } = Ember;

export default Service.extend({
  events: [],
  connectionViolations: computed.filterBy('events', 'violatedDirective', 'connect-src'),

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
