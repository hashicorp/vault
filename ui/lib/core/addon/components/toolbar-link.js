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
 *     <ToolbarLink @params={{array 'vault.cluster.policies.create'}} @type="add" @disabled={{true}} @disabledTooltip="This link is disabled">
 *       Create policy
 *     </ToolbarLink>
 *   </ToolbarActions>
 * </Toolbar>
 * ```
 *
 * @param {array} params - Array to pass to LinkTo
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
}
