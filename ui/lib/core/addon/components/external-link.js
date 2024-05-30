/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';

/**
 * @deprecated
 * @module ExternalLink
 * ExternalLink components are used to render anchor links to non-cluster links. Automatically opens in a new tab with noopener noreferrer.
 *
 * @example
 * <ExternalLink @href="https://hashicorp.com">Arbitrary Link</ExternalLink>
 *
 * * Use HDS links with @isHrefExternal={{true}} instead
 * <Hds::Link::Inline @icon="external-link" @isHrefExternal={{true}} @href="https://hashicorp.com">My link</Hds::Link::Inline>
 *
 * @param {string} href - The full href with protocol
 * @param {boolean} [sameTab=false] - by default, these links open in new tab. To override, pass @sameTab={{true}}
 *
 */
export default class ExternalLinkComponent extends Component {
  get href() {
    return this.args.href;
  }
}
