/**
 * @module ToolbarSecretLink
 * `ToolbarSecretLink` styles SecretLink for the Toolbar.
 * It should only be used inside of `Toolbar`.
 *
 * @example
 * ```js
 * <Toolbar>
 *   <ToolbarActions>
 *     <ToolbarDownloadButton @actionText="Download policy" @extension={{if (eq policyType "acl") model.format "sentinel"}} @filename={{model.name}} @data={{model.policy}} />
 *   </ToolbarActions>
 * </Toolbar>
 * ```
 *
 */

import DownloadButton from './download-button';

export default DownloadButton.extend({
  classNames: ['toolbar-link'],
});
