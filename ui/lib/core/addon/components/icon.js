/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { assert } from '@ember/debug';
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

export default class Icon extends Component {
  constructor(owner, args) {
    super(owner, args);
    assert('Icon component size argument must be either "16" or "24"', ['16', '24'].includes(this.size));
  }

  get size() {
    return this.args.size || '16';
  }

  get name() {
    return this.args.name || null;
  }

  // favor flight icon set and fall back to structure icons if not found
  get isFlightIcon() {
    return this.name ? flightIconNames.includes(this.name) : false;
  }

  get hsIconClass() {
    return this.size === '24' ? 'hs-icon-xlm' : 'hs-icon-l';
  }
}
