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

  get path() {
    return this.router.currentURL;
  }
}
