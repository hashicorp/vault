/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { assert } from '@ember/debug';
import flightIconMap from '@hashicorp/flight-icons/catalog.json';
const flightIconNames = flightIconMap.assets.map((asset) => asset.iconName).uniq();

/**
 * @module Icon
 * `Icon` components are used to display an icon.
 *
 * Flight icon documentation at https://helios.hashicorp.design/icons/usage-guidelines?tab=code#how-to-use-icons
 * Flight icon library at https://helios.hashicorp.design/icons/library
 *
 * @example
 * ```js
 * <Icon @name="x-square" @size="24" />
 * ```
 * @param {string} name - The name of the SVG to render inline. Required.
 * @param {string} [size=16] - size for flight icon, can be 16 or 24
 *
 */

export default class Icon extends Component {
  constructor(owner, args) {
    super(owner, args);
    assert('Icon component size argument must be either "16" or "24"', ['16', '24'].includes(this.size));
    assert('Icon name argument must be provided', this.args.name);
  }

  get size() {
    return this.args.size || '16';
  }

  // favor flight icon set and fall back to structure icons if not found
  get isFlightIcon() {
    return this.args.name ? flightIconNames.includes(this.args.name) : false;
  }

  get hsIconClass() {
    return this.size === '24' ? 'hs-icon-xlm' : 'hs-icon-l';
  }
}
