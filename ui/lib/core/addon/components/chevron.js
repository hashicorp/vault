/**
 * @module Chevron
 * Icon components are used to...
 *
 * @example
 * ```js
 * <Icon @param1={param1} @param2={param2} />
 * ```
 *
 * @param param1 {String} - param1 is...
 * @param [param2=value] {String} - param2 is... //brackets mean it is optional and = sets the default value
 */
import Icon from './icon';

const DEFAULT_GLYPH = 'chevron-right';

export default Icon.extend({
  'aria-hidden': 'true',
  direction: null,
  didReceiveAttrs() {
    this._super(...arguments);
    if (!this.glyph) {
      this.set('glyph', DEFAULT_GLYPH);
    }
    if (this.direction) {
      this.set('glyph', `chevron-${this.direction}`);
    }
  },
});
