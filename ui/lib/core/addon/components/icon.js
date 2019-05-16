/**
 * @module Icon
 * `Icon` components are glyphs used to indicate important information.
 *
 * @example
 * ```js
 * <Icon @glyph="cancel-square-outline" />
 * ```
 * @param glyph=null {String} - The name of the SVG to render inline.
 * @param [size='m'] {String} - The size of the Icon, can be one of 's', 'm', 'l', 'xl', 'xxl'. The default is 'm'.
 *
 */
import Component from '@ember/component';
import { computed } from '@ember/object';
import { assert } from '@ember/debug';
import layout from '../templates/components/icon';

const SIZES = ['s', 'm', 'l', 'xl', 'xxl'];

export default Component.extend({
  tagName: '',
  layout,
  glyph: null,
  size: 'm',
  sizeClass: computed('size', function() {
    let { size } = this;
    assert(
      `The size property of ${this.toString()} must be one of the following: ${SIZES.join(', ')}`,
      SIZES.includes(size)
    );
    if (size === 'm') {
      return '';
    }
    return `hs-icon-${size}`;
  }),
});
