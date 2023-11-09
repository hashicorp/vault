/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
/**
 * @module ToolbarLink
 * ToolbarLink components style links and buttons for the Toolbar
 * It should only be used inside of Toolbar component.
 *
 * @example
  <Toolbar>
   <ToolbarActions>
     <ToolbarLink @route="vault"  @disabled={{true}} @disabledTooltip="This link is disabled">
       Disabled link
     </ToolbarLink>
     <ToolbarLink @route="vault" @type="add">
       Create item
     </ToolbarLink>
   </ToolbarActions>
 </Toolbar>
 *
 *
 * @param {string} route - route to pass to LinkTo
 * @param {model} model - model to pass to LinkTo
 * @param {array} models - array of models to pass to LinkTo
 * @param {object} query - query params to pass to LinkTo
 * @param {boolean} replace - replace arg to pass to LinkTo
 * @param {string} type - Use "add" to change icon to plus sign, or pass in your own kind of icon.
 * @param {boolean} disabled - pass true to disable link
 * @param {string} disabledTooltip - tooltip to display on hover when disabled
 */

export default class ToolbarLinkComponent extends Component {
  get glyph() {
    // not ideal logic. Without refactoring, this allows us to add in our own icon type outside of chevron-right or plus.
    // For a later refactor we should remove the substitution for add to plus and just return type.
    const { type } = this.args;
    if (!type) return 'chevron-right';
    return type === 'add' ? 'plus' : type;
  }
  get models() {
    const { model, models } = this.args;
    if (model) {
      return [model];
    }
    return models || [];
  }
  get query() {
    return this.args.query || {};
  }
}
