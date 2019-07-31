import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  classNames: ['confirm-wrapper'],
  wormholeReference: null,
  wormholeId: computed(function() {
    return `confirm-${this.elementId}`;
  }),
  didInsertElement() {
    this.set('wormholeReference', this.element.querySelector(`#${this.wormholeId}`));
  },
});
