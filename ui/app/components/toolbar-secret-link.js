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
 * ```js
 * <Toolbar>
 *   <ToolbarActions>
 *     <ToolbarSecretLink
 *       @mode="create"
 *       @type="add"
 *       @secret="some-secret"
 *       @backend="mount-path"
 *       @queryParams={{hash tab="policy"}}
 *       @replace={{true}}
 *       @disabled={{false}}
 *       data-test-custom-tag
 *     >
 *       Create policy
 *     </ToolbarSecretLink>
 *   </ToolbarActions>
 * </Toolbar>
 * ```
 *
 * @param {string} type - use "add" to change icon from "chevron-right" to "plus"
 * @param {string} mode - *required* passed to secret-link, controls route
 * @param {string} backend - *required* backend path. Passed to secret-link
 * @param {string} secret - secret path. Passed to secret-link
 * @param {boolean} replace - passed to secret-link
 * @param {boolean} disabled - passed to secret-link
 * @param {object} queryParams - passed to secret-link
 */
export default class ToolbarSecretLink extends Component {
  get glyph() {
    return this.args.type === 'add' ? 'plus' : 'chevron-right';
  }
}
