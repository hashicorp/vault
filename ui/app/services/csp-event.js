import Ember from 'ember';
const { computed } = Ember;

export default Ember.Service.extend({
  init() {
    this._super(...arguments);
    this.handleCSP = this._handleCSP.bind(this);
  },
  events: [],
  _handleCSP(event) {
    this.get('events').addObject(event);
  },
  connectionViolations: computed.filterBy('events', 'violatedDirective', 'connect-src'),
  attach() {
    this.get('events').clear();
    window.document.addEventListener('securitypolicyviolation', this.handleCSP);
  },

  remove() {
    window.document.removeEventListener('securitypolicyviolation', this.handleCSP);
  },
});
