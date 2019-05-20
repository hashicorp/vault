/**
 * @module ToolbarActions
 * `ToolbarActions` is a container for toolbar links such as "Add item".
 * It should only be used inside of `Toolbar`.
 *
 * @example
 * ```js
 * <Toolbar>
 *   <ToolbarActions>
 *     <ToolbarLink @params={{array 'vault.cluster.policy.edit'}}>
 *       Edit policy
 *     </ToolbarLink>
 *   </ToolbarActions>
 * </Toolbar>
 * ```
 *
 */

import Component from '@ember/component';
import layout from '../templates/components/toolbar-actions';

export default Component.extend({
  tagName: '',
  layout,
});
