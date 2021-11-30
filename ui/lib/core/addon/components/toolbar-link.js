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
 * @param {string} type - Use "add" to change icon
 * @param {boolean} disabled - pass true to disable link
 * @param {string} disabledTooltip - tooltip to display on hover when disabled
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/toolbar-link';

export default Component.extend({
  layout,
  tagName: '',
  supportsDataTestProperties: true,
  type: null,
  disabled: false,
  disabledTooltip: null,
  glyph: computed('type', function() {
    if (this.type == 'add') {
      return 'plus-plain';
    } else {
      return 'chevron-right';
    }
  }),
});
