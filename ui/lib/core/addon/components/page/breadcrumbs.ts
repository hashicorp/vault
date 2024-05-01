/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { assert } from '@ember/debug';

interface Args {
  breadcrumbs: Array<Breadcrumb>;
}
interface Breadcrumb {
  label: string;
  route?: string; // Do not provide for current route
  icon?: string;
  model?: string;
  models?: string[];
  linkToExternal?: boolean;
}

/**
 * @module Page::Breadcrumbs
 * Page::Breadcrumbs components are used to display an array of breadcrumbs at the top of the page.
 *
 * @example
 * ```js
 * <Page::Breadcrumbs @breadcrumbs={{this.breadcrumbs}}  />
 * ```
 * @param {array} breadcrumbs - array of Breadcrumb objects, must contain a label key. If no route is provided, crumb is assumed to be the current page
 */

export default class Breadcrumbs extends Component<Args> {
  constructor(owner: unknown, args: Args) {
    super(owner, args);
    assert(
      'breadcrumb object must include a label key',
      this.args.breadcrumbs.every((crumb) => Object.keys(crumb).includes('label'))
    );
  }
}
