import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/toolbar-link';

export default Component.extend({
  layout,
  tagName: '',
  supportsDataTestProperties: true,
  type: null,
  glyph: computed('type', function() {
    if (this.type == 'add') {
      return 'plus-plain';
    } else {
      return 'chevron-right';
    }
  }),
});
