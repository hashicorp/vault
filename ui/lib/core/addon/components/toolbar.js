/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/**
 * @module Toolbar
 * `Toolbar` components are containers for Toolbar actions.
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
 */

import Component from '@ember/component';
import layout from '../templates/components/toolbar';

export default Component.extend({
  tagName: '',
  layout,
});
