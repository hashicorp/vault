/**
 * @module ToolbarLink
 * `ToolbarLink` components style links and buttons for the Toolbar
 *
 * @example
 * ```js
 * <ToolbarLink
 *   @params={{array 'vault.cluster.policies.create'}}
 *   @type="add"
 * >
 *   Create policy
 * </ToolbarLink>
 * ```
 *
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
