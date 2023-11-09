/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
/**
 * @module ToolbarSecretLink
 * `ToolbarSecretLink` styles SecretLink for the Toolbar.
 * It should only be used inside of `Toolbar`.
 *
 * @example
 * <Toolbar>
 *   <ToolbarActions>
 *     <ToolbarSecretLink @params={{array 'vault.cluster.policies.create'}} @type="add">
 *       Create policy
 *     </ToolbarSecretLink>
 *   </ToolbarActions>
 * </Toolbar>
 *
 *
 * @param type="" {String} - Use "add" to change icon
 */
export default class ToolbarSecretLink extends Component {
  get glyph() {
    return this.args.type === 'add' ? 'plus' : 'chevron-right';
  }
}
