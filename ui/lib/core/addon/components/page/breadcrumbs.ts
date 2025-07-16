/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { assert } from '@ember/debug';

import type { Breadcrumb } from 'vault/vault/app-types';

interface Args {
  breadcrumbs: Array<Breadcrumb>;
}

/**
 * @module Page::Breadcrumbs
 * Page::Breadcrumbs components are used to display an array of breadcrumbs at the top of the page.
 *
 * @example
 * <Page::Breadcrumbs @breadcrumbs={{array (hash label="Home" route="vault") (hash label="my-secret")}}  />
 *
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
