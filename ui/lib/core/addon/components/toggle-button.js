/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @module ToggleButton
 * `ToggleButton` components are used to expand and collapse content with a toggle.
 *
 * @example
 * <ToggleButton @isOpen={{this.showOptions}} @openLabel="Show stuff" @closedLabel="Hide the stuff" @onClick={{fn (mut this.showOptions) (not this.showOptions)}} />
 *  {{#if this.showOptions}}
 *     <div>
 *       <p>
 *         I will be toggled!
 *       </p>
 *     </div>
 *   {{/if}}
 *
 * @param {boolean} isOpen - determines whether to show open or closed label
 * @param {function} onClick - fired when button is clicked
 * @param {string} [openLabel=Hide options] - The message to display when the toggle is open.
 * @param {string} [closedLabel=More options] - The message to display when the toggle is closed.
 */
export default class ToggleButtonComponent extends Component {
  get openLabel() {
    return this.args.openLabel || 'Hide options';
  }
  get closedLabel() {
    return this.args.closedLabel || 'More options';
  }
}
