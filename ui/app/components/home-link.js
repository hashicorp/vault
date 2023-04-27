/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';

/**
 * @module HomeLink
 * `HomeLink` is a span that contains either the text `home` or the `LogoEdition` component.
 *
 * @example
 * ```js
 * <HomeLink @class="navbar-item splash-page-logo">
 *  <LogoEdition />
 * </HomeLink>
 * ```
 * @param {string} class - Classes attached to the the component.
 * @param {string} text - Text displayed instead of logo.
 *
 * @see {@link https://github.com/hashicorp/vault/search?l=Handlebars&q=HomeLink|Uses of HomeLink}
 * @see {@link https://github.com/hashicorp/vault/blob/main/ui/app/components/home-link.js|HomeLink Source Code}
 */

export default class HomeLink extends Component {
  get text() {
    return 'home';
  }

  get computedClasses() {
    return this.classNames.join(' ');
  }
}
