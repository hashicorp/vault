/**
 * @module Icon
 * `Icon` components are glyphs used to indicate important information.
 *
 * Flight icon documentation at https://flight-hashicorp.vercel.app/
 *
 * @example
 * ```js
 * <Icon @name="cancel-square-outline" />
 * ```
 * @param name=null {String} - The name of the SVG to render inline.
 * @param [sizeClass='m'] {String} - Used for sizing the hs-icon, can be one of 's', 'm', 'l', 'xlm', 'xl', 'xxl'. The default is 'm'.
 * @param [size="16"] {String} - size for flight icon, can be 16 or 24
 *
 */
import Component from '@ember/component';
import { computed } from '@ember/object';
import { assert } from '@ember/debug';
import layout from '../templates/components/icon';
import flightIconMap from '@hashicorp/flight-icons/catalog.json';

const flightIconNames = flightIconMap.assets.mapBy('iconName').uniq();
const SIZES = ['s', 'm', 'l', 'xlm', 'xl', 'xxl'];

export default Component.extend({
  tagName: '',
  layout,
  name: null,
  sizeClass: 'm', // hs-icon specific
  size: null, // fight icon specific

  // favor flight icon set and fall back to structure icons if not found
  isFlightIcon: computed('name', function() {
    return this.name ? flightIconNames.includes(this.name) : false;
  }),
  flightIconSize: computed('size', 'sizeClass', function() {
    if (!this.size) {
      // map sizeClass value to appropriate flight icon size if not provided
      return ['s', 'm', 'l'].includes(this.sizeClass) ? '16' : '24';
    }
    assert(
      `The size property of ${this.toString()} must be either '16' or '24'`,
      ['16', '24'].includes(this.size)
    );
    return this.size;
  }),
  iconClass: computed('sizeClass', function() {
    const size = this.sizeClass;
    assert(
      `The sizeClass property of ${this.toString()} must be one of the following: ${SIZES.join(', ')}`,
      SIZES.includes(size)
    );
    return `hs-icon-${size}`;
  }),
});
