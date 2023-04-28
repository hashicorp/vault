/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
/**
 * @module ToolbarLink
 * `ToolbarLink` components style links and buttons for the Toolbar
 * It should only be used inside of `Toolbar`.
 *
 * @example
 * ```js
 * <Toolbar>
 *   <ToolbarActions>
 *     <ToolbarLink @route="vault.cluster.policies.create" @type="add" @disabled={{true}} @disabledTooltip="This link is disabled">
 *       Create policy
 *     </ToolbarLink>
 *   </ToolbarActions>
 * </Toolbar>
 * ```
 *
 * @param {string} route - route to pass to LinkTo
 * @param {Model} model - model to pass to LinkTo
 * @param {Array} models - array of models to pass to LinkTo
 * @param {Object} query - query params to pass to LinkTo
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
