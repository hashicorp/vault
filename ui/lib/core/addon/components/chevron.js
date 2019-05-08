/**
 * @module Chevron
 * Icon components are used to...
 *
 * @example
 * ```js
 * <Chevron @param1={param1} @param2={param2} />
 * ```
 *
 * @param param1 {String} - param1 is...
 * @param [param2=value] {String} - param2 is... //brackets mean it is optional and = sets the default value
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
