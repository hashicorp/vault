/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import { assert } from '@ember/debug';
import layout from '../templates/components/icon';
import flightIconMap from '@hashicorp/flight-icons/catalog.json';
const flightIconNames = flightIconMap.assets.mapBy('iconName').uniq();
/**
 * @module Icon
 * `Icon` components are glyphs used to indicate important information.
 *
 * Flight icon documentation at https://flight-hashicorp.vercel.app/
 *
 * @example
 * ```js
 * <Icon @name="cancel-square-outline" @size="24" />
 * ```
 * @param {string} name=null - The name of the SVG to render inline.
 * @param {string} [size=16] - size for flight icon, can be 16 or 24
 *
 */
export default Component.extend({
  tagName: '',
  layout,
  name: null,
  size: '16',

  init() {
    this._super(...arguments);
    assert('Icon component size argument must be either "16" or "24"', ['16', '24'].includes(this.size));
  },

  // favor flight icon set and fall back to structure icons if not found
  isFlightIcon: computed('name', function () {
    return this.name ? flightIconNames.includes(this.name) : false;
  }),
  hsIconClass: computed('size', function () {
    return this.size === '24' ? 'hs-icon-xl' : 'hs-icon-l';
  }),
});
