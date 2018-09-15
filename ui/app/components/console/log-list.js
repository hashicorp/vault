import Component from '@ember/component';
import { computed } from '@ember/object';

export default Component.extend({
  content: null,
  list: computed('content', function() {
    return this.get('content').keys;
  }),
});
