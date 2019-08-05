import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/confirm';

export default Component.extend({
  layout,
  itemId: null,
  height: 0,
  wormholeReference: null,
  wormholeId: computed(function() {
    return `confirm-${this.elementId}`;
  }),
  didInsertElement() {
    this.set('wormholeReference', this.element.querySelector(`#${this.wormholeId}`));
  },
  updateHeight: function() {
    let height;
    height = this.openTrigger ? this.element.querySelector('.confirm-overlay').clientHeight : 0;
    this.set('height', height);
  },
  actions: {
    onTrigger: function(itemId) {
      this.set('openTrigger', itemId);
      this.updateHeight();
    },
    onCancel: function() {
      this.set('openTrigger', '');
      this.updateHeight();
    },
  },
});
