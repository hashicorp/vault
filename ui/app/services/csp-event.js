import Ember from 'ember';
const { computed } = Ember;

export default Ember.Service.extend({
  init() {
    this._super(...arguments);
    this.handleCSP = Ember.run.bind(this, '_handleCSP');
  },

  events: [],

  _handleCSP(event) {
    this.get('events').addObject(event);
  },

  connectionViolations: computed.filterBy('events', 'violatedDirective', 'connect-src'),

  attach() {
    this.get('events').clear();
    window.document.addEventListener('securitypolicyviolation', this.handleCSP, true);
  },

  remove() {
    window.document.removeEventListener('securitypolicyviolation', this.handleCSP, true);
  },
});
