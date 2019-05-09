import OuterHTML from './outer-html';
import { computed } from '@ember/object';

export default OuterHTML.extend({
  glyph: computed('type', function() {
    if (this.type == 'add') {
      return 'plus-plain';
    } else {
      return 'chevron-right';
    }
  }),
  tagName: '',
  type: null,
  supportsDataTestProperties: true,
});
