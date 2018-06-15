import Ember from 'ember';

const { computed } = Ember;

export default Ember.Component.extend({
  modalContainer: computed(function() {
    return document.getElementById('modal-wormhole');
  }),
  isAnimated: false,
  isActive: false,
  tagName: 'span',
  trackingSource: computed('pageName', function() {
    let trackingSource = 'vaultui';
    let pageName = this.get('pageName');
    if (pageName) {
      trackingSource = trackingSource + '_' + encodeURIComponent(pageName);
    }
    return trackingSource;
  }),
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
