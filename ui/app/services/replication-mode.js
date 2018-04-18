import Ember from 'ember';

export default Ember.Service.extend({
  mode: null,

  getMode() {
    this.get('mode');
  },

  setMode(mode) {
    this.set('mode', mode);
  },
});
