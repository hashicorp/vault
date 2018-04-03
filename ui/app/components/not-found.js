import Ember from 'ember';

const { computed, inject } = Ember;

export default Ember.Component.extend({
  // public
  model: null,

  tagName: '',
  routing: inject.service('-routing'),
  path: computed.alias('routing.router.currentURL'),
});
