import { later } from '@ember/runloop';
import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  modalContainer: computed('isActive', function() {
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
      later(
        this,
        function() {
          this.set('isAnimated', true);
        },
        10
      );
    },
    closeOverlay() {
      this.set('isAnimated', false);
      later(
        this,
        function() {
          this.set('isActive', false);
        },
        300
      );
    },
  },
});
