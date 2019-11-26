/**
 * @module ToolbarLink
 * `ToolbarLink` components style links and buttons for the Toolbar
 * It should only be used inside of `Toolbar`.
 *
 * @example
 * ```js
 * <Toolbar>
 *   <ToolbarActions>
 *     <ToolbarLink @params={{array 'vault.cluster.policies.create'}} @type="add">
 *       Create policy
 *     </ToolbarLink>
 *   </ToolbarActions>
 * </Toolbar>
 * ```
 *
 * @param params=''{Array} Array to pass to LinkTo
 * @param type=''{String} Use "add" to change icon
 */

import Component from '@ember/component';
import { computed } from '@ember/object';
import layout from '../templates/components/toolbar-link';

export default Component.extend({
  layout,
  tagName: '',
  supportsDataTestProperties: true,
  type: null,
  glyph: computed('type', function() {
    if (this.type == 'add') {
      return 'plus-plain';
    } else {
      return 'chevron-right';
    }
  }),
});
