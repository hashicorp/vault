import Component from '@glimmer/component';
import { assert } from '@ember/debug';

/**
 * @module Breadcrumbs
 * Breadcrumbs components are used to display a header with breadcrumb links and an optional title below
 *
 * @example
 * ```js
 * <Breadcrumbs @breadcrumbs={{this.breadcrumbs}}  />
 * ```
 * @param {array} breadcrumbs - array of objects with a label and path to display as breadcrumbs
 * @param {string} [pageTitle] - optional title to display in the header
 * @param {string} [icon] - icon name that displays to the left of the page title
 */

export default class Breadcrumbs extends Component {
  constructor() {
    super(...arguments);
    this.args.breadcrumbs.forEach((breadcrumb) => {
      assert('breadcrumb has a label key', Object.keys(breadcrumb).includes('label'));
    });
  }
}
