/**
 * @module ToolbarSecretLink
 * `ToolbarSecretLink` styles SecretLink for the Toolbar.
 * It should only be used inside of `Toolbar`.
 *
 * @example
 * ```js
 * <Toolbar>
 *   <ToolbarActions>
 *     <ToolbarSecretLink @params={{array 'vault.cluster.policies.create'}} @type="add">
 *       Create policy
 *     </ToolbarSecretLink>
 *   </ToolbarActions>
 * </Toolbar>
 * ```
 *
 * @param type="" {String} - Use "add" to change icon
 */

import OuterHTML from './outer-html';
import { computed } from '@ember/object';

export default OuterHTML.extend({
  glyph: computed('type', function() {
    if (this.type == 'add') {
      return 'plus-plain';
    } else {
      return 'chevron-right';
    }
  }),
  tagName: '',
  supportsDataTestProperties: true,
});
