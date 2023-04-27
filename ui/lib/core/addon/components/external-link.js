/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

/**
 * @module ExternalLinkComponent
 * `ExternalLink` components are used to render anchor links to non-cluster links. Automatically opens in a new tab with noopener noreferrer.
 *  To link to developer.hashicorp.com, use DocLink .
 *
 * @example
 * ```js
    <ExternalLink @href="https://hashicorp.com">Arbitrary Link</ExternalLink>
 * ```
 *
 * @param {string} href="https://example.com/" - The full href with protocol
 * @param {boolean} [sameTab=false] - by default, these links open in new tab. To override, pass @sameTab={{true}}
 *
 */
export default class ExternalLinkComponent extends Component {
  get href() {
    return this.args.href;
  }
}
