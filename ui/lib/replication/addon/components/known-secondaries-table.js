import Component from '@ember/component';
import { computed } from '@ember/object';

// TODO: add JSDOC comments

export default Component.extend({
  data: null,
  knownSecondaries: computed('data', function() {
    const { data } = this.data;
    return data.knownSecondaries;
  }),
});
