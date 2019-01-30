import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  tagName: '',

  text: computed(function() {
    return 'home';
  }),

  computedClasses: computed('classNames', function() {
    return this.get('classNames').join(' ');
  }),
});
