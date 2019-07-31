import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  classNames: ['confirm-wrapper'],
  itemId: null,
  transitionDirection: '',
  wormholeReference: null,
  wormholeId: computed(function() {
    return `confirm-${this.elementId}`;
  }),
  didInsertElement() {
    this.set('wormholeReference', this.element.querySelector(`#${this.wormholeId}`));
  },
  actions: {
    onTrigger: function(itemId) {
      this.set('openTrigger', itemId);
      this.set('transitionDirection', 'left');
    },
    onCancel: function() {
      this.set('transitionDirection', 'right');
      // this.set('openTrigger', '');
    },
  },
});
