import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  tagName: '',
  renderedTrigger: null,
  id: null,
  shouldYield: computed('id', 'renderedTrigger', function() {
    return this.id === this.renderedTrigger;
  }),
});
