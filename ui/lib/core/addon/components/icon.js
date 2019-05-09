/**
 * @module Icon
 * `Icon` components are glyphs used to indicate important information.
 *
 * @example
 * ```js
 * <ICon @glyph="cancel-square-outline" />
 * ```
 * @param glyph=null {String} - The name of the SVG to render inline.
 *
 */
import Component from '@ember/component';
import layout from '../templates/components/icon';

export default Component.extend({
  tagName: '',
  layout,
  glyph: null,
});
