import Ember from 'ember';

export default Ember.Component.extend({
  isAnimated: false,
  isActive: false,
  tagName: 'span',
  actions: {
    openOverlay() {
      this.set('isActive', true);
      Ember.run.later(
        this,
        function() {
          this.set('isAnimated', true);
        },
        10
      );
    },
    closeOverlay() {
      this.set('isAnimated', false);
      Ember.run.later(
        this,
        function() {
          this.set('isActive', false);
        },
        300
      );
    },
  },
});
