/**
 * Copyright IBM Corp. 2016, 2025
 * SPDX-License-Identifier: BUSL-1.1
 */

import FlashMessages from 'ember-cli-flash/services/flash-messages';

/**
 * we extend the ember-cli-flash service here so each ember engine can
 * import 'flash-messages' as a dependency giving it access to the
 * <FlashMessage> template in the main app's cluster.hbs file
 * @see https://github.com/adopted-ember-addons/ember-cli-flash
 *
 * To render links, pass the data in the options block. To standardized toast messages,
 * only one action is supported and the iconPosition will always be "trailing"
 * @example
 ```
 this.flashMessages.success('Policy saved successfully', {
  link: {
    text: 'View policy',
    route: 'vault.cluster.policy.show',
    models: ['acl', this.policyName],
    },
 });
 */

export default class FlashMessageService extends FlashMessages {}
