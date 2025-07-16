/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module InfoTooltip
 *
 * @example
 *  <InfoTooltip>Important info!</InfoTooltip>
 *
 * @param {string} [verticalPosition] - vertical position specification (above, below)
 * @param {string} [horizontalPosition] - horizontal position specification (center, auto-right)

 */

export default class InfoTooltip extends Component {
  @action
  preventSubmit(e) {
    e.preventDefault();
  }
}
