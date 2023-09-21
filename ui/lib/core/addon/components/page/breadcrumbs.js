/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { assert } from '@ember/debug';

/**
 * @module Page::Breadcrumbs
 * Page::Breadcrumbs components are used to display breadcrumb links. This is component will be replaced when HDS system is incorporated
 *
 * @example
 * ```js
 * <Page::Breadcrumbs @breadcrumbs={{this.breadcrumbs}}  />
 * ```
 * @param {array} breadcrumbs - array of objects with a label and route to display as breadcrumbs
 */

export default class Breadcrumbs extends Component {
  constructor() {
    super(...arguments);
    this.args.breadcrumbs.forEach((breadcrumb) => {
      assert('breadcrumb must have a label key', Object.keys(breadcrumb).includes('label'));
    });
  }
}
