/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

/**
 * @module ToggleButton
 * `ToggleButton` components are used to expand and collapse content with a toggle.
 *
 * @example
 * ```js
 *   <ToggleButton @isOpen={{this.showOptions}} @openLabel="Encrypt Output with PGP" @closedLabel="Encrypt Output with PGP" @onClick={{fn (mut this.showOptions}} />
 *  {{#if showOptions}}
 *     <div>
 *       <p>
 *         I will be toggled!
 *       </p>
 *     </div>
 *   {{/if}}
 * ```
 * @callback onClickCallback
 * @param {boolean} isOpen - determines whether to show open or closed label
 * @param {onClickCallback} onClick - fired when button is clicked
 * @param {string} [openLabel="Hide options"] - The message to display when the toggle is open.
 * @param {string} [closedLabel="More options"] - The message to display when the toggle is closed.
 */
export default class ToggleButtonComponent extends Component {
  get openLabel() {
    return this.args.openLabel || 'Hide options';
  }
  get closedLabel() {
    return this.args.closedLabel || 'More options';
  }
}
