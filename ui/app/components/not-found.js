/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { service } from '@ember/service';
import Component from '@glimmer/component';

/**
 * @module NotFound
 * NotFound components are used to show a message that the route was not found.
 *
 * @example
 * ```js
 * <NotFound @model={{this.model}} />
 * ```
 * @param {object} model - routes model passed into the component.
 */

export default class NotFound extends Component {
  @service router;

  constructor() {
    super(...arguments);
    const { currentURL } = this.router;
    if (currentURL.startsWith('/vault/settings/secrets/configure/')) {
      // vault.cluster.settings.secrets.configure was an old route that was removed.
      // Redirect to the new route and pass the id of the secret engine which was located at the end of the old path.
      this.router.transitionTo(
        'vault.cluster.secrets.backend.configuration.edit',
        currentURL.split('/').pop()
      );
    }
  }
}
