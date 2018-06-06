import Ember from 'ember';

export default Ember.Route.extend({
  version: Ember.inject.service(),
  beforeModel() {
    return this.get('version').fetchVersion();
  },

  afterModel() {
    //Ember.run.later(() => {
      //return this.replaceWith('vault');
    //}, 3000);
  },
});
