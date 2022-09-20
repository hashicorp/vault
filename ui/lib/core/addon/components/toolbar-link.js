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
    // not ideal logic, but updating old system to allow for other kinds of icons
    if (this.args.type === 'add') {
      return 'plus';
    }
    if (!this.args.type) {
      return 'chevron-right'; // this is the default if no type is passed in.
    } else {
      return this.args.type;
    }
  }
}
