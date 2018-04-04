import Ember from 'ember';

export default Ember.Controller.extend({
  queryParams: ['action'],
  action: '',
  reset() {
    this.set('action', '');
  },
});
