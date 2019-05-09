/**
 * @module Chevron
 * `Chevron` components render `Icon` with one of the `chevron-` glyphs.
 *
 * @example
 * ```js
 * <Chevron @direction="up" />
 * ```
 *
 * @param [direction="right"] {String} - the direction the chevron icon points. Accepted values are
 * "right", "down", "left", "up".
 * @param [isButton=false] {String} - if true, adjusts the CSS classes to push the icon closer to the right of a button.
 *
 */
import Component from '@ember/component';
import { computed } from '@ember/object';
import { assert } from '@ember/debug';

import layout from '../templates/components/chevron';

const DIRECTIONS = ['right', 'left', 'up', 'down'];

export default Component.extend({
  tagName: '',
  layout,
  direction: 'right',
  isButton: false,
  glyph: computed('direction', function() {
    let { direction } = this;
    assert(
      `The direction property of ${this.toString()} must be one of the following: ${DIRECTIONS.join(', ')}`,
      DIRECTIONS.includes(direction)
    );
    return `chevron-${direction}`;
  }),
});
