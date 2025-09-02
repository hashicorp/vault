/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { assert } from '@ember/debug';

const DIRECTIONS = ['right', 'left', 'up', 'down'];

/**
 * @deprecated
 * @module Chevron
 * * Use HDS icons instead
 * Chevron components render Icon with one of the "chevron-" glyphs.
 *
 * @example
 * <Chevron @direction="up" />
 *
 * @param {string} [direction="right"] - the direction the chevron icon points. Accepted values are "right", "down", "left", "up".
 * @param {string} [isButton=false] - if true, adjusts the CSS classes to push the icon closer to the right of a button.
 *
 */

export default class Chevron extends Component {
  get direction() {
    return this.args.direction || 'right';
  }

  get glyph() {
    const { direction } = this;

    assert(
      `The direction property of Chevron must be one of the following: ${DIRECTIONS.join(', ')}`,
      DIRECTIONS.includes(direction)
    );
    return `chevron-${direction}`;
  }
}
